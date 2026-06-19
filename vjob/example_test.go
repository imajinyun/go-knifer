package vjob_test

import (
	"context"
	"fmt"

	"github.com/imajinyun/go-knifer/vjob"
)

func ExampleNewBatch() {
	results := make([]string, 0)
	batch := vjob.NewBatch(func(_ context.Context, vals []int) (vjob.Merge, error) {
		for _, v := range vals {
			results = append(results, fmt.Sprintf("processed:%d", v))
		}
		return nil, nil
	}, []int{1, 2, 3})
	_ = vjob.Run(context.Background(), batch)
	fmt.Println(results)
	// Output: [processed:1 processed:2 processed:3]
}

func ExampleNewSlice() {
	visited := make([]int, 0)
	job := vjob.NewSlice(func(_ context.Context, start, end int) (vjob.Merge, error) {
		for i := start; i < end; i++ {
			visited = append(visited, i)
		}
		return nil, nil
	}, 3)

	_ = vjob.Run(context.Background(), job)
	fmt.Println(visited)
	// Output: [0 1 2]
}

func ExampleNewBatchSingle() {
	seen := make([]string, 0)
	batch := vjob.NewBatchSingle(func(_ context.Context, value string) (vjob.Merge, error) {
		seen = append(seen, value)
		return nil, nil
	}, []string{"a", "b"}).WithMaxConcurrency(1)

	_ = vjob.Run(context.Background(), batch)
	fmt.Println(seen)
	// Output: [a b]
}
