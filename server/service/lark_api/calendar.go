package lark_api

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2/log"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkcalendar "github.com/larksuite/oapi-sdk-go/v3/service/calendar/v4"
)

// 日历相关的 openapi
type CalendarAPI struct {
	client *LarkClient
}

func NewCalendarAPI(client *LarkClient) CalendarAPI {
	return CalendarAPI{client: client}
}

type PageParam struct {
	PageToken string `query:"page_token" json:"page_token"`
	PageSize  int    `query:"page_size" json:"page_size"`
}

// 获取主日历 ID
func (c CalendarAPI) GetPrimaryCalendars(ctx context.Context, userAccessToken string) ([]*UserCalendar, error) {
	// 发起请求
	resp, err := c.client.Calendar.Calendar.Primary(ctx,
		larkcalendar.NewPrimaryCalendarReqBuilder().Build(),
		larkcore.WithUserAccessToken(userAccessToken))

	if err != nil {
		log.Errorf("[CalendarAPI] GetPrimaryCalendars, err: %v", err)
		return nil, err
	}

	if !resp.Success() {
		log.Errorf("[CalendarAPI] GetPrimaryCalendars resp not success, code: %v, msg: %v, LOGID: %s", resp.Code, resp.Msg, resp.RequestId())
		return nil, resp.CodeError
	}
	if resp.Data == nil {
		return nil, NewErrResponseNotSuccess(resp.Code, resp.Msg)
	}
	result := make([]*UserCalendar, 0, len(resp.Data.Calendars))
	for _, c := range resp.Data.Calendars {
		if c.Calendar != nil {
			if c.Calendar.CalendarId == nil {
				continue
			}
			result = append(result, (*UserCalendar)(c))
		}
	}
	return result, nil
}

// 批量获取主日程 id, 用于使用应用身份来获取
// userIDType 可选值: open_id, union_id, user_id
func (c CalendarAPI) GetPrimaryCalendarsByLarkUserID(ctx context.Context, userIDType string, larkUserIDs []string) ([]*UserCalendar, error) {
	// 发起请求
	req := larkcalendar.NewPrimarysCalendarReqBuilder().
		UserIdType(userIDType).
		Body(larkcalendar.NewPrimarysCalendarReqBodyBuilder().
			UserIds(larkUserIDs).
			Build()).
		Build()

	resp, err := c.client.Calendar.V4.Calendar.Primarys(ctx, req)
	if err != nil {
		log.Errorf("[CalendarAPI] GetPrimaryCalendarIDs, err: %v", err)
		return nil, err
	}

	if !resp.Success() {
		log.Errorf("[CalendarAPI] GetPrimaryCalendarIDs resp not success, code: %v, msg: %v, LOGID: %s", resp.Code, resp.Msg, resp.RequestId())
		return nil, resp.CodeError
	}
	if resp.Data == nil || resp.Data.Calendars == nil {
		return nil, NewErrResponseNotSuccess(resp.Code, resp.Msg)
	}
	result := make([]*UserCalendar, 0, len(resp.Data.Calendars))
	for _, c := range resp.Data.Calendars {
		if c.Calendar != nil {
			if c.Calendar.CalendarId == nil {
				continue
			}
			result = append(result, (*UserCalendar)(c))
		}
	}
	return result, nil
}

// 搜索日程
func (c CalendarAPI) SearchCalendarEvents(ctx context.Context, calendarID, queryWord, userAccessToken string,
	pageParam PageParam) (*SearchCalendarEventRespData, error) {
	log.Infof("[SearchCalendarEvents] calendar ID: %s, queryWord: %s", calendarID, queryWord)
	req := larkcalendar.NewSearchCalendarEventReqBuilder().
		CalendarId(calendarID).
		PageSize(pageParam.PageSize).PageToken(pageParam.PageToken).
		Body(larkcalendar.NewSearchCalendarEventReqBodyBuilder().
			Query(queryWord).
			Filter(&larkcalendar.EventSearchFilter{
				StartTime: &larkcalendar.TimeInfo{
					Timestamp: GetPtr(fmt.Sprintf("%d", time.Now().Add(-90*24*time.Hour).Unix())),
				},
				EndTime: &larkcalendar.TimeInfo{
					Timestamp: GetPtr(fmt.Sprintf("%d", time.Now().Add(90*24*time.Hour).Unix())),
				},
			}).
			Build()).
		Build()

	resp, err := c.client.Calendar.CalendarEvent.Search(ctx, req, larkcore.WithUserAccessToken(userAccessToken))
	if err != nil {
		log.Errorf("[CalendarAPI] SearchCalendarEvents, err: %v", err)
		return nil, err
	}

	if !resp.Success() {
		log.Errorf("[CalendarAPI] SearchCalendarEvents resp not success, code: %v, msg: %v, LOGID: %s", resp.Code, resp.Msg, resp.RequestId())
		return nil, resp.CodeError
	}

	return (*SearchCalendarEventRespData)(resp.Data), nil
}

// 通过时间和群搜索
func (c CalendarAPI) SearchCalendarEventsByTimeAndChatIDs(ctx context.Context, calendarID, queryWord string,
	startTimeStamp, endTimeStamp string, userAccessToken string) (*SearchCalendarEventRespData, error) {
	log.Infof("[SearchCalendarEventsByTimeAndChatIDs] calendar ID: %s, queryWord: %s", calendarID, queryWord)
	req := larkcalendar.NewSearchCalendarEventReqBuilder().
		CalendarId(calendarID).
		Body(larkcalendar.NewSearchCalendarEventReqBodyBuilder().
			Query(queryWord).
			Filter(larkcalendar.NewEventSearchFilterBuilder().
				StartTime(larkcalendar.NewTimeInfoBuilder().
					Timestamp(startTimeStamp).
					Build()).
				EndTime(larkcalendar.NewTimeInfoBuilder().
					Timestamp(endTimeStamp).
					Build()).
				Build()).
			Build()).
		Build()

	var (
		resp *larkcalendar.SearchCalendarEventResp
		err  error
	)
	if len(userAccessToken) != 0 {
		resp, err = c.client.Calendar.CalendarEvent.Search(ctx, req, larkcore.WithUserAccessToken(userAccessToken))
	} else {
		resp, err = c.client.Calendar.CalendarEvent.Search(ctx, req)
	}
	if err != nil {
		log.Errorf("[CalendarAPI] SearchCalendarEventsByTimeAndChatIDs, err: %v", err)
		return nil, err
	}

	if !resp.Success() {
		log.Errorf("[CalendarAPI] SearchCalendarEventsByTimeAndChatIDs resp not success, code: %v, msg: %v, LOGID: %s", resp.Code, resp.Msg, resp.RequestId())
		return nil, resp.CodeError
	}

	return (*SearchCalendarEventRespData)(resp.Data), nil
}

// 获取日程详情
func (c CalendarAPI) GetCalendarEventDetail(ctx context.Context, calendarID, eventID, userAccessToken string) (*GetCalendarEventRespData, error) {
	req := larkcalendar.NewGetCalendarEventReqBuilder().
		CalendarId(calendarID).
		EventId(eventID).
		Build()
	resp, err := c.client.Calendar.CalendarEvent.Get(ctx, req, larkcore.WithUserAccessToken(userAccessToken))
	if err != nil {
		log.Errorf("[CalenderAPI] GetCalendarEventDetail, err: %v", err)
		return nil, err
	}
	if !resp.Success() {
		log.Errorf("[CalendarAPI] GetCalendarEventDetail resp not success, code: %v, msg: %v, LOGID: %s", resp.Code, resp.Msg, resp.RequestId())
		return nil, NewErrResponseNotSuccess(resp.Code, resp.Msg)
	}

	return (*GetCalendarEventRespData)(resp.Data), nil
}

// 订阅日程变更事件
func (c CalendarAPI) SubscriptionCalendarChangeEvent(ctx context.Context, calendarID, calendarEventID,
	userAccessToken string) (*larkcalendar.GetCalendarEventRespData, error) {
	req := larkcalendar.NewGetCalendarEventReqBuilder().
		CalendarId(calendarID).
		EventId(calendarEventID).
		Build()
	resp, err := c.client.Calendar.CalendarEvent.Get(ctx, req, larkcore.WithUserAccessToken(userAccessToken))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if !resp.Success() {
		log.Error(resp.Code, resp.Msg, resp.RequestId())
		return nil, NewErrResponseNotSuccess(resp.Code, resp.Msg)
	}
	return resp.Data, nil
}
