package zip

import (
	archivezip "archive/zip"
	"bytes"
	"compress/flate"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestArchiveOptions(t *testing.T) {
	tmp := t.TempDir()
	archive := filepath.Join(tmp, "entries.zip")
	entries := []EntryData{{Name: "a.txt", Data: []byte("abcd")}}
	if err := ZipEntriesWithOptions(archive, entries, WithFilePerm(0o600), WithCompressionLevel(flate.BestSpeed)); err != nil {
		t.Fatalf("ZipEntriesWithOptions: %v", err)
	}
	info, err := os.Stat(archive)
	if err != nil {
		t.Fatalf("stat archive: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("archive perm = %o, want 600", got)
	}
	if err := ZipEntriesWithOptions(archive, entries, WithOverwrite(false)); !errors.Is(err, os.ErrExist) {
		t.Fatalf("overwrite false err = %v, want exists", err)
	}
	if _, err := GetBytesWithOptions(archive, "a.txt", WithMaxBytes(3)); err == nil {
		t.Fatal("GetBytesWithOptions over limit error = nil")
	}
	dest := filepath.Join(tmp, "dest")
	if err := UnzipToWithOptions(archive, dest, WithDirPerm(0o700), WithFilePerm(0o600), WithPreserveMode(false)); err != nil {
		t.Fatalf("UnzipToWithOptions: %v", err)
	}
	fileInfo, err := os.Stat(filepath.Join(dest, "a.txt"))
	if err != nil {
		t.Fatalf("stat extracted: %v", err)
	}
	if got := fileInfo.Mode().Perm(); got != 0o600 {
		t.Fatalf("extracted perm = %o, want 600", got)
	}
	if err := UnzipToWithOptions(archive, dest, WithOverwrite(false)); !errors.Is(err, os.ErrExist) {
		t.Fatalf("unzip overwrite false err = %v, want exists", err)
	}
}

func TestArchiveProviderOptionsForZipEntries(t *testing.T) {
	var mkdirPath string
	var mkdirPerm os.FileMode
	var openPath string
	var openFlag int
	var openPerm os.FileMode
	var buf bytes.Buffer
	closer := &zipBufferWriteCloser{Buffer: &buf}

	err := ZipEntriesWithOptions("parent/archive.zip", []EntryData{{Name: "a.txt", Data: []byte("a")}},
		WithDirPerm(0o700),
		WithFilePerm(0o600),
		WithMkdirAll(func(path string, perm os.FileMode) error {
			mkdirPath = path
			mkdirPerm = perm
			return nil
		}),
		WithOpenFile(func(path string, flag int, perm os.FileMode) (io.WriteCloser, error) {
			openPath = path
			openFlag = flag
			openPerm = perm
			return closer, nil
		}),
	)
	if err != nil {
		t.Fatalf("ZipEntriesWithOptions() error = %v", err)
	}
	if mkdirPath != "parent" || mkdirPerm != 0o700 {
		t.Fatalf("mkdir = %q/%o, want parent/700", mkdirPath, mkdirPerm)
	}
	if openPath != "parent/archive.zip" || openPerm != 0o600 || openFlag&os.O_TRUNC == 0 {
		t.Fatalf("open = %q/%o/%#x", openPath, openPerm, openFlag)
	}
	if !closer.closed {
		t.Fatal("archive output was not closed")
	}
	r, err := archivezip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("NewReader(provider output): %v", err)
	}
	if len(r.File) != 1 || r.File[0].Name != "a.txt" {
		t.Fatalf("entries = %#v, want a.txt", r.File)
	}
}

func TestArchiveProviderOptionsForFileCompression(t *testing.T) {
	data := []byte("provider-data")
	openPath := ""
	statPath := ""
	open := func(path string) (io.ReadCloser, error) {
		openPath = path
		return io.NopCloser(bytes.NewReader(data)), nil
	}
	stat := func(path string) (os.FileInfo, error) {
		statPath = path
		return zipFakeFileInfo{name: path, size: int64(len(data))}, nil
	}
	gz, err := GzipFileWithOptions("virtual.txt", WithOpen(open), WithStat(stat), WithCompressionLevel(flate.BestSpeed))
	if err != nil {
		t.Fatalf("GzipFileWithOptions() error = %v", err)
	}
	if openPath != "virtual.txt" || statPath != "virtual.txt" {
		t.Fatalf("provider paths open=%q stat=%q", openPath, statPath)
	}
	out, err := UnGzip(gz)
	if err != nil || !bytes.Equal(out, data) {
		t.Fatalf("UnGzip(provider gzip) = %q, %v", out, err)
	}
	z, err := ZlibFileWithOptions("virtual.txt", flate.BestSpeed, WithOpen(open), WithStat(stat))
	if err != nil {
		t.Fatalf("ZlibFileWithOptions() error = %v", err)
	}
	out, err = UnZlib(z)
	if err != nil || !bytes.Equal(out, data) {
		t.Fatalf("UnZlib(provider zlib) = %q, %v", out, err)
	}
}

func TestArchiveProviderOptionsForReadAndExtract(t *testing.T) {
	tmp := t.TempDir()
	archive := filepath.Join(tmp, "entries.zip")
	if err := ZipEntries(archive, EntryData{Name: "a.txt", Data: []byte("a")}); err != nil {
		t.Fatalf("ZipEntries() error = %v", err)
	}
	opened := ""
	openZip := func(path string) (*archivezip.ReadCloser, error) {
		opened = path
		return archivezip.OpenReader(path)
	}
	data, err := GetBytesWithOptions(archive, "a.txt", WithOpenZipReader(openZip))
	if err != nil || string(data) != "a" || opened != archive {
		t.Fatalf("GetBytesWithOptions() = %q, %v, opened=%q", data, err, opened)
	}
	names, err := ListFileNamesWithOptions(archive, "", WithOpenZipReader(openZip))
	if err != nil || !reflect.DeepEqual(names, []string{"a.txt"}) {
		t.Fatalf("ListFileNamesWithOptions() = %v, %v", names, err)
	}
	seen := false
	if err := ReadWithOptions(archive, func(f *archivezip.File) error {
		seen = f.Name == "a.txt"
		return nil
	}, WithOpenZipReader(openZip)); err != nil || !seen {
		t.Fatalf("ReadWithOptions() = %v, seen=%v", err, seen)
	}

	r, err := archivezip.OpenReader(archive)
	if err != nil {
		t.Fatalf("OpenReader() error = %v", err)
	}
	defer func() { _ = r.Close() }()
	var extracted bytes.Buffer
	var mkdirs []string
	if err := UnzipReaderToWithOptions(&r.Reader, "dest",
		WithMkdirAll(func(path string, perm os.FileMode) error {
			mkdirs = append(mkdirs, path)
			return nil
		}),
		WithEvalSymlinks(func(path string) (string, error) {
			return path, nil
		}),
		WithOpenFile(func(path string, flag int, perm os.FileMode) (io.WriteCloser, error) {
			return &zipBufferWriteCloser{Buffer: &extracted}, nil
		}),
	); err != nil {
		t.Fatalf("UnzipReaderToWithOptions() error = %v", err)
	}
	if extracted.String() != "a" || len(mkdirs) == 0 {
		t.Fatalf("extracted=%q mkdirs=%v", extracted.String(), mkdirs)
	}
}
