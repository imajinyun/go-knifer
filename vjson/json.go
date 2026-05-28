package vjson

import jsonx "github.com/imajinyun/go-knifer/internal/json"

// JSONObject 是有序的 JSON 对象。
type JSONObject = jsonx.JSONObject

// JSONArray 是有序的 JSON 数组。
type JSONArray = jsonx.JSONArray

// JSONConfig 控制 JSON 序列化行为。
type JSONConfig = jsonx.Config

// JSONError 是 JSON 模块错误类型。
type JSONError = jsonx.JSONError

// JSONNull 是 JSON null 单例值。
var JSONNull = jsonx.Null

// NewJSONObject 创建一个空的有序 JSON 对象。
func NewJSONObject() *JSONObject { return jsonx.NewJSONObject() }

// NewJSONObjectWithConfig 使用指定配置创建 JSON 对象。
func NewJSONObjectWithConfig(cfg *JSONConfig) *JSONObject {
	return jsonx.NewJSONObjectWithConfig(cfg)
}

// NewJSONArray 创建空的有序 JSON 数组。
func NewJSONArray() *JSONArray { return jsonx.NewJSONArray() }

// NewJSONArrayWithConfig 使用指定配置创建 JSON 数组。
func NewJSONArrayWithConfig(cfg *JSONConfig) *JSONArray {
	return jsonx.NewJSONArrayWithConfig(cfg)
}

// NewJSONConfig 创建默认 JSON 配置。
func NewJSONConfig() *JSONConfig { return jsonx.NewConfig() }

// JSONIsNull 判断值是否为 nil 或 JSON null。
func JSONIsNull(v any) bool { return jsonx.IsNull(v) }

// JSONParse 自动判断并解析 JSON。
func JSONParse(src any) (any, error) { return jsonx.Parse(src) }

// JSONParseObj 强制解析为 JSONObject。
func JSONParseObj(src any) (*JSONObject, error) { return jsonx.ParseObj(src) }

// JSONParseArray 强制解析为 JSONArray。
func JSONParseArray(src any) (*JSONArray, error) { return jsonx.ParseArray(src) }

// JSONToStr 紧凑序列化任意值为 JSON。
func JSONToStr(v any) (string, error) { return jsonx.ToJSONStr(v) }

// JSONToPrettyStr 4 空格缩进序列化。
func JSONToPrettyStr(v any) (string, error) { return jsonx.ToJSONPrettyStr(v) }

// JSONToStrIndent 自定义缩进序列化。
func JSONToStrIndent(v any, indent int) (string, error) {
	return jsonx.ToJSONStrIndent(v, indent)
}

// JSONFormat 对原始 JSON 字符串重排版。
func JSONFormat(raw string) string { return jsonx.FormatJSONStr(raw) }

// JSONIsJSON 判断字符串是否为合法 JSON。
func JSONIsJSON(s string) bool { return jsonx.IsJSON(s) }

// JSONIsObj 判断字符串是否为 JSON 对象。
func JSONIsObj(s string) bool { return jsonx.IsJSONObj(s) }

// JSONIsArray 判断字符串是否为 JSON 数组。
func JSONIsArray(s string) bool { return jsonx.IsJSONArray(s) }

// JSONGetByPath 路径表达式取值。
func JSONGetByPath(root any, path string) any { return jsonx.GetByPath(root, path) }

// JSONGetByPathOr 路径表达式取值并提供默认值。
func JSONGetByPathOr(root any, path string, def any) any {
	return jsonx.GetByPathOr(root, path, def)
}

// JSONPutByPath 路径表达式写入。
func JSONPutByPath(root any, path string, value any) error {
	return jsonx.PutByPath(root, path, value)
}

// JSONQuote 给字符串添加 JSON 双引号并转义。
func JSONQuote(s string) string { return jsonx.Quote(s) }

// JSONToBean 将 JSON 反序列化到 dst（必须是指针）。
func JSONToBean(src any, dst any) error { return jsonx.ToBean(src, dst) }

// JSONToList 将 JSON 数组反序列化到 dst（必须是指向 slice 的指针）。
func JSONToList(src any, dst any) error { return jsonx.ToList(src, dst) }

// XMLToJSON 将 XML 字符串解析为 JSONObject。
func XMLToJSON(xmlStr string) (*JSONObject, error) { return jsonx.XMLToJSON(xmlStr) }

// JSONToXML 将 JSON 值序列化为 XML 字符串，rootTag 为空时直接拼接键。
func JSONToXML(root any, rootTag string) (string, error) {
	return jsonx.JSONToXML(root, rootTag)
}
