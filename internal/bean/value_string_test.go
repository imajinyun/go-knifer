package bean

import (
	"math"
	"reflect"
	"testing"
)

type valueStringSample struct {
	Name string
}

func TestValueString(t *testing.T) {
	// valid string value
	v := reflect.ValueOf("hello")
	if got := valueString(v); got != "hello" {
		t.Fatalf("valueString(string) = %q, want %q", got, "hello")
	}

	// int value (CanInterface)
	v = reflect.ValueOf(42)
	if got := valueString(v); got != "42" {
		t.Fatalf("valueString(int) = %q, want %q", got, "42")
	}

	// invalid value
	v = reflect.Value{}
	if got := valueString(v); got != "" {
		t.Fatalf("valueString(invalid) = %q, want empty", got)
	}

	// pointer to string
	s := "pointer"
	v = reflect.ValueOf(&s)
	if got := valueString(v); got != "pointer" {
		t.Fatalf("valueString(pointer) = %q, want %q", got, "pointer")
	}

	// struct pointer (CanInterface true)
	v = reflect.ValueOf(&valueStringSample{Name: "test"})
	if got := valueString(v); got == "" {
		t.Fatal("valueString(struct) should not be empty")
	}
}

func TestWeakScalarConversionBoundaries(t *testing.T) {
	cfg := NewOptions()

	boolCases := []struct {
		name string
		in   any
		want bool
	}{
		{name: "bool", in: true, want: true},
		{name: "string", in: "off", want: false},
		{name: "int", in: -1, want: true},
		{name: "uint", in: uint(0), want: false},
		{name: "float", in: 0.25, want: true},
		{name: "nil pointer", in: (*string)(nil), want: false},
	}
	for _, tt := range boolCases {
		t.Run("bool "+tt.name, func(t *testing.T) {
			got, err := valueBool(reflect.ValueOf(tt.in), cfg)
			if err != nil {
				t.Fatalf("valueBool(%T) error = %v", tt.in, err)
			}
			if got != tt.want {
				t.Fatalf("valueBool(%T) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
	if _, err := valueBool(reflect.ValueOf([]int{1}), cfg); err == nil {
		t.Fatal("valueBool unsupported kind error = nil")
	}

	intCases := []struct {
		name string
		in   any
		bits int
		want int64
	}{
		{name: "int", in: int64(-2), bits: 64, want: -2},
		{name: "uint", in: uint(3), bits: 64, want: 3},
		{name: "float", in: 4.9, bits: 64, want: 4},
		{name: "bool", in: true, bits: 64, want: 1},
		{name: "blank string", in: " ", bits: 64, want: 0},
		{name: "float string", in: "5.9", bits: 64, want: 5},
	}
	for _, tt := range intCases {
		t.Run("int "+tt.name, func(t *testing.T) {
			got, err := valueInt(reflect.ValueOf(tt.in), tt.bits, cfg)
			if err != nil {
				t.Fatalf("valueInt(%T) error = %v", tt.in, err)
			}
			if got != tt.want {
				t.Fatalf("valueInt(%T) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
	if _, err := valueInt(reflect.ValueOf(uint64(math.MaxUint64)), 64, cfg); err == nil {
		t.Fatal("valueInt uint64 overflow error = nil")
	}
	if _, err := valueInt(reflect.ValueOf(int16(128)), 8, cfg); err == nil {
		t.Fatal("valueInt int8 overflow error = nil")
	}
	if _, err := valueInt(reflect.ValueOf([]int{1}), 64, cfg); err == nil {
		t.Fatal("valueInt unsupported kind error = nil")
	}

	uintCases := []struct {
		name string
		in   any
		bits int
		want uint64
	}{
		{name: "uint", in: uint64(2), bits: 64, want: 2},
		{name: "int", in: int64(3), bits: 64, want: 3},
		{name: "float", in: 4.9, bits: 64, want: 4},
		{name: "bool", in: true, bits: 64, want: 1},
		{name: "blank string", in: " ", bits: 64, want: 0},
		{name: "float string", in: "5.9", bits: 64, want: 5},
	}
	for _, tt := range uintCases {
		t.Run("uint "+tt.name, func(t *testing.T) {
			got, err := valueUint(reflect.ValueOf(tt.in), tt.bits, cfg)
			if err != nil {
				t.Fatalf("valueUint(%T) error = %v", tt.in, err)
			}
			if got != tt.want {
				t.Fatalf("valueUint(%T) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
	if _, err := valueUint(reflect.ValueOf(int64(-1)), 64, cfg); err == nil {
		t.Fatal("valueUint negative int error = nil")
	}
	if _, err := valueUint(reflect.ValueOf(-1.5), 64, cfg); err == nil {
		t.Fatal("valueUint negative float error = nil")
	}
	if _, err := valueUint(reflect.ValueOf(uint16(256)), 8, cfg); err == nil {
		t.Fatal("valueUint uint8 overflow error = nil")
	}
	if _, err := valueUint(reflect.ValueOf([]int{1}), 64, cfg); err == nil {
		t.Fatal("valueUint unsupported kind error = nil")
	}

	floatCases := []struct {
		name string
		in   any
		want float64
	}{
		{name: "float", in: 1.25, want: 1.25},
		{name: "int", in: int64(-2), want: -2},
		{name: "uint", in: uint(3), want: 3},
		{name: "bool true", in: true, want: 1},
		{name: "bool false", in: false, want: 0},
		{name: "blank string", in: " ", want: 0},
		{name: "string", in: "4.5", want: 4.5},
	}
	for _, tt := range floatCases {
		t.Run("float "+tt.name, func(t *testing.T) {
			got, err := valueFloat(reflect.ValueOf(tt.in), 64, cfg)
			if err != nil {
				t.Fatalf("valueFloat(%T) error = %v", tt.in, err)
			}
			if got != tt.want {
				t.Fatalf("valueFloat(%T) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
	if _, err := valueFloat(reflect.ValueOf([]int{1}), 64, cfg); err == nil {
		t.Fatal("valueFloat unsupported kind error = nil")
	}
}

func TestFloatIntegerFitBoundaries(t *testing.T) {
	if !floatFitsIntBits(float64(math.MaxInt8), 8) {
		t.Fatal("MaxInt8 should fit int8")
	}
	if floatFitsIntBits(float64(math.MaxInt8)+1, 8) {
		t.Fatal("MaxInt8+1 should not fit int8")
	}
	if !floatFitsIntBits(float64(math.MinInt8), 8) {
		t.Fatal("MinInt8 should fit int8")
	}
	if floatFitsIntBits(float64(math.MinInt8)-1, 8) {
		t.Fatal("MinInt8-1 should not fit int8")
	}
	if floatFitsIntBits(math.NaN(), 64) || floatFitsIntBits(math.Inf(1), 64) {
		t.Fatal("NaN/Inf should not fit int64")
	}

	if !floatFitsUintBits(float64(math.MaxUint8), 8) {
		t.Fatal("MaxUint8 should fit uint8")
	}
	if floatFitsUintBits(float64(math.MaxUint8)+1, 8) {
		t.Fatal("MaxUint8+1 should not fit uint8")
	}
	if floatFitsUintBits(-1, 64) || floatFitsUintBits(math.Inf(1), 64) || floatFitsUintBits(math.NaN(), 64) {
		t.Fatal("negative/Inf/NaN should not fit uint")
	}
}
