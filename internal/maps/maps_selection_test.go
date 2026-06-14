package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
