package json

import (
	"slices"
	"strconv"
	"strings"
)

// JSONObject matches the utility JSONObject and preserves key insertion order.
type JSONObject struct {
	cfg    *Config
	keys   []string
	values map[string]any
	// keyMap stores lowercase-to-original key mappings when ignoreCase is enabled.
	keyMap map[string]string
}

// NewJSONObject creates an empty object.
func NewJSONObject() *JSONObject {
	return NewJSONObjectWithConfig(nil)
}

// NewJSONObjectWithConfig creates an object with the specified config.
func NewJSONObjectWithConfig(cfg *Config) *JSONObject {
	if cfg == nil {
		cfg = NewConfig()
	}
	o := &JSONObject{cfg: cfg, values: map[string]any{}}
	if cfg.IgnoreCase {
		o.keyMap = map[string]string{}
	}
	return o
}

// Config returns the config.
func (o *JSONObject) Config() *Config { return o.cfg }

// Len returns the key count.
func (o *JSONObject) Len() int { return len(o.keys) }

// Keys returns a copy of keys in insertion order.
func (o *JSONObject) Keys() []string {
	out := slices.Clone(o.keys)
	if out == nil {
		return []string{}
	}
	return out
}

// canonicalKey returns the real key when ignoreCase is enabled; otherwise returns the input key.
func (o *JSONObject) canonicalKey(key string) (string, bool) {
	if o.cfg.IgnoreCase {
		if real, ok := o.keyMap[strings.ToLower(key)]; ok {
			return real, true
		}
		return key, false
	}
	_, ok := o.values[key]
	return key, ok
}

// Has reports whether key exists.
func (o *JSONObject) Has(key string) bool {
	_, ok := o.canonicalKey(key)
	return ok
}

// Get gets the raw value and returns nil, false when absent.
func (o *JSONObject) Get(key string) (any, bool) {
	real, ok := o.canonicalKey(key)
	if !ok {
		return nil, false
	}
	v, ok := o.values[real]
	return v, ok
}

// GetOrDefault gets a value or the default.
func (o *JSONObject) GetOrDefault(key string, def any) any {
	if v, ok := o.Get(key); ok {
		return v
	}
	return def
}

// IsNull reports whether key exists and its value is JSON null.
func (o *JSONObject) IsNull(key string) bool {
	v, ok := o.Get(key)
	if !ok {
		return false
	}
	return IsNull(v)
}

// Set sets a key-value pair and returns itself for chaining.
func (o *JSONObject) Set(key string, value any) *JSONObject {
	value = wrap(value, o.cfg)
	if o.cfg.IgnoreNullValue && IsNull(value) {
		return o
	}
	if real, ok := o.canonicalKey(key); ok {
		o.values[real] = value
		return o
	}
	o.keys = append(o.keys, key)
	o.values[key] = value
	if o.cfg.IgnoreCase {
		o.keyMap[strings.ToLower(key)] = key
	}
	return o
}

// Put is the same as Set for utility toolkit naming compatibility.
func (o *JSONObject) Put(key string, value any) *JSONObject { return o.Set(key, value) }

// Remove removes a key and reports whether it was removed.
func (o *JSONObject) Remove(key string) bool {
	real, ok := o.canonicalKey(key)
	if !ok {
		return false
	}
	delete(o.values, real)
	for i, k := range o.keys {
		if k == real {
			o.keys = append(o.keys[:i], o.keys[i+1:]...)
			break
		}
	}
	if o.cfg.IgnoreCase {
		delete(o.keyMap, strings.ToLower(real))
	}
	return true
}

// ForEach iterates in insertion order.
func (o *JSONObject) ForEach(fn func(key string, value any) bool) {
	for _, k := range o.keys {
		if !fn(k, o.values[k]) {
			return
		}
	}
}

// ToMap converts to a regular map whose values are raw JSON values.
func (o *JSONObject) ToMap() map[string]any {
	out := make(map[string]any, len(o.keys))
	for _, k := range o.keys {
		out[k] = o.values[k]
	}
	return out
}

// Typed getters.

// GetString gets a string and returns "" when absent or mismatched.
func (o *JSONObject) GetString(key string) string { return o.GetStringOr(key, "") }

// GetStringOr gets a string or the default.
func (o *JSONObject) GetStringOr(key, def string) string {
	v, ok := o.Get(key)
	if !ok {
		return def
	}
	return toString(v, def, o.cfg)
}

// GetInt gets an int.
func (o *JSONObject) GetInt(key string) int { return int(o.GetInt64Or(key, 0)) }

// GetIntOr gets an int or the default.
func (o *JSONObject) GetIntOr(key string, def int) int {
	return int(o.GetInt64Or(key, int64(def)))
}

// GetInt64 gets an int64.
func (o *JSONObject) GetInt64(key string) int64 { return o.GetInt64Or(key, 0) }

// GetInt64Or gets an int64 or the default.
func (o *JSONObject) GetInt64Or(key string, def int64) int64 {
	v, ok := o.Get(key)
	if !ok {
		return def
	}
	return toInt64(v, def, o.cfg)
}

// GetFloat64 gets a float64.
func (o *JSONObject) GetFloat64(key string) float64 { return o.GetFloat64Or(key, 0) }

// GetFloat64Or gets a float64 or the default.
func (o *JSONObject) GetFloat64Or(key string, def float64) float64 {
	v, ok := o.Get(key)
	if !ok {
		return def
	}
	return toFloat64(v, def, o.cfg)
}

// GetBool gets a bool.
func (o *JSONObject) GetBool(key string) bool { return o.GetBoolOr(key, false) }

// GetBoolOr gets a bool or the default.
func (o *JSONObject) GetBoolOr(key string, def bool) bool {
	v, ok := o.Get(key)
	if !ok {
		return def
	}
	return toBool(v, def, o.cfg)
}

// GetJSONObject gets a nested object and returns nil when absent or not an object.
func (o *JSONObject) GetJSONObject(key string) *JSONObject {
	v, ok := o.Get(key)
	if !ok {
		return nil
	}
	if obj, ok := v.(*JSONObject); ok {
		return obj
	}
	return nil
}

// GetJSONArray gets a nested array and returns nil when absent or not an array.
func (o *JSONObject) GetJSONArray(key string) *JSONArray {
	v, ok := o.Get(key)
	if !ok {
		return nil
	}
	if arr, ok := v.(*JSONArray); ok {
		return arr
	}
	return nil
}

// String returns a compact JSON string.
func (o *JSONObject) String() string {
	s, _ := writeValue(o, 0)
	return s
}

// ToString returns compact output.
func (o *JSONObject) ToString() string { return o.String() }

// ToStringPretty returns output indented with 4 spaces.
func (o *JSONObject) ToStringPretty() string {
	s, _ := writeValue(o, defaultIndent(o.cfg))
	return s
}

// MarshalJSON implements encoding/json.Marshaler.
func (o *JSONObject) MarshalJSON() ([]byte, error) {
	s, err := writeValue(o, 0)
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}

// UnmarshalJSON implements encoding/json.Unmarshaler.
func (o *JSONObject) UnmarshalJSON(b []byte) error {
	v, err := parseBytes(b)
	if err != nil {
		return err
	}
	parsed, ok := v.(*JSONObject)
	if !ok {
		return NewJSONError("expect json object, got %T", v)
	}
	o.cfg = parsed.cfg
	o.keys = parsed.keys
	o.values = parsed.values
	o.keyMap = parsed.keyMap
	return nil
}

// GetByPath reads a value through a path expression.
func (o *JSONObject) GetByPath(path string) any { return getByPath(o, path) }

// PutByPath writes a value through a path expression.
func (o *JSONObject) PutByPath(path string, value any) error { return putByPath(o, path, value) }

// defaultIndent returns the configured indentation, or 4 when it is 0.
func defaultIndent(cfg *Config) int {
	if cfg != nil && cfg.IndentFactor > 0 {
		return cfg.IndentFactor
	}
	return 4
}

// indexKey converts numeric keys to int for interoperation with array operations.
func parseIndex(s string) (int, bool) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return n, true
}
