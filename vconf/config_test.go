package vconf_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/imajinyun/go-knifer/vconf"
)

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

func TestFacadeTypedGettersMutationAndSchemaMethods(t *testing.T) {
	s, err := vconf.Parse(`
name=demo
custom_int=custom
custom_bool=yes
[server]
port=8080
enabled=true
`)
	if err != nil {
		t.Fatal(err)
	}
	if got, ok := s.Lookup("", "name"); !ok || got != "demo" {
		t.Fatalf("Lookup default = %q, %v", got, ok)
	}
	if got := s.GetIntWithOptions("custom_int", 0, vconf.WithIntParser(func(value string) (int, error) {
		if value == "custom" {
			return 77, nil
		}
		return 0, errors.New("unexpected")
	})); got != 77 {
		t.Fatalf("GetIntWithOptions = %d", got)
	}
	if got := s.GetBoolByGroupWithOptions("", "custom_bool", false, vconf.WithBoolParser(func(value string) (bool, error) {
		return value == "yes", nil
	})); !got {
		t.Fatal("GetBoolByGroupWithOptions = false")
	}
	if got := s.GetIntByGroupWithOptions("server", "port", 0); got != 8080 {
		t.Fatalf("GetIntByGroupWithOptions = %d", got)
	}
	if got := s.GetBoolByGroupWithOptions("server", "enabled", false); !got {
		t.Fatal("GetBoolByGroupWithOptions server.enabled = false")
	}

	clone := s.Clone()
	clone.Set("name", "clone")
	if s.Get("name") != "demo" || clone.Get("name") != "clone" {
		t.Fatalf("Clone should not alias source: source=%q clone=%q", s.Get("name"), clone.Get("name"))
	}
	clone.Delete("name")
	if _, ok := clone.Lookup("", "name"); ok {
		t.Fatal("Delete did not remove default key")
	}
	clone.DeleteByGroup("server", "enabled")
	if _, ok := clone.Lookup("server", "enabled"); ok {
		t.Fatal("DeleteByGroup did not remove grouped key")
	}
	merged := clone.Merge(vconf.Merge(s))
	if got := merged.Get("name"); got != "demo" {
		t.Fatalf("Merge method name = %q", got)
	}

	schema := vconf.Schema{Fields: []vconf.FieldRule{
		{Group: "server", Key: "host", Default: "127.0.0.1"},
		{Group: "server", Key: "port", Required: true, Type: vconf.TypeInt},
	}}
	withDefaults := s.ApplyDefaults(schema)
	if got := withDefaults.GetByGroup("server", "host"); got != "127.0.0.1" {
		t.Fatalf("ApplyDefaults host = %q", got)
	}
	if err := withDefaults.ValidateSchema(schema); err != nil {
		t.Fatalf("ValidateSchema: %v", err)
	}
	type defaultConfig struct {
		CustomInt int `conf:"custom_int,required,int"`
	}
	if err := withDefaults.ValidateStructWithOptions(defaultConfig{}, vconf.WithSchemaIntParser(func(value string, base int, bitSize int) (int64, error) {
		if value == "custom" {
			return 77, nil
		}
		return 0, errors.New("unexpected")
	})); err != nil {
		t.Fatalf("ValidateStructWithOptions: %v", err)
	}
}

func TestFacadeBindAndSchemaParserOptions(t *testing.T) {
	s, err := vconf.Parse(`
flag=yes
count=custom-int
amount=custom-uint
ratio=custom-float
items=1,2,3
schema_bool=yes
schema_float=custom-float
choice=blue
`)
	if err != nil {
		t.Fatal(err)
	}

	type bindConfig struct {
		Flag   bool    `conf:"flag"`
		Count  int     `conf:"count"`
		Amount uint    `conf:"amount"`
		Ratio  float64 `conf:"ratio"`
		Items  []int   `conf:"items"`
	}
	var cfg bindConfig
	if err := s.BindWithOptions(&cfg,
		vconf.WithBindBoolParser(func(value string) (bool, error) {
			return value == "yes", nil
		}),
		vconf.WithBindIntParser(func(value string, base int, bitSize int) (int64, error) {
			if value == "custom-int" {
				return 42, nil
			}
			return 7, nil
		}),
		vconf.WithBindUintParser(func(value string, base int, bitSize int) (uint64, error) {
			if value == "custom-uint" {
				return 9, nil
			}
			return 3, nil
		}),
		vconf.WithBindFloatParser(func(value string, bitSize int) (float64, error) {
			if value == "custom-float" {
				return 1.5, nil
			}
			return 0, errors.New("unexpected float")
		}),
	); err != nil {
		t.Fatalf("BindWithOptions() error = %v", err)
	}
	if !cfg.Flag || cfg.Count != 42 || cfg.Amount != 9 || cfg.Ratio != 1.5 || !reflect.DeepEqual(cfg.Items, []int{7, 7, 7}) {
		t.Fatalf("BindWithOptions cfg = %#v", cfg)
	}

	err = s.ValidateSchemaWithOptions(vconf.Schema{Fields: []vconf.FieldRule{
		{Key: "schema_bool", Required: true, Type: vconf.TypeBool},
		{Key: "schema_float", Required: true, Type: vconf.TypeFloat},
		{Key: "choice", Required: true, Choices: []string{"red", "blue"}},
	}},
		vconf.WithSchemaBoolParser(func(value string) (bool, error) {
			return value == "yes", nil
		}),
		vconf.WithSchemaFloatParser(func(value string, bitSize int) (float64, error) {
			if value == "custom-float" {
				return 2.5, nil
			}
			return 0, errors.New("unexpected schema float")
		}),
	)
	if err != nil {
		t.Fatalf("ValidateSchemaWithOptions() error = %v", err)
	}
}
