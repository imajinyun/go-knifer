package conf

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

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
