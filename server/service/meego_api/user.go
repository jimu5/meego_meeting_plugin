package meego_api

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	sdk "github.com/larksuite/project-oapi-sdk-golang"
	"github.com/larksuite/project-oapi-sdk-golang/service/user"
)

type UserAPI struct {
	client *sdk.Client
}

func NewUserAPI(c *sdk.Client) UserAPI {
	return UserAPI{
		c,
	}
}

// TODO: 以后需要处理下分页, 最大 100 个 userKey
func (u UserAPI) GetUserInfo(ctx context.Context, userKeys []string) ([]*user.UserBasicInfo, error) {
	req := user.NewQueryUserDetailReqBuilder().UserKeys(userKeys).Build()
	resp, err := u.client.User.QueryUserDetail(ctx, req)
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
