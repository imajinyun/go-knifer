package cron

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestSchedulerListener(t *testing.T) {
	s := NewScheduler()
	s.SetMatchSecond(true)

	var started, succ, failed atomic.Int32
	s.AddListener(&testListener{
		started: &started, succ: &succ, failed: &failed,
	})

	_, err := s.ScheduleFunc("* * * * * *", func() {})
	if err != nil {
		t.Fatalf("schedule: %v", err)
	}
	_, err = s.ScheduleFunc("* * * * * *", func() { panic("boom") })
	if err != nil {
		t.Fatalf("schedule: %v", err)
	}
	if err := s.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer s.Stop()
	time.Sleep(1500 * time.Millisecond)
	if started.Load() < 2 {
		t.Fatalf("expect started >= 2, got %d", started.Load())
	}
	if succ.Load() < 1 {
		t.Fatalf("expect succ >= 1")
	}
	if failed.Load() < 1 {
		t.Fatalf("expect failed >= 1")
	}
}
