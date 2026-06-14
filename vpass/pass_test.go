package vpass

import "testing"

func TestFacade(t *testing.T) {
	analysis := Analyze("G0-Knifer#Pass2026")
	if analysis.Score == 0 {
		t.Fatal("Analyze() returned zero score for strong password")
	}
	if analysis.Strength != StrengthVeryStrong {
		t.Fatalf("Analyze().Strength = %v, want %v", analysis.Strength, StrengthVeryStrong)
	}
	if Score("password") > 10 {
		t.Fatal("Score() did not cap common weak password")
	}
	if StrengthOf("password") != StrengthVeryWeak {
		t.Fatal("StrengthOf() did not classify common weak password")
	}
	if !IsStrong("G0-Knifer#Pass2026") {
		t.Fatal("IsStrong() = false for strong password")
	}
	if !IsWeak("12345") {
		t.Fatal("IsWeak() = false for weak password")
	}
	if StrengthUnknown.String() != "unknown" {
		t.Fatal("Strength.String() alias failed")
	}
}
