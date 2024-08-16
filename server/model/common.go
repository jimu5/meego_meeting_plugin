package model

import "gorm.io/gorm"

type BaseModel struct {
	gorm.Model
	CreateBy string `json:"create_by" gorm:"size:100"` // 创建人 userKey
	UpdateBy string `json:"update_by" gorm:"size:100"` // 更新人 userKey
}
