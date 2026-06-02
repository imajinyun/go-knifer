package system

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// UserInfo describes current logged-in user information.
type UserInfo struct {
	Name       string
	HomeDir    string
	CurrentDir string
	TempDir    string
	Language   string
	Country    string
}

// NewUserInfo creates current user information.
func NewUserInfo() *UserInfo {
	u := &UserInfo{}

	if cur, err := user.Current(); err == nil && cur != nil {
		u.Name = cur.Username
		u.HomeDir = fixPath(cur.HomeDir)
	} else {
		u.Name = os.Getenv("USER")
		if u.Name == "" {
			u.Name = os.Getenv("USERNAME")
		}
		u.HomeDir = fixPath(os.Getenv("HOME"))
	}

	if dir, err := os.Getwd(); err == nil {
		u.CurrentDir = fixPath(dir)
	}
	u.TempDir = fixPath(os.TempDir())

	lang, country := parseLocale(os.Getenv("LANG"))
	if lang == "" {
		lang, country = parseLocale(os.Getenv("LC_ALL"))
	}
	u.Language = lang
	u.Country = country
	return u
}

// GetName returns the user name.
func (u *UserInfo) GetName() string { return u.Name }

// GetHomeDir returns the home directory.
func (u *UserInfo) GetHomeDir() string { return u.HomeDir }

// GetCurrentDir returns the current working directory.
func (u *UserInfo) GetCurrentDir() string { return u.CurrentDir }

// GetTempDir returns the temporary directory.
func (u *UserInfo) GetTempDir() string { return u.TempDir }

// GetLanguage returns the language, such as zh.
func (u *UserInfo) GetLanguage() string { return u.Language }

// GetCountry returns the country or region, such as CN.
func (u *UserInfo) GetCountry() string { return u.Country }

// String implements fmt.Stringer.
func (u *UserInfo) String() string {
	var b strings.Builder
	appendLine(&b, "User Name:        ", u.Name)
	appendLine(&b, "User Home Dir:    ", u.HomeDir)
	appendLine(&b, "User Current Dir: ", u.CurrentDir)
	appendLine(&b, "User Temp Dir:    ", u.TempDir)
	appendLine(&b, "User Language:    ", u.Language)
	appendLine(&b, "User Country:     ", u.Country)
	return b.String()
}

// fixPath appends a trailing path separator.
func fixPath(p string) string {
	if p == "" {
		return p
	}
	return addSuffixIfNot(p, string(filepath.Separator))
}

// parseLocale parses a LANG string such as "zh_CN.UTF-8" and returns language and country.
func parseLocale(locale string) (lang, country string) {
	if locale == "" {
		return "", ""
	}
	if i := strings.IndexByte(locale, '.'); i >= 0 {
		locale = locale[:i]
	}
	parts := strings.Split(locale, "_")
	switch len(parts) {
	case 0:
		return "", ""
	case 1:
		return parts[0], ""
	default:
		return parts[0], parts[1]
	}
}
