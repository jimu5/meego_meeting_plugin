package lark_api

import (
	"context"
	"encoding/json"

	"github.com/gofiber/fiber/v2/log"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/service/auth/v3"
	larkauthen "github.com/larksuite/oapi-sdk-go/v3/service/authen/v1"
)

type AuthenAPI struct {
	client *LarkClient
}

func NewAuthenAPI(client *LarkClient) AuthenAPI {
	return AuthenAPI{
		client: client,
	}
}

func (a AuthenAPI) GetUserAccessToken(ctx context.Context, appAccessToken, authCode string) (*larkauthen.CreateAccessTokenRespData, error) {
	req := larkauthen.NewCreateAccessTokenReqBuilder().
		Body(larkauthen.NewCreateAccessTokenReqBodyBuilder().
			GrantType(`authorization_code`).
			Code(authCode).
			Build()).
		Build()
	resp, err := a.client.Authen.AccessToken.Create(ctx, req, larkcore.WithTenantAccessToken(appAccessToken))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if resp.Data == nil {
		log.Error("err login resp empty")
		return nil, ErrResponseNotSuccess
	}
	if !resp.Success() {
		log.Error(resp.Code, resp.Msg, resp.RequestId())
		return nil, err
	}
	return resp.Data, err
}

func (a AuthenAPI) RefreshUserAccessToken(ctx context.Context, appAccessToken, refreshToken string) (*larkauthen.CreateRefreshAccessTokenRespData, error) {
	req := larkauthen.NewCreateRefreshAccessTokenReqBuilder().
		Body(larkauthen.NewCreateRefreshAccessTokenReqBodyBuilder().
			GrantType(`refresh_token`).
			RefreshToken(refreshToken).
			Build()).
		Build()
	resp, err := a.client.Authen.RefreshAccessToken.Create(ctx, req, larkcore.WithTenantAccessToken(appAccessToken))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if !resp.Success() {
		log.Error(resp.Code, resp.Msg, resp.RequestId())
		return nil, err
	}
	return resp.Data, nil
}

type AppAccessTokenResp struct {
	AppAccessToken string `json:"app_access_token,omitempty"`
	Expire         int    `json:"expire,omitempty"`
}

// TODO: 缓存
func (a AuthenAPI) GetAppAccessToken(ctx context.Context) (string, error) {
	req := larkauth.NewInternalAppAccessTokenReqBuilder().
		Body(larkauth.NewInternalAppAccessTokenReqBodyBuilder().
			AppId(a.client.appID).
			AppSecret(a.client.appSecret).
			Build()).
		Build()
	resp, err := a.client.Auth.AppAccessToken.Internal(ctx, req, larkcore.WithTenantAccessToken(""))
	if err != nil {
		log.Error(err)
		return "", err
	}
	if !resp.Success() {
		log.Error(resp.Code, resp.Msg, resp.RequestId())
		return "", resp.CodeError
	}
	var data AppAccessTokenResp
	err = json.Unmarshal(resp.ApiResp.RawBody, &data)
	if err != nil {
		log.Error(err)
		return "", err
	}
	return data.AppAccessToken, nil
}

func (a AuthenAPI) UserInfo(ctx context.Context, userAccessToken string) (*larkauthen.GetUserInfoRespData, error) {
	resp, err := a.client.Authen.UserInfo.Get(ctx, larkcore.WithUserAccessToken(userAccessToken))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if !resp.Success() {
		log.Error(resp.Code, resp.Msg, resp.RequestId())
		return nil, err
	}

	// 业务处理
	return resp.Data, nil
}
