package vstr

import (
	"reflect"
	"testing"
)

func TestStringFacade(t *testing.T) {
	if !IsEmpty("") || !IsNotEmpty("x") || !IsBlank(" \t") || !IsNotBlank("x") {
		t.Fatal("blank/empty checks failed")
	}
	if !HasEmpty("a", "") || !HasBlank("a", " ") || !IsAllEmpty("", "") || !IsAllBlank(" ", "") {
		t.Fatal("aggregate blank/empty checks failed")
	}
	if Trim(" x ") != "x" || TrimToEmpty(" x ") != "x" || TrimStart(" x") != "x" || TrimEnd("x ") != "x" {
		t.Fatal("trim helpers failed")
	}
	if Sub("你好世界", 1, 3) != "好世" || SubBefore("a.b.c", ".", false) != "a" || SubAfter("a.b.c", ".", true) != "c" {
		t.Fatal("substring helpers failed")
	}
	if got := Split("a,b", ","); !reflect.DeepEqual(got, []string{"a", "b"}) {
		t.Fatalf("Split failed: %v", got)
	}
	if got := SplitTrim(" a, ,b ", ","); !reflect.DeepEqual(got, []string{"a", "b"}) {
		t.Fatalf("SplitTrim failed: %v", got)
	}
	if Repeat("a", 3) != "aaa" || PadLeft("go", 4, '0') != "00go" || PadRight("go", 4, '0') != "go00" {
		t.Fatal("repeat/pad helpers failed")
	}
	if !Contains("hello", "ell") || !ContainsAny("hello", "x", "ell") || !ContainsAll("hello", "he", "lo") || !ContainsIgnoreCase("Hello", "he") {
		t.Fatal("contains helpers failed")
	}
	if !StartsWith("hello", "he") || !EndsWith("hello", "lo") || !EqualsIgnoreCase("Go", "go") {
		t.Fatal("prefix/suffix/equal helpers failed")
	}
	if Reverse("你好") != "好你" || Format("{},{}", "a", 1) != "a,1" {
		t.Fatal("reverse/format helpers failed")
	}
	if RemovePrefix("foobar", "foo") != "bar" || RemoveSuffix("foobar", "bar") != "foo" {
		t.Fatal("remove prefix/suffix failed")
	}
	if AddPrefixIfNot("bar", "foo") != "foobar" || AddSuffixIfNot("foo", "bar") != "foobar" {
		t.Fatal("add prefix/suffix failed")
	}
	if Length("你好") != 2 || RuneLen("你好") != 2 {
		t.Fatal("Length failed")
	}
	if !ContainsEmoji("hello😀") || RemoveEmoji("hello😀") != "hello" {
		t.Fatal("emoji helpers failed")
	}
	value := 5
	if DefaultIfNil(&value, 1) != 5 || DefaultIfEmpty("", "x") != "x" || DefaultIfBlank(" ", "x") != "x" {
		t.Fatal("default helpers failed")
	}
	if EscapeHTML(`<a>&`) != "&lt;a&gt;&amp;" || UnescapeHTML("&lt;a&gt;&amp;") != `<a>&` {
		t.Fatal("html helpers failed")
	}
	if ToCamelCase("hello_world") != "helloWorld" || ToPascalCase("hello_world") != "HelloWorld" || ToUnderlineCase("HelloWorld") != "hello_world" || ToKebabCase("HelloWorld") != "hello-world" {
		t.Fatal("naming helpers failed")
	}
}
