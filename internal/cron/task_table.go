package cron

import "sync"

// TaskTable 对应 hutool 的 TaskTable，按插入顺序保存任务。
type TaskTable struct {
	mu       sync.RWMutex
	ids      []string
	patterns []*Pattern
	tasks    []Task
}

// NewTaskTable 创建任务表。
func NewTaskTable() *TaskTable { return &TaskTable{} }

// Add 添加任务，id 重复返回错误。
func (t *TaskTable) Add(id string, p *Pattern, task Task) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, x := range t.ids {
		if x == id {
			return NewCronError("duplicate task id %q", id)
		}
	}
	t.ids = append(t.ids, id)
	t.patterns = append(t.patterns, p)
	t.tasks = append(t.tasks, task)
	return nil
}

// Remove 删除任务，返回是否删除成功。
func (t *TaskTable) Remove(id string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	for i, x := range t.ids {
		if x == id {
			t.ids = append(t.ids[:i], t.ids[i+1:]...)
			t.patterns = append(t.patterns[:i], t.patterns[i+1:]...)
			t.tasks = append(t.tasks[:i], t.tasks[i+1:]...)
			return true
		}
	}
	return false
}

// UpdatePattern 更新指定任务的 Cron 表达式。
func (t *TaskTable) UpdatePattern(id string, p *Pattern) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	for i, x := range t.ids {
		if x == id {
			t.patterns[i] = p
			return true
		}
	}
	return false
}

// GetTask 按 id 返回任务。
func (t *TaskTable) GetTask(id string) Task {
	t.mu.RLock()
	defer t.mu.RUnlock()
	for i, x := range t.ids {
		if x == id {
			return t.tasks[i]
		}
	}
	return nil
}

// GetPattern 按 id 返回 Pattern。
func (t *TaskTable) GetPattern(id string) *Pattern {
	t.mu.RLock()
	defer t.mu.RUnlock()
	for i, x := range t.ids {
		if x == id {
			return t.patterns[i]
		}
	}
	return nil
}

// Size 返回任务数量。
func (t *TaskTable) Size() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.ids)
}

// IsEmpty 判断是否为空。
func (t *TaskTable) IsEmpty() bool { return t.Size() == 0 }

// IDs 返回所有任务 id 的拷贝。
func (t *TaskTable) IDs() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]string, len(t.ids))
	copy(out, t.ids)
	return out
}

// executeIfMatch 按当前时间触发匹配任务。
func (t *TaskTable) executeIfMatch(s *Scheduler, fireTime int64) {
	t.mu.RLock()
	ids := append([]string(nil), t.ids...)
	pats := append([]*Pattern(nil), t.patterns...)
	tasks := append([]Task(nil), t.tasks...)
	t.mu.RUnlock()
	tt := timeFromMillisInLocation(fireTime, s.config.Location)
	for i, p := range pats {
		if p.Match(tt, s.config.MatchSecond) {
			ct := NewCronTask(ids[i], p, tasks[i])
			s.executorMgr.spawn(ct)
		}
	}
}
