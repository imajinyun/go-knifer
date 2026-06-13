package zip

import (
	"compress/flate"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestAppendWithOptions(t *testing.T) {
	tmp := t.TempDir()
	archive := filepath.Join(tmp, "entries.zip")
	if err := ZipEntries(archive, EntryData{Name: "a.txt", Data: []byte("a")}); err != nil {
		t.Fatalf("ZipEntries: %v", err)
	}
	keepFile := filepath.Join(tmp, "keep.txt")
	if err := os.WriteFile(keepFile, []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	skipFile := filepath.Join(tmp, "skip.log")
	if err := os.WriteFile(skipFile, []byte("skip"), 0o644); err != nil {
		t.Fatal(err)
	}
	filter := func(path string, info os.FileInfo) bool {
		return info.IsDir() || filepath.Ext(path) == ".txt"
	}
	if err := AppendWithOptions(archive, keepFile, WithFileFilter(filter), WithCompressionLevel(flate.BestSpeed)); err != nil {
		t.Fatalf("AppendWithOptions keep: %v", err)
	}
	if err := AppendWithOptions(archive, skipFile, WithFileFilter(filter)); err != nil {
		t.Fatalf("AppendWithOptions skip: %v", err)
	}
	data, err := GetBytes(archive, "keep.txt")
	if err != nil || string(data) != "keep" {
		t.Fatalf("appended keep.txt = %q, %v", data, err)
	}
	if _, err := GetBytes(archive, "skip.log"); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("filtered skip.log err = %v, want not exist", err)
	}
}
