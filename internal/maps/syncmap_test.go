package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSyncMapOperations(t *testing.T) {
	var m SyncMap[string, int]

	value, ok := m.Load("missing")
	assert.False(t, ok)
	assert.Zero(t, value)

	m.Store("a", 1)
	value, ok = m.Load("a")
	assert.True(t, ok)
	assert.Equal(t, 1, value)

	actual, loaded := m.LoadOrStore("a", 99)
	assert.True(t, loaded)
	assert.Equal(t, 1, actual)

	actual, loaded = m.LoadOrStore("b", 2)
	assert.False(t, loaded)
	assert.Equal(t, 2, actual)

	deleted, loaded := m.LoadAndDelete("b")
	assert.True(t, loaded)
	assert.Equal(t, 2, deleted)

	deleted, loaded = m.LoadAndDelete("missing")
	assert.False(t, loaded)
	assert.Zero(t, deleted)

	m.Delete("a")
	_, ok = m.Load("a")
	assert.False(t, ok)
}

func TestSyncMapRange(t *testing.T) {
	var m SyncMap[string, int]
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
