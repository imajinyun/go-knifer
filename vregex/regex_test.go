package vregex

import "testing"

func TestRegexFacade(t *testing.T) {
	if !Match(`^\d+$`, "123") || Match(`^\d+$`, "12a") || Match(`(`, "x") {
		t.Fatal("Match failed")
	}
	if Find(`\d+`, "ab123cd") != "123" || Find(`(`, "x") != "" {
		t.Fatal("Find failed")
	}
	if all := FindAll(`\d+`, "a1b22c333"); len(all) != 3 || all[2] != "333" {
		t.Fatalf("FindAll failed")
	}
	if Replace(`\d`, "a1b2", "*") != "a*b*" || Replace(`(`, "x", "*") != "x" {
		t.Fatal("Replace failed")
	}
}
