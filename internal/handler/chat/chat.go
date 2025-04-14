package chat

import (
	"github.com/gin-gonic/gin"
	"github.com/qianmianyao/parchment-server/internal/models/dot"
	"github.com/qianmianyao/parchment-server/internal/models/entity"
	"github.com/qianmianyao/parchment-server/internal/services/chat"
	"github.com/qianmianyao/parchment-server/internal/websocket"
	"github.com/qianmianyao/parchment-server/pkg/encryption"
	"github.com/qianmianyao/parchment-server/pkg/utils"
)

type WebSockerRouter struct {
	chatCreate *chat.Create
	chatFind   *chat.Find
}

type UpdateKeyPair struct {
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
// @Param username query string false "用户名"
// @Success 101 {string} string "Switching Protocols to websocket"
// @Router /chat/connect [get]
func (w *WebSockerRouter) WsHandler(hub *websocket.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		websocket.ServeWs(hub, c.Writer, c.Request)
	}
}

// CheckRoomPasswordRequired 检查房间密码是否需要
func (w *WebSockerRouter) CheckRoomPasswordRequired(c *gin.Context) {
	var data dot.JoinRoomData
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorWithDefault(c)
		return
	}

	roomExist := w.chatFind.IsRoomExist(data.RoomUUID)
	if roomExist == chat.RoomNotExist {
		utils.ErrorWithDefault(c)
		return
	}

	passwordExist := w.chatFind.IsRequirePassword(data.RoomUUID)
	if passwordExist == chat.NeedPassword {
		utils.FailWithDefault(c, "需要密码")
		return
	}
}

// CreateRoom 创建房间
func (w *WebSockerRouter) CreateRoom(c *gin.Context) {
	var data dot.CreateRoomData
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorWithDefault(c)
		return
	}

	var isprivate = false
	if data.Password != "" {
		isprivate = true
	}

	roomId, _ := encryption.GenerateUID("r_")
	if err := w.chatCreate.Room(data.RoomName, roomId, data.Password, isprivate); err != nil {
		utils.ErrorWithDefault(c)
		return
	}
	if err := w.chatCreate.RoomMembers(data.UserUUID, roomId); err != nil {
		utils.ErrorWithDefault(c)
		return
	}
	utils.SuccessWithDefault(c, map[string]string{"roomUUID": roomId})
	return
}

// JoinRoom 加入房间
func (w *WebSockerRouter) JoinRoom(c *gin.Context) {
	var data dot.JoinRoomData
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorWithDefault(c)
		return
	}
	verificationResults := w.chatFind.VerifyPassword(data.RoomUUID, data.Password)
	if verificationResults == chat.PasswordIncorrect {
		utils.FailWithDefault(c, "密码错误")
		return
	}

	if err := w.chatCreate.RoomMembers(data.UserUUID, data.RoomUUID); err != nil {
		utils.ErrorWithDefault(c)
		return
	}
	utils.SuccessWithDefault(c, nil)
	return
}

// SaveSignalKey 更新用户密钥对
func (w *WebSockerRouter) SaveSignalKey(c *gin.Context) {
	var data dot.SignalData
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorWithDefault(c)
		return
	}

	var signalIdentityKey = entity.SignalIdentityKey{
		ChatUserUUID:   data.Address.Name,
		RegistrationID: uint32(data.RegistrationId),
		IdentityKey:    data.IdentityKey,
	}
	// 创建身份密钥
	if err := w.chatCreate.SignalIdentityKey(signalIdentityKey); err != nil {
		utils.ErrorWithDefault(c)
		return
	}

	var signalSignedPreKey = entity.SignalSignedPreKey{
		ChatUserUUID:        data.Address.Name,
		PreKeyID:            uint32(data.PreKey.Id),
		PreKeyPublic:        data.PreKey.PublicKey,
		PreKeySignature:     data.SignedPreKey.Signature,
		ValidUntilTimestamp: 0,
		IsActive:            true,
	}
	// 创建预签名密钥
	if err := w.chatCreate.SignalSignedPreKey(signalSignedPreKey); err != nil {
		utils.ErrorWithDefault(c)
		return
	}

	var signalPreKey = entity.SignalPreKey{
		ChatUserUUID: data.Address.Name,
		PreKeyID:     uint32(data.PreKey.Id),
		PreKeyPublic: data.PreKey.PublicKey,
		IsUsed:       false,
	}
	// 创建一次性密钥
	if err := w.chatCreate.SignalPreKey(signalPreKey); err != nil {
		utils.ErrorWithDefault(c)
		return
	}

	utils.SuccessWithDefault(c, nil)
	return
}
