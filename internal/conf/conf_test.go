package conf

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	knifer "github.com/imajinyun/go-knifer"
)

func TestParseSetting(t *testing.T) {
	s, err := Parse("name = gokit\n[server]\nport=8080\ndebug=true")
	if err != nil {
		t.Fatal(err)
	}
	if got := s.Get("name"); got != "gokit" {
		t.Fatalf("Get(name) = %q", got)
	}
	if got := s.GetByGroup("server", "port"); got != "8080" {
		t.Fatalf("GetByGroup(server, port) = %q", got)
	}
	if got := s.GetOrDefault("missing", "def"); got != "def" {
		t.Fatalf("GetOrDefault() = %q", got)
	}
}

func TestParseYAML(t *testing.T) {
	s, err := ParseYAML("app: gokit\nserver:\n  port: 8080")
	if err != nil {
		t.Fatal(err)
	}
	if got := s.Get("app"); got != "gokit" {
		t.Fatalf("Get(app) = %q", got)
	}
	if got := s.GetByGroup("server", "port"); got != "8080" {
		t.Fatalf("GetByGroup(server, port) = %q", got)
	}
}

func TestNilConfReadMethodsAreEmptyAndSafe(t *testing.T) {
	var s *Conf

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

func TestGroupsKeysAndToMapKeepStableSemantics(t *testing.T) {
	s := New()
	s.Set("root", "value")
	s.SetByGroup("server", "port", "8080")
	s.SetByGroup("server", "host", "localhost")
	s.SetByGroup("app", "name", "gokit")

	if got, want := s.Groups(), []string{"", "app", "server"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Groups() = %v, want %v", got, want)
	}
	if got, want := s.Keys("server"), []string{"host", "port"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Keys(server) = %v, want %v", got, want)
	}

	m := s.ToMap()
	m["server"]["port"] = "9090"
	if got := s.GetByGroup("server", "port"); got != "8080" {
		t.Fatalf("ToMap() returned shallow copy, source port = %q", got)
	}
}

func TestConfErrorContract(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "missing.setting"))
	assertConfCode(t, err, knifer.ErrCodeNotFound)
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Load missing file should preserve os.ErrNotExist: %v", err)
	}

	_, err = Parse("invalid-line")
	assertConfCode(t, err, knifer.ErrCodeInvalidInput)

	_, err = Parse("=empty")
	assertConfCode(t, err, knifer.ErrCodeInvalidInput)

	_, err = ParseYAML("invalid-yaml-line")
	assertConfCode(t, err, knifer.ErrCodeInvalidInput)
}

func TestExpandVariables(t *testing.T) {
	t.Setenv("CONF_ENV_HOST", "envhost")
	s, err := Parse(`
host=localhost
base=http://${host}:8080
env=${ENV:CONF_ENV_HOST}
missing=${missing:fallback}
[db]
host=db.local
url=postgres://${db.host}/${name:app}
`)
	if err != nil {
		t.Fatal(err)
	}

	if got := s.GetExpanded("base"); got != "http://localhost:8080" {
		t.Fatalf("GetExpanded(base) = %q", got)
	}
	if got := s.GetExpanded("env"); got != "envhost" {
		t.Fatalf("GetExpanded(env) = %q", got)
	}
	if got := s.GetExpanded("missing"); got != "fallback" {
		t.Fatalf("GetExpanded(missing) = %q", got)
	}
	if got := s.GetByGroupExpanded("db", "url"); got != "postgres://db.local/app" {
		t.Fatalf("GetByGroupExpanded(db,url) = %q", got)
	}
	if got := s.Expand().Get("base"); got != "http://localhost:8080" {
		t.Fatalf("Expand().Get(base) = %q", got)
	}
}

func TestParseYAMLFullAndBind(t *testing.T) {
	s, err := ParseYAMLFull(`
app: demo
server:
  host: 127.0.0.1
  port: 8080
  debug: true
  tags: [api, admin]
`)
	if err != nil {
		t.Fatal(err)
	}
	if got := s.Get("app"); got != "demo" {
		t.Fatalf("Get(app) = %q", got)
	}
	if got := s.GetByGroup("server", "host"); got != "127.0.0.1" {
		t.Fatalf("GetByGroup(server,host) = %q", got)
	}

	type serverConf struct {
		Host  string   `conf:"host"`
		Port  int      `conf:"port"`
		Debug bool     `conf:"debug"`
		Tags  []string `conf:"tags"`
	}
	var cfg serverConf
	if err := s.BindGroup("server", &cfg); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(cfg, serverConf{Host: "127.0.0.1", Port: 8080, Debug: true, Tags: []string{"api", "admin"}}) {
		t.Fatalf("BindGroup() = %#v", cfg)
	}
}

func TestParseTOMLAndProfile(t *testing.T) {
	s, err := ParseTOML(`
name = "demo"
tags = ["a", "b"]
[server]
port = 8080
[profile.dev]
name = "dev-demo"
[profile.dev.server]
port = 9090
`)
	if err != nil {
		t.Fatal(err)
	}
	if got := s.Get("tags"); got != "a,b" {
		t.Fatalf("Get(tags) = %q", got)
	}
	dev := s.ApplyProfile("dev")
	if got := dev.Get("name"); got != "dev-demo" {
		t.Fatalf("ApplyProfile(dev).Get(name) = %q", got)
	}
	if got := dev.GetByGroup("server", "port"); got != "9090" {
		t.Fatalf("ApplyProfile(dev).server.port = %q", got)
	}
}

func TestLoadProfileAndParseByExt(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.toml")
	if err := os.WriteFile(path, []byte("name='base'\n[profile.test]\nname='test'"), 0o644); err != nil {
		t.Fatal(err)
	}
	c, err := LoadProfile(path, "test")
	if err != nil {
		t.Fatal(err)
	}
	if got := c.Get("name"); got != "test" {
		t.Fatalf("LoadProfile name = %q", got)
	}

	yamlConf, err := ParseByExt("app.yaml", []byte("app:\n  name: demo"))
	if err != nil {
		t.Fatal(err)
	}
	if got := yamlConf.GetByGroup("app", "name"); got != "demo" {
		t.Fatalf("ParseByExt yaml app.name = %q", got)
	}
}

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

func assertConfCode(t *testing.T, err error, code knifer.ErrCode) {
	t.Helper()
	if err == nil {
		t.Fatalf("err = nil, want %s", code)
	}
	if !errors.Is(err, code) {
		t.Fatalf("errors.Is(%v, %s) = false", err, code)
	}
	got, ok := knifer.CodeOf(err)
	if !ok || got != code {
		t.Fatalf("CodeOf(%v) = %q, %v; want %q, true", err, got, ok, code)
	}
}
