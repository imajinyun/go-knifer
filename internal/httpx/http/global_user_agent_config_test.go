package http

import "testing"

func TestGlobalUserAgent(t *testing.T) {
	old := GetGlobalUserAgent()
	defer SetGlobalUserAgent(old)

	SetGlobalUserAgent("gokit-test/1.0")
	if got := GetGlobalUserAgent(); got != "gokit-test/1.0" {
		t.Fatalf("ua: %q", got)
	}
}
