package hash

import "testing"

func TestHashFunctions(t *testing.T) {
	if MD5Hex("abc") != "900150983cd24fb0d6963f7d28e17f72" {
		t.Fatalf("MD5 failed")
	}
	if SHA1Hex("abc") != "a9993e364706816aba3e25717850c26c9cd0d89d" {
		t.Fatalf("SHA1 failed")
	}
	if SHA256Hex("abc") != "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad" {
		t.Fatalf("SHA256 failed")
	}
	if FnvHash("abc") == 0 {
		t.Fatalf("FnvHash zero")
	}
	if AdditiveHash("abc", 31) < 0 {
		t.Fatalf("AdditiveHash failed")
	}
}
