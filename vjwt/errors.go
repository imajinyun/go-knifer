package vjwt

import jwtimpl "github.com/imajinyun/go-knifer/internal/jwt"

// NewJWTError creates a JWT module error with invalid-input code.
func NewJWTError(msg string) *JWTError {
	return jwtimpl.NewJWTError(msg)
}

// JWTErrorf creates a formatted JWT module error with invalid-input code.
func JWTErrorf(format string, args ...any) *JWTError {
	return jwtimpl.JWTErrorf(format, args...)
}
