package api

import (
	"github.com/gin-gonic/gin"
	"qianmianyao/MistChat-Server/internal/handler/chat"
	"qianmianyao/MistChat-Server/internal/handler/hello"
	"qianmianyao/MistChat-Server/internal/websocket"
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

// RegisterWebSocketRoutes 使用提供的Gin引擎注册WebSocket路由
func RegisterWebSocketRoutes(r *gin.RouterGroup) {
	hub := websocket.NewHub()
	go hub.Run()
	r.POST("/register", chat.NewWebSockerRouter().Register)
	r.GET("/connect", chat.NewWebSockerRouter().WsHandler(hub))
	r.POST("/check_room_password", chat.NewWebSockerRouter().CheckRoomPasswordRequired)
	r.POST("/join_room", chat.NewWebSockerRouter().JoinRoom)
	r.POST("/create_room", chat.NewWebSockerRouter().CreateRoom)
	r.POST("/save_signal_prekey_bundle", chat.NewWebSockerRouter().SaveSignalKey)
	r.GET("/get_signal_prekey_bundle/:cuid", chat.NewWebSockerRouter().GetSignalKey)
	r.GET("/get_users_rooms", chat.NewWebSockerRouter().GetUsersRooms)
}
