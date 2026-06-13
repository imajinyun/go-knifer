package dfa

import "testing"

func TestFilter(t *testing.T) {
	Init([]string{"大", "大土豆", "土豆", "刚出锅", "出锅"})
	got := Filter(sampleText)
	want := "我有一颗$****，***的"
	if got != want {
		t.Fatalf("Filter() = %q, want %q", got, want)
	}
}

func TestFilterGreedyLongest(t *testing.T) {
	Init([]string{"赵", "赵阿", "赵阿三"})
	got := Filter("赵阿三在做什么。")
	want := "***在做什么。"
	if got != want {
		t.Fatalf("Filter() = %q, want %q", got, want)
	}
}

func TestCustomProcessor(t *testing.T) {
	tree := NewWordTree().AddWords("bad")
	got := tree.Filter("a bad word", true, func(word FoundWord) string {
		return "[" + word.Word + "]"
	})
	if got != "a [bad] word" {
		t.Fatalf("custom filter = %q", got)
	}
}
