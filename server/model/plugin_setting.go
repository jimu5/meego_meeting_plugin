package model

import "gorm.io/gorm"

// 暂不使用
type PluginSetting struct {
	gorm.Model
	ProjectKey string `json:"project_key" gorm:"size:100"` // 空间 key
	AppID      string `json:"app_id" gorm:"size:100"`      // 飞书自建应用 app id
	AppSecret  string `json:"app_secret" gorm:"size:100"`  // 飞书自建应用 app secret
}
