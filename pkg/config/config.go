package config

import (
	"parchment-server/internal/models/config"
	"sync"

	"parchment-server/pkg/global"

	"github.com/spf13/viper"
)

var (
	cfg  *config.Config
	v    *viper.Viper
	once sync.Once
)

// InitConfig 初始化配置并设置全局变量
func InitConfig() *config.Config {
	once.Do(func() {
		v = viper.New()
		v.SetConfigName("dev")
		v.SetConfigType("yaml")
		v.AddConfigPath("./config")

		if err := v.ReadInConfig(); err != nil {
			panic("Error reading config file: " + err.Error())
		}

		cfg = &config.Config{
			// 默认日志配置
			Log: config.LogConfig{
				Level:       "info",
				Format:      "console",
				OutputPaths: []string{"stdout"},
				Caller:      true,
				Stacktrace:  false,
			},
		}

		if err := v.Unmarshal(cfg); err != nil {
			panic("Unable to decode into struct: " + err.Error())
		}

		// 设置全局配置变量
		global.Config = v
	})
	return cfg
}

// GetConfig 获取配置
func GetConfig() *config.Config {
	if cfg == nil {
		InitConfig()
	}
	return cfg
}
