package json

import "testing"

func TestParseRejectsTrailingContent(t *testing.T) {
	if _, err := Parse(`{"a":1} {"b":2}`); err == nil {
		t.Fatal("expected trailing content error")
	}
}

func TestParseObjAndArrayErrors(t *testing.T) {
	if _, err := ParseObj(`[1,2]`); err == nil {
		t.Fatalf("expect error parsing array as obj")
	}
	if _, err := ParseArray(`{}`); err == nil {
		t.Fatalf("expect error parsing object as array")
	}
}
