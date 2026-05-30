package regex

import "regexp"

// ReMatch reports whether s matches pattern. Invalid patterns return false.
func ReMatch(pattern, s string) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(s)
}

// ReFind returns the first match, or an empty string when there is no match or the pattern is invalid.
func ReFind(pattern, s string) string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return ""
	}
	return re.FindString(s)
}

// ReFindAll returns all matches, or nil when the pattern is invalid.
func ReFindAll(pattern, s string) []string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil
	}
	return re.FindAllString(s, -1)
}

// ReReplace replaces matches of pattern with replacement. Invalid patterns return the original string.
func ReReplace(pattern, s, replacement string) string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return s
	}
	return re.ReplaceAllString(s, replacement)
}
