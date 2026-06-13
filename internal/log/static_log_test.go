package log

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestLoggerConsoleOptionsForStaticLog(t *testing.T) {
	prevLevel := GetConsoleLevel()
	SetConsoleLevel(LevelInfo)
	defer SetConsoleLevel(prevLevel)

	out := &bytes.Buffer{}
	fixed := time.Date(2024, 6, 7, 8, 9, 10, 0, time.UTC)
	InfoWithOptions([]LoggerOption{WithLoggerConsoleOptions(
		WithLogOutput(out, &bytes.Buffer{}),
		WithLogClock(func() time.Time { return fixed }),
		WithLogTimeLayout(time.RFC3339),
	)}, "static-option")
	if !strings.Contains(out.String(), "2024-06-07T08:09:10Z") || !strings.Contains(out.String(), "static-option") {
		t.Fatalf("static logger options not applied: %q", out.String())
	}
}

func TestStaticLogPipeline(t *testing.T) {
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	SetFactory(LogFactoryFunc(func(name string) Log {
		c := NewConsoleLog(name)
		c.SetOutput(out, errOut)
		return c
	}))
	defer SetFactory(LogFactoryFunc(func(name string) Log { return NewConsoleLog(name) }))

	prevLevel := GetConsoleLevel()
	SetConsoleLevel(LevelTrace)
	defer SetConsoleLevel(prevLevel)

	Tracef("trace {}", 1)
	Debugf("debug {}", 2)
	Infof("info {}", 3)
	Warnf("warn {}", 4)
	Errorf("error {}", 5)

	stdout := out.String()
	stderr := errOut.String()
	for _, want := range []string{"trace 1", "debug 2", "info 3"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("stdout missing %q in %q", want, stdout)
		}
	}
	for _, want := range []string{"warn 4", "error 5"} {
		if !strings.Contains(stderr, want) {
			t.Errorf("stderr missing %q in %q", want, stderr)
		}
	}
}

func TestLogAtAndLogAtE(t *testing.T) {
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	SetFactory(LogFactoryFunc(func(name string) Log {
		c := NewConsoleLog(name)
		c.SetOutput(out, errOut)
		return c
	}))
	defer SetFactory(LogFactoryFunc(func(name string) Log { return NewConsoleLog(name) }))

	prevLevel := GetConsoleLevel()
	SetConsoleLevel(LevelDebug)
	defer SetConsoleLevel(prevLevel)

	LogAt(LevelInfo, "hello {}", "world")
	LogAtE(LevelError, errors.New("oops"), "boom {}", "now")

	if !strings.Contains(out.String(), "hello world") {
		t.Errorf("LogAt info missing: %q", out.String())
	}
	if !strings.Contains(errOut.String(), "boom now") || !strings.Contains(errOut.String(), "error: oops") {
		t.Errorf("LogAtE error missing: %q", errOut.String())
	}
}
