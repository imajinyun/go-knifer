package cron

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestSchedulerLifecycle(t *testing.T) {
	s := NewScheduler()
	s.SetMatchSecond(true)

	var counter atomic.Int32
	id, err := s.ScheduleFunc("* * * * * *", func() { counter.Add(1) })
	if err != nil {
		t.Fatalf("schedule: %v", err)
	}
	if id == "" {
		t.Fatalf("empty id")
	}
	if s.Size() != 1 {
		t.Fatalf("expect 1 task, got %d", s.Size())
	}
	if err := s.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer s.Stop()
	// 等待至少 2 次执行
	time.Sleep(2500 * time.Millisecond)
	if counter.Load() < 1 {
		t.Fatalf("expect counter >= 1, got %d", counter.Load())
	}
	if !s.Deschedule(id) {
		t.Fatalf("expect deschedule ok")
	}
	if s.Size() != 0 {
		t.Fatalf("expect empty after deschedule")
	}
}

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

type testListener struct {
	started *atomic.Int32
	succ    *atomic.Int32
	failed  *atomic.Int32
}

func (l *testListener) OnStart(*TaskExecutor)       { l.started.Add(1) }
func (l *testListener) OnSucceeded(*TaskExecutor)   { l.succ.Add(1) }
func (l *testListener) OnFailed(*TaskExecutor, any) { l.failed.Add(1) }

func TestSchedulerUpdatePattern(t *testing.T) {
	s := NewScheduler()
	if err := s.ScheduleWithID("t1", "0 0 * * *", TaskFunc(func() {})); err != nil {
		t.Fatalf("schedule: %v", err)
	}
	if err := s.UpdatePattern("t1", "0 12 * * *"); err != nil {
		t.Fatalf("update: %v", err)
	}
	if got := s.GetPattern("t1").Raw(); got != "0 12 * * *" {
		t.Fatalf("expect updated pattern, got %q", got)
	}
	if err := s.UpdatePattern("nx", "* * * * *"); err == nil {
		t.Fatalf("expect error for unknown id")
	}
}

func TestSchedulerDuplicateID(t *testing.T) {
	s := NewScheduler()
	if err := s.ScheduleWithID("a", "* * * * *", TaskFunc(func() {})); err != nil {
		t.Fatalf("schedule: %v", err)
	}
	if err := s.ScheduleWithID("a", "* * * * *", TaskFunc(func() {})); err == nil {
		t.Fatalf("expect duplicate id error")
	}
}

func TestSchedulerStartTwice(t *testing.T) {
	s := NewScheduler()
	if err := s.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer s.Stop()
	if err := s.Start(); err == nil {
		t.Fatalf("expect error on second start")
	}
}
