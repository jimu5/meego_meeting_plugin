package model

import "gorm.io/gorm"

type JoinChatRecord struct {
	gorm.Model
	ProjectKey      string `json:"project_key"`
	WorkItemTypeKey string `json:"work_item_type_key"`
	WorkItemID      int64  `json:"work_item_id,omitempty" gorm:"uniqueIndex:uniq_wi"`
	Operator        string `json:"operator,omitempty"`
	ChatID          string `json:"chat_id"`
	Enable          bool   `json:"enable"`
}
