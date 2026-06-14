package jwt

import "testing"

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
