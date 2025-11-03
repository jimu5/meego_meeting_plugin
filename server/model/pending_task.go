package model

import (
	"gorm.io/gorm"
)

type TaskStatus int

const (
	TaskStatusPending   TaskStatus = 1
	TaskStatusProcessed TaskStatus = 2
	TaskStatusFailed    TaskStatus = 3
)

type PendingTask struct {
	gorm.Model
	MeegoUserKey string     `json:"meego_user_key" gorm:"size:256;index"`
	TaskType     string     `json:"task_type" gorm:"size:100"`
	Payload      string     `json:"payload" gorm:"type:text"`
	Status       TaskStatus `json:"status" gorm:"default:1"`
	RetryCount   int        `json:"retry_count" gorm:"default:0"`
	Remark       string     `json:"remark" gorm:"type:text"`
}
