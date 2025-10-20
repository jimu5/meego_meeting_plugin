package lark_api

import lark "github.com/larksuite/oapi-sdk-go/v3"

type LarkAPI struct {
	CalendarAPI CalendarAPI
	VChatAPI    VChatAPI
	AuthenAPI   AuthenAPI
	IMAPI       IMAPI
}

type LarkClient struct {
	*lark.Client
	appID     string
	appSecret string
}

func NewLarkAPI() LarkAPI {
	client := own_client
	return LarkAPI{
		CalendarAPI: NewCalendarAPI(client),
		VChatAPI:    NewVchatAPI(client),
		AuthenAPI:   NewAuthenAPI(client),
		IMAPI:       NewIMAPI(client),
	}
}
