package config

import "flag"

// 飞书开放平台配置
var (
	LarkAppID     string
	LarkAppSecret string
)

// Meego 开放平台配置
var (
	MeegoPluginID     string
	MeegoPluginSecret string
)

func InitConfig() {
	larkAppID := flag.String("lark_app_id", "", "lark app_id")
	larkAppSecret := flag.String("lark_app_secret", "", "lark app_secret")

	meegoPluginID := flag.String("meego_plugin_id", "", "meego_plugin_id")
	meegoPluginSecret := flag.String("meego_plugin_secret", "", "meego_plugin_secret")

	if larkAppID == nil || larkAppSecret == nil || meegoPluginID == nil || meegoPluginSecret == nil {
		panic("lark_app_id or lark_app_secret or meego_plugin_id or meego_plugin_secret not set")
	}

	LarkAppID = *larkAppID
	LarkAppSecret = *larkAppSecret
	MeegoPluginID = *meegoPluginID
	MeegoPluginSecret = *meegoPluginSecret

	if len(LarkAppID) == 0 || len(LarkAppSecret) == 0 || len(MeegoPluginID) == 0 || len(MeegoPluginSecret) == 0 {
		panic("lark_app_id or lark_app_secret or meego_plugin_id or meego_plugin_secret is empty")
	}

	return
}
