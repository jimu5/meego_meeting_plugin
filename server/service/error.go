package service

import "errors"

var (
	ErrPrimaryCalendar    = errors.New("primary calendar not found")
	ErrNilOpenApiResponse = errors.New("nil open api response")
	ErrNilCalendarTime    = errors.New("err nil calendar time")
	ErrNilMeeting         = errors.New("err nil meeting")
	ErrNilMeetingRecord   = errors.New("err nil meeting record")

	ErrNilUser = errors.New("err nil user")
	ErrToken   = errors.New("err token")

	ErrTryBotJoinChat = errors.New("邀请机器人加入群聊失败，请检查是否创建群聊，详细错误: ")
)
