package conf

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFilesMergesInOrder(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "base.setting")
	override := filepath.Join(dir, "override.toml")
	if err := os.WriteFile(base, []byte("name=base\nmode=dev"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(override, []byte("name='override'\n[server]\nport=9090"), 0o644); err != nil {
		t.Fatal(err)
	}
	c, err := LoadFiles(base, override)
	if err != nil {
		t.Fatal(err)
	}
	if got := c.Get("name"); got != "override" {
		t.Fatalf("LoadFiles merged name = %q", got)
	}
}
