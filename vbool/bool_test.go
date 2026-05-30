package vbool

import "testing"

func TestBoolFacade(t *testing.T) {
	if !Negate(false) || Negate(true) {
		t.Fatal("Negate failed")
	}
	if ToInt(true) != 1 || ToInt(false) != 0 {
		t.Fatal("ToInt failed")
	}
	if !And(true, true) || And(true, false) {
		t.Fatal("And failed")
	}
	if !Or(false, true) || Or(false, false) {
		t.Fatal("Or failed")
	}
}
