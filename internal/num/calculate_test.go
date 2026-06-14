package num

import (
	"math"
	"strconv"
	"testing"
)

func TestCalculateExpressions(t *testing.T) {
	result, err := Calculate("(0*1--3)-5/-4-(3*(-2.13))")
	if err != nil || math.Abs(result-10.64) > 1e-9 {
		t.Fatalf("Calculate: %v %v", result, err)
	}

	calcCases := map[string]float64{
		"1 + 2 * 3":   7,
		"(1 + 2) * 3": 9,
		"10 % 4":      2,
		"--2 + +3":    5,
		" 3.5 / 2 ":   1.75,
	}
	for expr, want := range calcCases {
		got, err := Calculate(expr)
		if err != nil || math.Abs(got-want) > 1e-9 {
			t.Fatalf("Calculate(%q) = %v, %v, want %v", expr, got, err, want)
		}
	}
}

func TestCalculateWithOptionsUsesProvider(t *testing.T) {
	calcParseCalled := false
	got, err := CalculateWithOptions("5 + 2", WithParseFloatFunc(func(s string, bitSize int) (float64, error) {
		calcParseCalled = true
		if s == "5" {
			return 5, nil
		}
		return strconv.ParseFloat(s, bitSize)
	}))
	if err != nil || got != 7 || !calcParseCalled {
		t.Fatalf("CalculateWithOptions = %v, %v called=%v", got, err, calcParseCalled)
	}
}

func TestCalculateInvalidExpressions(t *testing.T) {
	invalidExpressions := []string{"", "1+", "(1+2", "1 2", "abc"}
	for _, expr := range invalidExpressions {
		if got, err := Calculate(expr); err == nil {
			t.Fatalf("Calculate(%q) should fail, got %v", expr, got)
		}
	}
}
