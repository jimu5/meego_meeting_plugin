package mw

import (
	"meego_meeting_plugin/handler"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandle(c *fiber.Ctx) error {
	err := c.Next()
	if err != nil {
		r := handler.ErrMsgResp{
			Msg: err.Error(),
		}
		c.JSON(&r)
		c.SendStatus(500)
		return nil
	}
	return nil
}
