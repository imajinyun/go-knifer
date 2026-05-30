package cron

import "sync"

// defaultScheduler is the package-level scheduler aligned with hutool CronUtil.scheduler.
var (
	defaultMu        sync.Mutex
	defaultScheduler = NewScheduler()
)

// DefaultScheduler returns the package-level scheduler.
func DefaultScheduler() *Scheduler {
	return defaultScheduler
}

// SetMatchSecond sets whether the package-level scheduler matches seconds.
func SetMatchSecond(b bool) {
	defaultScheduler.SetMatchSecond(b)
}

// Schedule registers a task on the package-level scheduler and returns its id.
func Schedule(pattern string, task Task) (string, error) {
	return defaultScheduler.Schedule(pattern, task)
}

// ScheduleFunc registers a function task on the package-level scheduler.
func ScheduleFunc(pattern string, fn func()) (string, error) {
	return defaultScheduler.ScheduleFunc(pattern, fn)
}

// ScheduleWithID registers a task with the specified id on the package-level scheduler.
func ScheduleWithID(id, pattern string, task Task) error {
	return defaultScheduler.ScheduleWithID(id, pattern, task)
}

// Remove deletes a task from the package-level scheduler.
func Remove(id string) bool {
	return defaultScheduler.Deschedule(id)
}

// UpdatePattern updates a task expression on the package-level scheduler.
func UpdatePattern(id, pattern string) error {
	return defaultScheduler.UpdatePattern(id, pattern)
}

// Start starts the package-level scheduler.
func Start() error {
	defaultMu.Lock()
	defer defaultMu.Unlock()
	return defaultScheduler.Start()
}

// Stop stops the package-level scheduler.
func Stop() {
	defaultMu.Lock()
	defer defaultMu.Unlock()
	defaultScheduler.Stop()
}

// Restart restarts the package-level scheduler.
func Restart() error {
	defaultMu.Lock()
	defer defaultMu.Unlock()
	defaultScheduler.Stop()
	return defaultScheduler.Start()
}
