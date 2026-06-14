package num

import "testing"

func TestPartPowAndNumericHelperEdges(t *testing.T) {
	if !IsBeside(2, 1) || !IsBeside[int64](1, 2) || IsBeside(1, 3) {
		t.Fatal("IsBeside cases failed")
	}
	if PartValueWithMode(10, 0, true) != 0 || PartValueWithMode(10, 3, false) != 3 || PartValueWithMode(10, 3, true) != 4 {
		t.Fatal("PartValueWithMode cases failed")
	}
	if Pow(2, 3) != 8 || Pow(2, -3) != 0.13 || PowWithMode(2, -3, 2, RoundDown) != 0.12 {
		t.Fatal("Pow cases failed")
	}
	if IsPowerOfTwo(0) || IsPowerOfTwo(-2) || !IsPowerOfTwo(1) || !IsPowerOfTwo(1024) || IsPowerOfTwo(1023) {
		t.Fatal("IsPowerOfTwo cases failed")
	}
}
