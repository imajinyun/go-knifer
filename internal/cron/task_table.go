package cron

import "sync"

// TaskTable is aligned with hutool TaskTable and stores tasks in insertion order.
type TaskTable struct {
	mu       sync.RWMutex
	ids      []string
	patterns []*Pattern
	tasks    []Task
}

// NewTaskTable creates a task table.
func NewTaskTable() *TaskTable { return &TaskTable{} }

// Add adds a task and returns an error for duplicate ids.
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

// Remove deletes a task and reports whether it was removed.
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

// UpdatePattern updates the cron expression of the specified task.
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

// GetTask returns a task by id.
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

// GetPattern returns a Pattern by id.
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

// Size returns the task count.
func (t *TaskTable) Size() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.ids)
}

// IsEmpty reports whether the table is empty.
func (t *TaskTable) IsEmpty() bool { return t.Size() == 0 }

// IDs returns a copy of all task ids.
func (t *TaskTable) IDs() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]string, len(t.ids))
	copy(out, t.ids)
	return out
}

// executeIfMatch triggers matched tasks at the current fire time.
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
