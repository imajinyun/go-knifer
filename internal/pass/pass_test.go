package pass

import "testing"

func TestAnalyze(t *testing.T) {
	tests := []struct {
		name           string
		password       string
		wantStrength   Strength
		wantCommonWeak bool
		wantRepeated   bool
		wantSequential bool
		wantStrong     bool
		wantWeak       bool
		wantLower      bool
		wantUpper      bool
		wantDigit      bool
		wantSymbol     bool
		minScore       int
		maxScore       int
	}{
		{
			name:         "empty",
			password:     "",
			wantStrength: StrengthVeryWeak,
			wantWeak:     true,
			maxScore:     0,
		},
		{
			name:           "common weak",
			password:       "password",
			wantStrength:   StrengthVeryWeak,
			wantCommonWeak: true,
			wantWeak:       true,
			wantLower:      true,
			maxScore:       10,
		},
		{
			name:           "short digit only",
			password:       "12345",
			wantStrength:   StrengthVeryWeak,
			wantWeak:       true,
			wantDigit:      true,
			wantSequential: true,
			maxScore:       30,
		},
		{
			name:           "repeated characters",
			password:       "aaaBBB111",
			wantStrength:   StrengthWeak,
			wantRepeated:   true,
			wantSequential: false,
			wantWeak:       true,
			wantLower:      true,
			wantUpper:      true,
			wantDigit:      true,
			minScore:       25,
			maxScore:       49,
		},
		{
			name:           "sequential characters",
			password:       "abcDEF123!",
			wantStrength:   StrengthMedium,
			wantSequential: true,
			wantLower:      true,
			wantUpper:      true,
			wantDigit:      true,
			wantSymbol:     true,
			minScore:       50,
			maxScore:       69,
		},
		{
			name:         "strong mixed",
			password:     "G0-Knifer#Pass2026",
			wantStrength: StrengthVeryStrong,
			wantStrong:   true,
			wantLower:    true,
			wantUpper:    true,
			wantDigit:    true,
			wantSymbol:   true,
			minScore:     85,
			maxScore:     100,
		},
		{
			name:         "unicode strong",
			password:     "安全Pass-2026!",
			wantStrength: StrengthVeryStrong,
			wantStrong:   true,
			wantLower:    true,
			wantUpper:    true,
			wantDigit:    true,
			wantSymbol:   true,
			minScore:     70,
			maxScore:     100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Analyze(tt.password)
			if got.Strength != tt.wantStrength {
				t.Fatalf("Strength = %v, want %v (analysis=%+v)", got.Strength, tt.wantStrength, got)
			}
			if got.CommonWeak != tt.wantCommonWeak {
				t.Errorf("CommonWeak = %v, want %v", got.CommonWeak, tt.wantCommonWeak)
			}
			if got.Repeated != tt.wantRepeated {
				t.Errorf("Repeated = %v, want %v", got.Repeated, tt.wantRepeated)
			}
			if got.Sequential != tt.wantSequential {
				t.Errorf("Sequential = %v, want %v", got.Sequential, tt.wantSequential)
			}
			if got.HasLower != tt.wantLower {
				t.Errorf("HasLower = %v, want %v", got.HasLower, tt.wantLower)
			}
			if got.HasUpper != tt.wantUpper {
				t.Errorf("HasUpper = %v, want %v", got.HasUpper, tt.wantUpper)
			}
			if got.HasDigit != tt.wantDigit {
				t.Errorf("HasDigit = %v, want %v", got.HasDigit, tt.wantDigit)
			}
			if got.HasSymbol != tt.wantSymbol {
				t.Errorf("HasSymbol = %v, want %v", got.HasSymbol, tt.wantSymbol)
			}
			if got.Score < tt.minScore || got.Score > tt.maxScore {
				t.Errorf("Score = %d, want in [%d,%d]", got.Score, tt.minScore, tt.maxScore)
			}
			if IsStrong(tt.password) != tt.wantStrong {
				t.Errorf("IsStrong() = %v, want %v", IsStrong(tt.password), tt.wantStrong)
			}
			if IsWeak(tt.password) != tt.wantWeak {
				t.Errorf("IsWeak() = %v, want %v", IsWeak(tt.password), tt.wantWeak)
			}
			if Score(tt.password) != got.Score {
				t.Errorf("Score() = %d, want %d", Score(tt.password), got.Score)
			}
			if StrengthOf(tt.password) != got.Strength {
				t.Errorf("StrengthOf() = %v, want %v", StrengthOf(tt.password), got.Strength)
			}
		})
	}
}

func TestStrengthString(t *testing.T) {
	tests := []struct {
		name     string
		strength Strength
		want     string
	}{
		{name: "unknown", strength: StrengthUnknown, want: "unknown"},
		{name: "very weak", strength: StrengthVeryWeak, want: "very weak"},
		{name: "weak", strength: StrengthWeak, want: "weak"},
		{name: "medium", strength: StrengthMedium, want: "medium"},
		{name: "strong", strength: StrengthStrong, want: "strong"},
		{name: "very strong", strength: StrengthVeryStrong, want: "very strong"},
		{name: "invalid", strength: Strength(99), want: "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.strength.String(); got != tt.want {
				t.Fatalf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}
