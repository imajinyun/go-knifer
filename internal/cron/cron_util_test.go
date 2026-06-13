package cron

import (
	"errors"
	"testing"
)

func TestSetMatchSecondEWithOptionsReturnsStartedError(t *testing.T) {
	s := NewSchedulerWithOptions(WithMatchSecond(true))
	if err := s.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer s.Stop()

	err := SetMatchSecondEWithOptions(false, WithDefaultScheduler(s))
	if !errors.Is(err, ErrSchedulerStarted) {
		t.Fatalf("SetMatchSecondEWithOptions while started = %v, want ErrSchedulerStarted", err)
	}
	if !s.IsMatchSecond() {
		t.Fatal("started scheduler config should not be mutated")
	}
}

func TestDefaultSchedulerOperationOptions(t *testing.T) {
	global := ConfigureDefaultScheduler(WithIDGenerator(func() string { return "global-id" }))
	t.Cleanup(func() { ConfigureDefaultScheduler() })
	isolated := NewSchedulerWithOptions(WithIDGenerator(func() string { return "isolated-id" }))

	id, err := ScheduleFuncWithOptions("* * * * *", func() {}, WithDefaultScheduler(isolated))
	if err != nil {
		t.Fatalf("ScheduleFuncWithOptions: %v", err)
	}
	if id != "isolated-id" || isolated.Size() != 1 || global.Size() != 0 {
		t.Fatalf("default scheduler option not isolated: id=%q isolated=%d global=%d", id, isolated.Size(), global.Size())
	}
	if err := UpdatePatternWithOptions(id, "0 0 * * *", WithDefaultScheduler(isolated)); err != nil {
		t.Fatalf("UpdatePatternWithOptions: %v", err)
	}
	if !RemoveWithOptions(id, WithDefaultScheduler(isolated)) || isolated.Size() != 0 {
		t.Fatalf("RemoveWithOptions did not remove isolated task")
	}
}

func TestWithDefaultSchedulerOptionsCreatesPerCallScheduler(t *testing.T) {
	global := ConfigureDefaultScheduler(WithIDGenerator(func() string { return "global-id" }))
	t.Cleanup(func() { ConfigureDefaultScheduler() })

	id, err := ScheduleFuncWithOptions("* * * * *", func() {}, WithDefaultSchedulerOptions(WithIDGenerator(func() string { return "per-call-id" })))
	if err != nil {
		t.Fatalf("ScheduleFuncWithOptions: %v", err)
	}
	if id != "per-call-id" || global.Size() != 0 {
		t.Fatalf("per-call scheduler option leaked to global: id=%q globalSize=%d", id, global.Size())
	}
}
