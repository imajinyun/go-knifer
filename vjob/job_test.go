package vjob

import (
	"context"
	"errors"
	"reflect"
	"sort"
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
	sort.Ints(seen)
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

func TestFacadeRunAndErrors_BitsUT(t *testing.T) {
	if err := Run(context.Background(), nil); !errors.Is(err, ErrNilJob) {
		t.Fatalf("Run(nil) error = %v, want ErrNilJob", err)
	}

	var ranges []string
	j := NewSlice(func(ctx context.Context, start, end int) (Merge, error) {
		ranges = append(ranges, "single")
		return nil, nil
	}, 3)
	if err := Run(context.Background(), j); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if want := []string{"single"}; !reflect.DeepEqual(ranges, want) {
		t.Fatalf("ranges = %v, want %v", ranges, want)
	}
}

func TestFacadeMapKeys_BitsUT(t *testing.T) {
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
	sort.Strings(keys)
	if want := []string{"a", "b"}; !reflect.DeepEqual(keys, want) {
		t.Fatalf("keys = %v, want %v", keys, want)
	}
}

func TestFacadeNewMapE_BitsUT(t *testing.T) {
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
	if err := Run(context.Background(), j); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	sort.Ints(keys)
	if want := []int{1, 2}; !reflect.DeepEqual(keys, want) {
		t.Fatalf("keys = %v, want %v", keys, want)
	}

	if _, err := NewMapE(123, data); !errors.Is(err, ErrInvalidMapJob) {
		t.Fatalf("NewMapE() error = %v, want ErrInvalidMapJob", err)
	}
}
