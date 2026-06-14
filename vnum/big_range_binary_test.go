package vnum

import (
	"math/big"
	"reflect"
	"testing"
)

func TestNumBigRangeFactorialAndBinaryFacades(t *testing.T) {
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

	if got, err := Factorial(5); err != nil || got != 120 {
		t.Fatalf("Factorial = %d, %v", got, err)
	}
	if _, err := Factorial(21); err == nil {
		t.Fatal("Factorial overflow error = nil")
	}
	if got, err := FactorialRange(5, 2); err != nil || got != 60 {
		t.Fatalf("FactorialRange = %d, %v", got, err)
	}
	if got := FactorialBig(big.NewInt(5)); got.String() != "120" {
		t.Fatalf("FactorialBig = %s", got)
	}
	if got := FactorialBigRange(big.NewInt(5), big.NewInt(2)); got.String() != "60" {
		t.Fatalf("FactorialBigRange = %s", got)
	}
	if Sqrt(16) != 4 || ProcessMultiple(5, 2) != 10 || Divisor(24, 18) != 6 || Multiple(4, 6) != 12 {
		t.Fatal("factorial helpers failed")
	}
	if got := RangeClosed(1, 3, 0); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("RangeClosed = %v", got)
	}
	if got := AppendRange(3, 1, 1, []int{0}); !reflect.DeepEqual(got, []int{0, 3, 2, 1}) {
		t.Fatalf("AppendRange = %v", got)
	}
	if got := GenerateRandomNumberWithSeed(0, 4, 2, []int{1, 2, 3}); len(got) != 2 {
		t.Fatalf("GenerateRandomNumberWithSeed = %v", got)
	}
	if got := GenerateRandomNumber(0, 3, 2); len(got) != 2 {
		t.Fatalf("GenerateRandomNumber = %v", got)
	}
	if got := GenerateBySet(0, 3, 2); len(got) != 2 {
		t.Fatalf("GenerateBySet = %v", got)
	}

	if !EqualsExact(1, 1) || !EqualsFloat32Exact(1, 1) || !EqualsInt64(1, 1) || !EqualsBigDecimal(big.NewRat(1, 2), big.NewRat(2, 4)) || !EqualsChar('A', 'a', true) {
		t.Fatal("equals helpers failed")
	}
	if Compare(1, 2) != -1 || !IsGreater(2, 1) || !IsGreaterOrEqual(2, 2) || !IsLess(1, 2) || !IsLessOrEqual(2, 2) || !IsIn(2, 1, 3) {
		t.Fatal("comparison helpers failed")
	}
	if got, err := ToUnsignedByteArrayLen(2, big.NewInt(255)); err != nil || !reflect.DeepEqual(got, []byte{0, 255}) {
		t.Fatalf("ToUnsignedByteArrayLen = %v, %v", got, err)
	}
	if _, err := ToUnsignedByteArrayLen(1, big.NewInt(256)); err == nil {
		t.Fatal("ToUnsignedByteArrayLen overflow error = nil")
	}
	if got := FromUnsignedByteArray([]byte{1, 0}); got.Int64() != 256 {
		t.Fatalf("FromUnsignedByteArray = %s", got)
	}
	if got := FromUnsignedByteArrayRange([]byte{0, 1, 0}, 1, 2); got.Int64() != 256 {
		t.Fatalf("FromUnsignedByteArrayRange = %s", got)
	}
	if ToInt(ToBytes(12345)) != 12345 {
		t.Fatal("ToBytes/ToInt round trip failed")
	}
}
