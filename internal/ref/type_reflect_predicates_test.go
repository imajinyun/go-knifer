package ref

import (
	"reflect"
	"testing"
)

func TestReflectTypePredicates(t *testing.T) {
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
}
