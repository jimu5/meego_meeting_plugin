package service

import "meego_meeting_plugin/service/meego_api"

var Meego MeegoService

type MeegoService struct {
	MeegoAPI meego_api.MeegoAPI
}

func NewMeegoService() MeegoService {
	return MeegoService{
		MeegoAPI: meego_api.API,
	}
}
