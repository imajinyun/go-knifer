package vjob

import (
	"context"
	"reflect"
	"slices"
	"testing"
)

func TestZeroValueBatchEnsureInner_BitsUT(t *testing.T) {
	// A zero-value Batch must lazily build its inner implementation so that Len
	// and Run work without an explicit constructor.
	var b Batch[int]
	if got := b.Len(); got != 0 {
		t.Fatalf("zero Batch Len = %d, want 0", got)
	}
	if _, err := b.Run(context.Background(), 0, 0); err != nil {
		t.Fatalf("zero Batch Run error = %v", err)
	}
	if (Options{}) != b.JobOptions() {
		t.Fatalf("zero Batch JobOptions = %#v", b.JobOptions())
	}
}

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
