package db

import (
	"reflect"
	"testing"

	knifer "github.com/imajinyun/go-knifer"
)

func TestBuildCountSQL(t *testing.T) {
	sqlText, args, err := buildCountSQL(DialectPostgres, WrapperForDialect(DialectPostgres), []string{"users"}, Eq("status", "active"))
	if err != nil {
		t.Fatalf("buildCountSQL() error = %v", err)
	}
	if sqlText != `SELECT COUNT(*) FROM "users" WHERE "status" = $1` {
		t.Fatalf("sql = %q", sqlText)
	}
	if !reflect.DeepEqual(args, []any{"active"}) {
		t.Fatalf("args = %#v", args)
	}

	_, _, err = buildCountSQL(DialectQuestion, Wrapper{}, []string{"users; drop table users"})
	assertDBCode(t, err, knifer.ErrCodeInvalidInput)
}

func TestListColumnsSQLRejectsUnsafeTable(t *testing.T) {
	_, _, _, err := listColumnsSQL(DialectSQLite, "users; drop table users")
	assertDBCode(t, err, knifer.ErrCodeInvalidInput)
}

func TestScanAndMetaHelpersReportInvalidInputAndUnsupported(t *testing.T) {
	_, err := ScanRows(nil)
	assertDBCode(t, err, knifer.ErrCodeInvalidInput)

	err = AssignEntity(NewEntity("users"), nil)
	assertDBCode(t, err, knifer.ErrCodeInvalidInput)

	_, err = listTablesSQL(DialectOracle)
	assertDBCode(t, err, knifer.ErrCodeUnsupported)

	_, _, _, err = listColumnsSQL(DialectOracle, "users")
	assertDBCode(t, err, knifer.ErrCodeUnsupported)
}
