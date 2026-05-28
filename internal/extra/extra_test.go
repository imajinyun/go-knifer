package extra

import (
	"bytes"
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
	if !ContainsEmoji("hi😀") {
		t.Fatal("ContainsEmoji() = false")
	}
	if got := RemoveEmoji("hi😀"); got != "hi" {
		t.Fatalf("RemoveEmoji() = %q", got)
	}
	out, err := RenderTemplate("hello {{.Name}}", map[string]string{"Name": "gokit"})
	if err != nil {
		t.Fatal(err)
	}
	if out != "hello gokit" {
		t.Fatalf("RenderTemplate() = %q", out)
	}
	if !IsEmail("a@example.com") || !IsURL("https://example.com") || !IsBlank(" \t") {
		t.Fatal("validation helpers returned unexpected result")
	}
}
