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

func TestSetSignerNilIsSafe(t *testing.T) {
	j := New()
	j.SetSigner(nil)
	if j.Signer() != nil {
		t.Fatal("SetSigner(nil) should clear signer")
	}
	if got := j.Algorithm(); got != "" {
		t.Fatalf("SetSigner(nil) should not write alg, got %q", got)
	}
	if err := j.SetSignerE(nil); err == nil {
		t.Fatal("SetSignerE(nil) should return an error")
	}
}

func TestNeedSigner(t *testing.T) {
	j := New().SetPayload("sub", "x")
	if _, err := j.Sign(); err == nil {
		t.Fatalf("expected error when no signer set")
	}
}

func TestSignRejectsEmptySignature(t *testing.T) {
	if tok, err := New().SetPayload("sub", "x").SetSigner(emptySigner{alg: AlgPS256}).Sign(); err == nil || tok != "" {
		t.Fatalf("Sign should reject empty signature, token=%q err=%v", tok, err)
	}
}

type emptySigner struct{ alg string }

func (s emptySigner) Algorithm() string { return s.alg }

func (emptySigner) Sign(string, string) string { return "" }

func (emptySigner) Verify(string, string, string) bool { return false }

func TestVerifyMismatchKey(t *testing.T) {
	tok, err := New().SetPayload("a", 1).SetKey([]byte("right")).Sign()
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	j, _ := Of(tok)
	if j.SetKey([]byte("wrong")).Verify() {
		t.Fatalf("should fail with wrong key")
	}
}

func TestVerifyWithRejectsAlgorithmMismatch(t *testing.T) {
	tok, _ := New().SetKey([]byte("k")).SetPayload("a", 1).Sign()
	j, _ := Of(tok)
	hs, _ := NewHMACSigner(AlgHS512, []byte("k"))
	if j.VerifyWith(hs) {
		t.Fatalf("HS256 token with HS512 signer should fail")
	}
}
