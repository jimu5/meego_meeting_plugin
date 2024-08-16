package model

import (
	"gorm.io/gorm"
	"time"
)

// 目前这个模型并不是很安全
type User struct {
	gorm.Model
	MeegoUserKey                  string    `json:"meego_user_key,omitempty"`
	LarkUserAccessToken           string    `json:"lark_user_access_token,omitempty"`
	LarkUserRefreshToken          string    `json:"lark_user_refresh_token,omitempty"`
	LarkUserAccessTokenExpireAt   time.Time `json:"lark_user_access_token_expire_at"`
	LarkUserRefreshTokenExpiredAt time.Time `json:"lark_user_refresh_token_expired_at"`
	LarkUserID                    string    `json:"lark_user_id"`
	LarkUserInfo                  string    `json:"lark_user_info,omitempty"`
}
