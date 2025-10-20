package config

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2/log"
	"gopkg.in/yaml.v2"
)

type Configure struct {
	Database  *Database  `yaml:"database"`
	APPConfig *APPConfig `yaml:"app_config"`
}

func (c *Configure) Check() error {
	if c.APPConfig == nil {
		return fmt.Errorf("error app config is nil")
	}
	if err := c.APPConfig.Check(); err != nil {
		return err
	}

	if c.Database == nil {
		c.Database = &Database{}
	}
	if len(c.Database.DBName) == 0 {
		c.Database.DBName = "plugin"
	}

	return nil
}

type Database struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
}

type APPConfig struct {
	// 飞书开放平台配置
	LarkAppID     string `yaml:"lark_app_id"`
	LarkAppSecret string `yaml:"lark_app_secret"`

	// Meego 开放平台配置
	MeegoPluginID           string `yaml:"meego_plugin_id"`
	MeegoPluginSecret       string `yaml:"meego_plugin_secret"`
	MeegoEventCallbackToken string `yaml:"meego_event_callback_token"`
	DomainURL               string `yaml:"domain_url"`
}

func (a *APPConfig) Check() error {
	switch {
	case len(a.LarkAppID) == 0:
		return fmt.Errorf("lark_app_id is empty")
	case len(a.LarkAppSecret) == 0:
		return fmt.Errorf("lark_app_secret is empty")
	case len(a.MeegoPluginID) == 0:
		return fmt.Errorf("meego_plugin_id is empty")
	case len(a.MeegoPluginSecret) == 0:
		return fmt.Errorf("meego_plugin_secret is empty")
	case len(a.MeegoEventCallbackToken) == 0:
		return fmt.Errorf("meego_event_callback_token is empty")
	case len(a.DomainURL) == 0:
		return fmt.Errorf("domain_url is empty")
	default:

	}

	return nil
}

var Config Configure

func GetAPPConfig() *APPConfig {
	return Config.APPConfig
}

func readYAMLConfig() {
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Warn("not found config.yaml")
		return
	}
	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		panic(fmt.Sprintf("config.yaml file content error: %v", err))
	}
}
