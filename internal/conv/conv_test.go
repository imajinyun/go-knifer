package conv

import (
	"errors"
	"strconv"
	"testing"
)

// Tests cover the utility toolkit-core ConvertTest.

func TestToString(t *testing.T) {
	if ToString(nil) != "" {
		t.Fatalf("nil should be empty")
	}
	if ToString(123) != "123" {
		t.Fatalf("int")
	}
	if ToString(true) != "true" {
		t.Fatalf("bool")
	}
	if ToString(3.14) != "3.14" {
		t.Fatalf("float")
	}
	if ToString([]byte("hi")) != "hi" {
		t.Fatalf("bytes")
	}
	if ToStringDefault(nil, "x") != "x" {
		t.Fatalf("default")
	}
}

func TestToInt(t *testing.T) {
	if ToInt("123") != 123 {
		t.Fatalf("string int")
	}
	if ToInt("3.14") != 3 {
		t.Fatalf("string float")
	}
	if ToInt(int64(99)) != 99 {
		t.Fatalf("int64")
	}
	if ToInt(true) != 1 {
		t.Fatalf("bool true")
	}
	if ToIntDefault("abc", 42) != 42 {
		t.Fatalf("default")
	}
}

func TestToInt64AndFloat(t *testing.T) {
	if ToInt64("9999999999") != 9999999999 {
		t.Fatalf("ToInt64")
	}
	if ToFloat64("3.14") != 3.14 {
		t.Fatalf("ToFloat64")
	}
	if ToFloat64Default("x", 1.5) != 1.5 {
		t.Fatalf("ToFloat64 default")
	}
}

func TestToBool(t *testing.T) {
	cases := map[string]bool{
		"true": true, "yes": true, "y": true, "ok": true, "1": true, "on": true,
		"false": false, "no": false, "n": false, "0": false, "off": false,
	}
	for s, want := range cases {
		if ToBool(s) != want {
			t.Fatalf("ToBool(%q)", s)
		}
	}
	if ToBool(1) != true || ToBool(0) != false {
		t.Fatalf("ToBool int")
	}
	if ToBoolDefault("xx", true) != true {
		t.Fatalf("ToBool default")
	}
}

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

func TestToBytes(t *testing.T) {
	if string(ToBytes("ab")) != "ab" {
		t.Fatalf("ToBytes string")
	}
	if string(ToBytes([]byte("xy"))) != "xy" {
		t.Fatalf("ToBytes bytes")
	}
	if string(ToBytes(123)) != "123" {
		t.Fatalf("ToBytes int")
	}
	if ToBytes(nil) != nil {
		t.Fatalf("ToBytes nil")
	}
}
