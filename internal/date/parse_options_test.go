package date

import (
	"errors"
	"testing"
	"time"
)

func TestParseDateWithOptionsLocation(t *testing.T) {
	loc := time.FixedZone("biz", 8*60*60)
	parsed, err := ParseDateWithOptions("2024-07-15 10:20:30", WithLocation(loc))
	if err != nil {
		t.Fatalf("ParseDateWithOptions err: %v", err)
	}
	if parsed.Location() != loc || parsed.Format(NormPattern) != "2024-07-15 10:20:30" {
		t.Fatalf("ParseDateWithOptions location = %v, %s", parsed.Location(), parsed.Format(NormPattern))
	}

	parsed, err = ParseDateLayoutWithOptions("2024/07/15 10:20:30", "2006/01/02 15:04:05", WithLocation(loc))
	if err != nil {
		t.Fatalf("ParseDateLayoutWithOptions err: %v", err)
	}
	if parsed.Location() != loc || parsed.Format(NormPattern) != "2024-07-15 10:20:30" {
		t.Fatalf("ParseDateLayoutWithOptions location = %v, %s", parsed.Location(), parsed.Format(NormPattern))
	}
}

func TestParseDateWithOptionsParserProvider(t *testing.T) {
	loc := time.FixedZone("parser", 8*60*60)
	called := false
	parsed, err := ParseDateWithOptions("custom", WithLocation(loc), WithParseInLocationFunc(func(layout, value string, location *time.Location) (time.Time, error) {
		called = true
		if value == "custom" && location == loc {
			return time.Date(2026, 6, 7, 1, 2, 3, 0, location), nil
		}
		return time.Time{}, errors.New("unsupported")
	}))
	if err != nil {
		t.Fatalf("ParseDateWithOptions custom parser: %v", err)
	}
	if !called || parsed.Location() != loc || parsed.Format(NormPattern) != "2026-06-07 01:02:03" {
		t.Fatalf("custom parser called=%v parsed=%v", called, parsed)
	}

	parsed, err = ParseDateLayoutWithOptions("layout", "custom-layout", WithParseInLocationFunc(func(layout, value string, location *time.Location) (time.Time, error) {
		if layout != "custom-layout" || value != "layout" {
			return time.Time{}, errors.New("unexpected input")
		}
		return time.Date(2026, 1, 2, 0, 0, 0, 0, location), nil
	}))
	if err != nil || parsed.Format(NormDatePattern) != "2026-01-02" {
		t.Fatalf("ParseDateLayoutWithOptions custom parser = %v, %v", parsed, err)
	}
}
