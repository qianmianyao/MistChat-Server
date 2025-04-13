package chat

import (
	"github.com/gin-gonic/gin"
	"github.com/qianmianyao/parchment-server/internal/services/chat"
	"github.com/qianmianyao/parchment-server/internal/websocket"
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
func (w *WebSockerRouter) WsHandler(hub *websocket.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		websocket.ServeWs(hub, c.Writer, c.Request)
	}
}

// CheckRoomPasswordRequired checks if room password is required
// @Summary 检查房间是否需要密码
// @Description 检查房间是否需要密码
// @Tags Chat
// @Accept json
// @Produce json
// @Param room_uuid query string true "房间ID"
// @Param user_uuid query string true "用户ID"
// @Success 200 {object} utils.Response{data=string} "返回结果"
// @Router /chat/check_room_password [get]
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
		if err := w.chatCreate.Room(params.RoomUUID, params.RoomUUID); err != nil {
			utils.ErrorWithDefault(c)
			return
		}
		if err := w.chatCreate.RoomMembers(params.UserUUID, params.RoomUUID); err != nil {
			utils.ErrorWithDefault(c)
			return
		}
		utils.SuccessWithDefault(c, "加入房间成功")
		return
	}
	utils.SuccessWithDefault(c, "不需要密码")
	return
}

// JoinRoom handles joining a room
// @Summary 加入房间
// @Description 加入房间
// @Tags Chat
// @Accept json
// @Produce json
// @Param room_uuid body string true "房间ID"
// @Param user_uuid body string true "用户ID"
// @Success 200 {object} utils.Response{data=string} "返回结果"
// @Router /chat/join_room [post]
func (w *WebSockerRouter) JoinRoom(c *gin.Context) {
	var params JoinRoomParams
	if err := c.ShouldBindJSON(&params); err != nil {
		utils.ErrorWithDefault(c)
		return
	}

	if err := w.chatCreate.RoomMembers(params.UserUUID, params.RoomUUID); err != nil {
		utils.ErrorWithDefault(c)
		return
	}
	utils.SuccessWithDefault(c, nil)
	return
}
