package conv

import (
	"errors"
	"math"
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

func TestApplyOptionsFallsBackWhenCustomOptionClearsProviders(t *testing.T) {
	clearProviders := func(c *config) {
		c.parseBool = nil
		c.parseInt = nil
		c.parseFloat = nil
		c.formatBool = nil
		c.formatFloat = nil
	}

	if got := ToBoolWithOptions("on", clearProviders); !got {
		t.Fatal("nil parseBool provider should fall back to default")
	}
	if got := ToInt64WithOptions("12", clearProviders); got != 12 {
		t.Fatalf("nil parseInt provider should fall back to default, got %d", got)
	}
	if got := ToFloat64WithOptions("1.25", clearProviders); got != 1.25 {
		t.Fatalf("nil parseFloat provider should fall back to default, got %v", got)
	}
	if got := ToStringWithOptions(true, clearProviders); got != "true" {
		t.Fatalf("nil formatBool provider should fall back to default, got %q", got)
	}
	if got := ToStringWithOptions(1.25, clearProviders); got != "1.25" {
		t.Fatalf("nil formatFloat provider should fall back to default, got %q", got)
	}
}

func TestStrictIntegerConversionsCoverConcreteKinds(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want int64
	}{
		{name: "int", in: int(1), want: 1},
		{name: "int8", in: int8(2), want: 2},
		{name: "int16", in: int16(3), want: 3},
		{name: "int32", in: int32(4), want: 4},
		{name: "int64", in: int64(5), want: 5},
		{name: "uint", in: uint(6), want: 6},
		{name: "uint8", in: uint8(7), want: 7},
		{name: "uint16", in: uint16(8), want: 8},
		{name: "uint32", in: uint32(9), want: 9},
		{name: "uint64", in: uint64(10), want: 10},
		{name: "float32", in: float32(11), want: 11},
		{name: "float64", in: float64(12), want: 12},
		{name: "bool true", in: true, want: 1},
		{name: "bool false", in: false, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToInt64E(tt.in)
			if err != nil {
				t.Fatalf("ToInt64E(%T) unexpected error = %v", tt.in, err)
			}
			if got != tt.want {
				t.Fatalf("ToInt64E(%T) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestStrictIntegerConversionsCoverReflectAndStringFallbacks(t *testing.T) {
	type myUint uint
	type myUint64 uint64
	type myFloat float64
	type myString string

	tests := []struct {
		name    string
		in      any
		want    int64
		wantErr bool
	}{
		{name: "named uint", in: myUint(13), want: 13},
		{name: "named uint64 overflow", in: myUint64(uint64(math.MaxInt64) + 1), wantErr: true},
		{name: "named float", in: myFloat(14), want: 14},
		{name: "named string int", in: myString("15"), want: 15},
		{name: "strict float string truncates integral value", in: "16.0", want: 16},
		{name: "strict float string truncates fractional value", in: "16.5", want: 16},
		{name: "unsupported kind", in: []int{1}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToInt64E(tt.in)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidConversion) {
					t.Fatalf("ToInt64E(%T) error = %v, want ErrInvalidConversion", tt.in, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("ToInt64E(%T) unexpected error = %v", tt.in, err)
			}
			if got != tt.want {
				t.Fatalf("ToInt64E(%T) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestToBoolDefaultReflectAndProviderErrorPaths(t *testing.T) {
	type myBool bool
	type myString string

	if got := ToBoolDefaultWithOptions(myBool(true), false); !got {
		t.Fatal("named bool should use reflect bool branch")
	}
	if got := ToBoolDefaultWithOptions(myString("yes"), false); !got {
		t.Fatal("named string should parse through reflect string branch")
	}
	if got := ToBoolDefaultWithOptions(myString("bad"), true, WithBoolParser(func(string) (bool, error) {
		return false, errors.New("bad bool")
	})); !got {
		t.Fatal("named string parser errors should return default")
	}
	if got := ToBoolDefaultWithOptions(nil, true); !got {
		t.Fatal("nil input should return default")
	}
}

func TestToFloat64ReflectStringAndFallbackErrors(t *testing.T) {
	type myFloat float64
	type myString string

	if got := ToFloat64WithOptions(myFloat(2.5)); got != 2.5 {
		t.Fatalf("named float = %v, want 2.5", got)
	}
	if got := ToFloat64WithOptions(myString("3.5")); got != 3.5 {
		t.Fatalf("named string = %v, want 3.5", got)
	}
	if got := ToFloat64DefaultWithOptions(myString("   "), -1); got != -1 {
		t.Fatalf("blank named string = %v, want default", got)
	}
	if got := ToFloat64DefaultWithOptions(myString("bad"), -1); got != -1 {
		t.Fatalf("invalid named string = %v, want default", got)
	}
}
