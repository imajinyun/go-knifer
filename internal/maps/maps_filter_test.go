package maps

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
