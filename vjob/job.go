package vjob

import (
	"context"

	jobimpl "github.com/imajinyun/go-knifer/internal/job"
)

var (
	// ErrNilJob indicates that a nil Sliceable job was passed to a runner.
	ErrNilJob = jobimpl.ErrNilJob
	// ErrInvalidRange indicates that a Run call received an invalid half-open range.
	ErrInvalidRange = jobimpl.ErrInvalidRange
)

// Merge is called serially by the scheduler after a shard succeeds.
type Merge = jobimpl.Merge

// Sliceable describes work that can be split by half-open index ranges.
type Sliceable = jobimpl.Sliceable

// Options controls scheduling behavior. The zero value is valid.
type Options = jobimpl.Options

// Slice adapts index ranges to the Sliceable interface.
type Slice = jobimpl.Slice

// Batch adapts a typed slice to the Sliceable interface.
type Batch[T any] = jobimpl.Batch[T]

// Run executes job with the default Options.
func Run(ctx context.Context, job Sliceable) error { return jobimpl.Run(ctx, job) }

// RunWith executes job with explicit scheduling options.
func RunWith(ctx context.Context, job Sliceable, opts Options) error {
	return jobimpl.RunWith(ctx, job, opts)
}

// NewSlice creates a range-based job.
func NewSlice(run func(context.Context, int, int) (Merge, error), length int) *Slice {
	return jobimpl.NewSlice(run, length)
}

// NewSliceSingle creates a job that processes one index per shard.
func NewSliceSingle(run func(context.Context, int) (Merge, error), length int) *Slice {
	return jobimpl.NewSliceSingle(run, length)
}

// NewBatch creates a typed slice job.
func NewBatch[T any](run func(context.Context, []T) (Merge, error), vals []T) *Batch[T] {
	return jobimpl.NewBatch(run, vals)
}

// NewBatchSingle creates a typed slice job that processes one item per shard.
func NewBatchSingle[T any](run func(context.Context, T) (Merge, error), vals []T) *Batch[T] {
	return jobimpl.NewBatchSingle(run, vals)
}

// NewMap creates a single-item job over map keys.
func NewMap(run any, m any) *Slice { return jobimpl.NewMap(run, m) }

// NewMapKeys creates a single-item job over typed map keys.
func NewMapKeys[K comparable, V any](run func(context.Context, K) (Merge, error), m map[K]V) *Batch[K] {
	return jobimpl.NewMapKeys(run, m)
}
