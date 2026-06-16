package str

import (
	"strings"
	"testing"
)

var benchStringResult any

func stringBenchInputs() []struct {
	name  string
	input string
} {
	return []struct {
		name  string
		input string
	}{
		{name: "empty", input: ""},
		{name: "small", input: "go knifer"},
		{name: "medium", input: strings.Repeat("go knifer ", 128)},
		{name: "large", input: strings.Repeat("go knifer ", 2048)},
	}
}

func BenchmarkReverse(b *testing.B) {
	for _, tt := range stringBenchInputs() {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchStringResult = Reverse(tt.input)
			}
		})
	}
}

func BenchmarkFormat(b *testing.B) {
	for _, tt := range stringBenchInputs() {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchStringResult = Format("prefix={} payload={} suffix={}", tt.name, tt.input, len(tt.input))
			}
		})
	}
}

func BenchmarkContainsIgnoreCase(b *testing.B) {
	for _, tt := range stringBenchInputs() {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchStringResult = ContainsIgnoreCase(tt.input, "KNIFER")
			}
		})
	}
}

func BenchmarkLevenshteinDistance(b *testing.B) {
	for _, tt := range []struct {
		name  string
		left  string
		right string
	}{
		{name: "empty", left: "", right: ""},
		{name: "small", left: "kitten", right: "sitting"},
		{name: "medium", left: strings.Repeat("abcdef", 16), right: strings.Repeat("abcxef", 16)},
		{name: "large", left: strings.Repeat("abcdef", 64), right: strings.Repeat("abcxef", 64)},
	} {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchStringResult = LevenshteinDistance(tt.left, tt.right)
			}
		})
	}
}
