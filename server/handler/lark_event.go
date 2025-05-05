package handler

import (
	"context"
	"encoding/json"
	"time"

	"meego_meeting_plugin/common"
	"meego_meeting_plugin/dal"
	"meego_meeting_plugin/service"
	"meego_meeting_plugin/service/lark_api"

	"github.com/avast/retry-go"
	"github.com/gofiber/fiber/v2/log"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	larkvc "github.com/larksuite/oapi-sdk-go/v3/service/vc/v1"
)

var LarkEventHandler = dispatcher.NewEventDispatcher("", "").
	OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
		if event == nil {
			return nil
		}
		return handleChatCalendarMessage(ctx, event.Event)
	}).
	OnP2MeetingAllMeetingStartedV1(func(ctx context.Context, event *larkvc.P2MeetingAllMeetingStartedV1) error {
		if event == nil {
			return nil
		}
		return handleAllMeetingStarted(ctx, event.Event)
	}).
	OnP2MeetingAllMeetingEndedV1(func(ctx context.Context, event *larkvc.P2MeetingAllMeetingEndedV1) error {
		if event == nil {
			return nil
		}
		return handleAllMeetingEnd(ctx, event.Event)
	})

type calendarContent struct {
	Summary   string `json:"summary,omitempty"`
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
}

func handleChatCalendarMessage(ctx context.Context, eventBody *larkim.P2MessageReceiveV1Data) error {
	if eventBody == nil || eventBody.Message == nil || eventBody.Message.MessageType == nil {
		return nil
	}
	// 判断是否为日程卡片消息
	messageType := *eventBody.Message.MessageType
	if messageType != "share_calendar_event" && messageType != "calendar" && messageType != "general_calendar" {
		return nil
	}
	if eventBody.Message.ChatId == nil {
		return nil
	}
	if eventBody.Message.Content == nil {
		return nil
	}
	chatID := *eventBody.Message.ChatId
	content := *eventBody.Message.Content
	// 查询这个群是否被工作项绑定
	record, err := dal.JoinChatRecord.FirstByChatID(ctx, chatID)
	if record == nil || err != nil {
		log.Warnf("[handleChatCalendarMessage] chat id bind record not found, chatID: %s", chatID)
		return nil
	}
	if !record.Enable {
		// 没有启用
		return nil
	}
	// 使用 record 的 userKey 和开始结束时间以及关键字来搜索日程
	userInfo, err := service.Plugin.GetUserInfoByMeegoUserKey(ctx, record.Operator)
	if err != nil {
		log.Errorf("[handleChatCalendarMessage] GetUserInfoByMeegoUserKey err, userKey: %s, err: %v", userInfo, err)
		return err
	}
	calendarContentInfo := calendarContent{}
	err = json.Unmarshal([]byte(content), &calendarContentInfo)
	if err != nil {
		log.Errorf("[handleChatCalendarMessage] Unmarshal err, content: %s, err: %v", content, err)
		return err
	}
	// 需要把毫秒级时间戳处理成秒级, 两边的 range 再扩充一秒吧
	startTime, err := common.MillisecondToSecond(calendarContentInfo.StartTime)
	if err != nil {
		log.Errorf("[handleChatCalendarMessage] MillisecondToSecond startTime err, err: %v, sourceVal: %s", err, calendarContentInfo.StartTime)
		return err
	}
	endTime, err := common.MillisecondToSecond(calendarContentInfo.EndTime)
	if err != nil {
		log.Errorf("[handleChatCalendarMessage] MillisecondToSecond EndTime err, err: %v, sourceVal: %s", err, calendarContentInfo.EndTime)
		return err
	}
	var calendars *lark_api.SearchCalendarEventRespData
	// 这里可能有延迟, 等待 5 秒重试, 最多重试3次
	err = retry.Do(func() error {
		calendars, err = service.Lark.SearchCalendarByTimeAndChatIDs(ctx, calendarContentInfo.Summary, startTime, endTime, []string{chatID}, userInfo.LarkUserAccessToken)
		if err != nil || calendars == nil {
			log.Errorf("[handleChatCalendarMessage] err handle, err: %v, calendars: %v", err, calendars)
			return err
		}
		if len(calendars.Items) == 0 {
			log.Warnf("[handleChatCalendarMessage] not search any calendars ")
			return ErrEmptyCalendarSearch
		}
		return nil
	}, retry.Delay(time.Second*5), retry.Attempts(3), retry.DelayType(retry.FixedDelay))
	if err != nil && len(calendars.Items) == 0 {
		log.Errorf("[handleChatCalendarMessage] err get event, err: %v", err)
		return err
	}
	// 目前支取用第一个
	event := calendars.Items[0]
	err = service.Plugin.BindCalendar(ctx, service.BindCalendarParam{
		ProjectKey:      record.ProjectKey,
		WorkItemTypeKey: record.WorkItemTypeKey,
		WorkItemID:      record.WorkItemID,
		CalendarEventID: getPointerInfo(event.EventId),
	}, userInfo.LarkUserAccessToken, record.Operator)
	if err != nil {
		log.Errorf("[handleChatCalendarMessage] bind err, err: %v", err)
		return err
	}
	log.Infof("[handleChatCalendarMessage] source event info: %s", larkcore.Prettify(eventBody))
	return nil
}

func handleAllMeetingStarted(ctx context.Context, eventBody *larkvc.P2MeetingAllMeetingStartedV1Data) error {
	if eventBody == nil || eventBody.Meeting == nil {
		return nil
	}
	// 1. 这个会议是否有关联的日程
	if eventBody.Meeting.CalendarEventId == nil || len(*eventBody.Meeting.CalendarEventId) == 0 {
		return nil
	}
	calendarEventID := *eventBody.Meeting.CalendarEventId
	// 2. 搜索这个日程是否在绑定的日程中
	bind, err := dal.CalendarBind.GetBindByCalendarEventID(ctx, calendarEventID)
	if err != nil {
		log.Error(err)
		return nil
	}
	// 3. 获取绑定的用户 token
	userInfo, err := service.Plugin.GetUserInfoByMeegoUserKey(ctx, bind.UpdateBy)
	if err != nil {
		log.Error(err)
		return err
	}
	// 4. 将新的 meeing 绑定到这个 bind 上
	meeting, err := service.Lark.GetMeetingInfo(ctx, *eventBody.Meeting.Id, userInfo.LarkUserAccessToken)
	if err != nil {
		log.Error(err)
		return err
	}
	meetings := service.MeetingInfos2ModelVCMeeting([]*service.MeetingInfo{&meeting}, bind.CalendarID, bind.CalendarEventID)
	err = dal.CalendarBind.CreateOrUpdateCalendarMeetings(ctx, meetings, bind.UpdateBy)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func handleAllMeetingEnd(ctx context.Context, eventBody *larkvc.P2MeetingAllMeetingEndedV1Data) error {
	if eventBody == nil || eventBody.Meeting == nil {
		return nil
	}
	// 1. 这个会议是否有关联的日程
	if eventBody.Meeting.CalendarEventId == nil || len(*eventBody.Meeting.CalendarEventId) == 0 {
		return nil
	}
	calendarEventID := *eventBody.Meeting.CalendarEventId
	// 2. 搜索这个日程是否在绑定的日程中
	bind, err := dal.CalendarBind.GetBindByCalendarEventID(ctx, calendarEventID)
	if err != nil {
		log.Error(err)
		return nil
	}
	// 3. 获取绑定的用户 token
	userInfo, err := service.Plugin.GetUserInfoByMeegoUserKey(ctx, bind.UpdateBy)
	if err != nil {
		log.Error(err)
		return err
	}
	// 4. 将新的 meeing 更新到这个 bind 上
	meeting, err := service.Lark.GetMeetingInfo(ctx, *eventBody.Meeting.Id, userInfo.LarkUserAccessToken)
	if err != nil {
		log.Error(err)
		return err
	}
	meetings := service.MeetingInfos2ModelVCMeeting([]*service.MeetingInfo{&meeting}, bind.CalendarID, bind.CalendarEventID)
	err = dal.CalendarBind.CreateOrUpdateCalendarMeetings(ctx, meetings, bind.UpdateBy)
	if err != nil {
		log.Error(err)
		return err
	}
	go func() {
		time.Sleep(5 * time.Second)
		log.Infof("[handleAllMeetingEnd] start RetryRefreshMeetingRecordTask")
		service.Plugin.RetryRefreshMeetingRecordTask(ctx, []string{meeting.MeetingID}, userInfo.LarkUserAccessToken)
	}()
	return nil
}
