package log

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func newTestConsoleLog(name string) (*ConsoleLog, *bytes.Buffer, *bytes.Buffer) {
	c := NewConsoleLog(name)
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	c.SetOutput(out, errOut)
	return c, out, errOut
}

func TestConsoleLogLevels(t *testing.T) {
	prevLevel := GetConsoleLevel()
	SetConsoleLevel(LevelDebug)
	defer SetConsoleLevel(prevLevel)

	c, out, errOut := newTestConsoleLog("test.console")
	c.Debug("debug msg")
	c.Infof("user={}", "alice")
	c.Warnf("warn-{}-{}", 1, 2)
	c.Errorf("err=%d", 7)

	stdoutText := out.String()
	stderrText := errOut.String()

	if !strings.Contains(stdoutText, "[DEBUG]") || !strings.Contains(stdoutText, "debug msg") {
		t.Errorf("expected debug log in stdout, got %q", stdoutText)
	}
	if !strings.Contains(stdoutText, "[INFO ]") || !strings.Contains(stdoutText, "user=alice") {
		t.Errorf("expected info log with placeholder, got %q", stdoutText)
	}
	if !strings.Contains(stderrText, "[WARN ]") || !strings.Contains(stderrText, "warn-1-2") {
		t.Errorf("expected warn log in stderr, got %q", stderrText)
	}
	if !strings.Contains(stderrText, "[ERROR]") || !strings.Contains(stderrText, "err=7") {
		t.Errorf("expected error log in stderr, got %q", stderrText)
	}
}

func TestConsoleLogFiltering(t *testing.T) {
	prevLevel := GetConsoleLevel()
	SetConsoleLevel(LevelWarn)
	defer SetConsoleLevel(prevLevel)

	c, out, errOut := newTestConsoleLog("test.filter")
	c.Debug("should be filtered")
	c.Info("should be filtered")
	c.Warn("kept warn")

	if out.Len() != 0 {
		t.Errorf("expected no debug/info output, got %q", out.String())
	}
	if !strings.Contains(errOut.String(), "kept warn") {
		t.Errorf("expected warn output, got %q", errOut.String())
	}
	if !c.IsWarnEnabled() || c.IsDebugEnabled() {
		t.Errorf("level checks wrong: warn=%v debug=%v", c.IsWarnEnabled(), c.IsDebugEnabled())
	}
}

func TestConsoleLogInstanceLevelOverridesGlobalLevel(t *testing.T) {
	prevLevel := GetConsoleLevel()
	SetConsoleLevel(LevelError)
	defer SetConsoleLevel(prevLevel)

	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	c := NewConsoleLogWithOptions("test.instance.level", WithLogOutput(out, errOut), WithLogLevel(LevelInfo))
	c.Debug("filtered debug")
	c.Info("kept info")

	if strings.Contains(out.String(), "filtered debug") {
		t.Fatalf("debug should be filtered by instance level: %q", out.String())
	}
	if !strings.Contains(out.String(), "kept info") {
		t.Fatalf("info should be enabled by instance level despite global error threshold: %q", out.String())
	}
	if c.IsDebugEnabled() || !c.IsInfoEnabled() {
		t.Fatalf("instance level checks debug=%v info=%v", c.IsDebugEnabled(), c.IsInfoEnabled())
	}
}

func TestConsoleLogWithError(t *testing.T) {
	prevLevel := GetConsoleLevel()
	SetConsoleLevel(LevelDebug)
	defer SetConsoleLevel(prevLevel)

	c, _, errOut := newTestConsoleLog("test.err")
	c.LogE(LevelError, errors.New("boom"), "operation {} failed", "save")
	got := errOut.String()
	if !strings.Contains(got, "operation save failed") || !strings.Contains(got, "error: boom") {
		t.Errorf("error formatted output unexpected: %q", got)
	}
}

func TestConsoleLogWithOptions(t *testing.T) {
	prevLevel := GetConsoleLevel()
	SetConsoleLevel(LevelDebug)
	defer SetConsoleLevel(prevLevel)

	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	fixed := time.Date(2024, 2, 3, 4, 5, 6, 0, time.UTC)
	c := NewConsoleLogWithOptions("test.options",
		WithLogClock(func() time.Time { return fixed }),
		WithLogTimeLayout(time.RFC3339),
		WithLogOutput(out, errOut),
	)
	c.Info("hello")
	c.Warn("careful")

	if !strings.Contains(out.String(), "2024-02-03T04:05:06Z") || !strings.Contains(out.String(), "hello") {
		t.Fatalf("custom clock/layout/output not applied to stdout: %q", out.String())
	}
	if !strings.Contains(errOut.String(), "2024-02-03T04:05:06Z") || !strings.Contains(errOut.String(), "careful") {
		t.Fatalf("custom clock/layout/output not applied to stderr: %q", errOut.String())
	}
}
