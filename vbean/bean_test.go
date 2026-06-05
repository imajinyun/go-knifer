package vbean_test

import (
	"errors"
	"strconv"
	"testing"

	knifer "github.com/imajinyun/go-knifer"
	"github.com/imajinyun/go-knifer/vbean"
)

type userDTO struct {
	Name string `bean:"name,alias=full_name"`
	Age  string `bean:"age"`
}

type userModel struct {
	Name string `json:"full_name"`
	Age  int    `json:"age"`
}

func TestFacadeCopyProperties(t *testing.T) {
	var dst userModel
	if err := vbean.CopyProperties(userDTO{Name: "alice", Age: "18"}, &dst); err != nil {
		t.Fatalf("CopyProperties() error = %v", err)
	}
	if dst.Name != "alice" || dst.Age != 18 {
		t.Fatalf("dst = %+v", dst)
	}
}

func TestFacadeToMap(t *testing.T) {
	got, err := vbean.ToMap(userDTO{Name: "bob", Age: "20"})
	if err != nil {
		t.Fatalf("ToMap() error = %v", err)
	}
	if got["name"] != "bob" || got["age"] != "20" {
		t.Fatalf("map = %#v", got)
	}
}

func TestFacadeBeanErrorContract(t *testing.T) {
	_, err := vbean.ToMap(nil)
	assertFacadeBeanCode(t, err, knifer.ErrCodeInvalidInput)

	var dst userModel
	err = vbean.CopyProperties(map[string]any{"age": "not-a-number"}, &dst)
	assertFacadeBeanCode(t, err, knifer.ErrCodeInvalidInput)
	var numErr *strconv.NumError
	if !errors.As(err, &numErr) {
		t.Fatalf("CopyProperties should preserve strconv.NumError cause: %v", err)
	}
}

func assertFacadeBeanCode(t *testing.T, err error, code knifer.ErrCode) {
	t.Helper()
	if err == nil {
		t.Fatalf("err = nil, want %s", code)
	}
	if !errors.Is(err, code) {
		t.Fatalf("errors.Is(%v, %s) = false", err, code)
	}
	got, ok := knifer.CodeOf(err)
	if !ok || got != code {
		t.Fatalf("CodeOf(%v) = %q, %v; want %q, true", err, got, ok, code)
	}
	var beanErr *vbean.Error
	if !errors.As(err, &beanErr) {
		t.Fatalf("errors.As(err, *vbean.Error) = false: %v", err)
	}
}
