package model

import "gorm.io/gorm"

type VCMeetingUnBind struct {
	gorm.Model
	CalendarID      string `json:"calendar_id,omitempty"`
	CalendarEventID string `json:"calendar_event_id,omitempty"`
	WorkItemID      int64  `json:"work_item_id"`
	MeetingID       string `json:"meeting_id,omitempty"`
}
