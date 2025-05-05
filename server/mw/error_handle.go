package mw

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"meego_meeting_plugin/common"
	"meego_meeting_plugin/service"
)

type ErrMsgResp struct {
	Msg string `json:"msg"`
}

func ErrorHandle(c *fiber.Ctx) error {
	err := c.Next()
	if err != nil {
		err = beforeErrResp(c, err)
		r := ErrMsgResp{
			Msg: err.Error(),
		}
		c.JSON(&r)
		//
		switch {
		case errors.Is(err, common.ErrorNotLogin):
			c.SendStatus(401)
		default:
			c.SendStatus(500)
		}
		return nil
	}
	return nil
}

func beforeErrResp(c *fiber.Ctx, err error) error {
	if errDetail, ok := err.(larkcore.CodeError); ok {
		// 判断是否是用户token过期
		if errDetail.Code == 99991668 {
			userKey, okU := c.Locals(common.MeegoUserKey).(string)
			if !okU {
				log.Warnf("未获取到用户 user key")
			}
			log.Warnf("用户 token 异常, 刷新 token 时间, user key: %s", userKey)
			if len(userKey) != 0 {
				err = service.Plugin.ResetUserTokenExpired(c.Context(), userKey)
				if err != nil {
					log.Errorf("[beforeErrResp] err: %v", err)
				}
				err = common.ErrorNotLogin
			}
		}
	}
	return err
}
