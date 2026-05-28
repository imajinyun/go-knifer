package vcrypto_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/imajinyun/go-knifer/vcrypto"
)

func TestDigestAndHMAC(t *testing.T) {
	if got := vcrypto.MD5Hex("hello"); got != "5d41402abc4b2a76b9719d911017c592" {
		t.Fatalf("MD5Hex() = %s", got)
	}
	if got := vcrypto.MD5HexBytes([]byte("hello")); got != "5d41402abc4b2a76b9719d911017c592" {
		t.Fatalf("MD5HexBytes() = %s", got)
	}
	if got := vcrypto.SHA1Hex("hello"); got != "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d" {
		t.Fatalf("SHA1Hex() = %s", got)
	}
	if got := vcrypto.SHA256Hex("hello"); got != "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824" {
		t.Fatalf("SHA256Hex() = %s", got)
	}
	if got := vcrypto.SHA512Hex("hello"); got == "" {
		t.Fatal("SHA512Hex() is empty")
	}
	if got := vcrypto.HMACSHA256Hex([]byte("key"), []byte("hello")); got != "9307b3b915efb5171ff14d8cb55fbcc798c6c0ef1456d66ded1a6aa723a58b7b" {
		t.Fatalf("HMACSHA256Hex() = %s", got)
	}
}

func TestAESRoundTripAndErrors(t *testing.T) {
	key, err := vcrypto.GenerateAESKey(16)
	if err != nil {
		t.Fatal(err)
	}
	if len(key) != 16 {
		t.Fatalf("GenerateAESKey len = %d", len(key))
	}
	if _, err := vcrypto.GenerateAESKey(15); !errors.Is(err, vcrypto.ErrInvalidKey) {
		t.Fatalf("GenerateAESKey invalid error = %v", err)
	}
	iv := []byte("1234567890123456")
	plain := []byte("hutool crypto facade")
	cipherText, err := vcrypto.AESEncryptCBC(plain, key, iv)
	if err != nil {
		t.Fatal(err)
	}
	out, err := vcrypto.AESDecryptCBC(cipherText, key, iv)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out, plain) {
		t.Fatalf("AESDecryptCBC() = %q", out)
	}
	if _, err := vcrypto.AESEncryptCBC(plain, key, []byte("bad")); !errors.Is(err, vcrypto.ErrInvalidIV) {
		t.Fatalf("AESEncryptCBC invalid iv error = %v", err)
	}

	nonce := []byte("123456789012")
	cipherText, err = vcrypto.AESEncryptGCM(plain, key, nonce, []byte("aad"))
	if err != nil {
		t.Fatal(err)
	}
	out, err = vcrypto.AESDecryptGCM(cipherText, key, nonce, []byte("aad"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out, plain) {
		t.Fatalf("AESDecryptGCM() = %q", out)
	}
}

func TestRSAAndPEM(t *testing.T) {
	priv, err := vcrypto.GenerateRSAKey(1024)
	if err != nil {
		t.Fatal(err)
	}
	pubPEM, err := vcrypto.PublicKeyToPEM(&priv.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	pub, err := vcrypto.ParseRSAPublicKeyPEM(pubPEM)
	if err != nil {
		t.Fatal(err)
	}
	parsedPriv, err := vcrypto.ParseRSAPrivateKeyPEM(vcrypto.PrivateKeyToPEM(priv))
	if err != nil {
		t.Fatal(err)
	}
	plain := []byte("rsa message")
	cipherText, err := vcrypto.RSAEncryptOAEP(plain, pub, nil)
	if err != nil {
		t.Fatal(err)
	}
	out, err := vcrypto.RSADecryptOAEP(cipherText, parsedPriv, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out, plain) {
		t.Fatalf("RSADecryptOAEP() = %q", out)
	}
}
