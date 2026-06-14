package jwt

import (
	"math"
	"testing"
	"time"
)

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
