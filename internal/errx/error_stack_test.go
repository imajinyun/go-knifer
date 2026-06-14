package errx

import (
	"errors"
	"strings"
	"testing"
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
