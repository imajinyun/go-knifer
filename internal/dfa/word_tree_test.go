package dfa

import (
	"reflect"
	"testing"
)

func TestWordTreeMatchModes(t *testing.T) {
	tree := buildTestTree()
	tests := []struct {
		name    string
		density bool
		greed   bool
		want    []string
	}{
		{name: "standard", want: []string{"大", "土^豆", "刚出锅"}},
		{name: "density", density: true, want: []string{"大", "土^豆", "刚出锅", "出锅"}},
		{name: "greed", greed: true, want: []string{"大", "土^豆", "刚出锅"}},
		{name: "density greed", density: true, greed: true, want: []string{"大", "大土^豆", "土^豆", "刚出锅", "出锅"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tree.MatchAllMode(sampleText, -1, tt.density, tt.greed)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("MatchAllMode() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestWordTreeFoundWordMetadata(t *testing.T) {
	tree := NewWordTree().AddWords("赵", "赵阿", "赵阿三")
	got := tree.MatchAllWords("赵阿三在做什么", -1, true, true)
	if len(got) != 3 {
		t.Fatalf("len = %d", len(got))
	}
	assertFoundWord(t, got[0], "赵", "赵", 0, 0)
	assertFoundWord(t, got[1], "赵阿", "赵阿", 0, 1)
	assertFoundWord(t, got[2], "赵阿三", "赵阿三", 0, 2)
}

func TestWordTreeOptions(t *testing.T) {
	tree := NewWordTreeWithOptions(WithCharFilter(func(r rune) bool { return r != '-' })).AddWord("t-io")
	if got := tree.MatchAll("tio"); !reflect.DeepEqual(got, []string{"tio"}) {
		t.Fatalf("NewWordTreeWithOptions MatchAll() = %#v", got)
	}
}

func TestAddWordWithFilteredRune(t *testing.T) {
	tree := NewWordTree().AddWord("hello(")
	if got := tree.MatchAllLimit("hello", -1); !reflect.DeepEqual(got, []string{"hello"}) {
		t.Fatalf("trailing filtered rune match = %#v", got)
	}

	tree = NewWordTree().AddWord("he(llo")
	if got := tree.MatchAllLimit("hello", -1); !reflect.DeepEqual(got, []string{"hello"}) {
		t.Fatalf("middle filtered rune match = %#v", got)
	}
}

func TestClear(t *testing.T) {
	tree := NewWordTree().AddWord("黑")
	if !contains(tree.MatchAll("黑大衣"), "黑") {
		t.Fatalf("expected initial match")
	}
	tree.Clear()
	tree.AddWords("黑大衣", "红色大衣")
	if !contains(tree.MatchAll("黑大衣"), "黑大衣") {
		t.Fatalf("expected 黑大衣 after clear")
	}
	if contains(tree.MatchAll("黑大衣"), "黑") {
		t.Fatalf("did not expect stale 黑 match")
	}
	if !contains(tree.MatchAll("红色大衣"), "红色大衣") {
		t.Fatalf("expected 红色大衣")
	}
}
