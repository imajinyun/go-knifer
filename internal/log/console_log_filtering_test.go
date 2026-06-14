package log

import (
	"strings"
	"testing"
)

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
