package cron

// TaskListener is aligned with hutool cn.hutool.cron.listener.TaskListener.
type TaskListener interface {
	OnStart(executor *TaskExecutor)
	OnSucceeded(executor *TaskExecutor)
	OnFailed(executor *TaskExecutor, err any)
}

// SimpleTaskListener is aligned with hutool SimpleTaskListener and provides no-op implementations.
type SimpleTaskListener struct{}

// OnStart is a no-op default implementation.
func (SimpleTaskListener) OnStart(*TaskExecutor) {}

// OnSucceeded is a no-op default implementation.
func (SimpleTaskListener) OnSucceeded(*TaskExecutor) {}

// OnFailed is a no-op default implementation.
func (SimpleTaskListener) OnFailed(*TaskExecutor, any) {}
