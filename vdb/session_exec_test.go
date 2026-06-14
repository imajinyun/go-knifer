package vdb

import (
	"context"
	"database/sql"
	"testing"
)

func TestFacadeExec(t *testing.T) {
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
