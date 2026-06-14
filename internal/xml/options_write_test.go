package xml

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

func TestPerCallOptionsAndWriteErrors(t *testing.T) {
	doc, err := ParseXML(`<root xmlns:p="urn:p"><p:a>1</p:a></root>`, WithNamespaceAware(false))
	if err != nil {
		t.Fatalf("ParseXML with option failed: %v", err)
	}
	child := GetElement(doc.Root, "a")
	if child == nil || child.Name.Space != "" {
		t.Fatalf("per-call namespace option not applied: %#v", child)
	}

	doc, err = ParseXML(`<root xmlns:p="urn:p"><p:a>1</p:a></root>`)
	if err != nil {
		t.Fatalf("ParseXML failed: %v", err)
	}
	child = GetElement(doc.Root, "a")
	if child == nil || child.Name.Space != "urn:p" {
		t.Fatalf("default namespace awareness should remain enabled: %#v", child)
	}

	pretty, err := MarshalString(CreateXMLWithRoot("root"), WithCharset("GBK"), WithPretty())
	if err != nil || !strings.HasPrefix(pretty, `<?xml version="1.0" encoding="GBK"?>`) || !strings.Contains(pretty, "\n<root/>") {
		t.Fatalf("MarshalString options = %q, %v", pretty, err)
	}
	if err := WriteTo(nil, doc); err == nil || !strings.Contains(err.Error(), "nil writer") {
		t.Fatalf("WriteTo nil writer err = %v", err)
	}
	if err := WriteTo(io.Discard, "unsupported"); err == nil || !strings.Contains(err.Error(), "unsupported node") {
		t.Fatalf("WriteTo unsupported err = %v", err)
	}
	writeErr := errors.New("write failed")
	if err := WriteTo(failingWriter{err: writeErr}, doc); !errors.Is(err, writeErr) {
		t.Fatalf("WriteTo writer err = %v", err)
	}

	tmp := t.TempDir() + "/out.xml"
	if err := WriteFile(tmp, doc, WithOmitDeclaration(true)); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	if err := WriteFile(tmp, doc, WithOmitDeclaration(true), WithOverwrite(false)); err == nil {
		t.Fatalf("WriteFile should reject overwrite when disabled")
	}
	written, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(written), `<a>1</a>`) {
		t.Fatalf("WriteFile content = %q", written)
	}
	if _, err := ReadXMLReader(strings.NewReader(`<root><a>1</a></root>`), WithMaxBytes(8)); err == nil {
		t.Fatalf("ReadXMLReader should fail when max bytes truncates input")
	}
}
