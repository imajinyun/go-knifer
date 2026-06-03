package knifer_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/imajinyun/go-knifer"
)

func TestErrorCodeMatching(t *testing.T) {
	err := knifer.NewError(knifer.ErrCodeInvalidInput, "url empty")
	if !errors.Is(err, knifer.ErrCodeInvalidInput) {
		t.Fatal("errors.Is should match the error code")
	}
	if errors.Is(err, knifer.ErrCodeNotFound) {
		t.Fatal("errors.Is should not match a different code")
	}
}

func TestErrorWrapPreservesChain(t *testing.T) {
	cause := errors.New("disk full")
	err := knifer.WrapError(knifer.ErrCodeInternal, "write failed", cause)

	if !errors.Is(err, knifer.ErrCodeInternal) {
		t.Fatal("wrapped error should match its code")
	}
	if !errors.Is(err, cause) {
		t.Fatal("wrapped error should preserve the cause chain")
	}
	if got := errors.Unwrap(err); got != cause {
		t.Fatalf("Unwrap = %v, want %v", got, cause)
	}
}

func TestErrorAs(t *testing.T) {
	err := fmt.Errorf("context: %w", knifer.Errorf(knifer.ErrCodeTimeout, "deadline %ds", 3))
	var ke *knifer.Error
	if !errors.As(err, &ke) {
		t.Fatal("errors.As should extract *knifer.Error")
	}
	if ke.Code != knifer.ErrCodeTimeout {
		t.Fatalf("Code = %q, want %q", ke.Code, knifer.ErrCodeTimeout)
	}
}

func TestErrorString(t *testing.T) {
	if got := knifer.NewError(knifer.ErrCodeNotFound, "missing").Error(); got != "GK_NOT_FOUND: missing" {
		t.Fatalf("Error() = %q", got)
	}
	wrapped := knifer.WrapError(knifer.ErrCodeInternal, "boom", errors.New("root"))
	if got := wrapped.Error(); got != "GK_INTERNAL: boom: root" {
		t.Fatalf("Error() = %q", got)
	}
}
