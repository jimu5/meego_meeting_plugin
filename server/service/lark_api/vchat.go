package lark_api

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkvc "github.com/larksuite/oapi-sdk-go/v3/service/vc/v1"
)

type VChatAPI struct {
	client *LarkClient
}

func NewVchatAPI(client *LarkClient) VChatAPI {
	return VChatAPI{
		client: client,
	}
}

func (v VChatAPI) GetMeeting(ctx context.Context, meetingID string, userAccessToken string) (*GetMeetingRespData, error) {
	req := larkvc.NewGetMeetingReqBuilder().MeetingId(meetingID).Build()
	resp, err := v.client.Vc.Meeting.Get(ctx, req, larkcore.WithUserAccessToken(userAccessToken))
	if err != nil {
		log.Errorf("[VChatAPI] GetMeeting, err: %v")
		return nil, err
	}
	if !resp.Success() {
		log.Errorf("[VChatAPI] GetMeeting resp not success, code: %v, msg: %v, LOGID: %s", resp.Code, resp.Msg, resp.RequestId())
		return nil, NewErrResponseNotSuccess(resp.Code, resp.Msg)
	}
	return (*GetMeetingRespData)(resp.Data), nil
}

// time 都是 unix 时间, 需要遍历将全部查出
func (v VChatAPI) GetMeetingsListByNo(ctx context.Context, meetingNo string, startTime, endTime string, userAccessToken string, pageParam *PageParam) ([]*Meeting, error) {
	// 创建请求对象
	reqBuilder := larkvc.NewListByNoMeetingReqBuilder().
		MeetingNo(meetingNo).
		StartTime(startTime).
		EndTime(endTime)

	pageSize := 50
	pageToken := ""
	if pageParam != nil {
		if pageParam.PageSize != 0 {
			pageSize = pageParam.PageSize
		}
		pageToken = pageParam.PageToken
	}

	req := reqBuilder.PageSize(pageSize).PageToken(pageToken).Build()

	// 发起请求
	resp, err := v.client.Vc.Meeting.ListByNo(context.Background(), req, larkcore.WithUserAccessToken(userAccessToken))
	if err != nil {
		log.Errorf("[VChatAPI] GetMeeting, err: %v")
		return nil, err
	}
	if !resp.Success() {
		log.Errorf("[VChatAPI] GetMeeting resp not success, code: %v, msg: %v, LOGID: %s", resp.Code, resp.Msg, resp.RequestId())
		return nil, NewErrResponseNotSuccess(resp.Code, resp.Msg)
	}
	if len(resp.Data.MeetingBriefs) == 0 {
		return []*Meeting{}, nil
	}
	result := make([]*Meeting, 0)
	for _, m := range resp.Data.MeetingBriefs {
		if m != nil {
			result = append(result, (*Meeting)(m))
		}
	}
	if resp.Data != nil && resp.Data.HasMore != nil && *resp.Data.HasMore {
		if resp.Data.PageToken != nil {
			news, errN := v.GetMeetingsListByNo(ctx, meetingNo, startTime, endTime, userAccessToken, &PageParam{
				PageSize:  pageSize,
				PageToken: pageToken,
			})
			if errN != nil {
				return result, err
			}
			if len(news) == 0 {
				return result, nil
			}
			result = append(result, news...)
		}
	}
	return result, nil
}

func (v VChatAPI) GetMeetingRecord(ctx context.Context, meetingID string, userAccessToken string) (*GetMeetingRecordingRespData, error) {
	req := larkvc.NewGetMeetingRecordingReqBuilder().
		MeetingId(meetingID).
		Build()
	resp, err := v.client.Vc.MeetingRecording.Get(ctx, req, larkcore.WithUserAccessToken(userAccessToken))

	// 处理错误
	if err != nil {
		log.Errorf("[VChatAPI] GetMeetingRecord err, meeting ID: %s, err: %v", meetingID, err)
		return nil, err
	}

	// 服务端错误处理
	if !resp.Success() {
		log.Errorf("[VChatAPI] GetMeetingRecord resp not success, code: %v, msg: %v, LOGID: %s", resp.Code, resp.Msg, resp.RequestId())
		return nil, resp.CodeError
	}
	return (*GetMeetingRecordingRespData)(resp.Data), nil
}
