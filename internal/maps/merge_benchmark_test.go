package maps

import "testing"

func BenchmarkMerge_TwoMaps(b *testing.B) {
	a := makeBenchMap(1024)
	c := makeBenchMap(1024)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Merge(a, c)
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Merge(ms...)
	}
}

func BenchmarkMergeFunc_Sum(b *testing.B) {
	x := makeBenchMap(1024)
	y := makeBenchMap(1024)
	add := func(o, n int) int { return o + n }
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MergeFunc(add, x, y)
	}
}
