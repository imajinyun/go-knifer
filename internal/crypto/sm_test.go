package crypto

import (
	"bytes"
	"encoding/hex"
	"errors"
	"testing"

	knifer "github.com/imajinyun/knifer-go"
)

func TestSM3KnownVector(t *testing.T) {
	got := SM3Hex([]byte("abc"))
	const want = "66c7f0f462eeedd9d1f2d46bdc10e4e24167c4875cf2f7a2297da02b8f4ba8e0"
	if got != want {
		t.Fatalf("SM3Hex(abc) = %q, want %q", got, want)
	}
	sum := SM3([]byte("abc"))
	if !SM3Equal(sum, append([]byte(nil), sum...)) {
		t.Fatal("SM3Equal should match identical digest")
	}
	if SM3Equal(sum, bytes.Repeat([]byte{0}, len(sum))) {
		t.Fatal("SM3Equal should reject different digest")
	}
	if len(HMACSM3Bytes([]byte("key"), []byte("data"))) != 32 {
		t.Fatal("HMACSM3Bytes length should be 32")
	}
}

func TestSM4CBCAndECBRoundTrip(t *testing.T) {
	key := []byte("1234567890abcdef")
	iv := []byte("abcdef1234567890")
	plain := []byte("hello sm4")

	cbc, err := SM4EncryptCBC(plain, key, iv)
	if err != nil {
		t.Fatalf("SM4EncryptCBC() = %v", err)
	}
	got, err := SM4DecryptCBC(cbc, key, iv)
	if err != nil {
		t.Fatalf("SM4DecryptCBC() = %v", err)
	}
	if !bytes.Equal(got, plain) {
		t.Fatalf("SM4DecryptCBC() = %q, want %q", got, plain)
	}

	ecb, err := SM4EncryptECB(plain, key)
	if err != nil {
		t.Fatalf("SM4EncryptECB() = %v", err)
	}
	got, err = SM4DecryptECB(ecb, key)
	if err != nil {
		t.Fatalf("SM4DecryptECB() = %v", err)
	}
	if !bytes.Equal(got, plain) {
		t.Fatalf("SM4DecryptECB() = %q, want %q", got, plain)
	}
}

func TestSM4GCMRoundTripAndErrors(t *testing.T) {
	key := []byte("1234567890abcdef")
	nonce := []byte("123456789012")
	aad := []byte("aad")
	plain := []byte("hello sm4 gcm")

	cipherText, err := SM4EncryptGCM(plain, key, nonce, aad)
	if err != nil {
		t.Fatalf("SM4EncryptGCM() = %v", err)
	}
	got, err := SM4DecryptGCM(cipherText, key, nonce, aad)
	if err != nil {
		t.Fatalf("SM4DecryptGCM() = %v", err)
	}
	if !bytes.Equal(got, plain) {
		t.Fatalf("SM4DecryptGCM() = %q, want %q", got, plain)
	}
	if _, err := SM4DecryptGCM(cipherText, key, nonce, []byte("wrong")); !errors.Is(err, ErrInvalidCipherText) || !errors.Is(err, knifer.ErrCodeInvalidInput) {
		t.Fatalf("SM4DecryptGCM(wrong aad) = %v, want invalid cipher text/input", err)
	}
	if _, err := SM4EncryptGCM(plain, []byte("short"), nonce, aad); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("SM4EncryptGCM(short key) = %v, want invalid key", err)
	}
	if _, err := SM4EncryptGCM(plain, key, []byte("short"), aad); !errors.Is(err, ErrInvalidIV) {
		t.Fatalf("SM4EncryptGCM(short nonce) = %v, want invalid iv", err)
	}

	nonce, sealed, err := SM4SealGCMWithOptions(
		plain,
		key,
		aad,
		WithSM4RandomOptions(WithRandomReader(bytes.NewReader(bytes.Repeat([]byte{0x42}, 12)))),
	)
	if err != nil {
		t.Fatalf("SM4SealGCMWithOptions() = %v", err)
	}
	if hex.EncodeToString(nonce) != "424242424242424242424242" {
		t.Fatalf("nonce = %x", nonce)
	}
	if got, err := SM4DecryptGCM(sealed, key, nonce, aad); err != nil || !bytes.Equal(got, plain) {
		t.Fatalf("SM4SealGCM decrypt = %q, %v", got, err)
	}
}

func TestSM2RoundTripSignVerifyAndPEM(t *testing.T) {
	priv, err := GenSM2Key()
	if err != nil {
		t.Fatalf("GenSM2Key() = %v", err)
	}
	plain := []byte("hello sm2")
	cipherText, err := SM2Encrypt(plain, &priv.PublicKey)
	if err != nil {
		t.Fatalf("SM2Encrypt() = %v", err)
	}
	got, err := SM2Decrypt(cipherText, priv)
	if err != nil {
		t.Fatalf("SM2Decrypt() = %v", err)
	}
	if !bytes.Equal(got, plain) {
		t.Fatalf("SM2Decrypt() = %q, want %q", got, plain)
	}

	uid := []byte("knifer-go-sm2")
	sig, err := SM2SignWithOptions(plain, priv, WithSM2UID(uid))
	if err != nil {
		t.Fatalf("SM2SignWithOptions() = %v", err)
	}
	if err := SM2VerifyWithOptions(plain, sig, &priv.PublicKey, WithSM2UID(uid)); err != nil {
		t.Fatalf("SM2VerifyWithOptions() = %v", err)
	}
	if err := SM2VerifyWithOptions([]byte("tampered"), sig, &priv.PublicKey, WithSM2UID(uid)); !errors.Is(err, ErrInvalidSM2Signature) {
		t.Fatalf("SM2VerifyWithOptions(tampered) = %v, want invalid signature", err)
	}

	privPEM, err := SM2PrivateKeyToPEM(priv)
	if err != nil {
		t.Fatalf("SM2PrivateKeyToPEM() = %v", err)
	}
	parsedPriv, err := ParseSM2PrivateKeyPEM(privPEM)
	if err != nil {
		t.Fatalf("ParseSM2PrivateKeyPEM() = %v", err)
	}
	pubPEM, err := SM2PublicKeyToPEM(&parsedPriv.PublicKey)
	if err != nil {
		t.Fatalf("SM2PublicKeyToPEM() = %v", err)
	}
	parsedPub, err := ParseSM2PublicKeyPEM(pubPEM)
	if err != nil {
		t.Fatalf("ParseSM2PublicKeyPEM() = %v", err)
	}
	if err := SM2VerifyWithOptions(plain, sig, parsedPub, WithSM2UID(uid)); err != nil {
		t.Fatalf("SM2VerifyWithOptions(parsed pub) = %v", err)
	}
}

func TestSM2InvalidInputs(t *testing.T) {
	if _, err := SM2Encrypt([]byte("data"), nil); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("SM2Encrypt(nil) = %v, want invalid key", err)
	}
	if _, err := SM2Decrypt([]byte("bad"), nil); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("SM2Decrypt(nil) = %v, want invalid key", err)
	}
	if _, err := SM2Sign([]byte("data"), nil); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("SM2Sign(nil) = %v, want invalid key", err)
	}
	if err := SM2Verify([]byte("data"), []byte("sig"), nil); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("SM2Verify(nil) = %v, want invalid key", err)
	}
	if _, err := ParseSM2PrivateKeyPEM([]byte("bad")); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("ParseSM2PrivateKeyPEM(bad) = %v, want invalid key", err)
	}
	if _, err := ParseSM2PublicKeyPEM([]byte("bad")); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("ParseSM2PublicKeyPEM(bad) = %v, want invalid key", err)
	}
}
