package shared

import "testing"

func TestHeaderString(t *testing.T) {
	if got := HeaderContentDisposition.String(); got != "Content-Disposition" {
		t.Fatalf("HeaderContentDisposition.String = %q", got)
	}
}
