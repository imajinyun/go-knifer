package conf

import (
	"errors"
	"reflect"
	"strconv"
	"testing"

	knifer "github.com/imajinyun/go-knifer"
)

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

func TestCloneReturnsIndependentCopy(t *testing.T) {
	s := New()
	s.Set("root", "value")
	s.SetByGroup("server", "port", "8080")

	clone := s.Clone()
	clone.Set("root", "changed")
	clone.SetByGroup("server", "port", "9090")
	clone.SetByGroup("server", "host", "localhost")

	if got := s.Get("root"); got != "value" {
		t.Fatalf("source root changed to %q", got)
	}
	if got := s.GetByGroup("server", "port"); got != "8080" {
		t.Fatalf("source server.port changed to %q", got)
	}
	if got := s.GetByGroup("server", "host"); got != "" {
		t.Fatalf("source server.host changed to %q", got)
	}
	if got := clone.GetByGroup("server", "host"); got != "localhost" {
		t.Fatalf("clone server.host = %q", got)
	}
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

func TestExpandWithOptionsUsesCustomEnvLookup(t *testing.T) {
	s, err := Parse(`
host=${ENV:CONF_HOST}
base=http://${host}:${port:8080}
[db]
url=postgres://${ENV:CONF_DB_HOST}/${name:app}
`)
	if err != nil {
		t.Fatal(err)
	}
	lookup := func(key string) string {
		switch key {
		case "CONF_HOST":
			return "option.local"
		case "CONF_DB_HOST":
			return "db.option.local"
		default:
			return ""
		}
	}
	if got := s.GetExpandedWithOptions("host", WithEnvLookup(lookup)); got != "option.local" {
		t.Fatalf("GetExpandedWithOptions(host) = %q", got)
	}
	if got := s.GetExpandedWithOptions("base", WithEnvLookup(lookup)); got != "http://option.local:8080" {
		t.Fatalf("GetExpandedWithOptions(base) = %q", got)
	}
	if got := s.GetByGroupExpandedWithOptions("db", "url", WithEnvLookup(lookup)); got != "postgres://db.option.local/app" {
		t.Fatalf("GetByGroupExpandedWithOptions(db.url) = %q", got)
	}
	expanded := s.ExpandWithOptions(WithEnvLookup(lookup))
	if got := expanded.Get("host"); got != "option.local" {
		t.Fatalf("ExpandWithOptions host = %q", got)
	}
}

func TestTypedGettersWithOptionsUseParsers(t *testing.T) {
	s := New()
	s.Set("port", "custom-int")
	s.Set("debug", "custom-bool")
	s.SetByGroup("server", "port", "9090")
	s.SetByGroup("server", "debug", "true")

	intCalled := false
	if got := s.GetIntWithOptions("port", 10, WithIntParser(func(text string) (int, error) {
		intCalled = true
		if text != "custom-int" {
			t.Fatalf("int parser text = %q", text)
		}
		return 8080, nil
	})); got != 8080 || !intCalled {
		t.Fatalf("GetIntWithOptions = %d, called=%v", got, intCalled)
	}

	boolCalled := false
	if got := s.GetBoolWithOptions("debug", false, WithBoolParser(func(text string) (bool, error) {
		boolCalled = true
		if text != "custom-bool" {
			t.Fatalf("bool parser text = %q", text)
		}
		return true, nil
	})); !got || !boolCalled {
		t.Fatalf("GetBoolWithOptions = %v, called=%v", got, boolCalled)
	}
	if got := s.GetIntWithOptions("port", 10, WithIntParser(func(string) (int, error) {
		return 0, errors.New("invalid")
	})); got != 10 {
		t.Fatalf("GetIntWithOptions fallback = %d", got)
	}
	if got, err := s.GetIntEWithOptions("port", WithIntParser(func(string) (int, error) { return 7000, nil })); err != nil || got != 7000 {
		t.Fatalf("GetIntEWithOptions = %d, err=%v", got, err)
	}
	if got, err := s.GetBoolEWithOptions("debug", WithBoolParser(func(string) (bool, error) { return true, nil })); err != nil || !got {
		t.Fatalf("GetBoolEWithOptions = %v, err=%v", got, err)
	}
	if got, err := s.GetIntByGroupE("server", "port"); err != nil || got != 9090 {
		t.Fatalf("GetIntByGroupE = %d, err=%v", got, err)
	}
	if got, err := s.GetBoolByGroupE("server", "debug"); err != nil || !got {
		t.Fatalf("GetBoolByGroupE = %v, err=%v", got, err)
	}

	s.Set("bad-int", "abc")
	if _, err := s.GetIntE("missing"); !errors.Is(err, knifer.ErrCodeNotFound) {
		t.Fatalf("GetIntE missing err = %v, want not found", err)
	}
	if _, err := s.GetIntE("bad-int"); !errors.Is(err, knifer.ErrCodeInvalidInput) {
		t.Fatalf("GetIntE invalid err = %v, want invalid input", err)
	}
}

func TestBindWithOptionsUsesParsers(t *testing.T) {
	s := New()
	s.SetByGroup("server", "port", "custom-int")
	s.SetByGroup("server", "debug", "custom-bool")
	s.SetByGroup("server", "ratio", "custom-float")
	s.SetByGroup("server", "ids", "a,b")

	type serverConf struct {
		Port  int     `conf:"port"`
		Debug bool    `conf:"debug"`
		Ratio float64 `conf:"ratio"`
		IDs   []uint  `conf:"ids"`
	}
	var cfg serverConf
	var intCalled, boolCalled, floatCalled, uintCalled int
	err := s.BindGroupWithOptions("server", &cfg,
		WithBindIntParser(func(text string, base, bitSize int) (int64, error) {
			intCalled++
			if text == "custom-int" {
				return 8080, nil
			}
			return strconv.ParseInt(text, base, bitSize)
		}),
		WithBindBoolParser(func(text string) (bool, error) {
			boolCalled++
			return text == "custom-bool", nil
		}),
		WithBindFloatParser(func(text string, bitSize int) (float64, error) {
			floatCalled++
			if text == "custom-float" {
				return 0.75, nil
			}
			return strconv.ParseFloat(text, bitSize)
		}),
		WithBindUintParser(func(text string, base, bitSize int) (uint64, error) {
			uintCalled++
			switch text {
			case "a":
				return 1, nil
			case "b":
				return 2, nil
			default:
				return strconv.ParseUint(text, base, bitSize)
			}
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(cfg, serverConf{Port: 8080, Debug: true, Ratio: 0.75, IDs: []uint{1, 2}}) {
		t.Fatalf("BindGroupWithOptions = %#v", cfg)
	}
	if intCalled != 1 || boolCalled != 1 || floatCalled != 1 || uintCalled != 2 {
		t.Fatalf("parser calls int=%d bool=%d float=%d uint=%d", intCalled, boolCalled, floatCalled, uintCalled)
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
