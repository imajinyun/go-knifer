package str

import "testing"

// Tests cover the utility toolkit-core NamingCaseTest.

func TestNamingCase(t *testing.T) {
	cases := []struct {
		camel, pascal, under, kebab string
	}{
		{"helloWorld", "HelloWorld", "hello_world", "hello-world"},
		{"helloWorldFoo", "HelloWorldFoo", "hello_world_foo", "hello-world-foo"},
		{"a", "A", "a", "a"},
	}
	for _, c := range cases {
		if got := ToCamelCase(c.under); got != c.camel {
			t.Fatalf("ToCamelCase(%q)=%q want %q", c.under, got, c.camel)
		}
		if got := ToPascalCase(c.under); got != c.pascal {
			t.Fatalf("ToPascalCase(%q)=%q want %q", c.under, got, c.pascal)
		}
		if got := ToUnderlineCase(c.camel); got != c.under {
			t.Fatalf("ToUnderlineCase(%q)=%q want %q", c.camel, got, c.under)
		}
		if got := ToKebabCase(c.camel); got != c.kebab {
			t.Fatalf("ToKebabCase(%q)=%q want %q", c.camel, got, c.kebab)
		}
	}
}

func TestNamingFromKebab(t *testing.T) {
	if got := ToCamelCase("hello-world"); got != "helloWorld" {
		t.Fatalf("kebab->camel: %q", got)
	}
}
