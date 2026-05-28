package base

import (
	"path/filepath"
	"strings"
	"testing"
)

// 对应 hutool-core IoUtilTest / FileUtilTest / FileNameUtilTest。

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

func TestFileWriteRead(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "a.txt")
	if err := FileWriteString(path, "你好"); err != nil {
		t.Fatalf("write: %v", err)
	}
	if !FileExists(path) || !IsFile(path) {
		t.Fatalf("FileExists/IsFile failed")
	}
	if FileSize(path) <= 0 {
		t.Fatalf("FileSize failed")
	}
	got, err := FileReadString(path)
	if err != nil || got != "你好" {
		t.Fatalf("read: %v %q", err, got)
	}
	if err := FileAppendString(path, "X"); err != nil {
		t.Fatalf("append: %v", err)
	}
	got, _ = FileReadString(path)
	if !strings.HasSuffix(got, "X") {
		t.Fatalf("append result: %q", got)
	}
}

func TestFileLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "lines.txt")
	if err := FileWriteString(path, "x\ny\nz"); err != nil {
		t.Fatalf("write: %v", err)
	}
	lines, err := FileReadLines(path)
	if err != nil {
		t.Fatalf("read lines: %v", err)
	}
	if len(lines) != 3 || lines[2] != "z" {
		t.Fatalf("lines: %v", lines)
	}
}

func TestMkdirTouchDel(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "x", "y")
	if err := Mkdir(sub); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if !IsDirectory(sub) {
		t.Fatalf("IsDirectory failed")
	}
	f := filepath.Join(sub, "a.txt")
	if err := Touch(f); err != nil {
		t.Fatalf("touch: %v", err)
	}
	if !IsFile(f) {
		t.Fatalf("touch did not create file")
	}
	if err := Del(f); err != nil {
		t.Fatalf("del: %v", err)
	}
	if FileExists(f) {
		t.Fatalf("Del did not remove")
	}
}

func TestFileCopy(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "a.txt")
	dst := filepath.Join(dir, "sub", "b.txt")
	if err := FileWriteString(src, "hello"); err != nil {
		t.Fatalf("write src: %v", err)
	}
	if err := FileCopy(src, dst); err != nil {
		t.Fatalf("copy: %v", err)
	}
	got, _ := FileReadString(dst)
	if got != "hello" {
		t.Fatalf("copy content: %q", got)
	}
}

func TestMainNameAndExtension(t *testing.T) {
	if MainName("/x/y/foo.txt") != "foo" {
		t.Fatalf("MainName failed")
	}
	if Extension("/x/y/foo.txt") != "txt" {
		t.Fatalf("Extension failed")
	}
	if MainName("foo") != "foo" || Extension("foo") != "" {
		t.Fatalf("no-ext failed")
	}
}
