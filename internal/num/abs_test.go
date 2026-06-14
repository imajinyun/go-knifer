package num

import (
	"math"
	"testing"
)

func TestGenericNumberAbs(t *testing.T) {
	if got := AbsInteger[int](-12); got != 12 {
		t.Fatalf("AbsInteger[int] = %d", got)
	}
	if got := AbsInteger[uint](12); got != 12 {
		t.Fatalf("AbsInteger[uint] = %d", got)
	}
	if got := AbsInteger[int8](math.MinInt8); got != 0 {
		t.Fatalf("AbsInteger overflow = %d", got)
	}
	abs, err := AbsIntegerE[int8](math.MinInt8)
	if err == nil || abs != 0 {
		t.Fatalf("AbsIntegerE overflow = %d, %v", abs, err)
	}
	if got := AbsFloat32(float32(math.Copysign(0, -1))); math.Signbit(float64(got)) || got != 0 {
		t.Fatalf("AbsFloat32(-0) = %v", got)
	}
	if got := AbsFloat32(-3.5); got != 3.5 {
		t.Fatalf("AbsFloat32 = %v", got)
	}
	if got := AbsFloat64(-4.25); got != 4.25 {
		t.Fatalf("AbsFloat64 = %v", got)
	}
	if got := AbsFloat64(math.Inf(-1)); !math.IsInf(got, 1) {
		t.Fatalf("AbsFloat64(-Inf) = %v", got)
	}
}
