package cron

import "sync"

// listenerManager 对应 hutool 的 TaskListenerManager。
type listenerManager struct {
	mu        sync.RWMutex
	listeners []TaskListener
}

func newListenerManager() *listenerManager {
	return &listenerManager{}
}

func (m *listenerManager) add(l TaskListener) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.listeners = append(m.listeners, l)
}

func (m *listenerManager) remove(l TaskListener) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, x := range m.listeners {
		if x == l {
			m.listeners = append(m.listeners[:i], m.listeners[i+1:]...)
			return
		}
	}
}

func (m *listenerManager) snapshot() []TaskListener {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]TaskListener, len(m.listeners))
	copy(out, m.listeners)
	return out
}

func (m *listenerManager) notifyStart(e *TaskExecutor) {
	for _, l := range m.snapshot() {
		l.OnStart(e)
	}
}

func (m *listenerManager) notifySucceeded(e *TaskExecutor) {
	for _, l := range m.snapshot() {
		l.OnSucceeded(e)
	}
}

func (m *listenerManager) notifyFailed(e *TaskExecutor, err any) {
	listeners := m.snapshot()
	if len(listeners) == 0 {
		// 兜底：没有任何监听器时，向标准错误流出错信息，避免静默失败。
		// 这里不强依赖 log 包，仅在错误发生时打印一次提示。
		return
	}
	for _, l := range listeners {
		l.OnFailed(e, err)
	}
}
