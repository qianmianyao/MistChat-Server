package message_type

import (
	"encoding/json"
	"errors"

	"github.com/qianmianyao/parchment-server/internal/models/dot"
)

// Message 接口定义
// 不强制实现 StructureMessage，改为通过 SerializeWithArgs 间接调用
// 实现子类只需关注内容逻辑

type Message interface {
	GetType() dot.MessageType
	LoadFromEnvelope(dot.Envelope) error
	SerializeWithArgs(args ...any) ([]byte, error)
	Deserialize([]byte) (dot.Envelope, error)
}

// BaseMessage 提供泛型基础实现

type BaseMessage[T any] struct {
	MessageType dot.MessageType
	Data        T
	child       Message // 注入子类
}

func (bm *BaseMessage[T]) GetType() dot.MessageType {
	return bm.MessageType
}

func (bm *BaseMessage[T]) SerializeWithArgs(args ...any) ([]byte, error) {
	type structureCapable interface {
		StructureMessage(args ...any) *dot.Envelope
	}
	if msg, ok := bm.child.(structureCapable); ok {
		return json.Marshal(msg.StructureMessage(args...))
	}
	return nil, errors.New("child didn't come true StructureMessage(...any)")
}

func (bm *BaseMessage[T]) Deserialize(data []byte) (dot.Envelope, error) {
	var env dot.Envelope
	if err := json.Unmarshal(data, &env); err != nil {
		return env, err
	}
	err := bm.child.LoadFromEnvelope(env)
	return env, err
}
