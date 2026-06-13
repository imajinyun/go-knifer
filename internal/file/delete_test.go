package file

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestMkdirTouchDelDelete(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "x", "y")
	if err := Mkdir(sub); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	f := filepath.Join(sub, "a.txt")
	if err := Touch(f); err != nil {
		t.Fatalf("touch: %v", err)
	}
	if err := Del(f); err != nil {
		t.Fatalf("del: %v", err)
	}
	if FileExists(f) {
		t.Fatalf("Del did not remove")
	}
}

func TestDeleteProviderOptions(t *testing.T) {
	regular := fakeFileInfo{name: "file.txt", size: 42}
	stat := func(path string) (fs.FileInfo, error) {
		switch path {
		case "file.txt":
			return regular, nil
		default:
			return nil, os.ErrNotExist
		}
	}
	removed := ""
	if err := DelWithOptions("file.txt", WithStat(stat), WithRemoveAll(func(path string) error {
		removed = path
		return nil
	})); err != nil {
		t.Fatalf("DelWithOptions() error = %v", err)
	}
	if removed != "file.txt" {
		t.Fatalf("removed path = %q, want file.txt", removed)
	}
	removed = ""
	if err := DelWithOptions("missing.txt", WithStat(stat), WithRemoveAll(func(path string) error {
		removed = path
		return nil
	})); err != nil || removed != "" {
		t.Fatalf("DelWithOptions() missing = %v, removed %q", err, removed)
	}
}
