package str

import "cmp"

// DefaultIfEmpty returns def when s is empty.
func DefaultIfEmpty(s, def string) string {
	return cmp.Or(s, def)
}

// DefaultIfBlank returns def when s is blank.
func DefaultIfBlank(s, def string) string {
	if IsBlank(s) {
		return def
	}
	return s
}
