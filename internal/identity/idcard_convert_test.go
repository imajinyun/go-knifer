package identity

import "testing"

func TestConvert15To18(t *testing.T) {
	got, ok := Convert15To18("130503670401001")
	if !ok || got != "130503196704010016" {
		t.Fatalf("Convert15To18() = %q, %v", got, ok)
	}

	got15, ok := Convert18To15(got)
	if !ok || got15 != "130503670401001" {
		t.Fatalf("Convert18To15() = %q, %v", got15, ok)
	}
}
