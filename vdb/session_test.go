package vdb

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"io"
	"testing"
	"time"
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

func TestFacadeOpenProviderAndExec(t *testing.T) {
	opened := false
	customDB, err := Open("ignored", "ignored", WithSQLOpenFunc(func(driverName, dsn string) (*sql.DB, error) {
		opened = true
		return sql.Open("vdb_pool_test", "")
	}))
	if err != nil {
		t.Fatalf("Open with custom SQLOpen: %v", err)
	}
	defer func() { _ = customDB.Close() }()
	if !opened {
		t.Fatal("WithSQLOpenFunc provider was not called")
	}

	sqlDB, err := sql.Open("vdb_pool_test", "")
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()
	db := Use(sqlDB)
	if _, err := Exec(context.Background(), db, "UPDATE users SET name=?", "alice"); err != nil {
		t.Fatalf("Exec: %v", err)
	}
}
