package vdb

import (
	"reflect"
	"testing"
)

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
