package chat

import (
	"github.com/qianmianyao/parchment-server/internal/models/entity"
	"github.com/qianmianyao/parchment-server/pkg/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
	}
	if err := c.db.Create(user).Error; err != nil {
		global.Logger.Error("创建用户失败: ", zap.Error(err))
		return err
	}
	return nil
}

// Room create room
func (c *Create) Room(roomName, roomUUID string) error {
	room := entity.Room{
		UUID:      roomUUID,
		Name:      roomName,
		Isprivate: false,
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
