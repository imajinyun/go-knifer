package num

import (
	"math"
	"math/big"
	"testing"
)

func TestComparisonEqualityEdges(t *testing.T) {
	if Compare(2, 1) != 1 || Compare(1, 1) != 0 || Compare("a", "b") != -1 {
		t.Fatal("Compare cases failed")
	}
	if !IsGreaterOrEqual(2, 2) || !IsLess(1, 2) || !IsLessOrEqual(2, 2) || IsIn(0, 1, 3) {
		t.Fatal("ordered helper cases failed")
	}
	if !Equals(0.1+0.2, 0.3) || Equals(0.1, 0.2) {
		t.Fatal("Equals tolerance cases failed")
	}
	if EqualsExact(0.0, math.Copysign(0, -1)) || !EqualsFloat32Exact(float32(1), float32(1)) || !EqualsInt64(9, 9) {
		t.Fatal("exact equality cases failed")
	}
	if !EqualsBigDecimal(big.NewRat(10, 10), big.NewRat(1, 1)) || EqualsBigDecimal(big.NewRat(1, 1), nil) || !EqualsBigDecimal(nil, nil) {
		t.Fatal("big decimal equality cases failed")
	}
	if !EqualsChar('ß', 'ß', false) || EqualsChar('A', 'a', false) || !EqualsChar('A', 'a', true) {
		t.Fatal("char equality cases failed")
	}
}
