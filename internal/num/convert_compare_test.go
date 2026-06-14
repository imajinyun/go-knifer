package num

import (
	"math"
	"math/big"
	"testing"
)

func TestBinaryCompareAndConversion(t *testing.T) {
	if GetBinaryStr(int8(-1)) != "11111111" || GetBinaryStr(float32(1)) != "00111111100000000000000000000000" {
		t.Fatal("GetBinaryStr failed")
	}
	if got, err := BinaryToInt("1010"); err != nil || got != 10 {
		t.Fatalf("BinaryToInt: %d %v", got, err)
	}
	if got, err := BinaryToLong("1010"); err != nil || got != 10 {
		t.Fatalf("BinaryToLong: %d %v", got, err)
	}
	if Compare(1, 2) >= 0 || !IsGreater(3, 2) || !IsIn(2, 1, 3) || !EqualsExact(0.0, 0.0) || !EqualsChar('A', 'a', true) {
		t.Fatal("compare helpers failed")
	}
	if ToStr(5.0) != "5" || ToBigDecimal("1,234.50").FloatString(2) != "1234.50" || ToBigInteger("123").String() != "123" {
		t.Fatal("to string/big helpers failed")
	}
	if Count(10, 3) != 4 || Zero2One(0) != 1 || NullToZero[int](nil) != 0 || !IsBeside(1, 2) || PartValue(10, 3) != 4 {
		t.Fatal("small helpers failed")
	}
	bi, ok := NewBigInteger("0x10")
	if !ok || bi.Int64() != 16 {
		t.Fatal("NewBigInteger failed")
	}
}

func TestStringBigAndNullConversionEdges(t *testing.T) {
	if ToStrStrip(math.NaN(), true) != "" || ToStrStrip(12.3400, false) != "12.34" || ToStrStrip(12.0, true) != "12" {
		t.Fatal("ToStrStrip edge cases failed")
	}
	v := 12.0
	if ToStrDefault(nil, "fallback") != "fallback" || ToStrDefault(&v, "fallback") != "12" {
		t.Fatal("ToStrDefault edge cases failed")
	}
	if ToBigDecimal("").Sign() != 0 || ToBigDecimal("bad").Sign() != 0 || ToBigDecimal("1,234.25").FloatString(2) != "1234.25" {
		t.Fatal("ToBigDecimal edge cases failed")
	}
	if ToBigInteger("").Sign() != 0 || ToBigInteger("123.9").Int64() != 123 {
		t.Fatal("ToBigInteger edge cases failed")
	}
	if Count(0, 3) != 0 || Count(10, 0) != 0 || Count(9, 3) != 3 || Count(10, 3) != 4 {
		t.Fatal("Count edge cases failed")
	}
	if Null2Zero(nil).Sign() != 0 || NullBigIntToZero(nil).Sign() != 0 || NullBigDecimalToZero(nil).Sign() != 0 {
		t.Fatal("null to zero nil cases failed")
	}
	bi := big.NewInt(5)
	bd := big.NewRat(5, 2)
	if NullBigIntToZero(bi) != bi || NullBigDecimalToZero(bd) != bd || Null2Zero(bd) != bd {
		t.Fatal("null to zero should preserve non-nil pointers")
	}
	x := 7
	if NullToZero(&x) != 7 || Zero2One(9) != 9 {
		t.Fatal("small conversion helpers failed")
	}
}
