package file

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMkdirTouchDelDirectory(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "x", "y")
	if err := Mkdir(sub); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if !IsDirectory(sub) {
		t.Fatalf("IsDirectory failed")
	}
}

func TestMkdirAndCopyOptionsDirectory(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "private")
	if err := Mkdir(sub, WithMkdirPerm(0o700)); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	info, err := os.Stat(sub)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o700 {
		t.Fatalf("dir perm = %o, want 700", got)
	}
}
