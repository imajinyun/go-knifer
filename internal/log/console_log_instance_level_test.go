package log

import (
	"bytes"
	"strings"
	"testing"
)

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
