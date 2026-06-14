package cron

import (
	"testing"
	"time"
)

func TestPatternStepRange(t *testing.T) {
	// Every 15 minutes.
	p := mustPattern(t, "*/15 * * * *")
	for _, m := range []int{0, 15, 30, 45} {
		ts := time.Date(2024, 1, 1, 0, m, 0, 0, time.UTC)
		if !p.Match(ts, false) {
			t.Fatalf("expect match at minute %d", m)
		}
	}
	for _, m := range []int{1, 7, 14, 16, 31, 44, 46, 59} {
		ts := time.Date(2024, 1, 1, 0, m, 0, 0, time.UTC)
		if p.Match(ts, false) {
			t.Fatalf("should not match at minute %d", m)
		}
	}
}

func TestPatternDayOfWeekAlias(t *testing.T) {
	// Every Monday at 9:00.
	p := mustPattern(t, "0 9 * * mon")
	mon := time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC) // 2024-01-01 is Monday.
	tue := time.Date(2024, 1, 2, 9, 0, 0, 0, time.UTC)
	if !p.Match(mon, false) {
		t.Fatalf("expect match on monday")
	}
	if p.Match(tue, false) {
		t.Fatalf("should not match on tuesday")
	}
}

func TestPatternMonthAlias(t *testing.T) {
	p := mustPattern(t, "0 0 1 jan *")
	yes := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	no := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	if !p.Match(yes, false) || p.Match(no, false) {
		t.Fatalf("month alias mismatch")
	}
}

func TestPatternList(t *testing.T) {
	p := mustPattern(t, "0,15,30,45 * * * *")
	for _, m := range []int{0, 15, 30, 45} {
		ts := time.Date(2024, 1, 1, 0, m, 0, 0, time.UTC)
		if !p.Match(ts, false) {
			t.Fatalf("expect match minute %d", m)
		}
	}
}

func TestPatternRange(t *testing.T) {
	p := mustPattern(t, "0 9-17 * * *")
	for h := 9; h <= 17; h++ {
		ts := time.Date(2024, 1, 1, h, 0, 0, 0, time.UTC)
		if !p.Match(ts, false) {
			t.Fatalf("expect match hour %d", h)
		}
	}
	for _, h := range []int{0, 8, 18, 23} {
		ts := time.Date(2024, 1, 1, h, 0, 0, 0, time.UTC)
		if p.Match(ts, false) {
			t.Fatalf("should not match hour %d", h)
		}
	}
}
