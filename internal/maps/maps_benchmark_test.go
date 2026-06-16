package maps

import "testing"

var benchMapResult any

func makeBenchMap(n int) map[int]int {
	m := make(map[int]int, n)
	for i := 0; i < n; i++ {
		m[i] = i
	}
	return m
}

func BenchmarkFilter(b *testing.B) {
	m := makeBenchMap(4096)
	b.ReportAllocs()
	for b.Loop() {
		benchMapResult = Filter(m, func(_ int, v int) bool { return v%2 == 0 })
	}
}

func BenchmarkSortedKeys(b *testing.B) {
	m := makeBenchMap(4096)
	b.ReportAllocs()
	for b.Loop() {
		benchMapResult = SortedKeys(m)
	}
}

func BenchmarkClone(b *testing.B) {
	m := makeBenchMap(4096)
	b.ReportAllocs()
	for b.Loop() {
		benchMapResult = Clone(m)
	}
}
