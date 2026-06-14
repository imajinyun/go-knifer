package conf

import "testing"

func TestParseByExtYAML(t *testing.T) {
	yamlConf, err := ParseByExt("app.yaml", []byte("app:\n  name: demo"))
	if err != nil {
		t.Fatal(err)
	}
	if got := yamlConf.GetByGroup("app", "name"); got != "demo" {
		t.Fatalf("ParseByExt yaml app.name = %q", got)
	}
}

func TestParseByExtWithCustomParser(t *testing.T) {
	custom, err := ParseByExtWithOptions("app.custom", []byte("ignored"), WithParserForExt("custom", func([]byte) (*Conf, error) {
		c := New()
		c.Set("name", "custom-parser")
		return c, nil
	}))
	if err != nil {
		t.Fatal(err)
	}
	if got := custom.Get("name"); got != "custom-parser" {
		t.Fatalf("custom parser name = %q", got)
	}
}
