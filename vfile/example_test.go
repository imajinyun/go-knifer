package vfile_test

import (
	"fmt"
	"strings"

	"github.com/imajinyun/go-knifer/vfile"
)

func ExampleMainName() {
	fmt.Println(vfile.MainName("/tmp/report.csv"))
	// Output: report
}

func ExampleExtension() {
	fmt.Println(vfile.Extension("/tmp/report.csv"))
	// Output: csv
}

func ExampleReadString() {
	content, _ := vfile.ReadString(strings.NewReader("hello"))
	fmt.Println(content)
	// Output: hello
}
