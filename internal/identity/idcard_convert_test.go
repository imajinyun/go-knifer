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

func TestConvertRejectsMalformedInputs(t *testing.T) {
	if got, ok := Convert15To18("13050367040100A"); ok || got != "" {
		t.Fatalf("Convert15To18(non-digit) = %q, %v; want empty false", got, ok)
	}
	if got, ok := Convert15To18WithOptions("130503670401001", WithDigitsMatcher(func(string) bool { return false })); ok || got != "" {
		t.Fatalf("Convert15To18WithOptions(rejected) = %q, %v; want empty false", got, ok)
	}
	if got, ok := Convert18To15("110105194912310021"); ok || got != "" {
		t.Fatalf("Convert18To15(invalid check) = %q, %v; want empty false", got, ok)
	}
}
