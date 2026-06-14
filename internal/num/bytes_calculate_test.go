package num

import (
	"math"
	"math/big"
	"reflect"
	"strconv"
	"testing"
)

func TestParseBytesValidityAndCalculate(t *testing.T) {
	if ParseInt("0x10") != 16 || ParseInt("123.56") != 123 || ParseLong("123.56") != 123 {
		t.Fatal("parse integer helpers failed")
	}
	if ParseFloat(".125") != 0.125 || ParseDouble("1,234.5") != 1234.5 || ParseIntDefault("bad", 7) != 7 {
		t.Fatal("parse float/default helpers failed")
	}
	b := ToBytes(0x01020304)
	if !reflect.DeepEqual(b, []byte{1, 2, 3, 4}) || ToInt(b) != 0x01020304 {
		t.Fatal("byte conversion failed")
	}
	if got := ToInt(ToBytes(-1)); got != -1 {
		t.Fatalf("signed byte round trip failed: %d", got)
	}
	unsigned, err := ToUnsignedByteArrayLen(4, big.NewInt(255))
	if err != nil || !reflect.DeepEqual(unsigned, []byte{0, 0, 0, 255}) || FromUnsignedByteArray(unsigned).Int64() != 255 {
		t.Fatal("unsigned bytes failed")
	}
	if IsValid(math.Inf(1)) || IsValidFloat32(float32(math.NaN())) || !IsValidNumber(1) {
		t.Fatal("valid number helpers failed")
	}
	result, err := Calculate("(0*1--3)-5/-4-(3*(-2.13))")
	if err != nil || math.Abs(result-10.64) > 1e-9 {
		t.Fatalf("Calculate: %v %v", result, err)
	}
	if ToDouble(float32(1.23)) != 1.23 || !IsOdd(3) || !IsEven(4) || !IsPowerOfTwo(1024) {
		t.Fatal("misc helpers failed")
	}
}

func TestByteUnsignedValidityAndExpressionEdges(t *testing.T) {
	if ToInt(nil) != 0 || ToInt([]byte{1, 2, 3}) != 0 {
		t.Fatal("ToInt short input should be zero")
	}
	if ToUnsignedByteArray(nil) != nil {
		t.Fatal("ToUnsignedByteArray nil should be nil")
	}
	if got := ToUnsignedByteArray(big.NewInt(0)); len(got) != 0 {
		t.Fatalf("ToUnsignedByteArray zero should be empty: %v", got)
	}
	if _, err := ToUnsignedByteArrayLen(1, big.NewInt(256)); err == nil {
		t.Fatal("ToUnsignedByteArrayLen should reject values that exceed requested length")
	}
	if got, err := ToUnsignedByteArrayLen(0, big.NewInt(0)); err != nil || len(got) != 0 {
		t.Fatalf("ToUnsignedByteArrayLen zero length/value = %v, %v", got, err)
	}
	if FromUnsignedByteArray(nil).Sign() != 0 || FromUnsignedByteArrayRange([]byte{1, 2, 3, 4}, 1, 2).Int64() != 0x0203 {
		t.Fatal("FromUnsignedByteArray cases failed")
	}
	if FromUnsignedByteArrayRange([]byte{1, 2}, -1, 1).Sign() != 0 || FromUnsignedByteArrayRange([]byte{1, 2}, 1, 3).Sign() != 0 {
		t.Fatal("FromUnsignedByteArrayRange invalid ranges should be zero")
	}
	if IsValidNumber(nil) || IsValidNumber(math.NaN()) || IsValidNumber(float32(math.Inf(-1))) || !IsValidNumber("not-a-number") {
		t.Fatal("IsValidNumber cases failed")
	}
	if IsValid(math.NaN()) || IsValid(math.Inf(1)) || !IsValid(1.23) || IsValidFloat32(float32(math.Inf(1))) || !IsValidFloat32(1.23) {
		t.Fatal("valid finite checks failed")
	}
	toDoubleCases := []struct {
		input any
		want  float64
	}{
		{float32(1.25), 1.25},
		{float64(2.5), 2.5},
		{int(-3), -3},
		{int64(4), 4},
		{uint64(5), 5},
		{"bad", 0},
	}
	for _, tt := range toDoubleCases {
		if got := ToDouble(tt.input); got != tt.want {
			t.Fatalf("ToDouble(%T) = %v, want %v", tt.input, got, tt.want)
		}
	}
	formatCalled := false
	parseCalled := false
	if got := ToDoubleWithOptions(float32(1.25),
		WithDoubleFormatFloatFunc(func(v float64, fmtByte byte, prec, bitSize int) string {
			formatCalled = true
			return strconv.FormatFloat(v*2, fmtByte, prec, bitSize)
		}),
		WithDoubleParseFloatFunc(func(s string, bitSize int) (float64, error) {
			parseCalled = true
			return strconv.ParseFloat(s, bitSize)
		}),
	); got != 2.5 || !formatCalled || !parseCalled {
		t.Fatalf("ToDoubleWithOptions = %v format=%v parse=%v", got, formatCalled, parseCalled)
	}
	calcCases := map[string]float64{
		"1 + 2 * 3":   7,
		"(1 + 2) * 3": 9,
		"10 % 4":      2,
		"--2 + +3":    5,
		" 3.5 / 2 ":   1.75,
	}
	for expr, want := range calcCases {
		got, err := Calculate(expr)
		if err != nil || math.Abs(got-want) > 1e-9 {
			t.Fatalf("Calculate(%q) = %v, %v, want %v", expr, got, err, want)
		}
	}
	calcParseCalled := false
	got, err := CalculateWithOptions("5 + 2", WithParseFloatFunc(func(s string, bitSize int) (float64, error) {
		calcParseCalled = true
		if s == "5" {
			return 5, nil
		}
		return strconv.ParseFloat(s, bitSize)
	}))
	if err != nil || got != 7 || !calcParseCalled {
		t.Fatalf("CalculateWithOptions = %v, %v called=%v", got, err, calcParseCalled)
	}
	invalidExpressions := []string{"", "1+", "(1+2", "1 2", "abc"}
	for _, expr := range invalidExpressions {
		if got, err := Calculate(expr); err == nil {
			t.Fatalf("Calculate(%q) should fail, got %v", expr, got)
		}
	}
	if secureIntn(0) != 0 || secureIntn(-1) != 0 {
		t.Fatal("secureIntn non-positive max should be zero")
	}
	for i := 0; i < 20; i++ {
		if got := secureIntn(3); got < 0 || got >= 3 {
			t.Fatalf("secureIntn result out of range: %d", got)
		}
	}
	if !IsOdd(-3) || IsOdd(-2) || !IsEven(-2) || IsEven(-3) {
		t.Fatal("odd/even negative cases failed")
	}
}
