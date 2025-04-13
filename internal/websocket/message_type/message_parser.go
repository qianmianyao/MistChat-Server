package message_type

import (
	"encoding/json"
	"errors"

	"github.com/qianmianyao/parchment-server/internal/models/dot"
)

// MessageParser 用于处理消息解析
type MessageParser struct{}

// ParseMessage 解析原始JSON消息，返回对应类型的消息对象和解析出的Envelope
func ParseMessage(data []byte) (Message, dot.Envelope, error) {
	// 先将数据解析为通用的Envelope结构
	var envelope dot.Envelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, envelope, errors.New("无法解析消息信封: " + err.Error())
	}

	// 根据消息类型创建相应的消息对象
	var msg Message
	switch envelope.Message.Type {
	case dot.SystemMessage:
		msg = NewSystemMessage("")
	case dot.TextMessage:
		msg = NewTextMessage("")
	// 在这里添加其他消息类型的处理
	default:
		return nil, envelope, errors.New("未知的消息类型: " + string(envelope.Message.Type))
	}

	// 加载消息内容
	if err := msg.LoadFromEnvelope(envelope); err != nil {
		return nil, envelope, errors.New("加载消息内容失败: " + err.Error())
	}

	return msg, envelope, nil
}

// ParseMessageType 仅解析消息类型，不创建完整的消息对象
func ParseMessageType(data []byte) (dot.MessageType, error) {
	// 创建一个临时结构体只包含必要的字段，提高解析效率
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

// CreateMessage 根据消息类型创建空消息对象
func CreateMessage(msgType dot.MessageType) (Message, error) {
	switch msgType {
	case dot.SystemMessage:
		return NewSystemMessage(""), nil
	case dot.TextMessage:
		return NewTextMessage(""), nil
	// 在这里添加其他消息类型的处理
	default:
		return nil, errors.New("不支持的消息类型: " + string(msgType))
	}
}

// ParseMessageContent 解析消息内容到已存在的消息对象
func ParseMessageContent(data []byte, msg Message) (dot.Envelope, error) {
	env, err := msg.Deserialize(data)
	if err != nil {
		return env, errors.New("消息反序列化失败: " + err.Error())
	}
	return env, nil
}
