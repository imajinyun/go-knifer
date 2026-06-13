package ref

import (
	"context"
	"reflect"
	"testing"
	"unsafe"
)

func TestAdditionalTypeClassificationHelpers(t *testing.T) {
	var nilSlice []string
	var nilUnsafePointer unsafe.Pointer
	if !IsNilValue(reflect.Value{}) || !IsNilValue(reflect.ValueOf(nilSlice)) || !IsNilValue(reflect.ValueOf(nilUnsafePointer)) {
		t.Fatal("IsNilValue did not treat invalid or nil-able nil values as nil")
	}
	if IsNilValue(reflect.ValueOf(1)) {
		t.Fatal("IsNilValue returned true for non-nil int")
	}

	objectTests := []struct {
		name       string
		value      any
		function   bool
		iteratee   bool
		collection bool
		slice      bool
		array      bool
		mapValue   bool
	}{
		{name: "nil"},
		{name: "function", value: func() {}, function: true},
		{name: "nil function", value: (func())(nil), function: true},
		{name: "slice", value: []int{}, iteratee: true, collection: true, slice: true},
		{name: "nil slice", value: nilSlice, iteratee: true, collection: true, slice: true},
		{name: "array", value: [1]int{}, iteratee: true, collection: true, array: true},
		{name: "map", value: map[string]int{}, iteratee: true, mapValue: true},
		{name: "string", value: "value"},
	}
	for _, tt := range objectTests {
		t.Run("object/"+tt.name, func(t *testing.T) {
			if IsFunction(tt.value) != tt.function {
				t.Fatalf("IsFunction(%T) = %v, want %v", tt.value, !tt.function, tt.function)
			}
			if IsIteratee(tt.value) != tt.iteratee {
				t.Fatalf("IsIteratee(%T) = %v, want %v", tt.value, !tt.iteratee, tt.iteratee)
			}
			if IsCollection(tt.value) != tt.collection {
				t.Fatalf("IsCollection(%T) = %v, want %v", tt.value, !tt.collection, tt.collection)
			}
			if IsSlice(tt.value) != tt.slice {
				t.Fatalf("IsSlice(%T) = %v, want %v", tt.value, !tt.slice, tt.slice)
			}
			if IsArray(tt.value) != tt.array {
				t.Fatalf("IsArray(%T) = %v, want %v", tt.value, !tt.array, tt.array)
			}
			if IsMap(tt.value) != tt.mapValue {
				t.Fatalf("IsMap(%T) = %v, want %v", tt.value, !tt.mapValue, tt.mapValue)
			}
		})
	}

	tests := []struct {
		name       string
		typ        reflect.Type
		funcType   bool
		rangeable  bool
		collection bool
		sliceType  bool
		arrayType  bool
		mapType    bool
	}{
		{name: "nil"},
		{name: "function", typ: reflect.TypeOf(func() {}), funcType: true},
		{name: "slice", typ: reflect.TypeOf([]int{}), rangeable: true, collection: true, sliceType: true},
		{name: "array", typ: reflect.TypeOf([1]int{}), rangeable: true, collection: true, arrayType: true},
		{name: "map", typ: reflect.TypeOf(map[string]int{}), rangeable: true, mapType: true},
		{name: "string", typ: reflect.TypeOf("value")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if IsFuncType(tt.typ) != tt.funcType {
				t.Fatalf("IsFuncType(%v) = %v", tt.typ, !tt.funcType)
			}
			if IsRangeableType(tt.typ) != tt.rangeable {
				t.Fatalf("IsRangeableType(%v) = %v", tt.typ, !tt.rangeable)
			}
			if IsCollectionType(tt.typ) != tt.collection {
				t.Fatalf("IsCollectionType(%v) = %v", tt.typ, !tt.collection)
			}
			if IsSliceType(tt.typ) != tt.sliceType {
				t.Fatalf("IsSliceType(%v) = %v", tt.typ, !tt.sliceType)
			}
			if IsArrayType(tt.typ) != tt.arrayType {
				t.Fatalf("IsArrayType(%v) = %v", tt.typ, !tt.arrayType)
			}
			if IsMapType(tt.typ) != tt.mapType {
				t.Fatalf("IsMapType(%v) = %v", tt.typ, !tt.mapType)
			}
		})
	}

	if !ImplementsError(reflect.TypeOf(sampleError{})) || ImplementsError(nil) || ImplementsError(reflect.TypeOf("value")) {
		t.Fatal("ImplementsError returned unexpected result")
	}
	if !ImplementsContext(reflect.TypeOf(context.Background())) || ImplementsContext(nil) || ImplementsContext(reflect.TypeOf(sampleError{})) {
		t.Fatal("ImplementsContext returned unexpected result")
	}
	if got := GetPublicFieldNames(sample{}); !reflect.DeepEqual(got, []string{"Name", "Age"}) {
		t.Fatalf("GetPublicFieldNames = %#v", got)
	}
	if got := GetPublicFieldNames((*sample)(nil)); !reflect.DeepEqual(got, []string{"Name", "Age"}) {
		t.Fatalf("GetPublicFieldNames pointer = %#v", got)
	}
	if got := GetPublicFieldNames(123); got != nil {
		t.Fatalf("GetPublicFieldNames non-struct = %#v", got)
	}
}
