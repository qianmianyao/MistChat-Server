package message_type

import (
	"time"

	"github.com/qianmianyao/parchment-server/internal/models/dot"
)

type TextEnvelopeArgs struct {
	SenderUid   string
	SenderName  string
	Destination string
}

type TextMessage struct {
	BaseMessage[string]
	Text string `json:"text"`
}

func NewTextMessage(text string) *TextMessage {
	msg := &TextMessage{Text: text}
	msg.MessageType = dot.TextMessage
	msg.BaseMessage.child = msg
	return msg
}

func (t *TextMessage) StructureMessage(args ...any) *dot.Envelope {

	opt := args[0].(TextEnvelopeArgs) // 更安全地转换

	return &dot.Envelope{
		Source: dot.Source{
			Uid:  opt.SenderUid,
			Name: opt.SenderName,
		},
		Message: dot.DataMessage{
			Type: dot.TextMessage,
			Content: dot.Content{
				Text: t.Text,
			},
		},
		Destination: opt.Destination,
		Timestamp:   time.Now(),
	}
}

func (t *TextMessage) LoadFromEnvelope(env dot.Envelope) error {
	t.Text = env.Message.Content.Text
	return nil
}
