package lark_api

import (
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"meego_meeting_plugin/config"
)

var own_client = &LarkClient{
	Client:    lark.NewClient(config.LarkAppID, config.LarkAppSecret, lark.WithEnableTokenCache(true)),
	appID:     config.LarkAppID,
	appSecret: config.LarkAppSecret,
}
