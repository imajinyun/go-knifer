package mail

import (
	"strings"
	"testing"
)

func sequenceBoundary(values ...string) BoundaryGenerator {
	idx := 0
	return func() (string, error) {
		if idx >= len(values) {
			return "extra-boundary", nil
		}
		value := values[idx]
		idx++
		return value, nil
	}
}

func assertContains(t *testing.T, got, expected string) {
	t.Helper()
	if !strings.Contains(got, expected) {
		t.Fatalf("expected %q to contain %q", got, expected)
	}
}
