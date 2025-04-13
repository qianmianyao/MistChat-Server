package entity

import (
	"gorm.io/gorm"
	"time"
)

type ChatUser struct {
	gorm.Model
	UUID     string `gorm:"uniqueIndex;not null"`
	Username string `gorm:"not null"`
	IsOnline bool   `gorm:"not null;default:true"`
}

type Room struct {
	gorm.Model
	UUID      string `gorm:"uniqueIndex;not null"`
	Name      string `gorm:"not null"`
	Isprivate bool   `gorm:"not null,default:false"`
}

type RoomMembers struct {
	gorm.Model
	RoomUUID     string    `gorm:"index;not null"`
	ChatUserUUID string    `gorm:"index;not null"`
	JoinTime     time.Time `gorm:"not null"`
}
