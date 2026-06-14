package job

import (
	"context"
	"reflect"
	"testing"
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
