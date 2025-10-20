package mw

import (
	"meego_meeting_plugin/common"
	"meego_meeting_plugin/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type ReqHeaders struct {
	MeegoUserKey string `json:"meego_user_key,omitempty"`
}

func GetUser(c *fiber.Ctx) error {
	err := GetAndSetUserInfo(c)
	if err != nil {
		return err
	}
	return c.Next()
}

func GetAndSetUserInfo(c *fiber.Ctx) error {
	userKey := c.Get(common.MeegoUserKey)
	if len(userKey) == 0 {
		c.SendStatus(401)
		return common.ErrorNotLogin
	}
	c.Locals(common.MeegoUserKey, userKey)
	userInfo, err := service.Plugin.GetUserInfoByMeegoUserKey(c.Context(), userKey, false)
	if err != nil {
		log.Error(err)
		c.SendStatus(401)
		return common.ErrorNotLogin
	}
	c.Locals(common.LarkUserAccessToken, userInfo.LarkUserAccessToken)
	return nil
}
