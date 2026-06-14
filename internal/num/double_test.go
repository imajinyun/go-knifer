package num

import (
	"strconv"
	"testing"
)

func TestToDoubleHelpers(t *testing.T) {
	if ToDouble(float32(1.23)) != 1.23 {
		t.Fatal("ToDouble should preserve float32 textual precision")
	}

	toDoubleCases := []struct {
		input any
		want  float64
	}{
		{float32(1.25), 1.25},
		{float64(2.5), 2.5},
		{int(-3), -3},
		{int64(4), 4},
		{uint64(5), 5},
		{"bad", 0},
	}
	for _, tt := range toDoubleCases {
		if got := ToDouble(tt.input); got != tt.want {
			t.Fatalf("ToDouble(%T) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestToDoubleWithOptionsUsesProviders(t *testing.T) {
	formatCalled := false
	parseCalled := false
	if got := ToDoubleWithOptions(float32(1.25),
		WithDoubleFormatFloatFunc(func(v float64, fmtByte byte, prec, bitSize int) string {
			formatCalled = true
			return strconv.FormatFloat(v*2, fmtByte, prec, bitSize)
		}),
		WithDoubleParseFloatFunc(func(s string, bitSize int) (float64, error) {
			parseCalled = true
			return strconv.ParseFloat(s, bitSize)
		}),
	); got != 2.5 || !formatCalled || !parseCalled {
		t.Fatalf("ToDoubleWithOptions = %v format=%v parse=%v", got, formatCalled, parseCalled)
	}
}
