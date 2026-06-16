package job

import (
	"context"
	"reflect"
	"slices"
	"testing"
)

func TestBatchSliceAdapters_BitsUT(t *testing.T) {
	t.Run("batch", func(t *testing.T) {
		var seen []int
		j := NewBatch(func(ctx context.Context, vals []int) (Merge, error) {
			copied := append([]int(nil), vals...)
			return func() error {
				seen = append(seen, copied...)
				return nil
			}, nil
		}, []int{1, 2, 3}).WithBatchSize(2)

		if err := RunWith(context.Background(), j, j.Options); err != nil {
			t.Fatalf("RunWith() error = %v", err)
		}
		slices.Sort(seen)
		if want := []int{1, 2, 3}; !reflect.DeepEqual(seen, want) {
			t.Fatalf("seen = %v, want %v", seen, want)
		}
	})

	t.Run("single", func(t *testing.T) {
		var sum int
		j := NewBatchSingle(func(ctx context.Context, v int) (Merge, error) {
			return func() error {
				sum += v
				return nil
			}, nil
		}, []int{1, 2, 3})

		if err := RunWith(context.Background(), j, j.Options); err != nil {
			t.Fatalf("RunWith() error = %v", err)
		}
		if sum != 6 {
			t.Fatalf("sum = %d, want 6", sum)
		}
	})
}
