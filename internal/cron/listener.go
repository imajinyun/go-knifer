package cron

// TaskListener 对应 hutool 的 cn.hutool.cron.listener.TaskListener。
type TaskListener interface {
	OnStart(executor *TaskExecutor)
	OnSucceeded(executor *TaskExecutor)
	OnFailed(executor *TaskExecutor, err any)
}

// SimpleTaskListener 对应 hutool 的 SimpleTaskListener，提供空实现。
type SimpleTaskListener struct{}

// OnStart 默认空实现。
func (SimpleTaskListener) OnStart(*TaskExecutor) {}

// OnSucceeded 默认空实现。
func (SimpleTaskListener) OnSucceeded(*TaskExecutor) {}

// OnFailed 默认空实现。
func (SimpleTaskListener) OnFailed(*TaskExecutor, any) {}
