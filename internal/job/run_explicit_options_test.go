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
