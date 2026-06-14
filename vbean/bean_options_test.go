package vbean_test

import (
	"testing"

	"github.com/imajinyun/go-knifer/vbean"
)

func TestFacadeBeanOptions(t *testing.T) {
	type customTagged struct {
		Name string `db:"user_name"`
		Age  int    `db:"age"`
	}
	got, err := vbean.ToMap(customTagged{Name: "casey", Age: 0},
		vbean.WithTagNames("db"),
		vbean.WithIgnoreZero(true),
	)
	if err != nil {
		t.Fatalf("ToMap() with options error = %v", err)
	}
	if got["user_name"] != "casey" {
		t.Fatalf("ToMap() user_name = %#v", got["user_name"])
	}
	if _, ok := got["age"]; ok {
		t.Fatalf("ToMap() should skip zero age with WithIgnoreZero: %#v", got)
	}

	var dst userModel
	if err := vbean.ToStruct(map[string]any{"FULL_NAME": "drew", "age": "21"}, &dst,
		vbean.WithCaseInsensitive(true),
		vbean.WithWeaklyTyped(true),
	); err != nil {
		t.Fatalf("ToStruct() with options error = %v", err)
	}
	if dst.Name != "drew" || dst.Age != 21 {
		t.Fatalf("ToStruct() with options dst = %+v", dst)
	}

	dst = userModel{Name: "existing", Age: 30}
	if err := vbean.Copy(map[string]any{"full_name": "", "age": "22"}, &dst, vbean.WithIgnoreEmpty(true)); err != nil {
		t.Fatalf("Copy() with WithIgnoreEmpty error = %v", err)
	}
	if dst.Name != "existing" || dst.Age != 22 {
		t.Fatalf("Copy() WithIgnoreEmpty dst = %+v", dst)
	}

	var strict userModel
	if err := vbean.CopyProperties(map[string]any{"age": "23"}, &strict, vbean.WithWeaklyTyped(false)); err == nil {
		t.Fatal("CopyProperties() WithWeaklyTyped(false) error = nil, want strict assignment error")
	}
}
