package boolutil

// BoolNegate returns the logical negation of b.
func BoolNegate(b bool) bool { return !b }

// BoolToInt returns 1 for true and 0 for false.
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// BoolAnd returns true only when all inputs are true.
func BoolAnd(bs ...bool) bool {
	for _, b := range bs {
		if !b {
			return false
		}
	}
	return true
}

// BoolOr returns true when any input is true.
func BoolOr(bs ...bool) bool {
	for _, b := range bs {
		if b {
			return true
		}
	}
	return false
}
