package vdb

import "testing"

func TestFacadeNamedSQL(t *testing.T) {
	named, err := ParseNamed("select * from users where id=:id", map[string]any{"id": 1}, DialectQuestion)
	if err != nil {
		t.Fatalf("ParseNamed() error = %v", err)
	}
	if named.SQL != "select * from users where id=?" || named.Params[0] != 1 {
		t.Fatalf("named = %#v", named)
	}
}
