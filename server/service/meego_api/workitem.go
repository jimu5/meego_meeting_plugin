package meego_api

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	sdk "github.com/larksuite/project-oapi-sdk-golang"
	"github.com/larksuite/project-oapi-sdk-golang/core"
	"github.com/larksuite/project-oapi-sdk-golang/service/workitem"
)

// workitem 相关接口
type WorkItemAPI struct {
	client *sdk.Client
}

func NewWorkItemAPI(c *sdk.Client) WorkItemAPI {
	return WorkItemAPI{
		c,
	}
}

func (w WorkItemAPI) GetWorkItem(ctx context.Context, meegoUserKey, projectKey, workItemTypeKey string, workItemIDs []int64, fields []string) ([]*workitem.WorkItemInfo, error) {
	req := workitem.NewQueryWorkItemDetailReqBuilder().ProjectKey(projectKey).
		WorkItemTypeKey(workItemTypeKey).WorkItemIDs(workItemIDs).Fields(fields).Build()
	resp, err := w.client.WorkItem.QueryWorkItemDetail(ctx, req, core.WithUserKey(meegoUserKey))
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
