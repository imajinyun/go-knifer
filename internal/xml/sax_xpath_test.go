package xml

import (
	stdxml "encoding/xml"
	"errors"
	"io"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestReadVariantsSAXXPathAndFile(t *testing.T) {
	if CleanInvalid("a\x00b\x08c") != "abc" {
		t.Fatal("CleanInvalid failed")
	}
	if CleanComment("<a><!-- hidden --><b/></a>") != "<a><b/></a>" {
		t.Fatal("CleanComment failed")
	}
	if got := CleanInvalidWithOptions("aXb", WithInvalidRegexp(regexp.MustCompile(`X`))); got != "ab" {
		t.Fatalf("CleanInvalidWithOptions = %q", got)
	}
	if got := CleanCommentWithOptions("<a>[hidden]<b/></a>", WithCommentRegexp(regexp.MustCompile(`\[[^]]+\]`))); got != "<a><b/></a>" {
		t.Fatalf("CleanCommentWithOptions = %q", got)
	}
	if Escape(`<a&"'>`) != "&lt;a&amp;&#34;&#39;&gt;" || Unescape("&lt;a&amp;&gt;") != "<a&>" {
		t.Fatal("escape/unescape failed")
	}

	fromBytes, err := ReadXMLBytes([]byte(`<root><a>1</a></root>`))
	if err != nil || ElementText(fromBytes.Root, "a") != "1" {
		t.Fatalf("ReadXMLBytes doc=%#v err=%v", fromBytes, err)
	}
	fromReader, err := ReadXMLReader(strings.NewReader(`<root><b>2</b></root>`))
	if err != nil || ElementText(fromReader.Root, "b") != "2" {
		t.Fatalf("ReadXMLReader doc=%#v err=%v", fromReader, err)
	}
	_, err = ParseXML(`<root><unclosed></root>`)
	assertXMLInvalidInput(t, err)
	if _, err := ReadXML(""); err == nil {
		t.Fatal("ReadXML empty path should return error")
	}

	var starts []string
	if err := ReadBySAX(strings.NewReader(`<root><a>1</a></root>`), func(tok stdxml.Token) error {
		if start, ok := tok.(stdxml.StartElement); ok {
			starts = append(starts, start.Name.Local)
		}
		return nil
	}); err != nil || !reflect.DeepEqual(starts, []string{"root", "a"}) {
		t.Fatalf("ReadBySAX starts=%v err=%v", starts, err)
	}
	if err := ReadBySAX(strings.NewReader(`<root>`), nil); err != nil {
		t.Fatalf("ReadBySAX nil handler should ignore input: %v", err)
	}
	var nsStarts []stdxml.Name
	if err := ReadBySAXWithOptions(strings.NewReader(`<root xmlns:p="urn:p"><p:a>1</p:a></root>`), func(tok stdxml.Token) error {
		if start, ok := tok.(stdxml.StartElement); ok {
			nsStarts = append(nsStarts, start.Name)
		}
		return nil
	}, WithNamespaceAware(false)); err != nil {
		t.Fatalf("ReadBySAXWithOptions namespace-aware false: %v", err)
	}
	if !reflect.DeepEqual(nsStarts, []stdxml.Name{{Local: "root"}, {Local: "a"}}) {
		t.Fatalf("ReadBySAXWithOptions names = %#v", nsStarts)
	}
	assertXMLInvalidInput(t, ReadBySAXWithOptions(strings.NewReader(`<root><a/></root>`), func(stdxml.Token) error { return nil }, WithMaxBytes(6)))
	handlerErr := errors.New("handler stop")
	if err := ReadBySAX(strings.NewReader(`<root/>`), func(stdxml.Token) error { return handlerErr }); !errors.Is(err, handlerErr) {
		t.Fatalf("ReadBySAX handler err = %v", err)
	}
	assertXMLInvalidInput(t, ReadBySAX(strings.NewReader(`<root>`), func(stdxml.Token) error { return nil }))
	decoderCalls := 0
	withDecoder, err := ReadXMLReader(strings.NewReader(`<ignored/>`), WithDecoderFactory(func(io.Reader) *stdxml.Decoder {
		decoderCalls++
		return stdxml.NewDecoder(strings.NewReader(`<root><factory>ok</factory></root>`))
	}))
	if err != nil || decoderCalls != 1 || ElementText(withDecoder.Root, "factory") != "ok" {
		t.Fatalf("ReadXMLReader decoder factory doc=%#v calls=%d err=%v", withDecoder, decoderCalls, err)
	}
	saxDecoderCalls := 0
	var factoryStarts []string
	if err := ReadBySAXWithOptions(strings.NewReader(`<ignored/>`), func(tok stdxml.Token) error {
		if start, ok := tok.(stdxml.StartElement); ok {
			factoryStarts = append(factoryStarts, start.Name.Local)
		}
		return nil
	}, WithDecoderFactory(func(io.Reader) *stdxml.Decoder {
		saxDecoderCalls++
		return stdxml.NewDecoder(strings.NewReader(`<root><factory/></root>`))
	})); err != nil || saxDecoderCalls != 1 || !reflect.DeepEqual(factoryStarts, []string{"root", "factory"}) {
		t.Fatalf("ReadBySAXWithOptions decoder factory starts=%v calls=%d err=%v", factoryStarts, saxDecoderCalls, err)
	}
	assertXMLInvalidInput(t, ReadBySAXWithOptions(strings.NewReader(`<root/>`), func(stdxml.Token) error { return nil }, WithDecoderFactory(func(io.Reader) *stdxml.Decoder { return nil })))

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

	var out strings.Builder
	if err := TransformWith(strings.NewReader(`<root><a>1</a></root>`), &out, WithOmitDeclaration(true)); err != nil || out.String() != `<root><a>1</a></root>` {
		t.Fatalf("TransformWith = %q, %v", out.String(), err)
	}
	out.Reset()
	if err := TransformWithOptions(strings.NewReader(`<root xmlns:p="urn:p"><p:a>1</p:a></root>`), &out,
		WithTransformParseOptions(WithNamespaceAware(false)),
		WithTransformWriteOptions(WithOmitDeclaration(true)),
	); err != nil || strings.Contains(out.String(), `xmlns`) || !strings.Contains(out.String(), `<a>1</a>`) {
		t.Fatalf("TransformWithOptions = %q, %v", out.String(), err)
	}
}
