package vdb

import (
	"database/sql"
	"testing"
	"time"
)

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
