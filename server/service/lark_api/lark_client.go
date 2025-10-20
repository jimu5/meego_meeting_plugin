package lark_api

import (
	"meego_meeting_plugin/config"

	lark "github.com/larksuite/oapi-sdk-go/v3"
)

var own_client *LarkClient

func InitOwnClient() *LarkClient {
	l := &LarkClient{
		Client: lark.NewClient(
			config.GetAPPConfig().LarkAppID,
			config.GetAPPConfig().LarkAppSecret,
			lark.WithEnableTokenCache(true),
		),
		appID:     config.GetAPPConfig().LarkAppID,
		appSecret: config.GetAPPConfig().LarkAppSecret,
	}

	own_client = l

	return l
}
