package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRWMutexMapOperations(t *testing.T) {
	var zero RWMutexMap[string, int]

	value, ok := zero.Load("missing")
	assert.False(t, ok)
	assert.Zero(t, value)
	assert.Zero(t, zero.Len())
	assert.Empty(t, zero.Keys())
	assert.Empty(t, zero.Values())

	zero.Store("a", 1)
	zero.Store("b", 2)

	value, ok = zero.Load("a")
	assert.True(t, ok)
	assert.Equal(t, 1, value)
	assert.Equal(t, 2, zero.Len())
	assert.ElementsMatch(t, []string{"a", "b"}, zero.Keys())
	assert.ElementsMatch(t, []int{1, 2}, zero.Values())

	actual, loaded := zero.LoadOrStore("a", 99)
	assert.True(t, loaded)
	assert.Equal(t, 1, actual)

	actual, loaded = zero.LoadOrStore("c", 3)
	assert.False(t, loaded)
	assert.Equal(t, 3, actual)

	zero.Delete("b")
	_, ok = zero.Load("b")
	assert.False(t, ok)
	assert.Equal(t, 2, zero.Len())
}

func TestRWMutexMapConstructorAndRange(t *testing.T) {
	m := NewRWMutexMap[string, int]()
	assert.NotNil(t, m)

	m.Store("a", 1)
	m.Store("b", 2)

	visited := map[string]int{}
	m.Range(func(key string, value int) bool {
		visited[key] = value
		return true
	})
	assert.Equal(t, map[string]int{"a": 1, "b": 2}, visited)

	visited = map[string]int{}
	m.Range(func(key string, value int) bool {
		visited[key] = value
		return false
	})
	assert.Len(t, visited, 1)
}
