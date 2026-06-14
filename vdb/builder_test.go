package vdb

import "testing"

func TestFacadeBuilder(t *testing.T) {
	sqlText, args, err := NewBuilder(WithDialect(DialectPostgres), WithWrapper(WrapperForDialect(DialectPostgres))).
		Select("id").
		From("users").
		Where(Eq("name", "alice")).
		SQL()
	if err != nil {
		t.Fatalf("SQL() error = %v", err)
	}
	if sqlText != `SELECT "id" FROM "users" WHERE "name" = $1` {
		t.Fatalf("sql = %q", sqlText)
	}
	if len(args) != 1 || args[0] != "alice" {
		t.Fatalf("args = %#v", args)
	}
}

func TestFacadeBuilderOptionsWrapperPrecedence(t *testing.T) {
	sqlText, _, err := NewBuilder(WithDialect(DialectMySQL)).Select("id").From("users").SQL()
	if err != nil {
		t.Fatalf("SQL() with dialect option error = %v", err)
	}
	if sqlText != "SELECT `id` FROM `users`" {
		t.Fatalf("SQL() with dialect default wrapper = %q", sqlText)
	}

	sqlText, _, err = NewBuilder(WithDialect(DialectMySQL), WithWrapper(NewWrapper("\"", "\""))).Select("id").From("users").SQL()
	if err != nil {
		t.Fatalf("SQL() with wrapper option error = %v", err)
	}
	if sqlText != `SELECT "id" FROM "users"` {
		t.Fatalf("SQL() with explicit wrapper = %q", sqlText)
	}
}
