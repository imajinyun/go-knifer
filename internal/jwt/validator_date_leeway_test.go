package jwt

import (
	"testing"
	"time"
)

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
