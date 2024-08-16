package meego_api

import (
	"context"
	"fmt"
	"testing"
)

func TestMeegoAPI(t *testing.T) {
	ctx := context.Background()
	t.Run("test add chat", func(t *testing.T) {
		resp, err := API.Chat.BotJoinChat(ctx, BotJoinChatParam{
			ProjectKey:      "",
			WorkItemTypeKey: "",
			WorkItemID:      0,
			MeegoUserKey:    "",
			AppIDs:          []string{""},
		})
		fmt.Print(resp, err)
	})
	t.Run("test query user", func(t *testing.T) {
		resp, err := API.User.GetUserInfo(ctx, []string{""})
		fmt.Print(resp, err)
	})
}
