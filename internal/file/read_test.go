package file

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	knifer "github.com/imajinyun/knifer-go"
	strutil "github.com/imajinyun/knifer-go/internal/str"
)

func TestReadString(t *testing.T) {
	r := ReaderFromString("hello world")
	got, err := ReadString(r)
	if err != nil || got != "hello world" {
		t.Fatalf("ReadString: %v %q", err, got)
	}
}

func TestReadLines(t *testing.T) {
	r := ReaderFromString("a\nb\nc")
	lines, err := ReadLines(r)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(lines) != 3 || lines[0] != "a" || lines[2] != "c" {
		t.Fatalf("ReadLines: %v", lines)
	}
}

func TestReadOptions(t *testing.T) {
	if got, err := ReadStringWithOptions(ReaderFromString("abc"), WithMaxBytes(3)); err != nil || got != "abc" {
		t.Fatalf("ReadStringWithOptions exact limit = %q, %v", got, err)
	}
	if _, err := ReadStringWithOptions(ReaderFromString("abcd"), WithMaxBytes(3)); err == nil {
		t.Fatal("ReadStringWithOptions over limit error = nil")
	}

	lines, err := ReadLinesWithOptions(ReaderFromString("abc"), WithMaxBytes(3), WithInitialLineBuffer(1), WithMaxLineBytes(4))
	if err != nil {
		t.Fatalf("ReadLinesWithOptions exact limit: %v", err)
	}
	if len(lines) != 1 || lines[0] != "abc" {
		t.Fatalf("ReadLinesWithOptions lines = %v", lines)
	}
	if _, err := ReadLinesWithOptions(ReaderFromString("abcd"), WithMaxBytes(3), WithMaxLineBytes(4)); err == nil {
		t.Fatal("ReadLinesWithOptions over limit error = nil")
	}
}

func TestReadOptionsDefaultLimitAndExplicitUnlimited(t *testing.T) {
	if _, err := ReadStringWithOptions(ReaderFromString(strings.Repeat("x", int(DefaultMaxBytes)+1))); !errors.Is(err, knifer.ErrCodeInvalidInput) {
		t.Fatalf("ReadStringWithOptions default limit error = %v, want invalid input", err)
	}
	got, err := ReadStringWithOptions(ReaderFromString("abcd"), WithMaxBytes(3), WithUnlimitedRead())
	if err != nil || got != "abcd" {
		t.Fatalf("WithUnlimitedRead() = %q, %v", got, err)
	}
}

func TestReadOptionsCharset(t *testing.T) {
	gbk, err := strutil.FromUTF8([]byte("中文\n下一行"), "gbk")
	if err != nil {
		t.Fatalf("FromUTF8 error = %v", err)
	}

	text, err := ReadStringWithOptions(bytes.NewReader(gbk), WithCharset("gbk"))
	if err != nil {
		t.Fatalf("ReadStringWithOptions charset error = %v", err)
	}
	if text != "中文\n下一行" {
		t.Fatalf("ReadStringWithOptions charset = %q", text)
	}

	lines, err := ReadLinesWithOptions(bytes.NewReader(gbk), WithCharset("gbk"))
	if err != nil {
		t.Fatalf("ReadLinesWithOptions charset error = %v", err)
	}
	if len(lines) != 2 || lines[0] != "中文" || lines[1] != "下一行" {
		t.Fatalf("ReadLinesWithOptions charset = %v", lines)
	}
}

func TestFileReadOptionsCharset(t *testing.T) {
	gbk, err := strutil.FromUTF8([]byte("中文\n下一行"), "gbk")
	if err != nil {
		t.Fatalf("FromUTF8 error = %v", err)
	}
	path := filepath.Join(t.TempDir(), "gbk.txt")
	if err := os.WriteFile(path, gbk, 0o600); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	text, err := FileReadStringWithOptions(path, WithCharset("gbk"))
	if err != nil {
		t.Fatalf("FileReadStringWithOptions charset error = %v", err)
	}
	if text != "中文\n下一行" {
		t.Fatalf("FileReadStringWithOptions charset = %q", text)
	}

	lines, err := FileReadLinesWithOptions(path, WithCharset("gbk"))
	if err != nil {
		t.Fatalf("FileReadLinesWithOptions charset error = %v", err)
	}
	if len(lines) != 2 || lines[0] != "中文" || lines[1] != "下一行" {
		t.Fatalf("FileReadLinesWithOptions charset = %v", lines)
	}
}
