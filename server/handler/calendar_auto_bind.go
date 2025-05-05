package handler

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"meego_meeting_plugin/dal"
	"meego_meeting_plugin/mw"
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

	meegoUserKey := GetMeegoUserKey(c)
	err = service.Plugin.AutoBindCalendar(c.Context(), param.Enable, param.ProjectKey, param.WorkItemTypeKey, param.WorkItemID, meegoUserKey)

	if err != nil {
		log.Errorf("[ChatAutoBindCalendar] err: %v", err)
		return err
	}

	return nil
}

type GetAutoBindCalendarStatusReq struct {
	ProjectKey      string `query:"project_key" json:"project_key"`
	WorkItemTypeKey string `query:"work_item_type_key" json:"work_item_type_key"`
	WorkItemID      int64  `query:"work_item_id" json:"work_item_id"`
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
		log.Errorf("[GetAutoBindCalendarStatus] err: %v", err)
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
		err = mw.GetAndSetUserInfo(c)
		if err != nil {
			return err
		}
		meegoUserKey := GetMeegoUserKey(c)
		if len(param.ProjectKey) == 0 || len(param.WorkItemTypeKey) == 0 || param.WorkItemID == 0 {
			log.Warnf("[GetAutoBindCalendarStatus] error param empty, projectKey: %s, workItemTypeKey: %s, workItemID: %d",
				param.ProjectKey, param.WorkItemTypeKey, param.WorkItemID)
		} else {
			err = service.Plugin.AutoBindCalendar(c.Context(), true, param.ProjectKey, param.WorkItemTypeKey, param.WorkItemID, meegoUserKey)
			if err != nil {
				// 这里尝试自动绑定失败了不要降级
				log.Errorf("[GetAutoBindCalendarStatus] err auto bind, err: %v", err)
			} else {
				resp.Enable = true
			}
		}

	} else {
		resp.Enable = record.Enable
	}
	c.JSON(&resp)
	return nil
}
