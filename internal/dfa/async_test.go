package dfa

import "testing"

func TestAsyncRunnerCanBeConfiguredAndReset(t *testing.T) {
	ResetAsyncRunner()
	t.Cleanup(ResetAsyncRunner)

	runs := 0
	ConfigureAsyncRunner(func(fn func()) {
		runs++
		fn()
	})
	InitAsync([]string{"async"})
	if runs != 1 || !Contains("async word") {
		t.Fatalf("InitAsync with configured runner runs=%d contains=%v", runs, Contains("async word"))
	}
	InitStringAsync("runner", DefaultSeparator)
	if runs != 2 || !Contains("runner word") {
		t.Fatalf("InitStringAsync with configured runner runs=%d contains=%v", runs, Contains("runner word"))
	}
	InitStringAsyncWithOptions("a-b", DefaultSeparator, WithCharFilter(func(r rune) bool { return r != '-' }))
	if runs != 3 || !Contains("ab") {
		t.Fatalf("InitStringAsyncWithOptions with configured runner runs=%d contains=%v", runs, Contains("ab"))
	}

	ResetAsyncRunner()
	Init([]string{"reset"})
	if !Contains("reset") {
		t.Fatal("ResetAsyncRunner should preserve synchronous Init behavior")
	}
}
