package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeWithOverwrite(t *testing.T) {
	dst := map[string]int{"a": 1, "shared": 1}
	src1 := map[string]int{"b": 2, "shared": 2}
	src2 := map[string]int{"c": 3, "shared": 3}

	MergeWithOverwrite(dst, src1, nil, src2)

	assert.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3, "shared": 3}, dst)
	assert.Equal(t, map[string]int{"b": 2, "shared": 2}, src1, "source maps must not be mutated")
}

func TestMergeWithoutOverwrite(t *testing.T) {
	dst := map[string]int{"a": 1, "shared": 1}
	src1 := map[string]int{"b": 2, "shared": 2}
	src2 := map[string]int{"b": 20, "c": 3, "shared": 3}

	MergeWithoutOverwrite(dst, src1, nil, src2)

	assert.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3, "shared": 1}, dst)
}
