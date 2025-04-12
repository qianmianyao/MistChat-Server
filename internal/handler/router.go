package handler

import (
	"github.com/qianmianyao/parchment-server/internal/handler/hello"
	"github.com/qianmianyao/parchment-server/internal/websocket"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由组
func SetupRouter(r *gin.Engine) {

	v1 := r.Group("/api/v1")
	{
		// hello接口
		v1.GET("/example/hello_world", hello.Hello)
		wsGroup := v1.Group("/ws")
		{
			websocket.RegisterWebSocketRoutes(wsGroup)
		}
	}
}
