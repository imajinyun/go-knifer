package cron

import (
	"sync"
	"sync/atomic"
	"time"
)

// Scheduler 对应 hutool 的 Scheduler，是 gkcron 的核心调度器。
type Scheduler struct {
	mu          sync.Mutex
	config      *Config
	started     atomic.Bool
	timer       *cronTimer
	taskTable   *TaskTable
	launcherMgr *taskLauncherManager
	executorMgr *taskExecutorManager
	listenerMgr *listenerManager

	// executor 控制任务执行使用的 goroutine。允许外部替换为带有限并发的执行器。
	executor func(func())
}

// NewScheduler 创建 Scheduler。
func NewScheduler() *Scheduler {
	s := &Scheduler{
		config:    NewConfig(),
		taskTable: NewTaskTable(),
	}
	s.launcherMgr = newTaskLauncherManager(s)
	s.executorMgr = newTaskExecutorManager(s)
	s.listenerMgr = newListenerManager()
	s.executor = func(fn func()) { go fn() }
	return s
}

// Config 返回配置。
func (s *Scheduler) Config() *Config { return s.config }

// SetMatchSecond 设置是否匹配到秒（启动后修改无效）。
func (s *Scheduler) SetMatchSecond(b bool) *Scheduler {
	s.config.MatchSecond = b
	return s
}

// IsMatchSecond 是否匹配秒。
func (s *Scheduler) IsMatchSecond() bool { return s.config.MatchSecond }

// SetTimeZone 设置时区。
func (s *Scheduler) SetTimeZone(loc *time.Location) *Scheduler {
	if loc == nil {
		loc = time.Local
	}
	s.config.Location = loc
	return s
}

// SetExecutor 设置自定义执行函数。
func (s *Scheduler) SetExecutor(exec func(func())) *Scheduler {
	if exec != nil {
		s.executor = exec
	}
	return s
}

// AddListener 添加监听器。
func (s *Scheduler) AddListener(l TaskListener) *Scheduler {
	s.listenerMgr.add(l)
	return s
}

// RemoveListener 移除监听器。
func (s *Scheduler) RemoveListener(l TaskListener) *Scheduler {
	s.listenerMgr.remove(l)
	return s
}

// Schedule 通过表达式注册任务，自动生成 id。返回任务 id。
func (s *Scheduler) Schedule(pattern string, task Task) (string, error) {
	id := generateID()
	if err := s.ScheduleWithID(id, pattern, task); err != nil {
		return "", err
	}
	return id, nil
}

// ScheduleFunc 注册函数式任务。
func (s *Scheduler) ScheduleFunc(pattern string, fn func()) (string, error) {
	return s.Schedule(pattern, TaskFunc(fn))
}

// ScheduleWithID 使用指定 id 注册任务。
func (s *Scheduler) ScheduleWithID(id, pattern string, task Task) error {
	p, err := NewPattern(pattern)
	if err != nil {
		return err
	}
	return s.SchedulePattern(id, p, task)
}

// SchedulePattern 使用已解析的 Pattern 注册任务。
func (s *Scheduler) SchedulePattern(id string, p *Pattern, task Task) error {
	return s.taskTable.Add(id, p, task)
}

// Deschedule 删除任务。
func (s *Scheduler) Deschedule(id string) bool {
	return s.taskTable.Remove(id)
}

// UpdatePattern 更新任务表达式。
func (s *Scheduler) UpdatePattern(id, pattern string) error {
	p, err := NewPattern(pattern)
	if err != nil {
		return err
	}
	if !s.taskTable.UpdatePattern(id, p) {
		return NewCronError("task %q not found", id)
	}
	return nil
}

// TaskTable 返回任务表。
func (s *Scheduler) TaskTable() *TaskTable { return s.taskTable }

// GetPattern 按 id 返回 Pattern。
func (s *Scheduler) GetPattern(id string) *Pattern { return s.taskTable.GetPattern(id) }

// GetTask 按 id 返回 Task。
func (s *Scheduler) GetTask(id string) Task { return s.taskTable.GetTask(id) }

// IsEmpty 任务表是否为空。
func (s *Scheduler) IsEmpty() bool { return s.taskTable.IsEmpty() }

// Size 任务数量。
func (s *Scheduler) Size() int { return s.taskTable.Size() }

// Clear 清空任务。
func (s *Scheduler) Clear() {
	for _, id := range s.taskTable.IDs() {
		s.taskTable.Remove(id)
	}
}

// IsStarted 是否已启动。
func (s *Scheduler) IsStarted() bool { return s.started.Load() }

// Start 启动调度器。
func (s *Scheduler) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.started.CompareAndSwap(false, true) {
		return NewCronError("scheduler already started")
	}
	s.timer = newCronTimer(s)
	go s.timer.run()
	return nil
}

// Stop 停止调度器。clearTasks 为 true 时同时清空任务表。
func (s *Scheduler) Stop(clearTasks ...bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.started.CompareAndSwap(true, false) {
		return
	}
	if s.timer != nil {
		s.timer.stopTimer()
		s.timer = nil
	}
	if len(clearTasks) > 0 && clearTasks[0] {
		s.Clear()
	}
}

// submit 通过当前 executor 异步执行 fn。
func (s *Scheduler) submit(fn func()) {
	s.executor(fn)
}
