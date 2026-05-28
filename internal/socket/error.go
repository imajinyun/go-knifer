package socket

import "fmt"

// SocketRuntimeError 对应 hutool 的 SocketRuntimeException，
// 表示 Socket 通讯过程中的运行时错误。
type SocketRuntimeError struct {
	Msg   string
	Cause error
}

// Error 返回错误消息。
func (e *SocketRuntimeError) Error() string {
	if e == nil {
		return ""
	}
	if e.Cause != nil && e.Msg == "" {
		return e.Cause.Error()
	}
	if e.Cause != nil {
		return e.Msg + ": " + e.Cause.Error()
	}
	return e.Msg
}

// Unwrap 支持 errors.Is/As。
func (e *SocketRuntimeError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// NewSocketError 通过任意 error 创建。
func NewSocketError(err error) *SocketRuntimeError {
	if err == nil {
		return nil
	}
	return &SocketRuntimeError{Msg: err.Error(), Cause: err}
}

// NewSocketErrorMsg 直接通过消息创建。
func NewSocketErrorMsg(msg string) *SocketRuntimeError {
	return &SocketRuntimeError{Msg: msg}
}

// NewSocketErrorf 格式化创建。
func NewSocketErrorf(format string, args ...any) *SocketRuntimeError {
	return &SocketRuntimeError{Msg: fmt.Sprintf(format, args...)}
}

// WrapSocketError 包装一个底层 error 并附加消息。
func WrapSocketError(err error, msg string) *SocketRuntimeError {
	if err == nil {
		return nil
	}
	return &SocketRuntimeError{Msg: msg, Cause: err}
}
