package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/qianmianyao/parchment-server/internal/services/chat"
	"github.com/qianmianyao/parchment-server/pkg/utils"
)

type WebSockerRouter struct {
	chatCreate *chat.Create
	chatFind   *chat.Find
}

type JoinRoomParams struct {
	RoomUUID string `form:"room_uuid" binding:"required"`
	UserUUID string `form:"user_uuid" binding:"required"`
}

func NewWebSockerRouter() *WebSockerRouter {
	return &WebSockerRouter{
		chatCreate: chat.NewCreate(),
		chatFind:   chat.NewFind(),
	}
}

// WsHandler godoc
// @Summary WebSocket连接
// @Description 建立WebSocket连接
// @Tags Chat
// @Accept json
// @Produce json
// @Param uuid query string true "用户ID"
// @Param username query string false "用户名"
// @Success 101 {string} string "Switching Protocols to websocket"
// @Router /chat/connect [get]
func (w *WebSockerRouter) WsHandler(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		serveWs(hub, c.Writer, c.Request)
	}
}

// CheckRoomPasswordRequired checks if room password is required
func (w *WebSockerRouter) CheckRoomPasswordRequired(c *gin.Context) {
	var params JoinRoomParams
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.ErrorWithDefault(c)
		return
	}

	user := w.chatFind.IsUserExist(params.UserUUID)
	if user == chat.UserNotExist {
		utils.ErrorWithDefault(c)
		return
	}

	room := w.chatFind.IsRoomExist(params.RoomUUID)
	switch room {
	case chat.RoomExist:
		roomStatus := w.chatFind.IsTheUserIsInTheRoom(params.UserUUID, params.RoomUUID)
		if roomStatus == chat.NotInRoom {
			roomPassword := w.chatFind.IsRequirePassword(params.RoomUUID)
			if roomPassword == chat.NeedPassword {
				utils.FailWithDefault(c, "需要密码")
				return
			}
		}
	case chat.RoomNotExist:
		// 如果房间不存在就直接创建
		roomUuid := w.chatCreate.Room(params.RoomUUID, params.RoomUUID)
		w.chatCreate.RoomMembers(params.UserUUID, roomUuid)
		utils.SuccessWithDefault(c, "房间创建成功")
		return
	}

	utils.SuccessWithDefault(c, "不需要密码")
	return
}

// JoinRoom handles joining a room
func (w *WebSockerRouter) JoinRoom(c *gin.Context) {
	var params JoinRoomParams
	if err := c.ShouldBindJSON(&params); err != nil {
		utils.ErrorWithDefault(c)
		return
	}

	w.chatCreate.RoomMembers(params.UserUUID, params.RoomUUID)

	utils.SuccessWithDefault(c, nil)
	return
}
