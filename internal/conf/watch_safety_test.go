package conf

import (
	"os"
	"testing"
	"time"
)

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
