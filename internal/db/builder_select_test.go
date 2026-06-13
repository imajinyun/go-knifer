package db

import (
	"reflect"
	"testing"

	knifer "github.com/imajinyun/go-knifer"
)

func TestSQLBuilderSelectWherePage(t *testing.T) {
	sqlText, args, err := NewBuilder(WithDialect(DialectPostgres), WithWrapper(WrapperForDialect(DialectPostgres))).
		Select("id", "name").
		From("users").
		Where(Eq("status", "active"), OrWith(In("role", "admin", "owner"))).
		OrderBy(Desc("id")).
		Page(NewPage(2, 10)).
		SQL()
	if err != nil {
		t.Fatalf("SQL() error = %v", err)
	}
	wantSQL := `SELECT "id", "name" FROM "users" WHERE "status" = $1 OR "role" IN ($2, $3) ORDER BY "id" DESC LIMIT 10 OFFSET 10`
	if sqlText != wantSQL {
		t.Fatalf("sql = %q, want %q", sqlText, wantSQL)
	}
	if !reflect.DeepEqual(args, []any{"active", "admin", "owner"}) {
		t.Fatalf("args = %#v", args)
	}
}

func TestSQLBuilderRejectsUnsafeIdentifiers(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "select field", err: sqlErr(Select("id; drop table users").From("users").SQL())},
		{name: "from table", err: sqlErr(Select("id").From("users; drop table users").SQL())},
		{name: "where field", err: sqlErr(Select("id").From("users").Where(Eq("id OR 1=1", 1)).SQL())},
		{name: "order field", err: sqlErr(Select("id").From("users").OrderBy(Asc("id desc; drop table users")).SQL())},
		{name: "insert table", err: sqlErr(Insert(NewEntity("users; drop table users").Set("name", "alice")).SQL())},
		{name: "insert field", err: sqlErr(Insert(NewEntity("users").Set("name; drop", "alice")).SQL())},
		{name: "update field", err: sqlErr(Update(NewEntity("users").Set("name = hacked", "alice")).Where(Eq("id", 1)).SQL())},
		{name: "delete table", err: sqlErr(Delete("users; drop table users").Where(Eq("id", 1)).SQL())},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertDBCode(t, tt.err, knifer.ErrCodeInvalidInput)
		})
	}
}
