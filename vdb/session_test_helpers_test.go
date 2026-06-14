package vdb

import (
	"database/sql"
	"database/sql/driver"
	"io"
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
