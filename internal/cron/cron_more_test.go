package cron

import (
	"context"
	"errors"
	"testing"
	"time"

	knifer "github.com/imajinyun/go-knifer"
)

func TestCronErrorMessage(t *testing.T) {
	// without cause
	err := NewCronError("test error %d", 1)
	if got := err.Error(); got != "test error 1" {
		t.Fatalf("Error() = %q, want %q", got, "test error 1")
	}

	// with cause
	err2 := WrapCronError(errors.New("root cause"), "wrapped")
	if got := err2.Error(); got != "wrapped: root cause" {
		t.Fatalf("Error() with cause = %q", got)
	}
}

func TestCronErrorCode(t *testing.T) {
	err := NewCronError("invalid")
	if got := err.ErrorCode(); got != knifer.ErrCodeInvalidInput {
		t.Fatalf("ErrorCode = %v, want ErrCodeInvalidInput", got)
	}

	err2 := newSchedulerStartedError()
	if got := err2.ErrorCode(); got != knifer.ErrCodeUnsupported {
		t.Fatalf("ErrorCode started = %v, want ErrCodeUnsupported", got)
	}
}

func TestSimpleTaskListenerNoOps(t *testing.T) {
	l := SimpleTaskListener{}
	// These should not panic
	l.OnStart(nil)
	l.OnSucceeded(nil)
	l.OnFailed(nil, nil)
}

func TestCronTaskGettersAndSetters(t *testing.T) {
	p, err := NewPattern("* * * * *")
	if err != nil {
		t.Fatal(err)
	}
	task := TaskFunc(func() {})
	ct := NewCronTask("test-id", p, task)

	if got := ct.ID(); got != "test-id" {
		t.Fatalf("ID = %q, want %q", got, "test-id")
	}
	if got := ct.Pattern(); got != p {
		t.Fatal("Pattern mismatch")
	}
	if got := ct.Raw(); got == nil {
		t.Fatal("Raw task should not be nil")
	}

	p2, err := NewPattern("0 0 * * *")
	if err != nil {
		t.Fatal(err)
	}
	ct.SetPattern(p2)
	if got := ct.Pattern(); got != p2 {
		t.Fatal("SetPattern did not update")
	}
}

func TestCronTaskExecuteNilRaw(t *testing.T) {
	ct := &CronTask{}
	ct.Execute() // should not panic
}

func TestTaskExecutorGetters(t *testing.T) {
	s := NewScheduler()
	task := TaskFunc(func() {})
	ct := NewCronTask("exec-id", MustNewPattern("* * * * *"), task)

	e := &TaskExecutor{scheduler: s, task: ct}
	if got := e.CronTask(); got != ct {
		t.Fatal("CronTask getter mismatch")
	}
	if got := e.Task(); got == nil {
		t.Fatal("Task getter should not return nil")
	}
}

func TestBoolArrayMatcherMinMax(t *testing.T) {
	m := newBoolArrayMatcher([]int{2, 5, 8})
	if got := m.MinValue(); got != 2 {
		t.Fatalf("MinValue = %d, want 2", got)
	}
	if got := m.MaxValue(); got != 8 {
		t.Fatalf("MaxValue = %d, want 8", got)
	}

	// Match
	if !m.Match(5) {
		t.Fatal("Match(5) should be true")
	}
	if m.Match(3) {
		t.Fatal("Match(3) should be false")
	}

	// NextAfter
	if got := m.NextAfter(6); got != 8 {
		t.Fatalf("NextAfter(6) = %d, want 8", got)
	}
	if got := m.NextAfter(9); got != 2 {
		t.Fatalf("NextAfter(9) should wrap to %d, got %d", 2, got)
	}
}

func TestBoolArrayMatcherEmpty(t *testing.T) {
	m := newBoolArrayMatcher(nil)
	if m.Match(0) {
		t.Fatal("empty matcher should not match")
	}
	if got := m.MinValue(); got != 0 {
		t.Fatalf("empty MinValue = %d, want 0", got)
	}
}

func TestAlwaysTrueMatcher(t *testing.T) {
	if !AlwaysTrueMatcher.Match(0) {
		t.Fatal("AlwaysTrueMatcher should match 0")
	}
	if !AlwaysTrueMatcher.Match(100) {
		t.Fatal("AlwaysTrueMatcher should match 100")
	}
	if got := AlwaysTrueMatcher.NextAfter(42); got != 42 {
		t.Fatalf("NextAfter(42) = %d, want 42", got)
	}
}

func TestDefaultScheduler(t *testing.T) {
	s := DefaultScheduler()
	if s == nil {
		t.Fatal("DefaultScheduler returned nil")
	}
}

func TestDefaultSchedulerWithOptions(t *testing.T) {
	s := DefaultSchedulerWithOptions()
	if s == nil {
		t.Fatal("DefaultSchedulerWithOptions returned nil")
	}

	// With per-call option
	s2 := DefaultSchedulerWithOptions(WithDefaultSchedulerOptions())
	if s2 == nil {
		t.Fatal("isolated DefaultSchedulerWithOptions returned nil")
	}
}

func TestSetMatchSecondOnNewScheduler(t *testing.T) {
	s := NewScheduler()
	s.SetMatchSecond(true)
	if !s.IsMatchSecond() {
		t.Fatal("SetMatchSecond(true) failed")
	}
	SetMatchSecondWithOptions(false, WithDefaultScheduler(s))
	if s.IsMatchSecond() {
		t.Fatal("SetMatchSecondWithOptions(false) failed")
	}
}

func TestScheduleAndRunOnNewScheduler(t *testing.T) {
	s := NewScheduler()
	s.SetMatchSecond(true)
	done := make(chan struct{})
	id, err := ScheduleWithOptions("* * * * * *", TaskFunc(func() { close(done) }), WithDefaultScheduler(s))
	if err != nil {
		t.Fatal(err)
	}
	if id == "" {
		t.Fatal("ScheduleWithOptions returned empty ID")
	}
	if err := s.Start(); err != nil {
		t.Fatal(err)
	}
	defer s.Stop()

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for task execution")
	}
}

func TestRemoveAndUpdatePattern(t *testing.T) {
	s := NewScheduler()
	if err := s.SchedulePattern("t1", MustNewPattern("* * * * *"), TaskFunc(func() {})); err != nil {
		t.Fatal(err)
	}
	if s.Size() != 1 {
		t.Fatal("expected 1 task")
	}
	if !RemoveWithOptions("t1", WithDefaultScheduler(s)) {
		t.Fatal("RemoveWithOptions returned false")
	}
	if s.Size() != 0 {
		t.Fatal("expected 0 tasks after remove")
	}

	if err := s.SchedulePattern("t2", MustNewPattern("0 0 * * *"), TaskFunc(func() {})); err != nil {
		t.Fatal(err)
	}
	if err := UpdatePatternWithOptions("t2", "0 12 * * *", WithDefaultScheduler(s)); err != nil {
		t.Fatal(err)
	}
}

func TestApplyDefaultSchedulerOptionsNilOption(t *testing.T) {
	cfg := defaultSchedulerConfig{scheduler: NewScheduler()}
	WithDefaultScheduler(nil)(&cfg)
	if cfg.scheduler == nil {
		t.Fatal("nil WithDefaultScheduler should not replace scheduler")
	}
}

func TestShutdownWithOptions(t *testing.T) {
	s := NewScheduler()
	if err := ShutdownWithOptions(context.Background(), WithDefaultScheduler(s)); err != nil {
		t.Fatal(err)
	}
}

func TestDefaultSchedulerStop(t *testing.T) {
	s := DefaultScheduler()
	StopWithOptions(WithDefaultScheduler(s)) // should not panic
}

func TestDefaultSchedulerStartError(t *testing.T) {
	s := NewScheduler()
	if err := s.Start(); err != nil {
		t.Fatal(err)
	}
	// Starting an already started scheduler should return ErrSchedulerStarted
	StartWithOptions(WithDefaultScheduler(s))
}

func TestLastDayOfMonth(t *testing.T) {
	tests := []struct {
		month int
		leap  bool
		want  int
	}{
		{1, false, 31}, {2, false, 28}, {2, true, 29},
		{4, false, 30}, {12, false, 31},
	}
	for _, tc := range tests {
		if got := lastDayOfMonth(tc.month, tc.leap); got != tc.want {
			t.Fatalf("lastDayOfMonth(%d, %v) = %d, want %d", tc.month, tc.leap, got, tc.want)
		}
	}
}

func TestIsLeapYear(t *testing.T) {
	if !isLeapYear(2000) {
		t.Fatal("2000 should be leap")
	}
	if isLeapYear(1900) {
		t.Fatal("1900 should not be leap")
	}
	if !isLeapYear(2024) {
		t.Fatal("2024 should be leap")
	}
	if isLeapYear(2023) {
		t.Fatal("2023 should not be leap")
	}
}