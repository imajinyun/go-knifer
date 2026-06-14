package vdfa

import "testing"

func TestFacadeWordTree(t *testing.T) {
	tree := NewWordTree().AddWords("foo", "foobar")
	got := tree.MatchAllMode("foo foobar", -1, true, true)
	if len(got) != 3 || got[0] != "foo" || got[1] != "foo" || got[2] != "foobar" {
		t.Fatalf("MatchAllMode() = %#v", got)
	}
}

func TestFacadeWordTreeOptions(t *testing.T) {
	tree := NewWordTreeWithOptions(WithCharFilter(func(r rune) bool { return r != '-' })).AddWord("t-io")
	if got := tree.MatchAll("tio"); len(got) != 1 || got[0] != "tio" {
		t.Fatalf("NewWordTreeWithOptions MatchAll() = %#v", got)
	}
}
