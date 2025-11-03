package meego_api

type MeegoAPI struct {
	Chat     ChatAPI
	User     UserAPI
	WorkItem WorkItemAPI
}

func NewMeegoAPI() MeegoAPI {
	c := own_client
	return MeegoAPI{
		Chat:     NewChatAPI(c),
		User:     NewUserAPI(c),
		WorkItem: NewWorkItemAPI(c),
	}
}
