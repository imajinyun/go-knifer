package json

// jsonNull matches the utility JSONNull and represents JSON null.
type jsonNull struct{}

// Null is the singleton JSON null.
var Null = jsonNull{}

// String implements Stringer and returns "null".
func (jsonNull) String() string { return "null" }

// IsNull reports whether v is nil or JSON Null.
func IsNull(v any) bool {
	if v == nil {
		return true
	}
	_, ok := v.(jsonNull)
	return ok
}
