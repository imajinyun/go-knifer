package vpoi_test

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/imajinyun/go-knifer/vpoi"
	"github.com/xuri/excelize/v2"
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

func TestExcelFacadeWriteOptions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "book.xlsx")
	rows := [][]string{{"name", "score"}, {"go", "100"}}
	if err := vpoi.WriteRows(path, rows, vpoi.WithCreateParents(false)); err == nil {
		t.Fatal("WriteRows should fail when parent directory is missing and WithCreateParents(false) is set")
	}
	if err := vpoi.WriteRows(path, rows, vpoi.WithCreateParents(true), vpoi.WithFilePerm(0o600), vpoi.WithDirPerm(0o700)); err != nil {
		t.Fatalf("WriteRows with options: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat workbook: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("workbook file perm = %o, want 600", got)
	}
	if err := vpoi.WriteRows(path, rows, vpoi.WithOverwrite(false)); err == nil {
		t.Fatal("WriteRows should reject overwrite=false for existing workbook")
	}

	providerPath := filepath.Join(t.TempDir(), "nested", "provider.xlsx")
	var mkdirPath string
	var mkdirPerm fs.FileMode
	var chmodPath string
	var chmodPerm fs.FileMode
	if err := vpoi.WriteRows(providerPath, rows,
		vpoi.WithMkdirAll(func(path string, perm fs.FileMode) error {
			mkdirPath, mkdirPerm = path, perm
			return os.MkdirAll(path, perm)
		}),
		vpoi.WithChmod(func(path string, perm fs.FileMode) error {
			chmodPath, chmodPerm = path, perm
			return nil
		}),
		vpoi.WithDirPerm(0o700), vpoi.WithFilePerm(0o600),
	); err != nil {
		t.Fatalf("WriteRows provider options: %v", err)
	}
	if mkdirPath != filepath.Dir(providerPath) || mkdirPerm != 0o700 || chmodPath != providerPath || chmodPerm != 0o600 {
		t.Fatalf("providers mkdir=%q/%v chmod=%q/%v", mkdirPath, mkdirPerm, chmodPath, chmodPerm)
	}

	statErr := errors.New("stat denied")
	err = vpoi.WriteRows(providerPath, rows, vpoi.WithOverwrite(false), vpoi.WithStat(func(string) (os.FileInfo, error) {
		return nil, statErr
	}))
	if !errors.Is(err, statErr) {
		t.Fatalf("WriteRows should return custom stat error, got %v", err)
	}
}

func TestExcelFacadeSheetAndCellOptions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "book.xlsx")
	rows := [][]string{{"name", "score"}, {"go", "100"}}
	if err := vpoi.WriteRows(
		path, rows,
		vpoi.WithWriteSheet("Scores"),
		vpoi.WithStartCell(2, 2),
		vpoi.WithSaveOptions(excelize.Options{}),
	); err != nil {
		t.Fatalf("WriteRows with sheet/cell options: %v", err)
	}
	sheets, err := vpoi.SheetNames(path)
	if err != nil {
		t.Fatalf("SheetNames: %v", err)
	}
	if !reflect.DeepEqual(sheets, []string{"Scores"}) {
		t.Fatalf("SheetNames = %#v, want [Scores]", sheets)
	}
	f, err := excelize.OpenFile(path)
	if err != nil {
		t.Fatalf("open workbook: %v", err)
	}
	defer func() { _ = f.Close() }()
	if got, err := f.GetCellValue("Scores", "B2"); err != nil || got != "name" {
		t.Fatalf("Scores!B2 = %q, %v; want name, nil", got, err)
	}
	got, err := vpoi.ReadRows(path, vpoi.WithReadSheet("Scores"), vpoi.WithOpenOptions(excelize.Options{}))
	if err != nil {
		t.Fatalf("ReadRows with sheet/open options: %v", err)
	}
	if len(got) == 0 {
		t.Fatalf("ReadRows with sheet/open options returned no rows")
	}
}
