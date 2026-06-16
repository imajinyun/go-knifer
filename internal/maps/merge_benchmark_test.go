package maps

import "testing"

func BenchmarkMerge_TwoMaps(b *testing.B) {
	a := makeBenchMap(1024)
	c := makeBenchMap(1024)
	b.ReportAllocs()
	for b.Loop() {
		benchMapResult = Merge(a, c)
	}
}

func BenchmarkMerge_FiveMaps(b *testing.B) {
	ms := []map[int]int{
		makeBenchMap(256),
		makeBenchMap(256),
		makeBenchMap(256),
		makeBenchMap(256),
		makeBenchMap(256),
	}
	b.ReportAllocs()
	for b.Loop() {
		benchMapResult = Merge(ms...)
	}
}

func BenchmarkMergeFunc_Sum(b *testing.B) {
	x := makeBenchMap(1024)
	y := makeBenchMap(1024)
	add := func(o, n int) int { return o + n }
	b.ReportAllocs()
	for b.Loop() {
		benchMapResult = MergeFunc(add, x, y)
	}
}
