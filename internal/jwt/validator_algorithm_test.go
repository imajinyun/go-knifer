package jwt

import (
	"testing"
	"time"
)

// TestValidateAlgorithm should pass validation when algorithms match.
func TestValidateAlgorithm(t *testing.T) {
	tok, err := New().SetNotBefore(time.Now()).SetKey([]byte("123456")).Sign()
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	signer := MustHMACSigner(AlgHS256, []byte("123456"))
	if err := ValidateAlgorithm(tok, signer); err != nil {
		t.Fatalf("should pass: %v", err)
	}
}

// TestValidateAlgorithmMismatch should return an error when algorithms mismatch.
func TestValidateAlgorithmMismatch(t *testing.T) {
	tok, err := New().SetKey([]byte("123456")).Sign()
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	signer := MustHMACSigner(AlgHS512, []byte("123456"))
	if err := ValidateAlgorithm(tok, signer); err == nil {
		t.Fatalf("expected algorithm mismatch error")
	}
}

func TestValidateAlgorithmRejectsMissingSigner(t *testing.T) {
	tok, err := New().SetKey([]byte("secret")).SetPayload("sub", "public").Sign()
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if err := ValidateAlgorithm(tok, nil); err == nil {
		t.Fatal("ValidateAlgorithm should reject nil signer")
	}
}
