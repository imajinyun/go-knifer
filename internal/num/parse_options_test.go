package num

import (
	"errors"
	"testing"
)

func TestParseWithOptionsUsesProviders(t *testing.T) {
	var intCalled, floatCalled int
	parseInt := func(text string, base, bitSize int) (int64, error) {
		intCalled++
		switch text {
		case "custom-int":
			return 42, nil
		case "1010":
			if base != 2 {
				t.Fatalf("binary base = %d", base)
			}
			return 10, nil
		case "ff":
			if base != 16 {
				t.Fatalf("hex base = %d", base)
			}
			return 255, nil
		default:
			return 0, errors.New("custom int error")
		}
	}
	parseFloat := func(text string, bitSize int) (float64, error) {
		floatCalled++
		if text == "custom-float" {
			return 6.5, nil
		}
		return 0, errors.New("custom float error")
	}

	if got := ParseIntWithOptions("custom-int", WithParseIntFunc(parseInt)); got != 42 {
		t.Fatalf("ParseIntWithOptions = %d", got)
	}
	if got := ParseLongWithOptions("0xff", WithParseIntFunc(parseInt)); got != 255 {
		t.Fatalf("ParseLongWithOptions = %d", got)
	}
	if got := ParseDoubleWithOptions("custom-float", WithParseFloatFunc(parseFloat)); got != 6.5 {
		t.Fatalf("ParseDoubleWithOptions = %v", got)
	}
	if got, err := ParseNumberWithOptions("custom-float", WithParseFloatFunc(parseFloat)); err != nil || got != 6.5 {
		t.Fatalf("ParseNumberWithOptions = %v, %v", got, err)
	}
	if got, err := BinaryToIntWithOptions("1010", WithParseIntFunc(parseInt)); err != nil || got != 10 {
		t.Fatalf("BinaryToIntWithOptions = %d, %v", got, err)
	}
	if !IsNumberWithOptions("custom-float", WithParseFloatFunc(parseFloat)) || !IsIntegerWithOptions("custom-int", WithParseIntFunc(parseInt)) {
		t.Fatal("Is*WithOptions should use custom parsers")
	}
	if got := ParseFloatDefaultWithOptions("custom-float", 1, WithParseFloatFunc(parseFloat)); got != 6.5 {
		t.Fatalf("ParseFloatDefaultWithOptions = %v", got)
	}
	if got := ToBigDecimalWithOptions("not-rat", WithParseFloatFunc(parseFloat)); got.Sign() != 0 {
		t.Fatalf("ToBigDecimalWithOptions fallback = %s", got.String())
	}
	if intCalled == 0 || floatCalled == 0 {
		t.Fatalf("providers not called int=%d float=%d", intCalled, floatCalled)
	}
}

func TestBigIntegerPartPowAndParseEdges(t *testing.T) {
	bigIntCases := map[string]int64{"42": 42, "-0x10": -16, "#10": 16, "010": 8}
	for input, want := range bigIntCases {
		got, ok := NewBigInteger(input)
		if !ok || got.Int64() != want {
			t.Fatalf("NewBigInteger(%q) = %v/%v, want %d", input, got, ok, want)
		}
	}
	if got, ok := NewBigInteger(""); ok || got != nil {
		t.Fatalf("NewBigInteger blank = %v/%v", got, ok)
	}
	if got, ok := NewBigInteger("0xzz"); ok || got != nil {
		t.Fatalf("NewBigInteger invalid = %v/%v", got, ok)
	}
	if !IsBeside(2, 1) || !IsBeside[int64](1, 2) || IsBeside(1, 3) {
		t.Fatal("IsBeside cases failed")
	}
	if PartValueWithMode(10, 0, true) != 0 || PartValueWithMode(10, 3, false) != 3 || PartValueWithMode(10, 3, true) != 4 {
		t.Fatal("PartValueWithMode cases failed")
	}
	if Pow(2, 3) != 8 || Pow(2, -3) != 0.13 || PowWithMode(2, -3, 2, RoundDown) != 0.12 {
		t.Fatal("Pow cases failed")
	}
	if IsPowerOfTwo(0) || IsPowerOfTwo(-2) || !IsPowerOfTwo(1) || !IsPowerOfTwo(1024) || IsPowerOfTwo(1023) {
		t.Fatal("IsPowerOfTwo cases failed")
	}
	if ParseInt("") != 0 || ParseInt(".5") != 0 || ParseInt("1e3") != 0 || ParseInt("1,234.9") != 1234 {
		t.Fatal("ParseInt edge cases failed")
	}
	if ParseLong("") != 0 || ParseLong(".5") != 0 || ParseLong("0x7f") != 127 || ParseLong("1,234.9") != 1234 {
		t.Fatal("ParseLong edge cases failed")
	}
	if ParseDouble("") != 0 || ParseDouble("1,234.5") != 1234.5 || ParseFloat("2.5") != 2.5 {
		t.Fatal("ParseFloat/ParseDouble edge cases failed")
	}
	if got, err := ParseNumber("+1,234.5"); err != nil || got != 1234.5 {
		t.Fatalf("ParseNumber plus/comma failed: %v %v", got, err)
	}
	if got, err := ParseNumber("0x10"); err != nil || got != 16 {
		t.Fatalf("ParseNumber hex failed: %v %v", got, err)
	}
	if _, err := ParseNumber("bad"); err == nil {
		t.Fatal("ParseNumber should reject invalid input")
	}
	if ParseIntDefault("", 7) != 7 || ParseIntDefault("bad", 7) != 7 || ParseIntDefault("1,234", 7) != 1234 {
		t.Fatal("ParseIntDefault cases failed")
	}
	if ParseLongDefault("", 8) != 8 || ParseLongDefault("bad", 8) != 8 || ParseLongDefault("1,234", 8) != 1234 {
		t.Fatal("ParseLongDefault cases failed")
	}
	if ParseFloatDefault("", 1.5) != 1.5 || ParseFloatDefault("bad", 1.5) != 1.5 || ParseFloatDefault("1,234.5", 1.5) != 1234.5 {
		t.Fatal("ParseFloatDefault cases failed")
	}
	if ParseDoubleDefault("", 2.5) != 2.5 || ParseDoubleDefault("bad", 2.5) != 2.5 || ParseDoubleDefault("1,234.5", 2.5) != 1234.5 {
		t.Fatal("ParseDoubleDefault cases failed")
	}
}
