package http

import "testing"

func TestGlobalIgnoreEOFError(t *testing.T) {
	old := IsIgnoreEOFError()
	defer SetIgnoreEOFError(old)

	SetIgnoreEOFError(false)
	if IsIgnoreEOFError() {
		t.Fatal("expected false")
	}
}
