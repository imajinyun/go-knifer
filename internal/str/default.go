package str

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
