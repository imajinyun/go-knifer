package job

import (
	"context"
	"reflect"
	"slices"
	"testing"
)

func TestMapAdapters_BitsUT(t *testing.T) {
	t.Run("map keys", func(t *testing.T) {
		data := map[string]int{"b": 2, "a": 1}
		var keys []string
		j := NewMapKeys(func(ctx context.Context, key string) (Merge, error) {
			return func() error {
				keys = append(keys, key)
				return nil
			}, nil
		}, data)

		if err := RunWith(context.Background(), j, j.Options); err != nil {
			t.Fatalf("RunWith() error = %v", err)
		}
		slices.Sort(keys)
		if want := []string{"a", "b"}; !reflect.DeepEqual(keys, want) {
			t.Fatalf("keys = %v, want %v", keys, want)
		}
	})

	t.Run("reflect map", func(t *testing.T) {
		data := map[int]string{2: "b", 1: "a"}
		var keys []int
		j, err := NewMapE(func(ctx context.Context, key int) (Merge, error) {
			return func() error {
				keys = append(keys, key)
				return nil
			}, nil
		}, data)
		if err != nil {
			t.Fatalf("NewMapE() error = %v", err)
		}

		if err := RunWith(context.Background(), j, j.Options); err != nil {
			t.Fatalf("RunWith() error = %v", err)
		}
		slices.Sort(keys)
		if want := []int{1, 2}; !reflect.DeepEqual(keys, want) {
			t.Fatalf("keys = %v, want %v", keys, want)
		}
	})
}
