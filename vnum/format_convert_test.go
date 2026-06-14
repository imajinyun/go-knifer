package vnum

import (
	"strconv"
	"testing"
)

func TestNumFormatRoundAndConversionFacades(t *testing.T) {
	formatFloatCalls := 0
	formatIntCalls := 0
	formatFloat := func(f float64, fmtByte byte, prec int, bitSize int) string {
		formatFloatCalls++
		return strconv.FormatFloat(f, fmtByte, prec, bitSize)
	}
	formatInt := func(i int64, base int) string {
		formatIntCalls++
		return strconv.FormatInt(i, base)
	}
	formatOpts := []FormatOption{WithFormatFloatFunc(formatFloat), WithFormatIntFunc(formatInt)}

	if got := DecimalFormatWithOptions(",##0.00", 1234.5, formatOpts...); got != "1,234.50" {
		t.Fatalf("DecimalFormatWithOptions = %q", got)
	}
	if got := DecimalFormatMoney(1234.5); got != "1,234.50" {
		t.Fatalf("DecimalFormatMoney = %q", got)
	}
	if got := DecimalFormatMoneyWithOptions(12.3, formatOpts...); got != "12.30" {
		t.Fatalf("DecimalFormatMoneyWithOptions = %q", got)
	}
	if got := FormatPercent(0.125, 1); got != "12.5%" {
		t.Fatalf("FormatPercent = %q", got)
	}
	if got := FormatPercentWithOptions(0.5, 0, formatOpts...); got != "50%" {
		t.Fatalf("FormatPercentWithOptions = %q", got)
	}
	if got := ToStr(12.3400); got != "12.34" {
		t.Fatalf("ToStr = %q", got)
	}
	if got := ToStrWithOptions(10, formatOpts...); got != "10" {
		t.Fatalf("ToStrWithOptions = %q", got)
	}
	if got := ToStrDefault(nil, "empty"); got != "empty" {
		t.Fatalf("ToStrDefault nil = %q", got)
	}
	value := 7.5
	if got := ToStrDefaultWithOptions(&value, "empty", formatOpts...); got != "7.5" {
		t.Fatalf("ToStrDefaultWithOptions = %q", got)
	}
	if got := ToStrStripWithOptions(7.500, false, formatOpts...); got != "7.5" {
		t.Fatalf("ToStrStripWithOptions = %q", got)
	}
	if got := GetBinaryStrWithOptions(10, formatOpts...); got != "1010" {
		t.Fatalf("GetBinaryStrWithOptions = %q", got)
	}
	if _, err := BinaryToIntWithOptions("1010", WithParseIntFunc(strconv.ParseInt)); err != nil {
		t.Fatalf("BinaryToIntWithOptions: %v", err)
	}
	if got, err := BinaryToLongWithOptions("1010", WithParseIntFunc(strconv.ParseInt)); err != nil || got != 10 {
		t.Fatalf("BinaryToLongWithOptions = %d, %v", got, err)
	}
	if formatFloatCalls == 0 || formatIntCalls == 0 {
		t.Fatalf("custom formatters not called: float=%d int=%d", formatFloatCalls, formatIntCalls)
	}

	if got := RoundMode(1.25, 1, RoundHalfEven); got != 1.2 {
		t.Fatalf("RoundMode half even = %f", got)
	}
	if got := RoundStrWithOptions(1.234, 2, formatOpts...); got != "1.23" {
		t.Fatalf("RoundStrWithOptions = %q", got)
	}
	if got := RoundHalfEvenFloat(1.25, 1); got != 1.2 {
		t.Fatalf("RoundHalfEvenFloat = %f", got)
	}
	if got := RoundDownFloat(1.29, 1); got != 1.2 {
		t.Fatalf("RoundDownFloat = %f", got)
	}

	doubleCalls := 0
	if got := ToDoubleWithOptions(float32(1.25),
		WithDoubleFormatFloatFunc(strconv.FormatFloat),
		WithDoubleParseFloatFunc(func(s string, bitSize int) (float64, error) {
			doubleCalls++
			return strconv.ParseFloat(s, bitSize)
		}),
	); got != 1.25 || doubleCalls == 0 {
		t.Fatalf("ToDoubleWithOptions = %f calls=%d", got, doubleCalls)
	}
}
