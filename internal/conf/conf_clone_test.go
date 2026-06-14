package conf

import "testing"

func TestCloneReturnsIndependentCopy(t *testing.T) {
	s := New()
	s.Set("root", "value")
	s.SetByGroup("server", "port", "8080")

	clone := s.Clone()
	clone.Set("root", "changed")
	clone.SetByGroup("server", "port", "9090")
	clone.SetByGroup("server", "host", "localhost")

	if got := s.Get("root"); got != "value" {
		t.Fatalf("source root changed to %q", got)
	}
	if got := s.GetByGroup("server", "port"); got != "8080" {
		t.Fatalf("source server.port changed to %q", got)
	}
	if got := s.GetByGroup("server", "host"); got != "" {
		t.Fatalf("source server.host changed to %q", got)
	}
	if got := clone.GetByGroup("server", "host"); got != "localhost" {
		t.Fatalf("clone server.host = %q", got)
	}
}
