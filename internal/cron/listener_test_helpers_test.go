package cron

import "sync/atomic"

type testListener struct {
	started *atomic.Int32
	succ    *atomic.Int32
	failed  *atomic.Int32
}

func (l *testListener) OnStart(*TaskExecutor) {
	if l.started != nil {
		l.started.Add(1)
	}
}

func (l *testListener) OnSucceeded(*TaskExecutor) {
	if l.succ != nil {
		l.succ.Add(1)
	}
}

func (l *testListener) OnFailed(*TaskExecutor, any) {
	if l.failed != nil {
		l.failed.Add(1)
	}
}

type panicListener struct {
	onStart     bool
	onSucceeded bool
	onFailed    bool
}

func (l *panicListener) OnStart(*TaskExecutor) {
	if l.onStart {
		panic("start listener")
	}
}

func (l *panicListener) OnSucceeded(*TaskExecutor) {
	if l.onSucceeded {
		panic("success listener")
	}
}

func (l *panicListener) OnFailed(*TaskExecutor, any) {
	if l.onFailed {
		panic("failed listener")
	}
}
