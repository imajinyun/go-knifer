package vxml

import (
	stdxml "encoding/xml"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestFacadeXMLReadFileProviderOptions(t *testing.T) {
	openedPath := ""
	doc, err := ReadXMLFile("virtual.xml", WithOpenFile(func(path string) (io.ReadCloser, error) {
		openedPath = path
		return io.NopCloser(strings.NewReader(`<root><from>facade</from></root>`)), nil
	}))
	if err != nil || ElementText(doc.Root, "from") != "facade" || openedPath != "virtual.xml" {
		t.Fatalf("ReadXMLFile provider doc=%#v path=%q err=%v", doc, openedPath, err)
	}

	var starts []string
	openedPath = ""
	if err := ReadBySAXFileWithOptions("sax.xml", func(tok stdxml.Token) error {
		if start, ok := tok.(stdxml.StartElement); ok {
			starts = append(starts, start.Name.Local)
		}
		return nil
	}, WithOpenFile(func(path string) (io.ReadCloser, error) {
		openedPath = path
		return io.NopCloser(strings.NewReader(`<root><a/></root>`)), nil
	})); err != nil || !reflect.DeepEqual(starts, []string{"root", "a"}) || openedPath != "sax.xml" {
		t.Fatalf("ReadBySAXFileWithOptions provider starts=%v path=%q err=%v", starts, openedPath, err)
	}
}
