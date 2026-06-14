package conv

import "testing"

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
