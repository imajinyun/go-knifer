package vconf_test

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	knifer "github.com/imajinyun/go-knifer"
	"github.com/imajinyun/go-knifer/vconf"
)

func TestParseSettingFacade(t *testing.T) {
	s, err := vconf.Parse("name=gokit\ncount=42\nenabled=true\n[server]\nhost=127.0.0.1\nport=8080")
	if err != nil {
		t.Fatal(err)
	}
	if got := s.Get("name"); got != "gokit" {
		t.Fatalf("Get(name) = %q", got)
	}
	if got := s.GetInt("count", 0); got != 42 {
		t.Fatalf("GetInt(count) = %d", got)
	}
	if got := s.GetBool("enabled", false); !got {
		t.Fatal("GetBool(enabled) = false")
	}
	if got := s.GetByGroup("server", "host"); got != "127.0.0.1" {
		t.Fatalf("GetByGroup(server, host) = %q", got)
	}
	s.SetByGroup("server", "scheme", "http")
	if got := s.GetByGroup("server", "scheme"); got != "http" {
		t.Fatalf("SetByGroup() value = %q", got)
	}
	if !reflect.DeepEqual(s.Keys("server"), []string{"host", "port", "scheme"}) {
		t.Fatalf("Keys(server) = %#v", s.Keys("server"))
	}
}

func TestLoadAndParseYAMLFacade(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.setting")
	if err := os.WriteFile(path, []byte("app='demo'\n[db]\nuser=root"), 0o644); err != nil {
		t.Fatal(err)
	}
	s, err := vconf.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := s.Get("app"); got != "demo" {
		t.Fatalf("Load Get(app) = %q", got)
	}
	if got := s.GetByGroup("db", "user"); got != "root" {
		t.Fatalf("Load GetByGroup(db, user) = %q", got)
	}

	s, err = vconf.ParseYAML("app: gokit\nserver:\n  port: 8080\n  debug: true")
	if err != nil {
		t.Fatal(err)
	}
	if got := s.GetByGroup("server", "debug"); got != "true" {
		t.Fatalf("ParseYAML server.debug = %q", got)
	}
}

func TestNewAndParseBytesFacade(t *testing.T) {
	s := vconf.New()
	s.Set("k", "v")
	if got := s.GetOrDefault("k", "default"); got != "v" {
		t.Fatalf("GetOrDefault(k) = %q", got)
	}
	parsed, err := vconf.ParseBytes([]byte("x: 1"))
	if err != nil {
		t.Fatal(err)
	}
	if got := parsed.Get("x"); got != "1" {
		t.Fatalf("ParseBytes Get(x) = %q", got)
	}
}

func TestNilConfFacadeReadMethodsAreEmptyAndSafe(t *testing.T) {
	var s *vconf.Conf

	if got := s.Groups(); len(got) != 0 {
		t.Fatalf("Groups() = %v, want empty", got)
	}
	if got := s.Keys("missing"); len(got) != 0 {
		t.Fatalf("Keys(missing) = %v, want empty", got)
	}
	if got := s.ToMap(); len(got) != 0 {
		t.Fatalf("ToMap() = %v, want empty", got)
	}
}

func TestFacadeConfErrorContract(t *testing.T) {
	_, err := vconf.Parse("invalid-line")
	if err == nil {
		t.Fatal("Parse() error = nil, want invalid input")
	}
	if !errors.Is(err, knifer.ErrCodeInvalidInput) {
		t.Fatalf("errors.Is(err, ErrCodeInvalidInput) = false: %v", err)
	}
	code, ok := knifer.CodeOf(err)
	if !ok || code != knifer.ErrCodeInvalidInput {
		t.Fatalf("CodeOf(err) = %q, %v; want invalid input", code, ok)
	}
	var confErr *vconf.Error
	if !errors.As(err, &confErr) {
		t.Fatalf("errors.As(err, *vconf.Error) = false: %v", err)
	}
}

func TestAdvancedConfigFacade(t *testing.T) {
	t.Setenv("VCONF_HOST", "env.local")
	s, err := vconf.ParseTOML(`
name = "demo"
base = "http://${ENV:VCONF_HOST}"
[server]
port = 8080
debug = true
tags = ["api", "admin"]
[profile.prod.server]
port = 9090
`)
	if err != nil {
		t.Fatal(err)
	}
	if got := s.GetExpanded("base"); got != "http://env.local" {
		t.Fatalf("GetExpanded(base) = %q", got)
	}

	type serverConf struct {
		Port  int      `conf:"port"`
		Debug bool     `conf:"debug"`
		Tags  []string `conf:"tags"`
	}
	var cfg serverConf
	if err := s.BindGroup("server", &cfg); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(cfg, serverConf{Port: 8080, Debug: true, Tags: []string{"api", "admin"}}) {
		t.Fatalf("BindGroup() = %#v", cfg)
	}
	prod := s.ApplyProfile("prod")
	if got := prod.GetByGroup("server", "port"); got != "9090" {
		t.Fatalf("ApplyProfile(prod).server.port = %q", got)
	}
}

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
	changes := make(chan string, 1)
	stop, err := vconf.Watch(watchPath, 10*time.Millisecond, func(c *vconf.Conf, err error) {
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
	time.Sleep(20 * time.Millisecond)
	if err := os.WriteFile(watchPath, []byte("name=two"), 0o644); err != nil {
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
