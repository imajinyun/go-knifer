package zip

import (
	archivezip "archive/zip"
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	knifer "github.com/imajinyun/go-knifer"
)

func TestZipEntriesAppendReadAndLimit(t *testing.T) {
	tmp := t.TempDir()
	archive := filepath.Join(tmp, "entries.zip")
	if err := ZipEntries(archive, EntryData{Name: "a.txt", Data: []byte("a")}); err != nil {
		t.Fatalf("ZipEntries: %v", err)
	}
	appendFile := filepath.Join(tmp, "b.txt")
	if err := os.WriteFile(appendFile, []byte("b"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Append(archive, appendFile); err != nil {
		t.Fatalf("Append: %v", err)
	}
	seen := map[string]bool{}
	if err := Read(archive, func(f *archivezip.File) error {
		seen[f.Name] = true
		return nil
	}); err != nil {
		t.Fatalf("Read: %v", err)
	}
	if !seen["a.txt"] || !seen["b.txt"] {
		t.Fatalf("seen: %#v", seen)
	}
	assertZipCode(t, UnzipToLimit(archive, filepath.Join(tmp, "limited"), 1), knifer.ErrCodeInvalidInput)
	_, err := Get(archive, "missing.txt")
	assertZipCode(t, err, knifer.ErrCodeNotFound)
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("missing entry error should preserve os.ErrNotExist: %v", err)
	}
}

func TestGetBytesEnforcesConfiguredLimit(t *testing.T) {
	tmp := t.TempDir()
	archive := filepath.Join(tmp, "limit.zip")
	data := bytes.Repeat([]byte("x"), 16)
	if err := ZipEntries(archive, EntryData{Name: "a.txt", Data: data}); err != nil {
		t.Fatalf("ZipEntries: %v", err)
	}

	if got := applyDecompressOptions(nil).maxBytes; got != DefaultUnzipMaxBytes {
		t.Fatalf("default entry read max bytes = %d, want %d", got, DefaultUnzipMaxBytes)
	}
	if _, err := GetBytesWithOptions(archive, "a.txt", WithMaxBytes(8)); err == nil {
		t.Fatal("GetBytes should enforce configured entry read limit")
	}
	if out, err := GetBytesWithOptions(archive, "a.txt", WithMaxBytes(0)); err != nil || !bytes.Equal(out, data) {
		t.Fatalf("GetBytes explicit unlimited = %q, %v", out, err)
	}
}

func TestUnzipDefaultLimitCanBeOverridden(t *testing.T) {
	tmp := t.TempDir()
	archive := filepath.Join(tmp, "limit.zip")
	entries := []EntryData{{Name: "a.txt", Data: []byte("abcd")}}
	if err := ZipEntries(archive, entries...); err != nil {
		t.Fatalf("ZipEntries: %v", err)
	}
	if err := UnzipToWithOptions(archive, filepath.Join(tmp, "limited"), WithMaxBytes(3)); err == nil {
		t.Fatal("UnzipToWithOptions should reject archives larger than max bytes")
	}
	if err := UnzipToLimit(archive, filepath.Join(tmp, "unlimited"), -1); err != nil {
		t.Fatalf("UnzipToLimit with explicit unlimited limit: %v", err)
	}
	if got := applyUnzipOptions(nil).maxBytes; got != DefaultUnzipMaxBytes {
		t.Fatalf("default unzip max bytes = %d, want %d", got, DefaultUnzipMaxBytes)
	}
	if got := applyUnzipOptions([]ArchiveOption{WithMaxBytes(-1)}).maxBytes; got != -1 {
		t.Fatalf("explicit unlimited unzip max bytes = %d, want -1", got)
	}
}

func TestUnzipEnforcesActualCopiedBytes(t *testing.T) {
	var buf bytes.Buffer
	if err := ZipEntriesToWriter(&buf, EntryData{Name: "a.txt", Data: []byte("abcd")}); err != nil {
		t.Fatalf("ZipEntriesToWriter: %v", err)
	}
	r, err := archivezip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("NewReader: %v", err)
	}
	if len(r.File) != 1 {
		t.Fatalf("archive entries = %d, want 1", len(r.File))
	}
	// Simulate an archive whose declared uncompressed size is smaller than the
	// actual stream. Extraction must enforce the copy-time limit as a second line
	// of defense instead of trusting central-directory metadata only.
	r.File[0].UncompressedSize64 = 1
	dest := filepath.Join(t.TempDir(), "dest")
	if err := UnzipReaderToWithOptions(r, dest, WithMaxBytes(3)); err == nil {
		t.Fatal("UnzipReaderToWithOptions should reject streams exceeding the actual copy limit")
	}
	data, err := os.ReadFile(filepath.Join(dest, "a.txt"))
	if err != nil {
		t.Fatalf("read partial extraction: %v", err)
	}
	if len(data) > 3 {
		t.Fatalf("partial extraction wrote %d bytes, want at most 3", len(data))
	}
}
