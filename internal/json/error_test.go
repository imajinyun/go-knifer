package json

import (
	"errors"
	"testing"

	knifer "github.com/imajinyun/go-knifer"
)

func TestJSONErrorMatchesErrCode(t *testing.T) {
	if !errors.Is(NewJSONError("empty path"), knifer.ErrCodeInvalidInput) {
		t.Fatal("NewJSONError should match knifer.ErrCodeInvalidInput")
	}
	wrapped := WrapJSONError(errors.New("eof"), "parse failed")
	if !errors.Is(wrapped, knifer.ErrCodeInvalidInput) {
		t.Fatal("WrapJSONError should match knifer.ErrCodeInvalidInput")
	}
	if !errors.Is(wrapped, errors.Unwrap(wrapped)) {
		t.Fatal("WrapJSONError should preserve the cause chain")
	}
}
