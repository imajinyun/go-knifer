package vnum

import (
	"math/big"
	"testing"
)

func TestNumFactorialAndMathFacades(t *testing.T) {
	if got, err := Factorial(5); err != nil || got != 120 {
		t.Fatalf("Factorial = %d, %v", got, err)
	}
	if _, err := Factorial(21); err == nil {
		t.Fatal("Factorial overflow error = nil")
	}
	if got, err := FactorialRange(5, 2); err != nil || got != 60 {
		t.Fatalf("FactorialRange = %d, %v", got, err)
	}
	if got := FactorialBig(big.NewInt(5)); got.String() != "120" {
		t.Fatalf("FactorialBig = %s", got)
	}
	if got := FactorialBigRange(big.NewInt(5), big.NewInt(2)); got.String() != "60" {
		t.Fatalf("FactorialBigRange = %s", got)
	}
	if Sqrt(16) != 4 || ProcessMultiple(5, 2) != 10 || Divisor(24, 18) != 6 || Multiple(4, 6) != 12 {
		t.Fatal("factorial helpers failed")
	}
}
