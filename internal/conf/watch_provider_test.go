package conf

import (
	"os"
	"testing"
	"time"
)

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
