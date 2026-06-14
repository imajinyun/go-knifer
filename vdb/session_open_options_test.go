package vdb

import "testing"

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
