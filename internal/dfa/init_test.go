package dfa

import "testing"

func TestInitWithOptions(t *testing.T) {
	InitWithOptions([]string{"t-io"}, WithCharFilter(func(r rune) bool { return r != '-' }))
	if !Contains("tio") {
		t.Fatal("InitWithOptions should apply custom char filter")
	}
	InitStringWithOptions("a-b", 0, WithCharFilter(func(r rune) bool { return r != '-' }))
	if !Contains("ab") {
		t.Fatal("InitStringWithOptions should apply custom char filter")
	}
}

func TestMatcherOptionsBypassPackageMatcher(t *testing.T) {
	Init([]string{"global"})
	matcher := NewWordTree().AddWords("local")
	if ContainsWithOptions("local", WithMatcher(matcher)) != true {
		t.Fatal("ContainsWithOptions should use provided matcher")
	}
	if Contains("local") {
		t.Fatal("per-call matcher should not mutate package matcher")
	}
	if got := FilterWithOptions("a local word", WithMatcher(matcher)); got != "a ***** word" {
		t.Fatalf("FilterWithOptions = %q", got)
	}
	found, ok := GetFoundFirstWithOptions("a local word", WithMatcherWords([]string{"local"}))
	if !ok || found.Word != "local" {
		t.Fatalf("GetFoundFirstWithOptions = %#v ok=%v", found, ok)
	}
}
