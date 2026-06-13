package crypto

import (
	"crypto/sha256"
	"testing"
)

func TestHMAC(t *testing.T) {
	if got := HMACSHA256Hex([]byte("key"), []byte("hello")); got == "" {
		t.Fatal("HMACSHA256Hex() is empty")
	}
	mac := HMACBytes(sha256.New, []byte("key"), []byte("hello"))
	if !HMACEqual(mac, HMACBytes(sha256.New, []byte("key"), []byte("hello"))) {
		t.Fatal("HMACEqual() returned false for identical MAC values")
	}
	if !ConstantTimeEqual([]byte("same"), []byte("same")) || ConstantTimeEqual([]byte("same"), []byte("diff")) {
		t.Fatal("ConstantTimeEqual() returned unexpected result")
	}
}
