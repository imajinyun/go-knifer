package jwt

import (
	"math"
	"testing"
	"time"
)

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
