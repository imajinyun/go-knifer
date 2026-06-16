package maps

import (
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestToSlice(t *testing.T) {
	in := map[string]int{"a": 1, "b": 2}

	got := ToSlice(in, func(k string, v int) string {
		return k + strconv.Itoa(v)
	})
	slices.Sort(got)
	assert.Equal(t, []string{"a1", "b2"}, got)

	empty := ToSlice(map[string]int{}, func(k string, v int) string {
		return k + strconv.Itoa(v)
	})
	assert.NotNil(t, empty)
	assert.Empty(t, empty)

	nilMap := ToSlice(map[string]int(nil), func(k string, v int) string {
		return k + strconv.Itoa(v)
	})
	assert.NotNil(t, nilMap)
	assert.Empty(t, nilMap)
}
