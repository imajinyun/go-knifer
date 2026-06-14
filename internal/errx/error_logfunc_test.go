package errx

import (
	"context"
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestDefaultLogFuncCanBeConfiguredAndReset(t *testing.T) {
	ResetDefaultLogFunc()
	t.Cleanup(ResetDefaultLogFunc)

	want := errors.New("configured default logger")
	called := 0
	ConfigureDefaultLogFunc(func(ctx context.Context, level logrus.Level, err error, stack string, format string, args ...any) {
		called++
		if ctx == nil {
			t.Fatal("logger context is nil")
		}
		if level != logrus.ErrorLevel {
			t.Fatalf("logger level = %s, want error", level)
		}
		if !ErrorIs(err, want) {
			t.Fatalf("logger err = %v, want %v", err, want)
		}
		if format != "configured %s" || len(args) != 1 || args[0] != "logger" {
			t.Fatalf("logger format/args = %q/%v", format, args)
		}
	})
	if err := Recover(func() error { return want }, "configured %s", "logger"); !ErrorIs(err, want) {
		t.Fatalf("Recover() = %v, want %v", err, want)
	}
	if called != 1 {
		t.Fatalf("configured logger called %d times, want 1", called)
	}

	ResetDefaultLogFunc()
	called = 0
	if err := Recover(func() error { return want }, "reset logger"); !ErrorIs(err, want) {
		t.Fatalf("Recover() after reset = %v, want %v", err, want)
	}
	if called != 0 {
		t.Fatalf("configured logger called after reset %d times, want 0", called)
	}
}
