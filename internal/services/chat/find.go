package chat

import (
	"github.com/qianmianyao/parchment-server/internal/models/entity"
	"github.com/qianmianyao/parchment-server/pkg/global"
	"gorm.io/gorm"
)

type Status string
type PasswordType string
type RoomStatus string

type VerificationResults bool

const (
	UserExist         Status              = "exist"
	UserNotExist      Status              = "not_exist"
	RoomExist         Status              = "room_exist"
	RoomNotExist      Status              = "room_not_exist"
	NoPassword        PasswordType        = "no_password"
	NeedPassword      PasswordType        = "need_password"
	InRoom            RoomStatus          = "in_room"
	NotInRoom         RoomStatus          = "not_in_room"
	PasswordCorrect   VerificationResults = true
	PasswordIncorrect VerificationResults = false
)

type Find struct {
	db *gorm.DB
}

func NewFind() *Find {
	return &Find{
		db: global.DB,
	}
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

// AllUsersInTheRoom 获取房间内所有用户
func (f *Find) AllUsersInTheRoom(roomUUID string) []string {
	var roomMembers []entity.RoomMembers
	f.db.Model(&entity.RoomMembers{}).Where("room_uuid = ?", roomUUID).Find(&roomMembers)
	var usersUUID []string
	for _, roomMember := range roomMembers {
		usersUUID = append(usersUUID, roomMember.ChatUserUUID)
	}
	return usersUUID
}

// VerifyPassword 验证房间密码
func (f *Find) VerifyPassword(roomUUID, password string) VerificationResults {
	var room entity.Room
	if err := f.db.Where("uuid = ? AND password = ?", roomUUID, password).First(&room).Error; err != nil {
		return PasswordIncorrect
	}
	if room.Password == password {
		return PasswordCorrect
	}
	return PasswordIncorrect
}
