package cron

import (
	"errors"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	knifer "github.com/imajinyun/go-knifer"
)

func TestSchedulerConfigSettersIgnoredWhileStarted(t *testing.T) {
	loc := time.FixedZone("before", 3600)
	after := time.FixedZone("after", 7200)
	s := NewSchedulerWithOptions(WithLocation(loc), WithMatchSecond(true))
	if err := s.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	if err := s.SetMatchSecondE(false); !errors.Is(err, ErrSchedulerStarted) {
		t.Fatalf("SetMatchSecondE while started = %v, want ErrSchedulerStarted", err)
	} else if !errors.Is(err, knifer.ErrCodeUnsupported) {
		t.Fatalf("SetMatchSecondE while started = %v, want ErrCodeUnsupported", err)
	}
	if err := s.SetTimeZoneE(after); !errors.Is(err, ErrSchedulerStarted) {
		t.Fatalf("SetTimeZoneE while started = %v, want ErrSchedulerStarted", err)
	}
	s.SetMatchSecond(false).SetTimeZone(after)
	cfg := s.Config()
	if cfg.Location != loc || !cfg.MatchSecond {
		t.Fatalf("started scheduler config mutated: %#v", cfg)
	}
	s.Stop()
	if err := s.SetMatchSecondE(false); err != nil {
		t.Fatalf("SetMatchSecondE after stop: %v", err)
	}
	if err := s.SetTimeZoneE(after); err != nil {
		t.Fatalf("SetTimeZoneE after stop: %v", err)
	}
	s.SetMatchSecond(false).SetTimeZone(after)
	cfg = s.Config()
	if cfg.Location != after || cfg.MatchSecond {
		t.Fatalf("stopped scheduler config not updated: %#v", cfg)
	}
}

func TestSchedulerConfigReturnsSnapshot(t *testing.T) {
	s := NewSchedulerWithOptions(WithMatchSecond(true))
	cfg := s.Config()
	cfg.MatchSecond = false
	if !s.IsMatchSecond() {
		t.Fatal("mutating Config snapshot changed scheduler")
	}
}

func TestSchedulerOptions(t *testing.T) {
	loc := time.FixedZone("test", 8*3600)
	var submitted atomic.Int32
	var runnerCalls atomic.Int32
	var sleepCalls atomic.Int32
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	s := NewSchedulerWithOptions(
		WithLocation(loc),
		WithMatchSecond(true),
		WithIDGenerator(func() string { return "custom-id" }),
		WithClock(func() time.Time { return now }),
		WithSleeper(func(d time.Duration, stopCh <-chan struct{}) bool {
			sleepCalls.Add(1)
			now = now.Add(d)
			select {
			case <-stopCh:
				return false
			default:
				return true
			}
		}),
		WithExecutor(func(fn func()) {
			submitted.Add(1)
			fn()
		}),
		WithRunner(func(fn func()) {
			runnerCalls.Add(1)
			go fn()
		}),
	)
	if s.Config().Location != loc || !s.IsMatchSecond() {
		t.Fatalf("scheduler options not applied: %#v", s.Config())
	}
	id, err := s.ScheduleFunc("* * * * * *", func() {})
	if err != nil {
		t.Fatalf("schedule with custom id: %v", err)
	}
	if id != "custom-id" {
		t.Fatalf("custom id = %q", id)
	}
	s.submit(func() {})
	if submitted.Load() != 1 {
		t.Fatalf("custom executor not used")
	}
	if s.nowMillis() != now.UnixMilli() {
		t.Fatalf("custom clock not used")
	}
	if !s.sleep(time.Millisecond, make(chan struct{})) || sleepCalls.Load() != 1 {
		t.Fatalf("custom sleeper not used")
	}
	if err := s.Start(); err != nil {
		t.Fatalf("start with custom runner: %v", err)
	}
	if runnerCalls.Load() != 1 {
		t.Fatalf("custom runner calls = %d, want 1", runnerCalls.Load())
	}
	s.Stop()
}

func TestConfigOptions(t *testing.T) {
	loc := time.FixedZone("config", 9*3600)
	cfg := NewConfigWithOptions(WithConfigLocation(loc), WithConfigMatchSecond(true))
	if cfg.Location != loc || !cfg.MatchSecond {
		t.Fatalf("config options not applied: %#v", cfg)
	}
	cfg = NewConfigWithOptions(WithConfigLocation(nil))
	if cfg.Location == nil {
		t.Fatal("nil config location should fall back to local")
	}
}

func TestSchedulerPatternOptions(t *testing.T) {
	parseCalls := 0
	s := NewSchedulerWithOptions(
		WithIDGenerator(func() string { return "custom-pattern-id" }),
		WithSchedulerPatternOptions(WithPatternIntParser(func(text string) (int, error) {
			parseCalls++
			if text == "custom" {
				return 30, nil
			}
			return strconv.Atoi(text)
		})),
	)
	id, err := s.ScheduleFunc("custom * * * *", func() {})
	if err != nil {
		t.Fatalf("ScheduleFunc with pattern options: %v", err)
	}
	if id != "custom-pattern-id" || parseCalls == 0 {
		t.Fatalf("pattern options not used: id=%q parseCalls=%d", id, parseCalls)
	}
	if got := s.GetPattern(id); got == nil || !got.Match(time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC), false) {
		t.Fatalf("stored custom pattern = %#v", got)
	}
}
