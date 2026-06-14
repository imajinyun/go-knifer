package vcron_test

import (
	"testing"

	"github.com/imajinyun/go-knifer/vcron"
)

func TestFacadePatternParse(t *testing.T) {
	p, err := vcron.NewPattern("0 0 * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil pattern")
	}
}

func TestFacadePatternParseInvalid(t *testing.T) {
	_, err := vcron.NewPattern("invalid")
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestFacadePartConstants(t *testing.T) {
	if err := vcron.PartMinute.CheckValue(59); err != nil {
		t.Fatalf("PartMinute.CheckValue(59) = %v", err)
	}
	if err := vcron.PartMinute.CheckValue(60); err == nil {
		t.Fatal("PartMinute.CheckValue(60) should fail")
	}
	if !vcron.AlwaysTrueMatcher.Match(123) || vcron.AlwaysTrueMatcher.NextAfter(7) != 7 {
		t.Fatal("AlwaysTrueMatcher facade mismatch")
	}
}

func TestFacadeMustPattern(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for invalid pattern")
		}
	}()
	vcron.MustNewPattern("bad")
}
