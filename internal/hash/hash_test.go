package hash

import "testing"

func TestHashFunctions(t *testing.T) {
	if FnvHash("abc") == 0 {
		t.Fatalf("FnvHash zero")
	}
	if AdditiveHash("abc", 31) < 0 {
		t.Fatalf("AdditiveHash failed")
	}
}
