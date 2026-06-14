package jwt

import (
	"testing"
	"time"
)

// TestValidateExpired should return false for expired tokens when validating overall validity with leeway=0.
func TestValidateExpired(t *testing.T) {
	// Same as validateTest in the utility toolkit tests.
	token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9." +
		"eyJpc3MiOiJNb0xpIiwiZXhwIjoxNjI0OTU4MDk0NTI4LCJpYXQiOjE2MjQ5NTgwMzQ1MjAsInVzZXIiOiJ1c2VyIn0." +
		"L0uB38p9sZrivbmP0VlDe--j_11YUXTu3TfHhfQhRKc"
	j, err := Of(token)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	// Note that the exp field in this utility toolkit token appears to be milliseconds (1624958094528).
	// validate(0) returns false in the utility toolkit, so it should fail here too.
	if j.SetKey([]byte("1234567890")).Validate(0) {
		t.Fatalf("expected validate=false")
	}
}

func TestValidateWithOptions(t *testing.T) {
	now := time.Now()
	tok, err := New().
		SetPayload("sub", "options").
		SetIssuedAt(now).
		SetNotBefore(now.Add(3 * time.Second)).
		SetExpiresAt(now.Add(3 * time.Second)).
		SetKey([]byte("123456")).
		Sign()
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	j, err := Of(tok)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	j.SetKey([]byte("123456"))
	if j.ValidateWithOptions(WithValidateTime(now.Add(-4 * time.Second))) {
		t.Fatal("ValidateWithOptions should reject token before nbf without leeway")
	}
	if !j.ValidateWithOptions(
		WithValidateClock(func() time.Time { return now.Add(-4 * time.Second) }),
		WithValidateLeeway(10),
	) {
		t.Fatal("ValidateWithOptions should accept token within configured leeway")
	}
}
