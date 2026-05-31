package obj

import (
	"math"
	"reflect"
	"testing"
)

type sample struct {
	Name string
	Tags []string
}

func TestEqualLengthContainsAndEmpty(t *testing.T) {
	if !Equal(1, int64(1)) || NotEqual("a", "a") {
		t.Fatal("numeric or string equality failed")
	}
	if Length([]int{1, 2, 3}) != 3 || Length(10) != -1 {
		t.Fatal("length failed")
	}
	if !Contains([]int{1, 2, 3}, int64(2)) || !Contains("hello", "ell") {
		t.Fatal("contains failed")
	}
	if !IsEmpty(map[string]int{}) || IsEmpty(1) || !IsNotEmpty([]int{1}) {
		t.Fatal("empty checks failed")
	}
}

func TestDefaultsApplyAcceptAndAggregates(t *testing.T) {
	value := "go"
	if DefaultIfNil(&value, "x") != "go" || DefaultIfNil[string](nil, "x") != "x" {
		t.Fatal("DefaultIfNil failed")
	}
	if DefaultIfEmpty("", "x") != "x" || DefaultIfBlank("  ", "x") != "x" {
		t.Fatal("string defaults failed")
	}
	if got := Apply(&value, func(s string) int { return len(s) }); got != 2 {
		t.Fatalf("Apply: %d", got)
	}
	called := false
	Accept(&value, func(string) { called = true })
	if !called {
		t.Fatal("Accept not called")
	}
	if EmptyCount(nil, "", []int{}, 1) != 3 || !HasNil(1, nil) || !HasEmpty(1, "") {
		t.Fatal("aggregate checks failed")
	}
	if !IsAllEmpty(nil, "") || !IsAllNotEmpty(1, "x") {
		t.Fatal("all checks failed")
	}
}

func TestCloneSerializeCompareAndType(t *testing.T) {
	src := sample{Name: "n", Tags: []string{"a"}}
	clone, err := Clone(src)
	if err != nil {
		t.Fatalf("Clone: %v", err)
	}
	clone.Tags[0] = "b"
	if src.Tags[0] != "a" {
		t.Fatal("clone is not independent")
	}
	data, err := Serialize(src)
	if err != nil {
		t.Fatalf("Serialize: %v", err)
	}
	var out sample
	if err := Deserialize(data, &out); err != nil || !reflect.DeepEqual(out, src) {
		t.Fatalf("Deserialize: %#v %v", out, err)
	}
	a, b := 1, 2
	if Compare(&a, &b) >= 0 || CompareNull[int](nil, &b, true) <= 0 {
		t.Fatal("compare failed")
	}
	if TypeName(src) == "" || ToString(nil) != "null" {
		t.Fatal("type or string failed")
	}
}

func TestBasicAndValidNumber(t *testing.T) {
	if !IsBasicType("x") || IsBasicType(sample{}) {
		t.Fatal("basic type check failed")
	}
	if !IsValidIfNumber(1) || IsValidIfNumber(math.NaN()) || IsValidIfNumber(math.Inf(1)) {
		t.Fatal("valid number check failed")
	}
}
