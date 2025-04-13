package message_type

import (
	"time"

	"github.com/qianmianyao/parchment-server/internal/models/dot"
)

type SystemMessage struct {
	BaseMessage[string]
	Text string `json:"text"`
}

func NewSystemMessage(text string) *SystemMessage {
	msg := &SystemMessage{Text: text}
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
				Text: sm.Text,
			},
		},
		Destination: "all",
		Timestamp:   time.Now(),
	}
}

func (sm *SystemMessage) LoadFromEnvelope(env dot.Envelope) error {
	sm.Text = env.Message.Content.Text
	return nil
}
