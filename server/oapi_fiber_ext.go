package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	"github.com/valyala/fasthttp"
	"net/http"
)

func doProcess(writer *fasthttp.Response, req *fasthttp.Request, reqHandler larkevent.IReqHandler, options ...larkevent.OptionFunc) {
	// 转换http请求对象为标准请求对象
	ctx := context.Background()
	eventReq, err := translate(ctx, req)
	if err != nil {
		writer.Header.SetStatusCode(http.StatusInternalServerError)
		writer.SetBodyString(err.Error())
		return
	}

	//处理请求
	eventResp := reqHandler.Handle(ctx, eventReq)

	// 回写结果
	err = write(ctx, writer, eventResp)
	if err != nil {
		reqHandler.Logger().Error(ctx, fmt.Sprintf("write resp result error:%s", err.Error()))
	}
}

func NewCardActionHandlerFunc(cardActionHandler *larkcard.CardActionHandler, options ...larkevent.OptionFunc) func(c *fiber.Ctx) error {

	// 构建模板类
	cardActionHandler.InitConfig(options...)
	return func(c *fiber.Ctx) error {
		doProcess(c.Response(), c.Request(), cardActionHandler, options...)
		return nil
	}
}

func NewEventHandlerFunc(eventDispatcher *dispatcher.EventDispatcher, options ...larkevent.OptionFunc) func(c *fiber.Ctx) error {
	eventDispatcher.InitConfig(options...)
	return func(c *fiber.Ctx) error {
		doProcess(c.Response(), c.Request(), eventDispatcher, options...)
		return nil
	}
}

func processError(ctx context.Context, logger larkcore.Logger, path string, err error) *larkevent.EventResp {
	header := map[string][]string{}
	header[larkevent.ContentTypeHeader] = []string{larkevent.DefaultContentType}
	eventResp := &larkevent.EventResp{
		Header:     header,
		Body:       []byte(fmt.Sprintf(larkevent.WebhookResponseFormat, err.Error())),
		StatusCode: http.StatusInternalServerError,
	}
	logger.Error(ctx, fmt.Sprintf("event handle err:%s, %v", path, err))
	return eventResp
}

func write(ctx context.Context, writer *fasthttp.Response, eventResp *larkevent.EventResp) error {
	writer.SetStatusCode(eventResp.StatusCode)
	for k, vs := range eventResp.Header {
		for _, v := range vs {
			writer.Header.Add(k, v)
		}
	}

	if len(eventResp.Body) > 0 {
		fmt.Println(eventResp.Body)
		writer.SetBodyRaw(eventResp.Body)
		return nil
	}
	return nil
}

func translate(ctx context.Context, req *fasthttp.Request) (*larkevent.EventReq, error) {
	headers := make(map[string][]string)
	req.Header.VisitAll(func(key, value []byte) {
		keyStr := string(key)
		valueStr := string(value)
		headers[keyStr] = append(headers[keyStr], valueStr)
	})
	eventReq := &larkevent.EventReq{
		Header:     headers,
		Body:       req.Body(),
		RequestURI: req.URI().String(),
	}

	return eventReq, nil
}
