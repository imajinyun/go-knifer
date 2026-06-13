package db

import (
	"reflect"
	"testing"

	knifer "github.com/imajinyun/go-knifer"
)

func TestSQLBuilderRejectsUnsafeConditionOperator(t *testing.T) {
	_, _, err := NewBuilder(WithDialect(DialectQuestion)).
		Select("id").
		From("users").
		Where(Condition{Field: "name", Op: "= ? OR 1=1 --", Value: "alice"}).
		SQL()
	assertDBCode(t, err, knifer.ErrCodeInvalidInput)

	sqlText, args, err := NewBuilder(WithDialect(DialectQuestion)).
		Select("id").
		From("users").
		Where(Condition{Field: "name", Op: "NOT LIKE", Value: "%bot%"}).
		SQL()
	if err != nil {
		t.Fatalf("SQL() NOT LIKE error = %v", err)
	}
	if sqlText != "SELECT id FROM users WHERE name NOT LIKE ?" || !reflect.DeepEqual(args, []any{"%bot%"}) {
		t.Fatalf("NOT LIKE sql=%q args=%#v", sqlText, args)
	}
}

func TestBuildLikeValue(t *testing.T) {
	if got := BuildLikeValue("go", "contains"); got != "%go%" {
		t.Fatalf("contains like = %q", got)
	}
}
