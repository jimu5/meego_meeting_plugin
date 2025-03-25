package service

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"meego_meeting_plugin/service/lark_api"
	"meego_meeting_plugin/util"

	"github.com/gofiber/fiber/v2/log"
	larkcalendar "github.com/larksuite/oapi-sdk-go/v3/service/calendar/v4"
	larkvc "github.com/larksuite/oapi-sdk-go/v3/service/vc/v1"
)

var Lark LarkService

// accessToken 应该从 ctx 中拿到, 不从参数里面取
type LarkService struct {
	LarkAPI lark_api.LarkAPI
}

func NewLarkService() LarkService {
	return LarkService{
		LarkAPI: lark_api.API,
	}
}

type CalendarMeetingInfo struct {
	CalendarID      string                     `json:"calendar_id,omitempty"`
	CalendarEventID string                     `json:"calendar_event_id,omitempty"`
	EventData       larkcalendar.CalendarEvent `json:"event_data"`
	HasVChat        bool                       `json:"has_v_chat,omitempty"`
	MeetingNo       string                     `json:"meeting_no,omitempty"`

	Meetings []*lark_api.Meeting `json:"meetings,omitempty"`
}

func (c CalendarMeetingInfo) GetEventStartTime() *time.Time {
	if c.EventData.StartTime == nil || c.EventData.StartTime.Timestamp == nil {
		return nil
	}
	timeStamp, err := strconv.ParseInt(*c.EventData.StartTime.Timestamp, 10, 64)
	if err != nil {
		return nil
	}
	t := time.Unix(timeStamp, 0)
	return &t
}

// 搜索日程
func (l LarkService) SearchCalendar(ctx context.Context, queryWord, userAccessToken string,
	pageParam lark_api.PageParam) (*lark_api.SearchCalendarEventRespData, error) {
	calendars, err := lark_api.API.CalendarAPI.GetPrimaryCalendars(ctx, userAccessToken)
	if err != nil {
		return nil, err
	}
	// 获取属于自己的主日程id
	if len(calendars) == 0 {
		return nil, ErrPrimaryCalendar
	}
	calendarID := l.getFirstPrimaryCalendarID(calendars)
	if len(calendarID) == 0 {
		return nil, ErrPrimaryCalendar
	}
	// 搜索日程
	searchRes, err := lark_api.API.CalendarAPI.SearchCalendarEvents(ctx, calendarID, queryWord, userAccessToken, pageParam)
	if err != nil {
		log.Errorf("[LarkService] SearchCalendar err: %v", err)
		return nil, err
	}
	return searchRes, nil
}

func (l LarkService) SearchCalendarByTimeAndChatIDs(ctx context.Context, queryWord string,
	startTimeStamp, endTimeStamp string, chatIDs []string, userAccessToken string) (
	*lark_api.SearchCalendarEventRespData, error) {
	calendars, err := lark_api.API.CalendarAPI.GetPrimaryCalendars(ctx, userAccessToken)
	if err != nil {
		return nil, err
	}
	// 获取属于自己的主日程id
	if len(calendars) == 0 {
		return nil, ErrPrimaryCalendar
	}
	calendarID := l.getFirstPrimaryCalendarID(calendars)
	if len(calendarID) == 0 {
		return nil, ErrPrimaryCalendar
	}
	searchRes, err := lark_api.API.CalendarAPI.SearchCalendarEventsByTimeAndChatIDs(ctx, calendarID, queryWord, startTimeStamp, endTimeStamp,
		chatIDs, userAccessToken)
	if err != nil {
		log.Errorf("[LarkService] SearchCalendar err: %v", err)
		return nil, err
	}
	return searchRes, nil
}

// 获取会议的信息
func (l LarkService) GetMeetingRecordInfoByCalendar(ctx context.Context, calendarID, calendarEventID, userAccessToken string) (CalendarMeetingInfo, error) {
	res := CalendarMeetingInfo{
		CalendarID:      calendarID,
		CalendarEventID: calendarEventID,
	}
	// 1. 获取第一个主日程 ID
	if len(calendarID) == 0 {
		calendars, err := lark_api.API.CalendarAPI.GetPrimaryCalendars(ctx, userAccessToken)
		if err != nil {
			log.Errorf("[LarkService] GetMeetingRecordInfoByCalendar GetPrimaryCalendars, err: %v", err)
			return res, err
		}
		calendarID = l.getFirstPrimaryCalendarID(calendars)
		if len(calendarID) == 0 {
			return res, ErrPrimaryCalendar
		}
	}
	res.CalendarID = calendarID
	// 2. 获取日程详情
	calendarEvent, err := lark_api.API.CalendarAPI.GetCalendarEventDetail(ctx, calendarID, calendarEventID, userAccessToken)
	if err != nil {
		log.Errorf("[LarkService] GetMeetingRecordInfoByCalendar GetCalendarEventDetail, err: %v", err)
		return res, err
	}
	if calendarEvent == nil || calendarEvent.Event == nil {
		log.Errorf("[LarkService] GetMeetingRecordInfoByCalendar GetMeetingRecordInfo, calendarEvent is nil")
		return res, ErrNilOpenApiResponse
	}
	res.EventData = *calendarEvent.Event
	// 3. 判断日程是否存在视频会议
	if calendarEvent.Event.Vchat == nil {
		res.HasVChat = false
	} else {
		res.HasVChat = (*lark_api.CalendarEvent)(calendarEvent.Event).HasMeeting()
	}
	if !res.HasVChat {
		return res, nil
	}
	meetingNo := getMeetingNOByMeetingUrl(*calendarEvent.Event.Vchat.MeetingUrl)
	if len(meetingNo) == 0 {
		log.Infof("[LarkService] GetMeetingRecordInfoByCalendar, parsed meetingNo is empty")
		return res, nil
	}
	res.MeetingNo = meetingNo

	// 4. 根据会议号获取关联的会议
	meetingStartTime, meetingEndTime, err := l.getMeetingUnixTimeByCalendarEvent((*lark_api.CalendarEvent)(calendarEvent.Event))
	if err != nil {
		log.Errorf("[LarkService] GetMeetingRecordInfoByCalendar, err: %v", err)
		return res, nil
	}

	meetings, err := lark_api.API.VChatAPI.GetMeetingsListByNo(ctx, meetingNo, meetingStartTime, meetingEndTime, userAccessToken, nil)
	if err != nil {
		log.Errorf("[LarkService] GetMeetingsListByNo, err: %v", err)
		return res, nil
	}
	res.Meetings = meetings
	return res, nil
}

type MeetingInfo struct {
	MeetingID   string
	MeetingData larkvc.Meeting
	RecordInfo  larkvc.MeetingRecording
}

func (l LarkService) GetMeetingInfo(ctx context.Context, meetingID string, userToken string) (MeetingInfo, error) {
	res := MeetingInfo{}
	res.MeetingID = meetingID
	meetingInfo, err := lark_api.API.VChatAPI.GetMeeting(ctx, meetingID, userToken)
	if err != nil {
		log.Errorf("[LarkService] GetMeetingInfo err, meetingID: %s, err: %v", meetingID, err)
		return res, err
	}
	if meetingInfo != nil && meetingInfo.Meeting != nil {
		res.MeetingData = *meetingInfo.Meeting
	}
	recordInfo, err := lark_api.API.VChatAPI.GetMeetingRecord(ctx, meetingID, userToken)
	if err != nil {
		log.Errorf("[LarkService] Get Meeting Record err, meetingID: %s, err: %v", meetingID, err)
	}
	if recordInfo != nil && recordInfo.Recording != nil {
		res.RecordInfo = *recordInfo.Recording
	}
	return res, nil
}

func (l LarkService) MGetMeetingInfo(ctx context.Context, meetingIDs []string, userAccessToken string) ([]*MeetingInfo, error) {
	var wg sync.WaitGroup
	meetings := make([]*MeetingInfo, 0, len(meetingIDs))
	for index := range meetingIDs {
		wg.Add(1)
		meetingID := meetingIDs[index]
		go func() {
			defer wg.Done()
			meeting, err := l.GetMeetingInfo(ctx, meetingID, userAccessToken)
			if err != nil {
				log.Errorf("[LarkService] MGetMeetingInfo err, err: %v", err)
			} else {
				meetings = append(meetings, &meeting)
			}
		}()
	}
	wg.Wait()
	return meetings, nil
}

type UserTokenInfo struct {
	AccessToken        string
	AccessTokenExpire  int64 // unix 时间
	RefreshToken       string
	RefreshTokenExpire int64 // unix 时间
}

func (l LarkService) GetUserAccessToken(ctx context.Context, userAuthCode string) (*UserTokenInfo, error) {
	if len(userAuthCode) == 0 {
		return nil, ErrToken
	}
	// 1. 获取 app token
	appToken, err := lark_api.API.AuthenAPI.GetAppAccessToken(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	// 2. 根据 userAuthCode 获取 accessToken
	userTokenInfo := UserTokenInfo{}
	if len(userAuthCode) != 0 {
		res, errG := lark_api.API.AuthenAPI.GetUserAccessToken(ctx, appToken, userAuthCode)
		if errG != nil {
			log.Error(errG)
			return nil, err
		}
		userTokenInfo.AccessToken = util.GetPointerInfo(res.AccessToken)
		userTokenInfo.RefreshToken = util.GetPointerInfo(res.RefreshToken)
		userTokenInfo.AccessTokenExpire = time.Now().Unix() + int64(util.GetPointerInfo(res.ExpiresIn))
		userTokenInfo.RefreshTokenExpire = time.Now().Unix() + int64(util.GetPointerInfo(res.RefreshExpiresIn))
	}
	return &userTokenInfo, nil
}

func (l LarkService) RefreshUserAccessToken(ctx context.Context, userRefreshToken string) (*UserTokenInfo, error) {
	if len(userRefreshToken) == 0 {
		return nil, ErrToken
	}
	// 1. 获取 app token
	appToken, err := lark_api.API.AuthenAPI.GetAppAccessToken(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	// 2. 刷新 token
	userTokenInfo := UserTokenInfo{}
	if len(userRefreshToken) != 0 {
		res, errG := lark_api.API.AuthenAPI.RefreshUserAccessToken(ctx, appToken, userRefreshToken)
		if errG != nil {
			log.Error(errG)
			return nil, err
		}
		userTokenInfo.AccessToken = util.GetPointerInfo(res.AccessToken)
		userTokenInfo.RefreshToken = util.GetPointerInfo(res.RefreshToken)
		userTokenInfo.AccessTokenExpire = time.Now().Unix() + int64(util.GetPointerInfo(res.ExpiresIn))
		userTokenInfo.RefreshTokenExpire = time.Now().Unix() + int64(util.GetPointerInfo(res.RefreshExpiresIn))
	}
	return &userTokenInfo, nil
}

func (l LarkService) getFirstPrimaryCalendarID(calendars []*lark_api.UserCalendar) string {
	var calendarID string
	for _, calendar := range calendars {
		if calendar == nil || calendar.Calendar == nil || calendar.Calendar.Type == nil || calendar.Calendar.CalendarId == nil {
			continue
		}
		if *calendar.Calendar.Type == CalendarTypePrimary {
			calendarID = *calendar.Calendar.CalendarId
		}
	}
	return calendarID
}

func (l LarkService) getMeetingUnixTimeByCalendarEvent(calendarEvent *lark_api.CalendarEvent) (string, string, error) {
	if calendarEvent == nil {
		return "", "", ErrNilCalendarTime
	}
	var (
		startTime string
		endTime   string
		err       error
	)
	calendarStartTime, err := strconv.ParseInt(*calendarEvent.StartTime.Timestamp, 10, 64)
	if err != nil {
		return "", "", err
	}
	calendarEndTime, err := strconv.ParseInt(*calendarEvent.EndTime.Timestamp, 10, 64)
	if err != nil {
		return "", "", err
	}

	weekDay := time.Hour * 24 * 7
	startTime = fmt.Sprintf("%d", time.Unix(calendarStartTime, 0).Add(-weekDay).Unix())
	endTime = fmt.Sprintf("%d", time.Unix(calendarEndTime, 0).Add(weekDay).Unix())
	return startTime, endTime, nil
}
