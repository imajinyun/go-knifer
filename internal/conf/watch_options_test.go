package conf

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWatchWithOptionsUsesRunner(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app.setting")
	if err := os.WriteFile(path, []byte("name=one"), 0o644); err != nil {
		t.Fatal(err)
	}
	ticks := make(chan time.Time)
	ticker := &watchTestTicker{stopped: make(chan struct{})}
	runnerCalled := make(chan struct{}, 1)
	stop, err := WatchWithOptions(path, WatchOptions{
		Interval: 10 * time.Second,
		TickerFactory: func(delay time.Duration) (<-chan time.Time, WatchTicker) {
			if delay != 10*time.Second {
				t.Fatalf("ticker delay = %s, want 10s", delay)
			}
			return ticks, ticker
		},
		Runner: func(fn func()) {
			runnerCalled <- struct{}{}
			go fn()
		},
	}, func(*Conf, error) {})
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-runnerCalled:
	case <-time.After(time.Second):
		t.Fatal("watch runner was not used")
	}
	stop()
	select {
	case <-ticker.stopped:
	case <-time.After(time.Second):
		t.Fatal("watch ticker was not stopped")
	}
	stop()
}

func TestWatchWithOptionsRejectsNilCallback(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app.setting")
	if err := os.WriteFile(path, []byte("name=one"), 0o644); err != nil {
		t.Fatal(err)
	}
	if stop, err := WatchWithOptions(path, WatchOptions{}, nil); err == nil || stop != nil {
		t.Fatalf("WatchWithOptions nil callback stop nil=%v err=%v, want error", stop == nil, err)
	}
}
