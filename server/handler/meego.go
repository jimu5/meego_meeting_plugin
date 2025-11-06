package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"runtime/debug"
	"time"

	"meego_meeting_plugin/dal"
	"meego_meeting_plugin/model"
	"meego_meeting_plugin/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type MeegoLarkLoginParam struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

// ShowAccount godoc
//
//	@Summary		飞书授权登录
//	@Description	飞书 code 授权登录
//	@Tags			Plugin
//	@Param			MeegoLarkLoginParam	query	MeegoLarkLoginParam	true	"参数"
//	@Success		302
//	@Router			/api/v1/meego/lark/auth	[get]
func MeegoLarkLogin(c *fiber.Ctx) error {
	// TODO: 这个接口需要搞一些安全机制, 比如说搞一个类似 csrf 的机制
	param := MeegoLarkLoginParam{}
	err := c.QueryParser(&param)
	if err != nil {
		log.Errorf("[MeegoLarkLogin] err parser, err: %v", err)
		return err
	}
	// 先尝试 url 解码
	state, err := url.QueryUnescape(param.State)
	if err != nil {
		log.Warn(err)
		state = param.State
	}
	// 解析 state 链接
	stateString, err := base64.StdEncoding.DecodeString(state)
	if err != nil {
		return err
	}
	stateURL, err := url.Parse(string(stateString))
	if err != nil {
		log.Errorf("[MeegoLarkLogin] err parser state url, err: %v", err)
		return err
	}
	fmt.Println(stateURL.RawQuery)
	// 从 stateURL 中解析出 userKey
	meegoUserKey := stateURL.Query().Get(MeegoUserKey)
	if len(meegoUserKey) == 0 {
		log.Infof("meegoUserKey: %s", meegoUserKey)
		meegoUserKey = stateURL.Query().Get("meego-user-key")
		if len(meegoUserKey) == 0 {
			fmt.Println(stateURL.Query(), meegoUserKey)
			log.Errorf("[MeegoLarkLogin] parse empty meego user key")
			return ErrInvalidParam
		}
	}
	userTokenInfo, err := service.Lark.GetUserAccessToken(c.Context(), param.Code)
	if err != nil {
		log.Error(err)
		return err
	}
	if userTokenInfo == nil {
		return ErrInvalidUserInfo
	}
	// 先查询
	userModelInfo, err := dal.User.QueryByMeegoUserKey(c.Context(), meegoUserKey)
	if err != nil || userModelInfo == nil {
		userModelInfo = &model.User{}
	}
	userModelInfo.MeegoUserKey = meegoUserKey
	userModelInfo.LarkUserAccessToken = userTokenInfo.AccessToken
	userModelInfo.LarkUserRefreshToken = userTokenInfo.RefreshToken
	log.Info("login", meegoUserKey, userTokenInfo.AccessTokenExpire, userTokenInfo.RefreshTokenExpire)
	userModelInfo.LarkUserAccessTokenExpireAt = time.Unix(userTokenInfo.AccessTokenExpire, 0)
	userModelInfo.LarkUserRefreshTokenExpiredAt = time.Unix(userTokenInfo.RefreshTokenExpire, 0)
	userData, err := service.Lark.LarkAPI.AuthenAPI.UserInfo(c.Context(), userModelInfo.LarkUserAccessToken)
	if userData != nil {
		userModelInfo.LarkUserID = *userData.UserId
	}
	data, _ := json.Marshal(userData)
	userModelInfo.LarkUserInfo = string(data)
	err = service.User.SaveUser(c.Context(), userModelInfo)
	if err != nil {
		log.Error(err)
		return err
	}
	// 302 跳转到 state
	err = c.Redirect(string(stateString), 302)
	if err != nil {
		return err
	}
	// 异步触发这个用户的任务
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("Recovered in ProcessTasksByUser: %v, stack: %s", r, debug.Stack())
			}
		}()
		ctx := context.Background()
		err := service.Cron.ProcessTasksByUser(ctx, *userModelInfo)
		if err != nil {
			log.Error(err)
		}
	}()

	return nil
}
