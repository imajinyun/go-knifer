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

func TestSchedulerListenerPanicsAreIsolated(t *testing.T) {
	s := NewSchedulerWithOptions(WithExecutor(func(fn func()) { fn() }))
	defer s.Stop()

	var starts, successes, failures atomic.Int32
	s.AddListener(&testListener{started: &starts})
	s.AddListener(&panicListener{onStart: true})
	s.AddListener(&testListener{succ: &successes})
	s.AddListener(&panicListener{onSucceeded: true})
	s.AddListener(&testListener{failed: &failures})
	s.AddListener(&panicListener{onFailed: true})

	s.executorMgr.spawn(NewCronTask("ok", MustNewPattern("* * * * *"), TaskFunc(func() {})))
	s.executorMgr.spawn(NewCronTask("bad", MustNewPattern("* * * * *"), TaskFunc(func() { panic("task") })))

	if got := starts.Load(); got != 2 {
		t.Fatalf("starts = %d, want 2", got)
	}
	if got := successes.Load(); got != 1 {
		t.Fatalf("successes = %d, want 1", got)
	}
	if got := failures.Load(); got != 1 {
		t.Fatalf("failures = %d, want 1", got)
	}
}

type testListener struct {
	started *atomic.Int32
	succ    *atomic.Int32
	failed  *atomic.Int32
}

func (l *testListener) OnStart(*TaskExecutor) {
	if l.started != nil {
		l.started.Add(1)
	}
}

func (l *testListener) OnSucceeded(*TaskExecutor) {
	if l.succ != nil {
		l.succ.Add(1)
	}
}

func (l *testListener) OnFailed(*TaskExecutor, any) {
	if l.failed != nil {
		l.failed.Add(1)
	}
}

type panicListener struct {
	onStart     bool
	onSucceeded bool
	onFailed    bool
}

func (l *panicListener) OnStart(*TaskExecutor) {
	if l.onStart {
		panic("start listener")
	}
}

func (l *panicListener) OnSucceeded(*TaskExecutor) {
	if l.onSucceeded {
		panic("success listener")
	}
}

func (l *panicListener) OnFailed(*TaskExecutor, any) {
	if l.onFailed {
		panic("failed listener")
	}
}
