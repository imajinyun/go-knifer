package cron

import "fmt"

// CronError is aligned with hutool CronException and represents cron-related errors.
type CronError struct {
	Msg   string
	Cause error
}

// NewCronError creates an error with a formatted message.
func NewCronError(format string, args ...any) *CronError {
	return &CronError{Msg: fmt.Sprintf(format, args...)}
}

// WrapCronError wraps an underlying error with a formatted message.
func WrapCronError(cause error, format string, args ...any) *CronError {
	return &CronError{Msg: fmt.Sprintf(format, args...), Cause: cause}
}

func (e *CronError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Msg, e.Cause)
	}
	return e.Msg
}

// Unwrap supports errors.Is and errors.As.
func (e *CronError) Unwrap() error { return e.Cause }
