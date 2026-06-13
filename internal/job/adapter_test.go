package job

import (
	"context"
	"errors"
	"reflect"
	"sort"
	"testing"

	knifer "github.com/imajinyun/go-knifer"
)

func TestBatchSliceAndMapAdapters_BitsUT(t *testing.T) {
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
		sort.Ints(seen)
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
		sort.Strings(keys)
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
		sort.Ints(keys)
		if want := []int{1, 2}; !reflect.DeepEqual(keys, want) {
			t.Fatalf("keys = %v, want %v", keys, want)
		}
	})
}

func TestAdaptersValidateRangesAndInputs_BitsUT(t *testing.T) {
	t.Run("invalid slice range", func(t *testing.T) {
		_, err := NewSlice(func(ctx context.Context, start, end int) (Merge, error) { return nil, nil }, 1).Run(context.Background(), 1, 2)
		if !errors.Is(err, ErrInvalidRange) {
			t.Fatalf("Slice.Run() error = %v, want ErrInvalidRange", err)
		}
		if !errors.Is(err, knifer.ErrCodeInvalidInput) {
			t.Fatalf("Slice.Run() error = %v, want ErrCodeInvalidInput", err)
		}
	})

	t.Run("invalid batch range", func(t *testing.T) {
		_, err := NewBatch(func(ctx context.Context, vals []int) (Merge, error) { return nil, nil }, []int{1}).Run(context.Background(), -1, 1)
		if !errors.Is(err, ErrInvalidRange) {
			t.Fatalf("Batch.Run() error = %v, want ErrInvalidRange", err)
		}
	})

	t.Run("invalid reflect map input", func(t *testing.T) {
		tests := []struct {
			name string
			run  any
			data any
		}{
			{name: "run is not func", run: 123, data: map[string]int{"a": 1}},
			{name: "invalid signature", run: func(context.Context, string) error { return nil }, data: map[string]int{"a": 1}},
			{name: "invalid return", run: func(context.Context, string) (error, error) { return nil, nil }, data: map[string]int{"a": 1}},
			{name: "data is not map", run: func(context.Context, string) (Merge, error) { return nil, nil }, data: []string{"a"}},
			{name: "key type mismatch", run: func(context.Context, int) (Merge, error) { return nil, nil }, data: map[string]int{"a": 1}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if _, err := NewMapE(tt.run, tt.data); !errors.Is(err, ErrInvalidMapJob) {
					t.Fatalf("NewMapE() error = %v, want ErrInvalidMapJob", err)
				}
				expectPanic(t, func() { NewMap(tt.run, tt.data) })
			})
		}
	})
}
