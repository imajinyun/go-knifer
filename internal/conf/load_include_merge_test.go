package conf

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadWithOptionsIncludesAndMerges(t *testing.T) {
	dir := t.TempDir()
	common := filepath.Join(dir, "common.setting")
	main := filepath.Join(dir, "main.setting")
	if err := os.WriteFile(common, []byte("name=common\n[server]\nhost=127.0.0.1"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(main, []byte("include=common.setting\nname=main\n[server]\nport=8080"), 0o644); err != nil {
		t.Fatal(err)
	}

	c, err := LoadWithOptions(main, LoadOptions{AllowInclude: true})
	if err != nil {
		t.Fatal(err)
	}
	if got := c.Get("name"); got != "main" {
		t.Fatalf("merged name = %q", got)
	}
	if got := c.GetByGroup("server", "host"); got != "127.0.0.1" {
		t.Fatalf("included server.host = %q", got)
	}
	if _, ok := c.Lookup("", "include"); ok {
		t.Fatal("include key should be removed after loading")
	}
}
