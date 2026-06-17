package obj

import (
	"math"
	"testing"
)

func TestCompareTypeAndString(t *testing.T) {
	src := sample{Name: "n", Tags: []string{"a"}}
	a, b := 1, 2
	if Compare(&a, &b) >= 0 || CompareNull[int](nil, &b, true) <= 0 {
		t.Fatal("compare failed")
	}
	if TypeName(src) == "" || ToString(nil) != "null" {
		t.Fatal("type or string failed")
	}
}

func TestCompareNullAndTypeBoundaries(t *testing.T) {
	a, b := 2, 1
	if CompareNull[int](nil, nil, true) != 0 || CompareNull[int](nil, &b, false) >= 0 || CompareNull(&a, nil, false) <= 0 {
		t.Fatal("CompareNull nil ordering failed")
	}
	if CompareNull(&a, &b, true) <= 0 || CompareNull(&b, &a, true) >= 0 || CompareNull(&a, &a, true) != 0 {
		t.Fatal("CompareNull value ordering failed")
	}
	if TypeOf(nil) != nil || TypeName(nil) != "" || ToString(42) != "42" {
		t.Fatal("type or string boundary failed")
	}
}

func TestBasicAndValidNumber(t *testing.T) {
	if !IsBasicType("x") || IsBasicType(sample{}) {
		t.Fatal("basic type check failed")
	}
	if !IsValidIfNumber(1) || IsValidIfNumber(math.NaN()) || IsValidIfNumber(math.Inf(1)) {
		t.Fatal("valid number check failed")
	}
	if IsBasicType(nil) || IsValidIfNumber(float32(math.NaN())) || IsValidIfNumber(float32(math.Inf(-1))) {
		t.Fatal("basic or float32 validity boundary failed")
	}
}
