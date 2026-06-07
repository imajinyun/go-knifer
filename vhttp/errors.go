package vhttp

import (
	knifer "github.com/imajinyun/go-knifer"
	httpx "github.com/imajinyun/go-knifer/internal/httpx/http"
)

// NewError delegates to the internal httpx implementation.
func NewError(msg string, cause error) *Error {
	return httpx.NewHTTPError(msg, cause)
}

// NewErrorWithCode creates an HTTP error with an explicit go-knifer code.
func NewErrorWithCode(code knifer.ErrCode, msg string, cause error) *Error {
	return httpx.NewHTTPErrorWithCode(code, msg, cause)
}

// Errorf delegates to the internal httpx implementation.
func Errorf(format string, args ...any) *Error {
	return httpx.HTTPErrorf(format, args...)
}

// ErrorfWithCode creates an HTTP error with an explicit code and formatted message.
func ErrorfWithCode(code knifer.ErrCode, format string, args ...any) *Error {
	return httpx.HTTPErrorfWithCode(code, format, args...)
}
