package dal

import (
	"context"

	"meego_meeting_plugin/model"

	"github.com/gofiber/fiber/v2/log"
)

type PendingTaskDao struct{}

func NewPendingTaskDao() PendingTaskDao {
	return PendingTaskDao{}
}

func (d PendingTaskDao) Create(ctx context.Context, task *model.PendingTask) error {
	if task == nil {
		return nil
	}
	err := db.WithContext(ctx).Create(task).Error
	if err != nil {
		log.Errorf("failed to create pending task: %v", err)
	}
	return err
}

func (d PendingTaskDao) GetUnprocessedTasks(ctx context.Context) ([]model.PendingTask, error) {
	var tasks []model.PendingTask
	err := db.WithContext(ctx).
		Where("status != ?", model.TaskStatusProcessed).
		Order("created_at asc").
		Find(&tasks).Error
	if err != nil {
		log.Errorf("failed to get pending tasks: %v", err)
		return nil, err
	}
	return tasks, nil
}

// 获取Meego用户相关的所有 pending task
func (d PendingTaskDao) GetUnprocessedTasksByMeegoUserKey(ctx context.Context, meegoUserKey string) ([]model.PendingTask, error) {
	var tasks []model.PendingTask
	err := db.WithContext(ctx).
		Where("meego_user_key = ? AND status != ?", meegoUserKey, model.TaskStatusProcessed).
		Order("created_at asc").
		Find(&tasks).Error
	if err != nil {
		log.Errorf("failed to get pending tasks for meego user %s: %v", meegoUserKey, err)
		return nil, err
	}
	return tasks, nil
}

func (d PendingTaskDao) Update(ctx context.Context, task *model.PendingTask) error {
	if task == nil || task.ID == 0 {
		return nil
	}
	err := db.WithContext(ctx).Save(task).Error
	if err != nil {
		log.Errorf("failed to update pending task: %v", err)
	}
	return err
}
