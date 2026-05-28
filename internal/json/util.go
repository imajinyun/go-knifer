package json

import (
	"encoding/json"
	"strings"
)

// Parse 自动判断 JSON 类型：对象/数组/基础值。
func Parse(src any) (any, error) { return ParseWithConfig(src, nil) }

// ParseWithConfig 解析并使用配置。
func ParseWithConfig(src any, cfg *Config) (any, error) {
	switch x := src.(type) {
	case nil:
		return Null, nil
	case []byte:
		return parseBytesWithConfig(x, cfg)
	case string:
		return parseBytesWithConfig([]byte(x), cfg)
	case *JSONObject, *JSONArray:
		return x, nil
	}
	// 复杂类型：先 wrap 再返回
	return wrap(src, configOrDefault(cfg)), nil
}

// ParseObj 强制解析为 JSONObject。
func ParseObj(src any) (*JSONObject, error) { return ParseObjWithConfig(src, nil) }

// ParseObjWithConfig 解析为 JSONObject。
func ParseObjWithConfig(src any, cfg *Config) (*JSONObject, error) {
	v, err := ParseWithConfig(src, cfg)
	if err != nil {
		return nil, err
	}
	if obj, ok := v.(*JSONObject); ok {
		return obj, nil
	}
	return nil, NewJSONError("expect json object, got %T", v)
}

// ParseArray 强制解析为 JSONArray。
func ParseArray(src any) (*JSONArray, error) { return ParseArrayWithConfig(src, nil) }

// ParseArrayWithConfig 解析为 JSONArray。
func ParseArrayWithConfig(src any, cfg *Config) (*JSONArray, error) {
	v, err := ParseWithConfig(src, cfg)
	if err != nil {
		return nil, err
	}
	if arr, ok := v.(*JSONArray); ok {
		return arr, nil
	}
	return nil, NewJSONError("expect json array, got %T", v)
}

// ToJSONStr 紧凑序列化。
func ToJSONStr(v any) (string, error) {
	w := wrap(v, NewConfig())
	return writeValue(w, 0)
}

// ToJSONPrettyStr 4 空格缩进序列化。
func ToJSONPrettyStr(v any) (string, error) {
	w := wrap(v, NewConfig())
	return writeValue(w, 4)
}

// ToJSONStrIndent 自定义缩进序列化。
func ToJSONStrIndent(v any, indent int) (string, error) {
	w := wrap(v, NewConfig())
	return writeValue(w, indent)
}

// IsJSON 检查字符串是否合法 JSON。
func IsJSON(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	return json.Valid([]byte(s))
}

// IsJSONObj 检查字符串是否是 JSON 对象。
func IsJSONObj(s string) bool {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "{") || !strings.HasSuffix(s, "}") {
		return false
	}
	return IsJSON(s)
}

// IsJSONArray 检查字符串是否是 JSON 数组。
func IsJSONArray(s string) bool {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "[") || !strings.HasSuffix(s, "]") {
		return false
	}
	return IsJSON(s)
}

// GetByPath 顶层路径查询。
func GetByPath(root any, path string) any { return getByPath(root, path) }

// GetByPathOr 顶层路径查询，缺省回退。
func GetByPathOr(root any, path string, def any) any {
	if v := getByPath(root, path); v != nil && !IsNull(v) {
		return v
	}
	return def
}

// PutByPath 顶层路径写入。
func PutByPath(root any, path string, value any) error { return putByPath(root, path, value) }

// Quote 在 JSON 字符串两侧加引号并进行必要转义。
func Quote(s string) string {
	var sb strings.Builder
	writeQuoted(&sb, s)
	return sb.String()
}

// configOrDefault 返回非空配置。
func configOrDefault(cfg *Config) *Config {
	if cfg == nil {
		return NewConfig()
	}
	return cfg
}
