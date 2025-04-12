package chat

import (
	"github.com/qianmianyao/parchment-server/internal/models/entity"
	"github.com/qianmianyao/parchment-server/pkg/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Update struct {
	db *gorm.DB
}

func NewUpdate() *Update {
	return &Update{
		db: global.DB,
	}
}

func (u *Update) UserOnlineStatus(uuid string, isOnline bool) error {
	err := u.db.Model(&entity.ChatUser{}).Where("uuid = ?", uuid).Update("IsOnline", isOnline).Error
	if err != nil {
		global.Logger.Error("更新用户在线状态失败: ", zap.Error(err))
	}
	return nil
}
