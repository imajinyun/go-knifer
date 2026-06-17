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

func TestFacadeWordTreeBoundaries(t *testing.T) {
	var nilTree *WordTree
	if !nilTree.IsEmpty() || nilTree.IsMatch("anything") {
		t.Fatal("nil WordTree should be empty and not match")
	}

	tree := NewWordTree().AddWords("foo", "foo", "foobar")
	if got, ok := tree.Match("foo-bar"); !ok || got != "foo" {
		t.Fatalf("Match() = %q, %v", got, ok)
	}
	if got := tree.MatchAllLimit("foo foobar foo", 2); len(got) != 2 || got[0] != "foo" || got[1] != "foo" {
		t.Fatalf("MatchAllLimit() = %#v", got)
	}
	tree.Clear()
	if !tree.IsEmpty() || len(tree.MatchAll("foo")) != 0 {
		t.Fatal("Clear should remove all words")
	}
}

func TestFacadeStopRunesAndProcessor(t *testing.T) {
	if !IsStopChar(' ') || IsStopChar('A') || !IsNotStopChar('A') || IsNotStopChar(' ') {
		t.Fatal("stop rune facade helpers returned inconsistent values")
	}
	if got := DefaultProcessor(FoundWord{FoundWord: "敏感"}); got != "**" {
		t.Fatalf("DefaultProcessor = %q, want **", got)
	}
}
