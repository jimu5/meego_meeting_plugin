package meego_api

import (
	sdk "github.com/larksuite/project-oapi-sdk-golang"
	"github.com/larksuite/project-oapi-sdk-golang/core"
	"meego_meeting_plugin/config"
	"net/http"
	"time"
)

var own_client = initOwnClient()

func initOwnClient() *sdk.Client {
	c := sdk.NewClient(config.MeegoPluginID, config.MeegoPluginSecret, sdk.WithLogLevel(core.LogLevelDebug),
		sdk.WithReqTimeout(3*time.Second),
		sdk.WithEnableTokenCache(true),
		sdk.WithHttpClient(http.DefaultClient),
	)
	return c
}
