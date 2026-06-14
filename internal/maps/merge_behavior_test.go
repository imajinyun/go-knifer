package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// snapshot deep-copies a map to detect mutation of the input.
func snapshot[K comparable, V any](m map[K]V) map[K]V {
	out := make(map[K]V, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func TestMerge(t *testing.T) {
	t.Run("zero arg", func(t *testing.T) {
		out := Merge[string, int]()
		assert.NotNil(t, out)
		assert.Empty(t, out)
	})

	t.Run("nil single arg", func(t *testing.T) {
		var nilMap map[string]int
		out := Merge(nilMap)
		assert.NotNil(t, out)
		assert.Empty(t, out)
	})

	t.Run("single arg returns clone", func(t *testing.T) {
		src := map[string]int{"a": 1}
		out := Merge(src)
		assert.Equal(t, src, out)
		// mutating the result must not affect the input
		out["a"] = 99
		assert.Equal(t, 1, src["a"])
	})

	t.Run("later overrides earlier", func(t *testing.T) {
		a := map[string]int{"k": 1, "x": 10}
		b := map[string]int{"k": 2, "y": 20}
		c := map[string]int{"k": 3}
		got := Merge(a, b, c)
		assert.Equal(t, map[string]int{"k": 3, "x": 10, "y": 20}, got)
	})

	t.Run("inputs are not mutated", func(t *testing.T) {
		a := map[string]int{"k": 1}
		b := map[string]int{"k": 2}
		snapA := snapshot(a)
		snapB := snapshot(b)
		_ = Merge(a, b)
		assert.Equal(t, snapA, a)
		assert.Equal(t, snapB, b)
	})
}

func TestMergeIsAssociativeForLastWins(t *testing.T) {
	a := map[string]int{"x": 1, "y": 1}
	b := map[string]int{"y": 2, "z": 2}
	c := map[string]int{"z": 3}

	left := Merge(Merge(a, b), c)
	right := Merge(a, Merge(b, c))
	assert.Equal(t, left, right)
}
