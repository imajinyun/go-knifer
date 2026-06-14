package vref

import (
	"reflect"
	"testing"
)

func TestFacadeInstantiationAndInvocationHelpers(t *testing.T) {
	created, err := NewInstance(newFacadeSample, "constructed")
	if err != nil {
		t.Fatalf("NewInstance constructor: %v", err)
	}
	if got, ok := created.(facadeSample); !ok || got.Name != "constructed" {
		t.Fatalf("NewInstance constructor = %#v", created)
	}
	created, err = NewInstance(facadeSample{})
	if err != nil {
		t.Fatalf("NewInstance struct: %v", err)
	}
	if got, ok := created.(facadeSample); !ok || got.Name != "" {
		t.Fatalf("NewInstance struct = %#v", created)
	}
	if got := NewInstanceIfPossible((*facadeSample)(nil)); reflect.TypeOf(got) != reflect.TypeOf(&facadeSample{}) {
		t.Fatalf("NewInstanceIfPossible pointer type = %#v", got)
	}

	s := &facadeSample{Name: "alice"}
	method, ok := GetMethodOfObj(s, "Add", int32(2), int32(3))
	if !ok {
		t.Fatal("GetMethodOfObj Add with convertible args = false")
	}
	if got, err := InvokeWithCheck(s, method, int32(2), int32(3)); err != nil || got != 5 {
		t.Fatalf("InvokeWithCheck = %v, %v", got, err)
	}
	if got, err := InvokeMethod(s, method, 4, 5); err != nil || got != 9 {
		t.Fatalf("InvokeMethod = %v, %v", got, err)
	}
	if got, err := InvokeRaw(func(a, b int) int { return a + b }, 6, 7); err != nil || got != 13 {
		t.Fatalf("InvokeRaw = %v, %v", got, err)
	}
	if got, err := InvokeStatic(func() string { return "static" }); err != nil || got != "static" {
		t.Fatalf("InvokeStatic = %v, %v", got, err)
	}
	if got, err := InvokeFunc(func(a int) (int, string) { return a, "ok" }, 8); err != nil || !reflect.DeepEqual(got, []any{8, "ok"}) {
		t.Fatalf("InvokeFunc multi = %#v, %v", got, err)
	}
	if _, err := InvokeFunc("not-func"); err == nil {
		t.Fatal("InvokeFunc non-func error = nil")
	}
	if _, err := Invoke(s, "Missing"); err == nil {
		t.Fatal("Invoke missing method error = nil")
	}
	if _, err := InvokeWithCheck(s, reflect.Method{}); err == nil {
		t.Fatal("InvokeWithCheck invalid method error = nil")
	}
	if got := SetAccessible(s); got.Name != s.Name {
		t.Fatalf("SetAccessible = %#v", got)
	}
	RemoveFinalModify(&s)
}
