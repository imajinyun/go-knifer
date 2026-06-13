package date

import (
	"errors"
	"testing"

	knifer "github.com/imajinyun/go-knifer"
)

func TestDateErrorContract(t *testing.T) {
	_, err := ParseDate("")
	assertDateCode(t, err, knifer.ErrCodeInvalidInput)

	_, err = ParseDate("not-a-date")
	assertDateCode(t, err, knifer.ErrCodeInvalidInput)

	_, err = ParseDateLayout("2026-06-05", "bad-layout")
	assertDateCode(t, err, knifer.ErrCodeInvalidInput)
}

func assertDateCode(t *testing.T, err error, code knifer.ErrCode) {
	t.Helper()
	if err == nil {
		t.Fatalf("err = nil, want %s", code)
	}
	if !errors.Is(err, code) {
		t.Fatalf("errors.Is(%v, %s) = false", err, code)
	}
	got, ok := knifer.CodeOf(err)
	if !ok || got != code {
		t.Fatalf("CodeOf(%v) = %q, %v; want %q, true", err, got, ok, code)
	}
	var dateErr *DateError
	if !errors.As(err, &dateErr) {
		t.Fatalf("errors.As(err, *DateError) = false: %v", err)
	}
}
