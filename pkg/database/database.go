package database

import (
	"fmt"
	"log"
	"qianmianyao/MistChat-Server/internal/models/entity"
	"sync"

	"qianmianyao/MistChat-Server/pkg/config"
	"qianmianyao/MistChat-Server/pkg/global"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

// InitDB 初始化数据库连接并设置全局变量
func InitDB() *gorm.DB {
	once.Do(func() {
		cfg := config.GetConfig().Database
		// PostgreSQL 连接格式
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		models := []interface{}{
			&entity.ChatUser{},
			&entity.Room{},
			&entity.RoomMembers{},
			&entity.SignalIdentityKey{},
			&entity.SignalSignedPreKey{},
			&entity.SignalPreKey{},
		}

		if err := db.AutoMigrate(models...); err != nil {
			log.Fatalf("Databses failed to migrate: %v", err)
		}

		// 设置全局DB变量
		global.DB = db
	})
	return db
}
