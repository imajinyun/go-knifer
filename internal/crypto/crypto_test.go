package crypto

import (
	"bytes"
	"testing"
)

func TestDigestAndHMAC(t *testing.T) {
	if got := MD5Hex([]byte("hello")); got != "5d41402abc4b2a76b9719d911017c592" {
		t.Fatalf("MD5Hex() = %s", got)
	}
	if got := SHA256Hex([]byte("hello")); got != "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824" {
		t.Fatalf("SHA256Hex() = %s", got)
	}
	if got := HMACSHA256Hex([]byte("key"), []byte("hello")); got == "" {
		t.Fatal("HMACSHA256Hex() is empty")
	}
}

func TestAESCBCAndGCM(t *testing.T) {
	key := []byte("1234567890123456")
	iv := []byte("abcdefghijklmnop")
	plain := []byte("hello hutool")
	cipherText, err := AESEncryptCBC(plain, key, iv)
	if err != nil {
		t.Fatal(err)
	}
	out, err := AESDecryptCBC(cipherText, key, iv)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out, plain) {
		t.Fatalf("AESDecryptCBC() = %q", out)
	}

	nonce := []byte("123456789012")
	cipherText, err = AESEncryptGCM(plain, key, nonce, nil)
	if err != nil {
		t.Fatal(err)
	}
	out, err = AESDecryptGCM(cipherText, key, nonce, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out, plain) {
		t.Fatalf("AESDecryptGCM() = %q", out)
	}
}

func TestRSAOAEPAndPEM(t *testing.T) {
	priv, err := GenerateRSAKey(1024)
	if err != nil {
		t.Fatal(err)
	}
	plain := []byte("hello")
	cipherText, err := RSAEncryptOAEP(plain, &priv.PublicKey, nil)
	if err != nil {
		t.Fatal(err)
	}
	out, err := RSADecryptOAEP(cipherText, priv, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out, plain) {
		t.Fatalf("RSADecryptOAEP() = %q", out)
	}
	parsed, err := ParseRSAPrivateKeyPEM(PrivateKeyToPEM(priv))
	if err != nil {
		t.Fatal(err)
	}
	if parsed.N.Cmp(priv.N) != 0 {
		t.Fatal("parsed private key mismatch")
	}
}
