package lark_api

import (
	larkcalendar "github.com/larksuite/oapi-sdk-go/v3/service/calendar/v4"
	larkvc "github.com/larksuite/oapi-sdk-go/v3/service/vc/v1"
)

type (
	SearchCalendarEventRespData larkcalendar.SearchCalendarEventRespData
	GetCalendarEventRespData    larkcalendar.GetCalendarEventRespData
	UserCalendar                larkcalendar.UserCalendar
	CalendarEvent               larkcalendar.CalendarEvent
	VChat                       larkcalendar.Vchat
)

// 是否为无需视频会议的类型
func (ce *CalendarEvent) HasMeeting() bool {
	if ce.Vchat == nil {
		return false
	}
	return (*VChat)(ce.Vchat).HasMeeting()
}

// 是否无需视频会议
func (v *VChat) HasMeeting() bool {
	if v.VcType != nil {
		if *v.VcType == "no_meeting" {
			return false
		}
	}
	if v.MeetingUrl == nil {
		return false
	}
	if v.MeetingUrl != nil && len(*v.MeetingUrl) == 0 {
		return false
	}

	return true
}

type (
	GetMeetingRespData          larkvc.GetMeetingRespData
	GetMeetingRecordingRespData larkvc.GetMeetingRecordingRespData
	Meeting                     larkvc.Meeting
)
