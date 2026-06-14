package num

import "testing"

func TestSecureIntnBoundaries(t *testing.T) {
	if secureIntn(0) != 0 || secureIntn(-1) != 0 {
		t.Fatal("secureIntn non-positive max should be zero")
	}
	for i := 0; i < 20; i++ {
		if got := secureIntn(3); got < 0 || got >= 3 {
			t.Fatalf("secureIntn result out of range: %d", got)
		}
	}
}

func TestParityHelpers(t *testing.T) {
	if !IsOdd(3) || !IsEven(4) || !IsPowerOfTwo(1024) {
		t.Fatal("misc helpers failed")
	}
	if !IsOdd(-3) || IsOdd(-2) || !IsEven(-2) || IsEven(-3) {
		t.Fatal("odd/even negative cases failed")
	}
}
