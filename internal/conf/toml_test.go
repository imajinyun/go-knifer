package conf

import "testing"

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

func TestParseTOMLNestedDottedKeys(t *testing.T) {
	c, err := ParseTOML(`
title = "demo"
[database]
ports = [8000, 8001, 8002]
enabled = true
connection.max = 5000
[servers.alpha]
ip = "10.0.0.1"
`)
	if err != nil {
		t.Fatal(err)
	}
	if got := c.GetByGroup("database", "ports"); got != "8000,8001,8002" {
		t.Fatalf("database.ports = %q", got)
	}
	if got := c.GetByGroup("database.connection", "max"); got != "5000" {
		t.Fatalf("database.connection.max = %q", got)
	}
	if got := c.GetByGroup("servers.alpha", "ip"); got != "10.0.0.1" {
		t.Fatalf("servers.alpha.ip = %q", got)
	}
}

func TestParseTOMLWithOptionsUsesProvider(t *testing.T) {
	called := false
	c, err := ParseTOMLWithOptions("ignored", WithTOMLUnmarshalFunc(func(data []byte, out any) error {
		called = true
		root, ok := out.(*map[string]any)
		if !ok {
			t.Fatalf("toml unmarshal output = %T, want *map[string]any", out)
		}
		*root = map[string]any{"app": map[string]any{"name": "provider"}}
		return nil
	}))
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("custom TOML unmarshal provider was not called")
	}
	if got := c.GetByGroup("app", "name"); got != "provider" {
		t.Fatalf("provider app.name = %q", got)
	}
}
