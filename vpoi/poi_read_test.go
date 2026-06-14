package vpoi_test

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/imajinyun/go-knifer/vpoi"
	"github.com/xuri/excelize/v2"
)

func TestExcelFacadeReadRoundTrip(t *testing.T) {
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
	sheets, err = vpoi.SheetNamesWithOptions(path, vpoi.WithOpenOptions(excelize.Options{}))
	if err != nil {
		t.Fatalf("SheetNamesWithOptions: %v", err)
	}
	if !reflect.DeepEqual(sheets, []string{vpoi.DefaultSheetName}) {
		t.Fatalf("SheetNamesWithOptions = %#v", sheets)
	}

	got, err := vpoi.ReadRows(path)
	if err != nil {
		t.Fatalf("ReadRows: %v", err)
	}
	if !reflect.DeepEqual(got, rows) {
		t.Fatalf("ReadRows = %#v, want %#v", got, rows)
	}
	got, err = vpoi.ReadSheetRowsWithOptions(path, vpoi.DefaultSheetName, vpoi.WithOpenOptions(excelize.Options{}))
	if err != nil {
		t.Fatalf("ReadSheetRowsWithOptions: %v", err)
	}
	if !reflect.DeepEqual(got, rows) {
		t.Fatalf("ReadSheetRowsWithOptions = %#v, want %#v", got, rows)
	}
}
