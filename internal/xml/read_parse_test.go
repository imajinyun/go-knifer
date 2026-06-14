package xml

import (
	stdxml "encoding/xml"
	"io"
	"strings"
	"testing"
)

func TestReadXMLBytesAndReader(t *testing.T) {
	fromBytes, err := ReadXMLBytes([]byte(`<root><a>1</a></root>`))
	if err != nil || ElementText(fromBytes.Root, "a") != "1" {
		t.Fatalf("ReadXMLBytes doc=%#v err=%v", fromBytes, err)
	}
	fromReader, err := ReadXMLReader(strings.NewReader(`<root><b>2</b></root>`))
	if err != nil || ElementText(fromReader.Root, "b") != "2" {
		t.Fatalf("ReadXMLReader doc=%#v err=%v", fromReader, err)
	}
}

func TestReadXMLReaderWithDecoderFactory(t *testing.T) {
	decoderCalls := 0
	withDecoder, err := ReadXMLReader(strings.NewReader(`<ignored/>`), WithDecoderFactory(func(io.Reader) *stdxml.Decoder {
		decoderCalls++
		return stdxml.NewDecoder(strings.NewReader(`<root><factory>ok</factory></root>`))
	}))
	if err != nil || decoderCalls != 1 || ElementText(withDecoder.Root, "factory") != "ok" {
		t.Fatalf("ReadXMLReader decoder factory doc=%#v calls=%d err=%v", withDecoder, decoderCalls, err)
	}
}

func TestParseXMLInvalidInputAndEmptyPath(t *testing.T) {
	_, err := ParseXML(`<root><unclosed></root>`)
	assertXMLInvalidInput(t, err)
	if _, err := ReadXML(""); err == nil {
		t.Fatal("ReadXML empty path should return error")
	}
}
