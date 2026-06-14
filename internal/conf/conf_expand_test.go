package conf

import "testing"

func TestExpandVariables(t *testing.T) {
	t.Setenv("CONF_ENV_HOST", "envhost")
	s, err := Parse(`
host=localhost
base=http://${host}:8080
env=${ENV:CONF_ENV_HOST}
missing=${missing:fallback}
[db]
host=db.local
url=postgres://${db.host}/${name:app}
`)
	if err != nil {
		t.Fatal(err)
	}

	if got := s.GetExpanded("base"); got != "http://localhost:8080" {
		t.Fatalf("GetExpanded(base) = %q", got)
	}
	if got := s.GetExpanded("env"); got != "envhost" {
		t.Fatalf("GetExpanded(env) = %q", got)
	}
	if got := s.GetExpanded("missing"); got != "fallback" {
		t.Fatalf("GetExpanded(missing) = %q", got)
	}
	if got := s.GetByGroupExpanded("db", "url"); got != "postgres://db.local/app" {
		t.Fatalf("GetByGroupExpanded(db,url) = %q", got)
	}
	if got := s.Expand().Get("base"); got != "http://localhost:8080" {
		t.Fatalf("Expand().Get(base) = %q", got)
	}
}

func TestExpandWithOptionsUsesCustomEnvLookup(t *testing.T) {
	s, err := Parse(`
host=${ENV:CONF_HOST}
base=http://${host}:${port:8080}
[db]
url=postgres://${ENV:CONF_DB_HOST}/${name:app}
`)
	if err != nil {
		t.Fatal(err)
	}
	lookup := func(key string) string {
		switch key {
		case "CONF_HOST":
			return "option.local"
		case "CONF_DB_HOST":
			return "db.option.local"
		default:
			return ""
		}
	}
	if got := s.GetExpandedWithOptions("host", WithEnvLookup(lookup)); got != "option.local" {
		t.Fatalf("GetExpandedWithOptions(host) = %q", got)
	}
	if got := s.GetExpandedWithOptions("base", WithEnvLookup(lookup)); got != "http://option.local:8080" {
		t.Fatalf("GetExpandedWithOptions(base) = %q", got)
	}
	if got := s.GetByGroupExpandedWithOptions("db", "url", WithEnvLookup(lookup)); got != "postgres://db.option.local/app" {
		t.Fatalf("GetByGroupExpandedWithOptions(db.url) = %q", got)
	}
	expanded := s.ExpandWithOptions(WithEnvLookup(lookup))
	if got := expanded.Get("host"); got != "option.local" {
		t.Fatalf("ExpandWithOptions host = %q", got)
	}
}
