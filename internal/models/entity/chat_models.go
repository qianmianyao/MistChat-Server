package entity

import (
	"gorm.io/gorm"
	"time"
)

type ChatUser struct {
	gorm.Model
	UUID     string `gorm:"uniqueIndex;not null"`
	Username string `gorm:"not null"`
	IsOnline bool   `gorm:"not null;default:false"`
}

type Room struct {
	gorm.Model
	UUID      string `gorm:"uniqueIndex;not null"`
	Name      string `gorm:"not null"`
	Password  string `gorm:"column:password"`
	Isprivate bool   `gorm:"not null,default:false"`
}

type RoomMembers struct {
	gorm.Model
	RoomUUID     string    `gorm:"index;not null"`
	ChatUserUUID string    `gorm:"index;not null"`
	JoinTime     time.Time `gorm:"not null"`
}

type SignalIdentityKey struct {
	gorm.Model
	ChatUserUUID   string `gorm:"type:varchar(64);not null;uniqueIndex"`
	RegistrationID uint32 `gorm:"not null"`
	IdentityKey    string `gorm:"type:text;not null"` // Base64 编码
}

type SignalSignedPreKey struct {
	gorm.Model
	ChatUserUUID        string `gorm:"type:varchar(64);not null;index"`
	PreKeyID            uint32 `gorm:"not null"`
	PreKeyPublic        string `gorm:"type:text;not null"`
	PreKeySignature     string `gorm:"type:text;not null"`
	ValidUntilTimestamp int64  `gorm:"not null"`     // 可选：用于前端判断是否需要刷新
	IsActive            bool   `gorm:"default:true"` // 用于轮换时标记是否为当前生效的
}

type SignalPreKey struct {
	gorm.Model
	ChatUserUUID string `gorm:"type:varchar(64);not null;index"`
	PreKeyID     uint32 `gorm:"not null;index"`
	PreKeyPublic string `gorm:"type:text;not null"`
	IsUsed       bool   `gorm:"default:false"` // 是否已被取用
}
