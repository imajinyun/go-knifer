package conv

import "testing"

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
