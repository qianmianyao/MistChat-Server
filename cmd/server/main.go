package main

import (
	"fmt"

	"parchment-server/internal/handler"
	"parchment-server/pkg/config"
	"parchment-server/pkg/database"
	"parchment-server/pkg/global"
	"parchment-server/pkg/logger"

	// 导入swagger文档
	_ "parchment-server/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Parchment API
// @version 1.0
// @description Parchment服务器API文档
// @host localhost:8080
// @BasePath /api/v1

// 初始化所有全局组件
func initComponents() {
	// 初始化配置（必须第一个初始化）
	config.InitConfig()

	// 初始化日志
	logger.InitLogger()

	// 初始化数据库
	database.InitDB()

	// 所有组件都已初始化，现在可以通过 global 包访问
	global.Logger.Info("所有组件初始化完成")
}

func main() {
	// 初始化所有组件
	initComponents()
	router := gin.Default()

	// Swagger文档路由
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	handler.SetupRouter(router) // 设置路由组

	err := router.Run()
	if err != nil {
		return
	}

	// 模拟启动服务器
	fmt.Println("服务器已启动...")
}
