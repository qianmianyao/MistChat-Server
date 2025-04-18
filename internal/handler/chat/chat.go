package chat

import (
	"fmt"
	"github.com/qianmianyao/parchment-server/pkg/global"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qianmianyao/parchment-server/internal/models/dot"
	"github.com/qianmianyao/parchment-server/internal/models/entity"
	"github.com/qianmianyao/parchment-server/internal/services/chat"
	"github.com/qianmianyao/parchment-server/internal/websocket"
	"github.com/qianmianyao/parchment-server/pkg/encryption"
	"github.com/qianmianyao/parchment-server/pkg/utils"
)

// WebSockerRouter 定义了处理聊天相关 WebSocket 请求的路由结构。
type WebSockerRouter struct {
	chatCreate *chat.Create
	chatFind   *chat.Find
	chatUpdate *chat.Update
}

// NewWebSockerRouter 创建并返回一个新的 WebSockerRouter 实例。
func NewWebSockerRouter() *WebSockerRouter {
	return &WebSockerRouter{
		chatCreate: chat.NewCreate(),
		chatFind:   chat.NewFind(),
		chatUpdate: chat.NewUpdate(),
	}
}

// WsHandler 处理 WebSocket 连接请求。
// @Summary WebSocket连接
// @Description 建立WebSocket连接，升级HTTP连接为WebSocket。
// @Tags Chat
// @Accept json
// @Produce json
// @Param username query string false "用户名 (可选)"
// @Success 101 {string} string "Switching Protocols" "成功切换协议到WebSocket"
// @Router /chat/connect [get]
func (w *WebSockerRouter) WsHandler(hub *websocket.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		websocket.ServeWs(hub, c.Writer, c.Request)
	}
}

// Register 处理用户注册请求。
func (w *WebSockerRouter) Register(c *gin.Context) {
	var data dot.RegisterData
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorWithDefault(c)
		return
	}
	// 生成 uuid
	uuid, err := encryption.GenerateUID("u_")
	if err != nil {
		global.Logger.Error(fmt.Sprintf("Failed to generate UID: %v", err))
		utils.ErrorWithDefault(c)
		return
	}
	// 验证 uuid 格式。
	if ok, err := encryption.ValidateUID(uuid, "u_"); err != nil || !ok {
		global.Logger.Warn(fmt.Sprintf("Invalid UID provided or generated: %s, validation error: %v", uuid, err))
		utils.ErrorWithDefault(c)
		return
	}
	// 创建用户
	if err := w.chatCreate.User(data.Username, uuid); err != nil {
		utils.ErrorWithDefault(c)
		return
	}
	utils.Success(c, dot.RegisterResponse{Username: data.Username, UUID: uuid}, "注册成功")
}

// CheckRoomPasswordRequired 检查加入房间是否需要密码。
// @Summary 检查房间密码要求
// @Description 根据房间UUID检查该房间是否存在以及是否需要密码才能加入。
// @Tags Chat
// @Accept json
// @Produce json
// @Param room body dot.JoinRoomData true "包含房间UUID的数据"
// @Success 200 {object} utils.Response "房间存在且不需要密码"
// @Failure 400 {object} utils.Response "请求参数错误或房间不存在"
// @Failure 401 {object} utils.Response "需要密码"
// @Router /chat/check-password [post]
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

// CreateRoom 处理创建新聊天房间的请求。
// @Summary 创建聊天房间
// @Description 创建一个新的聊天房间，可以设置房间名和可选的密码。
// @Tags Chat
// @Accept json
// @Produce json
// @Param room body dot.CreateRoomData true "创建房间所需的数据 (房间名, 用户UUID, 可选密码)"
// @Success 200 {object} utils.Response{data=map[string]string} "成功创建房间，返回房间UUID"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 500 {object} utils.Response "服务器内部错误 (创建房间或添加成员失败)"
// @Router /chat/create-room [post]
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
}

// JoinRoom 处理用户加入现有聊天房间的请求。
// @Summary 加入聊天房间
// @Description 用户根据房间UUID和可选的密码加入一个已存在的聊天房间。
// @Tags Chat
// @Accept json
// @Produce json
// @Param join body dot.JoinRoomData true "加入房间所需的数据 (房间UUID, 用户UUID, 可选密码)"
// @Success 200 {object} utils.Response "成功加入房间"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 401 {object} utils.Response "密码错误"
// @Failure 500 {object} utils.Response "服务器内部错误 (添加成员失败)"
// @Router /chat/join-room [post]
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
}

// SaveSignalKey 处理上传用户 Signal 协议密钥束的请求。
// @Summary 保存Signal密钥
// @Description接收并存储用户的 Signal 协议密钥，包括身份密钥、预签名密钥和一次性密钥。
// @Tags Signal
// @Accept json
// @Produce json
// @Param keys body dot.SignalData true "包含用户地址和密钥束的数据"
// @Success 200 {object} utils.Response "成功保存密钥"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 500 {object} utils.Response "服务器内部错误 (保存密钥失败)"
// @Router /chat/save-signal-key [post]
func (w *WebSockerRouter) SaveSignalKey(c *gin.Context) {
	var data dot.SignalData
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorWithDefault(c)
		return
	}

	var signalIdentityKey = entity.SignalIdentityKey{
		ChatUserUUID:   data.Address.UUID,
		RegistrationID: uint32(data.RegistrationId),
		IdentityKey:    data.IdentityKey,
	}
	// 创建身份密钥
	if err := w.chatCreate.SignalIdentityKey(signalIdentityKey); err != nil {
		utils.ErrorWithDefault(c)
		return
	}

	var signalSignedPreKey = entity.SignalSignedPreKey{
		ChatUserUUID:        data.Address.UUID,
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
		ChatUserUUID: data.Address.UUID,
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
}

// GetSignalKey 处理获取指定用户 Signal 协议密钥束的请求。
// @Summary 获取Signal密钥
// @Description 根据用户的聊天ID (cuid) 查询并返回其 Signal 协议密钥束，并将使用的一次性密钥标记为已用。
// @Tags Signal
// @Accept json
// @Produce json
// @Param cuid path uint true "用户的聊天ID"
// @Success 200 {object} utils.Response{data=dot.SignalData} "成功获取密钥束"
// @Failure 400 {object} utils.Response "无效的用户ID格式"
// @Failure 500 {object} utils.Response "服务器内部错误 (查询或更新密钥失败)"
// @Router /chat/get-signal-key/{cuid} [get]
func (w *WebSockerRouter) GetSignalKey(c *gin.Context) {
	cuid := c.Param("cuid")
	num, err := strconv.ParseUint(cuid, 10, 0)
	if err != nil {
		utils.ErrorWithDefault(c)
		return
	}
	uuid := w.chatFind.ChatUserUUIDByID(uint(num))

	signalIdentityKey := w.chatFind.SignalIdentityKey(uuid)
	signalSignedPreKey := w.chatFind.SignalSignedPreKey(uuid)
	signalPreKey, err := w.chatFind.SignalPreKey(uuid)
	if err != nil {
		utils.ErrorWithDefault(c)
		return
	}
	var data = dot.SignalData{
		Address:        nil,
		RegistrationId: int(signalIdentityKey.RegistrationID),
		IdentityKey:    signalIdentityKey.IdentityKey,
		SignedPreKey: dot.SignedPreKey{
			Id:        int(signalSignedPreKey.PreKeyID),
			PublicKey: signalSignedPreKey.PreKeyPublic,
			Signature: signalSignedPreKey.PreKeySignature,
		},
		PreKey: dot.PreKey{
			Id:        int(signalPreKey.PreKeyID),
			PublicKey: signalPreKey.PreKeyPublic,
		},
	}
	if err := w.chatUpdate.MarkUsed(signalPreKey.PreKeyID); err != nil {
		utils.ErrorWithDefault(c)
		return
	}
	utils.SuccessWithDefault(c, &data)
}

func (w *WebSockerRouter) GetUsersRooms(c *gin.Context) {
	var params dot.GetUsersRoomsParams
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.Error(c, "参数错误")
		return
	}
	rooms, err := w.chatFind.UsersRooms(params.UserUUID)
	if err != nil {
		utils.ErrorWithDefault(c)
		return
	}

	utils.SuccessWithDefault(c, rooms)
}
