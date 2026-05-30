package vrand

import (
	"strings"
	"testing"
)

func TestRandFacade(t *testing.T) {
	for i := 0; i < 50; i++ {
		if got := IntRange(10, 20); got < 10 || got >= 20 {
			t.Fatalf("IntRange out of bounds: %d", got)
		}
	}
	if Int(0) != 0 || Long() < 0 || Float() < 0 || Float() >= 1 {
		t.Fatal("basic random helpers failed")
	}
	_ = Bool()
	if len(Bytes(8)) != 8 {
		t.Fatal("Bytes length failed")
	}
	if s := String(8); len(s) != 8 {
		t.Fatalf("String length failed: %q", s)
	}
	if s := Numbers(6); len(s) != 6 {
		t.Fatalf("Numbers length failed: %q", s)
	}
	upper := StringUpper(8)
	for _, r := range upper {
		if !strings.ContainsRune(BaseCharNumberUC, r) {
			t.Fatalf("StringUpper charset failed: %q", upper)
		}
	}
	if s := StringFrom("ab", 4); len(s) != 4 {
		t.Fatalf("StringFrom failed: %q", s)
	}
	if got := Ele([]string{"x"}); got != "x" {
		t.Fatalf("Ele failed: %q", got)
	}
}
