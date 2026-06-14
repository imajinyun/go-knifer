package vdfa

import "testing"

func TestFacadePackageMatcher(t *testing.T) {
	Init([]string{"bad", "badword"})
	if !Contains("a badword") {
		t.Fatalf("expected package matcher to contain word")
	}
	got := Filter("a badword")
	if got != "a *******" {
		t.Fatalf("Filter() = %q", got)
	}
}

func TestFacadeInitWithOptions(t *testing.T) {
	InitWithOptions([]string{"a-b"}, WithCharFilter(func(r rune) bool { return r != '-' }))
	if !Contains("ab") {
		t.Fatal("InitWithOptions should apply custom char filter")
	}
	InitStringWithOptions("x-y", DefaultSeparator, WithCharFilter(func(r rune) bool { return r != '-' }))
	if !Contains("xy") {
		t.Fatal("InitStringWithOptions should apply custom char filter")
	}
}
