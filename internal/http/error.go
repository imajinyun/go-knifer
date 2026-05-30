package http

import "fmt"

// HTTPError represents an error during HTTP operations, aligned with hutool-http HttpException.
type HTTPError struct {
	Msg   string
	Cause error
}

// Error returns the error message.
func (e *HTTPError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Msg, e.Cause)
	}
	return e.Msg
}

// Unwrap returns the underlying error.
func (e *HTTPError) Unwrap() error { return e.Cause }

// NewHTTPError creates an HTTP error.
func NewHTTPError(msg string, cause error) *HTTPError {
	return &HTTPError{Msg: msg, Cause: cause}
}

// HTTPErrorf creates an HTTP error with a formatted message.
func HTTPErrorf(format string, args ...any) *HTTPError {
	return &HTTPError{Msg: fmt.Sprintf(format, args...)}
}
