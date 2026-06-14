package vdb

import (
	"database/sql"
	"testing"
)

func TestFacadeOpenProvider(t *testing.T) {
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
}
