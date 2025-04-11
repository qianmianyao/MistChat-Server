package global

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	// DB Config 全局配置
	DB *gorm.DB

	// Logger 全局日志
	Logger *zap.Logger

	// Config 全局配置
	Config *viper.Viper
)
