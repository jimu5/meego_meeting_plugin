package model

import "gorm.io/gorm"

type VCMeetingUnBind struct {
	gorm.Model
	CalendarID      string `json:"calendar_id,omitempty" gorm:"size:256"`
	CalendarEventID string `json:"calendar_event_id,omitempty" gorm:"size:256"`
	WorkItemID      int64  `json:"work_item_id"`
	MeetingID       string `json:"meeting_id,omitempty" gorm:"size:256"`
}
