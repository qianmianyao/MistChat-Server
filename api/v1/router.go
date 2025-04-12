package api

import (
	"github.com/gin-gonic/gin"
	"github.com/qianmianyao/parchment-server/internal/handler/hello"
	"github.com/qianmianyao/parchment-server/internal/websocket"
)

// SetupRouter 设置路由组
func SetupRouter(r *gin.Engine) {

	v1 := r.Group("/api/v1")
	{
		v1.GET("/example/hello_world", hello.Hello)
		wsGroup := v1.Group("/chat")
		{
			RegisterWebSocketRoutes(wsGroup)
		}
	}
}

// RegisterWebSocketRoutes registers the WebSocket routes with the provided Gin engine.
func RegisterWebSocketRoutes(r *gin.RouterGroup) {
	hub := websocket.NewHub()
	go hub.Run()
	r.GET("/connect", websocket.NewWebSockerRouter().WsHandler(hub))
	r.GET("/check_room_password", websocket.NewWebSockerRouter().CheckRoomPasswordRequired)
	r.GET("/join_room", websocket.NewWebSockerRouter().JoinRoom)
}
