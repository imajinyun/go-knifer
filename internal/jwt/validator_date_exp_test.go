package jwt

import (
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
