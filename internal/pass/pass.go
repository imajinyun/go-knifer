package pass

import (
	"strings"
	"unicode"
)

const (
	minStrongScore = 70
	maxScore       = 100
)

// Strength classifies a password score into broad buckets.
type Strength int

const (
	StrengthUnknown Strength = iota
	StrengthVeryWeak
	StrengthWeak
	StrengthMedium
	StrengthStrong
	StrengthVeryStrong
)

var commonWeakPasswords = map[string]struct{}{
	"000000":    {},
	"111111":    {},
	"123123":    {},
	"123456":    {},
	"12345678":  {},
	"123456789": {},
	"admin":     {},
	"iloveyou":  {},
	"letmein":   {},
	"password":  {},
	"password1": {},
	"qwerty":    {},
	"qwerty123": {},
	"welcome":   {},
	"woaini":    {},
	"zxcvbn":    {},
}

// Analysis describes the rule-level result of a password strength check.
type Analysis struct {
	Score      int
	Strength   Strength
	Length     int
	HasLower   bool
	HasUpper   bool
	HasDigit   bool
	HasSymbol  bool
	Repeated   bool
	Sequential bool
	CommonWeak bool
}

// Analyze evaluates password strength using local deterministic rules.
func Analyze(password string) Analysis {
	a := classify(password)
	a.Score = score(a)
	a.Strength = strengthFromScore(a.Score)
	return a
}

// Score returns a password strength score in the range 0..100.
func Score(password string) int {
	return Analyze(password).Score
}

// StrengthOf returns the strength bucket for password.
func StrengthOf(password string) Strength {
	return Analyze(password).Strength
}

// IsStrong reports whether password reaches the strong threshold.
func IsStrong(password string) bool {
	return Score(password) >= minStrongScore
}

// IsWeak reports whether password is weak or very weak.
func IsWeak(password string) bool {
	strength := StrengthOf(password)
	return strength == StrengthVeryWeak || strength == StrengthWeak
}

// String returns a stable lowercase label for s.
func (s Strength) String() string {
	switch s {
	case StrengthVeryWeak:
		return "very weak"
	case StrengthWeak:
		return "weak"
	case StrengthMedium:
		return "medium"
	case StrengthStrong:
		return "strong"
	case StrengthVeryStrong:
		return "very strong"
	default:
		return "unknown"
	}
}

func classify(password string) Analysis {
	runes := []rune(password)
	a := Analysis{
		Length:     len(runes),
		Repeated:   hasRepeatedRunes(runes),
		Sequential: hasSequentialRunes(runes),
		CommonWeak: isCommonWeak(password),
	}

	for _, r := range runes {
		switch {
		case unicode.IsLower(r):
			a.HasLower = true
		case unicode.IsUpper(r):
			a.HasUpper = true
		case unicode.IsDigit(r):
			a.HasDigit = true
		default:
			a.HasSymbol = true
		}
	}

	return a
}

func score(a Analysis) int {
	if a.Length == 0 {
		return 0
	}

	score := lengthScore(a.Length)
	score += categoryCount(a) * 10
	if categoryCount(a) >= 3 {
		score += 10
	}
	if a.Length >= 16 {
		score += 5
	}
	if a.Repeated {
		score -= 30
	}
	if a.Sequential {
		score -= 15
	}
	if categoryCount(a) == 1 && score > 35 {
		score = 35
	}
	if a.Length < 8 && score > 30 {
		score = 30
	}
	if a.CommonWeak && score > 10 {
		score = 10
	}

	return clamp(score, 0, maxScore)
}

func lengthScore(length int) int {
	switch {
	case length < 8:
		return length * 4
	case length < 12:
		return 30
	case length < 16:
		return 45
	default:
		return 55
	}
}

func categoryCount(a Analysis) int {
	count := 0
	if a.HasLower {
		count++
	}
	if a.HasUpper {
		count++
	}
	if a.HasDigit {
		count++
	}
	if a.HasSymbol {
		count++
	}
	return count
}

func strengthFromScore(score int) Strength {
	switch {
	case score >= 85:
		return StrengthVeryStrong
	case score >= 70:
		return StrengthStrong
	case score >= 50:
		return StrengthMedium
	case score >= 25:
		return StrengthWeak
	default:
		return StrengthVeryWeak
	}
}

func hasRepeatedRunes(runes []rune) bool {
	if len(runes) < 3 {
		return false
	}
	runLength := 1
	for i := 1; i < len(runes); i++ {
		if runes[i] == runes[i-1] {
			runLength++
			if runLength >= 3 {
				return true
			}
			continue
		}
		runLength = 1
	}
	return false
}

func hasSequentialRunes(runes []rune) bool {
	if len(runes) < 3 {
		return false
	}
	normalized := make([]rune, len(runes))
	for i, r := range runes {
		normalized[i] = unicode.ToLower(r)
	}
	for i := 2; i < len(normalized); i++ {
		prev2 := normalized[i-2]
		prev1 := normalized[i-1]
		cur := normalized[i]
		if prev1 == prev2+1 && cur == prev1+1 {
			return true
		}
		if prev1 == prev2-1 && cur == prev1-1 {
			return true
		}
	}
	return false
}

func isCommonWeak(password string) bool {
	normalized := strings.ToLower(strings.TrimSpace(password))
	_, ok := commonWeakPasswords[normalized]
	return ok
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
