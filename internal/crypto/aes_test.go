package crypto

import (
	"bytes"
	"testing"
)

func TestAESGCM(t *testing.T) {
	key := []byte("1234567890123456")
	plain := []byte("hello crypto")
	nonce := []byte("123456789012")
	generatedNonce, sealed, err := AESSealGCMWithOptions(
		plain,
		key,
		[]byte("aad"),
		WithGCMRandomOptions(WithRandomReader(bytes.NewReader([]byte("abcdefghijkl")))),
	)
	if err != nil {
		t.Fatalf("AESSealGCMWithOptions() error = %v", err)
	}
	if !bytes.Equal(generatedNonce, []byte("abcdefghijkl")) {
		t.Fatalf("AESSealGCMWithOptions() nonce = %q", generatedNonce)
	}
	opened, err := AESOpenGCM(sealed, key, generatedNonce, []byte("aad"))
	if err != nil {
		t.Fatalf("AESOpenGCM() error = %v", err)
	}
	if !bytes.Equal(opened, plain) {
		t.Fatalf("AESOpenGCM() = %q", opened)
	}

	cipherText, err := AESEncryptGCM(plain, key, nonce, nil)
	if err != nil {
		t.Fatal(err)
	}
	out, err := AESDecryptGCM(cipherText, key, nonce, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out, plain) {
		t.Fatalf("AESDecryptGCM() = %q", out)
	}

	customNonce := []byte("1234567890123456")
	cipherText, err = AESEncryptGCMWithOptions(plain, key, customNonce, []byte("aad"), WithGCMNonceSize(len(customNonce)))
	if err != nil {
		t.Fatalf("AESEncryptGCMWithOptions() error = %v", err)
	}
	out, err = AESDecryptGCMWithOptions(cipherText, key, customNonce, []byte("aad"), WithGCMNonceSize(len(customNonce)))
	if err != nil {
		t.Fatalf("AESDecryptGCMWithOptions() error = %v", err)
	}
	if !bytes.Equal(out, plain) {
		t.Fatalf("AESDecryptGCMWithOptions() = %q", out)
	}
}
