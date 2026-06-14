package vref

import (
	"context"
	"reflect"
	"testing"
	"unsafe"
)

func TestFacadeTypeHelpers(t *testing.T) {
	var nilSlice []string
	if !IsNil(nilSlice) || IsNil(1) {
		t.Fatal("IsNil facade returned unexpected result")
	}
	var nilUnsafePointer unsafe.Pointer
	if !IsNilValue(reflect.Value{}) || !IsNilValue(reflect.ValueOf(nilSlice)) || !IsNilValue(reflect.ValueOf(nilUnsafePointer)) {
		t.Fatal("IsNilValue facade returned unexpected result")
	}
	if IsNilValue(reflect.ValueOf(1)) {
		t.Fatal("IsNilValue facade returned true for non-nil int")
	}
	if !IsFunction(newFacadeSample) || IsFunction(nil) {
		t.Fatal("IsFunction facade returned unexpected result")
	}
	if !IsIteratee(map[string]int{}) || !IsIteratee([]int{}) || !IsIteratee([1]int{}) || IsIteratee("value") || IsIteratee(nil) {
		t.Fatal("IsIteratee facade returned unexpected result")
	}
	if !IsCollection([]int{}) || !IsCollection([1]int{}) || IsCollection(map[string]int{}) || IsCollection(nil) {
		t.Fatal("IsCollection facade returned unexpected result")
	}
	if !IsSlice(nilSlice) || !IsArray([1]int{}) || !IsMap(map[string]int{}) {
		t.Fatal("specific object predicate facade returned unexpected result")
	}
	if IsSlice(nil) || IsArray(nil) || IsMap(nil) {
		t.Fatal("specific object predicate facade should reject nil")
	}
	if !IsFuncType(reflect.TypeOf(newFacadeSample)) || IsFuncType(nil) {
		t.Fatal("IsFuncType facade returned unexpected result")
	}
	if !IsRangeableType(reflect.TypeOf(map[string]int{})) || !IsRangeableType(reflect.TypeOf([]int{})) || IsRangeableType(reflect.TypeOf(1)) {
		t.Fatal("IsRangeableType facade returned unexpected result")
	}
	if !IsCollectionType(reflect.TypeOf([1]int{})) || !IsCollectionType(reflect.TypeOf([]int{})) || IsCollectionType(reflect.TypeOf(map[string]int{})) {
		t.Fatal("IsCollectionType facade returned unexpected result")
	}
	if !IsSliceType(reflect.TypeOf([]int{})) || !IsArrayType(reflect.TypeOf([1]int{})) || !IsMapType(reflect.TypeOf(map[string]int{})) {
		t.Fatal("specific type predicate facade returned unexpected result")
	}
	if !ImplementsError(reflect.TypeOf(facadeError{})) || ImplementsError(nil) || !ImplementsContext(reflect.TypeOf(context.Background())) || ImplementsContext(nil) {
		t.Fatal("interface implementation facade returned unexpected result")
	}
	if got := ValueOf(nil); got.IsValid() {
		t.Fatalf("ValueOf(nil).IsValid() = true: %v", got)
	}
	if got := IndirectValue(reflect.ValueOf(&facadeSample{Name: "alice"})); !got.IsValid() || got.FieldByName("Name").String() != "alice" {
		t.Fatalf("IndirectValue = %v", got)
	}

	ctor := GetConstructor(newFacadeSample)
	if !ctor.IsValid() || len(GetConstructors(newFacadeSample)) != 1 || len(GetConstructorsDirectly(newFacadeSample)) != 1 {
		t.Fatal("constructor helpers did not expose function target")
	}
	if got := GetConstructor(facadeSample{}); got.IsValid() {
		t.Fatalf("GetConstructor(non-func).IsValid() = true: %v", got)
	}
}
