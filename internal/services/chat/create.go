package chat

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"qianmianyao/MistChat-Server/internal/models/entity"
	"qianmianyao/MistChat-Server/pkg/global"
	"time"
)

type Create struct {
	db *gorm.DB
}

func NewCreate() *Create {
	return &Create{
		db: global.DB,
	}
}

func (c *Create) User(username, uuid string) error {
	user := &entity.ChatUser{
		Username: username,
		UUID:     uuid,
		IsOnline: false,
	}
	if err := c.db.Create(user).Error; err != nil {
		global.Logger.Error("创建用户失败: ", zap.Error(err))
		return err
	}
	return nil
}

// Room create room
func (c *Create) Room(roomName, roomUUID, password string, isprivate bool) error {
	room := entity.Room{
		UUID:      roomUUID,
		Name:      roomName,
		Password:  password,
		Isprivate: isprivate,
	}
	if err := c.db.Create(&room).Error; err != nil {
		global.Logger.Error("创建房间失败: ", zap.Error(err))
		return err
	}
	return nil
}

// RoomMembers join room members
func (c *Create) RoomMembers(uuid, roomUUID string) error {
	roomMembers := entity.RoomMembers{
		ChatUserUUID: uuid,
		RoomUUID:     roomUUID,
		JoinTime:     time.Now(),
	}
	if err := c.db.Create(&roomMembers).Error; err != nil {
		global.Logger.Error("加入房间失败: ", zap.Error(err))
		return err
	}
	return nil
}

// SignalIdentityKey 身份密钥
func (c *Create) SignalIdentityKey(signalIdentityKey entity.SignalIdentityKey) error {
	if err := c.db.Create(&signalIdentityKey).Error; err != nil {
		global.Logger.Error("创建 SignalIdentityKey 失败: ", zap.Error(err))
		return err
	}
	return nil
}

// SignalSignedPreKey 预签名密钥
func (c *Create) SignalSignedPreKey(signalSignedPreKey entity.SignalSignedPreKey) error {
	if err := c.db.Create(&signalSignedPreKey).Error; err != nil {
		global.Logger.Error("创建 SignalSignedPreKey 失败: ", zap.Error(err))
		return err
	}
	return nil
}

// SignalPreKey 一次性密钥
func (c *Create) SignalPreKey(signalPreKey entity.SignalPreKey) error {
	if err := c.db.Create(&signalPreKey).Error; err != nil {
		global.Logger.Error("创建 SignalPreKey 失败: ", zap.Error(err))
		return err
	}
	return nil
}
