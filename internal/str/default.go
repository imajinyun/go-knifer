package str

// DefaultIfNil returns def when v is nil; otherwise it returns *v.
func DefaultIfNil[T any](v *T, def T) T {
	if v == nil {
		return def
	}
	return *v
}

// DefaultIfEmpty returns def when s is empty.
func DefaultIfEmpty(s, def string) string {
	if IsEmpty(s) {
		return def
	}
	return s
}

// DefaultIfBlank returns def when s is blank.
func DefaultIfBlank(s, def string) string {
	if IsBlank(s) {
		return def
	}
	return s
}
