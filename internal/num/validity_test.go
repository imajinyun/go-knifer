package num

import (
	"math"
	"testing"
)

func TestValidityHelpers(t *testing.T) {
	if IsValid(math.Inf(1)) || IsValidFloat32(float32(math.NaN())) || !IsValidNumber(1) {
		t.Fatal("valid number helpers failed")
	}
	if IsValidNumber(nil) || IsValidNumber(math.NaN()) || IsValidNumber(float32(math.Inf(-1))) || !IsValidNumber("not-a-number") {
		t.Fatal("IsValidNumber cases failed")
	}
	if IsValid(math.NaN()) || IsValid(math.Inf(1)) || !IsValid(1.23) || IsValidFloat32(float32(math.Inf(1))) || !IsValidFloat32(1.23) {
		t.Fatal("valid finite checks failed")
	}
}
