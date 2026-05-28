package socket

import "bytes"

// MsgDecoder 对应 hutool MsgDecoder，从读缓冲中解码出业务对象。
// 返回 nil 表示数据尚不足以解码出完整消息。
type MsgDecoder[T any] interface {
	Decode(session *AioSession, readBuffer *bytes.Buffer) (T, bool)
}

// MsgEncoder 对应 hutool MsgEncoder，将业务对象编码到写缓冲。
type MsgEncoder[T any] interface {
	Encode(session *AioSession, writeBuffer *bytes.Buffer, data T)
}

// Protocol 对应 hutool Protocol，是 MsgEncoder + MsgDecoder 的组合。
type Protocol[T any] interface {
	MsgEncoder[T]
	MsgDecoder[T]
}

// FuncDecoder 函数式 MsgDecoder。
type FuncDecoder[T any] func(session *AioSession, readBuffer *bytes.Buffer) (T, bool)

func (f FuncDecoder[T]) Decode(session *AioSession, readBuffer *bytes.Buffer) (T, bool) {
	return f(session, readBuffer)
}

// FuncEncoder 函数式 MsgEncoder。
type FuncEncoder[T any] func(session *AioSession, writeBuffer *bytes.Buffer, data T)

func (f FuncEncoder[T]) Encode(session *AioSession, writeBuffer *bytes.Buffer, data T) {
	f(session, writeBuffer, data)
}
