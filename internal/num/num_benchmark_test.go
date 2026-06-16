package num

import "testing"

var benchNumResult any

func makeBenchFloat64s(n int) []float64 {
	out := make([]float64, n)
	for i := range out {
		out[i] = float64(i%97) + 0.125
	}
	return out
}

func numBenchSizes() []struct {
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

func BenchmarkAdd(b *testing.B) {
	for _, tt := range numBenchSizes() {
		values := makeBenchFloat64s(tt.size)
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchNumResult = Add(values...)
			}
		})
	}
}

func BenchmarkDecimalFormat(b *testing.B) {
	for _, tt := range []struct {
		name   string
		format string
		value  float64
	}{
		{name: "empty", format: "", value: 0},
		{name: "small", format: "#.##", value: 123.456},
		{name: "medium", format: "#,###.0000", value: 123456.789},
		{name: "large", format: "#,###,###,###.000000", value: 123456789.123456},
	} {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchNumResult = DecimalFormat(tt.format, tt.value)
			}
		})
	}
}

func BenchmarkFormatPercent(b *testing.B) {
	for _, tt := range []struct {
		name  string
		scale int
		value float64
	}{
		{name: "empty", scale: 0, value: 0},
		{name: "small", scale: 2, value: 0.1234},
		{name: "medium", scale: 4, value: 0.123456},
		{name: "large", scale: 8, value: 0.123456789},
	} {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchNumResult = FormatPercent(tt.value, tt.scale)
			}
		})
	}
}

func BenchmarkRoundStr(b *testing.B) {
	for _, tt := range []struct {
		name  string
		scale int
		value float64
	}{
		{name: "empty", scale: 0, value: 0},
		{name: "small", scale: 2, value: 123.456},
		{name: "medium", scale: 6, value: 123456.789123},
		{name: "large", scale: 12, value: 123456789.123456789},
	} {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchNumResult = RoundStr(tt.value, tt.scale)
			}
		})
	}
}
