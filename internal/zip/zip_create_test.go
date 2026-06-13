package zip

import (
	archivezip "archive/zip"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestZipFilesUnzipGetAndList(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "src")
	if err := os.MkdirAll(filepath.Join(src, "nested"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(src, "empty"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "a.txt"), []byte("a"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "nested", "b.txt"), []byte("b"), 0o644); err != nil {
		t.Fatal(err)
	}
	archive := filepath.Join(tmp, "out.zip")
	if err := ZipFiles(archive, false, src); err != nil {
		t.Fatalf("ZipFiles: %v", err)
	}
	data, err := GetBytes(archive, "a.txt")
	if err != nil || string(data) != "a" {
		t.Fatalf("GetBytes: %q %v", data, err)
	}
	names, err := ListFileNames(archive, "")
	if err != nil {
		t.Fatalf("ListFileNames: %v", err)
	}
	sort.Strings(names)
	if !reflect.DeepEqual(names, []string{"a.txt"}) {
		t.Fatalf("names: %#v", names)
	}
	dest := filepath.Join(tmp, "dest")
	if err := UnzipTo(archive, dest); err != nil {
		t.Fatalf("UnzipTo: %v", err)
	}
	if got, err := os.ReadFile(filepath.Join(dest, "nested", "b.txt")); err != nil || string(got) != "b" {
		t.Fatalf("unzipped: %q %v", got, err)
	}
	if info, err := os.Stat(filepath.Join(dest, "empty")); err != nil || !info.IsDir() {
		t.Fatalf("empty directory was not restored, info=%v err=%v", info, err)
	}
}

func TestZipStreamsEnforcesInputLimit(t *testing.T) {
	var buf bytes.Buffer
	err := ZipStreamsToWriterWithOptions(&buf, []StreamEntry{{Name: "a.txt", Reader: bytes.NewReader([]byte("abcd"))}}, WithMaxBytes(3))
	if err == nil {
		t.Fatal("ZipStreamsToWriterWithOptions should reject stream input over max bytes")
	}
}

func TestZipCreationUsingOptions(t *testing.T) {
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

	archive := filepath.Join(tmp, "filtered.zip")
	filter := func(path string, info os.FileInfo) bool {
		return info.IsDir() || filepath.Ext(path) == ".txt"
	}
	if err := ZipFilesWithOptions(archive, []string{src}, WithSourceDir(true), WithFileFilter(filter)); err != nil {
		t.Fatalf("ZipFilesWithOptions: %v", err)
	}
	data, err := GetBytes(archive, "src/keep.txt")
	if err != nil || string(data) != "keep" {
		t.Fatalf("GetBytes keep = %q, %v", data, err)
	}
	if _, err := GetBytes(archive, "src/skip.log"); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("skip.log err = %v, want not exist", err)
	}

	var buf bytes.Buffer
	if err := ZipToWriterWithOptions(&buf, []string{src}, WithFileFilter(filter)); err != nil {
		t.Fatalf("ZipToWriterWithOptions: %v", err)
	}
	r, err := archivezip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("NewReader: %v", err)
	}
	for _, f := range r.File {
		if f.Name == "skip.log" {
			t.Fatal("ZipToWriterWithOptions should filter skip.log")
		}
		if f.Name == "keep.txt" {
			rc, err := f.Open()
			if err != nil {
				t.Fatalf("open keep.txt: %v", err)
			}
			content, err := io.ReadAll(rc)
			_ = rc.Close()
			if err != nil || string(content) != "keep" {
				t.Fatalf("keep.txt = %q, %v", content, err)
			}
			return
		}
	}
	t.Fatal("keep.txt not found in writer archive")
}
