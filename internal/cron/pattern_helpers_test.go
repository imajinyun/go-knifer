package cron

import "testing"

func mustPattern(t *testing.T, expr string) *Pattern {
	t.Helper()
	p, err := NewPattern(expr)
	if err != nil {
		t.Fatalf("parse %q: %v", expr, err)
	}
	return p
}
