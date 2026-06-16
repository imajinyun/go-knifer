package slice

import (
	"fmt"
	"testing"
)

var benchSliceResult any

func makeBenchSlice(n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = i % max(n/2, 1)
	}
	return out
}

func sliceBenchSizes() []struct {
	name string
	size int
} {
	return []struct {
		name string
		size int
	}{
		{name: "empty", size: 0},
		{name: "small", size: 16},
		{name: "medium", size: 1024},
		{name: "large", size: 16384},
	}
}

func BenchmarkMap(b *testing.B) {
	for _, tt := range sliceBenchSizes() {
		input := makeBenchSlice(tt.size)
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSliceResult = Map(input, func(v int) int { return v * 2 })
			}
		})
	}
}

func BenchmarkFilter(b *testing.B) {
	for _, tt := range sliceBenchSizes() {
		input := makeBenchSlice(tt.size)
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSliceResult = Filter(input, func(v int) bool { return v%2 == 0 })
			}
		})
	}
}

func BenchmarkDistinct(b *testing.B) {
	for _, tt := range sliceBenchSizes() {
		input := makeBenchSlice(tt.size)
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSliceResult = Distinct(input)
			}
		})
	}
}

func BenchmarkGroupBy(b *testing.B) {
	for _, tt := range sliceBenchSizes() {
		input := makeBenchSlice(tt.size)
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSliceResult = GroupBy(input, func(v int) string { return fmt.Sprintf("bucket-%d", v%8) })
			}
		})
	}
}

func BenchmarkSetOperations(b *testing.B) {
	for _, tt := range sliceBenchSizes() {
		left := makeBenchSlice(tt.size)
		right := makeBenchSlice(tt.size / 2)
		b.Run(tt.name+"/union", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSliceResult = Union(left, right)
			}
		})
		b.Run(tt.name+"/intersection", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSliceResult = Intersection(left, right)
			}
		})
		b.Run(tt.name+"/subtract", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSliceResult = Subtract(left, right)
			}
		})
	}
}
