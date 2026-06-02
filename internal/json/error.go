package json

import "fmt"

// JSONError 对应 the utility JSONException。
type JSONError struct {
	Msg   string
	Cause error
}

// NewJSONError 使用消息构造错误。
func NewJSONError(format string, args ...any) *JSONError {
	return &JSONError{Msg: fmt.Sprintf(format, args...)}
}

// WrapJSONError 包装一个底层错误。
func WrapJSONError(cause error, format string, args ...any) *JSONError {
	return &JSONError{Msg: fmt.Sprintf(format, args...), Cause: cause}
}

func (e *JSONError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Msg, e.Cause)
	}
	return e.Msg
}

// Unwrap 支持 errors.Is/As。
func (e *JSONError) Unwrap() error { return e.Cause }
