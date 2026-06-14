package vnum

import (
	"math"
	"testing"
)

func TestGenericNumberFacade(t *testing.T) {
	if got := SumNumber[int](-2, 5, 7); got != 10 {
		t.Fatalf("SumNumber[int] = %v", got)
	}
	if got := SumNumber[float64](1.25, 2.5, -0.75); got != 3 {
		t.Fatalf("SumNumber[float64] = %v", got)
	}
	if got := AvgNumber[uint](2, 5, 8); got != 5 {
		t.Fatalf("AvgNumber[uint] = %v", got)
	}
	if got := MinInteger[int](-2, 5); got != -2 {
		t.Fatalf("MinInteger[int] = %d", got)
	}
	if got := MinIntegers[int](4, -8, 2); got != -8 {
		t.Fatalf("MinIntegers[int] = %d", got)
	}
	if got := MaxInteger[int](-2, 5); got != 5 {
		t.Fatalf("MaxInteger[int] = %d", got)
	}
	if got := MaxIntegers[int](4, -8, 2); got != 4 {
		t.Fatalf("MaxIntegers[int] = %d", got)
	}
	if got := MinFloat64s(3.5, -1.25, 2); got != -1.25 {
		t.Fatalf("MinFloat64s = %v", got)
	}
	if got := MaxFloat64s(3.5, -1.25, 2); got != 3.5 {
		t.Fatalf("MaxFloat64s = %v", got)
	}
	if got := AvgNumber[int](); got != 0 {
		t.Fatalf("AvgNumber empty = %v", got)
	}
	if got := MinIntegers[int](); got != 0 {
		t.Fatalf("MinIntegers empty = %d", got)
	}
	if got := MaxIntegers[int](); got != 0 {
		t.Fatalf("MaxIntegers empty = %d", got)
	}
	if got := AbsInteger[int](-12); got != 12 {
		t.Fatalf("AbsInteger[int] = %d", got)
	}
	if got := AbsInteger[int8](math.MinInt8); got != 0 {
		t.Fatalf("AbsInteger overflow = %d", got)
	}
	abs, err := AbsIntegerE[int8](math.MinInt8)
	if err == nil || abs != 0 {
		t.Fatalf("AbsIntegerE overflow = %d, %v", abs, err)
	}
	if got := AbsFloat32(-3.5); got != 3.5 {
		t.Fatalf("AbsFloat32 = %v", got)
	}
	if got := AbsFloat64(math.Inf(-1)); !math.IsInf(got, 1) {
		t.Fatalf("AbsFloat64(-Inf) = %v", got)
	}
}
