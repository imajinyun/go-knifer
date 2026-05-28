package system

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// UserInfo 对应 hutool UserInfo，代表当前登录用户信息。
type UserInfo struct {
	Name       string
	HomeDir    string
	CurrentDir string
	TempDir    string
	Language   string
	Country    string
}

// NewUserInfo 构造当前用户信息。
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

// GetName 取得用户名。
func (u *UserInfo) GetName() string { return u.Name }

// GetHomeDir 取得 home 目录。
func (u *UserInfo) GetHomeDir() string { return u.HomeDir }

// GetCurrentDir 取得当前工作目录。
func (u *UserInfo) GetCurrentDir() string { return u.CurrentDir }

// GetTempDir 取得临时目录。
func (u *UserInfo) GetTempDir() string { return u.TempDir }

// GetLanguage 取得语言（如 zh）。
func (u *UserInfo) GetLanguage() string { return u.Language }

// GetCountry 取得国家/区域（如 CN）。
func (u *UserInfo) GetCountry() string { return u.Country }

// String 实现 fmt.Stringer。
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

// fixPath 路径末尾追加分隔符。
func fixPath(p string) string {
	if p == "" {
		return p
	}
	return addSuffixIfNot(p, string(filepath.Separator))
}

// parseLocale 解析形如 "zh_CN.UTF-8" 的 LANG 字符串，返回语言与国家。
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
