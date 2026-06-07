package cron

import (
	"context"
	"sync"
)

// TaskExecutor is aligned with the utility toolkit TaskExecutor and executes a single CronTask.
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

// taskExecutorManager is aligned with the utility toolkit TaskExecutorManager.
type taskExecutorManager struct {
	scheduler *Scheduler
	mu        sync.Mutex
	cond      *sync.Cond
	executors []*TaskExecutor
	idleCh    chan struct{}
}

func newTaskExecutorManager(s *Scheduler) *taskExecutorManager {
	idleCh := make(chan struct{})
	close(idleCh)
	m := &taskExecutorManager{scheduler: s, idleCh: idleCh}
	m.cond = sync.NewCond(&m.mu)
	return m
}

// spawn creates a TaskExecutor for a CronTask and submits it to the scheduler executor.
func (m *taskExecutorManager) spawn(task *CronTask) *TaskExecutor {
	e := &TaskExecutor{scheduler: m.scheduler, task: task}
	m.mu.Lock()
	if len(m.executors) == 0 {
		m.idleCh = make(chan struct{})
	}
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
			m.cond.Broadcast()
			if len(m.executors) == 0 {
				close(m.idleCh)
			}
			return
		}
	}
}

func (m *taskExecutorManager) runningCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.executors)
}

func (m *taskExecutorManager) wait() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for len(m.executors) > 0 {
		m.cond.Wait()
	}
}

func (m *taskExecutorManager) waitContext(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	m.mu.Lock()
	idleCh := m.idleCh
	m.mu.Unlock()
	select {
	case <-idleCh:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// taskLauncher is aligned with the utility toolkit TaskLauncher; one launcher is submitted per firing instant.
type taskLauncher struct {
	scheduler *Scheduler
	millis    int64
}

func (l *taskLauncher) run() {
	defer l.scheduler.launcherMgr.completed(l)
	l.scheduler.taskTable.executeIfMatch(l.scheduler, l.millis)
}

// taskLauncherManager is aligned with the utility toolkit TaskLauncherManager.
type taskLauncherManager struct {
	scheduler *Scheduler
	mu        sync.Mutex
	cond      *sync.Cond
	launchers []*taskLauncher
	idleCh    chan struct{}
}

func newTaskLauncherManager(s *Scheduler) *taskLauncherManager {
	idleCh := make(chan struct{})
	close(idleCh)
	m := &taskLauncherManager{scheduler: s, idleCh: idleCh}
	m.cond = sync.NewCond(&m.mu)
	return m
}

func (m *taskLauncherManager) spawn(millis int64) *taskLauncher {
	l := &taskLauncher{scheduler: m.scheduler, millis: millis}
	m.mu.Lock()
	if len(m.launchers) == 0 {
		m.idleCh = make(chan struct{})
	}
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
			m.cond.Broadcast()
			if len(m.launchers) == 0 {
				close(m.idleCh)
			}
			return
		}
	}
}

func (m *taskLauncherManager) runningCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.launchers)
}

func (m *taskLauncherManager) wait() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for len(m.launchers) > 0 {
		m.cond.Wait()
	}
}

func (m *taskLauncherManager) waitContext(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	m.mu.Lock()
	idleCh := m.idleCh
	m.mu.Unlock()
	select {
	case <-idleCh:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
