package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestIterators(t *testing.T) {
	in := map[string]int{"a": 1, "b": 2}

	entries := make(map[string]int, len(in))
	for key, value := range Iter(in) {
		entries[key] = value
	}
	assert.Equal(t, in, entries)

	keys := make([]string, 0, len(in))
	for key := range IterKeys(in) {
		keys = append(keys, key)
	}
	assert.ElementsMatch(t, []string{"a", "b"}, keys)

	values := make([]int, 0, len(in))
	for value := range IterValues(in) {
		values = append(values, value)
	}
	assert.ElementsMatch(t, []int{1, 2}, values)
}
