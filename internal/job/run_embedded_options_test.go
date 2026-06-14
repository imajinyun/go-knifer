package job

import (
	"context"
	"reflect"
	"sort"
	"testing"
)

func TestRunWithEmbeddedOptions_BitsUT(t *testing.T) {
	wrapped := &embeddedOptionsJob{
		Options: Options{BatchSize: 2, MaxConcurrency: 2},
		vals:    []int{1, 2, 3, 4},
	}
	if err := RunWith(context.Background(), wrapped, wrapped.Options); err != nil {
		t.Fatalf("RunWith() error = %v", err)
	}
	sort.Ints(wrapped.seen)
	if want := []int{1, 2, 3, 4}; !reflect.DeepEqual(wrapped.seen, want) {
		t.Fatalf("seen = %v, want %v", wrapped.seen, want)
	}
}

func TestRunUsesEmbeddedOptions_BitsUT(t *testing.T) {
	var ranges []string
	j := NewSlice(func(ctx context.Context, start, end int) (Merge, error) {
		ranges = append(ranges, formatRange(start, end))
		return nil, nil
	}, 5).WithBatchSize(2)

	if err := Run(context.Background(), j); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if want := []string{"0:2", "2:4", "4:5"}; !reflect.DeepEqual(ranges, want) {
		t.Fatalf("ranges = %v, want %v", ranges, want)
	}
}

type embeddedOptionsJob struct {
	Options
	vals []int
	seen []int
}

func (j *embeddedOptionsJob) Len() int { return len(j.vals) }

func (j *embeddedOptionsJob) Run(ctx context.Context, start, end int) (Merge, error) {
	_ = ctx
	batch := append([]int(nil), j.vals[start:end]...)
	return func() error {
		j.seen = append(j.seen, batch...)
		return nil
	}, nil
}
