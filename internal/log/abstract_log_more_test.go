package log

import (
	"sync"
	"testing"
)

func TestAbstractLogIsTraceEnabled(t *testing.T) {
	var enabled bool
	a := &AbstractLog{IsEnabledFn: func(level Level) bool {
		return level == LevelTrace && enabled
	}}
	if a.IsTraceEnabled() {
		t.Error("IsTraceEnabled should be false when IsEnabledFn returns false")
	}
	enabled = true
	if !a.IsTraceEnabled() {
		t.Error("IsTraceEnabled should be true when IsEnabledFn returns true")
	}
}

func TestAbstractLogIsErrorEnabled(t *testing.T) {
	a := &AbstractLog{IsEnabledFn: func(level Level) bool {
		return level == LevelError
	}}
	if !a.IsErrorEnabled() {
		t.Error("IsErrorEnabled should be true for LevelError")
	}
}

func TestAbstractLogTraceAndError(t *testing.T) {
	var mu sync.Mutex
	type call struct {
		level Level
	}
	var calls []call
	a := &AbstractLog{
		IsEnabledFn: func(level Level) bool { return true },
		Core: func(level Level, err error, format string, args ...any) {
			mu.Lock()
			calls = append(calls, call{level: level})
			mu.Unlock()
		},
	}
	a.Trace("msg", 1)
	a.Error("msg", 2)

	if len(calls) != 2 {
		t.Fatalf("expected 2 calls, got %d", len(calls))
	}
	if calls[0].level != LevelTrace {
		t.Errorf("Trace: expected LevelTrace, got %v", calls[0].level)
	}
	if calls[1].level != LevelError {
		t.Errorf("Error: expected LevelError, got %v", calls[1].level)
	}
}
