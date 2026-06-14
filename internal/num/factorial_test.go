package num

import (
	"math/big"
	"testing"
)

func TestFactorial(t *testing.T) {
	if got, err := Factorial(5); err != nil || got != 120 {
		t.Fatalf("Factorial: %d %v", got, err)
	}
	if got, err := FactorialRange(5, 2); err != nil || got != 60 {
		t.Fatalf("FactorialRange: %d %v", got, err)
	}
	if FactorialBig(big.NewInt(20)).String() != "2432902008176640000" {
		t.Fatal("FactorialBig failed")
	}
}

func TestFactorialEdges(t *testing.T) {
	if got, err := Factorial(0); err != nil || got != 1 {
		t.Fatalf("Factorial(0) = %d, %v", got, err)
	}
	if got, err := Factorial(21); err == nil || got != 0 {
		t.Fatalf("Factorial overflow should fail: %d, %v", got, err)
	}
	if got, err := FactorialRange(0, 10); err != nil || got != 1 {
		t.Fatalf("FactorialRange start zero = %d, %v", got, err)
	}
	if got, err := FactorialRange(2, 5); err != nil || got != 0 {
		t.Fatalf("FactorialRange start smaller should be zero: %d, %v", got, err)
	}
	if got, err := FactorialRange(21, 0); err == nil || got != 0 {
		t.Fatalf("FactorialRange overflow should fail: %d, %v", got, err)
	}
	if FactorialBig(nil).String() != "1" || FactorialBig(big.NewInt(-1)).String() != "1" {
		t.Fatal("FactorialBig nil/negative should be one")
	}
	if FactorialBigRange(big.NewInt(5), big.NewInt(2)).String() != "60" {
		t.Fatal("FactorialBigRange normal case failed")
	}
	if FactorialBigRange(nil, big.NewInt(0)).String() != "1" || FactorialBigRange(big.NewInt(2), big.NewInt(5)).String() != "1" {
		t.Fatal("FactorialBigRange guard cases failed")
	}
}
