package conv

import (
	"errors"
	"strconv"
	"testing"
)

func TestConversionWithOptionsUsesProviders(t *testing.T) {
	if got := ToStringWithOptions(true, WithFormatBoolFunc(func(bool) string { return "YES" })); got != "YES" {
		t.Fatalf("ToStringWithOptions bool = %q", got)
	}
	if got := ToStringWithOptions(1.25, WithFormatFloatFunc(func(f float64, fmt byte, prec, bitSize int) string {
		if fmt != 'f' || prec != -1 || bitSize != 64 {
			t.Fatalf("float formatter args = %v %q %d %d", f, fmt, prec, bitSize)
		}
		return "float"
	})); got != "float" {
		t.Fatalf("ToStringWithOptions float = %q", got)
	}

	parseInt := func(text string, base, bitSize int) (int64, error) {
		if text == "custom-int" {
			return 42, nil
		}
		return strconv.ParseInt(text, base, bitSize)
	}
	parseFloat := func(text string, bitSize int) (float64, error) {
		if text == "custom-float" {
			return 6.5, nil
		}
		return strconv.ParseFloat(text, bitSize)
	}
	parseBool := func(text string) (bool, error) {
		if text == "enabled" {
			return true, nil
		}
		return false, errors.New("bad bool")
	}

	if got := ToIntWithOptions("custom-int", WithParseIntFunc(parseInt)); got != 42 {
		t.Fatalf("ToIntWithOptions = %d", got)
	}
	if got := ToInt64WithOptions("custom-float", WithParseFloatFunc(parseFloat)); got != 6 {
		t.Fatalf("ToInt64WithOptions = %d", got)
	}
	if got := ToFloat64WithOptions("custom-float", WithParseFloatFunc(parseFloat)); got != 6.5 {
		t.Fatalf("ToFloat64WithOptions = %v", got)
	}
	if got := ToBoolWithOptions("enabled", WithBoolParser(parseBool)); !got {
		t.Fatalf("ToBoolWithOptions = %v", got)
	}
	if got := ToBoolDefaultWithOptions("disabled", true, WithBoolParser(parseBool)); !got {
		t.Fatalf("ToBoolDefaultWithOptions should return default")
	}
	if got := string(ToBytesWithOptions(1.25, WithFormatFloatFunc(func(float64, byte, int, int) string { return "bytes" }))); got != "bytes" {
		t.Fatalf("ToBytesWithOptions = %q", got)
	}
}
