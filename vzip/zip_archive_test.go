package vzip_test

import (
	archivezip "archive/zip"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/imajinyun/go-knifer/vzip"
)

func TestFacadeZipCreationUsingOptions(t *testing.T) {
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
	if err := vzip.ZipFilesUsingOptions(archive, []string{src}, vzip.WithSourceDir(true), vzip.WithFileFilter(filter)); err != nil {
		t.Fatalf("ZipFilesUsingOptions: %v", err)
	}
	data, err := vzip.GetBytes(archive, "src/keep.txt")
	if err != nil || string(data) != "keep" {
		t.Fatalf("GetBytes keep = %q, %v", data, err)
	}
	if _, err := vzip.GetBytes(archive, "src/skip.log"); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("skip.log err = %v, want not exist", err)
	}

	var buf bytes.Buffer
	if err := vzip.ZipToWriterUsingOptions(&buf, []string{src}, vzip.WithFileFilter(filter)); err != nil {
		t.Fatalf("ZipToWriterUsingOptions: %v", err)
	}
	bufReader, err := archivezip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("NewReader: %v", err)
	}
	if len(bufReader.File) != 1 || bufReader.File[0].Name != "keep.txt" {
		t.Fatalf("writer archive entries = %#v", bufReader.File)
	}
	entry, err := vzip.GetStream(bufReader.File[0])
	if err != nil {
		t.Fatalf("GetStream: %v", err)
	}
	defer func() { _ = entry.Close() }()
	if _, err := io.ReadAll(entry); err != nil {
		t.Fatalf("read entry: %v", err)
	}
}

func TestFacadeZipDefaultFileAndWriterHelpers(t *testing.T) {
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

	autoArchive, err := vzip.Zip(filepath.Join(src, "keep.txt"))
	if err != nil {
		t.Fatalf("Zip: %v", err)
	}
	if got, err := vzip.GetBytes(autoArchive, "keep.txt"); err != nil || string(got) != "keep" {
		t.Fatalf("Zip content = %q, %v", got, err)
	}

	toArchive := filepath.Join(tmp, "to.zip")
	if err := vzip.ZipTo(src, toArchive, true); err != nil {
		t.Fatalf("ZipTo: %v", err)
	}
	if got, err := vzip.GetBytes(toArchive, "src/keep.txt"); err != nil || string(got) != "keep" {
		t.Fatalf("ZipTo content = %q, %v", got, err)
	}

	filesArchive := filepath.Join(tmp, "files.zip")
	if err := vzip.ZipFiles(filesArchive, false, filepath.Join(src, "keep.txt")); err != nil {
		t.Fatalf("ZipFiles: %v", err)
	}
	if got, err := vzip.GetBytes(filesArchive, "keep.txt"); err != nil || string(got) != "keep" {
		t.Fatalf("ZipFiles content = %q, %v", got, err)
	}

	filterArchive := filepath.Join(tmp, "filter.zip")
	filter := func(path string, info os.FileInfo) bool {
		return info.IsDir() || filepath.Ext(path) == ".txt"
	}
	if err := vzip.ZipFilesFilter(filterArchive, false, filter, src); err != nil {
		t.Fatalf("ZipFilesFilter: %v", err)
	}
	if _, err := vzip.GetBytes(filterArchive, "skip.log"); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("ZipFilesFilter skip err = %v, want not exist", err)
	}

	var buf bytes.Buffer
	if err := vzip.ZipToWriter(&buf, false, filter, src); err != nil {
		t.Fatalf("ZipToWriter: %v", err)
	}
	zr, err := archivezip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil || len(zr.File) != 1 || zr.File[0].Name != "keep.txt" {
		t.Fatalf("ZipToWriter archive = %#v, %v", zr, err)
	}
}
