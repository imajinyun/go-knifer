package dfa

import "testing"

const sampleText = "我有一颗$大土^豆，刚出锅的"

func buildTestTree() *WordTree {
	return NewWordTree().AddWords("大", "大土豆", "土豆", "刚出锅", "出锅")
}

func assertFoundWord(t *testing.T, got FoundWord, word, found string, start, end int) {
	t.Helper()
	if got.Word != word || got.FoundWord != found || got.Start != start || got.End != end {
		t.Fatalf("FoundWord = %#v, want word=%q found=%q start=%d end=%d", got, word, found, start, end)
	}
}

func contains(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}
