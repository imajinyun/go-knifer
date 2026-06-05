package vdb

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"testing"
	"time"

	knifer "github.com/imajinyun/go-knifer"
)

func init() {
	sql.Register("vdb_pool_test", poolTestDriver{})
}

type poolTestDriver struct{}

func (poolTestDriver) Open(string) (driver.Conn, error) { return poolTestConn{}, nil }

type poolTestConn struct{}

func (poolTestConn) Prepare(string) (driver.Stmt, error) { return poolTestStmt{}, nil }
func (poolTestConn) Close() error                        { return nil }
func (poolTestConn) Begin() (driver.Tx, error)           { return poolTestTx{}, nil }

type poolTestStmt struct{}

func (poolTestStmt) Close() error                               { return nil }
func (poolTestStmt) NumInput() int                              { return -1 }
func (poolTestStmt) Exec([]driver.Value) (driver.Result, error) { return driver.RowsAffected(0), nil }
func (poolTestStmt) Query([]driver.Value) (driver.Rows, error)  { return poolTestRows{}, nil }

type poolTestRows struct{}

func (poolTestRows) Columns() []string         { return []string{"id"} }
func (poolTestRows) Close() error              { return nil }
func (poolTestRows) Next([]driver.Value) error { return io.EOF }

type poolTestTx struct{}

func (poolTestTx) Commit() error   { return nil }
func (poolTestTx) Rollback() error { return nil }

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

func TestFacadeNamedSQL(t *testing.T) {
	named, err := ParseNamed("select * from users where id=:id", map[string]any{"id": 1}, DialectQuestion)
	if err != nil {
		t.Fatalf("ParseNamed() error = %v", err)
	}
	if named.SQL != "select * from users where id=?" || named.Params[0] != 1 {
		t.Fatalf("named = %#v", named)
	}
}

func TestFacadeDBErrorContract(t *testing.T) {
	_, _, err := NewBuilder().SQL()
	if err == nil {
		t.Fatal("SQL() error = nil, want invalid input error")
	}
	if !errors.Is(err, knifer.ErrCodeInvalidInput) {
		t.Fatalf("errors.Is(err, ErrCodeInvalidInput) = false: %v", err)
	}
	code, ok := knifer.CodeOf(err)
	if !ok || code != knifer.ErrCodeInvalidInput {
		t.Fatalf("CodeOf(err) = %q, %v; want invalid input", code, ok)
	}
	var dbErr *DBError
	if !errors.As(err, &dbErr) {
		t.Fatalf("errors.As(err, *DBError) = false: %v", err)
	}
}

func TestFacadePoolOptionsApplyToWrappedDB(t *testing.T) {
	sqlDB, err := sql.Open("vdb_pool_test", "")
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	db := Use(sqlDB,
		WithMaxOpenConns(7),
		WithMaxIdleConns(3),
		WithConnMaxLifetime(2*time.Minute),
		WithConnMaxIdleTime(time.Minute),
	)
	if db.SQLDB() != sqlDB {
		t.Fatal("Use should preserve the wrapped *sql.DB")
	}
	stats := sqlDB.Stats()
	if stats.MaxOpenConnections != 7 {
		t.Fatalf("MaxOpenConnections = %d", stats.MaxOpenConnections)
	}
	if got := sqlDB.Stats().MaxOpenConnections; got != db.SQLDB().Stats().MaxOpenConnections {
		t.Fatalf("wrapped DB stats mismatch = %d", got)
	}
}

func TestFacadePoolOptionsApplyWhenOpeningDB(t *testing.T) {
	db, err := Open("vdb_pool_test", "", WithMaxOpenConns(5), WithMaxIdleConns(2))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() { _ = db.Close() }()

	stats := db.SQLDB().Stats()
	if stats.MaxOpenConnections != 5 {
		t.Fatalf("MaxOpenConnections = %d", stats.MaxOpenConnections)
	}
	if db.Dialect() != DialectQuestion {
		t.Fatalf("Dialect = %q", db.Dialect())
	}
}
