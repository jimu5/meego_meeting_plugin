package handler

import (
	"sort"
	"strconv"

	"meego_meeting_plugin/common"
	"meego_meeting_plugin/dal"
	"meego_meeting_plugin/model"
	"meego_meeting_plugin/service"
	"meego_meeting_plugin/service/lark_api"
	"meego_meeting_plugin/util"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	larkcalendar "github.com/larksuite/oapi-sdk-go/v3/service/calendar/v4"
	larkvc "github.com/larksuite/oapi-sdk-go/v3/service/vc/v1"
	"github.com/larksuite/project-oapi-sdk-golang/service/user"
	"github.com/samber/lo"
)

type CalendarSearchParam struct {
	QueryWord string `query:"query_word" json:"query_word"`
	lark_api.PageParam
}

// ShowAccount godoc
//
//	@Summary		搜索日程
//	@Description	根据关键字搜索, 返回的结构信息可以参考 https://open.feishu.cn/document/server-docs/calendar-v4/calendar-event/search
//	@Tags			calendar
//	@Produce		json
//	@Param			param	query		CalendarSearchParam	true	"Account ID"
//	@Success		200	{object}	lark_api.SearchCalendarEventRespData
//	@Router			/api/v1/lark/calendar/search	[get]
func CalendarSearch(c *fiber.Ctx) error {
	param := CalendarSearchParam{}
	err := c.QueryParser(&param)
	if err != nil {
		return err
	}
	token := c.Locals(common.LarkUserAccessToken).(string)
	res, err := service.Lark.SearchCalendar(c.Context(), param.QueryWord, token, param.PageParam)
	if err != nil {
		log.Errorf("[CalendarSearch]search calendar failed, err: %v", err)
		return err
	}
	err = c.JSON(res)
	if err != nil {
		return err
	}
	return nil
}

type OperateCalendarEventWithWorkItemParam struct {
	ProjectKey              string `json:"project_key"`
	WorkItemTypeKey         string `json:"work_item_type_key"`
	WorkItemID              int64  `json:"work_item_id"`
	CalendarID              string `json:"calendar_id"`                // 日历 ID
	CalendarEventID         string `json:"calendar_event_id"`          // 日程 ID
	MeetingID               string `json:"meeting_id"`                 // 会议 id (解绑的时候使用)
	WithAfterRecurringEvent bool   `json:"with_after_recurring_event"` // 处理之后的重复日程 (解绑的时候使用)
}

// ShowAccount godoc
//
//	@Summary		操作日程和 meego 实例
//	@Description	将日程和实例绑定或者解绑
//	@Tags			Plugin
//	@Produce		json
//	@Param			OperateCalendarEventWithWorkItemParam	body	OperateCalendarEventWithWorkItemParam	true	"参数"
//	@Success		200
//	@Router			/api/v1/meego/calendar_event/bind	[post]
func BindCalendarEventWithWorkItem(c *fiber.Ctx) error {
	param := OperateCalendarEventWithWorkItemParam{}
	err := c.BodyParser(&param)
	if err != nil {
		log.Errorf("[BindCalendarEventWithWorkItem] body parser error, err: %v", err)
		return err
	}
	if len(param.ProjectKey) == 0 || len(param.WorkItemTypeKey) == 0 || param.WorkItemID == 0 ||
		len(param.CalendarEventID) == 0 {
		return ErrInvalidParam
	}
	token := c.Locals(common.LarkUserAccessToken).(string)
	operator := c.Locals(common.MeegoUserKey).(string)
	err = service.Plugin.BindCalendar(c.Context(), service.BindCalendarParam{
		ProjectKey:      param.ProjectKey,
		WorkItemTypeKey: param.WorkItemTypeKey,
		WorkItemID:      param.WorkItemID,
		CalendarEventID: param.CalendarEventID,
	}, token, operator)
	if err != nil {
		log.Errorf("BindCalendarEventWithWorkItem err, err: %v", err)
		return err
	}
	err = c.JSON(&DefaultResp)
	if err != nil {
		return err
	}
	return nil
}

// ShowAccount godoc
//
//	@Summary		操作日程和 meego 实例解绑
//	@Description	将日程和实例解绑
//	@Tags			Plugin
//	@Produce		json
//	@Param			OperateCalendarEventWithWorkItemParam	body	OperateCalendarEventWithWorkItemParam	true	"参数"
//	@Success		200
//	@Router			/api/v1/meego/calendar_event/unbind	[post]
func UnBindCalendarEventWithWorkItem(c *fiber.Ctx) error {
	param := OperateCalendarEventWithWorkItemParam{}
	err := c.BodyParser(&param)
	if err != nil {
		log.Errorf("[UnBindCalendarEventWithWorkItem] body parser error, err: %v", err)
		return err
	}
	if len(param.ProjectKey) == 0 || len(param.WorkItemTypeKey) == 0 || param.WorkItemID == 0 ||
		len(param.CalendarID) == 0 || len(param.CalendarEventID) == 0 {
		return ErrInvalidParam
	}
	// 1. 查询出原来的记录来
	bind, err := dal.CalendarBind.GetCalendarBindByWorkItemIDAndCalendarEventID(c.Context(), param.WorkItemID, param.CalendarEventID)
	if err != nil {
		return err
	}
	// 2. 修改记录为未绑定状态
	needUnbind := false
	if len(util.GetPointerInfo(bind.CalendarEventData.Recurrence)) == 0 || len(param.MeetingID) == 0 {
		// 非循环的直接解绑
		needUnbind = true
	} else {
		// 循环的需要判断下时间啥的
		meetings, errQ := dal.CalendarBind.MGetCalendarMeetingsByCalendarEventID(c.Context(), bind.CalendarEventID)
		if errQ != nil {
			log.Error(errQ)
			return err
		}
		if len(meetings) == 0 {
			needUnbind = true
		}
		meetingIDMap := lo.SliceToMap(meetings, func(item *model.VCMeeting) (string, *model.VCMeeting) {
			if item != nil {
				return item.MeetingID, item
			}
			return "", nil
		})
		unbindVCMeetings := make([]*model.VCMeeting, 0)
		currentMeeting := meetingIDMap[param.MeetingID]
		if currentMeeting != nil {
			unbindVCMeetings = append(unbindVCMeetings, currentMeeting)
			withAfter := param.WithAfterRecurringEvent
			curStartTime, errP := strconv.ParseInt(util.GetPointerInfo(currentMeeting.MeetingData.StartTime), 10, 64)
			if errP == nil && withAfter {
				for _, m := range meetings {
					if m == nil {
						continue
					}
					mStartTime, errM := strconv.ParseInt(util.GetPointerInfo(m.MeetingData.StartTime), 10, 64)
					if errM != nil {
						continue
					}
					// 晚于这个的开始时间的都需要解绑
					if mStartTime > curStartTime {
						unbindVCMeetings = append(unbindVCMeetings, m)
					}
				}
			}
		}
		if len(unbindVCMeetings) > 0 {
			err = dal.VCMeetingUnBind.SaveUnbindVCMeetings(c.Context(), param.WorkItemID, unbindVCMeetings)
			if err != nil {
				log.Error(err)
				return err
			}
		}
		// FIXME: 写后查询场景
		unbinds, errQ := dal.VCMeetingUnBind.GetVCMeetingUnbindInfoByWorkItemID(c.Context(), param.WorkItemID)
		if errQ != nil {
			log.Error(err)
			return err
		}
		unbindsMap := lo.SliceToMap(unbinds, func(i *model.VCMeetingUnBind) (string, struct{}) {
			return i.MeetingID, struct{}{}
		})
		if len(meetings) == len(unbindsMap) {
			needUnbind = true
		}
	}
	if needUnbind {
		bind.Bind = false
		// 解绑的话把 meeting 绑定信息也清理
		err = dal.VCMeetingUnBind.DeleteMeetingsByWorkItemID(c.Context(), param.WorkItemID)
		if err != nil {
			log.Error(err)
			return err
		}
		err = dal.CalendarBind.CreateOrUpdateCalendarBind(c.Context(), &bind, GetMeegoUserKey(c))
		if err != nil {
			return err
		}
	}

	err = c.JSON(&DefaultResp)
	if err != nil {
		return err
	}
	return nil
}

type ListWorkItemMeetingsParam struct {
	ProjectKey      string `json:"project_key,omitempty"`
	WorkItemTypeKey string `json:"work_item_type_key,omitempty"`
	WorkItemID      int64  `json:"work_item_id,omitempty"`
	PageParam              // 分页
}

// TODO: 还没实现
type MeegoUserInfo struct {
	MeegoUserKey string `json:"meego_user_key"` // Meego 的UserKey
	NameCN       string `json:"name_cn"`        // 中文名
	NameEN       string `json:"name_en"`        // 英文名
	Email        string `json:"email"`          // 邮箱
	AvatarUrl    string `json:"avatar_url"`     // 头像链接
}

type WorkItemMeeting struct {
	ProjectKey                         string                       `json:"project_key,omitempty"`
	WorkItemTypeKey                    string                       `json:"work_item_type_key,omitempty"`
	WorkItemID                         int64                        `json:"work_item_id,omitempty"`
	CalendarID                         string                       `json:"calendar_id,omitempty"`
	CalendarEventID                    string                       `json:"calendar_event_id,omitempty"`
	CalendarEventName                  string                       `json:"calendar_event_name,omitempty"`
	CalendarEventRecurrence            string                       `json:"calendar_event_recurrence"`
	CalendarEventDesc                  string                       `json:"calendar_event_desc,omitempty"`                   // 日程描述(应该是对应的会议描述, 因为会议没有描述)
	CalendarEventAppLink               string                       `json:"calendar_event_app_link"`                         // 日程跳转 APP 链接
	CalendarEventOrganizer             *larkcalendar.EventOrganizer `json:"calendar_event_organizer"`                        // 日程组织者
	MeetingID                          string                       `json:"meeting_id"`                                      // 会议 ID
	MeetingTopic                       string                       `json:"meeting_topic,omitempty"`                         // 会议主题
	MeetingHostUser                    larkvc.MeetingUser           `json:"meeting_host_user"`                               // 会议主持人
	MeetingTime                        MeetingTime                  `json:"meeting_time"`                                    // 会议时间, 都是 unix 时间, 单位 sec
	MeetingRecordURL                   string                       `json:"meeting_record_url,omitempty"`                    // 会议录制链接
	MeetingMinuteURL                   string                       `json:"meeting_minute_url,omitempty"`                    // 会议纪要链接(看着openapi接口好像没有)
	MeetingStatus                      int                          `json:"meeting_status,omitempty"`                        // 会议状态, 可选值 1(呼叫中), 2(进行中) 3(已结束)
	MeetingParticipantCount            string                       `json:"meeting_participant_count,omitempty"`             // 参会峰值人数
	MeetingParticipantCountAccumulated string                       `json:"meeting_participant_count_accumulated,omitempty"` // 参会累计人数
	BindOperator                       string                       `json:"bind_operator,omitempty"`                         // 关联日程操作人(meego 的 userKey)
	BindOperatorInfo                   MeegoUserInfo                `json:"bind_operator_info"`                              // 关联日程操作人的信息
}

func (m *WorkItemMeeting) ApplyModelMeeting(meeting *model.VCMeeting) {
	m.MeetingID = meeting.MeetingID
	m.MeetingTopic = util.GetPointerInfo(meeting.MeetingData.Topic)
	m.MeetingHostUser = util.GetPointerInfo(meeting.MeetingData.HostUser)
	m.MeetingTime = MeetingTime{
		CreateTime: util.GetPointerInfo(meeting.MeetingData.CreateTime),
		StartTime:  util.GetPointerInfo(meeting.MeetingData.StartTime),
		EndTime:    util.GetPointerInfo(meeting.MeetingData.EndTime),
	}
	m.MeetingRecordURL = util.GetPointerInfo(meeting.RecordInfo.Url)
	m.MeetingMinuteURL = util.GetPointerInfo(meeting.RecordInfo.Url)
	m.MeetingStatus = util.GetPointerInfo(meeting.MeetingData.Status)
	m.MeetingParticipantCount = util.GetPointerInfo(meeting.MeetingData.ParticipantCount)
	m.MeetingParticipantCountAccumulated = util.GetPointerInfo(meeting.MeetingData.ParticipantCountAccumulated)
	// 如果日程没有名字, 用会议的名
	if len(m.CalendarEventName) == 0 {
		m.CalendarEventName = util.GetPointerInfo(meeting.MeetingData.Topic)
	}
}

// 都是 unix 时间, 单位 sec
type MeetingTime struct {
	CreateTime string `json:"create_time,omitempty"` // 创建时间
	StartTime  string `json:"start_time,omitempty"`  // 开始时间
	EndTime    string `json:"end_time,omitempty"`    // 结束时间
}

type ListWorkItemMeetingsResp struct {
	Meetings []*WorkItemMeeting `json:"meetings"`
	Total    int64              `json:"total"`
}

// ShowAccount godoc
//
//	@Summary		分页获取实例关联的会议
//	@Description	分页获取实例关联的会议
//	@Tags			Plugin
//	@Produce		json
//	@Param			ListWorkItemMeetingsParam	query	ListWorkItemMeetingsParam	true	"参数"
//	@Success		200	{object}	ListWorkItemMeetingsResp
//	@Router			/api/v1/meego/work_item_meetings	[post]
func ListWorkItemMeetings(c *fiber.Ctx) error {
	param := ListWorkItemMeetingsParam{}
	err := c.BodyParser(&param)
	if err != nil {
		log.Error(err)
		return err
	}
	resp := ListWorkItemMeetingsResp{}
	// 查询出总数量
	binds, err := dal.CalendarBind.GetBindCalendarByProjectKeyAndWorkItemTypeKeyAndWorkItemID(c.Context(), param.ProjectKey, param.WorkItemTypeKey, param.WorkItemID)
	if err != nil {
		return err
	}
	userMap, err := service.Plugin.GetUserInfoByBinds(c.Context(), binds)
	if err != nil {
		log.Errorf("[ListWorkItemMeetings] err, err: %v", err)
	}
	sortBindsByCalendarEventStartTime(binds)
	calendarEventIDMap := make(map[string]*model.CalendarBind, len(binds))
	calenderEvent2Meetings := make(map[string][]*model.VCMeeting)
	for _, bind := range binds {
		calendarEventIDMap[bind.CalendarEventID] = bind
		calenderEvent2Meetings[bind.CalendarEventID] = make([]*model.VCMeeting, 0)
	}
	_, err = dal.CalendarBind.CountMeetingByCalendarEventID(c.Context(), lo.Keys(calendarEventIDMap))
	if err != nil {
		return err
	}
	meetings, err := dal.CalendarBind.MGetMeetingByCalendarEventID(c.Context(), lo.Keys(calendarEventIDMap))
	if err != nil {
		log.Errorf("[GetBindCalendarByProjectKeyAndWorkItemTypeKeyAndWorkItemID] err, err: %v", err)
		return err
	}
	for _, meeting := range meetings {
		if meeting != nil {
			calenderEvent2Meetings[meeting.CalendarEventID] = append(calenderEvent2Meetings[meeting.CalendarEventID], meeting)
		}
	}
	unBindMeetings, err := dal.VCMeetingUnBind.GetVCMeetingUnbindInfoByWorkItemID(c.Context(), param.WorkItemID)
	if err != nil {
		log.Error(err)
		return err
	}
	unBindMeetingIDs := make(map[string]struct{}, len(unBindMeetings))
	for _, m := range unBindMeetings {
		if m != nil {
			unBindMeetingIDs[m.MeetingID] = struct{}{}
		}
	}
	resp.Meetings = make([]*WorkItemMeeting, 0, len(meetings))
	for _, bind := range binds {
		if bind == nil {
			continue
		}
		userInfo := userMap[bind.UpdateBy]
		vcMeetings := calenderEvent2Meetings[bind.CalendarEventID]
		bindMeetings := make([]*WorkItemMeeting, 0, len(vcMeetings))
		if len(vcMeetings) == 0 {
			workItemMeeting := genWorkItemMeetingByBind(bind, userInfo)
			bindMeetings = append(bindMeetings, &workItemMeeting)
		}
		for _, meeting := range vcMeetings {
			if _, ok := unBindMeetingIDs[meeting.MeetingID]; ok {
				continue
			}
			if meeting.MeetingData == nil {
				meeting.MeetingData = &model.Meeting{}
			}
			if meeting.RecordInfo == nil {
				meeting.RecordInfo = &model.RecordInfo{}
			}
			workItemMeeting := genWorkItemMeetingByBind(bind, userInfo)
			workItemMeeting.ApplyModelMeeting(meeting)
			bindMeetings = append(bindMeetings, &workItemMeeting)
		}
		sortWorkItemMeetingsByMeetingStartTime(bindMeetings)
		resp.Meetings = append(resp.Meetings, bindMeetings...)
	}
	resp.Total = int64(len(resp.Meetings))
	err = c.JSON(&resp)
	if err != nil {
		return err
	}
	return nil
}

func sortBindsByCalendarEventStartTime(binds []*model.CalendarBind) {
	sort.Slice(binds, func(i, j int) bool {
		iStartTime := binds[i].CalendarEventStartTime
		jStartTime := binds[j].CalendarEventStartTime
		if iStartTime == nil {
			iStartTime = binds[i].CalendarEventData.GetStartTime()
		}
		if jStartTime == nil {
			jStartTime = binds[j].CalendarEventData.GetStartTime()
		}
		return util.GetPointerInfo(iStartTime).Unix() > util.GetPointerInfo(jStartTime).Unix()
	})
}

func genWorkItemMeetingByBind(bind *model.CalendarBind, userInfo *user.UserBasicInfo) WorkItemMeeting {
	if bind == nil {
		return WorkItemMeeting{}
	}
	calendarEventData := &model.CalendarEventData{}
	if bind.CalendarEventData != nil {
		calendarEventData = bind.CalendarEventData
	}
	return WorkItemMeeting{
		ProjectKey:              bind.ProjectKey,
		WorkItemTypeKey:         bind.WorkItemTypeKey,
		WorkItemID:              bind.WorkItemID,
		CalendarID:              bind.CalendarID,
		CalendarEventID:         bind.CalendarEventID,
		BindOperator:            bind.UpdateBy,
		CalendarEventName:       util.GetPointerInfo(calendarEventData.Summary),
		CalendarEventDesc:       util.GetPointerInfo(calendarEventData.Description),
		CalendarEventAppLink:    util.GetPointerInfo(calendarEventData.AppLink),
		CalendarEventOrganizer:  calendarEventData.EventOrganizer,
		CalendarEventRecurrence: util.GetPointerInfo(calendarEventData.Recurrence),
		BindOperatorInfo:        genUserInfo(userInfo),
	}
}

func genUserInfo(u *user.UserBasicInfo) MeegoUserInfo {
	if u == nil {
		return MeegoUserInfo{}
	}
	return MeegoUserInfo{
		MeegoUserKey: u.UserKey,
		NameCN:       u.NameCn,
		NameEN:       u.NameEn,
		Email:        u.Email,
		AvatarUrl:    u.AvatarUrl,
	}
}

func sortWorkItemMeetingsByMeetingStartTime(meetings []*WorkItemMeeting) {
	sort.Slice(meetings, func(i, j int) bool {
		iStartTime, _ := strconv.ParseInt(meetings[i].MeetingTime.StartTime, 10, 64)
		jStartTime, _ := strconv.ParseInt(meetings[j].MeetingTime.StartTime, 10, 64)

		return iStartTime > jStartTime
	})
}
