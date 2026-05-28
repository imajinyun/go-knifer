package cron

import (
	"sync"
)

// TaskExecutor 对应 hutool 的 TaskExecutor，负责执行单个 CronTask。
type TaskExecutor struct {
	scheduler *Scheduler
	task      *CronTask
}

// Task 返回包装的 CronTask 内部 Task。
func (e *TaskExecutor) Task() Task { return e.task.Raw() }

// CronTask 返回 CronTask。
func (e *TaskExecutor) CronTask() *CronTask { return e.task }

// run 执行任务，触发 listener 回调，并通知 manager。
func (e *TaskExecutor) run() {
	defer e.scheduler.executorMgr.completed(e)
	e.scheduler.listenerMgr.notifyStart(e)
	defer func() {
		if r := recover(); r != nil {
			e.scheduler.listenerMgr.notifyFailed(e, r)
		}
	}()
	e.task.Execute()
	e.scheduler.listenerMgr.notifySucceeded(e)
}

// taskExecutorManager 对应 hutool 的 TaskExecutorManager。
type taskExecutorManager struct {
	scheduler *Scheduler
	mu        sync.Mutex
	executors []*TaskExecutor
}

func newTaskExecutorManager(s *Scheduler) *taskExecutorManager {
	return &taskExecutorManager{scheduler: s}
}

// spawn 根据 CronTask 创建 TaskExecutor 并提交到调度器线程池执行。
func (m *taskExecutorManager) spawn(task *CronTask) *TaskExecutor {
	e := &TaskExecutor{scheduler: m.scheduler, task: task}
	m.mu.Lock()
	m.executors = append(m.executors, e)
	m.mu.Unlock()
	m.scheduler.submit(e.run)
	return e
}

// completed 任务执行完毕后从执行列表中移除。
func (m *taskExecutorManager) completed(e *TaskExecutor) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, x := range m.executors {
		if x == e {
			m.executors = append(m.executors[:i], m.executors[i+1:]...)
			return
		}
	}
}

// taskLauncher 对应 hutool 的 TaskLauncher，每个调度时刻投递一个 launcher。
type taskLauncher struct {
	scheduler *Scheduler
	millis    int64
}

func (l *taskLauncher) run() {
	defer l.scheduler.launcherMgr.completed(l)
	l.scheduler.taskTable.executeIfMatch(l.scheduler, l.millis)
}

// taskLauncherManager 对应 hutool 的 TaskLauncherManager。
type taskLauncherManager struct {
	scheduler *Scheduler
	mu        sync.Mutex
	launchers []*taskLauncher
}

func newTaskLauncherManager(s *Scheduler) *taskLauncherManager {
	return &taskLauncherManager{scheduler: s}
}

func (m *taskLauncherManager) spawn(millis int64) *taskLauncher {
	l := &taskLauncher{scheduler: m.scheduler, millis: millis}
	m.mu.Lock()
	m.launchers = append(m.launchers, l)
	m.mu.Unlock()
	m.scheduler.submit(l.run)
	return l
}

func (m *taskLauncherManager) completed(l *taskLauncher) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, x := range m.launchers {
		if x == l {
			m.launchers = append(m.launchers[:i], m.launchers[i+1:]...)
			return
		}
	}
}
