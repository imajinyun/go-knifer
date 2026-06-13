package db

import "testing"

func TestUtilityTypes(t *testing.T) {
	page := NewPage(0, 0, Asc("name"))
	if page.Number != 1 || page.Size != 20 || page.Offset() != 0 || page.Limit() != 20 {
		t.Fatalf("page = %#v", page)
	}
	result := NewPageResult(page, 41, []Entity{NewEntity("users")})
	if result.TotalPage != 3 || !result.IsFirst() || result.IsLast() {
		t.Fatalf("result = %#v", result)
	}
	if got := WrapperForDialect(DialectMySQL).Wrap("users.name"); got != "`users`.`name`" {
		t.Fatalf("wrapped = %q", got)
	}
	if !IsSafeIdentifier("users.name") || !IsSafeIdentifier("users.*") || !IsSafeIdentifier("`users`.`name`") {
		t.Fatal("expected safe identifiers to be accepted")
	}
	if IsSafeIdentifier("users; drop table users") || IsSafeIdentifier("COUNT(*)") || IsSafeIdentifier("users name") {
		t.Fatal("expected unsafe identifiers to be rejected")
	}
}
