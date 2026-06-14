package conf

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadProfile(t *testing.T) {
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
}

func TestApplyProfileFromYAML(t *testing.T) {
	profileYAML, err := ParseYAMLFull(`
app:
  name: base
server:
  port: 8080
profile:
  dev:
    app:
      name: dev
    server:
      port: 9090
`)
	if err != nil {
		t.Fatal(err)
	}
	dev := profileYAML.ApplyProfile("dev")
	if got := dev.GetByGroup("app", "name"); got != "dev" {
		t.Fatalf("YAML profile app.name = %q", got)
	}
	if got := dev.GetByGroup("server", "port"); got != "9090" {
		t.Fatalf("YAML profile server.port = %q", got)
	}
}
