package maps

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
