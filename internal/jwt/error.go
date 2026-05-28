package jwt

import "fmt"

// JWTError JWT 相关错误。
type JWTError struct {
	Msg string
	Err error
}

// Error 实现 error 接口。
func (e *JWTError) Error() string {
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

// Unwrap 返回内部错误。
func (e *JWTError) Unwrap() error { return e.Err }

// NewJWTError 构造错误。
func NewJWTError(msg string) *JWTError { return &JWTError{Msg: msg} }

// JWTErrorf 格式化构造错误。
func JWTErrorf(format string, args ...any) *JWTError {
	return &JWTError{Msg: fmt.Sprintf(format, args...)}
}
