package vzip_test

import (
	"os"
	"path/filepath"
	"testing"
)

func newZipArchiveSource(t *testing.T) (string, string) {
	t.Helper()
	tmp := t.TempDir()
	src := filepath.Join(tmp, "src")
	if err := os.MkdirAll(src, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "keep.txt"), []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "skip.log"), []byte("skip"), 0o644); err != nil {
		t.Fatal(err)
	}
	return tmp, src
}

func zipArchiveTextFilter(path string, info os.FileInfo) bool {
	return info.IsDir() || filepath.Ext(path) == ".txt"
}
