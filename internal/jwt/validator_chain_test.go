package jwt

import (
	"testing"
	"time"
)

func TestJWTValidator_Chain(t *testing.T) {
	signer := HS256([]byte("secret"))
	now := time.Now()
	token, err := New().
		AddPayloads(map[string]any{
			PayloadIssuer:    "alice",
			PayloadIssuedAt:  now.Unix(),
			PayloadExpiresAt: now.Add(time.Hour).Unix(),
		}).
		SetSigner(signer).
		Sign()
	if err != nil {
		t.Fatal(err)
	}

	if err := OfValidator(token).
		ValidateAlgorithm(signer).
		ValidateDate(now, 0).
		Err(); err != nil {
		t.Fatalf("validator should pass: %v", err)
	}

	// Algorithm mismatch.
	if err := OfValidator(token).ValidateAlgorithm(HS384([]byte("secret"))).Err(); err == nil {
		t.Fatal("expected algorithm mismatch error")
	}

	// Expiration scenario.
	expired, _ := New().
		AddPayloads(map[string]any{PayloadExpiresAt: now.Add(-time.Hour).Unix()}).
		SetSigner(signer).
		Sign()
	if err := OfValidator(expired).
		ValidateAlgorithm(signer).
		ValidateDate(now, 0).
		Err(); err == nil {
		t.Fatal("expected expired error")
	}
}
