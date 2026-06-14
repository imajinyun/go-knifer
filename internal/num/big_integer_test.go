package num

import "testing"

func TestBigIntegerParsingEdges(t *testing.T) {
	bigIntCases := map[string]int64{"42": 42, "-0x10": -16, "#10": 16, "010": 8}
	for input, want := range bigIntCases {
		got, ok := NewBigInteger(input)
		if !ok || got.Int64() != want {
			t.Fatalf("NewBigInteger(%q) = %v/%v, want %d", input, got, ok, want)
		}
	}
	if got, ok := NewBigInteger(""); ok || got != nil {
		t.Fatalf("NewBigInteger blank = %v/%v", got, ok)
	}
	if got, ok := NewBigInteger("0xzz"); ok || got != nil {
		t.Fatalf("NewBigInteger invalid = %v/%v", got, ok)
	}
}
