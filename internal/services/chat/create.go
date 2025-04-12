package chat

import (
	"github.com/qianmianyao/parchment-server/internal/models/entity"
	"github.com/qianmianyao/parchment-server/pkg/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
func (c *Create) Room(roomName, roomUUID string) string {
	return ""
}

// RoomMembers join room members
func (c *Create) RoomMembers(uuid, roomUUID string) {

}
