package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/qianmianyao/parchment-server/pkg/config"
	"github.com/qianmianyao/parchment-server/pkg/global"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	once   sync.Once
)

// InitLogger 初始化日志并设置全局变量
func InitLogger() *zap.Logger {
	once.Do(func() {
		initLoggerImpl()
	})
	return logger
}

// initLoggerImpl 初始化日志的具体实现
func initLoggerImpl() {
	// 从配置中获取日志配置
	logConfig := config.GetConfig().Log

	// 设置日志级别
	level := zap.InfoLevel
	switch logConfig.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "fatal":
		level = zap.FatalLevel
	}

	// 创建自定义编码器配置，适合控制台输出
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 区分控制台和文件输出的编码配置
	var core zapcore.Core
	if logConfig.Format == "console" {
		// 控制台输出
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		consoleOutput := zapcore.AddSync(os.Stdout)
		consoleLevel := zap.NewAtomicLevelAt(level)
		consoleCore := zapcore.NewCore(consoleEncoder, consoleOutput, consoleLevel)

		// 如果有文件输出，创建文件输出的Core
		var fileCores []zapcore.Core // 存储所有文件输出的Core

		for _, path := range logConfig.OutputPaths {
			if path != "stdout" && path != "stderr" {
				// 文件输出使用JSON格式，没有颜色代码
				fileEncoderConfig := encoderConfig
				fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
				fileEncoder := zapcore.NewJSONEncoder(fileEncoderConfig)

				// 确保日志目录存在
				dir := filepath.Dir(path)
				if _, err := os.Stat(dir); os.IsNotExist(err) {
					err = os.MkdirAll(dir, 0755)
					if err != nil {
						fmt.Printf("无法创建日志目录 %s: %v\n", dir, err)
						continue
					}
					fmt.Printf("已创建日志目录: %s\n", dir)
				}

				// 打开日志文件
				fileOutput, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Printf("无法打开日志文件 %s: %v\n", path, err)
					continue
				}

				// 创建文件输出的Core并添加到集合
				fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(fileOutput), level)
				fileCores = append(fileCores, fileCore)
			}
		}

		// 组合所有Core
		if len(fileCores) > 0 {
			// 将控制台输出Core和所有文件输出Core组合在一起
			allCores := append([]zapcore.Core{consoleCore}, fileCores...)
			core = zapcore.NewTee(allCores...)
		} else {
			core = consoleCore
			fmt.Println("警告: 没有配置文件日志输出或全部配置失败，只使用控制台输出")
		}
	} else {
		// 使用标准配置创建Core
		zapConfig := zap.Config{
			Level:            zap.NewAtomicLevelAt(level),
			Development:      false,
			Sampling:         &zap.SamplingConfig{Initial: 100, Thereafter: 100},
			Encoding:         logConfig.Format,
			EncoderConfig:    encoderConfig,
			OutputPaths:      logConfig.OutputPaths,
			ErrorOutputPaths: []string{"stderr"},
		}

		var err error
		logger, err = zapConfig.Build()
		if err != nil {
			panic("Failed to initialize logger: " + err.Error())
		}

		// 设置全局Logger变量
		global.Logger = logger

		return
	}

	// 创建日志器实例
	logger = zap.New(core)

	// 添加选项
	var opts []zap.Option

	// 添加调用者信息
	if logConfig.Caller {
		opts = append(opts, zap.AddCaller())
	}

	// 添加堆栈跟踪
	if logConfig.Stacktrace {
		opts = append(opts, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	if len(opts) > 0 {
		logger = logger.WithOptions(opts...)
	}

	// 设置全局Logger变量
	global.Logger = logger

	// 替换全局 logger
	zap.ReplaceGlobals(logger)

	logger.Info("Logger initialized successfully")
}
