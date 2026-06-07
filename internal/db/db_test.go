package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"io"
	"sync/atomic"
	"testing"
)

var txPanicRollbackCount atomic.Int32

func init() {
	sql.Register("goknifer_tx_panic_test", txPanicTestDriver{})
}

type txPanicTestDriver struct{}

func (txPanicTestDriver) Open(string) (driver.Conn, error) { return txPanicTestConn{}, nil }

type txPanicTestConn struct{}

func (txPanicTestConn) Prepare(string) (driver.Stmt, error) { return txPanicTestStmt{}, nil }
func (txPanicTestConn) Close() error                        { return nil }
func (txPanicTestConn) Begin() (driver.Tx, error)           { return txPanicTestTx{}, nil }

type txPanicTestStmt struct{}

func (txPanicTestStmt) Close() error  { return nil }
func (txPanicTestStmt) NumInput() int { return -1 }
func (txPanicTestStmt) Exec([]driver.Value) (driver.Result, error) {
	return driver.RowsAffected(0), nil
}
func (txPanicTestStmt) Query([]driver.Value) (driver.Rows, error) { return txPanicTestRows{}, nil }

type txPanicTestRows struct{}

func (txPanicTestRows) Columns() []string         { return []string{"id"} }
func (txPanicTestRows) Close() error              { return nil }
func (txPanicTestRows) Next([]driver.Value) error { return io.EOF }

type txPanicTestTx struct{}

func (txPanicTestTx) Commit() error { return nil }
func (txPanicTestTx) Rollback() error {
	txPanicRollbackCount.Add(1)
	return nil
}

func TestTxRollsBackOnPanic(t *testing.T) {
	txPanicRollbackCount.Store(0)
	sqlDB, err := sql.Open("goknifer_tx_panic_test", "")
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	db := Use(sqlDB)
	defer func() {
		p := recover()
		if p != "boom" {
			t.Fatalf("panic = %#v, want boom", p)
		}
		if got := txPanicRollbackCount.Load(); got != 1 {
			t.Fatalf("rollback count = %d, want 1", got)
		}
	}()

	_ = db.Tx(context.Background(), nil, func(*Session) error {
		panic("boom")
	})
}
