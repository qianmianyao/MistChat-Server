package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/qianmianyao/parchment-server/internal/services/chat"
	"github.com/qianmianyao/parchment-server/pkg/global"
	"github.com/qianmianyao/parchment-server/pkg/utils"
)

type WebSockerRouter struct {
	chatCreate *chat.Create
	chatFind   *chat.Find
}

type JoinRoomParams struct {
	RoomUUID string `json:"room_uuid" binding:"required"`
	UserUUID string `json:"user_uuid" binding:"required"`
}

func NewWebSockerRouter() *WebSockerRouter {
	return &WebSockerRouter{
		chatCreate: chat.NewCreate(),
		chatFind:   chat.NewFind(),
	}
}

// WsHandler returns a gin.HandlerFunc that handles websocket connections
func (w *WebSockerRouter) WsHandler(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		serveWs(hub, c.Writer, c.Request)
	}
}

// CheckRoomPasswordRequired checks if room password is required
func (w *WebSockerRouter) CheckRoomPasswordRequired(c *gin.Context) {
	var params JoinRoomParams
	if err := c.ShouldBindJSON(&params); err != nil {
		utils.Error(c, 404, "参数错误")
	}

	user := w.chatFind.IsUserExist(params.UserUUID)
	if user == chat.UserNotExist {
		global.Logger.Error("用户不存在!")
		return
	}

	room := w.chatFind.IsRoomExist(params.RoomUUID)
	switch room {
	case chat.RoomExist:
		roomStatus := w.chatFind.IsTheUserIsInTheRoom(params.UserUUID, params.RoomUUID)
		if roomStatus == chat.NotInRoom {
			roomPassword := w.chatFind.IsRequirePassword(params.RoomUUID)
			if roomPassword == chat.NeedPassword {
				utils.SuccessWithDefault(c, "需要密码")
			}
		}
	case chat.RoomNotExist:
		// 如果房间不存在就直接创建
		roomUuid := w.chatCreate.Room(params.RoomUUID, params.RoomUUID)
		w.chatCreate.RoomMembers(params.UserUUID, roomUuid)
	}
}

// JoinRoom handles joining a room
func (w *WebSockerRouter) JoinRoom(c *gin.Context) {
	var params JoinRoomParams
	if err := c.ShouldBindJSON(&params); err != nil {
		utils.Error(c, 404, "参数错误")
	}

	w.chatCreate.RoomMembers(params.UserUUID, params.RoomUUID)

	utils.SuccessWithDefault(c, nil)
}
