package log

import (
	"errors"
	"strings"
	"testing"
)

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
