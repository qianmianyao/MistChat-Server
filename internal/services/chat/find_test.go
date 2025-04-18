package chat

import (
	"gorm.io/gorm"
	"qianmianyao/MistChat-Server/pkg/config"
	"qianmianyao/MistChat-Server/pkg/database"
	"qianmianyao/MistChat-Server/pkg/global"
	"qianmianyao/MistChat-Server/pkg/logger"
	"testing"
)

func setUpTest(t *testing.T) {
	// 初始化配置（必须第一个初始化）
	config.InitConfig()

	// 初始化日志
	logger.InitLogger()

	// 初始化数据库
	database.InitDB()
}

func TestFind_IsUserExist(t *testing.T) {

	setUpTest(t)

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		uuid string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Status
	}{
		{
			name:   "用户存在",
			args:   args{uuid: "d2520bcf-e268-4579-9bfc-6786e027ff47"},
			want:   UserExist,
			fields: fields{db: global.DB},
		},
		{
			name:   "用户不存在",
			args:   args{uuid: "non-existent-uuid"},
			want:   UserNotExist,
			fields: fields{db: global.DB},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Find{
				db: tt.fields.db,
			}
			if got := f.IsUserExist(tt.args.uuid); got != tt.want {
				t.Errorf("IsUserExist() = %v, want %v", got, tt.want)
			}
		})
	}
}
