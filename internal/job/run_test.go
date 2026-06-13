package job

import (
	"context"
	"reflect"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRunUsesDefaultOptionsAsSingleSerialShard_BitsUT(t *testing.T) {
	var (
		ranges []string
		merged []string
	)

	j := NewSlice(func(ctx context.Context, start, end int) (Merge, error) {
		if ctx == nil {
			t.Fatal("ctx should not be nil")
		}
		ranges = append(ranges, formatRange(start, end))
		return func() error {
			merged = append(merged, formatRange(start, end))
			return nil
		}, nil
	}, 5)

	if err := Run(context.Background(), j); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if want := []string{"0:5"}; !reflect.DeepEqual(ranges, want) {
		t.Fatalf("ranges = %v, want %v", ranges, want)
	}
	if want := []string{"0:5"}; !reflect.DeepEqual(merged, want) {
		t.Fatalf("merged = %v, want %v", merged, want)
	}
}

func TestRunWithUsesExplicitOptionsAndSerialMergeOrder_BitsUT(t *testing.T) {
	var (
		mu              sync.Mutex
		ranges          []string
		merged          []string
		active, maxSeen atomic.Int32
	)

	j := NewSlice(func(ctx context.Context, start, end int) (Merge, error) {
		current := active.Add(1)
		for {
			max := maxSeen.Load()
			if current <= max || maxSeen.CompareAndSwap(max, current) {
				break
			}
		}
		defer active.Add(-1)

		if start == 0 {
			time.Sleep(20 * time.Millisecond)
		}
		mu.Lock()
		ranges = append(ranges, formatRange(start, end))
		mu.Unlock()

		return func() error {
			merged = append(merged, formatRange(start, end))
			return nil
		}, nil
	}, 5)

	err := RunWith(context.Background(), j, Options{BatchSize: 2, MaxConcurrency: 2})
	if err != nil {
		t.Fatalf("RunWith() error = %v", err)
	}

	sort.Strings(ranges)
	if want := []string{"0:2", "2:4", "4:5"}; !reflect.DeepEqual(ranges, want) {
		t.Fatalf("ranges = %v, want %v", ranges, want)
	}
	if want := []string{"0:2", "2:4", "4:5"}; !reflect.DeepEqual(merged, want) {
		t.Fatalf("merged = %v, want %v", merged, want)
	}
	if got := maxSeen.Load(); got > 2 {
		t.Fatalf("max concurrency = %d, want <= 2", got)
	}
}

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
