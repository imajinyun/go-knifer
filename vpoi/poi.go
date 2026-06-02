package vpoi

import (
	"bytes"
	"io"

	poiimpl "github.com/imajinyun/go-knifer/internal/poi"
)

const (
	// DefaultSheetName is the default worksheet name used for read/write helpers.
	DefaultSheetName = poiimpl.DefaultSheetName
)

var (
	// ErrNoSheet indicates that a workbook does not contain any worksheet.
	ErrNoSheet = poiimpl.ErrNoSheet
	// ErrEmptySheetName indicates an empty worksheet name.
	ErrEmptySheetName = poiimpl.ErrEmptySheetName
)

// SheetNames returns all worksheet names in path.
func SheetNames(path string) ([]string, error) { return poiimpl.SheetNames(path) }

// ReadRows reads rows from the first worksheet in path.
func ReadRows(path string) ([][]string, error) { return poiimpl.ReadRows(path) }

// ReadSheetRows reads rows from sheet in path.
func ReadSheetRows(path, sheet string) ([][]string, error) {
	return poiimpl.ReadSheetRows(path, sheet)
}

// ReadRowsFromReader reads rows from the first worksheet in r.
func ReadRowsFromReader(r io.Reader) ([][]string, error) { return poiimpl.ReadRowsFromReader(r) }

// WriteRows writes rows into path using the default worksheet name.
func WriteRows(path string, rows [][]string) error { return poiimpl.WriteRows(path, rows) }

// WriteSheetRows writes rows into path using sheet.
func WriteSheetRows(path, sheet string, rows [][]string) error {
	return poiimpl.WriteSheetRows(path, sheet, rows)
}

// WriteSheets writes multiple worksheets into path.
func WriteSheets(path string, sheets map[string][][]string) error {
	return poiimpl.WriteSheets(path, sheets)
}

// WriteRowsToBuffer writes rows into an in-memory XLSX workbook.
func WriteRowsToBuffer(sheet string, rows [][]string) (*bytes.Buffer, error) {
	return poiimpl.WriteRowsToBuffer(sheet, rows)
}
