package cron

import (
	"testing"
	"time"
)

func mustPattern(t *testing.T, expr string) *Pattern {
	t.Helper()
	p, err := NewPattern(expr)
	if err != nil {
		t.Fatalf("parse %q: %v", expr, err)
	}
	return p
}

func TestPatternBasic(t *testing.T) {
	// 5 字段：每分钟
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

func TestPatternStepRange(t *testing.T) {
	// 每 15 分钟
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
	// 每周一 9:00
	p := mustPattern(t, "0 9 * * mon")
	mon := time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC) // 2024-01-01 是周一
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

func TestPatternLastDay(t *testing.T) {
	p := mustPattern(t, "0 0 L * *")
	// 2024-02 是闰年，最后一天为 29
	yes := time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)
	no := time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC)
	if !p.Match(yes, false) {
		t.Fatalf("expect last day of feb 2024 = 29")
	}
	if p.Match(no, false) {
		t.Fatalf("28 is not last day of feb 2024")
	}
	// 非闰年 2 月最后一天
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
	mon9 := time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)   // 周一
	fri17 := time.Date(2024, 1, 5, 17, 0, 0, 0, time.UTC) // 周五
	tue := time.Date(2024, 1, 2, 9, 0, 0, 0, time.UTC)
	if !p.Match(mon9, false) || !p.Match(fri17, false) {
		t.Fatalf("OR expression should match both")
	}
	if p.Match(tue, false) {
		t.Fatalf("should not match tuesday")
	}
}

func TestPatternYear(t *testing.T) {
	// 7 字段：second minute hour dom month dow year
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

func TestPatternInvalid(t *testing.T) {
	cases := []string{
		"",
		"* *",
		"60 * * * *",
		"* * * 13 *",
		"* * * * 8",
	}
	for _, c := range cases {
		if _, err := NewPattern(c); err == nil {
			t.Fatalf("expected error for %q", c)
		}
	}
}
