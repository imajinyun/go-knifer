package vconv

import "testing"

func TestConvFacade(t *testing.T) {
	if ToString(12) != "12" || ToStringDefault(nil, "x") != "x" {
		t.Fatal("string conversion failed")
	}
	if ToInt("12") != 12 || ToIntDefault("bad", 7) != 7 {
		t.Fatal("int conversion failed")
	}
	if ToInt64(true) != 1 || ToFloat64("3.5") != 3.5 {
		t.Fatal("number conversion failed")
	}
	if !ToBool("yes") || ToBoolDefault("bad", true) != true {
		t.Fatal("bool conversion failed")
	}
	if string(ToBytes("go")) != "go" {
		t.Fatal("bytes conversion failed")
	}
}
