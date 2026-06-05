package errx

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-multierror"
	knifer "github.com/imajinyun/go-knifer"
)

type stackError struct{ stack string }

func (e stackError) Error() string { return "stack error" }
func (e stackError) Stack() string { return e.stack }

func TestGetStackUsesAttachedStackWhenAvailable(t *testing.T) {
	const want = "attached stack"
	if got := GetStack(stackError{stack: want}); got != want {
		t.Fatalf("GetStack() = %q, want %q", got, want)
	}
}

func TestGetStackFallsBackToRuntimeStack(t *testing.T) {
	got := GetStack(errors.New("plain"))
	if !strings.Contains(got, "goroutine") || !strings.Contains(got, "TestGetStackFallsBackToRuntimeStack") {
		t.Fatalf("GetStack() fallback does not look like a runtime stack: %q", got)
	}
	if got := GetStack(nil); got != "" {
		t.Fatalf("GetStack(nil) = %q, want empty", got)
	}
}

func TestErrorIsHandlesNestedMultierror(t *testing.T) {
	target := errors.New("target")
	nested := multierror.Append(nil, errors.New("other"), fmt.Errorf("wrapped: %w", target))
	err := multierror.Append(nil, errors.New("top"), nested)

	if !ErrorIs(err, target) {
		t.Fatalf("ErrorIs() = false, want true for nested multierror")
	}
	if ErrorIs(err, errors.New("missing")) {
		t.Fatal("ErrorIs() = true for an unrelated error")
	}
	if !ErrorIs(nil, nil) {
		t.Fatal("ErrorIs(nil, nil) should be true")
	}
	if ErrorIs(err, nil) {
		t.Fatal("ErrorIs(non-nil, nil) should be false")
	}
}

func TestPanicErrorPreservesErrorValues(t *testing.T) {
	want := errors.New("panic error")
	got := panicError(want)
	if !errors.Is(got, want) {
		t.Fatalf("panicError(error) = %v, want wrapping original error", got)
	}
	var pe *PanicError
	if !errors.As(got, &pe) {
		t.Fatalf("panicError(error) type = %T, want *PanicError", got)
	}
	if pe.Value != want || pe.Cause != want {
		t.Fatalf("PanicError value/cause = (%v, %v), want original error", pe.Value, pe.Cause)
	}
	if pe.Stack() == "" {
		t.Fatal("PanicError should capture a stack")
	}
	if got := panicError("panic string"); got == nil || got.Error() != "panic string" {
		t.Fatalf("panicError(string) = %v, want converted error", got)
	} else if !errors.Is(got, knifer.ErrCodeInternal) {
		t.Fatalf("panicError(string) = %v, want ErrCodeInternal", got)
	}

	coded := panicError(knifer.NewError(knifer.ErrCodeInvalidInput, "bad input"))
	if !errors.Is(coded, knifer.ErrCodeInvalidInput) {
		t.Fatalf("panicError(coded error) = %v, want ErrCodeInvalidInput", coded)
	}
	if code, ok := knifer.CodeOf(coded); !ok || code != knifer.ErrCodeInvalidInput {
		t.Fatalf("CodeOf(panicError(coded error)) = %q, %v; want invalid input", code, ok)
	}
}

func TestPanicErrorNilReceiver(t *testing.T) {
	var pe *PanicError
	if pe.Error() != "<nil>" || pe.Unwrap() != nil || pe.Stack() != "" || pe.ErrorCode() != "" {
		t.Fatalf("nil PanicError methods returned unexpected values")
	}
}
