package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"meego_meeting_plugin/service"
)

type RefreshCalendarParam struct {
	ProjectKey      string `json:"project_key,omitempty"`
	WorkItemTypeKey string `json:"work_item_type_key,omitempty"`
	WorkItemID      int64  `json:"work_item_id,omitempty"`
}

// ShowAccount godoc
//
//	@Summary		刷新工作项关联的会议最新信息
//	@Description	刷新工作项关联的会议最新信息
//	@Tags			Plugin
//	@Produce		json
//	@Param			RefreshCalendarParam	body	RefreshCalendarParam	true	"参数"
//	@Success		200
//	@Router			/api/v1/meego/work_item_meetings/refresh	[post]
func RefreshCalendar(c *fiber.Ctx) error {
	param := RefreshCalendarParam{}
	err := c.BodyParser(&param)
	if err != nil {
		log.Error(err)
		return err
	}
	//token := c.Locals(common.LarkUserAccessToken).(string)
	//operator := c.Locals(common.MeegoUserKey).(string)
	err = service.Plugin.RefreshBind(c.Context(), param.WorkItemID)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
