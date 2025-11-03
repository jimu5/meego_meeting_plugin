package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"meego_meeting_plugin/dal"
	"meego_meeting_plugin/model"

	"github.com/gofiber/fiber/v2/log"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

var Cron CronService

type CronService struct {
	scheduler *cron.Cron
}

func NewCronService() CronService {
	c := cron.New(cron.WithSeconds())
	return CronService{
		scheduler: c,
	}
}

func (c *CronService) Start() {
	_, err := c.scheduler.AddFunc("@every 1d", c.processPendingTasks)
	if err != nil {
		log.Errorf("Failed to schedule pending task processing job: %v", err)
		return
	}

	c.scheduler.Start()
	log.Info("Cron service started successfully")
}

func (c *CronService) Stop() {
	if c.scheduler != nil {
		c.scheduler.Stop()
		log.Info("Cron service stopped")
	}
}

func (c *CronService) processPendingTasks() {
	ctx := context.Background()
	log.Info("Starting pending task processing...")

	tasks, err := dal.PendingTask.GetUnprocessedTasks(ctx)
	if err != nil {
		log.Errorf("Failed to query pending tasks: %v", err)
		return
	}

	if len(tasks) > 0 {
		log.Infof("Found %d pending tasks", len(tasks))
	}

	for _, task := range tasks {
		user, err := Plugin.GetUserInfoByMeegoUserKey(ctx, task.MeegoUserKey, false)
		if err != nil || user.ID == 0 {
			log.Infof("User %s for task %d not ready, skipping", task.MeegoUserKey, task.ID)
			continue
		}

		if err := c.processSingleTask(ctx, task, user); err != nil {
			log.Errorf("Failed to process task %d for user %s: %v", task.ID, task.MeegoUserKey, err)
			task.Status = model.TaskStatusFailed
			task.Remark = err.Error()
		} else {
			task.Status = model.TaskStatusProcessed
		}

		task.RetryCount++
		if err := dal.PendingTask.Update(ctx, &task); err != nil {
			log.Errorf("Failed to update task status for task %d: %v", task.ID, err)
		}
	}
}

// 通过 lark 用户来触发这个用户相关的任务执行
func (c *CronService) ProcessTasksByLarkUserInfo(ctx context.Context, larkUserInfo larkim.UserId) error {
	var larkUserID string
	if larkUserInfo.UserId != nil {
		larkUserID = *larkUserInfo.UserId
	}
	// 先通过 lark user id 来查询到 meego user key
	meegoUserKey, err := User.GetMeegoUserKeyByLarkUserInfo(ctx, larkUserInfo)
	if err != nil || meegoUserKey == "" {
		return fmt.Errorf("user %s not ready, skipping", larkUserID)
	}
	// 获取这个用户的相关信息
	user, err := Plugin.GetUserInfoByMeegoUserKey(ctx, meegoUserKey, false)
	if err != nil || user.ID == 0 {
		return fmt.Errorf("user %s not ready, skipping", larkUserID)
	}

	// 再通过 meego user key 来查询到所有相关的任务
	tasks, err := dal.PendingTask.GetUnprocessedTasksByMeegoUserKey(ctx, meegoUserKey)
	if err != nil {
		return fmt.Errorf("failed to query pending tasks for lark user %s: %w", larkUserID, err)
	}
	if len(tasks) > 0 {
		log.Infof("Found %d pending tasks for lark user %s", len(tasks), larkUserID)
	}
	for _, t := range tasks {
		if err := c.processSingleTask(ctx, t, user); err != nil {
			log.Errorf("Failed to process task %d for user %s: %v", t.ID, larkUserID, err)
			t.Status = model.TaskStatusFailed
			t.Remark = err.Error()
		} else {
			t.Status = model.TaskStatusProcessed
		}
		t.RetryCount++
		if err := dal.PendingTask.Update(ctx, &t); err != nil {
			log.Errorf("Failed to update task status for task %d: %v", t.ID, err)
		}
	}
	return nil
}

// 通过 meego 用户信息来触发这个任务的执行
func (c *CronService) ProcessTasksByUser(ctx context.Context, user model.User) error {
	if user.MeegoUserKey == "" {
		return fmt.Errorf("user %s not ready, skipping", user.LarkUserID)
	}
	// 再通过 meego user key 来查询到所有相关的任务
	tasks, err := dal.PendingTask.GetUnprocessedTasksByMeegoUserKey(ctx, user.MeegoUserKey)
	if err != nil {
		return fmt.Errorf("failed to query pending tasks for lark user %s: %w", user.MeegoUserKey, err)
	}
	if len(tasks) > 0 {
		log.Infof("Found %d pending tasks for lark user %s", len(tasks), user.MeegoUserKey)
	}
	for _, t := range tasks {
		if err := c.processSingleTask(ctx, t, user); err != nil {
			log.Errorf("Failed to process task %d for user %s: %v", t.ID, user.MeegoUserKey, err)
			t.Status = model.TaskStatusFailed
			t.Remark = err.Error()
		} else {
			t.Status = model.TaskStatusProcessed
		}
		t.RetryCount++
		if err := dal.PendingTask.Update(ctx, &t); err != nil {
			log.Errorf("Failed to update task status for task %d: %v", t.ID, err)
		}
	}
	return nil
}

// 处理单个任务的执行
func (c *CronService) processSingleTask(ctx context.Context, task model.PendingTask, user model.User) error {
	log.Infof("Processing task %d of type %s for user: %s", task.ID, task.TaskType, user.MeegoUserKey)

	switch task.TaskType {
	case TaskTypeBindCalendar:
		var param BindCalendarParam
		if err := json.Unmarshal([]byte(task.Payload), &param); err != nil {
			return fmt.Errorf("failed to unmarshal payload for task %d: %w", task.ID, err)
		}
		// The BindCalendar function is what we need to call.
		// We need to decide if we should call BindCalendar or RefreshBind.
		// Since the task was created from RefreshBind, let's stick to the logic within it.
		// A simpler approach for now is to call BindCalendar directly.
		return Plugin.BindCalendar(ctx, param, user.LarkUserAccessToken, user.MeegoUserKey)
	case TaskTypeHandleMeetingBindByUserKey:
		var param HandleMeetingBindByUserKeyParam
		if err := json.Unmarshal([]byte(task.Payload), &param); err != nil {
			return fmt.Errorf("failed to unmarshal payload for task %d: %w", task.ID, err)
		}
		// 重新查询下 record
		record, err := dal.JoinChatRecord.FirstByChatID(ctx, param.Record.ChatID)
		if err != nil {
			// 如果是没有找到 error 则直接返回
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return fmt.Errorf("failed to get join chat record %s: %w", param.Record.ChatID, err)
		}
		param.Record = record
		return Plugin.HandleMeetingBindByUserKey(ctx, param)
	default:
		return fmt.Errorf("unknown task type: %s", task.TaskType)
	}
}
