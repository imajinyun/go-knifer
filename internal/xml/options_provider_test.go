package xml

import (
	stdxml "encoding/xml"
	"errors"
	"io"
	"io/fs"
	"os"
	"reflect"
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

func TestXMLFileProviderOptions(t *testing.T) {
	openedRead := ""
	doc, err := ReadXMLFile("virtual.xml", WithOpenFile(func(path string) (io.ReadCloser, error) {
		openedRead = path
		return io.NopCloser(strings.NewReader(`<root><from>provider</from></root>`)), nil
	}))
	if err != nil || ElementText(doc.Root, "from") != "provider" || openedRead != "virtual.xml" {
		t.Fatalf("ReadXMLFile provider doc=%#v path=%q err=%v", doc, openedRead, err)
	}

	var saxStarts []string
	openedRead = ""
	if err := ReadBySAXFileWithOptions("sax.xml", func(tok stdxml.Token) error {
		if start, ok := tok.(stdxml.StartElement); ok {
			saxStarts = append(saxStarts, start.Name.Local)
		}
		return nil
	}, WithOpenFile(func(path string) (io.ReadCloser, error) {
		openedRead = path
		return io.NopCloser(strings.NewReader(`<root><a/></root>`)), nil
	})); err != nil || !reflect.DeepEqual(saxStarts, []string{"root", "a"}) || openedRead != "sax.xml" {
		t.Fatalf("ReadBySAXFileWithOptions provider starts=%v path=%q err=%v", saxStarts, openedRead, err)
	}

	var mkdirPath string
	var mkdirPerm fs.FileMode
	var openPath string
	var openFlag int
	var openPerm fs.FileMode
	var written strings.Builder
	closer := nopWriteCloser{Writer: &written}
	err = WriteFile("/virtual/out.xml", CreateXMLWithRoot("root"), WithOmitDeclaration(true),
		WithMkdirAll(func(path string, perm fs.FileMode) error {
			mkdirPath, mkdirPerm = path, perm
			return nil
		}),
		WithOpenWriteFile(func(path string, flag int, perm fs.FileMode) (io.WriteCloser, error) {
			openPath, openFlag, openPerm = path, flag, perm
			return closer, nil
		}),
		WithDirPerm(0o700), WithFilePerm(0o600),
	)
	if err != nil {
		t.Fatalf("WriteFile provider: %v", err)
	}
	if mkdirPath != "/virtual" || mkdirPerm != 0o700 || openPath != "/virtual/out.xml" || openPerm != 0o600 || openFlag&os.O_CREATE == 0 || written.String() != `<root/>` {
		t.Fatalf("WriteFile providers mkdir=%q/%v open=%q flag=%#x perm=%v content=%q", mkdirPath, mkdirPerm, openPath, openFlag, openPerm, written.String())
	}
}
