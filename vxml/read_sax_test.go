package vxml

import (
	stdxml "encoding/xml"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestFacadeSAXReadOptions(t *testing.T) {
	var starts []string
	if err := ReadBySAX(strings.NewReader(`<root><a>1</a></root>`), func(tok stdxml.Token) error {
		if start, ok := tok.(stdxml.StartElement); ok {
			starts = append(starts, start.Name.Local)
		}
		return nil
	}); err != nil || !reflect.DeepEqual(starts, []string{"root", "a"}) {
		t.Fatalf("ReadBySAX facade starts=%v err=%v", starts, err)
	}
	var nsStarts []stdxml.Name
	if err := ReadBySAXWithOptions(strings.NewReader(`<root xmlns:p="urn:p"><p:a>1</p:a></root>`), func(tok stdxml.Token) error {
		if start, ok := tok.(stdxml.StartElement); ok {
			nsStarts = append(nsStarts, start.Name)
		}
		return nil
	}, WithNamespaceAware(false)); err != nil {
		t.Fatalf("ReadBySAXWithOptions facade: %v", err)
	}
	if !reflect.DeepEqual(nsStarts, []stdxml.Name{{Local: "root"}, {Local: "a"}}) {
		t.Fatalf("ReadBySAXWithOptions facade names=%#v", nsStarts)
	}
	tmp := t.TempDir() + "/sax.xml"
	if err := os.WriteFile(tmp, []byte(`<root><a>1</a></root>`), 0o600); err != nil {
		t.Fatal(err)
	}
	starts = nil
	if err := ReadBySAXFile(tmp, func(tok stdxml.Token) error {
		if start, ok := tok.(stdxml.StartElement); ok {
			starts = append(starts, start.Name.Local)
		}
		return nil
	}); err != nil || !reflect.DeepEqual(starts, []string{"root", "a"}) {
		t.Fatalf("ReadBySAXFile facade starts=%v err=%v", starts, err)
	}
	starts = nil
	if err := ReadBySAXFileWithOptions(tmp, func(tok stdxml.Token) error {
		if start, ok := tok.(stdxml.StartElement); ok {
			starts = append(starts, start.Name.Local)
		}
		return nil
	}, WithStrict(true)); err != nil || !reflect.DeepEqual(starts, []string{"root", "a"}) {
		t.Fatalf("ReadBySAXFileWithOptions facade starts=%v err=%v", starts, err)
	}
}
