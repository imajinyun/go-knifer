package conf

import (
	"reflect"
	"testing"
)

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

func TestParseYAMLFullWithOptionsUsesProvider(t *testing.T) {
	called := false
	s, err := ParseYAMLFullWithOptions("ignored", WithYAMLUnmarshalFunc(func(data []byte, out any) error {
		called = true
		root, ok := out.(*any)
		if !ok {
			t.Fatalf("unmarshal output = %T, want *any", out)
		}
		*root = map[string]any{"app": map[string]any{"name": "provider"}}
		return nil
	}))
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("custom YAML unmarshal provider was not called")
	}
	if got := s.GetByGroup("app", "name"); got != "provider" {
		t.Fatalf("provider app.name = %q", got)
	}
}
