package meego_api

import (
	"net/http"
	"time"

	"meego_meeting_plugin/config"

	sdk "github.com/larksuite/project-oapi-sdk-golang"
	"github.com/larksuite/project-oapi-sdk-golang/core"
)

var own_client *sdk.Client

func InitOwnClient() *sdk.Client {
	c := sdk.NewClient(config.GetAPPConfig().MeegoPluginID, config.GetAPPConfig().MeegoPluginSecret, sdk.WithLogLevel(core.LogLevelDebug),
		sdk.WithReqTimeout(3*time.Second),
		sdk.WithEnableTokenCache(true),
		sdk.WithHttpClient(http.DefaultClient),
	)

	own_client = c

	return c
}
