package message_type

import (
	"errors"
	"time"

	"qianmianyao/MistChat-Server/internal/models/dot"
)

// TextEnvelopeArgs 定义了构建文本消息 Envelope 所需的参数结构。
type TextEnvelopeArgs struct {
	SenderUid   string // 发送者的唯一标识符
	SenderName  string // 发送者的名称
	Destination string // 消息的目标地址或标识符
}

// TextMessage 代表一个纯文本消息。
// 它嵌入了 BaseMessage，并包含一个表示文本内容的 Text 字段。
type TextMessage struct {
	BaseMessage[string]        // 嵌入基础消息结构，泛型参数 string 对应 Text 字段。
	Text                string `json:"text"` // 文本消息的具体内容。
}

// NewTextMessage 创建并返回一个新的 TextMessage 实例。
// text 是要包含在文本消息中的内容。
func NewTextMessage(text string) *TextMessage {
	msg := &TextMessage{Text: text}
	msg.MessageType = dot.TextMessage
	msg.BaseMessage.child = msg
	return msg
}

// StructureMessage 根据 TextMessage 的数据和传入的参数构建一个 dot.Envelope 结构。
// args 参数列表应包含一个 TextEnvelopeArgs 类型的参数，用于填充 Envelope 的 Source 和 Destination 字段。
func (t *TextMessage) StructureMessage(args ...any) (*dot.Envelope, error) {
	// 检查参数数量
	if len(args) != 1 {
		return nil, errors.New("StructureMessage for TextMessage expects exactly one argument of type TextEnvelopeArgs")
	}
	// 类型断言，获取 TextEnvelopeArgs 参数
	opt, ok := args[0].(TextEnvelopeArgs)
	if !ok {
		return nil, errors.New("invalid argument type: expected TextEnvelopeArgs")
	}

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
	}, nil
}

// LoadFromEnvelope 从给定的 dot.Envelope 中加载数据到 TextMessage。
func (t *TextMessage) LoadFromEnvelope(env dot.Envelope) error {
	t.Text = env.Message.Content.Text
	return nil
}
