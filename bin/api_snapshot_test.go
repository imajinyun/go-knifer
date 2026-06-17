package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateSnapshotIncludesCompatibilityDetails(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "go.mod", "module github.com/imajinyun/go-knifer\n\ngo 1.25.0\n")
	writeTestFile(t, root, "vcompat/compat.go", `package vcompat

const DefaultName = "demo"

var DefaultResult Result

type Alias = Result

type Result struct {
	Name   string
	Count  int
	hidden bool
}

func Build(input string) (Result, error) {
	return Result{Name: input}, nil
}

func (r Result) String() string { return r.Name }

func (r *Result) SetName(name string) { r.Name = name }

type Validator interface {
	Validate(Result) error
}
`)
	writeTestFile(t, root, "internal/hidden/hidden.go", `package hidden

func Hidden() {}
`)
	writeTestFile(t, root, "notfacade/notfacade.go", `package notfacade

func Hidden() {}
`)

	lines, err := generateSnapshot(root)
	if err != nil {
		t.Fatalf("generateSnapshot() error = %v", err)
	}
	snapshot := strings.Join(lines, "\n")
	for _, want := range []string{
		"github.com/imajinyun/go-knifer/vcompat",
		`const DefaultName untyped string = "demo"`,
		"var DefaultResult Result",
		"func Build(input string) (Result, error)",
		"type Alias = Result",
		"type Result struct{ Count int; Name string }",
		"method (*Result) SetName(name string)",
		"method (Result) String() string",
		"type Validator interface{ Validate(Result) error }",
	} {
		if !strings.Contains(snapshot, want) {
			t.Fatalf("snapshot missing %q:\n%s", want, snapshot)
		}
	}
	for _, unwanted := range []string{"hidden bool", "internal/hidden", "notfacade"} {
		if strings.Contains(snapshot, unwanted) {
			t.Fatalf("snapshot unexpectedly contains %q:\n%s", unwanted, snapshot)
		}
	}
}

func writeTestFile(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}
