package base

import "testing"

// 对应 hutool-core ConvertTest。

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
