package jwt

import (
	"strings"
	"testing"
	"time"
)

// Matches the utility toolkit-jwt JWTTest.

func TestCreateHS256(t *testing.T) {
	key := []byte("1234567890")
	j := New().
		SetPayload("sub", "1234567890").
		SetPayload("name", "looly").
		SetPayload("admin", true).
		SetExpiresAt(time.Unix(1640966400, 0)).
		SetKey(key)

	tok, err := j.Sign()
	if err != nil {
		t.Fatalf("sign err: %v", err)
	}
	parts := strings.Split(tok, ".")
	if len(parts) != 3 {
		t.Fatalf("token parts: %d", len(parts))
	}
	// It is enough that the parsed token verifies successfully.
	parsed, err := Of(tok)
	if err != nil {
		t.Fatalf("parse err: %v", err)
	}
	if !parsed.SetKey(key).Verify() {
		t.Fatalf("verify failed")
	}
	if parsed.Payload("name") != "looly" {
		t.Fatalf("payload name: %v", parsed.Payload("name"))
	}
	if parsed.Algorithm() != AlgHS256 {
		t.Fatalf("alg: %s", parsed.Algorithm())
	}
	if parsed.Type() != "JWT" {
		t.Fatalf("typ: %s", parsed.Type())
	}
}

func TestParseAndVerifyKnownToken(t *testing.T) {
	// Fixed test token from the utility toolkit.
	rightToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9." +
		"eyJzdWIiOiIxMjM0NTY3ODkwIiwiYWRtaW4iOnRydWUsIm5hbWUiOiJsb29seSJ9." +
		"U2aQkC2THYV9L0fTN-yBBI7gmo5xhmvMhATtu8v0zEA"

	j, err := Of(rightToken)
	if err != nil {
		t.Fatalf("parse err: %v", err)
	}
	if !j.SetKey([]byte("1234567890")).Verify() {
		t.Fatalf("verify failed")
	}
	if j.Header(HeaderType) != "JWT" {
		t.Fatalf("type: %v", j.Header(HeaderType))
	}
	if j.Header(HeaderAlgorithm) != "HS256" {
		t.Fatalf("alg: %v", j.Header(HeaderAlgorithm))
	}
	if j.Header(HeaderContentType) != nil {
		t.Fatalf("cty should be nil")
	}
	if j.Payload("sub") != "1234567890" {
		t.Fatalf("sub: %v", j.Payload("sub"))
	}
	if j.Payload("name") != "looly" {
		t.Fatalf("name: %v", j.Payload("name"))
	}
	if j.Payload("admin") != true {
		t.Fatalf("admin: %v", j.Payload("admin"))
	}
}
