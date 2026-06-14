package cron

import (
	"testing"
	"time"
)

func TestPatternLastDay(t *testing.T) {
	p := mustPattern(t, "0 0 L * *")
	// 2024-02 is a leap-year month, so the last day is 29.
	yes := time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)
	no := time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC)
	if !p.Match(yes, false) {
		t.Fatalf("expect last day of feb 2024 = 29")
	}
	if p.Match(no, false) {
		t.Fatalf("28 is not last day of feb 2024")
	}
	// Last day of February in a non-leap year.
	yes2 := time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC)
	if !p.Match(yes2, false) {
		t.Fatalf("expect last day of feb 2023 = 28")
	}
}

func TestPatternWithSecondField(t *testing.T) {
	p := mustPattern(t, "30 0 12 * * *")
	yes := time.Date(2024, 1, 1, 12, 0, 30, 0, time.UTC)
	no := time.Date(2024, 1, 1, 12, 0, 31, 0, time.UTC)
	if !p.Match(yes, true) {
		t.Fatalf("expect match at 12:00:30")
	}
	if p.Match(no, true) {
		t.Fatalf("should not match at 12:00:31")
	}
}

func TestPatternOrExpression(t *testing.T) {
	p := mustPattern(t, "0 9 * * mon | 0 17 * * fri")
	mon9 := time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)   // Monday.
	fri17 := time.Date(2024, 1, 5, 17, 0, 0, 0, time.UTC) // Friday.
	tue := time.Date(2024, 1, 2, 9, 0, 0, 0, time.UTC)
	if !p.Match(mon9, false) || !p.Match(fri17, false) {
		t.Fatalf("OR expression should match both")
	}
	if p.Match(tue, false) {
		t.Fatalf("should not match tuesday")
	}
}

func TestPatternYear(t *testing.T) {
	// 7 fields: second minute hour dom month dow year.
	p := mustPattern(t, "0 0 0 1 1 * 2024")
	yes := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	no := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	if !p.Match(yes, false) {
		t.Fatalf("expect match year 2024")
	}
	if p.Match(no, false) {
		t.Fatalf("should not match year 2025")
	}
}
