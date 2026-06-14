package vcron_test

import (
	"testing"

	"github.com/imajinyun/go-knifer/vcron"
)

func TestFacadeSchedulerLifecycle(t *testing.T) {
	s := vcron.NewScheduler()
	if s == nil {
		t.Fatal("expected non-nil scheduler")
	}

	id, err := vcron.CronScheduleFunc("* * * * * *", func() {})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == "" {
		t.Fatal("expected non-empty task id")
	}

	if !vcron.CronRemove(id) {
		t.Fatal("expected task to be removed")
	}
}
