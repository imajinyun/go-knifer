package vconf_test

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

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

func TestFacadeParserProviderOptions(t *testing.T) {
	parsed, err := vconf.ParseByExtWithOptions("app.custom", []byte("ignored"), vconf.WithParserForExt(".custom", func(data []byte) (*vconf.Conf, error) {
		c := vconf.New()
		c.Set("from", string(data))
		return c, nil
	}))
	if err != nil || parsed.Get("from") != "ignored" {
		t.Fatalf("ParseByExtWithOptions = %#v, %v", parsed, err)
	}

	yamlCalled := false
	_, err = vconf.ParseYAMLFullWithOptions("ignored", vconf.WithYAMLUnmarshalFunc(func(data []byte, out any) error {
		yamlCalled = true
		return errors.New("yaml provider failed")
	}))
	if err == nil || !yamlCalled {
		t.Fatalf("ParseYAMLFullWithOptions err=%v called=%v", err, yamlCalled)
	}

	tomlCalled := false
	_, err = vconf.ParseTOMLWithOptions("ignored", vconf.WithTOMLUnmarshalFunc(func(data []byte, out any) error {
		tomlCalled = true
		return errors.New("toml provider failed")
	}))
	if err == nil || !tomlCalled {
		t.Fatalf("ParseTOMLWithOptions err=%v called=%v", err, tomlCalled)
	}
}
