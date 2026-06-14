package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeFunc(t *testing.T) {
	t.Run("nil resolver falls back to last-wins", func(t *testing.T) {
		a := map[string]int{"k": 1}
		b := map[string]int{"k": 2}
		assert.Equal(t, map[string]int{"k": 2}, MergeFunc[string, int](nil, a, b))
	})

	t.Run("sum resolver", func(t *testing.T) {
		a := map[string]int{"x": 1, "y": 2}
		b := map[string]int{"x": 10, "z": 30}
		c := map[string]int{"x": 100}
		got := MergeFunc(func(old, new int) int { return old + new }, a, b, c)
		assert.Equal(t, map[string]int{"x": 111, "y": 2, "z": 30}, got)
	})

	t.Run("keep old resolver", func(t *testing.T) {
		base := map[string]string{"theme": "dark"}
		override := map[string]string{"theme": "light", "lang": "zh"}
		got := MergeFunc(func(old, _ string) string { return old }, base, override)
		assert.Equal(t, map[string]string{"theme": "dark", "lang": "zh"}, got)
	})

	t.Run("slice append resolver", func(t *testing.T) {
		a := map[string][]int{"k": {1, 2}}
		b := map[string][]int{"k": {3, 4}, "x": {9}}
		got := MergeFunc(
			func(old, new []int) []int { return append(old, new...) },
			a, b,
		)
		assert.Equal(t, []int{1, 2, 3, 4}, got["k"])
		assert.Equal(t, []int{9}, got["x"])
	})
}
