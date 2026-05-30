package cron

import (
	"sync"
)

// TaskExecutor is aligned with hutool TaskExecutor and executes a single CronTask.
type TaskExecutor struct {
	scheduler *Scheduler
	task      *CronTask
}

// Task returns the inner Task wrapped by the CronTask.
func (e *TaskExecutor) Task() Task { return e.task.Raw() }

// CronTask returns the wrapped CronTask.
func (e *TaskExecutor) CronTask() *CronTask { return e.task }

// run executes the task, triggers listener callbacks, and notifies the manager.
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

// taskExecutorManager is aligned with hutool TaskExecutorManager.
type taskExecutorManager struct {
	scheduler *Scheduler
	mu        sync.Mutex
	executors []*TaskExecutor
}

func newTaskExecutorManager(s *Scheduler) *taskExecutorManager {
	return &taskExecutorManager{scheduler: s}
}

// spawn creates a TaskExecutor for a CronTask and submits it to the scheduler executor.
func (m *taskExecutorManager) spawn(task *CronTask) *TaskExecutor {
	e := &TaskExecutor{scheduler: m.scheduler, task: task}
	m.mu.Lock()
	m.executors = append(m.executors, e)
	m.mu.Unlock()
	m.scheduler.submit(e.run)
	return e
}

// completed removes an executor from the running list after task execution completes.
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

// taskLauncher is aligned with hutool TaskLauncher; one launcher is submitted per firing instant.
type taskLauncher struct {
	scheduler *Scheduler
	millis    int64
}

func (l *taskLauncher) run() {
	defer l.scheduler.launcherMgr.completed(l)
	l.scheduler.taskTable.executeIfMatch(l.scheduler, l.millis)
}

// taskLauncherManager is aligned with hutool TaskLauncherManager.
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
