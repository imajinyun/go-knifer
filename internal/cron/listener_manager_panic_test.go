package cron

import (
	"sync/atomic"
	"testing"
)

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
