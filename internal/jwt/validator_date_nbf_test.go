package jwt

import (
	"testing"
	"time"
)

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
