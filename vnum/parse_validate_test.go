package vnum

import (
	"math"
	"strconv"
	"testing"
)

func TestNumParseFormatAndValidateFacades(t *testing.T) {
	parseIntCalls := 0
	parseFloatCalls := 0
	parseInt := func(s string, base int, bitSize int) (int64, error) {
		parseIntCalls++
		return strconv.ParseInt(s, base, bitSize)
	}
	parseFloat := func(s string, bitSize int) (float64, error) {
		parseFloatCalls++
		return strconv.ParseFloat(s, bitSize)
	}
	opts := []ParseOption{WithParseIntFunc(parseInt), WithParseFloatFunc(parseFloat)}

	if got := ParseInt("0x10"); got != 16 {
		t.Fatalf("ParseInt = %d", got)
	}
	if got := ParseIntWithOptions("42", opts...); got != 42 {
		t.Fatalf("ParseIntWithOptions = %d", got)
	}
	if got := ParseLongWithOptions("1,234", opts...); got != 1234 {
		t.Fatalf("ParseLongWithOptions = %d", got)
	}
	if got := ParseFloatWithOptions("3.5", opts...); got != 3.5 {
		t.Fatalf("ParseFloatWithOptions = %f", got)
	}
	if got := ParseDoubleWithOptions("6.25", opts...); got != 6.25 {
		t.Fatalf("ParseDoubleWithOptions = %f", got)
	}
	if got, err := ParseNumberWithOptions("0x2a", opts...); err != nil || got != 42 {
		t.Fatalf("ParseNumberWithOptions = %f, %v", got, err)
	}
	if parseIntCalls == 0 || parseFloatCalls == 0 {
		t.Fatalf("custom parsers not called: int=%d float=%d", parseIntCalls, parseFloatCalls)
	}

	if got := ParseIntDefaultWithOptions("bad", 7, opts...); got != 7 {
		t.Fatalf("ParseIntDefaultWithOptions = %d", got)
	}
	if got := ParseLongDefaultWithOptions("", 9, opts...); got != 9 {
		t.Fatalf("ParseLongDefaultWithOptions = %d", got)
	}
	if got := ParseFloatDefaultWithOptions("bad", 1.5, opts...); got != 1.5 {
		t.Fatalf("ParseFloatDefaultWithOptions = %f", got)
	}
	if got := ParseDoubleDefaultWithOptions("bad", 2.5, opts...); got != 2.5 {
		t.Fatalf("ParseDoubleDefaultWithOptions = %f", got)
	}
	if !IsNumberWithOptions("0x10", opts...) || !IsIntegerWithOptions("42", opts...) || !IsLongWithOptions("42", opts...) || !IsDoubleWithOptions("3.14", opts...) {
		t.Fatal("numeric validation with options failed")
	}
	if !IsValidNumber(1.25) || IsValid(math.Inf(1)) || IsValidFloat32(float32(math.NaN())) || !IsOdd(3) || !IsEven(4) || !IsPowerOfTwo(64) {
		t.Fatal("validation helpers failed")
	}
}
