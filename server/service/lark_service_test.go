package service

import (
	"context"
	"fmt"
	"testing"

	"meego_meeting_plugin/service/lark_api"
)

const TestUserToken = ""

func TestLarkService(t *testing.T) {
	larkService := NewLarkService()
	ctx := context.Background()
	userToken := TestUserToken

	t.Run("SearchCalendar", func(t *testing.T) {
		res, err := larkService.SearchCalendar(ctx, "", userToken, lark_api.PageParam{})
		fmt.Println(res, err)
	})

	t.Run("GetMeetingRecordInfoByCalendar", func(t *testing.T) {
		eventID := ""
		res, err := larkService.GetMeetingRecordInfoByCalendar(ctx, eventID, userToken)
		fmt.Println(res, err)
	})

	t.Run("GetMeetingInfo", func(t *testing.T) {
		meetingID := "7313131366149226524"
		res, err := larkService.GetMeetingInfo(ctx, meetingID, userToken)
		fmt.Println(res, err)
	})

	t.Run("GetUserAccessToken", func(t *testing.T) {
		code := ""
		res, err := larkService.GetUserAccessToken(ctx, code)
		fmt.Println(res, err)
	})
}
