package str

import "testing"

func TestLevenshteinDistanceAndSimilarity(t *testing.T) {
	tests := []struct {
		name       string
		a          string
		b          string
		distance   int
		similarity float64
	}{
		{name: "both empty", a: "", b: "", distance: 0, similarity: 1},
		{name: "ascii edits", a: "kitten", b: "sitting", distance: 3, similarity: 1 - 3.0/7.0},
		{name: "unicode edits", a: "你好世界", b: "你好呀", distance: 2, similarity: 0.5},
		{name: "insertions", a: "go", b: "go语言", distance: 2, similarity: 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LevenshteinDistance(tt.a, tt.b); got != tt.distance {
				t.Fatalf("LevenshteinDistance(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.distance)
			}
			if got := LevenshteinSimilarity(tt.a, tt.b); got != tt.similarity {
				t.Fatalf("LevenshteinSimilarity(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.similarity)
			}
		})
	}
}
