package vzip_test

import (
	"bytes"
	"compress/flate"
	"os"
	"path/filepath"
	"testing"

	"github.com/imajinyun/go-knifer/vzip"
)

func TestFacadeZipAndCompression(t *testing.T) {
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
	gz, err := vzip.GzipString("hello")
	if err != nil {
		t.Fatalf("GzipString: %v", err)
	}
	text, err := vzip.UnGzipString(gz)
	if err != nil || text != "hello" {
		t.Fatalf("UnGzipString: %q %v", text, err)
	}
	dataBytes := []byte("hello the utility toolkit zip facade")
	gzipBytes, err := vzip.Gzip(dataBytes)
	if err != nil {
		t.Fatalf("Gzip: %v", err)
	}
	out, err := vzip.Gunzip(gzipBytes)
	if err != nil || !bytes.Equal(out, dataBytes) {
		t.Fatalf("Gunzip: %q %v", out, err)
	}
	zlibBytes, err := vzip.Zlib(dataBytes)
	if err != nil {
		t.Fatalf("Zlib: %v", err)
	}
	out, err = vzip.Unzlib(zlibBytes)
	if err != nil || !bytes.Equal(out, dataBytes) {
		t.Fatalf("Unzlib: %q %v", out, err)
	}
}

func TestFacadeZipAppendUnzipAndCompressionOptions(t *testing.T) {
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

	payload := []byte("compression option payload")
	source := filepath.Join(tmp, "payload.txt")
	if err := os.WriteFile(source, payload, 0o644); err != nil {
		t.Fatal(err)
	}
	gz, err := vzip.GzipFileWithOptions(source, vzip.WithMaxBytes(int64(len(payload))))
	if err != nil {
		t.Fatalf("GzipFileWithOptions: %v", err)
	}
	if out, err := vzip.UnGzipReaderWithOptions(bytes.NewReader(gz), len(payload), vzip.WithMaxBytes(int64(len(payload)))); err != nil || !bytes.Equal(out, payload) {
		t.Fatalf("UnGzipReaderWithOptions = %q, %v", out, err)
	}
	if _, err := vzip.GzipFileWithOptions(source, vzip.WithMaxBytes(1)); err == nil {
		t.Fatal("GzipFileWithOptions max bytes error = nil")
	}

	zlibBytes, err := vzip.ZlibFileWithOptions(source, flate.BestSpeed, vzip.WithMaxBytes(int64(len(payload))))
	if err != nil {
		t.Fatalf("ZlibFileWithOptions: %v", err)
	}
	if out, err := vzip.UnZlibReaderWithOptions(bytes.NewReader(zlibBytes), len(payload), vzip.WithMaxBytes(int64(len(payload)))); err != nil || !bytes.Equal(out, payload) {
		t.Fatalf("UnZlibReaderWithOptions = %q, %v", out, err)
	}
	if got, err := vzip.ZlibLevelWithOptions(payload, flate.NoCompression, vzip.WithMaxBytes(int64(len(payload)))); err != nil || len(got) == 0 {
		t.Fatalf("ZlibLevelWithOptions len=%d err=%v", len(got), err)
	}
	if got, err := vzip.ZlibReaderWithOptions(bytes.NewReader(payload), flate.BestSpeed, len(payload), vzip.WithMaxBytes(int64(len(payload)))); err != nil || len(got) == 0 {
		t.Fatalf("ZlibReaderWithOptions len=%d err=%v", len(got), err)
	}
}
