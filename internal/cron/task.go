package cron

// Task 对应 hutool 的 cn.hutool.cron.task.Task，定时任务接口。
type Task interface {
	Execute()
}

// TaskFunc 函数式 Task 适配器，对应 RunnableTask。
type TaskFunc func()

// Execute 实现 Task 接口。
func (f TaskFunc) Execute() {
	if f != nil {
		f()
	}
}

// CronTask 对应 hutool 的 CronTask，封装 id、Pattern 与原始 Task。
type CronTask struct {
	id      string
	pattern *Pattern
	raw     Task
}

// NewCronTask 创建 CronTask。
func NewCronTask(id string, pattern *Pattern, task Task) *CronTask {
	return &CronTask{id: id, pattern: pattern, raw: task}
}

// Execute 委托给原始 Task。
func (t *CronTask) Execute() {
	if t.raw != nil {
		t.raw.Execute()
	}
}

// ID 返回任务 ID。
func (t *CronTask) ID() string { return t.id }

// Pattern 返回 Cron 表达式对象。
func (t *CronTask) Pattern() *Pattern { return t.pattern }

// SetPattern 更新表达式。
func (t *CronTask) SetPattern(p *Pattern) { t.pattern = p }

// Raw 返回包装的原始 Task。
func (t *CronTask) Raw() Task { return t.raw }
