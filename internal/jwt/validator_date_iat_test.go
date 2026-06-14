package jwt

import (
	"testing"
	"time"
)

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
