package vpoi_test

import (
	"bytes"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/imajinyun/go-knifer/vpoi"
)

func TestExcelFacadeRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "book.xlsx")
	rows := [][]string{{"name", "score"}, {"go", "100"}}

	if err := vpoi.WriteRows(path, rows); err != nil {
		t.Fatalf("WriteRows: %v", err)
	}

	sheets, err := vpoi.SheetNames(path)
	if err != nil {
		t.Fatalf("SheetNames: %v", err)
	}
	if !reflect.DeepEqual(sheets, []string{vpoi.DefaultSheetName}) {
		t.Fatalf("SheetNames = %#v", sheets)
	}

	got, err := vpoi.ReadRows(path)
	if err != nil {
		t.Fatalf("ReadRows: %v", err)
	}
	if !reflect.DeepEqual(got, rows) {
		t.Fatalf("ReadRows = %#v, want %#v", got, rows)
	}
}

func TestExcelFacadeBufferRoundTrip(t *testing.T) {
	rows := [][]string{{"id", "name"}, {"1", "alice"}}
	buf, err := vpoi.WriteRowsToBuffer("Users", rows)
	if err != nil {
		t.Fatalf("WriteRowsToBuffer: %v", err)
	}

	got, err := vpoi.ReadRowsFromReader(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("ReadRowsFromReader: %v", err)
	}
	if !reflect.DeepEqual(got, rows) {
		t.Fatalf("ReadRowsFromReader = %#v, want %#v", got, rows)
	}
}
