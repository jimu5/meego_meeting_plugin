package service

import (
	"context"
	"meego_meeting_plugin/dal"
	"meego_meeting_plugin/model"
	"meego_meeting_plugin/service/meego_api"

	"github.com/gofiber/fiber/v2/log"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

var User = UserService{}

type UserService struct {
}

func (u UserService) GetUserInfoByMeegoUserKey(ctx context.Context, meegoUserKey string) (model.User, error) {
	user, err := dal.User.QueryByMeegoUserKey(ctx, meegoUserKey)
	if err != nil {
		log.Error(err)
		return model.User{}, err
	}
	if user == nil {
		return model.User{}, ErrNilUser
	}
	return *user, nil
}

func (u UserService) GetUserInfoByLarkUserID(ctx context.Context, larkUserID string) (model.User, error) {
	user, err := dal.User.QueryByLarkUserID(ctx, larkUserID)
	if err != nil {
		log.Error(err)
		return model.User{}, err
	}
	if user == nil {
		return model.User{}, ErrNilUser
	}
	return *user, nil
}

// 通过 lark user id 获取 meego_user_key
func (u UserService) GetMeegoUserKeyByLarkUserInfo(ctx context.Context, larkUserInfo larkim.UserId) (string, error) {
	var larkUserKey string
	if larkUserInfo.UserId != nil {
		larkUserKey = *larkUserInfo.UserId
	}
	var meegoUserKey string
	if len(larkUserKey) != 0 {
		larkUser, errU := User.GetUserInfoByLarkUserID(ctx, larkUserKey)
		if errU != nil {
			log.Infof("cant find lark user related meego user, lark user id: %s", larkUserKey)
			// 尝试使用 api 来获取用户
			var larkUnionID string
			if larkUserInfo.UnionId != nil {
				larkUnionID = *larkUserInfo.UnionId
			}
			meegoUserInfos, err := meego_api.API.User.GetUserInfoByLarkUnionID(ctx, []string{larkUnionID})
			if err != nil {
				log.Errorf("[handleChatCalendarMessage] GetUserInfoByLarkUnionID err, unionID: %s, err: %v", larkUnionID, err)
				return meegoUserKey, err
			}
			if len(meegoUserInfos) > 0 {
				meegoUserKey = meegoUserInfos[0].UserKey
			}

		} else {
			meegoUserKey = larkUser.MeegoUserKey
		}
	}
	return meegoUserKey, nil
}

func (u UserService) SaveUser(ctx context.Context, user *model.User) error {
	log.Info("save user")
	return dal.User.Save(ctx, user)
}
