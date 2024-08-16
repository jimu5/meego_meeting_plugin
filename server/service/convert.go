package service

import (
	"meego_meeting_plugin/model"
)

func MeetingInfos2ModelVCMeeting(meetings []*MeetingInfo, calendarID, calendarEventID string) []*model.VCMeeting {
	result := make([]*model.VCMeeting, 0, len(meetings))
	for _, m := range meetings {
		if m == nil {
			continue
		}
		result = append(result, &model.VCMeeting{
			CalendarID:      calendarID,
			CalendarEventID: calendarEventID,
			MeetingID:       m.MeetingID,
			MeetingData:     (*model.Meeting)(&m.MeetingData),
			RecordInfo:      (*model.RecordInfo)(&m.RecordInfo),
		})
	}
	return result
}
