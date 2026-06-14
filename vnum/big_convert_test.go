package vnum

import "testing"

func TestNumBigAndNullConversionFacades(t *testing.T) {
	if got := ToBigDecimal("1,234.50"); got.String() != "2469/2" {
		t.Fatalf("ToBigDecimal = %s", got)
	}
	if got := ToBigDecimalWithOptions("bad", WithParseFloatFunc(func(string, int) (float64, error) { return 2.5, nil })); got.String() != "5/2" {
		t.Fatalf("ToBigDecimalWithOptions fallback = %s", got)
	}
	if got := ToBigInteger("42"); got.String() != "42" {
		t.Fatalf("ToBigInteger = %s", got)
	}
	if got := Null2Zero(nil); got.Sign() != 0 {
		t.Fatalf("Null2Zero = %s", got)
	}
	if Zero2One(0) != 1 || Count(10, 3) != 4 {
		t.Fatal("Zero2One/Count failed")
	}
	var n int
	if NullToZero(&n) != 0 || NullToZero[int](nil) != 0 {
		t.Fatal("NullToZero failed")
	}
	if NullBigIntToZero(nil).Sign() != 0 || NullBigDecimalToZero(nil).Sign() != 0 {
		t.Fatal("null big conversions failed")
	}
	if got, ok := NewBigInteger("0x10"); !ok || got.Int64() != 16 {
		t.Fatalf("NewBigInteger = %v, %v", got, ok)
	}
	if got := PartValue(10, 3); got != 4 {
		t.Fatalf("PartValue = %d", got)
	}
	if got := PartValueWithMode(10, 3, false); got != 3 {
		t.Fatalf("PartValueWithMode = %d", got)
	}
}
