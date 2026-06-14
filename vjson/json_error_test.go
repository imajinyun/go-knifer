package vjson_test

import (
	"errors"
	"testing"

	knifer "github.com/imajinyun/go-knifer"
	"github.com/imajinyun/go-knifer/vjson"
)

func TestFacadeErrorNameWithoutJSONPrefix(t *testing.T) {
	_, err := vjson.ParseObj(`[1,2]`)
	var jsonErr *vjson.Error
	if !errors.As(err, &jsonErr) {
		t.Fatalf("ParseObj() error type = %T, want *vjson.Error", err)
	}
	if !errors.Is(err, knifer.ErrCodeInvalidInput) {
		t.Fatalf("errors.Is(err, ErrCodeInvalidInput) = false: %v", err)
	}
	code, ok := knifer.CodeOf(err)
	if !ok || code != knifer.ErrCodeInvalidInput {
		t.Fatalf("CodeOf(err) = %q, %v; want invalid input", code, ok)
	}

	_, err = vjson.XMLToJSON(`<root><unclosed></root>`)
	if !errors.Is(err, knifer.ErrCodeInvalidInput) {
		t.Fatalf("XMLToJSON malformed XML code = %v, want invalid input", err)
	}
}
