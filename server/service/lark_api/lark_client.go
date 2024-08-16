package lark_api

import (
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"meego_meeting_plugin/config"
)

var own_client = &LarkClient{
	Client:    lark.NewClient(config.APPID, config.APPSECRET, lark.WithEnableTokenCache(true)),
	appID:     config.APPID,
	appSecret: config.APPSECRET,
}
