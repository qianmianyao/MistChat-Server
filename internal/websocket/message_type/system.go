package message_type

import (
	"encoding/json"
	"github.com/qianmianyao/parchment-server/internal/models/dot"
	"time"
)

// Message 接口定义
type Message interface {
	GetType() dot.MessageType
	StructureMessage() *dot.Envelope
	Serialize() ([]byte, error) // 序列化方法
	Deserialize([]byte) error   // 反序列化方法
}

type SystemMessage struct {
	Text string `json:"text"`
}

func NewSystemMessage(text string) *SystemMessage {
	return &SystemMessage{Text: text}
}

func (sm *SystemMessage) GetType() dot.MessageType {
	return dot.SystemMessage
}

func (sm *SystemMessage) StructureMessage() *dot.Envelope {
	return &dot.Envelope{
		Source: dot.Source{
			Uid:  "system",
			Name: "System",
		},
		Message: dot.DataMessage{
			Type: dot.SystemMessage,
			Content: dot.Content{
				Text: sm.Text,
			},
		},
		Destination: "all",
		Timestamp:   time.Now(),
	}
}

// Serialize 实现序列化方法 - 序列化整个 Envelope
func (sm *SystemMessage) Serialize() ([]byte, error) {
	envelope := sm.StructureMessage()
	return json.Marshal(envelope)
}

// Deserialize 实现反序列化方法 - 从 Envelope 中提取数据
func (sm *SystemMessage) Deserialize(data []byte) error {
	var envelope dot.Envelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return err
	}

	// 从 envelope 中提取系统消息内容
	sm.Text = envelope.Message.Content.Text
	return nil
}
