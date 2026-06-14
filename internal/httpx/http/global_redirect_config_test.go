package http

import "testing"

func TestGlobalFollowRedirects(t *testing.T) {
	old := GetGlobalFollowRedirects()
	defer SetGlobalFollowRedirects(old)

	SetGlobalFollowRedirects(false)
	if GetGlobalFollowRedirects() {
		t.Fatal("expected false")
	}
}

func TestGlobalMaxRedirects(t *testing.T) {
	old := GetGlobalMaxRedirects()
	defer SetGlobalMaxRedirects(old)

	SetGlobalMaxRedirects(3)
	if got := GetGlobalMaxRedirects(); got != 3 {
		t.Fatalf("max: %d", got)
	}
}
