package boolean

import "testing"

func TestBoolean(t *testing.T) {
	if !Negate(false) || Negate(true) {
		t.Fatalf("Negate failed")
	}
	if ToInt(true) != 1 || ToInt(false) != 0 {
		t.Fatalf("ToInt failed")
	}
	if !And(true, true) || And(true, false) {
		t.Fatalf("And failed")
	}
	if !Or(false, true) || Or(false, false) {
		t.Fatalf("Or failed")
	}
}
