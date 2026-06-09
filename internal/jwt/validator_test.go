package jwt

import (
	"math"
	"testing"
	"time"
)

// Matches the utility toolkit-jwt JWTValidatorTest.

// TestExpiredAt should return a validation error for an expired token.
func TestExpiredAt(t *testing.T) {
	// Use the same token as the utility toolkit test; exp=1477592 is expired.
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0Nzc1OTJ9.isvT0Pqx0yjnZk53mUFSeYFJLDs-Ls9IsNAm86gIdZo"
	j, err := Of(token)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if err := ValidateDate(j, time.Now(), 0); err == nil {
		t.Fatalf("expected expired error")
	}
}

// TestIssueAt should fail when issued-at time is after the reference time.
func TestIssueAt(t *testing.T) {
	now := time.Now()
	tok, err := New().SetIssuedAt(now).SetKey([]byte("123456")).Sign()
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	j, _ := Of(tok)
	yesterday := now.AddDate(0, 0, -1)
	if err := ValidateDate(j, yesterday, 0); err == nil {
		t.Fatalf("expected error: iat in future of yesterday")
	}
}

// TestIssueAtPass should pass when issued-at time is not after the reference time.
func TestIssueAtPass(t *testing.T) {
	now := time.Now()
	tok, err := New().SetIssuedAt(now).SetKey([]byte("123456")).Sign()
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	j, _ := Of(tok)
	if err := ValidateDate(j, now, 0); err != nil {
		t.Fatalf("should pass: %v", err)
	}
}

// TestNotBefore should fail when nbf is after the reference time.
func TestNotBefore(t *testing.T) {
	now := time.Now()
	j := New().SetNotBefore(now)
	yesterday := now.AddDate(0, 0, -1)
	if err := ValidateDate(j, yesterday, 0); err == nil {
		t.Fatalf("expected error: nbf later than now")
	}
}

// TestNotBeforePass should pass when nbf is not after the reference time.
func TestNotBeforePass(t *testing.T) {
	now := time.Now()
	j := New().SetNotBefore(now)
	if err := ValidateDate(j, now, 0); err != nil {
		t.Fatalf("should pass: %v", err)
	}
}

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

// TestValidateDateExpired should reject a directly constructed expired JWT.
func TestValidateDateExpired(t *testing.T) {
	exp, _ := time.Parse("2006-01-02 15:04:05", "2021-10-13 09:59:00")
	j := New().
		SetPayload("id", 123).
		SetPayload("username", "the utility toolkit").
		SetExpiresAt(exp)
	if err := ValidateDate(j, time.Now(), 0); err == nil {
		t.Fatalf("expected expired error")
	}
}

// TestValidateLeeway should pass within the leeway tolerance window.
func TestValidateLeeway(t *testing.T) {
	now := time.Now()
	expired := now.Add(3 * time.Second)
	tok, err := New().
		SetPayload("sub", "blue-light").
		SetIssuedAt(now).
		SetNotBefore(expired).
		SetExpiresAt(expired).
		SetKey([]byte("123456")).
		Sign()
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	j, _ := Of(tok)
	// 4 seconds ago: exceeds the nbf boundary by 1 second, but leeway=10 permits it.
	before := now.Add(-4 * time.Second)
	if err := ValidateDate(j, before, 10); err != nil {
		t.Fatalf("should pass with leeway: %v", err)
	}
	// 4 seconds later: exceeds exp by 1 second, but leeway=10 still permits it.
	after := now.Add(4 * time.Second)
	if err := ValidateDate(j, after, 10); err != nil {
		t.Fatalf("should pass with leeway: %v", err)
	}
}

func TestValidateDateRejectsMalformedTimeClaims(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name  string
		claim string
		value any
	}{
		{name: "nbf string", claim: PayloadNotBefore, value: "not-a-number"},
		{name: "exp fractional", claim: PayloadExpiresAt, value: float64(123.45)},
		{name: "iat unsupported", claim: PayloadIssuedAt, value: []string{"bad"}},
		{name: "exp infinite", claim: PayloadExpiresAt, value: math.Inf(1)},
		{name: "iat overflow", claim: PayloadIssuedAt, value: uint64(math.MaxInt64) + 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := New().SetPayload(tt.claim, tt.value)
			if err := ValidateDate(j, now, 0); err == nil {
				t.Fatal("ValidateDate should reject malformed time claim")
			}
		})
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
