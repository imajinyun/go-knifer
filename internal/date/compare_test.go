package date

import (
	"testing"
	"time"
)

func TestBetweenAndSameDay(t *testing.T) {
	a := time.Date(2024, 7, 15, 0, 0, 0, 0, time.Local)
	b := time.Date(2024, 7, 20, 0, 0, 0, 0, time.Local)
	if BetweenDays(a, b) != 5 {
		t.Fatalf("BetweenDays failed")
	}
	if !IsSameDay(a, time.Date(2024, 7, 15, 23, 0, 0, 0, time.Local)) {
		t.Fatalf("IsSameDay failed")
	}
}
