package poi

import (
	"bytes"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/xuri/excelize/v2"
)

func TestWriteAnyRowsAndReadCells(t *testing.T) {
	path := filepath.Join(t.TempDir(), "typed.xlsx")
	createdAt := time.Date(2026, 6, 24, 9, 30, 0, 0, time.UTC)
	rows := [][]any{
		{"name", "score", "active", "created_at", "empty"},
		{"go", 100, true, createdAt, nil},
	}

	if err := WriteAnyRows(path, rows); err != nil {
		t.Fatalf("WriteAnyRows: %v", err)
	}

	cells, err := ReadCells(path, WithReadStartCell(2, 2), WithReadLimit(1, 3))
	if err != nil {
		t.Fatalf("ReadCells: %v", err)
	}
	if len(cells) != 1 || len(cells[0]) != 3 {
		t.Fatalf("ReadCells shape = %#v, want one row with three cells", cells)
	}
	if cells[0][0].Axis != "B2" || cells[0][0].Row != 2 || cells[0][0].Col != 2 {
		t.Fatalf("first cell position = %#v, want B2 row=2 col=2", cells[0][0])
	}
	if cells[0][0].Value != "100" || cells[0][0].Type != excelize.CellTypeNumber {
		t.Fatalf("score cell = %#v, want value 100 numeric type", cells[0][0])
	}
	if cells[0][1].Value != "TRUE" || cells[0][1].Type != excelize.CellTypeBool {
		t.Fatalf("active cell = %#v, want TRUE bool type", cells[0][1])
	}
	if cells[0][2].Value == "" || cells[0][2].Type != excelize.CellTypeNumber {
		t.Fatalf("date cell = %#v, want formatted value with numeric type", cells[0][2])
	}

	stringRows, err := ReadRows(path)
	if err != nil {
		t.Fatalf("ReadRows: %v", err)
	}
	if got, want := stringRows[1][1], "100"; got != want {
		t.Fatalf("ReadRows typed numeric value = %q, want %q", got, want)
	}
}

func TestWriteSheetAnyRowsAndReadCellsFromReader(t *testing.T) {
	rows := [][]any{{"id", "ok"}, {1, false}}
	buf, err := WriteAnyRowsToBuffer("Typed", rows)
	if err != nil {
		t.Fatalf("WriteAnyRowsToBuffer: %v", err)
	}

	cells, err := ReadCellsFromReader(bytes.NewReader(buf.Bytes()), WithReadSheet("Typed"))
	if err != nil {
		t.Fatalf("ReadCellsFromReader: %v", err)
	}
	if got, want := cells[1][0].Type, excelize.CellTypeNumber; got != want {
		t.Fatalf("id type = %v, want %v", got, want)
	}
	if got, want := cells[1][1].Type, excelize.CellTypeBool; got != want {
		t.Fatalf("ok type = %v, want %v", got, want)
	}

	path := filepath.Join(t.TempDir(), "sheet-typed.xlsx")
	if err := WriteSheetAnyRows(path, "Typed", rows); err != nil {
		t.Fatalf("WriteSheetAnyRows: %v", err)
	}
	got, err := ReadSheetCellsWithOptions(path, "Typed", WithReadStartCell(2, 1), WithReadLimit(1, 2))
	if err != nil {
		t.Fatalf("ReadSheetCellsWithOptions: %v", err)
	}
	if !reflect.DeepEqual([]string{got[0][0].Value, got[0][1].Value}, []string{"1", "FALSE"}) {
		t.Fatalf("ReadSheetCellsWithOptions values = %#v, want [1 FALSE]", got)
	}
}
