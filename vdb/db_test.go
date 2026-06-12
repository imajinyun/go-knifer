package vdb

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"reflect"
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

func TestFacadeTopLevelBuildersAndConditions(t *testing.T) {
	rawSQL, rawArgs, err := Raw("SELECT ? AS id", 7).SQL()
	if err != nil || rawSQL != "SELECT ? AS id" || !reflect.DeepEqual(rawArgs, []any{7}) {
		t.Fatalf("Raw SQL=%q args=%#v err=%v", rawSQL, rawArgs, err)
	}

	insertSQL, insertArgs, err := Insert(EntityFromMap("users", map[string]any{"name": "alice", "age": 18})).SQL()
	if err != nil {
		t.Fatalf("Insert SQL: %v", err)
	}
	if insertSQL != "INSERT INTO users (age, name) VALUES (?, ?)" || !reflect.DeepEqual(insertArgs, []any{18, "alice"}) {
		t.Fatalf("Insert SQL=%q args=%#v", insertSQL, insertArgs)
	}

	updateSQL, updateArgs, err := Update(NewEntity("users").Set("name", "bob")).
		Where(AndGroup(Gt("id", 10), Lte("id", 20)), OrWith(IsNull("deleted_at"))).
		SQL()
	if err != nil {
		t.Fatalf("Update SQL: %v", err)
	}
	if updateSQL != "UPDATE users SET name = ? WHERE (id > ? AND id <= ?) OR deleted_at IS NULL" {
		t.Fatalf("Update SQL = %q", updateSQL)
	}
	if !reflect.DeepEqual(updateArgs, []any{"bob", 10, 20}) {
		t.Fatalf("Update args = %#v", updateArgs)
	}

	deleteSQL, deleteArgs, err := Delete("users").
		Where(OrGroup(Ne("status", "active"), Between("created_at", 1, 9), IsNotNull("blocked_at"))).
		SQL()
	if err != nil {
		t.Fatalf("Delete SQL: %v", err)
	}
	if deleteSQL != "DELETE FROM users WHERE (status <> ? AND created_at BETWEEN ? AND ? AND blocked_at IS NOT NULL)" {
		t.Fatalf("Delete SQL = %q", deleteSQL)
	}
	if !reflect.DeepEqual(deleteArgs, []any{"active", 1, 9}) {
		t.Fatalf("Delete args = %#v", deleteArgs)
	}

	conditionSQL, conditionArgs, err := BuildConditions(Like("name", BuildLikeValue("go", "prefix")), In("role", "admin", "owner"))
	if err != nil {
		t.Fatalf("BuildConditions: %v", err)
	}
	if conditionSQL != "name LIKE ? AND role IN (?, ?)" || !reflect.DeepEqual(conditionArgs, []any{"go%", "admin", "owner"}) {
		t.Fatalf("BuildConditions SQL=%q args=%#v", conditionSQL, conditionArgs)
	}

	conds := ConditionsFromEntity(NewEntity("users").Set("id", 1).Set("name", "alice"))
	if len(conds) != 2 || conds[0].Field != "id" || conds[1].Field != "name" {
		t.Fatalf("ConditionsFromEntity = %#v", conds)
	}
}

func TestFacadeDialectPageOptionsAndExec(t *testing.T) {
	if NormalizeDialect("postgresql") != DialectPostgres {
		t.Fatalf("NormalizeDialect(postgresql) = %q", NormalizeDialect("postgresql"))
	}
	if got := NewWrapper("[", "]").Wrap("users.name"); got != "[users].[name]" {
		t.Fatalf("NewWrapper.Wrap = %q", got)
	}
	if !IsInClause("select * from t where id in (?, ?)") || IsInClause("select * from t") {
		t.Fatal("IsInClause result mismatch")
	}
	if got := RemoveOuterOrderBy("select * from users order by id desc"); got != "select * from users" {
		t.Fatalf("RemoveOuterOrderBy = %q", got)
	}

	page := NewPage(2, 5, Desc("id"), Asc("name"))
	if page.Number != 2 || page.Size != 5 || page.Offset() != 5 || len(page.Orders) != 2 {
		t.Fatalf("NewPage = %#v", page)
	}
	result := NewPageResult(page, 12, []string{"a", "b"})
	if result.TotalPage != 3 || result.IsFirst() || result.IsLast() || len(result.Items) != 2 {
		t.Fatalf("NewPageResult = %#v", result)
	}

	opts := NewOptions()
	if opts.Dialect != DialectQuestion {
		t.Fatalf("NewOptions = %#v", opts)
	}
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
