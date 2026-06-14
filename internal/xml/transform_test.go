package xml

import (
	"strings"
	"testing"
)

func TestTransformWith(t *testing.T) {
	var out strings.Builder
	if err := TransformWith(strings.NewReader(`<root><a>1</a></root>`), &out, WithOmitDeclaration(true)); err != nil || out.String() != `<root><a>1</a></root>` {
		t.Fatalf("TransformWith = %q, %v", out.String(), err)
	}
}

func TestTransformWithOptions(t *testing.T) {
	var out strings.Builder
	if err := TransformWithOptions(strings.NewReader(`<root xmlns:p="urn:p"><p:a>1</p:a></root>`), &out,
		WithTransformParseOptions(WithNamespaceAware(false)),
		WithTransformWriteOptions(WithOmitDeclaration(true)),
	); err != nil || strings.Contains(out.String(), `xmlns`) || !strings.Contains(out.String(), `<a>1</a>`) {
		t.Fatalf("TransformWithOptions = %q, %v", out.String(), err)
	}
}
