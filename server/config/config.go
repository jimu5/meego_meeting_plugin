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

	flag.Parse()

	if larkAppID == nil || larkAppSecret == nil || meegoPluginID == nil || meegoPluginSecret == nil {
		panic("lark_app_id or lark_app_secret or meego_plugin_id or meego_plugin_secret not set")
	}

	switch {
	case larkAppID == nil:
		panic("lark_app_id is required")
	case larkAppSecret == nil:
		panic("lark_app_secret is required")
	case meegoPluginID == nil:
		panic("meego_plugin_id is required")
	case meegoPluginSecret == nil:
		panic("meego_plugin_secret is required")
	}

	LarkAppID = *larkAppID
	LarkAppSecret = *larkAppSecret
	MeegoPluginID = *meegoPluginID
	MeegoPluginSecret = *meegoPluginSecret

	switch {
	case len(LarkAppID) == 0:
		panic("lark_app_id is empty")
	case len(LarkAppSecret) == 0:
		panic("lark_app_secret is empty")
	case len(MeegoPluginID) == 0:
		panic("meego_plugin_id is empty")
	case len(MeegoPluginSecret) == 0:
		panic("meego_plugin_secret is empty")
	}

	return
}
