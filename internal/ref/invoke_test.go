package ref

import (
	"reflect"
	"testing"
)

func TestInvokeHelpers(t *testing.T) {
	s := &sample{Name: "alice"}
	got, err := Invoke(s, "Add", int8(1), int8(2))
	if err != nil || got != 3 {
		t.Fatalf("Invoke Add = %v, %v", got, err)
	}
	if _, err := Invoke(s, "Missing"); err == nil {
		t.Fatal("Invoke missing should fail")
	}
	if got, err := InvokeStatic(func(a int, b int) int { return a * b }, 2, 3); err != nil || got != 6 {
		t.Fatalf("InvokeStatic = %v, %v", got, err)
	}
	if got, err := InvokeFunc(func() (int, string) { return 1, "a" }); err != nil || !reflect.DeepEqual(got, []any{1, "a"}) {
		t.Fatalf("InvokeFunc multi return = %#v, %v", got, err)
	}
	if _, err := InvokeRaw(123); err == nil {
		t.Fatal("InvokeRaw non-func should fail")
	}
}
