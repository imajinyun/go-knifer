package cron

import "fmt"

// CronError 对应 hutool 中的 CronException，表示定时任务相关错误。
type CronError struct {
	Msg   string
	Cause error
}

// NewCronError 使用消息构造错误。
func NewCronError(format string, args ...any) *CronError {
	return &CronError{Msg: fmt.Sprintf(format, args...)}
}

// WrapCronError 包装一个底层错误。
func WrapCronError(cause error, format string, args ...any) *CronError {
	return &CronError{Msg: fmt.Sprintf(format, args...), Cause: cause}
}

func (e *CronError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Msg, e.Cause)
	}
	return e.Msg
}

// Unwrap 支持 errors.Is/As。
func (e *CronError) Unwrap() error { return e.Cause }
