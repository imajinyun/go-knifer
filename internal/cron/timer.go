package cron

import (
	"sync/atomic"
	"time"
)

const (
	timerUnitSecond = int64(time.Second / time.Millisecond)
	timerUnitMinute = int64(time.Minute / time.Millisecond)
)

// cronTimer 对应 hutool 的 CronTimer，按时间步长触发 launcher。
type cronTimer struct {
	scheduler *Scheduler
	stop      atomic.Bool
	stopCh    chan struct{}
}

func newCronTimer(s *Scheduler) *cronTimer {
	return &cronTimer{scheduler: s, stopCh: make(chan struct{})}
}

// run 启动循环，按 unit 触发任务，并自动校准时钟漂移。
func (t *cronTimer) run() {
	unit := timerUnitMinute
	if t.scheduler.config.MatchSecond {
		unit = timerUnitSecond
	}
	thisTime := nowMillis()
	for !t.stop.Load() {
		thisTime += unit
		now := nowMillis()
		sleep := thisTime - now
		if sleep > 0 {
			if !t.sleep(time.Duration(sleep) * time.Millisecond) {
				return
			}
		}
		// 时间倒退或大幅跳跃：重置基准
		now = nowMillis()
		if now-thisTime > 2*unit || thisTime-now > 2*unit {
			thisTime = now
		}
		if t.stop.Load() {
			return
		}
		t.scheduler.launcherMgr.spawn(thisTime)
	}
}

// sleep 在被停止时立即返回 false。
func (t *cronTimer) sleep(d time.Duration) bool {
	if d <= 0 {
		return true
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-t.stopCh:
		return false
	case <-timer.C:
		return true
	}
}

func (t *cronTimer) stopTimer() {
	if t.stop.CompareAndSwap(false, true) {
		close(t.stopCh)
	}
}

// nowMillis 返回当前 Unix 毫秒。
func nowMillis() int64 {
	return time.Now().UnixMilli()
}

// timeFromMillisInLocation 将毫秒转为指定时区的 time.Time。
func timeFromMillisInLocation(ms int64, loc *time.Location) time.Time {
	if loc == nil {
		loc = time.Local
	}
	return time.UnixMilli(ms).In(loc)
}
