package vcron_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/imajinyun/go-knifer/vcron"
)

func TestFacadeSchedulerWithOptions(t *testing.T) {
	loc := time.FixedZone("facade", 8*60*60)
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	s := vcron.NewSchedulerWithOptions(
		vcron.WithLocation(loc),
		vcron.WithMatchSecond(true),
		vcron.WithIDGenerator(func() string { return "facade-task" }),
		vcron.WithClock(func() time.Time { return now }),
		vcron.WithSleeper(func(d time.Duration, stopCh <-chan struct{}) bool {
			now = now.Add(d)
			return true
		}),
		vcron.WithExecutor(func(fn func()) { fn() }),
	)
	if s.Config().Location != loc {
		t.Fatalf("scheduler location = %v, want %v", s.Config().Location, loc)
	}
	if !s.IsMatchSecond() {
		t.Fatal("scheduler should match seconds")
	}
	id, err := s.ScheduleFunc("* * * * * *", func() {})
	if err != nil {
		t.Fatalf("ScheduleFunc with options: %v", err)
	}
	if id != "facade-task" {
		t.Fatalf("scheduled id = %q, want facade-task", id)
	}
}

func TestFacadeSchedulerIDRandomReaderOption(t *testing.T) {
	s := vcron.NewSchedulerWithOptions(vcron.WithIDRandomReader(bytes.NewReader([]byte{8, 7, 6, 5, 4, 3, 2, 1})))
	id, err := s.ScheduleFunc("* * * * *", func() {})
	if err != nil {
		t.Fatalf("ScheduleFunc: %v", err)
	}
	if id != "0807060504030201" {
		t.Fatalf("id = %q, want 0807060504030201", id)
	}
}
