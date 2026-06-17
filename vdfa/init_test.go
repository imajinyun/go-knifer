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

func TestFacadePackageMatcherAccessors(t *testing.T) {
	InitWithOptions([]string{"foo", "foobar"})
	if !IsInited() {
		t.Fatal("IsInited should report initialized matcher")
	}
	first, ok := GetFoundFirst("foo foobar")
	if !ok || first.Word != "foo" || first.Start != 0 || first.End != 2 {
		t.Fatalf("GetFoundFirst = %#v, %v", first, ok)
	}
	all := GetFoundAll("foo foobar")
	if len(all) != 2 || all[0].FoundWord != "foo" || all[1].FoundWord != "foo" {
		t.Fatalf("GetFoundAll = %#v", all)
	}
	denseGreedy := GetFoundAllMode("foo foobar", true, true)
	if len(denseGreedy) != 3 || denseGreedy[2].FoundWord != "foobar" {
		t.Fatalf("GetFoundAllMode dense/greedy = %#v", denseGreedy)
	}
	got := FilterMode("foo", true, func(word FoundWord) string { return "[" + word.FoundWord + "]" })
	if got != "[foo]" {
		t.Fatalf("FilterMode = %q", got)
	}
}

func TestFacadeSetCharFilter(t *testing.T) {
	Init([]string{"ab"})
	SetCharFilter(func(r rune) bool { return r != '-' })
	if !Contains("a-b") {
		t.Fatal("SetCharFilter should update package-level matcher filter")
	}
	SetCharFilter(nil)
	if !Contains("a-b") {
		t.Fatal("SetCharFilter(nil) should leave current filter unchanged")
	}
}
