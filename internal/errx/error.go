package errx

import (
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/hashicorp/go-multierror"
)

// WithStack is implemented by errors that can expose a string stack trace.
type WithStack interface {
	Stack() string
}

// PanicError wraps a recovered panic value and records the stack captured at the
// recovery point. If the panic value is an error, Unwrap exposes it for errors.Is
// and errors.As.
type PanicError struct {
	Value      any
	Cause      error
	StackTrace StackTrace
}

// Error returns the recovered panic value as an error message.
func (e *PanicError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return fmt.Sprint(e.Value)
}

// Unwrap returns the panic value when it is an error.
func (e *PanicError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// Stack returns the stack captured when the panic was recovered.
func (e *PanicError) Stack() string {
	if e == nil || len(e.StackTrace) == 0 {
		return ""
	}
	return fmt.Sprintf("%+v", e.StackTrace)
}

// GetStack returns the stack attached to err, or the current goroutine stack.
func GetStack(err error) string {
	if err == nil {
		return ""
	}
	var ws WithStack
	if errors.As(err, &ws) {
		return ws.Stack()
	}
	return string(debug.Stack())
}

// ErrorIs is like errors.Is, but it also checks each member of a multierror.
func ErrorIs(err error, target error) bool {
	if target == nil {
		return err == nil
	}
	if errors.Is(err, target) {
		return true
	}
	var merr *multierror.Error
	if !errors.As(err, &merr) {
		return false
	}
	for _, item := range merr.Errors {
		if ErrorIs(item, target) {
			return true
		}
	}
	return false
}
