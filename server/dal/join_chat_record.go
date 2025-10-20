package dal

import (
	"context"
	"errors"

	"meego_meeting_plugin/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type JoinChatRecordDao struct {
}

func NewJoinChatRecordDao() JoinChatRecordDao {
	return JoinChatRecordDao{}
}

func (j JoinChatRecordDao) Save(ctx context.Context, record *model.JoinChatRecord) error {
	if record == nil {
		return nil
	}
	err := db.WithContext(ctx).Save(record).Error
	if err != nil {
		log.Error(err)
	}
	return err
}

func (j JoinChatRecordDao) CreateOrUpdate(ctx context.Context, record *model.JoinChatRecord) error {
	if record == nil {
		return nil
	}
	err := db.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "work_item_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"operator", "enable"}),
		}).Create(record).Error
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (j JoinChatRecordDao) FirstByWorkItemID(ctx context.Context, workItemID int64) (*model.JoinChatRecord, error) {
	result := model.JoinChatRecord{}
	err := db.WithContext(ctx).Where("work_item_id = ?", workItemID).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error(err)
		return nil, err
	}
	return &result, err
}

func (j JoinChatRecordDao) FirstByChatID(ctx context.Context, chatID string) (*model.JoinChatRecord, error) {
	result := model.JoinChatRecord{}
	err := db.WithContext(ctx).Where("chat_id = ?", chatID).First(&result).Error
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &result, err
}
