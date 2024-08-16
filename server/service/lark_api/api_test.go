package lark_api

import (
	"context"
	"fmt"
	"testing"
)

func TestAPI(t *testing.T) {
	ctx := context.Background()
	t.Run("AuthenTest", func(t *testing.T) {
		resp, err := API.AuthenAPI.GetAppAccessToken(ctx)
		fmt.Println(resp, err)
	})
}
