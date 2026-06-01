// Package vxml is the public facade over the internal XML implementation.
// All APIs are concurrent-safe and use per-call options instead of global state.
package vxml

import (
	stdxml "encoding/xml"
	"io"

	xmlimpl "github.com/imajinyun/go-knifer/internal/xml"
)

// ---------------------------------------------------------------------------
// Constants
// ---------------------------------------------------------------------------

const (
	NBSP           = xmlimpl.NBSP
	AMP            = xmlimpl.AMP
	QUOTE          = xmlimpl.QUOTE
	APOS           = xmlimpl.APOS
	LT             = xmlimpl.LT
	GT             = xmlimpl.GT
	InvalidRegex   = xmlimpl.InvalidRegex
	CommentRegex   = xmlimpl.CommentRegex
	IndentDefault  = xmlimpl.IndentDefault
	ContentKey     = xmlimpl.ContentKey
	DefaultCharset = xmlimpl.DefaultCharset
)

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

type (
	Document       = xmlimpl.Document
	Element        = xmlimpl.Element
	TokenHandler   = xmlimpl.TokenHandler
	NamespaceCache = xmlimpl.NamespaceCache
	ParseOption    = xmlimpl.ParseOption
	WriteOption    = xmlimpl.WriteOption
)

// ---------------------------------------------------------------------------
// Options
// ---------------------------------------------------------------------------

// WithNamespaceAware controls whether parsed element names keep namespace URIs.
func WithNamespaceAware(b bool) ParseOption { return xmlimpl.WithNamespaceAware(b) }

// WithCharset sets the XML declaration charset.
func WithCharset(s string) WriteOption { return xmlimpl.WithCharset(s) }

// WithIndent sets the indentation width in spaces (0 disables pretty printing).
func WithIndent(n int) WriteOption { return xmlimpl.WithIndent(n) }

// WithPretty enables pretty printing with the default indentation.
func WithPretty() WriteOption { return xmlimpl.WithPretty() }

// WithOmitDeclaration controls whether the <?xml ... ?> prolog is emitted.
func WithOmitDeclaration(b bool) WriteOption { return xmlimpl.WithOmitDeclaration(b) }

// WithIgnoreNullFields skips struct fields whose value is a typed nil.
func WithIgnoreNullFields(b bool) WriteOption { return xmlimpl.WithIgnoreNullFields(b) }

// WithRootName overrides the synthesized root element name for MarshalMap / MarshalBean.
func WithRootName(s string) WriteOption { return xmlimpl.WithRootName(s) }

// WithNamespace sets the xmlns attribute on the synthesized root element.
func WithNamespace(s string) WriteOption { return xmlimpl.WithNamespace(s) }

// ---------------------------------------------------------------------------
// Reading
// ---------------------------------------------------------------------------

// ReadXML parses XML content directly, or treats the input as a file path when
// it does not start with '<'.
func ReadXML(pathOrContent string, opts ...ParseOption) (*Document, error) {
	return xmlimpl.ReadXML(pathOrContent, opts...)
}

// ReadXMLFile parses an XML file.
func ReadXMLFile(path string, opts ...ParseOption) (*Document, error) {
	return xmlimpl.ReadXMLFile(path, opts...)
}

// ReadXMLBytes parses XML bytes.
func ReadXMLBytes(data []byte, opts ...ParseOption) (*Document, error) {
	return xmlimpl.ReadXMLBytes(data, opts...)
}

// ReadXMLReader parses XML from reader.
func ReadXMLReader(r io.Reader, opts ...ParseOption) (*Document, error) {
	return xmlimpl.ReadXMLReader(r, opts...)
}

// ParseXML parses an XML string.
func ParseXML(xmlStr string, opts ...ParseOption) (*Document, error) {
	return xmlimpl.ParseXML(xmlStr, opts...)
}

// ReadBySAX streams XML tokens from reader to handler.
func ReadBySAX(r io.Reader, handler TokenHandler) error { return xmlimpl.ReadBySAX(r, handler) }

// ReadBySAXFile streams XML tokens from file.
func ReadBySAXFile(path string, handler TokenHandler) error {
	return xmlimpl.ReadBySAXFile(path, handler)
}

// ---------------------------------------------------------------------------
// Writing
// ---------------------------------------------------------------------------

// WriteTo serializes a document or element to writer.
func WriteTo(w io.Writer, v any, opts ...WriteOption) error {
	return xmlimpl.WriteTo(w, v, opts...)
}

// MarshalString serializes a document or element to string.
func MarshalString(v any, opts ...WriteOption) (string, error) {
	return xmlimpl.MarshalString(v, opts...)
}

// WriteFile writes a document or element to path.
func WriteFile(path string, v any, opts ...WriteOption) error {
	return xmlimpl.WriteFile(path, v, opts...)
}

// MarshalMap serializes map data to an XML string.
func MarshalMap(data map[string]any, opts ...WriteOption) (string, error) {
	return xmlimpl.MarshalMap(data, opts...)
}

// MarshalBean serializes a struct or map-like value to an XML string.
func MarshalBean(bean any, opts ...WriteOption) (string, error) {
	return xmlimpl.MarshalBean(bean, opts...)
}

// TransformWith copies XML from source to result with per-call options.
func TransformWith(source io.Reader, result io.Writer, opts ...WriteOption) error {
	return xmlimpl.TransformWith(source, result, opts...)
}

// Format pretty prints XML content.
func Format(xmlStr string) (string, error) { return xmlimpl.Format(xmlStr) }

// ---------------------------------------------------------------------------
// Element construction & traversal
// ---------------------------------------------------------------------------

// CreateXML creates an empty XML document.
func CreateXML() *Document { return xmlimpl.CreateXML() }

// CreateXMLWithRoot creates an XML document with root element.
func CreateXMLWithRoot(rootElementName string) *Document {
	return xmlimpl.CreateXMLWithRoot(rootElementName)
}

// CreateXMLWithRootNS creates an XML document with root element and namespace URI.
func CreateXMLWithRootNS(rootElementName, namespace string) *Document {
	return xmlimpl.CreateXMLWithRootNS(rootElementName, namespace)
}

// GetRootElement returns the document root element.
func GetRootElement(doc *Document) *Element { return xmlimpl.GetRootElement(doc) }

// GetOwnerDocument returns the document that owns node by walking to the root.
func GetOwnerDocument(node *Element) *Document { return xmlimpl.GetOwnerDocument(node) }

// CleanInvalid removes XML 1.0 invalid control characters.
func CleanInvalid(xmlContent string) string { return xmlimpl.CleanInvalid(xmlContent) }

// CleanComment removes XML comments.
func CleanComment(xmlContent string) string { return xmlimpl.CleanComment(xmlContent) }

// GetElements returns child elements with tag name. Empty tagName returns all direct children.
func GetElements(element *Element, tagName string) []*Element {
	return xmlimpl.GetElements(element, tagName)
}

// GetElement returns the first child element with tag name.
func GetElement(element *Element, tagName string) *Element {
	return xmlimpl.GetElement(element, tagName)
}

// ElementText returns child text or defaultValue when missing.
func ElementText(element *Element, tagName string, defaultValue ...string) string {
	return xmlimpl.ElementText(element, tagName, defaultValue...)
}

// TransElements returns the input list without nil elements.
func TransElements(nodes []*Element) []*Element { return xmlimpl.TransElements(nodes) }

// IsElement reports whether node is not nil.
func IsElement(node *Element) bool { return xmlimpl.IsElement(node) }

// AppendChild appends and returns a child element.
func AppendChild(node *Element, tagName string, namespace ...string) *Element {
	return xmlimpl.AppendChild(node, tagName, namespace...)
}

// AppendText appends text to an element.
func AppendText(node *Element, text any) *Element { return xmlimpl.AppendText(node, text) }

// Append appends map, slice, struct, or scalar data to node.
func Append(node *Element, data any) { xmlimpl.Append(node, data) }

// XMLName builds a stdxml.Name from local name.
func XMLName(local string) stdxml.Name { return stdxml.Name{Local: local} }

// ---------------------------------------------------------------------------
// XPath (simple subset)
// ---------------------------------------------------------------------------

// GetElementByXPath returns the first element matched by a simple expression.
func GetElementByXPath(expression string, source any) *Element {
	return xmlimpl.GetElementByXPath(expression, source)
}

// GetNodeListByXPath returns all elements matched by a simple expression.
func GetNodeListByXPath(expression string, source any) []*Element {
	return xmlimpl.GetNodeListByXPath(expression, source)
}

// GetNodeByXPath returns the first node matched by a simple expression.
func GetNodeByXPath(expression string, source any) *Element {
	return xmlimpl.GetNodeByXPath(expression, source)
}

// GetByXPath returns matched text, element, or list based on returnType.
func GetByXPath(expression string, source any, returnType string) any {
	return xmlimpl.GetByXPath(expression, source, returnType)
}

// ---------------------------------------------------------------------------
// Escape / unescape
// ---------------------------------------------------------------------------

// Escape escapes XML text.
func Escape(s string) string { return xmlimpl.Escape(s) }

// Unescape unescapes XML/HTML entities.
func Unescape(s string) string { return xmlimpl.Unescape(s) }

// ---------------------------------------------------------------------------
// XML <-> Map / Bean
// ---------------------------------------------------------------------------

// XMLToMap parses XML into a nested map. Repeated sibling tags become []any.
func XMLToMap(xmlStr string) (map[string]any, error) { return xmlimpl.XMLToMap(xmlStr) }

// XMLNodeToMap converts an element into a nested map value.
func XMLNodeToMap(node *Element) map[string]any { return xmlimpl.XMLNodeToMap(node) }

// XMLToMapInto parses XML and merges values into result.
func XMLToMapInto(xmlStr string, result map[string]any) (map[string]any, error) {
	return xmlimpl.XMLToMapInto(xmlStr, result)
}

// XMLNodeToMapInto converts an element to map and merges values into result.
func XMLNodeToMapInto(node *Element, result map[string]any) map[string]any {
	return xmlimpl.XMLNodeToMapInto(node, result)
}

// XMLToBean parses XML and decodes the generated map into dst.
func XMLToBean(xmlStr string, dst any) error { return xmlimpl.XMLToBean(xmlStr, dst) }

// XMLNodeToBean converts an element tree to a map and decodes it into dst.
func XMLNodeToBean(node *Element, dst any) error { return xmlimpl.XMLNodeToBean(node, dst) }

// ---------------------------------------------------------------------------
// Namespace cache
// ---------------------------------------------------------------------------

// NewNamespaceCache collects namespace declarations from doc.
func NewNamespaceCache(doc *Document) *NamespaceCache { return xmlimpl.NewNamespaceCache(doc) }
