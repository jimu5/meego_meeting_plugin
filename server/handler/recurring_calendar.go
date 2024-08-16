package handler

import (
	"context"
	"meego_meeting_plugin/model"
)

func handleRecurring_calendar(ctx context.Context, bindInfo *model.CalendarBind) error {
	if bindInfo == nil || bindInfo.CalendarEventData == nil {
		return nil
	}
	//eventData := bindInfo.CalendarEventData
	return nil
}
