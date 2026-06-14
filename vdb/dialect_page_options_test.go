package vdb

import "testing"

func TestFacadeDialectPageAndOptions(t *testing.T) {
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
}
