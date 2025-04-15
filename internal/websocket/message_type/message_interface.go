package message_type

import (
	"encoding/json"
	"errors"

	"github.com/qianmianyao/parchment-server/internal/models/dot"
)

// Message 是所有 WebSocket 消息类型的通用接口。
// 定义了消息的基本行为，如获取类型、序列化和反序列化。
type Message interface {
	GetType() dot.MessageType
	LoadFromEnvelope(dot.Envelope) error
	SerializeWithArgs(args ...any) ([]byte, error)
	Deserialize([]byte) (dot.Envelope, error)
}

// BaseMessage 为实现 Message 接口的类型提供了一个泛型基础结构。
type BaseMessage[T any] struct {
	MessageType dot.MessageType // 消息的具体类型。
	Data        T               // 消息携带的数据。
	child       Message         // 指向实现此基础结构的具体子类实例，用于方法委托。
}

// GetType 返回 BaseMessage 的消息类型。
func (bm *BaseMessage[T]) GetType() dot.MessageType {
	return bm.MessageType
}

// SerializeWithArgs 将消息序列化为 JSON 格式的 []byte。
// 委托给子类实现的 StructureMessage 方法。
func (bm *BaseMessage[T]) SerializeWithArgs(args ...any) ([]byte, error) {
	// structureCapable 定义了一个内部接口，用于检查子类是否实现了 StructureMessage 方法。
	type structureCapable interface {
		StructureMessage(args ...any) *dot.Envelope
	}

	if msg, ok := bm.child.(structureCapable); ok {
		return json.Marshal(msg.StructureMessage(args...))
	}

	return nil, errors.New("child didn't implement StructureMessage(...any)")
}

// Deserialize 将 JSON 数据反序列化为 dot.Envelope 结构，
// 并委托给子类的 LoadFromEnvelope 方法加载消息数据。
func (bm *BaseMessage[T]) Deserialize(data []byte) (dot.Envelope, error) {
	var env dot.Envelope

	if err := json.Unmarshal(data, &env); err != nil {
		return env, err
	}

	err := bm.child.LoadFromEnvelope(env)
	return env, err
}
