package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"meego_meeting_plugin/config"
	"meego_meeting_plugin/dal"
	"meego_meeting_plugin/model"
	"meego_meeting_plugin/service/lark_api"
	"meego_meeting_plugin/service/meego_api"

	"github.com/avast/retry-go"
	"github.com/gofiber/fiber/v2/log"
	"github.com/larksuite/project-oapi-sdk-golang/service/user"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

var Plugin = NewPluginService()

type PluginService struct {
}

func NewPluginService() PluginService {
	return PluginService{}
}

func (p PluginService) List() {

}

type BindCalendarParam struct {
	ProjectKey      string `json:"project_key"`
	WorkItemTypeKey string `json:"work_item_type_key"`
	WorkItemID      int64  `json:"work_item_id"`
	CalendarEventID string `json:"calendar_event_id"` // 日程 ID
}

func (p PluginService) BindCalendar(ctx context.Context, param BindCalendarParam, userToken, meegoUserKey string) error {
	operator := meegoUserKey
	meetingInfo, err := Lark.GetMeetingRecordInfoByCalendar(ctx, param.CalendarEventID, userToken)
	if err != nil {
		log.Errorf("[PluginService] BindCalendar GetMeetingRecordInfoByCalendar err, err: %v", err)
		return err
	}

	meetingsIDs := make([]string, 0, len(meetingInfo.Meetings))
	for _, meeting := range meetingInfo.Meetings {
		if meeting == nil {
			continue
		}
		meetingsIDs = append(meetingsIDs, *meeting.Id)
	}
	meetingInfos, err := Lark.MGetMeetingInfo(ctx, meetingsIDs, userToken)
	if err != nil {
		log.Errorf("[PluginService] BindCalendar err, err: %v", err)
		return err
	}
	modelCalendarBindInfo := model.CalendarBind{
		CalendarID:             meetingInfo.CalendarID,
		CalendarEventID:        param.CalendarEventID,
		WorkItemID:             param.WorkItemID,
		WorkItemTypeKey:        param.WorkItemTypeKey,
		ProjectKey:             param.ProjectKey,
		CalendarEventStartTime: meetingInfo.GetEventStartTime(),
		CalendarEventData:      (*model.CalendarEventData)(&meetingInfo.EventData),
		Bind:                   true,
	}
	err = dal.CalendarBind.CreateOrUpdateCalendarBind(ctx, &modelCalendarBindInfo, operator)
	if err != nil {
		return err
	}
	modelMeetings := MeetingInfos2ModelVCMeeting(meetingInfos, meetingInfo.CalendarID, param.CalendarEventID)
	err = dal.CalendarBind.CreateOrUpdateCalendarMeetings(ctx, modelMeetings, operator)
	if err != nil {
		return err
	}
	// p.TryJoinChatBycBindFirstCalendar(ctx, param.ProjectKey, param.WorkItemTypeKey, param.WorkItemID, operator)
	go func() {
		errG := p.RetryRefreshMeetingRecordTask(ctx, meetingsIDs, userToken)
		if errG != nil {
			log.Errorf("BindCalendar go task RetryRefreshMeetingRecordTask err: %v", err)
		}
	}()
	return nil
}

// 获取或刷新 user
func (p PluginService) GetUserInfoByMeegoUserKey(ctx context.Context, meegoUserKey string, sendACLMsg bool) (model.User, error) {
	userKey := meegoUserKey
	userInfo, err := User.GetUserInfoByMeegoUserKey(ctx, userKey)
	if err != nil {
		// 没有找到对应的用户, 尝试引导用户去给应用授权
		// TODO: 重新授权之后应当重新执行任务
		if sendACLMsg {
			err = p.SendMsgForACL(ctx, meegoUserKey)
			if err != nil {
				log.Errorf("SendMsgForACL err: %v", err)
				return model.User{}, err
			}
		}
		return model.User{}, err
	}
	refreshTag := false
	if userInfo.LarkUserAccessTokenExpireAt.Add(-10*time.Second).UnixMilli() < time.Now().UnixMilli() {
		refreshTag = true
	}
	userData, err := Lark.LarkAPI.AuthenAPI.UserInfo(ctx, userInfo.LarkUserAccessToken)
	if err != nil {
		refreshTag = true
	}
	if refreshTag {
		log.Info("refreshToken: " + userKey)
		userTokenInfo, errQ := Lark.RefreshUserAccessToken(ctx, userInfo.LarkUserRefreshToken)
		if errQ != nil {
			return model.User{}, errQ
		}
		userInfo.LarkUserAccessToken = userTokenInfo.AccessToken
		userInfo.LarkUserRefreshToken = userTokenInfo.RefreshToken
		userInfo.LarkUserAccessTokenExpireAt = time.Unix(userTokenInfo.AccessTokenExpire, 0)
		userInfo.LarkUserRefreshTokenExpiredAt = time.Unix(userTokenInfo.RefreshTokenExpire, 0)
		if userData != nil {
			userInfo.LarkUserID = *userData.UserId
		}
		data, _ := json.Marshal(userData)
		userInfo.LarkUserInfo = string(data)
		err = User.SaveUser(ctx, &userInfo)
		if err != nil {
			log.Error(err)
			return model.User{}, err
		}
	}

	return userInfo, nil
}

func (p PluginService) ResetUserTokenExpired(ctx context.Context, meegoUserKey string) error {
	userInfo, err := User.GetUserInfoByMeegoUserKey(ctx, meegoUserKey)
	if err != nil {
		return err
	}

	userInfo.LarkUserRefreshTokenExpiredAt = time.Now().Add(-1 * time.Second)
	userInfo.LarkUserAccessTokenExpireAt = time.Now().Add(-1 * time.Second)
	err = User.SaveUser(ctx, &userInfo)
	return err
}

func (p PluginService) RefreshBind(ctx context.Context, workItemID int64) error {
	binds, err := dal.CalendarBind.MGetCalendarBindByWorkItemIDs(ctx, []int64{workItemID})
	if err != nil {
		log.Error(err)
		return err
	}
	if len(binds) == 0 {
		return nil
	}
	// 对每个都重新进行绑定
	var wg sync.WaitGroup
	for _, bind := range binds {
		if bind.Bind == false {
			continue
		}
		// 取用最后一个更新的人
		userInfo, errG := p.GetUserInfoByMeegoUserKey(ctx, bind.UpdateBy, true)
		if errG != nil {
			log.Errorf("[RefreshBind] GetUserInfo err, err: %v, userKey: %s", bind.UpdateBy)
			continue
		}

		wg.Add(1)

		param := BindCalendarParam{
			ProjectKey:      bind.ProjectKey,
			WorkItemTypeKey: bind.WorkItemTypeKey,
			WorkItemID:      bind.WorkItemID,
			CalendarEventID: bind.CalendarEventID,
		}
		go func() {
			defer wg.Done()
			count, errQ := dal.CalendarBind.CountMeetingByCalendarEventID(ctx, []string{param.CalendarEventID})
			if errQ != nil {
				log.Error(errQ)
				return
			}
			unbindMeetings, errQ := dal.VCMeetingUnBind.GetVCMeetingUnbindInfoByWorkItemID(ctx, workItemID)
			if errQ != nil {
				log.Error(errQ)
				return
			}
			meetingIDMap := make(map[string]struct{}, count)
			for _, um := range unbindMeetings {
				if um != nil {
					meetingIDMap[um.MeetingID] = struct{}{}
				}
			}
			// 清理下, 说明有脏数据, 这种数据就不用去绑定了
			if len(meetingIDMap) == int(count) && count != 0 {
				err = dal.CalendarBind.UnbindByCalendarEventIDAndWorkItemID(ctx, param.CalendarEventID, workItemID)
				if err != nil {
					log.Error(err)
				}
				err = dal.VCMeetingUnBind.DeleteMeetingsByWorkItemID(ctx, workItemID)
				if err != nil {
					log.Error(err)
				}
				return
			}

			err = p.BindCalendar(ctx, param, userInfo.LarkUserAccessToken, userInfo.MeegoUserKey)
			if err != nil {
				log.Error(err)
				return
			}
		}()
	}
	wg.Wait()
	return err
}

func (p PluginService) AutoBindCalendar(ctx context.Context, enable bool, projectKey, workItemTypeKey string, workItemID int64,
	meegoUserKey string) error {
	record, err := dal.JoinChatRecord.FirstByWorkItemID(ctx, workItemID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error(err)
			return err
		}
	}
	if record == nil && enable {
		err = p.tryJoinChatBycBindFirstCalendar(ctx, projectKey, workItemTypeKey, workItemID, meegoUserKey)
		if err != nil {
			log.Errorf("[PluginService.AutoBindCalendar] err tryJoinChatBycBindFirstCalendar, detail: %v", err)
			return err
		}
		// FIXME: 写后读场景, 不应该有 error, 但是需要读主库
		record, err = dal.JoinChatRecord.FirstByWorkItemID(ctx, workItemID)
		if err != nil {
			log.Error("[PluginService.AutoBindCalendar] err FirstByWorkItemID: %v", err)
			return err
		}
	}

	if record == nil {
		return nil
	}

	record.Enable = enable
	record.Operator = meegoUserKey
	err = dal.JoinChatRecord.Save(ctx, record)
	if err != nil {
		log.Errorf("[PluginService.AutoBindCalendar] error save join chat: %v", err)
		return err
	}

	return nil
}

// 尝试将机器人拉进群
func (p PluginService) tryJoinChatBycBindFirstCalendar(ctx context.Context, projectKey, workItemTypeKey string, workItemID int64,
	meegoUserKey string) error {
	record, err := dal.JoinChatRecord.FirstByWorkItemID(ctx, workItemID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error(err)
			return err
		}
	}
	if record != nil {
		return nil
	}
	resp, err := Meego.MeegoAPI.Chat.BotJoinChat(ctx, meego_api.BotJoinChatParam{
		ProjectKey:      projectKey,
		WorkItemTypeKey: workItemTypeKey,
		AppIDs:          []string{config.GetAPPConfig().LarkAppID},
		WorkItemID:      workItemID,
		MeegoUserKey:    meegoUserKey,
	})
	if err != nil {
		log.Errorf("[PluginService.tryJoinChatBycBindFirstCalendar] err BotJoinChat: %v", err)
		return errors.New(ErrTryBotJoinChat.Error() + err.Error())
	}
	record = &model.JoinChatRecord{
		WorkItemID:      workItemID,
		ProjectKey:      projectKey,
		WorkItemTypeKey: workItemTypeKey,
		Operator:        meegoUserKey,
		ChatID:          resp.ChatID,
		Enable:          true,
	}
	err = dal.JoinChatRecord.Save(ctx, record)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// 非常临时的逻辑, 不能长期使用
func (p PluginService) RefreshMeetingRecordTask(ctx context.Context, meetingIDs []string, userAccessToken string) error {
	noRecordMeetings, err := dal.CalendarBind.GetRealNoRecordMeetingByMeetingIDs(ctx, meetingIDs)
	if err != nil {
		log.Errorf("[RefreshMeetingRecordTask] GetRealNoRecordMeetingByMeetingIDs err, err: %v", err)
		return err
	}
	if len(noRecordMeetings) == 0 {
		log.Infof("[PluginService] RefreshMeetingRecordTask, no meetings record len 0, meetingIDs: %s", meetingIDs)
		return nil
	}
	var resultErr error
	for _, m := range noRecordMeetings {
		if m == nil {
			continue
		}
		recordInfo, err := lark_api.API.VChatAPI.GetMeetingRecord(ctx, m.MeetingID, userAccessToken)
		if err != nil {
			log.Errorf("[PluginService] RefreshMeetingRecordTask Get Meeting Record err, meetingID: %s, err: %v", m.MeetingID, err)
			resultErr = err
		}
		if recordInfo != nil && recordInfo.Recording != nil {
			if recordInfo.Recording.Url != nil && len(*recordInfo.Recording.Url) != 0 {
				m.RecordInfo = (*model.RecordInfo)(recordInfo.Recording)
				log.Infof("[PluginService] RefreshMeetingRecordTask start save, meetingID: %s", m.MeetingID)
				err = dal.CalendarBind.UpdateMeetingsRecordInfo(ctx, []*model.VCMeeting{m})
				if err != nil {
					resultErr = err
				}
				// 如果这段代码没有走到, 说明有问题, 所以最后的时候加上 Error
				continue
			}
		}
		resultErr = ErrNilMeetingRecord
	}
	return resultErr
}

func (p PluginService) RetryRefreshMeetingRecordTask(ctx context.Context, meetingIDs []string, userAccessToken string) error {
	err := retry.Do(func() error {
		err := p.RefreshMeetingRecordTask(ctx, meetingIDs, userAccessToken)
		if err != nil {
			log.Errorf("RefreshMeetingRecordTask err, err: %v", err)
			return err
		}
		return nil
	}, retry.Delay(time.Second*6), retry.Attempts(3))
	if err != nil {
		log.Errorf("RetryRefreshMeetingRecordTask finally error, err: %v", err)
		return err
	}
	return nil
}

func (p PluginService) GetUserInfoByBinds(ctx context.Context, binds []*model.CalendarBind) (map[string]*user.UserBasicInfo, error) {
	userKeysMap := make(map[string]struct{})
	for _, bind := range binds {
		if bind != nil {
			userKeysMap[bind.UpdateBy] = struct{}{}
		}
	}
	if len(userKeysMap) == 0 {
		return map[string]*user.UserBasicInfo{}, nil
	}
	userKeys := lo.Keys(userKeysMap)
	userInfos, err := meego_api.API.User.GetUserInfo(ctx, userKeys)
	if err != nil {
		log.Errorf("[PluginService] GetUserInfoByBinds, userKeys: %s, err: %v", userKeys, err)
		return nil, err
	}
	result := make(map[string]*user.UserBasicInfo, len(userInfos))
	for _, u := range userInfos {
		if u != nil {
			result[u.UserKey] = u
		}
	}
	return result, nil
}

// 给用户发送消息尝试让对应用户给当前应用授权
func (p PluginService) SendMsgForACL(ctx context.Context, meegoUserKey string) error {
	// 查询 lark 用户信息
	meegoUsers, err := meego_api.API.User.GetUserInfo(ctx, []string{meegoUserKey})
	if err != nil {
		log.Errorf("[PluginService.SendMsgForACL] GetUserInfoByMeegoUserKey err, meegoUserKey: %s, err: %v", meegoUserKey, err)
		return err
	}

	if len(meegoUsers) == 0 {
		log.Infof("[PluginService.SendMsgForACL] GetUserInfoByMeegoUserKey, meegoUserKey: %s, not found", meegoUserKey)
		return nil
	}
	// 私聊发送消息
	larkUnionID := meegoUsers[0].OutID

	_, err = lark_api.API.IMAPI.CreateTextMessage(ctx, lark_api.LarkUserIDTypeUnionID, larkUnionID, p.genACLMessage(meegoUserKey))
	if err != nil {
		log.Errorf("[PluginService.SendMsgForACL] CreateTextMessage err, meegoUserKey: %s, err: %v", meegoUserKey, err)
		return err
	}
	return nil
}

func (p PluginService) genACLMessage(meegoUserKey string) string {
	redirectUrl := fmt.Sprintf("%s/api/v1/meego/lark/auth", config.GetAPPConfig().DomainURL)
	stateUrl := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("https://project.feishu.cn?meego_user_key=%s", meegoUserKey)))
	aclUrl := fmt.Sprintf("%s?redirect_uri=%s&app_id=%s&state=%s",
		LarkAuthURL,
		redirectUrl,
		config.GetAPPConfig().LarkAppID,
		stateUrl)
	text := fmt.Sprintf(`您需要给当前应用授权, 才能正常使用「会议管理」插件相关功能，<b>点击链接进行应用授权</b>：[授权链接](%s)`, aclUrl)
	return text
}
