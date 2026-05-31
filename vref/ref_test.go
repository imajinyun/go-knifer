package vref

import (
	"reflect"
	"testing"
)

type facadeSample struct {
	Name string `json:"name"`
}

func (s facadeSample) GetName() string      { return s.Name }
func (s facadeSample) Add(a int, b int) int { return a + b }

func TestFacadeReflectionHelpers(t *testing.T) {
	s := &facadeSample{Name: "alice"}
	if !HasField(s, "name") || GetFieldValue(s, "name") != "alice" {
		t.Fatal("field facade failed")
	}
	if err := SetFieldValue(s, "Name", "bob"); err != nil || s.Name != "bob" {
		t.Fatalf("SetFieldValue facade = %v name=%s", err, s.Name)
	}
	if _, ok := GetMethod(s, false, "Add", reflect.TypeOf(1), reflect.TypeOf(2)); !ok {
		t.Fatal("method facade failed")
	}
	got, err := Invoke(s, "Add", 2, 3)
	if err != nil || got != 5 {
		t.Fatalf("Invoke facade = %v, %v", got, err)
	}
	if TypeOf(s).Kind() != reflect.Pointer || IndirectType(TypeOf(s)).Name() != "facadeSample" {
		t.Fatal("type facade failed")
	}
}
