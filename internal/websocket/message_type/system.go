package message_type

import (
	"time"

	"qianmianyao/MistChat-Server/internal/models/dot"
)

// SystemMessage 代表系统生成的消息。
type SystemMessage struct {
	BaseMessage[string]     // 嵌入基础消息结构
	Data                any `json:"data"` // 系统消息携带的具体数据
}

// NewSystemMessage 创建并返回一个新的 SystemMessage 实例。
func NewSystemMessage(data any) *SystemMessage {
	msg := &SystemMessage{Data: data}
	msg.MessageType = dot.SystemMessage
	msg.BaseMessage.child = msg
	return msg
}

// StructureMessage 根据 SystemMessage 的数据构建一个 dot.Envelope 结构。
func (sm *SystemMessage) StructureMessage(args ...any) *dot.Envelope {
	return &dot.Envelope{
		Source: dot.Source{
			Uid:  "system",
			Name: "System",
		},
		Message: dot.DataMessage{
			Type: dot.SystemMessage,
			Content: dot.Content{
				Data: sm.Data,
			},
		},
		Destination: "all",
		Timestamp:   time.Now(),
	}
}

// LoadFromEnvelope 从给定的 dot.Envelope 中加载数据到 SystemMessage。
func (sm *SystemMessage) LoadFromEnvelope(env dot.Envelope) error {
	sm.Data = env.Message.Content.Text
	return nil
}
