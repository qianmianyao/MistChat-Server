package websocket

import "github.com/gin-gonic/gin"

// RegisterWebSocketRoutes registers the WebSocket routes with the provided Gin engine.
func RegisterWebSocketRoutes(r *gin.RouterGroup) {
	hub := NewHub()
	go hub.run()
	r.GET("/connect", NewWebSockerRouter().WsHandler(hub))
	r.GET("/check_room_password", NewWebSockerRouter().CheckRoomPasswordRequired)
	r.GET("/join_room", NewWebSockerRouter().JoinRoom)
}
