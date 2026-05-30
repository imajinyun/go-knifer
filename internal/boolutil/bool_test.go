package boolutil

import "testing"

func TestBooleanUtil(t *testing.T) {
	if !BoolNegate(false) || BoolNegate(true) {
		t.Fatalf("Negate failed")
	}
	if BoolToInt(true) != 1 || BoolToInt(false) != 0 {
		t.Fatalf("BoolToInt failed")
	}
	if !BoolAnd(true, true) || BoolAnd(true, false) {
		t.Fatalf("BoolAnd failed")
	}
	if !BoolOr(false, true) || BoolOr(false, false) {
		t.Fatalf("BoolOr failed")
	}
}
