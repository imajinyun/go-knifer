package vconf_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/imajinyun/go-knifer/vconf"
)

func TestLoadProfileAndWatchFacade(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.yaml")
	if err := os.WriteFile(path, []byte("app:\n  name: base\nprofile:\n  dev:\n    app:\n      name: dev"), 0o644); err != nil {
		t.Fatal(err)
	}
	c, err := vconf.LoadProfile(path, "dev")
	if err != nil {
		t.Fatal(err)
	}
	if got := c.GetByGroup("app", "name"); got != "dev" {
		t.Fatalf("LoadProfile yaml app.name = %q", got)
	}

	watchPath := filepath.Join(dir, "watch.setting")
	if err := os.WriteFile(watchPath, []byte("name=one"), 0o644); err != nil {
		t.Fatal(err)
	}
	ticks := make(chan time.Time, 1)
	changes := make(chan string, 1)
	stop, err := vconf.WatchWithOptions(watchPath, vconf.WatchOptions{
		Interval:       time.Hour,
		CompareContent: true,
		TickerFactory: func(time.Duration) (<-chan time.Time, vconf.WatchTicker) {
			return ticks, facadeWatchTicker{}
		},
	}, func(c *vconf.Conf, err error) {
		if err != nil {
			changes <- "err"
			return
		}
		changes <- c.Get("name")
	})
	if err != nil {
		t.Fatal(err)
	}
	defer stop()
	if err := os.WriteFile(watchPath, []byte("name=two"), 0o644); err != nil {
		t.Fatal(err)
	}
	ticks <- time.Now()
	select {
	case got := <-changes:
		if got != "two" {
			t.Fatalf("watch change = %q", got)
		}
	case <-time.After(time.Second):
		t.Fatal("watch did not report change")
	}
}

type facadeWatchTicker struct{}

func (facadeWatchTicker) Stop() {}

func TestWatchOptionsProviderTypesFacade(t *testing.T) {
	ticks := make(chan time.Time)
	var factory vconf.WatchTickerFactory = func(time.Duration) (<-chan time.Time, vconf.WatchTicker) {
		return ticks, facadeWatchTicker{}
	}
	_ = vconf.WatchOptions{TickerFactory: factory}
}

func TestWatchWithOptionsFacade(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app.setting")
	if err := os.WriteFile(path, []byte("name=one"), 0o644); err != nil {
		t.Fatal(err)
	}
	ticks := make(chan time.Time, 1)
	changes := make(chan string, 1)
	stop, err := vconf.WatchWithOptions(path, vconf.WatchOptions{
		Interval:       time.Hour,
		CompareContent: true,
		TickerFactory: func(time.Duration) (<-chan time.Time, vconf.WatchTicker) {
			return ticks, facadeWatchTicker{}
		},
	}, func(c *vconf.Conf, err error) {
		if err != nil {
			changes <- "err"
			return
		}
		changes <- c.Get("name")
	})
	if err != nil {
		t.Fatal(err)
	}
	defer stop()
	if err := os.WriteFile(path, []byte("name=two"), 0o644); err != nil {
		t.Fatal(err)
	}
	ticks <- time.Now()
	select {
	case got := <-changes:
		if got != "two" {
			t.Fatalf("watch change = %q", got)
		}
	case <-time.After(time.Second):
		t.Fatal("watch did not report change")
	}
}
