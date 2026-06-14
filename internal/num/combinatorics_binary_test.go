package num

import (
	"math/big"
	"testing"
)

func TestCombinatoricsGcdLcmAndBinaryEdges(t *testing.T) {
	if Sqrt(81) != 9 || ProcessMultiple(7, 5) != 21 || Divisor(24, 18) != 6 || Multiple(4, 6) != 12 {
		t.Fatal("math helpers failed")
	}
	if ProcessMultiple(3, 5) != 0 || ProcessMultiple(5, -1) != 0 || ProcessMultiple(5, 5) != 1 || ProcessMultiple(5, 2) != 10 {
		t.Fatal("ProcessMultiple edge cases failed")
	}
	if Divisor(-24, 18) != 6 || Divisor(7, 0) != 7 || Multiple(-4, 6) != 12 || Multiple(0, 6) != 0 {
		t.Fatal("gcd/lcm edge cases failed")
	}
	binaryCases := map[any]string{
		int(5):            "101",
		int8(-2):          "11111110",
		int16(-2):         "1111111111111110",
		int32(-2):         "-10",
		int64(9):          "1001",
		uint8(7):          "111",
		float64(1):        "0011111111110000000000000000000000000000000000000000000000000000",
		big.NewInt(10):    "1010",
		(*big.Int)(nil):   "",
		complex64(1 + 2i): "",
	}
	for input, want := range binaryCases {
		if got := GetBinaryStr(input); got != want {
			t.Fatalf("GetBinaryStr(%T %[1]v) = %q, want %q", input, got, want)
		}
	}
	if _, err := BinaryToInt("102"); err == nil {
		t.Fatal("BinaryToInt should reject invalid binary strings")
	}
	if _, err := BinaryToLong("2"); err == nil {
		t.Fatal("BinaryToLong should reject invalid binary strings")
	}
}
