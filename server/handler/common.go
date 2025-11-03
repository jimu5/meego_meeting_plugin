package handler

import (
	"meego_meeting_plugin/common"
	"meego_meeting_plugin/mw"

	"github.com/gofiber/fiber/v2"
)

type PageParam struct {
	PageSize   int `json:"page_size,omitempty"`
	PageNumber int `json:"page_number,omitempty"`
}

var DefaultResp = map[string]string{"status": "ok"}

type ErrMsgResp = mw.ErrMsgResp

func GetMeegoUserKey(c *fiber.Ctx) string {
	return c.Locals(common.MeegoUserKey).(string)
}
