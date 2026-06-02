package poi

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/xuri/excelize/v2"
)

const (
	// DefaultSheetName is the default worksheet name used for read/write helpers.
	DefaultSheetName = "Sheet1"
)

var (
	// ErrNoSheet indicates that a workbook does not contain any worksheet.
	ErrNoSheet = errors.New("poi: workbook has no sheet")
	// ErrEmptySheetName indicates an empty worksheet name.
	ErrEmptySheetName = errors.New("poi: sheet name is empty")
)

// SheetNames returns all worksheet names in path.
func SheetNames(path string) ([]string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	return f.GetSheetList(), nil
}

// ReadRows reads rows from the first worksheet in path.
func ReadRows(path string) ([][]string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	return readFirstSheetRows(f)
}

// ReadSheetRows reads rows from sheet in path.
func ReadSheetRows(path, sheet string) ([][]string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	return readSheetRows(f, sheet)
}

// ReadRowsFromReader reads rows from the first worksheet in r.
func ReadRowsFromReader(r io.Reader) ([][]string, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	return readFirstSheetRows(f)
}

// WriteRows writes rows into path using the default worksheet name.
func WriteRows(path string, rows [][]string) error {
	return WriteSheetRows(path, DefaultSheetName, rows)
}

// WriteSheetRows writes rows into path using sheet.
func WriteSheetRows(path, sheet string, rows [][]string) error {
	if sheet == "" {
		return ErrEmptySheetName
	}
	f := excelize.NewFile()
	defer func() { _ = f.Close() }()
	if err := replaceDefaultSheet(f, sheet); err != nil {
		return err
	}
	if err := setRows(f, sheet, rows); err != nil {
		return err
	}
	if err := ensureParentDir(path); err != nil {
		return err
	}
	return f.SaveAs(path)
}

// WriteSheets writes multiple worksheets into path.
func WriteSheets(path string, sheets map[string][][]string) error {
	f := excelize.NewFile()
	defer func() { _ = f.Close() }()

	if len(sheets) == 0 {
		if err := ensureParentDir(path); err != nil {
			return err
		}
		return f.SaveAs(path)
	}

	first := true
	for sheet, rows := range sheets {
		if sheet == "" {
			return ErrEmptySheetName
		}
		if first {
			if err := replaceDefaultSheet(f, sheet); err != nil {
				return err
			}
			first = false
		} else if _, err := f.NewSheet(sheet); err != nil {
			return err
		}
		if err := setRows(f, sheet, rows); err != nil {
			return err
		}
	}
	if err := ensureParentDir(path); err != nil {
		return err
	}
	return f.SaveAs(path)
}

// WriteRowsToBuffer writes rows into an in-memory XLSX workbook.
func WriteRowsToBuffer(sheet string, rows [][]string) (*bytes.Buffer, error) {
	if sheet == "" {
		return nil, ErrEmptySheetName
	}
	f := excelize.NewFile()
	defer func() { _ = f.Close() }()
	if err := replaceDefaultSheet(f, sheet); err != nil {
		return nil, err
	}
	if err := setRows(f, sheet, rows); err != nil {
		return nil, err
	}
	return f.WriteToBuffer()
}

func readFirstSheetRows(f *excelize.File) ([][]string, error) {
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, ErrNoSheet
	}
	return readSheetRows(f, sheets[0])
}

func readSheetRows(f *excelize.File, sheet string) ([][]string, error) {
	if sheet == "" {
		return nil, ErrEmptySheetName
	}
	return f.GetRows(sheet)
}

func replaceDefaultSheet(f *excelize.File, sheet string) error {
	if sheet == DefaultSheetName {
		return nil
	}
	if err := f.SetSheetName(DefaultSheetName, sheet); err != nil {
		return err
	}
	return nil
}

func setRows(f *excelize.File, sheet string, rows [][]string) error {
	for rowIndex, row := range rows {
		for colIndex, value := range row {
			cell, err := excelize.CoordinatesToCellName(colIndex+1, rowIndex+1)
			if err != nil {
				return fmt.Errorf("poi: cell coordinates row=%d col=%d: %w", rowIndex+1, colIndex+1, err)
			}
			if err := f.SetCellStr(sheet, cell, value); err != nil {
				return err
			}
		}
	}
	return nil
}

func ensureParentDir(path string) error {
	dir := filepath.Dir(path)
	if dir == "." || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o750)
}
