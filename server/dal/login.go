package dal

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"meego_meeting_plugin/model"
)

type UserDao struct {
}

func NewUserDao() UserDao {
	return UserDao{}
}

func (u UserDao) Save(ctx context.Context, user *model.User) error {
	if user == nil {
		return nil
	}
	err := db.WithContext(ctx).Save(user).Error
	if err != nil {
		log.Error(err)
	}
	return err
}

func (u UserDao) QueryByMeegoUserKey(ctx context.Context, meegoUserKey string) (*model.User, error) {
	// 不够严谨
	if len(meegoUserKey) == 0 {
		return nil, nil
	}
	result := &model.User{}
	err := db.WithContext(ctx).Where("meego_user_key =?", meegoUserKey).First(result).Error
	if err != nil {
		log.Error(err)
	}
	return result, err
}

func (u UserDao) QueryByLarkUserID(ctx context.Context, larkUserID string) (*model.User, error) {
	// 不够严谨
	if len(larkUserID) == 0 {
		return nil, nil
	}
	result := &model.User{}
	err := db.WithContext(ctx).Where("lark_user_id = ?", larkUserID).First(result).Error
	if err != nil {
		log.Error(err)
	}
	return result, err
}
