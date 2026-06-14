package vdfa

import "testing"

func TestFacadeMatcherOptions(t *testing.T) {
	Init([]string{"global"})
	matcher := NewWordTree().AddWord("local")
	if !ContainsWithOptions("a local word", WithMatcher(matcher)) {
		t.Fatal("ContainsWithOptions should use provided matcher")
	}
	if Contains("local") {
		t.Fatal("per-call matcher should not mutate package matcher")
	}
	if got := FilterWithOptions("a local word", WithMatcherWords([]string{"local"})); got != "a ***** word" {
		t.Fatalf("FilterWithOptions = %q", got)
	}
}
