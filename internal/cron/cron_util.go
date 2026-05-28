package cron

import "sync"

// 全局调度器，对应 hutool 的 CronUtil.scheduler。
var (
	defaultMu        sync.Mutex
	defaultScheduler = NewScheduler()
)

// DefaultScheduler 获取全局调度器。
func DefaultScheduler() *Scheduler {
	return defaultScheduler
}

// SetMatchSecond 设置全局调度器是否匹配到秒。
func SetMatchSecond(b bool) {
	defaultScheduler.SetMatchSecond(b)
}

// Schedule 在全局调度器上注册任务，返回任务 id。
func Schedule(pattern string, task Task) (string, error) {
	return defaultScheduler.Schedule(pattern, task)
}

// ScheduleFunc 在全局调度器上注册函数任务。
func ScheduleFunc(pattern string, fn func()) (string, error) {
	return defaultScheduler.ScheduleFunc(pattern, fn)
}

// ScheduleWithID 在全局调度器上注册指定 id 的任务。
func ScheduleWithID(id, pattern string, task Task) error {
	return defaultScheduler.ScheduleWithID(id, pattern, task)
}

// Remove 在全局调度器上删除任务。
func Remove(id string) bool {
	return defaultScheduler.Deschedule(id)
}

// UpdatePattern 在全局调度器上更新任务表达式。
func UpdatePattern(id, pattern string) error {
	return defaultScheduler.UpdatePattern(id, pattern)
}

// Start 启动全局调度器。
func Start() error {
	defaultMu.Lock()
	defer defaultMu.Unlock()
	return defaultScheduler.Start()
}

// Stop 停止全局调度器。
func Stop() {
	defaultMu.Lock()
	defer defaultMu.Unlock()
	defaultScheduler.Stop()
}

// Restart 重启全局调度器。
func Restart() error {
	defaultMu.Lock()
	defer defaultMu.Unlock()
	defaultScheduler.Stop()
	return defaultScheduler.Start()
}
