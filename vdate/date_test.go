package vdate

import (
	"testing"
	"time"
)

func TestDateFacade(t *testing.T) {
	base := time.Date(2026, 5, 30, 1, 2, 3, 0, time.Local)
	if Format(base, "") != "2026-05-30 01:02:03" || FormatNorm(base) != "2026-05-30 01:02:03" {
		t.Fatal("format failed")
	}
	if FormatDateOnly(base) != "2026-05-30" || FormatTimeOnly(base) != "01:02:03" {
		t.Fatal("date/time format failed")
	}
	if got, err := Parse("2026-05-30"); err != nil || got.Year() != 2026 {
		t.Fatalf("Parse = %v, %v", got, err)
	}
	if got, err := ParseLayout("2026/05/30", "2006/01/02"); err != nil || got.Day() != 30 {
		t.Fatalf("ParseLayout = %v, %v", got, err)
	}
	if BeginOfDay(base).Hour() != 0 || EndOfDay(base).Hour() != 23 {
		t.Fatal("begin/end day failed")
	}
	if BeginOfMonth(base).Day() != 1 || EndOfMonth(base).Day() != 31 {
		t.Fatal("begin/end month failed")
	}
	if BeginOfYear(base).Month() != time.January || EndOfYear(base).Month() != time.December {
		t.Fatal("begin/end year failed")
	}
	if OffsetDay(base, 1).Day() != 31 || OffsetMonth(base, 1).Month() != time.June || OffsetYear(base, 1).Year() != 2027 {
		t.Fatal("date offset failed")
	}
	if OffsetHour(base, 1).Hour() != 2 || OffsetMinute(base, 1).Minute() != 3 || OffsetSecond(base, 1).Second() != 4 {
		t.Fatal("time offset failed")
	}
	if BetweenDays(base, base.Add(48*time.Hour)) != 2 || !IsSameDay(base, base.Add(time.Hour)) {
		t.Fatal("comparison failed")
	}
}
