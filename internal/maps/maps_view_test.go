package maps

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sortedStrings(s []string) []string {
	out := append([]string(nil), s...)
	sort.Strings(out)
	return out
}

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

func TestKeysOf(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 1}
	got := sortedStrings(KeysOf(m, 1))
	assert.Equal(t, []string{"a", "c"}, got)

	assert.Empty(t, KeysOf(m, 99))
}
