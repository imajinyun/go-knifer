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

func TestMergeWithOverwrite(t *testing.T) {
	dst := map[string]int{"a": 1, "shared": 1}
	src1 := map[string]int{"b": 2, "shared": 2}
	src2 := map[string]int{"c": 3, "shared": 3}

	MergeWithOverwrite(dst, src1, nil, src2)

	assert.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3, "shared": 3}, dst)
	assert.Equal(t, map[string]int{"b": 2, "shared": 2}, src1, "source maps must not be mutated")
}

func TestMergeWithoutOverwrite(t *testing.T) {
	dst := map[string]int{"a": 1, "shared": 1}
	src1 := map[string]int{"b": 2, "shared": 2}
	src2 := map[string]int{"b": 20, "c": 3, "shared": 3}

	MergeWithoutOverwrite(dst, src1, nil, src2)

	assert.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3, "shared": 1}, dst)
}

func TestMergeIsAssociativeForLastWins(t *testing.T) {
	a := map[string]int{"x": 1, "y": 1}
	b := map[string]int{"y": 2, "z": 2}
	c := map[string]int{"z": 3}

	left := Merge(Merge(a, b), c)
	right := Merge(a, Merge(b, c))
	assert.Equal(t, left, right)
}

func BenchmarkMerge_TwoMaps(b *testing.B) {
	a := makeBenchMap(1024)
	c := makeBenchMap(1024)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Merge(a, c)
	}
}

func BenchmarkMerge_FiveMaps(b *testing.B) {
	ms := []map[int]int{
		makeBenchMap(256),
		makeBenchMap(256),
		makeBenchMap(256),
		makeBenchMap(256),
		makeBenchMap(256),
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Merge(ms...)
	}
}

func BenchmarkMergeFunc_Sum(b *testing.B) {
	x := makeBenchMap(1024)
	y := makeBenchMap(1024)
	add := func(o, n int) int { return o + n }
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MergeFunc(add, x, y)
	}
}
