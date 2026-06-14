package jwt

import (
	"strings"
	"testing"
)

func TestSetKeyRejectsNone(t *testing.T) {
	tok := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJzdWIiOiJwdWJsaWMifQ."
	parsed, err := Of(tok)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if parsed.SetKey([]byte("ignored")).Verify() {
		t.Fatal("SetKey should reject alg=none")
	}
	if err := parsed.SetKeyStrict([]byte("ignored")); err == nil {
		t.Fatal("SetKeyStrict should reject alg=none")
	}
	if Verify(tok, []byte("ignored")) {
		t.Fatal("Verify should reject alg=none")
	}
}

func TestSetKeyEReturnsSignerCreationError(t *testing.T) {
	j := New().SetHeader(HeaderAlgorithm, AlgPS256)
	if err := j.SetKeyE([]byte("hmac-key")); err == nil {
		t.Fatal("SetKeyE should return signer creation error for non-HMAC header alg")
	}
}

func TestSetKeyStrictWithMinLength(t *testing.T) {
	weak := []byte("weak")
	j := New().SetHeader(HeaderAlgorithm, AlgHS256)
	if err := j.SetKeyStrictWithMinLength(weak); err == nil {
		t.Fatal("SetKeyStrictWithMinLength should reject weak HMAC key")
	}

	strong := []byte(strings.Repeat("k", MinHMACKeyBytesHS256))
	if err := j.SetKeyStrictWithMinLength(strong); err != nil {
		t.Fatalf("SetKeyStrictWithMinLength strong key: %v", err)
	}
	if j.Signer() == nil {
		t.Fatal("SetKeyStrictWithMinLength should set signer")
	}
	if j.Algorithm() != AlgHS256 {
		t.Fatalf("algorithm = %q, want HS256", j.Algorithm())
	}
}
