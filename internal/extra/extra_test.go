package extra

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestGzipAndZlib(t *testing.T) {
	data := []byte("hello hutool extra")
	gz, err := Gzip(data)
	if err != nil {
		t.Fatal(err)
	}
	out, err := Gunzip(gz)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out, data) {
		t.Fatalf("Gunzip() = %q", out)
	}
	z, err := Zlib(data)
	if err != nil {
		t.Fatal(err)
	}
	out, err = Unzlib(z)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out, data) {
		t.Fatalf("Unzlib() = %q", out)
	}
}

func TestEmojiTemplateValidation(t *testing.T) {
	if !ContainsEmoji("hi😀") || !ContainsEmoji("go❤️") || !ContainsEmoji("1️⃣") {
		t.Fatal("ContainsEmoji() = false")
	}
	if got := RemoveEmoji("hi😀 go❤️ 1️⃣"); got != "hi go " {
		t.Fatalf("RemoveEmoji() = %q", got)
	}
	out, err := RenderTemplate("hello {{.Name}}", map[string]string{"Name": "gokit"})
	if err != nil {
		t.Fatal(err)
	}
	if out != "hello gokit" {
		t.Fatalf("RenderTemplate() = %q", out)
	}
	if !IsEmail("a@example.com") || IsEmail("Alice <a@example.com>") || IsEmail("bad@@example") {
		t.Fatal("email validation helpers returned unexpected result")
	}
	if !IsURL("https://example.com") || !IsURL("ftp://example.com/a") || IsURL("/relative/path") || IsURL(" https://example.com") {
		t.Fatal("URL validation helpers returned unexpected result")
	}
	if !IsBlank(" \t") || RuneLen("你好") != 2 {
		t.Fatal("validation helpers returned unexpected result")
	}
}

func TestZipFilesRoundTripPreservesEmptyDirectory(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src")
	if err := os.MkdirAll(filepath.Join(src, "empty"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "a.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	archive := filepath.Join(dir, "archives", "out.zip")
	if err := ZipFiles(archive, src); err != nil {
		t.Fatal(err)
	}
	dest := filepath.Join(dir, "dest")
	if err := Unzip(archive, dest); err != nil {
		t.Fatal(err)
	}

	if info, err := os.Stat(filepath.Join(dest, "src", "empty")); err != nil || !info.IsDir() {
		t.Fatalf("empty directory was not restored, info=%v err=%v", info, err)
	}
	data, err := os.ReadFile(filepath.Join(dest, "src", "a.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Fatalf("unzipped content = %q", data)
	}
}

func TestUnzipRejectsPathTraversal(t *testing.T) {
	dir := t.TempDir()
	archive := filepath.Join(dir, "evil.zip")
	out, err := os.Create(archive)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(out)
	w, err := zw.Create("../evil.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write([]byte("evil")); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := out.Close(); err != nil {
		t.Fatal(err)
	}

	dest := filepath.Join(dir, "dest")
	if err := Unzip(archive, dest); err == nil {
		t.Fatal("Unzip() should reject path traversal entries")
	}
	if _, err := os.Stat(filepath.Join(dir, "evil.txt")); !os.IsNotExist(err) {
		t.Fatalf("path traversal wrote outside destination, err=%v", err)
	}
}
