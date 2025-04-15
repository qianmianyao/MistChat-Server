package message_type

import (
	"encoding/json"
	"errors"

	"github.com/qianmianyao/parchment-server/internal/models/dot"
)

// MessageParser 提供了用于解析和处理 WebSocket 消息的工具函数。
type MessageParser struct{}

// ParseMessage 将原始的 JSON 格式数据解析为具体的 Message 对象和通用的 Envelope 结构。
func ParseMessage(data []byte) (Message, dot.Envelope, error) {
	var envelope dot.Envelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, envelope, errors.New("无法解析消息信封: " + err.Error())
	}

	var msg Message
	switch envelope.Message.Type {
	case dot.SystemMessage:
		msg = NewSystemMessage("")
	case dot.TextMessage:
		msg = NewTextMessage("")
	// TODO: 在这里添加其他消息类型的处理
	default:
		return nil, envelope, errors.New("未知的消息类型: " + string(envelope.Message.Type))
	}

	if err := msg.LoadFromEnvelope(envelope); err != nil {
		return nil, envelope, errors.New("加载消息内容失败: " + err.Error())
	}

	return msg, envelope, nil
}

// ParseMessageType 从原始的 JSON 格式数据中仅解析出消息类型。
func ParseMessageType(data []byte) (dot.MessageType, error) {
	type EnvelopeType struct {
		Message struct {
			Type dot.MessageType `json:"type"`
		} `json:"message"`
	}

	var envType EnvelopeType
	if err := json.Unmarshal(data, &envType); err != nil {
		return "", errors.New("无法解析消息类型: " + err.Error())
	}

	return envType.Message.Type, nil
}

// CreateMessage 根据给定的消息类型创建对应类型的空 Message 对象实例。
func CreateMessage(msgType dot.MessageType) (Message, error) {
	switch msgType {
	case dot.SystemMessage:
		return NewSystemMessage(""), nil
	case dot.TextMessage:
		return NewTextMessage(""), nil
	// TODO: 在这里添加其他消息类型的处理
	default:
		return nil, errors.New("不支持的消息类型: " + string(msgType))
	}
}

// ParseMessageContent 将原始的 JSON 格式数据反序列化到已存在的 Message 对象中。
func ParseMessageContent(data []byte, msg Message) (dot.Envelope, error) {
	env, err := msg.Deserialize(data)
	if err != nil {
		return env, errors.New("消息反序列化失败: " + err.Error())
	}
	return env, nil
}
