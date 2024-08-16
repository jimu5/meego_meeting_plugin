package handler

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"meego_meeting_plugin/dal"
	"meego_meeting_plugin/service"
)

type ChatAutoBindCalendarParam struct {
	ProjectKey      string `json:"project_key,omitempty"`
	WorkItemTypeKey string `json:"work_item_type_key,omitempty"`
	WorkItemID      int64  `json:"work_item_id,omitempty"`
	Enable          bool   `json:"enable,omitempty"` // 功能启用
}

// ShowAccount godoc
//
//	@Summary		群自动关联日程
//	@Description	群自动关联日程
//	@Tags			Plugin
//	@Produce		json
//	@Param			ChatAutoBindCalendarParam	body	ChatAutoBindCalendarParam	true	"参数"
//	@Success		200
//	@Failure	400	{object}	ErrMsgResp
//	@Router			/api/v1/meego/work_item_meetings/chat_auto_bind	[post]
func ChatAutoBindCalendar(c *fiber.Ctx) error {
	param := ChatAutoBindCalendarParam{}
	err := c.BodyParser(&param)
	if err != nil {
		log.Error(err)
		return err
	}
	record, err := dal.JoinChatRecord.FirstByWorkItemID(c.Context(), param.WorkItemID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error(err)
			errMsg := ErrMsgResp{Msg: err.Error()}
			c.JSON(&errMsg)
			c.SendStatus(400)
			return nil
		}
	}
	meegoUserKey := GetMeegoUserKey(c)
	if record == nil && param.Enable {
		err = service.Plugin.TryJoinChatBycBindFirstCalendar(c.Context(), param.ProjectKey, param.WorkItemTypeKey, param.WorkItemID, meegoUserKey)
		if err != nil {
			log.Error(err)
			errMsg := ErrMsgResp{Msg: "请先创建当前工作项实例的飞书群聊"}
			c.JSON(&errMsg)
			c.SendStatus(400)
			return nil
		}
		// FIXME: 写后读场景, 不应该有 error, 但是需要读主库
		record, err = dal.JoinChatRecord.FirstByWorkItemID(c.Context(), param.WorkItemID)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	if record == nil {
		return nil
	}

	record.Enable = param.Enable
	record.Operator = meegoUserKey
	err = dal.JoinChatRecord.Save(c.Context(), record)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

type GetAutoBindCalendarStatusReq struct {
	WorkItemID int64 `query:"work_item_id" json:"work_item_id"`
}

type GetAutoBindCalendarStatusResp struct {
	Enable bool `json:"enable"`
}

// ShowAccount godoc
//
//	@Summary		获取群自动关联日程的状态
//	@Description	获取群自动关联日程的状态
//	@Tags			Plugin
//	@Produce		json
//	@Param			work_item_id	query	string	true	"查询某个工作项的自动加群状态"
//	@Success		200	{object}	GetAutoBindCalendarStatusResp
//	@Router			/api/v1/meego/work_item_meetings/chat_auto_bind	[get]
func GetAutoBindCalendarStatus(c *fiber.Ctx) error {
	var param GetAutoBindCalendarStatusReq
	err := c.QueryParser(&param)
	if err != nil {
		log.Error(err)
		return err
	}

	resp := GetAutoBindCalendarStatusResp{}
	record, err := dal.JoinChatRecord.FirstByWorkItemID(c.Context(), param.WorkItemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Enable = false
		} else {
			log.Error(err)
			return err
		}
	}
	if record == nil {
		resp.Enable = false
	} else {
		resp.Enable = record.Enable
	}
	c.JSON(&resp)
	return nil
}
