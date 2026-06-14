package vjson_test

import (
	"strings"
	"testing"

	"github.com/imajinyun/go-knifer/vjson"
)

func TestFacadeJSONOptions(t *testing.T) {
	compact, err := vjson.ToStr(map[string]any{"name": "go", "empty": nil}, vjson.WithIgnoreNullValue(true))
	if err != nil {
		t.Fatalf("ToStr with options: %v", err)
	}
	if strings.Contains(compact, "empty") {
		t.Fatalf("ToStr WithIgnoreNullValue should omit null field: %s", compact)
	}
	formatted := vjson.FormatWithOptions(`{"a":1,"b":{"c":2}}`, vjson.WithFormatIndentWidth(2), vjson.WithFormatSpaceAfterKey(false))
	if !strings.Contains(formatted, "\n  \"a\":1") || strings.Contains(formatted, "\": 1") {
		t.Fatalf("FormatWithOptions = %q", formatted)
	}
	cfg := vjson.NewConfig()
	cfg.IgnoreNullValue = true
	obj, err := vjson.ParseObjWithOptions(map[string]any{"name": "go", "empty": nil}, vjson.WithParseConfig(cfg))
	if err != nil {
		t.Fatalf("ParseObjWithOptions: %v", err)
	}
	if got := obj.ToString(); strings.Contains(got, "empty") {
		t.Fatalf("ParseObjWithOptions should apply config: %s", got)
	}
}
