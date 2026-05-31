package zip

import (
	archivezip "archive/zip"
	"bytes"
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
	if err := UnzipToLimit(archive, filepath.Join(tmp, "limited"), 1); err == nil {
		t.Fatal("expected size limit error")
	}
}

func TestGzipAndZlib(t *testing.T) {
	data := []byte("hello compression")
	gz, err := Gzip(data)
	if err != nil {
		t.Fatalf("Gzip: %v", err)
	}
	out, err := UnGzip(gz)
	if err != nil || !bytes.Equal(out, data) {
		t.Fatalf("UnGzip: %q %v", out, err)
	}
	z, err := ZlibLevel(data, 6)
	if err != nil {
		t.Fatalf("ZlibLevel: %v", err)
	}
	out, err = UnZlib(z)
	if err != nil || !bytes.Equal(out, data) {
		t.Fatalf("UnZlib: %q %v", out, err)
	}
}

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
	if err := UnzipTo(archive, filepath.Join(tmp, "dest")); err == nil {
		t.Fatal("expected path traversal error")
	}
}
