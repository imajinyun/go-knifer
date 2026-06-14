package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOf(t *testing.T) {
	m := Of[string, int]("a", 1, "b", 2, "a", 3) // last wins
	assert.Equal(t, map[string]int{"a": 3, "b": 2}, m)

	empty := Of[string, int]()
	assert.NotNil(t, empty)
	assert.Empty(t, empty)
}

func TestOf_PanicOnOddArgs(t *testing.T) {
	t.Run("odd args panics", func(t *testing.T) {
		args := []any{"a", 1, "b"}
		assert.PanicsWithValue(t, "maps.Of: odd number of arguments", func() {
			Of[string, int](args...)
		})
	})
}

func TestOfE(t *testing.T) {
	tests := []struct {
		name    string
		args    []any
		want    map[string]int
		wantErr bool
	}{
		{
			name: "builds map",
			args: []any{"a", 1, "b", 2, "a", 3},
			want: map[string]int{"a": 3, "b": 2},
		},
		{
			name:    "rejects odd args",
			args:    []any{"a", 1, "b"},
			wantErr: true,
		},
		{
			name:    "rejects invalid key type",
			args:    []any{1, 1},
			wantErr: true,
		},
		{
			name:    "rejects invalid value type",
			args:    []any{"a", "bad"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OfE[string, int](tt.args...)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
