package num

import "testing"

func TestNumberChecksAndFormat(t *testing.T) {
	if !IsNumber("0x1A") || !IsNumber("123E3") || !IsNumber("123D") || !IsNumber("+123") {
		t.Fatal("IsNumber should support hex, scientific, type suffix and signs")
	}
	if IsNumber("0x") || IsNumber("1E-") || !IsLong("9223372036854775807") || !IsDouble("1.23") || IsDouble("123") {
		t.Fatal("number checks failed")
	}
	if !IsPrimes(97) || IsPrimes(100) {
		t.Fatal("prime check failed")
	}
	if DecimalFormatMoney(12345.6) != "12,345.60" || FormatPercent(0.1234, 2) != "12.34%" {
		t.Fatalf("format failed: %s %s", DecimalFormatMoney(12345.6), FormatPercent(0.1234, 2))
	}
	if CeilDiv(10, 3) != 4 || RoundStr(1.2, 2) != "1.20" || RoundHalfEvenFloat(2.5, 0) != 2 || RoundDownFloat(1.29, 1) != 1.2 {
		t.Fatal("round helpers failed")
	}
}

func TestFormatWithOptionsUsesProviders(t *testing.T) {
	floatCalls := 0
	floatFormatter := func(v float64, fmt byte, prec, bitSize int) string {
		floatCalls++
		if fmt == 'f' && prec == 2 && bitSize == 64 {
			return "custom-float"
		}
		return "fallback-float"
	}
	if got := RoundStrWithOptions(1.2, 2, WithFormatFloatFunc(floatFormatter)); got != "custom-float" {
		t.Fatalf("RoundStrWithOptions = %q", got)
	}
	if got := DecimalFormatWithOptions("0.00", 1.2, WithFormatFloatFunc(floatFormatter)); got != "custom-float" {
		t.Fatalf("DecimalFormatWithOptions = %q", got)
	}
	if got := ToStrStripWithOptions(1.2, false, WithFormatFloatFunc(func(v float64, fmt byte, prec, bitSize int) string {
		if fmt != 'f' || prec != -1 || bitSize != 64 {
			t.Fatalf("format args fmt=%q prec=%d bitSize=%d", fmt, prec, bitSize)
		}
		return "1.200"
	})); got != "1.200" {
		t.Fatalf("ToStrStripWithOptions = %q", got)
	}
	intCalls := 0
	if got := GetBinaryStrWithOptions(int64(5), WithFormatIntFunc(func(v int64, base int) string {
		intCalls++
		if v != 5 || base != 2 {
			t.Fatalf("format int args v=%d base=%d", v, base)
		}
		return "custom-int"
	})); got != "custom-int" || intCalls != 1 {
		t.Fatalf("GetBinaryStrWithOptions = %q intCalls=%d", got, intCalls)
	}
	if floatCalls < 2 {
		t.Fatalf("float formatter calls = %d", floatCalls)
	}
}
