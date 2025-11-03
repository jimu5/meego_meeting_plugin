package handler

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"meego_meeting_plugin/config"
	"meego_meeting_plugin/dal"
	"meego_meeting_plugin/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

const (
	MeegoEventSourceNormal     = "normal"
	MeegoEventSourceOpenAPI    = "openapi"
	MeegoEventSourceSystem     = "system"
	MeegoEventSourceAutomation = "automation"
)

type WorkItemHandleEventBody struct {
	LogID       string `json:"log_id,omitempty"`
	RequestTime int64  `json:"request_time,omitempty"`
	Signature   string `json:"signature,omitempty"`

	SourcePluginID   string `json:"source_plugin_id,omitempty"`
	SourcePluginName string `json:"source_plugin_name,omitempty"`

	Data WorkItemEventData `json:"data"`
}

type WorkItemEventData struct {
	EventType     int    `json:"event_type,omitempty"`
	IdempotentKey string `json:"idempotent_key,omitempty"`
	ProjectKey    string `json:"project_key,omitempty"`
	ProjectName   string `json:"project_name,omitempty"`
	// normal,openapi,system,automation 区分来源，normal代表该操作来自普通用户操作，openapi代表该操作来自于OpenAPI，system代表该操作来自于系统行为，automation代表该操作来自自动化。
	Source string `json:"source,omitempty"`

	FieldInfo    []FieldInfo    `json:"field_info,omitempty"`
	WorkItemInfo []WorkItemInfo `json:"work_item_info"`
	UserInfo     UserInfo       `json:"user_info"`
}

type FieldInfo struct {
	FieldKey        string `json:"field_key,omitempty"`
	AfterFieldValue string `json:"after_field_value,omitempty"`
}

type WorkItemInfo struct {
	WorkItemID      int64  `json:"work_item_id,omitempty"`
	WorkItemName    string `json:"work_item_name,omitempty"`
	WorkItemTypeKey string `json:"work_item_type_key,omitempty"`
}

type UserInfo struct {
	UserKey string `json:"user_key,omitempty"`
	Email   string `json:"email,omitempty"`
}

func MeegoEventHandler(c *fiber.Ctx) error {
	// 获取请求body
	const fn = "MeegoEventHandler"
	body := WorkItemHandleEventBody{}
	// 解析请求body
	if err := c.BodyParser(&body); err != nil {
		log.Errorf("%s, err: %v", fn, err)
		return err
	}
	// 校验签名:(如果后续有多个不同 url 监听逻辑, 可以改成中间件)
	exceptedSignature := calculateSignature(body.SourcePluginID, fmt.Sprintf("%d", body.RequestTime), config.GetAPPConfig().MeegoEventCallbackToken)
	if body.Signature != exceptedSignature {
		log.Errorf("%s, exceptedSignature: %s, body.Signature: %s", fn, exceptedSignature, body.Signature)
		return nil
	}

	// 处理业务逻辑
	// 自动绑定
	err := UseMeegoEventAutoBind(c, &body)
	if err != nil {
		log.Errorf("%s, err: %v", fn, err)
		return nil
	}

	return nil
}

func UseMeegoEventAutoBind(c *fiber.Ctx, body *WorkItemHandleEventBody) error {
	const fn = "UseMeegoEventAutoBind"
	if body == nil || len(body.Data.WorkItemInfo) == 0 {
		log.Warnf("%s, body is nil", fn)
		return nil
	}
	var (
		workItemID      = body.Data.WorkItemInfo[0].WorkItemID
		workItemTypeKey = body.Data.WorkItemInfo[0].WorkItemTypeKey
		projectKey      = body.Data.ProjectKey
		meegoUserKey    = body.Data.UserInfo.UserKey
		meegoUserEmail  = body.Data.UserInfo.Email
	)
	if len(projectKey) == 0 || len(workItemTypeKey) == 0 || workItemID == 0 {
		log.Warnf("%s, projectKey: %s, workItemTypeKey: %s, workItemID: %d is empty", fn, projectKey, workItemTypeKey, workItemID)
		return nil
	}
	if len(meegoUserKey) == 0 {
		log.Warnf("%s, meegoUserKey: %s, meegoUserEmail: %s is empty", fn, meegoUserKey, meegoUserEmail)
		return nil
	}

	// 判断当前变更的对象是否有群组
	var (
		hadGroup bool
		//hadGroupField bool
	)

	for _, f := range body.Data.FieldInfo {
		if f.FieldKey == "chat_group" {
			//hadGroupField = true
			if f.AfterFieldValue != "" {
				hadGroup = true
			}
		}
	}

	log.Infof("%s, workItemID: %d hadGroup: %v, field len: %d", fn, workItemID, hadGroup, len(body.Data.FieldInfo))
	if !hadGroup {
		log.Infof("%s, workItemID: %d not has group", fn, workItemID)
		return nil
	}

	record, err := dal.JoinChatRecord.FirstByWorkItemID(c.Context(), workItemID)
	if err != nil {
		log.Errorf("%s, err: %v", fn, err)
		return err
	}

	// 只有首次发现这个记录的时候才会尝试绑定
	if record != nil {
		log.Infof("%s, workItemID: %d already bind", fn, workItemID)
		return nil
	}

	err = service.Plugin.AutoBindCalendar(c.Context(), true, projectKey, workItemTypeKey, workItemID, meegoUserKey)
	if err != nil {
		log.Errorf("%s, err: %v", fn, err)
		return err
	}
	return nil
}

func calculateSignature(pluginId, requestTime, token string) string {
	var b strings.Builder
	b.WriteString(pluginId)
	b.WriteString(requestTime)
	b.WriteString(token)
	bs := []byte(b.String())
	h := sha256.New()
	h.Write(bs)
	bs = h.Sum(nil)
	sig := fmt.Sprintf("%x", bs)
	return sig
}
