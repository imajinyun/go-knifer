package json

import (
	"slices"
	"strings"
)

// JSONArray matches the utility JSONArray as an ordered list of JSON values.
type JSONArray struct {
	cfg    *Config
	values []any
}

// NewJSONArray creates an empty array.
func NewJSONArray() *JSONArray { return NewJSONArrayWithConfig(nil) }

// NewJSONArrayWithConfig creates a value with config.
func NewJSONArrayWithConfig(cfg *Config) *JSONArray {
	if cfg == nil {
		cfg = NewConfig()
	}
	return &JSONArray{cfg: cfg}
}

// Config returns the config.
func (a *JSONArray) Config() *Config { return a.cfg }

// Len returns the element count.
func (a *JSONArray) Len() int { return len(a.values) }

// Get accesses an index and returns nil, false when out of range.
func (a *JSONArray) Get(i int) (any, bool) {
	if i < 0 || i >= len(a.values) {
		return nil, false
	}
	return a.values[i], true
}

// GetOrDefault returns the default when out of range.
func (a *JSONArray) GetOrDefault(i int, def any) any {
	if v, ok := a.Get(i); ok {
		return v
	}
	return def
}

// IsNull reports whether the indexed value is JSON null.
func (a *JSONArray) IsNull(i int) bool {
	v, ok := a.Get(i)
	if !ok {
		return false
	}
	return IsNull(v)
}

// Add appends an element.
func (a *JSONArray) Add(value any) *JSONArray {
	v := wrap(value, a.cfg)
	if a.cfg.IgnoreNullValue && IsNull(v) {
		return a
	}
	a.values = append(a.values, v)
	return a
}

// AddAll appends multiple elements.
func (a *JSONArray) AddAll(values ...any) *JSONArray {
	for _, v := range values {
		a.Add(v)
	}
	return a
}

// Set writes at the given index and fills gaps with nil.
func (a *JSONArray) Set(i int, value any) *JSONArray {
	v := wrap(value, a.cfg)
	for len(a.values) <= i {
		a.values = append(a.values, Null)
	}
	a.values[i] = v
	return a
}

// Insert inserts at index i.
func (a *JSONArray) Insert(i int, value any) *JSONArray {
	v := wrap(value, a.cfg)
	if i < 0 {
		i = 0
	}
	if i >= len(a.values) {
		a.values = append(a.values, v)
		return a
	}
	a.values = append(a.values, nil)
	copy(a.values[i+1:], a.values[i:])
	a.values[i] = v
	return a
}

// Remove removes the indexed element and returns false when out of range.
func (a *JSONArray) Remove(i int) bool {
	if i < 0 || i >= len(a.values) {
		return false
	}
	a.values = append(a.values[:i], a.values[i+1:]...)
	return true
}

// Range iterates elements in order.
func (a *JSONArray) Range(fn func(i int, v any) bool) {
	for i, v := range a.values {
		if !fn(i, v) {
			return
		}
	}
}

// ToSlice converts to []any.
func (a *JSONArray) ToSlice() []any {
	out := slices.Clone(a.values)
	if out == nil {
		return []any{}
	}
	return out
}

// Join joins all elements with separator after converting values through toString.
func (a *JSONArray) Join(sep string) string {
	if len(a.values) == 0 {
		return ""
	}
	var b strings.Builder
	for i, v := range a.values {
		if i > 0 {
			b.WriteString(sep)
		}
		b.WriteString(toString(v, "", a.cfg))
	}
	return b.String()
}

// Typed getters.

// GetString gets a string by index.
func (a *JSONArray) GetString(i int) string { return a.GetStringOr(i, "") }

// GetStringOr gets a string by index or the default.
func (a *JSONArray) GetStringOr(i int, def string) string {
	v, ok := a.Get(i)
	if !ok {
		return def
	}
	return toString(v, def, a.cfg)
}

// GetInt gets an int by index.
func (a *JSONArray) GetInt(i int) int { return int(a.GetInt64Or(i, 0)) }

// GetInt64 gets an int64 by index.
func (a *JSONArray) GetInt64(i int) int64 { return a.GetInt64Or(i, 0) }

// GetInt64Or gets an int64 by index or the default.
func (a *JSONArray) GetInt64Or(i int, def int64) int64 {
	v, ok := a.Get(i)
	if !ok {
		return def
	}
	return toInt64(v, def, a.cfg)
}

// GetFloat64 gets a float64 by index.
func (a *JSONArray) GetFloat64(i int) float64 { return a.GetFloat64Or(i, 0) }

// GetFloat64Or gets a float64 by index or the default.
func (a *JSONArray) GetFloat64Or(i int, def float64) float64 {
	v, ok := a.Get(i)
	if !ok {
		return def
	}
	return toFloat64(v, def, a.cfg)
}

// GetBool gets a bool by index.
func (a *JSONArray) GetBool(i int) bool { return a.GetBoolOr(i, false) }

// GetBoolOr gets a bool by index or the default.
func (a *JSONArray) GetBoolOr(i int, def bool) bool {
	v, ok := a.Get(i)
	if !ok {
		return def
	}
	return toBool(v, def, a.cfg)
}

// GetJSONObject gets a JSONObject by index.
func (a *JSONArray) GetJSONObject(i int) *JSONObject {
	v, ok := a.Get(i)
	if !ok {
		return nil
	}
	if obj, ok := v.(*JSONObject); ok {
		return obj
	}
	return nil
}

// GetJSONArray gets a JSONArray by index.
func (a *JSONArray) GetJSONArray(i int) *JSONArray {
	v, ok := a.Get(i)
	if !ok {
		return nil
	}
	if arr, ok := v.(*JSONArray); ok {
		return arr
	}
	return nil
}

// String returns compact output.
func (a *JSONArray) String() string {
	s, _ := writeValue(a, 0)
	return s
}

// ToString returns compact output.
func (a *JSONArray) ToString() string { return a.String() }

// ToStringPretty returns output indented with 4 spaces.
func (a *JSONArray) ToStringPretty() string {
	s, _ := writeValue(a, defaultIndent(a.cfg))
	return s
}

// MarshalJSON implements encoding/json.Marshaler.
func (a *JSONArray) MarshalJSON() ([]byte, error) {
	s, err := writeValue(a, 0)
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}

// UnmarshalJSON implements encoding/json.Unmarshaler.
func (a *JSONArray) UnmarshalJSON(b []byte) error {
	v, err := parseBytes(b)
	if err != nil {
		return err
	}
	arr, ok := v.(*JSONArray)
	if !ok {
		return NewJSONError("expect json array, got %T", v)
	}
	a.cfg = arr.cfg
	a.values = arr.values
	return nil
}

// GetByPath reads by path.
func (a *JSONArray) GetByPath(path string) any { return getByPath(a, path) }

// PutByPath writes by path.
func (a *JSONArray) PutByPath(path string, value any) error { return putByPath(a, path, value) }
