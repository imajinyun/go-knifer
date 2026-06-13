package conf

import "testing"

func TestParseSetting(t *testing.T) {
	s, err := Parse("name = gokit\n[server]\nport=8080\ndebug=true")
	if err != nil {
		t.Fatal(err)
	}
	if got := s.Get("name"); got != "gokit" {
		t.Fatalf("Get(name) = %q", got)
	}
	if got := s.GetByGroup("server", "port"); got != "8080" {
		t.Fatalf("GetByGroup(server, port) = %q", got)
	}
	if got := s.GetOrDefault("missing", "def"); got != "def" {
		t.Fatalf("GetOrDefault() = %q", got)
	}
}
