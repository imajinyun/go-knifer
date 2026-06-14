package log

import (
	"strings"
	"testing"
)

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
