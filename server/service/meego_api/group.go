package meego_api

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	sdk "github.com/larksuite/project-oapi-sdk-golang"
	"github.com/larksuite/project-oapi-sdk-golang/core"
	"github.com/larksuite/project-oapi-sdk-golang/service/chat"
)

// 群组相关操作
type ChatAPI struct {
	client *sdk.Client
}

func NewChatAPI(c *sdk.Client) ChatAPI {
	return ChatAPI{c}
}

type BotJoinChatParam struct {
	ProjectKey      string
	WorkItemTypeKey string
	AppIDs          []string
	WorkItemID      int64
	MeegoUserKey    string
}

func (g ChatAPI) BotJoinChat(ctx context.Context, param BotJoinChatParam) (*chat.BotJoinChatInfo, error) {
	req := chat.NewBotJoinChatReqBuilder().
		ProjectKey(param.ProjectKey).
		WorkItemTypeKey(param.WorkItemTypeKey).
		WorkItemID(param.WorkItemID).
		AppIDs(param.AppIDs).Build()
	resp, err := g.client.Chat.BotJoinChat(ctx, req, core.WithUserKey(param.MeegoUserKey))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if !resp.Success() {
		log.Error(resp.Code(), resp.ErrMsg, resp.RequestId())
		return nil, ErrRespNotSuccess
	}
	return resp.Data, nil
}
