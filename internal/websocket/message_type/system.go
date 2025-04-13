package message_type

import (
	"github.com/qianmianyao/parchment-server/internal/models/dot"
	"time"
)

type SystemMessage struct {
	BaseMessage[string]
	Data any `json:"data"`
}

func NewSystemMessage(data any) *SystemMessage {
	msg := &SystemMessage{Data: data}
	msg.MessageType = dot.SystemMessage
	msg.BaseMessage.child = msg
	return msg
}

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

func (sm *SystemMessage) LoadFromEnvelope(env dot.Envelope) error {
	sm.Data = env.Message.Content.Text
	return nil
}
