package model

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/log"
	larkcalendar "github.com/larksuite/oapi-sdk-go/v3/service/calendar/v4"
	larkvc "github.com/larksuite/oapi-sdk-go/v3/service/vc/v1"
)

type (
	CalendarEventData larkcalendar.CalendarEvent
)

func (c *CalendarEventData) Scan(input any) error {
	if input == nil {
		return nil
	}
	var err error
	if data, ok := input.([]byte); ok {
		err = json.Unmarshal(data, c)
	}
	if data, ok := input.(string); ok {
		err = json.Unmarshal([]byte(data), c)
	}
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (c *CalendarEventData) Value() (driver.Value, error) {
	ms, err := json.Marshal(c)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return string(ms), nil
}

func (c *CalendarEventData) GetStartTime() *time.Time {
	if c.StartTime == nil || c.StartTime.Timestamp == nil {
		return nil
	}
	timeStamp, err := strconv.ParseInt(*c.StartTime.Timestamp, 10, 64)
	if err != nil {
		return nil
	}
	t := time.Unix(timeStamp, 0)
	return &t
}

type CalendarBind struct {
	BaseModel
	ProjectKey             string             `json:"project_key,omitempty" gorm:"size:256"`
	WorkItemTypeKey        string             `json:"work_item_type_key,omitempty" gorm:"size:256"`
	WorkItemID             int64              `json:"work_item_id,omitempty" gorm:"uniqueIndex:uniq_wi_cei"`
	CalendarID             string             `json:"calendar_id,omitempty" gorm:"size:256"`
	CalendarEventID        string             `json:"calendar_event_id,omitempty" gorm:"uniqueIndex:uniq_wi_cei;size:256"`
	CalendarEventData      *CalendarEventData `json:"calendar_event_data" gorm:"type:longtext"`
	Bind                   bool               `json:"bind"`
	CalendarEventStartTime *time.Time         `json:"calendar_event_start_time"`
}

type (
	Meeting    larkvc.Meeting
	RecordInfo larkvc.MeetingRecording
)

func (m *Meeting) Scan(input any) error {
	if input == nil {
		return nil
	}
	var err error
	if data, ok := input.([]byte); ok {
		err = json.Unmarshal(data, m)
	}
	if data, ok := input.(string); ok {
		err = json.Unmarshal([]byte(data), m)
	}
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (m *Meeting) Value() (driver.Value, error) {
	ms, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return string(ms), nil
}

func (r *RecordInfo) Scan(input any) error {
	if input == nil {
		return nil
	}
	var err error
	if data, ok := input.([]byte); ok {
		err = json.Unmarshal(data, r)
	}
	if data, ok := input.(string); ok {
		err = json.Unmarshal([]byte(data), r)
	}
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (r *RecordInfo) Value() (driver.Value, error) {
	ms, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return string(ms), nil
}

type VCMeeting struct {
	BaseModel
	CalendarID      string      `json:"calendar_id,omitempty" gorm:"size:256"`
	CalendarEventID string      `json:"calendar_event_id,omitempty" gorm:"uniqueIndex:uniq_cei_mi;size:256"`
	MeetingID       string      `json:"meeting_id" gorm:"size:256;uniqueIndex:uniq_cei_mi"`
	MeetingData     *Meeting    `json:"meeting_data,omitempty" gorm:"type:longtext"`
	RecordInfo      *RecordInfo `json:"record_info,omitempty" gorm:"type:longtext"`
}
