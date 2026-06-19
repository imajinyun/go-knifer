package conv

import (
	"errors"
	"strconv"
	"testing"
)

// TestToInt64AllNumericKinds exercises every concrete numeric branch of toInt64.
func TestToInt64AllNumericKinds(t *testing.T) {
	cases := []struct {
		name string
		in   any
		want int64
	}{
		{"int", int(1), 1},
		{"int8", int8(2), 2},
		{"int16", int16(3), 3},
		{"int32", int32(4), 4},
		{"int64", int64(5), 5},
		{"uint", uint(6), 6},
		{"uint8", uint8(7), 7},
		{"uint16", uint16(8), 8},
		{"uint32", uint32(9), 9},
		{"uint64", uint64(10), 10},
		{"float32", float32(11.9), 11},
		{"float64", float64(12.9), 12},
		{"bool-true", true, 1},
		{"bool-false", false, 0},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToInt64(tt.in); got != tt.want {
				t.Fatalf("ToInt64(%v) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

// TestToInt64ReflectKinds covers the reflect fallback for named numeric types.
func TestToInt64ReflectKinds(t *testing.T) {
	type myInt int
	type myUint uint
	type myFloat float64
	if got := ToInt64(myInt(5)); got != 5 {
		t.Fatalf("named int = %d", got)
	}
	if got := ToInt64(myUint(6)); got != 6 {
		t.Fatalf("named uint = %d", got)
	}
	if got := ToInt64(myFloat(7.5)); got != 7 {
		t.Fatalf("named float = %d", got)
	}
	// Unsupported kind falls back to default.
	if got := ToInt64Default([]int{1}, -1); got != -1 {
		t.Fatalf("unsupported kind = %d", got)
	}
}

func TestToInt64StringEdgeCases(t *testing.T) {
	if got := ToInt64Default(nil, -1); got != -1 {
		t.Fatalf("nil = %d", got)
	}
	if got := ToInt64Default("   ", -1); got != -1 {
		t.Fatalf("blank = %d", got)
	}
	// Non-integer string that parses as float.
	if got := ToInt64("42.7"); got != 42 {
		t.Fatalf("float string = %d", got)
	}
	// Completely invalid string.
	if got := ToInt64Default("not-a-number", -1); got != -1 {
		t.Fatalf("invalid string = %d", got)
	}
}

func TestToFloat64Branches(t *testing.T) {
	if got := ToFloat64(float32(1.5)); got != 1.5 {
		t.Fatalf("float32 = %v", got)
	}
	if got := ToFloat64Default(nil, -1); got != -1 {
		t.Fatalf("nil = %v", got)
	}
	if got := ToFloat64Default("   ", -1); got != -1 {
		t.Fatalf("blank = %v", got)
	}
	if got := ToFloat64Default("nope", -1); got != -1 {
		t.Fatalf("invalid = %v", got)
	}
	// Integer value flows through the toInt64 fallback.
	if got := ToFloat64(int64(7)); got != 7 {
		t.Fatalf("int fallback = %v", got)
	}
	// Unsupported value falls back to default.
	if got := ToFloat64Default([]int{1}, -1); got != -1 {
		t.Fatalf("unsupported = %v", got)
	}
}

func TestToStringSpecialKinds(t *testing.T) {
	if got := ToStringWithOptions(errors.New("boom")); got != "boom" {
		t.Fatalf("error = %q", got)
	}
	if got := ToStringWithOptions(float32(1.5)); got != "1.5" {
		t.Fatalf("float32 = %q", got)
	}
	if got := ToStringWithOptions([]byte("bytes")); got != "bytes" {
		t.Fatalf("bytes = %q", got)
	}
	if got := ToStringDefaultWithOptions(nil, "def"); got != "def" {
		t.Fatalf("nil default = %q", got)
	}
	if got := ToStringDefaultWithOptions(7, "def"); got != "7" {
		t.Fatalf("non-nil default = %q", got)
	}
}

func TestToBoolDefaultFromInt(t *testing.T) {
	if !ToBoolDefault(int64(5), false) {
		t.Fatal("non-zero int should be true")
	}
	if ToBoolDefault(int64(0), true) {
		t.Fatal("zero int should be false")
	}
	if !ToBoolDefault("yes", false) {
		t.Fatal("yes should parse true")
	}
	if ToBoolDefault("garbage", false) {
		t.Fatal("garbage should fall back to def")
	}
	// Unsupported value uses default.
	if !ToBoolDefault([]int{1}, true) {
		t.Fatal("unsupported should use default")
	}
}

// TestAllOptionSettersApplied ensures every With* option installs a custom
// provider that is actually invoked.
func TestAllOptionSettersApplied(t *testing.T) {
	intCalled, floatCalled, boolCalled := false, false, false
	formatBoolCalled, formatFloatCalled := false, false

	opts := []Option{
		WithParseIntFunc(func(s string, base, bit int) (int64, error) {
			intCalled = true
			return strconv.ParseInt(s, base, bit)
		}),
		WithParseFloatFunc(func(s string, bit int) (float64, error) {
			floatCalled = true
			return strconv.ParseFloat(s, bit)
		}),
		WithBoolParser(func(string) (bool, error) {
			boolCalled = true
			return true, nil
		}),
		WithFormatBoolFunc(func(bool) string {
			formatBoolCalled = true
			return "B"
		}),
		WithFormatFloatFunc(func(float64, byte, int, int) string {
			formatFloatCalled = true
			return "F"
		}),
	}

	if got := ToInt64WithOptions("10", opts...); got != 10 || !intCalled {
		t.Fatalf("custom parseInt not used (got=%d called=%v)", got, intCalled)
	}
	if got := ToFloat64WithOptions("1.5", opts...); got != 1.5 || !floatCalled {
		t.Fatalf("custom parseFloat not used (got=%v called=%v)", got, floatCalled)
	}
	if got := ToBoolWithOptions("anything", opts...); !got || !boolCalled {
		t.Fatalf("custom parseBool not used (got=%v called=%v)", got, boolCalled)
	}
	if got := ToStringWithOptions(true, opts...); got != "B" || !formatBoolCalled {
		t.Fatalf("custom formatBool not used (got=%q)", got)
	}
	if got := ToStringWithOptions(1.5, opts...); got != "F" || !formatFloatCalled {
		t.Fatalf("custom formatFloat not used (got=%q)", got)
	}

	// nil options are ignored and defaults are retained.
	if got := ToInt64WithOptions("3", nil, WithParseIntFunc(nil)); got != 3 {
		t.Fatalf("nil options should keep defaults, got %d", got)
	}
}
