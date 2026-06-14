package cron

import (
	"testing"
	"time"
)

func TestPatternBasic(t *testing.T) {
	// 5 fields: every minute.
	p := mustPattern(t, "* * * * *")
	now := time.Date(2024, 1, 1, 12, 30, 45, 0, time.UTC)
	if !p.Match(now, false) {
		t.Fatalf("'* * * * *' should match any minute")
	}
}

func TestPatternMinute(t *testing.T) {
	p := mustPattern(t, "30 12 * * *")
	yes := time.Date(2024, 1, 2, 12, 30, 0, 0, time.UTC)
	no := time.Date(2024, 1, 2, 13, 30, 0, 0, time.UTC)
	if !p.Match(yes, false) {
		t.Fatalf("expect match at 12:30")
	}
	if p.Match(no, false) {
		t.Fatalf("should not match at 13:30")
	}
}
