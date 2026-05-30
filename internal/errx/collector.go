// Package errx provides small error handling and panic-recovery helpers used by
// internal packages.
package errx

import (
	"context"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/sirupsen/logrus"
)

// Collector runs functions, recovers panics, logs failures, and aggregates
// returned errors. It is safe for concurrent use.
type Collector struct {
	level logrus.Level
	ctx   context.Context

	swg sync.WaitGroup
	mux sync.Mutex
	err []error
}

// NewCollector creates a Collector that logs failures at error level.
func NewCollector() *Collector {
	return &Collector{
		level: logrus.ErrorLevel,
		ctx:   context.Background(),
	}
}

// WithContext sets the context attached to log entries.
func (c *Collector) WithContext(ctx context.Context) *Collector {
	if ctx == nil {
		ctx = context.Background()
	}
	c.mux.Lock()
	defer c.mux.Unlock()
	c.ctx = ctx
	return c
}

// WithLevel sets the log level used for recovered or returned errors.
func (c *Collector) WithLevel(level logrus.Level) *Collector {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.level = level
	return c
}

// Collect stores err for the final aggregated result.
func (c *Collector) Collect(err error) {
	if err == nil {
		return
	}
	c.mux.Lock()
	c.err = append(c.err, err)
	c.mux.Unlock()
}

// Error waits for all launched functions and returns all collected errors.
func (c *Collector) Error() error {
	c.swg.Wait()
	return c.error()
}

// WaitUntil waits until all launched functions finish or duration expires.
// It returns whether all functions completed and the aggregated error, if any.
func (c *Collector) WaitUntil(duration time.Duration) (bool, error) {
	if duration <= 0 {
		return false, nil
	}
	done := make(chan struct{})
	go func() {
		c.swg.Wait()
		close(done)
	}()

	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-done:
		return true, c.error()
	case <-timer.C:
		return false, nil
	}
}

// Recover executes f in the current goroutine, recovers panics, logs failures,
// and stores non-nil errors in the collector.
func (c *Collector) Recover(f func() error, format string, args ...any) error {
	c.swg.Add(1)
	defer c.swg.Done()

	err := c.run(f, format, args...)
	c.Collect(err)
	return err
}

// GoRun executes f in a new goroutine and stores any panic or returned error.
func (c *Collector) GoRun(f func() error, format string, args ...any) {
	c.swg.Add(1)
	go func() {
		defer c.swg.Done()
		c.Collect(c.run(f, format, args...))
	}()
}

// CollectError is kept as a compatibility alias for Recover.
func (c *Collector) CollectError(f func() error, format string, args ...any) {
	_ = c.Recover(f, format, args...)
}

func (c *Collector) run(f func() error, format string, args ...any) (err error) {
	defer func() {
		if v := recover(); v != nil {
			err = multierror.Append(err, panicError(v))
		}
		if err != nil {
			c.log(err, format, args...)
		}
	}()
	if f == nil {
		return nil
	}
	return f()
}

func (c *Collector) log(err error, format string, args ...any) {
	if format == "" {
		format = "operation failed"
	}
	c.mux.Lock()
	ctx, level := c.ctx, c.level
	c.mux.Unlock()
	logrus.WithContext(ctx).
		WithError(err).
		WithField("stack", GetStack(err)).
		Logf(level, format, args...)
}

func (c *Collector) error() error {
	c.mux.Lock()
	defer c.mux.Unlock()
	if len(c.err) == 0 {
		return nil
	}
	return multierror.Append(nil, c.err...)
}

func panicError(v any) error {
	pe := &PanicError{
		Value:      v,
		StackTrace: GetStackTrace(4),
	}
	if err, ok := v.(error); ok {
		pe.Cause = err
	}
	return pe
}
