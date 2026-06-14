package resty

import "testing"

func TestGlobalHeadersArePlainValues(t *testing.T) {
	SetGlobalHeader("X-Resty-Plain", "one")
	AddGlobalHeader("X-Resty-Plain", "two")
	defer RemoveGlobalHeader("X-Resty-Plain")

	headers := CloneGlobalHeaders()
	if got := headers["X-Resty-Plain"]; len(got) != 2 || got[0] != "one" || got[1] != "two" {
		t.Fatalf("CloneGlobalHeaders()[X-Resty-Plain] = %v, want [one two]", got)
	}
}
