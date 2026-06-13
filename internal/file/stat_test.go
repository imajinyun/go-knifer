package file

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestFileWriteReadStat(t *testing.T) {
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
}

func TestMkdirTouchDelStat(t *testing.T) {
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
}

func TestFileStatProviderOptions(t *testing.T) {
	regular := fakeFileInfo{name: "file.txt", size: 42}
	dir := fakeFileInfo{name: "dir", dir: true}
	stat := func(path string) (fs.FileInfo, error) {
		switch path {
		case "file.txt":
			return regular, nil
		case "dir":
			return dir, nil
		default:
			return nil, os.ErrNotExist
		}
	}
	if !FileExistsWithOptions("file.txt", WithStat(stat)) || !IsFileWithOptions("file.txt", WithStat(stat)) {
		t.Fatal("file stat providers failed")
	}
	if !IsDirectoryWithOptions("dir", WithStat(stat)) {
		t.Fatal("directory stat provider failed")
	}
	if got := FileSizeWithOptions("file.txt", WithStat(stat)); got != 42 {
		t.Fatalf("FileSizeWithOptions() = %d, want 42", got)
	}
}
