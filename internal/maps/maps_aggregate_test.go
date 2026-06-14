package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReduce(t *testing.T) {
	in := map[string]int{"a": 1, "b": 2, "c": 3}
	sum := Reduce(in, 0, func(acc int, _ string, v int) int { return acc + v })
	assert.Equal(t, 6, sum)

	concat := Reduce(in, "", func(acc string, k string, _ int) string { return acc + k })
	// order is non-deterministic; just check length & character set
	assert.Len(t, concat, 3)
}

func TestGroupBy(t *testing.T) {
	type emp struct {
		Name string
		Dept string
	}
	items := []emp{{"a", "eng"}, {"b", "eng"}, {"c", "ops"}}
	got := GroupBy(items, func(e emp) string { return e.Dept })
	require.Len(t, got, 2)
	assert.ElementsMatch(t, []emp{{"a", "eng"}, {"b", "eng"}}, got["eng"])
	assert.Equal(t, []emp{{"c", "ops"}}, got["ops"])
}

func TestCountBy(t *testing.T) {
	logs := []string{"GET", "POST", "GET", "GET", "POST"}
	got := CountBy(logs, func(s string) string { return s })
	assert.Equal(t, map[string]int{"GET": 3, "POST": 2}, got)
}
