package handler

import (
	"github.com/gofiber/fiber/v2"
	"meego_meeting_plugin/common"
	"meego_meeting_plugin/mw"
	"meego_meeting_plugin/util"
)

type PageParam struct {
	PageSize   int `json:"page_size,omitempty"`
	PageNumber int `json:"page_number,omitempty"`
}

func getPointerInfo[T any](s *T) T {
	return util.GetPointerInfo(s)
}

var DefaultResp = map[string]string{"status": "ok"}

type ErrMsgResp = mw.ErrMsgResp

func GetMeegoUserKey(c *fiber.Ctx) string {
	return c.Locals(common.MeegoUserKey).(string)
}
