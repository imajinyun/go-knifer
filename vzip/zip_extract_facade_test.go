package vzip_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/imajinyun/go-knifer/vzip"
)

func TestFacadeZipAndExtraction(t *testing.T) {
	tmp := t.TempDir()
	archive := filepath.Join(tmp, "data.zip")
	if err := vzip.ZipEntries(archive, vzip.EntryData{Name: "hello.txt", Data: []byte("hello")}); err != nil {
		t.Fatalf("ZipEntries: %v", err)
	}
	data, err := vzip.GetBytes(archive, "hello.txt")
	if err != nil || string(data) != "hello" {
		t.Fatalf("GetBytes: %q %v", data, err)
	}
	dest := filepath.Join(tmp, "dest")
	if err := vzip.UnzipTo(archive, dest); err != nil {
		t.Fatalf("UnzipTo: %v", err)
	}
	if got, err := os.ReadFile(filepath.Join(dest, "hello.txt")); err != nil || string(got) != "hello" {
		t.Fatalf("unzipped: %q %v", got, err)
	}
}

func TestFacadeZipAppendAndUnzipOptions(t *testing.T) {
	tmp := t.TempDir()
	archive := filepath.Join(tmp, "append-default.zip")
	if err := vzip.ZipEntries(archive, vzip.EntryData{Name: "base.txt", Data: []byte("base")}); err != nil {
		t.Fatalf("ZipEntries: %v", err)
	}
	extra := filepath.Join(tmp, "extra.txt")
	if err := os.WriteFile(extra, []byte("extra"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := vzip.Append(archive, extra); err != nil {
		t.Fatalf("Append: %v", err)
	}
	if got, err := vzip.GetBytes(archive, "extra.txt"); err != nil || string(got) != "extra" {
		t.Fatalf("Append content = %q, %v", got, err)
	}

	defaultDest, err := vzip.Unzip(archive)
	if err != nil {
		t.Fatalf("Unzip: %v", err)
	}
	if got, err := os.ReadFile(filepath.Join(defaultDest, "base.txt")); err != nil || string(got) != "base" {
		t.Fatalf("Unzip default output = %q, %v", got, err)
	}
	if err := vzip.UnzipToLimit(archive, filepath.Join(tmp, "limit"), 1); err == nil {
		t.Fatal("UnzipToLimit should reject content larger than limit")
	}
}
