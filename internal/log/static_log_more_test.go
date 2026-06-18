package log

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestStaticLogSimpleMethods(t *testing.T) {
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

	Trace("trace-msg")
	Debug("debug-msg")
	Info("info-msg")
	Warn("warn-msg")
	ErrorLog("error-msg")

	stdout := out.String()
	stderr := errOut.String()
	for _, want := range []string{"trace-msg", "debug-msg", "info-msg"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("stdout missing %q in %q", want, stdout)
		}
	}
	for _, want := range []string{"warn-msg", "error-msg"} {
		if !strings.Contains(stderr, want) {
			t.Errorf("stderr missing %q in %q", want, stderr)
		}
	}
}

func TestStaticLogWithOptionsVariants(t *testing.T) {
	out := &bytes.Buffer{}
	SetFactory(LogFactoryFunc(func(name string) Log {
		c := NewConsoleLog(name)
		// Only test stdout.
		var buf bytes.Buffer
		c.SetOutput(&buf, &buf)
		return c
	}))
	defer SetFactory(LogFactoryFunc(func(name string) Log { return NewConsoleLog(name) }))

	prevLevel := GetConsoleLevel()
	SetConsoleLevel(LevelTrace)
	defer SetConsoleLevel(prevLevel)

	// Test each WithOptions variant.
	// Use the same buffer for both stdout and stderr so that warn/error level messages are captured.
	opts := []LoggerOption{WithLoggerConsoleOptions(
		WithLogOutput(out, out),
	)}

	TraceWithOptions(opts, "trace-opt")
	TracefWithOptions(opts, "tracef-%s", "opt")
	DebugWithOptions(opts, "debug-opt")
	DebugfWithOptions(opts, "debugf-%s", "opt")
	InfofWithOptions(opts, "infof-%s", "opt")
	WarnWithOptions(opts, "warn-opt")
	WarnfWithOptions(opts, "warnf-%s", "opt")
	ErrorLogWithOptions(opts, "error-opt")
	ErrorfWithOptions(opts, "errorf-%s", "opt")
	LogAtWithOptions(opts, LevelInfo, "logat-%s", "opt")
	LogAtEWithOptions(opts, LevelError, errors.New("e"), "logate-%s", "opt")

	output := out.String()
	for _, want := range []string{
		"trace-opt", "tracef-opt",
		"debug-opt", "debugf-opt",
		"infof-opt",
		"warn-opt", "warnf-opt",
		"error-opt", "errorf-opt",
		"logat-opt", "logate-opt",
	} {
		if !strings.Contains(output, want) {
			t.Errorf("output missing %q in %q", want, output)
		}
	}
}
