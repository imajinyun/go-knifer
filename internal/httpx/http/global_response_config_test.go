package http

import "testing"

func TestGlobalMaxResponseBytes(t *testing.T) {
	old := GetGlobalMaxResponseBytes()
	defer SetGlobalMaxResponseBytes(old)

	SetGlobalMaxResponseBytes(123)
	if got := GetGlobalMaxResponseBytes(); got != 123 {
		t.Fatalf("max response bytes: %d", got)
	}
}
