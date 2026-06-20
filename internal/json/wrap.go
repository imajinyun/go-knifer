package json

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sort"
	"time"
)

// wrap converts any Go value into a JSON-compatible value: primitive, *JSONObject, *JSONArray, or Null.
func wrap(v any, cfg *Config) any {
	if cfg == nil {
		cfg = NewConfig()
	}
	if v == nil {
		return Null
	}
	switch x := v.(type) {
	case jsonNull, *jsonNull:
		return Null
	case *JSONObject, *JSONArray:
		return x
	case string, bool:
		return x
	case int:
		return int64(x)
	case int8:
		return int64(x)
	case int16:
		return int64(x)
	case int32:
		return int64(x)
	case int64:
		return x
	case uint:
		return int64(x)
	case uint8:
		return int64(x)
	case uint16:
		return int64(x)
	case uint32:
		return int64(x)
	case uint64:
		return int64(x)
	case float32:
		return float64(x)
	case float64:
		return x
	case json.Number:
		if i, err := x.Int64(); err == nil {
			return i
		}
		if f, err := x.Float64(); err == nil {
			return f
		}
		return x.String()
	case time.Time:
		if cfg.DateFormat == "" {
			return x.UnixMilli()
		}
		return x.Format(cfg.DateFormat)
	case []byte:
		// Treat []byte as a string.
		return string(x)
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer, reflect.Interface:
		if rv.IsNil() {
			return Null
		}
		return wrap(rv.Elem().Interface(), cfg)
	case reflect.Map:
		obj := NewJSONObjectWithConfig(cfg)
		// Only string keys are supported.
		keys := rv.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return cfg.sprint(keys[i].Interface()) < cfg.sprint(keys[j].Interface())
		})
		for _, key := range keys {
			k := cfg.sprint(key.Interface())
			obj.Set(k, wrap(rv.MapIndex(key).Interface(), cfg))
		}
		return obj
	case reflect.Slice, reflect.Array:
		arr := NewJSONArrayWithConfig(cfg)
		for i := 0; i < rv.Len(); i++ {
			arr.Add(wrap(rv.Index(i).Interface(), cfg))
		}
		return arr
	case reflect.Struct:
		// Unmarshal through encoding/json into generic structures before wrapping, so tags take effect.
		marshal := json.Marshal
		if cfg != nil && cfg.MarshalFunc != nil {
			marshal = cfg.MarshalFunc
		}
		b, err := marshal(v)
		if err != nil {
			return cfg.sprint(v)
		}
		if cfg != nil && cfg.UnmarshalFunc != nil {
			var raw any
			if err := cfg.UnmarshalFunc(b, &raw); err != nil {
				return cfg.sprint(v)
			}
			return wrap(raw, cfg)
		}
		dec := newDecoderWithConfig(bytes.NewReader(b), cfg)
		if dec == nil {
			return cfg.sprint(v)
		}
		var raw any
		if err := dec.Decode(&raw); err != nil {
			return cfg.sprint(v)
		}
		return wrap(raw, cfg)
	}
	return cfg.sprint(v)
}

// toString converts any JSON value to a string.
func toString(v any, def string, cfg *Config) string {
	cfg = configOrDefault(cfg)
	if IsNull(v) {
		return def
	}
	switch x := v.(type) {
	case string:
		return x
	case bool:
		if x {
			return "true"
		}
		return "false"
	case int64:
		return cfg.formatInt(x, 10)
	case float64:
		return cfg.formatFloat(x, 'f', -1, 64)
	case *JSONObject:
		return x.String()
	case *JSONArray:
		return x.String()
	}
	return cfg.sprint(v)
}

// toInt64 converts to int64 and returns def on failure.
func toInt64(v any, def int64, cfg *Config) int64 {
	cfg = configOrDefault(cfg)
	if IsNull(v) {
		return def
	}
	switch x := v.(type) {
	case int64:
		return x
	case float64:
		return int64(x)
	case bool:
		if x {
			return 1
		}
		return 0
	case string:
		n, err := cfg.parseInt(x, 10, 64)
		if err == nil {
			return n
		}
		f, err := cfg.parseFloat(x, 64)
		if err == nil {
			return int64(f)
		}
	}
	return def
}

// toFloat64 converts to float64 and returns def on failure.
func toFloat64(v any, def float64, cfg *Config) float64 {
	cfg = configOrDefault(cfg)
	if IsNull(v) {
		return def
	}
	switch x := v.(type) {
	case float64:
		return x
	case int64:
		return float64(x)
	case bool:
		if x {
			return 1
		}
		return 0
	case string:
		f, err := cfg.parseFloat(x, 64)
		if err == nil {
			return f
		}
	}
	return def
}

// toBool converts to bool.
func toBool(v any, def bool, cfg *Config) bool {
	cfg = configOrDefault(cfg)
	if IsNull(v) {
		return def
	}
	switch x := v.(type) {
	case bool:
		return x
	case int64:
		return x != 0
	case float64:
		return x != 0
	case string:
		if b, err := cfg.parseBool(x); err == nil {
			return b
		}
	}
	return def
}
