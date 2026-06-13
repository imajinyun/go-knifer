package log

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestConsoleColorLog(t *testing.T) {
	prevLevel := GetConsoleLevel()
	SetConsoleLevel(LevelDebug)
	defer SetConsoleLevel(prevLevel)

	c := NewConsoleColorLog("test.color")
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	c.SetOutput(out, errOut)

	c.Info("hi")
	c.Warn("careful")

	if !strings.Contains(out.String(), "hi") || !strings.Contains(out.String(), "\033[") {
		t.Errorf("color info expected ANSI codes, got %q", out.String())
	}
	if !strings.Contains(errOut.String(), "careful") {
		t.Errorf("color warn output unexpected: %q", errOut.String())
	}
}

func TestConsoleColorLogWithOptions(t *testing.T) {
	prevLevel := GetConsoleLevel()
	SetConsoleLevel(LevelDebug)
	defer SetConsoleLevel(prevLevel)

	out := &bytes.Buffer{}
	fixed := time.Date(2024, 3, 4, 5, 6, 7, 0, time.UTC)
	c := NewConsoleColorLogWithOptions("test.color.options",
		WithLogClock(func() time.Time { return fixed }),
		WithLogTimeLayout("15:04:05"),
		WithLogOutput(out, &bytes.Buffer{}),
	)
	c.Info("hi")
	got := out.String()
	if !strings.Contains(got, "05:06:07") || !strings.Contains(got, "\033[") || !strings.Contains(got, "hi") {
		t.Fatalf("color options not applied: %q", got)
	}
}

func TestConsoleColorLogInstanceColorFactory(t *testing.T) {
	prevLevel := GetConsoleLevel()
	SetConsoleLevel(LevelDebug)
	defer SetConsoleLevel(prevLevel)

	SetColorFactory(func(Level) string { return ansiBlue })
	defer SetColorFactory(defaultColorFactory)

	out := &bytes.Buffer{}
	called := false
	c := NewConsoleColorLogWithOptions("test.color.instance",
		WithLogOutput(out, &bytes.Buffer{}),
		WithLogColorFactory(func(Level) string {
			called = true
			return ansiCyan
		}),
	)
	c.Info("hi")

	if !called {
		t.Fatal("instance color factory was not called")
	}
	if got := out.String(); !strings.Contains(got, ansiCyan) || strings.Contains(got, ansiBlue) {
		t.Fatalf("instance color factory should override global factory, got %q", got)
	}
}

func TestSetColorFactory(t *testing.T) {
	called := false
	SetColorFactory(func(level Level) string {
		called = true
		return ansiBlue
	})
	defer SetColorFactory(defaultColorFactory)

	prevLevel := GetConsoleLevel()
	SetConsoleLevel(LevelDebug)
	defer SetConsoleLevel(prevLevel)

	c := NewConsoleColorLog("test.colorfactory")
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	c.SetOutput(out, errOut)
	c.Info("x")

	if !called {
		t.Error("custom color factory was not called")
	}
	if !strings.Contains(out.String(), ansiBlue) {
		t.Errorf("custom color expected, got %q", out.String())
	}
}
