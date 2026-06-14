package vnum

import (
	"math/big"
	"testing"
)

func TestNumCompareAndEqualityFacades(t *testing.T) {
	if !EqualsExact(1, 1) || !EqualsFloat32Exact(1, 1) || !EqualsInt64(1, 1) || !EqualsBigDecimal(big.NewRat(1, 2), big.NewRat(2, 4)) || !EqualsChar('A', 'a', true) {
		t.Fatal("equals helpers failed")
	}
	if Compare(1, 2) != -1 || !IsGreater(2, 1) || !IsGreaterOrEqual(2, 2) || !IsLess(1, 2) || !IsLessOrEqual(2, 2) || !IsIn(2, 1, 3) {
		t.Fatal("comparison helpers failed")
	}
}
