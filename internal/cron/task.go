package cron

// Task is the cron task interface aligned with hutool cn.hutool.cron.task.Task.
type Task interface {
	Execute()
}

// TaskFunc adapts a function into Task and is aligned with RunnableTask.
type TaskFunc func()

// Execute implements Task.
func (f TaskFunc) Execute() {
	if f != nil {
		f()
	}
}

// CronTask is aligned with hutool CronTask and wraps an id, Pattern, and raw Task.
type CronTask struct {
	id      string
	pattern *Pattern
	raw     Task
}

// NewCronTask creates a CronTask.
func NewCronTask(id string, pattern *Pattern, task Task) *CronTask {
	return &CronTask{id: id, pattern: pattern, raw: task}
}

// Execute delegates to the raw Task.
func (t *CronTask) Execute() {
	if t.raw != nil {
		t.raw.Execute()
	}
}

// ID returns the task ID.
func (t *CronTask) ID() string { return t.id }

// Pattern returns the cron expression object.
func (t *CronTask) Pattern() *Pattern { return t.pattern }

// SetPattern updates the expression.
func (t *CronTask) SetPattern(p *Pattern) { t.pattern = p }

// Raw returns the wrapped raw Task.
func (t *CronTask) Raw() Task { return t.raw }
