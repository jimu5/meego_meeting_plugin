package handler

import (
	"context"
	"time"

	"meego_meeting_plugin/dal"
	"meego_meeting_plugin/service"

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
	// 获取发消息的人信息
	var larkUserInfo *larkim.UserId
	if eventBody.Sender != nil && eventBody.Sender.SenderId != nil {
		larkUserInfo = eventBody.Sender.SenderId
	}
	param := service.HandleMeetingBindByUserKeyParam{
		Content:      content,
		LarkUserInfo: larkUserInfo,
		Record:       record,
	}
	err = service.Plugin.HandleMeetingBindByUserKey(ctx, param)
	if err != nil {
		log.Errorf("[handleChatCalendarMessage] handleMeetingBindByUserKey err, err: %v", err)
		return nil
	}
	// 目前没法使用应用身份进行日程绑定，应用身份获取用户主日历上的日程时，会默认没有权限
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
	userInfo, err := service.Plugin.GetUserInfoByMeegoUserKey(ctx, bind.UpdateBy, true)
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
	userInfo, err := service.Plugin.GetUserInfoByMeegoUserKey(ctx, bind.UpdateBy, true)
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
