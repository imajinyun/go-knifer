package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
