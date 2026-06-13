package xml

import (
	"strings"
	"testing"
)

func TestReadCreateWriteAndFormat(t *testing.T) {
	doc, err := ParseXML(`<root><name id="1">alice</name><tags>a</tags><tags>b</tags></root>`)
	if err != nil {
		t.Fatalf("ParseXML failed: %v", err)
	}
	root := GetRootElement(doc)
	if root == nil || root.Name.Local != "root" {
		t.Fatalf("root mismatch: %#v", root)
	}
	if got := ElementText(root, "name"); got != "alice" {
		t.Fatalf("ElementText = %q", got)
	}
	if got := GetElements(root, ""); len(got) != 3 {
		t.Fatalf("GetElements all children length = %d", len(got))
	}
	if got := GetElements(root, "tags"); len(got) != 2 {
		t.Fatalf("GetElements tags length = %d", len(got))
	}
	if GetOwnerDocument(GetElement(root, "name")).Root != root {
		t.Fatal("GetOwnerDocument should walk to root")
	}
	plain, err := MarshalString(doc, WithOmitDeclaration(true))
	if err != nil || plain != `<root><name id="1">alice</name><tags>a</tags><tags>b</tags></root>` {
		t.Fatalf("MarshalString = %q, %v", plain, err)
	}
	formatted, err := Format(plain)
	if err != nil || !strings.Contains(formatted, `<?xml version="1.0" encoding="UTF-8"?>`) || !strings.Contains(formatted, "\n  <name") {
		t.Fatalf("Format = %q, %v", formatted, err)
	}
	created := CreateXMLWithRootNS("user", "urn:test")
	AppendChild(created.Root, "name")
	AppendText(GetElement(created.Root, "name"), "bob")
	createdStr, err := MarshalString(created, WithOmitDeclaration(true))
	if err != nil || !strings.Contains(createdStr, `xmlns="urn:test"`) || !strings.Contains(createdStr, `<name>bob</name>`) {
		t.Fatalf("CreateXMLWithRootNS serialized = %q, %v", createdStr, err)
	}
}

func TestFormatWithOptions(t *testing.T) {
	formatted, err := FormatWithOptions(`<root xmlns:p="urn:p"><p:a>1</p:a></root>`,
		WithFormatParseOptions(WithNamespaceAware(false)),
		WithFormatWriteOptions(WithOmitDeclaration(true), WithIndent(4)),
	)
	if err != nil || strings.Contains(formatted, `xmlns`) || !strings.Contains(formatted, "\n    <a>") {
		t.Fatalf("FormatWithOptions = %q, %v", formatted, err)
	}
}
