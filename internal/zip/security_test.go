package zip

import (
	archivezip "archive/zip"
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	knifer "github.com/imajinyun/go-knifer"
)

func TestUnzipRejectsPathTraversal(t *testing.T) {
	tmp := t.TempDir()
	archive := filepath.Join(tmp, "bad.zip")
	var buf bytes.Buffer
	zw := archivezip.NewWriter(&buf)
	w, err := zw.Create("../evil.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write([]byte("bad")); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(archive, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	assertZipCode(t, UnzipTo(archive, filepath.Join(tmp, "dest")), knifer.ErrCodeInvalidInput)
}

func TestUnzipRejectsSymlinkEscape(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink behavior differs on windows")
	}

	tmp := t.TempDir()
	dest := filepath.Join(tmp, "dest")
	outside := filepath.Join(tmp, "outside")
	if err := os.MkdirAll(filepath.Join(dest, "link"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(dest, "link")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(dest, "link")); err != nil {
		t.Fatalf("create symlink: %v", err)
	}

	var buf bytes.Buffer
	zw := archivezip.NewWriter(&buf)
	w, err := zw.Create("link/evil.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write([]byte("bad")); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	r, err := archivezip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("NewReader: %v", err)
	}

	assertZipCode(t, UnzipReaderTo(r, dest), knifer.ErrCodeInvalidInput)
	if _, err := os.Stat(filepath.Join(outside, "evil.txt")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("symlink escape wrote outside file, stat err=%v", err)
	}
}
