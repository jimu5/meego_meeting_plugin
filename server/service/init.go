package service

import (
	"meego_meeting_plugin/service/lark_api"
	"meego_meeting_plugin/service/meego_api"
)

func InitClient() {
	lark_api.InitOwnClient()
	meego_api.InitOwnClient()

	lark_api.API = lark_api.NewLarkAPI()
	meego_api.API = meego_api.NewMeegoAPI()

	Lark = NewLarkService()
}
