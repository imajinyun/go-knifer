package vdb_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"

	"github.com/imajinyun/knifer-go/vdb"
)

func init() {
	sql.Register("vdb_example", exampleDriver{})
}

func ExampleSelect() {
	b := vdb.Select("id", "name").From("users").Where(vdb.Gt("age", 18))
	sql, args, _ := b.SQL()
	fmt.Println(sql)
	fmt.Println(args)
	// Output:
	// SELECT id, name FROM users WHERE age > ?
	// [18]
}

func ExampleNewEntity() {
	e := vdb.NewEntity("users")
	e.Values["name"] = "Alice"
	e.Values["age"] = 30
	b := vdb.Insert(e)
	sql, args, _ := b.SQL()
	fmt.Println(sql)
	fmt.Println(args)
	// Output:
	// INSERT INTO users (age, name) VALUES (?, ?)
	// [30 Alice]
}

func ExampleEntityFromMap() {
	entity := vdb.EntityFromMap("users", map[string]any{"id": 7, "name": "alice"})
	fmt.Println(entity.Table)
	fmt.Println(entity.Values["name"])
	// Output:
	// users
	// alice
}

func ExampleAssignEntity() {
	entity := vdb.EntityFromMap("users", map[string]any{"id": int64(7), "full_name": "Alice"})
	var dst struct {
		ID       int64
		FullName string `db:"full_name"`
	}
	err := vdb.AssignEntity(entity, &dst)
	fmt.Println(err == nil)
	fmt.Println(dst.ID, dst.FullName)
	// Output:
	// true
	// 7 Alice
}

func ExampleBuildConditions() {
	sql, args, err := vdb.BuildConditions(
		vdb.Like("name", vdb.BuildLikeValue("go", "prefix")),
		vdb.In("role", "admin", "owner"),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(sql)
	fmt.Println(args)
	// Output:
	// name LIKE ? AND role IN (?, ?)
	// [go% admin owner]
}

func ExampleDelete() {
	b := vdb.Delete("users").Where(vdb.IsNull("deleted_at"))
	sql, args, _ := b.SQL()
	fmt.Println(sql)
	fmt.Println(args)
	// Output:
	// DELETE FROM users WHERE deleted_at IS NULL
	// []
}

func ExampleBuildLikeValue() {
	fmt.Println(vdb.BuildLikeValue("go", "contains"))
	// Output: %go%
}

func ExampleNewPage() {
	sql, args, _ := vdb.NewBuilder(vdb.WithDialect(vdb.DialectMySQL)).
		Select("id", "created_at").
		From("orders").
		Page(vdb.NewPage(2, 10, vdb.Desc("created_at"))).
		SQL()
	fmt.Println(sql)
	fmt.Println(args)
	// Output:
	// SELECT `id`, `created_at` FROM `orders` ORDER BY `created_at` DESC LIMIT 10 OFFSET 10
	// []
}

func ExampleIsSafeIdentifier() {
	fmt.Println(vdb.IsSafeIdentifier("orders.created_at"))
	fmt.Println(vdb.IsSafeIdentifier("orders; drop table orders"))
	// Output:
	// true
	// false
}

func ExampleParseNamed() {
	named, err := vdb.ParseNamed(
		"SELECT * FROM users WHERE id = :id::int AND status = :status",
		map[string]any{"id": 7, "status": "active"},
		vdb.DialectPostgres,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(named.SQL)
	fmt.Println(named.Params)
	fmt.Println(named.Names)
	// Output:
	// SELECT * FROM users WHERE id = $1::int AND status = $2
	// [7 active]
	// [id status]
}

func ExampleUpdate() {
	entity := vdb.EntityFromMap("users", map[string]any{"active": false})
	sql, args, _ := vdb.Update(entity).Where(vdb.Eq("id", 7)).SQL()
	fmt.Println(sql)
	fmt.Println(args)
	// Output:
	// UPDATE users SET active = ? WHERE id = ?
	// [false 7]
}

func ExampleOrGroup() {
	sql, args, _ := vdb.BuildConditions(
		vdb.AndGroup(vdb.Eq("tenant_id", 42)),
		vdb.OrGroup(vdb.Eq("status", "active"), vdb.Eq("status", "pending")),
	)
	fmt.Println(sql)
	fmt.Println(args)
	// Output:
	// (tenant_id = ?) OR (status = ? AND status = ?)
	// [42 active pending]
}

func ExampleConditionsFromEntity() {
	conditions := vdb.ConditionsFromEntity(vdb.EntityFromMap("users", map[string]any{"id": 7, "name": "alice"}))
	sql, args, _ := vdb.BuildConditions(conditions...)
	fmt.Println(sql)
	fmt.Println(args)
	// Output:
	// id = ? AND name = ?
	// [7 alice]
}

func ExampleExec() {
	sqlDB, err := sql.Open("vdb_example", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() { _ = sqlDB.Close() }()

	db := vdb.Use(sqlDB)
	result, err := vdb.Exec(context.Background(), db, "UPDATE users SET active = ? WHERE id = ?", true, 7)
	if err != nil {
		fmt.Println(err)
		return
	}
	affected, _ := result.RowsAffected()
	fmt.Println(affected)
	// Output:
	// 1
}

func ExampleRaw() {
	sql, args, _ := vdb.Raw("COUNT(*) FILTER (WHERE status = ?)", "active").SQL()
	fmt.Println(sql)
	fmt.Println(args)
	// Output:
	// COUNT(*) FILTER (WHERE status = ?)
	// [active]
}

func ExampleNewWrapper() {
	wrapper := vdb.NewWrapper("[", "]")
	fmt.Println(wrapper.Wrap("users.name"))
	// Output:
	// [users].[name]
}

func ExampleNormalizeDialect() {
	fmt.Println(vdb.NormalizeDialect("pgx"))
	fmt.Println(vdb.NormalizeDialect("sqlite3"))
	// Output:
	// postgres
	// sqlite
}

func ExampleWrapperForDialect() {
	fmt.Println(vdb.WrapperForDialect(vdb.DialectMySQL).Wrap("users.name"))
	fmt.Println(vdb.WrapperForDialect(vdb.DialectPostgres).Wrap("users.name"))
	// Output:
	// `users`.`name`
	// "users"."name"
}

func ExampleNewPageResult() {
	page := vdb.NewPage(2, 10)
	result := vdb.NewPageResult(page, 25, []string{"order-11", "order-12"})
	fmt.Println(result.Page, result.PageSize, result.Total, result.TotalPage)
	fmt.Println(result.IsFirst(), result.IsLast())
	fmt.Println(result.Items)
	// Output:
	// 2 10 25 3
	// false false
	// [order-11 order-12]
}

func ExampleDB_Upsert() {
	sqlDB, err := sql.Open("vdb_example", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() { _ = sqlDB.Close() }()

	db := vdb.Use(sqlDB, vdb.WithDialect(vdb.DialectPostgres), vdb.WithWrapper(vdb.WrapperForDialect(vdb.DialectPostgres)))
	entity := vdb.EntityFromMap("users", map[string]any{"id": 7, "name": "alice"})
	result, err := db.Upsert(context.Background(), entity, []string{"id"})
	if err != nil {
		fmt.Println(err)
		return
	}
	affected, _ := result.RowsAffected()
	fmt.Println(affected)
	// Output:
	// 1
}

func ExampleDB_ExecBatch() {
	sqlDB, err := sql.Open("vdb_example", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() { _ = sqlDB.Close() }()

	db := vdb.Use(sqlDB)
	results, err := db.ExecBatch(context.Background(), "INSERT INTO users(name) VALUES (?)", []any{"alice"}, []any{"bob"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(results))
	// Output:
	// 2
}

func ExampleDB_Tx() {
	sqlDB, err := sql.Open("vdb_example", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() { _ = sqlDB.Close() }()

	db := vdb.Use(sqlDB)
	err = db.Tx(context.Background(), nil, func(s *vdb.Session) error {
		_, err := s.Exec(context.Background(), "UPDATE users SET active = ? WHERE id = ?", true, 7)
		return err
	})
	fmt.Println(err == nil)
	// Output:
	// true
}

func ExampleScanRows() {
	rows := exampleRows()
	entities, err := vdb.ScanRows(rows)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(entities[0].Values["id"])
	fmt.Println(entities[0].Values["name"])
	// Output:
	// 7
	// alice
}

func ExampleScanOne() {
	rows := exampleRows()
	entity, ok, err := vdb.ScanOne(rows)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ok)
	fmt.Println(entity.Values["name"])
	// Output:
	// true
	// alice
}

func exampleRows() *sql.Rows {
	sqlDB, err := sql.Open("vdb_example", "")
	if err != nil {
		panic(err)
	}
	rows, err := sqlDB.Query("SELECT example")
	if err != nil {
		_ = sqlDB.Close()
		panic(err)
	}
	return rows
}

type exampleDriver struct{}

func (exampleDriver) Open(string) (driver.Conn, error) {
	return exampleConn{}, nil
}

type exampleConn struct{}

func (exampleConn) Prepare(string) (driver.Stmt, error) {
	return nil, fmt.Errorf("prepare not supported")
}
func (exampleConn) Close() error              { return nil }
func (exampleConn) Begin() (driver.Tx, error) { return exampleTx{}, nil }

func (exampleConn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	return exampleTx{}, nil
}

func (exampleConn) ExecContext(context.Context, string, []driver.NamedValue) (driver.Result, error) {
	return exampleResult{}, nil
}

func (exampleConn) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	return &exampleRowsResult{
		columns: []string{"id", "name"},
		values:  [][]driver.Value{{int64(7), []byte("alice")}},
	}, nil
}

type exampleTx struct{}

func (exampleTx) Commit() error   { return nil }
func (exampleTx) Rollback() error { return nil }

type exampleResult struct{}

func (exampleResult) LastInsertId() (int64, error) { return 0, nil }
func (exampleResult) RowsAffected() (int64, error) { return 1, nil }

type exampleRowsResult struct {
	columns []string
	values  [][]driver.Value
	pos     int
}

func (r *exampleRowsResult) Columns() []string {
	return r.columns
}

func (r *exampleRowsResult) Close() error {
	return nil
}

func (r *exampleRowsResult) Next(dest []driver.Value) error {
	if r.pos >= len(r.values) {
		return io.EOF
	}
	copy(dest, r.values[r.pos])
	r.pos++
	return nil
}
