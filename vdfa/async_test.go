package vdfa

import "testing"

func TestFacadeAsyncRunner(t *testing.T) {
	ResetAsyncRunner()
	t.Cleanup(ResetAsyncRunner)

	runs := 0
	ConfigureAsyncRunner(func(fn func()) {
		runs++
		fn()
	})
	InitAsync([]string{"facade-async"})
	if runs != 1 || !Contains("facade-async word") {
		t.Fatalf("InitAsync runner runs=%d contains=%v", runs, Contains("facade-async word"))
	}
}
