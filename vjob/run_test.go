package vjob

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

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
