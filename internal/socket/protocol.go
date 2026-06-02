package socket

import "bytes"

// MsgDecoder is aligned with the utility toolkit MsgDecoder and decodes business objects from the read buffer.
// The boolean return value indicates whether a complete message was decoded.
type MsgDecoder[T any] interface {
	Decode(session *AioSession, readBuffer *bytes.Buffer) (T, bool)
}

// MsgEncoder is aligned with the utility toolkit MsgEncoder and encodes business objects into the write buffer.
type MsgEncoder[T any] interface {
	Encode(session *AioSession, writeBuffer *bytes.Buffer, data T)
}

// Protocol is aligned with the utility toolkit Protocol and combines MsgEncoder with MsgDecoder.
type Protocol[T any] interface {
	MsgEncoder[T]
	MsgDecoder[T]
}

// FuncDecoder is a function adapter for MsgDecoder.
type FuncDecoder[T any] func(session *AioSession, readBuffer *bytes.Buffer) (T, bool)

func (f FuncDecoder[T]) Decode(session *AioSession, readBuffer *bytes.Buffer) (T, bool) {
	return f(session, readBuffer)
}

// FuncEncoder is a function adapter for MsgEncoder.
type FuncEncoder[T any] func(session *AioSession, writeBuffer *bytes.Buffer, data T)

func (f FuncEncoder[T]) Encode(session *AioSession, writeBuffer *bytes.Buffer, data T) {
	f(session, writeBuffer, data)
}
