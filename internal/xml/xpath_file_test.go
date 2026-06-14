package xml

import (
	stdxml "encoding/xml"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestReadXMLFileAndXPathHelpers(t *testing.T) {
	tmp := t.TempDir() + "/x.xml"
	if err := os.WriteFile(tmp, []byte(`<root><a>1</a><a>2</a></root>`), 0o600); err != nil {
		t.Fatal(err)
	}
	doc, err := ReadXML(tmp)
	if err != nil {
		t.Fatalf("ReadXML file failed: %v", err)
	}
	if got := GetByXPath("/root/a", doc, "string"); got != "1" {
		t.Fatalf("GetByXPath string = %v", got)
	}
	if got := GetByXPath("/root/a", doc, "nodes"); len(got.([]*Element)) != 2 {
		t.Fatalf("GetByXPath nodes = %#v", got)
	}
	if got := GetElementByXPath("/root/a", doc); got == nil || strings.TrimSpace(got.Text) != "1" {
		t.Fatalf("GetElementByXPath = %#v", got)
	}
	if got := GetNodeByXPath("/root/missing", doc); got != nil {
		t.Fatalf("missing XPath should be nil: %#v", got)
	}
	if got := GetNodeListByXPath("//a", doc); len(got) != 2 {
		t.Fatalf("GetNodeListByXPath = %d", len(got))
	}
}

func TestReadBySAXFileVariants(t *testing.T) {
	tmp := t.TempDir() + "/x.xml"
	if err := os.WriteFile(tmp, []byte(`<root><a>1</a><a>2</a></root>`), 0o600); err != nil {
		t.Fatal(err)
	}
	var saxFileStarts []string
	if err := ReadBySAXFile(tmp, func(tok stdxml.Token) error {
		if start, ok := tok.(stdxml.StartElement); ok {
			saxFileStarts = append(saxFileStarts, start.Name.Local)
		}
		return nil
	}); err != nil || !reflect.DeepEqual(saxFileStarts, []string{"root", "a", "a"}) {
		t.Fatalf("ReadBySAXFile starts=%v err=%v", saxFileStarts, err)
	}
	saxFileStarts = nil
	if err := ReadBySAXFileWithOptions(tmp, func(tok stdxml.Token) error {
		if start, ok := tok.(stdxml.StartElement); ok {
			saxFileStarts = append(saxFileStarts, start.Name.Local)
		}
		return nil
	}, WithStrict(true)); err != nil || !reflect.DeepEqual(saxFileStarts, []string{"root", "a", "a"}) {
		t.Fatalf("ReadBySAXFileWithOptions starts=%v err=%v", saxFileStarts, err)
	}
}
