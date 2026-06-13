package xml

import (
	stdxml "encoding/xml"
	"fmt"
	"strings"
	"testing"
)

func TestInternalHelpers(t *testing.T) {
	if got := typeName(&sampleBean{}); got != "samplebean" {
		t.Fatalf("typeName pointer = %q", got)
	}
	if got := typeName(map[string]any{}); got != "" {
		t.Fatalf("typeName map = %q", got)
	}
	if got, ok := normalizeMap(map[string]string{"a": "b"}); !ok || !gotOK(got, "a", "b") {
		t.Fatalf("normalizeMap string map = %#v", got)
	}
	if got, ok := normalizeMap(map[any]any{"a": 1}); !ok || !gotOK(got, "a", 1) {
		t.Fatalf("normalizeMap any map = %#v", got)
	}
	if got, ok := normalizeMap(map[int]string{1: "one"}); !ok || !gotOK(got, "1", "one") {
		t.Fatalf("normalizeMap typed map = %#v", got)
	}
	if m, ok := normalizeMap(nil); ok || m != nil {
		t.Fatalf("normalizeMap nil = %#v %v", m, ok)
	}
	if !isNilValue((*string)(nil)) || isNilValue("") {
		t.Fatal("isNilValue mismatch")
	}
	if !isStruct(sampleBean{}) || !isStruct(&sampleBean{}) || isStruct(nil) || isStruct(map[string]any{}) {
		t.Fatal("isStruct mismatch")
	}
	cfg := applyParse(nil)
	if got := parseScalar("false", cfg); got != false {
		t.Fatalf("parseScalar false = %#v", got)
	}
	if got := parseScalar("plain", cfg); got != "plain" {
		t.Fatalf("parseScalar string = %#v", got)
	}
	if got := fmt.Sprint(structToMap(struct {
		XMLName stdxml.Name `xml:"root"`
	}{}, true, false)); strings.Contains(got, "XMLName") {
		t.Fatalf("structToMap should skip XMLName: %s", got)
	}
}
