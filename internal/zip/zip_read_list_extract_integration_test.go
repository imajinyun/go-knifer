package zip

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestZipReadListAndExtractCreatedArchive(t *testing.T) {
	tmp, src := newZipCreateSource(t)
	archive := filepath.Join(tmp, "out.zip")
	if err := ZipFiles(archive, false, src); err != nil {
		t.Fatalf("ZipFiles: %v", err)
	}

	data, err := GetBytes(archive, "a.txt")
	if err != nil || string(data) != "a" {
		t.Fatalf("GetBytes: %q %v", data, err)
	}
	names, err := ListFileNames(archive, "")
	if err != nil {
		t.Fatalf("ListFileNames: %v", err)
	}
	sort.Strings(names)
	if !reflect.DeepEqual(names, []string{"a.txt"}) {
		t.Fatalf("names: %#v", names)
	}
	dest := filepath.Join(tmp, "dest")
	if err := UnzipTo(archive, dest); err != nil {
		t.Fatalf("UnzipTo: %v", err)
	}
	if got, err := os.ReadFile(filepath.Join(dest, "nested", "b.txt")); err != nil || string(got) != "b" {
		t.Fatalf("unzipped: %q %v", got, err)
	}
	if info, err := os.Stat(filepath.Join(dest, "empty")); err != nil || !info.IsDir() {
		t.Fatalf("empty directory was not restored, info=%v err=%v", info, err)
	}
}
