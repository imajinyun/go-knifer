package conf

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWatchReloadsOnChange(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app.setting")
	if err := os.WriteFile(path, []byte("name=one"), 0o644); err != nil {
		t.Fatal(err)
	}
	changes := make(chan string, 1)
	stop, err := Watch(path, 10*time.Millisecond, func(c *Conf, err error) {
		if err != nil {
			changes <- "err:" + err.Error()
			return
		}
		changes <- c.Get("name")
	})
	if err != nil {
		t.Fatal(err)
	}
	defer stop()

	time.Sleep(20 * time.Millisecond)
	if err := os.WriteFile(path, []byte("name=two"), 0o644); err != nil {
		t.Fatal(err)
	}
	select {
	case got := <-changes:
		if got != "two" {
			t.Fatalf("watch change = %q", got)
		}
	case <-time.After(time.Second):
		t.Fatal("watch did not report change")
	}
}

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

func TestWatchCallbacksPanicsAreIsolated(t *testing.T) {
	tick := make(chan time.Time)
	ticker := &watchTestTicker{stopped: make(chan struct{})}
	info := fakeFileInfo{name: "app.setting", size: int64(len("name=one")), modTime: time.Unix(1, 0)}
	content := []byte("name=one")
	started := make(chan struct{})

	stop, err := WatchWithOptions("app.setting", WatchOptions{
		Interval:       time.Hour,
		CompareContent: true,
		TickerFactory:  func(time.Duration) (<-chan time.Time, WatchTicker) { return tick, ticker },
		Runner: func(fn func()) {
			close(started)
			go fn()
		},
		Stat: func(string) (os.FileInfo, error) { return info, nil },
		ReadFile: func(string, int64) ([]byte, error) {
			return content, nil
		},
		OnEvent: func(WatchEvent) { panic("event") },
	}, func(*Conf, error) { panic("change") })
	if err != nil {
		t.Fatal(err)
	}
	<-started
	info.size = int64(len("name=two"))
	info.modTime = time.Unix(2, 0)
	content = []byte("name=two")
	tick <- time.Unix(2, 0)

	done := make(chan struct{})
	go func() {
		stop()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("stop blocked after callback panic")
	}
}

func TestWatchWithOptionsCompareContentAndEvent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app.setting")
	if err := os.WriteFile(path, []byte("name=one"), 0o644); err != nil {
		t.Fatal(err)
	}
	changes := make(chan string, 1)
	events := make(chan WatchEvent, 1)
	stop, err := WatchWithOptions(path, WatchOptions{
		Interval:       10 * time.Millisecond,
		Debounce:       5 * time.Millisecond,
		CompareContent: true,
		OnEvent: func(event WatchEvent) {
			events <- event
		},
	}, func(c *Conf, err error) {
		if err != nil {
			changes <- "err:" + err.Error()
			return
		}
		changes <- c.Get("name")
	})
	if err != nil {
		t.Fatal(err)
	}
	defer stop()

	time.Sleep(20 * time.Millisecond)
	if err := os.WriteFile(path, []byte("name=two"), 0o644); err != nil {
		t.Fatal(err)
	}
	select {
	case got := <-changes:
		if got != "two" {
			t.Fatalf("watch change = %q", got)
		}
	case <-time.After(time.Second):
		t.Fatal("watch did not report change")
	}
	select {
	case event := <-events:
		if event.Path != path || event.Size == 0 {
			t.Fatalf("watch event = %#v", event)
		}
	case <-time.After(time.Second):
		t.Fatal("watch did not report event")
	}
}

func TestWatchWithOptionsProviders(t *testing.T) {
	tick := make(chan time.Time)
	ticker := &watchTestTicker{stopped: make(chan struct{})}
	reads := 0
	changes := make(chan string, 1)
	events := make(chan WatchEvent, 1)
	info := fakeFileInfo{name: "app.setting", size: int64(len("name=one")), modTime: time.Unix(1, 0)}
	content := []byte("name=one")

	stop, err := WatchWithOptions("app.setting", WatchOptions{
		Interval:       time.Hour,
		Debounce:       time.Nanosecond,
		CompareContent: true,
		TickerFactory: func(delay time.Duration) (<-chan time.Time, WatchTicker) {
			if delay != time.Hour {
				t.Fatalf("ticker delay = %s, want %s", delay, time.Hour)
			}
			return tick, ticker
		},
		After: func(delay time.Duration) <-chan time.Time {
			if delay != time.Nanosecond {
				t.Fatalf("debounce delay = %s, want %s", delay, time.Nanosecond)
			}
			ch := make(chan time.Time, 1)
			ch <- time.Unix(3, 0)
			return ch
		},
		Stat: func(path string) (os.FileInfo, error) {
			if path != "app.setting" {
				t.Fatalf("stat path = %q", path)
			}
			return info, nil
		},
		ReadFile: func(path string, maxBytes int64) ([]byte, error) {
			reads++
			if path != "app.setting" || maxBytes != 64 {
				t.Fatalf("read path=%q maxBytes=%d", path, maxBytes)
			}
			return content, nil
		},
		LoadOptions: LoadOptions{MaxBytes: 64},
		OnEvent: func(event WatchEvent) {
			events <- event
		},
	}, func(c *Conf, err error) {
		if err != nil {
			changes <- "err:" + err.Error()
			return
		}
		changes <- c.Get("name")
	})
	if err != nil {
		t.Fatal(err)
	}

	info.size = int64(len("name=two"))
	info.modTime = time.Unix(2, 0)
	content = []byte("name=two")
	tick <- time.Unix(2, 0)
	select {
	case got := <-changes:
		if got != "two" {
			t.Fatalf("watch change = %q", got)
		}
	case <-time.After(time.Second):
		t.Fatal("watch did not report provider-driven change")
	}
	select {
	case event := <-events:
		if event.Path != "app.setting" || event.Size != int64(len("name=two")) {
			t.Fatalf("watch event = %#v", event)
		}
	case <-time.After(time.Second):
		t.Fatal("watch did not report provider-driven event")
	}
	stop()
	select {
	case <-ticker.stopped:
	case <-time.After(time.Second):
		t.Fatal("watch ticker was not stopped")
	}
	if reads < 3 {
		t.Fatalf("read count = %d, want at least 3", reads)
	}
}

type watchTestTicker struct {
	stopped chan struct{}
}

func (t *watchTestTicker) Stop() { close(t.stopped) }

type fakeFileInfo struct {
	name    string
	size    int64
	modTime time.Time
}

func (f fakeFileInfo) Name() string       { return f.name }
func (f fakeFileInfo) Size() int64        { return f.size }
func (f fakeFileInfo) Mode() fs.FileMode  { return 0o644 }
func (f fakeFileInfo) ModTime() time.Time { return f.modTime }
func (f fakeFileInfo) IsDir() bool        { return false }
func (f fakeFileInfo) Sys() any           { return nil }
