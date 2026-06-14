package vcron_test

import (
	"testing"

	"github.com/imajinyun/go-knifer/vcron"
)

func TestFacadeDefaultSchedulerOptions(t *testing.T) {
	global := vcron.ConfigureDefaultScheduler(vcron.WithIDGenerator(func() string { return "global-id" }))
	t.Cleanup(func() { vcron.ConfigureDefaultScheduler() })
	isolated := vcron.NewSchedulerWithOptions(vcron.WithIDGenerator(func() string { return "facade-isolated" }))

	id, err := vcron.CronScheduleFuncWithOptions("* * * * *", func() {}, vcron.WithDefaultScheduler(isolated))
	if err != nil {
		t.Fatalf("CronScheduleFuncWithOptions: %v", err)
	}
	if id != "facade-isolated" || isolated.Size() != 1 || global.Size() != 0 {
		t.Fatalf("default scheduler option not isolated: id=%q isolated=%d global=%d", id, isolated.Size(), global.Size())
	}
	if !vcron.CronRemoveWithOptions(id, vcron.WithDefaultScheduler(isolated)) {
		t.Fatal("CronRemoveWithOptions should remove isolated task")
	}
}
