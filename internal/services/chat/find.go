package chat

import (
	"fmt"
	"github.com/qianmianyao/parchment-server/internal/models/entity"
	"github.com/qianmianyao/parchment-server/pkg/global"

	"gorm.io/gorm"
)

type Status string
type PasswordType string
type RoomStatus string

const (
	UserExist    Status       = "exist"
	UserNotExist Status       = "not_exist"
	RoomExist    Status       = "room_exist"
	RoomNotExist Status       = "room_not_exist"
	NoPassword   PasswordType = "no_password"
	NeedPassword PasswordType = "need_password"
	InRoom       RoomStatus   = "in_room"
	NotInRoom    RoomStatus   = "not_in_room"
)

type Find struct {
	db *gorm.DB
}

func NewFind() *Find {
	return &Find{
		db: global.DB,
	}
}

// GetUserRooms 获取当前用户的房间列表
func (f *Find) GetUserRooms(userUUID string) ([]entity.Room, error) {
	// 先找到用户
	var user entity.ChatUser
	if err := f.db.Where("uuid = ?", userUUID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("用户不存在: %w", err)
	}

	// 查找用户所在的房间成员记录
	var roomMembers []entity.RoomMembers
	if err := f.db.Where("chat_user_uuid = ?", user.ID).Find(&roomMembers).Error; err != nil {
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
	if err := f.db.Where("uuid IN ?", roomUUIDs).Find(&rooms).Error; err != nil {
		return nil, fmt.Errorf("查询房间详情失败: %w", err)
	}

	return rooms, nil
}

// IsUserExist 检查用户是否存在
func (f *Find) IsUserExist(uuid string) Status {
	var count int64
	f.db.Model(&entity.ChatUser{}).Where("uuid = ?", uuid).Count(&count)
	if count > 0 {
		return UserExist
	}
	return UserNotExist
}

// IsRoomExist 检查房间是否存在
func (f *Find) IsRoomExist(roomUUID string) Status {
	var count int64
	f.db.Model(&entity.Room{}).Where("uuid = ?", roomUUID).Count(&count)
	if count > 0 {
		return RoomExist
	}
	return RoomNotExist
}

// IsTheUserIsInTheRoom 检查用户是否在房间内
func (f *Find) IsTheUserIsInTheRoom(uuid, roomUUID string) RoomStatus {
	var count int64
	f.db.Model(&entity.RoomMembers{}).Where("chat_user_uuid = ? AND room_uuid = ?", uuid, roomUUID).Count(&count)
	if count > 0 {
		return InRoom
	}
	return NotInRoom
}

// IsRequirePassword 检查房间是否需要密码
func (f *Find) IsRequirePassword(uuid string) PasswordType {
	var room entity.Room
	if err := f.db.Where("uuid = ?", uuid).First(&room).Error; err != nil {
		return NoPassword
	}
	if !room.Isprivate {
		return NoPassword
	}
	return NeedPassword
}
