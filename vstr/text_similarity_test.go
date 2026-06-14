package vstr

import "testing"

func TestFacadeLevenshteinSimilarity(t *testing.T) {
	if got := LevenshteinDistance("kitten", "sitting"); got != 3 {
		t.Fatalf("LevenshteinDistance = %d, want 3", got)
	}
	if got := LevenshteinSimilarity("你好世界", "你好呀"); got != 0.5 {
		t.Fatalf("LevenshteinSimilarity = %v, want 0.5", got)
	}
}
