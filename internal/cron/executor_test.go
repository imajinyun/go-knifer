package cron

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSchedulerLauncherPanicsAreIsolated(t *testing.T) {
	s := NewSchedulerWithOptions(WithExecutor(func(fn func()) { fn() }))
	if err := s.taskTable.Add("bad", nil, TaskFunc(func() {})); err != nil {
		t.Fatalf("add invalid task: %v", err)
	}
	s.launcherMgr.spawn(time.Now().UnixMilli())
	if got := s.LaunchingCount(); got != 0 {
		t.Fatalf("LaunchingCount = %d, want 0", got)
	}
}

func TestSchedulerShutdownWaitsForRunningTasks(t *testing.T) {
	start := make(chan struct{})
	finish := make(chan struct{})
	s := NewSchedulerWithOptions(WithExecutor(func(fn func()) { go fn() }))
	s.executorMgr.spawn(NewCronTask("slow", MustNewPattern("* * * * *"), TaskFunc(func() {
		close(start)
		<-finish
	})))
	<-start
	if got := s.RunningCount(); got != 1 {
		t.Fatalf("RunningCount = %d, want 1", got)
	}
	done := make(chan error, 1)
	go func() { done <- s.Shutdown(context.Background()) }()
	select {
	case err := <-done:
		t.Fatalf("Shutdown returned before task finished: %v", err)
	case <-time.After(20 * time.Millisecond):
	}
	close(finish)
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Shutdown error: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Shutdown did not return after task finished")
	}
	if got := s.RunningCount(); got != 0 {
		t.Fatalf("RunningCount after shutdown = %d, want 0", got)
	}
}

func TestSchedulerShutdownContextTimeout(t *testing.T) {
	start := make(chan struct{})
	finish := make(chan struct{})
	s := NewSchedulerWithOptions(WithExecutor(func(fn func()) { go fn() }))
	s.executorMgr.spawn(NewCronTask("slow", MustNewPattern("* * * * *"), TaskFunc(func() {
		close(start)
		<-finish
	})))
	<-start
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	if err := s.Shutdown(ctx); err == nil {
		close(finish)
		t.Fatal("Shutdown should return context timeout")
	}
	close(finish)
	s.Wait()
}

func TestSchedulerShutdownWaitsForLaunchersBeforeExecutors(t *testing.T) {
	launcherStarted := make(chan struct{})
	allowLauncher := make(chan struct{})
	taskDone := make(chan struct{})
	var launcherSeen atomic.Bool

	s := NewSchedulerWithOptions(
		WithMatchSecond(true),
		WithClock(func() time.Time { return time.Unix(1, 0) }),
		WithSleeper(func(time.Duration, <-chan struct{}) bool { return true }),
		WithExecutor(func(fn func()) {
			if !launcherSeen.Swap(true) {
				close(launcherStarted)
				<-allowLauncher
			}
			go fn()
		}),
	)
	if _, err := s.ScheduleFunc("* * * * * *", func() { close(taskDone) }); err != nil {
		t.Fatalf("schedule: %v", err)
	}
	if err := s.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}

	select {
	case <-launcherStarted:
	case <-time.After(time.Second):
		t.Fatal("launcher did not start")
	}
	if got := s.LaunchingCount(); got != 1 {
		t.Fatalf("LaunchingCount = %d, want 1", got)
	}

	shutdownDone := make(chan error, 1)
	go func() { shutdownDone <- s.Shutdown(context.Background()) }()
	select {
	case err := <-shutdownDone:
		t.Fatalf("Shutdown returned before launcher was released: %v", err)
	case <-time.After(20 * time.Millisecond):
	}

	close(allowLauncher)
	select {
	case err := <-shutdownDone:
		if err != nil {
			t.Fatalf("Shutdown error = %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Shutdown did not finish after launcher was released")
	}
	select {
	case <-taskDone:
	case <-time.After(time.Second):
		t.Fatal("task did not run before Shutdown returned")
	}
	if got := s.LaunchingCount(); got != 0 {
		t.Fatalf("LaunchingCount after shutdown = %d, want 0", got)
	}
}

func TestSchedulerExecutorAndRunnerConcurrentReplacement(t *testing.T) {
	s := NewSchedulerWithOptions(WithExecutor(func(fn func()) { fn() }), WithRunner(func(fn func()) { fn() }))
	var wg sync.WaitGroup
	for i := 0; i < 64; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			s.SetExecutor(func(fn func()) { fn() })
			s.SetRunner(func(fn func()) { fn() })
		}()
		go func() {
			defer wg.Done()
			s.submit(func() {})
			s.run(func() {})
		}()
	}
	wg.Wait()
}
