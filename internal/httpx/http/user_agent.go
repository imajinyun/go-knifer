package http

import "strings"

// UserAgent represents parsed User-Agent information, aligned with the utility toolkit-http UserAgent.
type UserAgent struct {
	IsMobile bool
	Browser  string
	Version  string
	OS       string
	Engine   string
	Platform string
	Raw      string
}

// ParseUserAgent parses a UA string with lightweight rules covering common browsers and systems.
func ParseUserAgent(ua string) *UserAgent {
	u := &UserAgent{Raw: ua}
	low := strings.ToLower(ua)

	// Mobile devices.
	if strings.Contains(low, "mobile") || strings.Contains(low, "android") || strings.Contains(low, "iphone") || strings.Contains(low, "ipad") {
		u.IsMobile = true
	}

	// Browsers.
	switch {
	case strings.Contains(low, "edg/"):
		u.Browser = "Edge"
		u.Version = sliceAfter(low, "edg/")
	case strings.Contains(low, "edge/"):
		u.Browser = "Edge"
		u.Version = sliceAfter(low, "edge/")
	case strings.Contains(low, "chrome/") && !strings.Contains(low, "chromium"):
		u.Browser = "Chrome"
		u.Version = sliceAfter(low, "chrome/")
	case strings.Contains(low, "firefox/"):
		u.Browser = "Firefox"
		u.Version = sliceAfter(low, "firefox/")
	case strings.Contains(low, "safari/"):
		u.Browser = "Safari"
		u.Version = sliceAfter(low, "version/")
	case strings.Contains(low, "msie") || strings.Contains(low, "trident"):
		u.Browser = "MSIE"
	default:
		u.Browser = "Unknown"
	}

	// Engines.
	switch {
	case strings.Contains(low, "webkit"):
		u.Engine = "WebKit"
	case strings.Contains(low, "gecko"):
		u.Engine = "Gecko"
	case strings.Contains(low, "trident"):
		u.Engine = "Trident"
	case strings.Contains(low, "presto"):
		u.Engine = "Presto"
	default:
		u.Engine = "Unknown"
	}

	// OS
	switch {
	case strings.Contains(low, "windows nt 10"):
		u.OS = "Windows 10/11"
	case strings.Contains(low, "windows nt 6.3"):
		u.OS = "Windows 8.1"
	case strings.Contains(low, "windows nt 6.1"):
		u.OS = "Windows 7"
	case strings.Contains(low, "windows"):
		u.OS = "Windows"
	case strings.Contains(low, "iphone"):
		u.OS = "iOS"
	case strings.Contains(low, "ipad"):
		u.OS = "iPadOS"
	case strings.Contains(low, "mac os x"):
		u.OS = "macOS"
	case strings.Contains(low, "android"):
		u.OS = "Android"
	case strings.Contains(low, "linux"):
		u.OS = "Linux"
	default:
		u.OS = "Unknown"
	}

	// Platforms.
	switch {
	case strings.Contains(low, "windows"):
		u.Platform = "Windows"
	case strings.Contains(low, "macintosh"):
		u.Platform = "Macintosh"
	case strings.Contains(low, "android"):
		u.Platform = "Android"
	case strings.Contains(low, "iphone"):
		u.Platform = "iPhone"
	case strings.Contains(low, "ipad"):
		u.Platform = "iPad"
	case strings.Contains(low, "linux"):
		u.Platform = "Linux"
	default:
		u.Platform = "Unknown"
	}
	return u
}

func sliceAfter(s, prefix string) string {
	idx := strings.Index(s, prefix)
	if idx < 0 {
		return ""
	}
	rest := s[idx+len(prefix):]
	end := len(rest)
	for i, c := range rest {
		if c == ' ' || c == ';' || c == ')' {
			end = i
			break
		}
	}
	return rest[:end]
}
