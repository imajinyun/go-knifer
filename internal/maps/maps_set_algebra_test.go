package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInverse(t *testing.T) {
	in := map[string]int{"a": 1, "b": 2}
	inv := Inverse(in)
	assert.Equal(t, map[int]string{1: "a", 2: "b"}, inv)
}

func TestIntersect(t *testing.T) {
	a := map[string]int{"a": 1, "b": 2, "c": 3}
	b := map[string]int{"b": 20, "c": 30, "d": 40}
	c := map[string]int{"c": 300, "d": 400}

	got := Intersect(a, b, c)
	assert.Equal(t, map[string]int{"c": 300}, got)

	// edge: zero / one input
	assert.Empty(t, Intersect[string, int]())
	assert.Equal(t, a, Intersect(a))

	// empty intersection
	assert.Empty(t, Intersect(
		map[string]int{"a": 1},
		map[string]int{"b": 2},
	))
}

func TestDiff(t *testing.T) {
	a := map[string]int{"a": 1, "b": 2, "c": 3}
	b := map[string]int{"a": 10}
	c := map[string]int{"b": 20}
	assert.Equal(t, map[string]int{"c": 3}, Diff(a, b, c))

	// no others → returns clone of a
	assert.Equal(t, a, Diff(a))
}

func TestSymmetricDiff(t *testing.T) {
	a := map[string]int{"a": 1, "b": 2}
	b := map[string]int{"b": 20, "c": 3}
	assert.Equal(t, map[string]int{"a": 1, "c": 3}, SymmetricDiff(a, b))
}
