package cron

import (
	"bytes"
	"testing"
)

func TestSchedulerIDRandomReaderOption(t *testing.T) {
	s := NewSchedulerWithOptions(WithIDRandomReader(bytes.NewReader([]byte{0, 1, 2, 3, 4, 5, 6, 7})))
	id, err := s.ScheduleFunc("* * * * *", func() {})
	if err != nil {
		t.Fatalf("ScheduleFunc: %v", err)
	}
	if id != "0001020304050607" {
		t.Fatalf("id = %q, want 0001020304050607", id)
	}
}
