package date

import (
	"testing"
	"time"
)

// Tests aligned with hutool-core DateUtilTest.

func TestFormatAndParse(t *testing.T) {
	tt := time.Date(2024, 7, 15, 10, 20, 30, 0, time.Local)
	if got := FormatDateNorm(tt); got != "2024-07-15 10:20:30" {
		t.Fatalf("FormatDateNorm: %q", got)
	}
	if got := FormatDateOnly(tt); got != "2024-07-15" {
		t.Fatalf("FormatDateOnly: %q", got)
	}
	parsed, err := ParseDate("2024-07-15 10:20:30")
	if err != nil {
		t.Fatalf("ParseDate err: %v", err)
	}
	if !parsed.Equal(tt) {
		t.Fatalf("Parsed mismatch: %v", parsed)
	}
	if _, err := ParseDate("2024/07/15"); err != nil {
		t.Fatalf("ParseDate slash: %v", err)
	}
	if _, err := ParseDate("20240715"); err != nil {
		t.Fatalf("ParseDate pure: %v", err)
	}
}

func TestBeginEndOf(t *testing.T) {
	tt := time.Date(2024, 7, 15, 10, 20, 30, 123, time.Local)
	if FormatDateNorm(BeginOfDay(tt)) != "2024-07-15 00:00:00" {
		t.Fatalf("BeginOfDay failed")
	}
	if FormatDateOnly(EndOfDay(tt)) != "2024-07-15" || EndOfDay(tt).Hour() != 23 {
		t.Fatalf("EndOfDay failed")
	}
	if FormatDateNorm(BeginOfMonth(tt)) != "2024-07-01 00:00:00" {
		t.Fatalf("BeginOfMonth failed")
	}
	if FormatDateOnly(EndOfMonth(tt)) != "2024-07-31" {
		t.Fatalf("EndOfMonth failed")
	}
	if FormatDateNorm(BeginOfYear(tt)) != "2024-01-01 00:00:00" {
		t.Fatalf("BeginOfYear failed")
	}
	if FormatDateOnly(EndOfYear(tt)) != "2024-12-31" {
		t.Fatalf("EndOfYear failed")
	}
}

func TestOffsets(t *testing.T) {
	tt := time.Date(2024, 7, 15, 10, 0, 0, 0, time.Local)
	if FormatDateOnly(OffsetDay(tt, 1)) != "2024-07-16" {
		t.Fatalf("OffsetDay failed")
	}
	if FormatDateOnly(OffsetMonth(tt, 1)) != "2024-08-15" {
		t.Fatalf("OffsetMonth failed")
	}
	if FormatDateOnly(OffsetYear(tt, -1)) != "2023-07-15" {
		t.Fatalf("OffsetYear failed")
	}
	if OffsetHour(tt, 2).Hour() != 12 {
		t.Fatalf("OffsetHour failed")
	}
}

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
