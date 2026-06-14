package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromPairs(t *testing.T) {
	got := FromPairs(
		Pair[string, int]{Key: "a", Value: 1},
		Pair[string, int]{Key: "b", Value: 2},
		Pair[string, int]{Key: "a", Value: 3},
	)
	assert.Equal(t, map[string]int{"a": 3, "b": 2}, got)
}
