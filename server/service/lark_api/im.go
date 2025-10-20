package lark_api

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type IMAPI struct {
	client *LarkClient
}

func NewIMAPI(client *LarkClient) IMAPI {
	return IMAPI{
		client: client,
	}
}

// 发送消息
func (i *IMAPI) CreateTextMessage(ctx context.Context, recIDType, recID, text string) (*larkim.CreateMessageRespData, error) {
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(recIDType).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(recID).
			MsgType(`text`).
			Content(fmt.Sprintf(`{"text":"%s"}`, text)).
			Build()).
		Build()
	// 发起请求
	resp, err := i.client.Im.V1.Message.Create(ctx, req)
	if err != nil {
		log.Errorf("[IMAPI] CreateTextMessage err: %v", err)
		return nil, err
	}
	if resp == nil {
		log.Errorf("[IMAPI] CreateTextMessage resp is nil")
		return nil, ErrInvalidResponse
	}
	if !resp.Success() {
		log.Errorf("[IMAPI] CreateTextMessage resp is not success, code: %d, msg: %s", resp.Code, resp.Msg)
		return nil, NewErrResponseNotSuccess(resp.Code, resp.Msg)
	}
	return resp.Data, nil
}
