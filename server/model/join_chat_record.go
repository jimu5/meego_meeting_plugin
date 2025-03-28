package model

import "gorm.io/gorm"

type JoinChatRecord struct {
	gorm.Model
	ProjectKey      string `json:"project_key" gorm:"size:256"`
	WorkItemTypeKey string `json:"work_item_type_key" gorm:"size:256"`
	WorkItemID      int64  `json:"work_item_id,omitempty" gorm:"uniqueIndex:uniq_wi"`
	Operator        string `json:"operator,omitempty" gorm:"size:256"`
	ChatID          string `json:"chat_id" gorm:"size:256"`
	Enable          bool   `json:"enable"`
}
