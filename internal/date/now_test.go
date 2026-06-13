package date

import (
	"testing"
	"time"
)

func TestNowWithOptionsClock(t *testing.T) {
	fixed := time.Date(2026, 6, 6, 12, 34, 56, 0, time.FixedZone("fixed", 8*60*60))
	if got := NowWithOptions(WithClock(func() time.Time { return fixed })); !got.Equal(fixed) {
		t.Fatalf("NowWithOptions = %v, want %v", got, fixed)
	}
	today := TodayWithOptions(WithClock(func() time.Time { return fixed }))
	if !today.Equal(time.Date(2026, 6, 6, 0, 0, 0, 0, fixed.Location())) {
		t.Fatalf("TodayWithOptions = %v", today)
	}
	if got := NowWithOptions(WithClock(nil)); got.IsZero() {
		t.Fatal("NowWithOptions nil clock should fall back to time.Now")
	}
}
