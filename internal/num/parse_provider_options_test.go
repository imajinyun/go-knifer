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
