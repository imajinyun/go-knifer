package vpoi_test

import (
	"bytes"
	"fmt"

	"github.com/imajinyun/go-knifer/vpoi"
)

func ExampleSheetNames() {
	names, err := vpoi.SheetNames("nonexistent.xlsx")
	fmt.Println(names)
	fmt.Println(err != nil)
	// Output:
	// []
	// true
}

func ExampleWriteRowsToBuffer() {
	rows := [][]string{
		{"Name", "Age"},
		{"Alice", "30"},
	}

	buf, err := vpoi.WriteRowsToBuffer("Users", rows)
	fmt.Println(buf.Len() > 0)
	fmt.Println(err)
	// Output:
	// true
	// <nil>
}

func ExampleReadRowsFromReader() {
	rows := [][]string{
		{"Name", "Age"},
		{"Alice", "30"},
	}
	buf, _ := vpoi.WriteRowsToBuffer("Users", rows)
	got, err := vpoi.ReadRowsFromReader(bytes.NewReader(buf.Bytes()))

	fmt.Println(got)
	fmt.Println(err)
	// Output:
	// [[Name Age] [Alice 30]]
	// <nil>
}
