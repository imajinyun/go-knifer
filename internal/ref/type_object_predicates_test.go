package ref

import "testing"

func TestObjectTypePredicates(t *testing.T) {
	var nilSlice []string
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
		t.Run(tt.name, func(t *testing.T) {
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
}
