package vnum

import (
	"strconv"
	"testing"
)

func TestNumCalculateFacadeWithCustomParser(t *testing.T) {
	calls := 0
	got, err := CalculateWithOptions("1 + 2 * 3", WithParseFloatFunc(func(s string, bitSize int) (float64, error) {
		calls++
		return strconv.ParseFloat(s, bitSize)
	}))
	if err != nil || got != 7 || calls == 0 {
		t.Fatalf("CalculateWithOptions = %f, %v calls=%d", got, err, calls)
	}
	if _, err := Calculate("1 + * 2"); err == nil {
		t.Fatal("Calculate invalid expression error = nil")
	}
}
