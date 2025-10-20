package config

import "flag"

func InitConfig() {
	readYAMLConfig()

	if Config.APPConfig == nil {
		Config.APPConfig = &APPConfig{}
	}

	larkAppID := flag.String("lark_app_id", "", "lark app_id")
	larkAppSecret := flag.String("lark_app_secret", "", "lark app_secret")
	meegoPluginID := flag.String("meego_plugin_id", "", "meego_plugin_id")
	meegoPluginSecret := flag.String("meego_plugin_secret", "", "meego_plugin_secret")
	meegoEventCallbackToken := flag.String("meego_event_callback_token", "", "meego_event_callback_token")
	domainURL := flag.String("domain_url", "", "domain_url")

	flag.Parse()

	if larkAppID != nil && len(*larkAppID) != 0 {
		Config.APPConfig.LarkAppID = *larkAppID
	}
	if larkAppSecret != nil && len(*larkAppSecret) != 0 {
		Config.APPConfig.LarkAppSecret = *larkAppSecret
	}
	if meegoPluginID != nil && len(*meegoPluginID) != 0 {
		Config.APPConfig.MeegoPluginID = *meegoPluginID
	}
	if meegoPluginSecret != nil && len(*meegoPluginSecret) != 0 {
		Config.APPConfig.MeegoPluginSecret = *meegoPluginSecret
	}
	if meegoEventCallbackToken != nil && len(*meegoEventCallbackToken) != 0 {
		Config.APPConfig.MeegoEventCallbackToken = *meegoEventCallbackToken
	}
	if domainURL != nil && len(*domainURL) != 0 {
		Config.APPConfig.DomainURL = *domainURL
	}

	if err := Config.Check(); err != nil {
		panic(err)
	}

	return
}
