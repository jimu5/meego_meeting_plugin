package model

import (
	"time"

	"gorm.io/gorm"
)

// 目前这个模型并不是很安全
type User struct {
	gorm.Model
	MeegoUserKey                  string    `json:"meego_user_key,omitempty" gorm:"size:256"`
	LarkUserAccessToken           string    `json:"lark_user_access_token,omitempty" gorm:"size:256"`
	LarkUserRefreshToken          string    `json:"lark_user_refresh_token,omitempty" gorm:"size:256"`
	LarkUserAccessTokenExpireAt   time.Time `json:"lark_user_access_token_expire_at"`
	LarkUserRefreshTokenExpiredAt time.Time `json:"lark_user_refresh_token_expired_at"`
	LarkUserID                    string    `json:"lark_user_id" gorm:"size:256"`
	LarkUserInfo                  string    `json:"lark_user_info,omitempty" gorm:"type:text"`
}
