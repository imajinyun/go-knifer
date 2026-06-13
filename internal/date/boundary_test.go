package date

import (
	"testing"
	"time"
)

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
