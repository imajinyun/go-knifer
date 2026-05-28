package vextra_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/imajinyun/go-knifer/vextra"
)

func TestCompressFacade(t *testing.T) {
	data := []byte("hello hutool extra facade")
	gz, err := vextra.Gzip(data)
	if err != nil {
		t.Fatal(err)
	}
	out, err := vextra.Gunzip(gz)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out, data) {
		t.Fatalf("Gunzip() = %q", out)
	}
	z, err := vextra.Zlib(data)
	if err != nil {
		t.Fatal(err)
	}
	out, err = vextra.Unzlib(z)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out, data) {
		t.Fatalf("Unzlib() = %q", out)
	}
}

func TestZipFacade(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src")
	if err := os.Mkdir(src, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "a.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	archive := filepath.Join(dir, "out.zip")
	if err := vextra.ZipFiles(archive, src); err != nil {
		t.Fatal(err)
	}
	dest := filepath.Join(dir, "dest")
	if err := vextra.Unzip(archive, dest); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(dest, "src", "a.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "hello" {
		t.Fatalf("unzipped content = %q", b)
	}
}

func TestEmojiTemplateAndValidationFacade(t *testing.T) {
	if !vextra.ContainsEmoji("hello😀") {
		t.Fatal("ContainsEmoji() = false")
	}
	if got := vextra.RemoveEmoji("hello😀"); got != "hello" {
		t.Fatalf("RemoveEmoji() = %q", got)
	}
	out, err := vextra.RenderTemplate("hi {{.Name}}", map[string]string{"Name": "gokit"})
	if err != nil {
		t.Fatal(err)
	}
	if out != "hi gokit" {
		t.Fatalf("RenderTemplate() = %q", out)
	}
	if !vextra.IsEmail("a@example.com") || vextra.IsEmail("bad@@example") {
		t.Fatal("IsEmail returned unexpected result")
	}
	if !vextra.IsURL("https://example.com/a") || vextra.IsURL("/relative/path") {
		t.Fatal("IsURL returned unexpected result")
	}
	if !vextra.IsBlank(" \n\t") || vextra.RuneLen("你好") != 2 {
		t.Fatal("string helpers returned unexpected result")
	}
}
