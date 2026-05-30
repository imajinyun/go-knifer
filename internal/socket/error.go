package socket

import "fmt"

// SocketRuntimeError represents a runtime error during socket communication.
type SocketRuntimeError struct {
	Msg   string
	Cause error
}

// Error returns the error message.
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

// Unwrap supports errors.Is and errors.As.
func (e *SocketRuntimeError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// NewSocketError creates a SocketRuntimeError from any error.
func NewSocketError(err error) *SocketRuntimeError {
	if err == nil {
		return nil
	}
	return &SocketRuntimeError{Msg: err.Error(), Cause: err}
}

// NewSocketErrorMsg creates a SocketRuntimeError from a message.
func NewSocketErrorMsg(msg string) *SocketRuntimeError {
	return &SocketRuntimeError{Msg: msg}
}

// NewSocketErrorf creates a formatted SocketRuntimeError.
func NewSocketErrorf(format string, args ...any) *SocketRuntimeError {
	return &SocketRuntimeError{Msg: fmt.Sprintf(format, args...)}
}

// WrapSocketError wraps an underlying error with an additional message.
func WrapSocketError(err error, msg string) *SocketRuntimeError {
	if err == nil {
		return nil
	}
	return &SocketRuntimeError{Msg: msg, Cause: err}
}
