package vjob

import (
	"context"
	"reflect"
	"slices"
	"testing"
)

func TestFacadeRunWithBatch_BitsUT(t *testing.T) {
	var seen []int
	j := NewBatch(func(ctx context.Context, vals []int) (Merge, error) {
		copied := append([]int(nil), vals...)
		return func() error {
			seen = append(seen, copied...)
			return nil
		}, nil
	}, []int{3, 1, 2}).WithBatchSize(2).WithMaxConcurrency(2)

	if err := RunWith(context.Background(), j, j.Options); err != nil {
		t.Fatalf("RunWith() error = %v", err)
	}
	slices.Sort(seen)
	if want := []int{1, 2, 3}; !reflect.DeepEqual(seen, want) {
		t.Fatalf("seen = %v, want %v", seen, want)
	}
}

func TestFacadeRunUsesBatchOptions_BitsUT(t *testing.T) {
	var sizes []int
	j := NewBatch(func(ctx context.Context, vals []int) (Merge, error) {
		sizes = append(sizes, len(vals))
		return nil, nil
	}, []int{1, 2, 3, 4, 5}).WithBatchSize(2)

	if err := Run(context.Background(), j); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if want := []int{2, 2, 1}; !reflect.DeepEqual(sizes, want) {
		t.Fatalf("sizes = %v, want %v", sizes, want)
	}
}
