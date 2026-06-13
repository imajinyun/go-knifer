package shared

import "testing"

func TestMethodString(t *testing.T) {
	if got := MethodPatch.String(); got != "PATCH" {
		t.Fatalf("MethodPatch.String = %q", got)
	}
}
