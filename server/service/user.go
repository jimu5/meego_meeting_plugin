package service

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"meego_meeting_plugin/dal"
	"meego_meeting_plugin/model"
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

func (u UserService) SaveUser(ctx context.Context, user *model.User) error {
	log.Info("save user")
	return dal.User.Save(ctx, user)
}
