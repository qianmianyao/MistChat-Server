package websocket

import "github.com/gin-gonic/gin"

// RegisterWebSocketRoutes registers the WebSocket routes with the provided Gin engine.
func RegisterWebSocketRoutes(r *gin.RouterGroup) {
	hub := NewHub()
	go hub.run()
	r.GET("/chat", NewWebSockerRouter().WsHandler(hub))
	r.GET("/join_room", NewWebSockerRouter().JoinRoom)
	r.GET("/check_room_password", NewWebSockerRouter().CheckRoomPasswordRequired)
}
