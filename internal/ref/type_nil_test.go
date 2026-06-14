package ref

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestIsNilValue(t *testing.T) {
	var nilSlice []string
	var nilUnsafePointer unsafe.Pointer
	if !IsNilValue(reflect.Value{}) || !IsNilValue(reflect.ValueOf(nilSlice)) || !IsNilValue(reflect.ValueOf(nilUnsafePointer)) {
		t.Fatal("IsNilValue did not treat invalid or nil-able nil values as nil")
	}
	if IsNilValue(reflect.ValueOf(1)) {
		t.Fatal("IsNilValue returned true for non-nil int")
	}
}
