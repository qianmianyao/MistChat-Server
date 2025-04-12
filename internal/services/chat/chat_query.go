package chat

import (
	"fmt"
	"github.com/qianmianyao/parchment-server/internal/models/entity"
	"github.com/qianmianyao/parchment-server/pkg/global"

	"gorm.io/gorm"
)

type Query struct {
	db *gorm.DB
}

func NewQuery(db *gorm.DB) *Query {
	return &Query{
		db: global.DB,
	}
}

// GetUserRooms 获取当前用户的房间列表
func (q *Query) GetUserRooms(userUUID string) ([]entity.Room, error) {
	// 先找到用户
	var user entity.ChatUser
	if err := q.db.Where("uuid = ?", userUUID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("用户不存在: %w", err)
	}

	// 查找用户所在的房间成员记录
	var roomMembers []entity.RoomMembers
	if err := q.db.Where("chat_user_uuid = ?", user.ID).Find(&roomMembers).Error; err != nil {
		return nil, fmt.Errorf("查询用户房间关系失败: %w", err)
	}

	// 如果用户不在任何房间中
	if len(roomMembers) == 0 {
		return []entity.Room{}, nil
	}

	// 提取所有房间的 UUID
	var roomUUIDs []string
	for _, member := range roomMembers {
		roomUUIDs = append(roomUUIDs, member.RoomUUID)
	}

	// 查询这些房间的详细信息
	var rooms []entity.Room
	if err := q.db.Where("uuid IN ?", roomUUIDs).Find(&rooms).Error; err != nil {
		return nil, fmt.Errorf("查询房间详情失败: %w", err)
	}

	return rooms, nil
}
