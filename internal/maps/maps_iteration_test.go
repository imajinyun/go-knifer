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
