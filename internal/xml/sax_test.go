package xml

import (
	stdxml "encoding/xml"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestReadBySAXStartElements(t *testing.T) {
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
}

func TestReadBySAXOptionsAndErrors(t *testing.T) {
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
}

func TestReadBySAXWithDecoderFactory(t *testing.T) {
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
}
