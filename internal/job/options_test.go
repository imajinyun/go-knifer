package job

import (
	"strconv"
	"strings"
	"testing"
)

func TestOptionsAndHelpers_BitsUT(t *testing.T) {
	if got := (Options{}).normalized(7); got.BatchSize != 7 || got.MaxConcurrency != 1 {
		t.Fatalf("Options{}.normalized(7) = %+v, want batch 7 concurrency 1", got)
	}
	if got := (Options{BatchSize: 3, MaxConcurrency: 2}).normalized(7); got.BatchSize != 3 || got.MaxConcurrency != 2 {
		t.Fatalf("Options.normalized(7) = %+v, want batch 3 concurrency 2", got)
	}
	if got := chunks(7, 3); got != 3 {
		t.Fatalf("chunks(7, 3) = %d, want 3", got)
	}
	if got := chunks(7, 0); got != 0 {
		t.Fatalf("chunks(7, 0) = %d, want 0", got)
	}
}

func formatRange(start, end int) string {
	return strings.Join([]string{itoa(start), itoa(end)}, ":")
}

func itoa(v int) string {
	return strconv.Itoa(v)
}
