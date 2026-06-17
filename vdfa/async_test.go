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

func TestFacadeStringAsyncWithOptions(t *testing.T) {
	ResetAsyncRunner()
	t.Cleanup(ResetAsyncRunner)

	runs := 0
	ConfigureAsyncRunner(func(fn func()) {
		runs++
		fn()
	})
	InitStringAsyncWithOptions("a-b|c-d", '|', WithCharFilter(func(r rune) bool { return r != '-' }))
	if runs != 1 || !Contains("ab") || !Contains("cd") {
		t.Fatalf("InitStringAsyncWithOptions runs=%d contains ab=%v cd=%v", runs, Contains("ab"), Contains("cd"))
	}
	InitStringAsync("plain", 0)
	if runs != 2 || !Contains("plain") {
		t.Fatalf("InitStringAsync runs=%d contains=%v", runs, Contains("plain"))
	}
}
