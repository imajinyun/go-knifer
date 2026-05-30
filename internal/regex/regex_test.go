package regex

import "testing"

func TestRegex(t *testing.T) {
	if !ReMatch(`^\d+$`, "123") || ReMatch(`^\d+$`, "12a") {
		t.Fatalf("ReMatch failed")
	}
	if ReFind(`\d+`, "ab123cd") != "123" {
		t.Fatalf("ReFind failed")
	}
	all := ReFindAll(`\d+`, "a1b22c333")
	if len(all) != 3 || all[2] != "333" {
		t.Fatalf("ReFindAll failed: %v", all)
	}
	if ReReplace(`\d`, "a1b2", "*") != "a*b*" {
		t.Fatalf("ReReplace failed")
	}
}
