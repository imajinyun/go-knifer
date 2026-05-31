package vxml

import (
	stdxml "encoding/xml"
	"io"

	xmlimpl "github.com/imajinyun/go-knifer/internal/xml"
)

const (
	NBSP          = xmlimpl.NBSP
	AMP           = xmlimpl.AMP
	QUOTE         = xmlimpl.QUOTE
	APOS          = xmlimpl.APOS
	LT            = xmlimpl.LT
	GT            = xmlimpl.GT
	InvalidRegex  = xmlimpl.InvalidRegex
	CommentRegex  = xmlimpl.CommentRegex
	IndentDefault = xmlimpl.IndentDefault
	ContentKey    = xmlimpl.ContentKey
)

type (
	Document       = xmlimpl.Document
	Element        = xmlimpl.Element
	TokenHandler   = xmlimpl.TokenHandler
	NamespaceCache = xmlimpl.NamespaceCache
	ParseOption    = xmlimpl.ParseOption
)

func DisableDefaultDocumentBuilderFactory() { xmlimpl.DisableDefaultDocumentBuilderFactory() }

// SetNamespaceAware records whether parsed element names should keep namespace URIs.
//
// Deprecated: prefer ReadXMLReaderWithOptions or ParseXMLWithOptions to avoid
// package-level mutable state in concurrent callers.
func SetNamespaceAware(isNamespaceAware bool) {
	xmlimpl.SetNamespaceAwareForCompatibility(isNamespaceAware)
}

func WithNamespaceAware(isNamespaceAware bool) ParseOption {
	return xmlimpl.WithNamespaceAware(isNamespaceAware)
}
func ReadXML(pathOrContent string) (*Document, error) { return xmlimpl.ReadXML(pathOrContent) }
func ReadXMLFile(path string) (*Document, error)      { return xmlimpl.ReadXMLFile(path) }
func ReadXMLBytes(data []byte) (*Document, error)     { return xmlimpl.ReadXMLBytes(data) }
func ReadXMLReader(r io.Reader) (*Document, error)    { return xmlimpl.ReadXMLReader(r) }
func ReadXMLReaderWithOptions(r io.Reader, opts ...ParseOption) (*Document, error) {
	return xmlimpl.ReadXMLReaderWithOptions(r, opts...)
}
func ParseXML(xmlStr string) (*Document, error) { return xmlimpl.ParseXML(xmlStr) }
func ParseXMLWithOptions(xmlStr string, opts ...ParseOption) (*Document, error) {
	return xmlimpl.ParseXMLWithOptions(xmlStr, opts...)
}
func ReadBySAX(r io.Reader, handler TokenHandler) error { return xmlimpl.ReadBySAX(r, handler) }
func ReadBySAXFile(path string, handler TokenHandler) error {
	return xmlimpl.ReadBySAXFile(path, handler)
}
func ToStr(v any) string                    { return xmlimpl.ToStr(v) }
func ToStrPretty(v any, pretty bool) string { return xmlimpl.ToStrPretty(v, pretty) }
func ToStrCharset(v any, charset string, pretty bool, omitXMLDeclaration bool) (string, error) {
	return xmlimpl.ToStrCharset(v, charset, pretty, omitXMLDeclaration)
}
func Format(xmlStr string) (string, error)            { return xmlimpl.Format(xmlStr) }
func ToFile(v any, path string, charset string) error { return xmlimpl.ToFile(v, path, charset) }
func Write(v any, w io.Writer, charset string, indentSize int, omitXMLDeclaration bool) error {
	return xmlimpl.Write(v, w, charset, indentSize, omitXMLDeclaration)
}

func Transform(source io.Reader, result io.Writer, charset string, indentSize int, omitXMLDeclaration bool) error {
	return xmlimpl.Transform(source, result, charset, indentSize, omitXMLDeclaration)
}
func CreateXML() *Document { return xmlimpl.CreateXML() }
func CreateXMLWithRoot(rootElementName string) *Document {
	return xmlimpl.CreateXMLWithRoot(rootElementName)
}

func CreateXMLWithRootNS(rootElementName, namespace string) *Document {
	return xmlimpl.CreateXMLWithRootNS(rootElementName, namespace)
}
func CreateDocumentBuilder() struct{}          { return xmlimpl.CreateDocumentBuilder() }
func CreateDocumentBuilderFactory() struct{}   { return xmlimpl.CreateDocumentBuilderFactory() }
func GetRootElement(doc *Document) *Element    { return xmlimpl.GetRootElement(doc) }
func GetOwnerDocument(node *Element) *Document { return xmlimpl.GetOwnerDocument(node) }
func CleanInvalid(xmlContent string) string    { return xmlimpl.CleanInvalid(xmlContent) }
func CleanComment(xmlContent string) string    { return xmlimpl.CleanComment(xmlContent) }
func GetElements(element *Element, tagName string) []*Element {
	return xmlimpl.GetElements(element, tagName)
}

func GetElement(element *Element, tagName string) *Element {
	return xmlimpl.GetElement(element, tagName)
}

func ElementText(element *Element, tagName string, defaultValue ...string) string {
	return xmlimpl.ElementText(element, tagName, defaultValue...)
}
func TransElements(nodes []*Element) []*Element    { return xmlimpl.TransElements(nodes) }
func WriteObjectAsXML(path string, bean any) error { return xmlimpl.WriteObjectAsXML(path, bean) }
func CreateXPath() struct{}                        { return xmlimpl.CreateXPath() }
func GetElementByXPath(expression string, source any) *Element {
	return xmlimpl.GetElementByXPath(expression, source)
}

func GetNodeListByXPath(expression string, source any) []*Element {
	return xmlimpl.GetNodeListByXPath(expression, source)
}

func GetNodeByXPath(expression string, source any) *Element {
	return xmlimpl.GetNodeByXPath(expression, source)
}

func GetByXPath(expression string, source any, returnType string) any {
	return xmlimpl.GetByXPath(expression, source, returnType)
}
func Escape(s string) string                         { return xmlimpl.Escape(s) }
func Unescape(s string) string                       { return xmlimpl.Unescape(s) }
func XMLToMap(xmlStr string) (map[string]any, error) { return xmlimpl.XMLToMap(xmlStr) }
func XMLNodeToMap(node *Element) map[string]any      { return xmlimpl.XMLNodeToMap(node) }
func XMLToMapInto(xmlStr string, result map[string]any) (map[string]any, error) {
	return xmlimpl.XMLToMapInto(xmlStr, result)
}

func XMLNodeToMapInto(node *Element, result map[string]any) map[string]any {
	return xmlimpl.XMLNodeToMapInto(node, result)
}
func XMLToBean(xmlStr string, dst any) error     { return xmlimpl.XMLToBean(xmlStr, dst) }
func XMLNodeToBean(node *Element, dst any) error { return xmlimpl.XMLNodeToBean(node, dst) }
func MapToXMLStr(data map[string]any, rootName ...string) (string, error) {
	return xmlimpl.MapToXMLStr(data, rootName...)
}

func MapToXMLStrOptions(data map[string]any, rootName, namespace string, pretty, omitXMLDeclaration bool, charset string) (string, error) {
	return xmlimpl.MapToXMLStrOptions(data, rootName, namespace, pretty, omitXMLDeclaration, charset)
}

func MapToXML(data map[string]any, rootName string, namespace ...string) *Document {
	return xmlimpl.MapToXML(data, rootName, namespace...)
}
func BeanToXML(bean any, namespace ...string) *Document { return xmlimpl.BeanToXML(bean, namespace...) }
func BeanToXMLWithOptions(bean any, namespace string, ignoreNull bool) *Document {
	return xmlimpl.BeanToXMLWithOptions(bean, namespace, ignoreNull)
}
func IsElement(node *Element) bool { return xmlimpl.IsElement(node) }
func AppendChild(node *Element, tagName string, namespace ...string) *Element {
	return xmlimpl.AppendChild(node, tagName, namespace...)
}
func AppendText(node *Element, text any) *Element     { return xmlimpl.AppendText(node, text) }
func Append(node *Element, data any)                  { xmlimpl.Append(node, data) }
func NewNamespaceCache(doc *Document) *NamespaceCache { return xmlimpl.NewNamespaceCache(doc) }

func XMLName(local string) stdxml.Name { return stdxml.Name{Local: local} }
