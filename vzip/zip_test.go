package vzip_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	knifer "github.com/imajinyun/go-knifer"
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

func TestFacadeZipErrorContract(t *testing.T) {
	_, err := vzip.GetStream(nil)
	if err == nil {
		t.Fatal("GetStream(nil) error = nil, want invalid input")
	}
	if !errors.Is(err, knifer.ErrCodeInvalidInput) {
		t.Fatalf("errors.Is(err, ErrCodeInvalidInput) = false: %v", err)
	}
	code, ok := knifer.CodeOf(err)
	if !ok || code != knifer.ErrCodeInvalidInput {
		t.Fatalf("CodeOf(err) = %q, %v; want invalid input", code, ok)
	}
	var zipErr *vzip.Error
	if !errors.As(err, &zipErr) {
		t.Fatalf("errors.As(err, *vzip.Error) = false: %v", err)
	}
}
