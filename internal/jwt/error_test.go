package jwt

import (
	"errors"
	"testing"

	knifer "github.com/imajinyun/go-knifer"
)

func TestJWTErrorMatchesErrCode(t *testing.T) {
	if !errors.Is(NewJWTError("bad token"), knifer.ErrCodeInvalidInput) {
		t.Fatal("NewJWTError should match knifer.ErrCodeInvalidInput")
	}
	if !errors.Is(JWTErrorf("bad %s", "alg"), knifer.ErrCodeInvalidInput) {
		t.Fatal("JWTErrorf should match knifer.ErrCodeInvalidInput")
	}
	if errors.Is(NewJWTError("bad token"), knifer.ErrCodeNotFound) {
		t.Fatal("should not match an unrelated code")
	}
}
