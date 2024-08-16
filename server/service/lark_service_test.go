package service

import (
	"context"
	"fmt"
	"meego_meeting_plugin/service/lark_api"
	"testing"
)

const TestUserToken = "u-fEnSq3wjRdJpPEjbG4xY2Z54mjD05kP9jww0hlY004VA"

func TestLarkService(t *testing.T) {
	larkService := NewLarkService()
	ctx := context.Background()
	userToken := TestUserToken

	t.Run("SearchCalendar", func(t *testing.T) {
		res, err := larkService.SearchCalendar(ctx, "", userToken, lark_api.PageParam{})
		fmt.Println(res, err)
	})

	t.Run("GetMeetingRecordInfoByCalendar", func(t *testing.T) {
		eventID := "aefedb6a-637d-4cc9-b877-ecea2fac505e_0"
		res, err := larkService.GetMeetingRecordInfoByCalendar(ctx, "", eventID, userToken)
		fmt.Println(res, err)
	})

	t.Run("GetMeetingInfo", func(t *testing.T) {
		meetingID := "7313131366149226524"
		res, err := larkService.GetMeetingInfo(ctx, meetingID, userToken)
		fmt.Println(res, err)
	})

	t.Run("GetUserAccessToken", func(t *testing.T) {
		code := "bdbm88dc0ebf4ba093bea885d5002216"
		res, err := larkService.GetUserAccessToken(ctx, code)
		fmt.Println(res, err)
	})
}
