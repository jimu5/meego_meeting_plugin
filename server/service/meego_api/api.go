package meego_api

type MeegoAPI struct {
	Chat ChatAPI
	User UserAPI
}

func NewMeegoAPI() MeegoAPI {
	c := own_client
	return MeegoAPI{
		Chat: NewChatAPI(c),
		User: NewUserAPI(c),
	}
}
