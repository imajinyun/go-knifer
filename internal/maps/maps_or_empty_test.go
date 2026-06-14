package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrEmpty(t *testing.T) {
	var nilMap map[string]int
	got := OrEmpty(nilMap)
	assert.NotNil(t, got)
	assert.Empty(t, got)

	src := map[string]int{"x": 1}
	returned := OrEmpty(src)
	returned["y"] = 2
	assert.Equal(t, 2, src["y"], "OrEmpty should return the original non-nil map")
}
