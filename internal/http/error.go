package http

import "fmt"

// HTTPError 表示 HTTP 操作过程中发生的异常（对应 hutool-http HttpException）。
type HTTPError struct {
	Msg   string
	Cause error
}

// Error 返回错误描述。
func (e *HTTPError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Msg, e.Cause)
	}
	return e.Msg
}

// Unwrap 返回底层错误。
func (e *HTTPError) Unwrap() error { return e.Cause }

// NewHTTPError 创建 HTTP 错误。
func NewHTTPError(msg string, cause error) *HTTPError {
	return &HTTPError{Msg: msg, Cause: cause}
}

// HTTPErrorf 通过 format 创建 HTTP 错误。
func HTTPErrorf(format string, args ...any) *HTTPError {
	return &HTTPError{Msg: fmt.Sprintf(format, args...)}
}
