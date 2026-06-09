package maps

import (
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func sortedStrings(s []string) []string {
	out := append([]string(nil), s...)
	sort.Strings(out)
	return out
}

// snapshot deep-copies a map to detect mutation of the input.
func snapshot[K comparable, V any](m map[K]V) map[K]V {
	out := make(map[K]V, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

// ---------------------------------------------------------------------------
// Construction
// ---------------------------------------------------------------------------

func TestNewAndNewWithCap(t *testing.T) {
	m := New[string, int]()
	assert.NotNil(t, m)
	assert.Empty(t, m)

	m2 := NewWithCap[string, int](128)
	assert.NotNil(t, m2)
	assert.Empty(t, m2)

	// negative hint is normalized to 0
	m3 := NewWithCap[string, int](-5)
	assert.NotNil(t, m3)
}

func TestOf(t *testing.T) {
	m := Of[string, int]("a", 1, "b", 2, "a", 3) // last wins
	assert.Equal(t, map[string]int{"a": 3, "b": 2}, m)

	empty := Of[string, int]()
	assert.NotNil(t, empty)
	assert.Empty(t, empty)
}

func TestOf_PanicOnOddArgs(t *testing.T) {
	t.Run("odd args panics", func(t *testing.T) {
		args := []any{"a", 1, "b"}
		assert.PanicsWithValue(t, "maps.Of: odd number of arguments", func() {
			Of[string, int](args...)
		})
	})
}

func TestOfE(t *testing.T) {
	tests := []struct {
		name    string
		args    []any
		want    map[string]int
		wantErr bool
	}{
		{
			name: "builds map",
			args: []any{"a", 1, "b", 2, "a", 3},
			want: map[string]int{"a": 3, "b": 2},
		},
		{
			name:    "rejects odd args",
			args:    []any{"a", 1, "b"},
			wantErr: true,
		},
		{
			name:    "rejects invalid key type",
			args:    []any{1, 1},
			wantErr: true,
		},
		{
			name:    "rejects invalid value type",
			args:    []any{"a", "bad"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OfE[string, int](tt.args...)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFromPairs(t *testing.T) {
	got := FromPairs(
		Pair[string, int]{Key: "a", Value: 1},
		Pair[string, int]{Key: "b", Value: 2},
		Pair[string, int]{Key: "a", Value: 3},
	)
	assert.Equal(t, map[string]int{"a": 3, "b": 2}, got)
}

func TestOrEmpty(t *testing.T) {
	var nilMap map[string]int
	got := OrEmpty(nilMap)
	assert.NotNil(t, got)
	assert.Empty(t, got)

	src := map[string]int{"x": 1}
	returned := OrEmpty(src)
	returned["y"] = 2
	assert.Equal(t, 2, src["y"], "OrEmpty should return the original non-nil map")
}

// ---------------------------------------------------------------------------
// Predicates
// ---------------------------------------------------------------------------

func TestIsEmptyAndIsNotEmpty(t *testing.T) {
	var nilMap map[int]int
	assert.True(t, IsEmpty(nilMap))
	assert.False(t, IsNotEmpty(nilMap))

	assert.True(t, IsEmpty(map[int]int{}))
	assert.False(t, IsNotEmpty(map[int]int{}))

	assert.False(t, IsEmpty(map[int]int{1: 1}))
	assert.True(t, IsNotEmpty(map[int]int{1: 1}))
}

func TestContainsKey(t *testing.T) {
	m := map[string]int{"a": 1}
	assert.True(t, ContainsKey(m, "a"))
	assert.False(t, ContainsKey(m, "z"))
}

func TestContainsValue(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	assert.True(t, ContainsValue(m, 1))
	assert.False(t, ContainsValue(m, 99))
}

func TestSome(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	assert.True(t, Some(m, func(_ string, v int) bool { return v > 2 }))
	assert.False(t, Some(m, func(_ string, v int) bool { return v > 100 }))
	assert.False(t, Some(map[string]int{}, func(_ string, _ int) bool { return true }))
}

func TestEvery(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	assert.True(t, Every(m, func(_ string, v int) bool { return v > 0 }))
	assert.False(t, Every(m, func(_ string, v int) bool { return v > 1 }))
	// empty map → vacuously true
	assert.True(t, Every(map[string]int{}, func(_ string, _ int) bool { return false }))
}

// ---------------------------------------------------------------------------
// Lookup
// ---------------------------------------------------------------------------

func TestGetAndGetOr(t *testing.T) {
	m := map[string]int{"a": 1}
	assert.Equal(t, 1, Get(m, "a"))
	assert.Equal(t, 0, Get(m, "missing"))
	assert.Equal(t, 1, GetOr(m, "a", 99))
	assert.Equal(t, 99, GetOr(m, "missing", 99))
}

func TestGetAny(t *testing.T) {
	headers := map[string]string{"X-Username": "alice"}
	v, ok := GetAny(headers, "X-User", "X-Username", "User")
	assert.True(t, ok)
	assert.Equal(t, "alice", v)

	v2, ok2 := GetAny(headers, "missing-1", "missing-2")
	assert.False(t, ok2)
	assert.Equal(t, "", v2)
}

func TestFindAndFindKey(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	k, v, ok := Find(m, func(_ string, v int) bool { return v == 2 })
	assert.True(t, ok)
	assert.Equal(t, "b", k)
	assert.Equal(t, 2, v)

	_, _, ok = Find(m, func(_ string, v int) bool { return v < 0 })
	assert.False(t, ok)

	fk, ok := FindKey(m, func(v int) bool { return v == 3 })
	assert.True(t, ok)
	assert.Equal(t, "c", fk)

	_, ok = FindKey(m, func(v int) bool { return v > 999 })
	assert.False(t, ok)
}

// ---------------------------------------------------------------------------
// Collection views
// ---------------------------------------------------------------------------

func TestKeysAndValues(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	assert.ElementsMatch(t, []string{"a", "b", "c"}, Keys(m))
	assert.ElementsMatch(t, []int{1, 2, 3}, Values(m))

	// nil-safe
	assert.Empty(t, Keys[string, int](nil))
	assert.Empty(t, Values[string, int](nil))
}

func TestSortedKeysAndValues(t *testing.T) {
	m := map[string]int{"c": 3, "a": 1, "b": 2}
	assert.Equal(t, []string{"a", "b", "c"}, SortedKeys(m))
	assert.Equal(t, []int{1, 2, 3}, SortedValues(m))

	descending := SortedKeysFunc(m, func(a, b string) bool { return a > b })
	assert.Equal(t, []string{"c", "b", "a"}, descending)
}

func TestKeysOf(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 1}
	got := sortedStrings(KeysOf(m, 1))
	assert.Equal(t, []string{"a", "c"}, got)

	assert.Empty(t, KeysOf(m, 99))
}

// ---------------------------------------------------------------------------
// Transformation
// ---------------------------------------------------------------------------

func TestMap(t *testing.T) {
	in := map[string]int{"a": 1, "b": 2}
	out := Map(in, func(k string, v int) (string, string) {
		return strings.ToUpper(k), strconv.Itoa(v)
	})
	assert.Equal(t, map[string]string{"A": "1", "B": "2"}, out)
}

func TestMapKeysAndMapValues(t *testing.T) {
	in := map[string]int{"a": 1, "b": 2}

	mk := MapKeys(in, func(k string, _ int) string { return strings.ToUpper(k) })
	assert.Equal(t, map[string]int{"A": 1, "B": 2}, mk)

	mv := MapValues(in, func(_ string, v int) int { return v * 10 })
	assert.Equal(t, map[string]int{"a": 10, "b": 20}, mv)
}

func TestFilterAndReject(t *testing.T) {
	in := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	keep := Filter(in, func(_ string, v int) bool { return v%2 == 0 })
	drop := Reject(in, func(_ string, v int) bool { return v%2 == 0 })

	assert.Equal(t, map[string]int{"b": 2, "d": 4}, keep)
	assert.Equal(t, map[string]int{"a": 1, "c": 3}, drop)
}

func TestFilterKeysAndFilterValues(t *testing.T) {
	in := map[string]int{"alpha": 1, "beta": 2, "gamma": 3}
	fk := FilterKeys(in, func(k string) bool { return strings.HasPrefix(k, "a") })
	assert.Equal(t, map[string]int{"alpha": 1}, fk)

	fv := FilterValues(in, func(v int) bool { return v > 1 })
	assert.Equal(t, map[string]int{"beta": 2, "gamma": 3}, fv)
}

func TestPartition(t *testing.T) {
	in := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	yes, no := Partition(in, func(_ string, v int) bool { return v >= 3 })
	assert.Equal(t, map[string]int{"c": 3, "d": 4}, yes)
	assert.Equal(t, map[string]int{"a": 1, "b": 2}, no)
}

func TestForEach(t *testing.T) {
	in := map[string]int{"a": 1, "b": 2}
	sum := 0
	keys := make([]string, 0, 2)
	ForEach(in, func(k string, v int) {
		sum += v
		keys = append(keys, k)
	})
	assert.Equal(t, 3, sum)
	assert.ElementsMatch(t, []string{"a", "b"}, keys)
}

// ---------------------------------------------------------------------------
// Aggregation
// ---------------------------------------------------------------------------

func TestReduce(t *testing.T) {
	in := map[string]int{"a": 1, "b": 2, "c": 3}
	sum := Reduce(in, 0, func(acc int, _ string, v int) int { return acc + v })
	assert.Equal(t, 6, sum)

	concat := Reduce(in, "", func(acc string, k string, _ int) string { return acc + k })
	// order is non-deterministic; just check length & character set
	assert.Len(t, concat, 3)
}

func TestGroupBy(t *testing.T) {
	type emp struct {
		Name string
		Dept string
	}
	items := []emp{{"a", "eng"}, {"b", "eng"}, {"c", "ops"}}
	got := GroupBy(items, func(e emp) string { return e.Dept })
	require.Len(t, got, 2)
	assert.ElementsMatch(t, []emp{{"a", "eng"}, {"b", "eng"}}, got["eng"])
	assert.Equal(t, []emp{{"c", "ops"}}, got["ops"])
}

func TestCountBy(t *testing.T) {
	logs := []string{"GET", "POST", "GET", "GET", "POST"}
	got := CountBy(logs, func(s string) string { return s })
	assert.Equal(t, map[string]int{"GET": 3, "POST": 2}, got)
}

// ---------------------------------------------------------------------------
// Set algebra
// ---------------------------------------------------------------------------

func TestInverse(t *testing.T) {
	in := map[string]int{"a": 1, "b": 2}
	inv := Inverse(in)
	assert.Equal(t, map[int]string{1: "a", 2: "b"}, inv)
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

// ---------------------------------------------------------------------------
// Selection
// ---------------------------------------------------------------------------

func TestPick(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	got := Pick(m, "a", "c", "missing")
	assert.Equal(t, map[string]int{"a": 1, "c": 3}, got)

	assert.Empty(t, Pick(m))
	assert.Empty(t, Pick[string, int](nil, "a"))
}

func TestOmit(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	got := Omit(m, "b", "missing")
	assert.Equal(t, map[string]int{"a": 1, "c": 3}, got)

	assert.Equal(t, m, Omit(m))
}

// ---------------------------------------------------------------------------
// Mutation helpers
// ---------------------------------------------------------------------------

func TestClear(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	Clear(m)
	assert.Empty(t, m)
	assert.NotNil(t, m) // still the same allocated map
}

func TestUpdate(t *testing.T) {
	dst := map[string]int{"a": 1}
	src := map[string]int{"a": 10, "b": 2}
	got := Update(dst, src)
	assert.Equal(t, map[string]int{"a": 10, "b": 2}, dst)
	got["c"] = 3
	assert.Equal(t, 3, dst["c"], "Update should return dst for chaining")

	// nil dst is allocated
	got2 := Update[string, int](nil, src)
	assert.NotNil(t, got2)
	assert.Equal(t, src, got2)
}

func TestClone(t *testing.T) {
	m := map[string]int{"a": 1}
	c := Clone(m)
	assert.Equal(t, m, c)

	c["a"] = 999
	assert.Equal(t, 1, m["a"], "Clone must not share storage with the input")

	// nil input → empty non-nil
	cn := Clone[string, int](nil)
	assert.NotNil(t, cn)
	assert.Empty(t, cn)
}

// ---------------------------------------------------------------------------
// Comparison
// ---------------------------------------------------------------------------

func TestEqual(t *testing.T) {
	a := map[string]int{"a": 1, "b": 2}
	b := map[string]int{"b": 2, "a": 1}
	c := map[string]int{"a": 1}
	d := map[string]int{"a": 1, "b": 99}

	assert.True(t, Equal(a, b))
	assert.False(t, Equal(a, c))
	assert.False(t, Equal(a, d))
	assert.True(t, Equal[string, int](nil, nil))
	assert.True(t, Equal(map[string]int{}, map[string]int{}))
}

func TestEqualFunc(t *testing.T) {
	a := map[string]int{"a": 1, "b": 2}
	b := map[string]string{"a": "1", "b": "2"}

	got := EqualFunc(a, b, func(x int, y string) bool {
		return strconv.Itoa(x) == y
	})
	assert.True(t, got)

	notMatching := EqualFunc(a, b, func(x int, _ string) bool { return x < 0 })
	assert.False(t, notMatching)
}

// ---------------------------------------------------------------------------
// Property-style sanity: shape consistency between functions
// ---------------------------------------------------------------------------

func TestKeysValuesShape(t *testing.T) {
	m := map[int]int{1: 10, 2: 20, 3: 30}
	keys := SortedKeys(m)
	values := SortedValues(m)
	require.Len(t, keys, len(m))
	require.Len(t, values, len(m))
	for i, k := range keys {
		assert.Equal(t, m[k], values[i])
	}
}

func TestMergeIsAssociativeForLastWins(t *testing.T) {
	a := map[string]int{"x": 1, "y": 1}
	b := map[string]int{"y": 2, "z": 2}
	c := map[string]int{"z": 3}

	left := Merge(Merge(a, b), c)
	right := Merge(a, Merge(b, c))
	assert.Equal(t, left, right)
}

// ---------------------------------------------------------------------------
// Benchmarks
// ---------------------------------------------------------------------------

func makeBenchMap(n int) map[int]int {
	m := make(map[int]int, n)
	for i := 0; i < n; i++ {
		m[i] = i
	}
	return m
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

func BenchmarkFilter(b *testing.B) {
	m := makeBenchMap(4096)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Filter(m, func(_ int, v int) bool { return v%2 == 0 })
	}
}

func BenchmarkSortedKeys(b *testing.B) {
	m := makeBenchMap(4096)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SortedKeys(m)
	}
}

func BenchmarkClone(b *testing.B) {
	m := makeBenchMap(4096)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Clone(m)
	}
}
