package net

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// URLBuilder builds URLs from scheme, host, path, query, and fragment parts.
type URLBuilder struct {
	scheme   string
	host     string
	port     int
	path     []string
	query    url.Values
	fragment string
	endSlash bool
}

// NewURLBuilder creates an empty URL builder.
func NewURLBuilder() *URLBuilder { return &URLBuilder{port: -1, query: url.Values{}} }

// NewHTTPURLBuilder creates an HTTP URL builder.
func NewHTTPURLBuilder(host string) *URLBuilder {
	return NewURLBuilder().SetScheme("http").SetHost(host)
}

// ParseURLBuilder parses raw into a URL builder.
func ParseURLBuilder(raw string) (*URLBuilder, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	b := NewURLBuilder().SetScheme(u.Scheme).SetHost(u.Hostname()).SetFragment(u.Fragment)
	if p := u.Port(); p != "" {
		if port, err := strconv.Atoi(p); err == nil {
			b.SetPort(port)
		}
	}
	b.SetPath(u.EscapedPath())
	b.query = u.Query()
	return b, nil
}

// SetScheme sets the URL scheme.
func (b *URLBuilder) SetScheme(scheme string) *URLBuilder { b.scheme = scheme; return b }

// Scheme returns the URL scheme.
func (b *URLBuilder) Scheme() string { return b.scheme }

// SchemeWithDefault returns the scheme or defaultScheme when empty.
func (b *URLBuilder) SchemeWithDefault(defaultScheme string) string {
	if b.scheme == "" {
		return defaultScheme
	}
	return b.scheme
}

// SetHost sets the URL host.
func (b *URLBuilder) SetHost(host string) *URLBuilder { b.host = host; return b }

// Host returns the URL host.
func (b *URLBuilder) Host() string { return b.host }

// SetPort sets the URL port.
func (b *URLBuilder) SetPort(port int) *URLBuilder { b.port = port; return b }

// Port returns the URL port.
func (b *URLBuilder) Port() int { return b.port }

// PortWithDefault returns the port or defaultPort when unset.
func (b *URLBuilder) PortWithDefault(defaultPort int) int {
	if b.port < 0 {
		return defaultPort
	}
	return b.port
}

// Authority returns host[:port].
func (b *URLBuilder) Authority() string {
	if b.port >= 0 {
		return b.host + ":" + strconv.Itoa(b.port)
	}
	return b.host
}

// SetWithEndTag controls whether the built path ends with a slash.
func (b *URLBuilder) SetWithEndTag(withEndTag bool) *URLBuilder { b.endSlash = withEndTag; return b }

// SetPath replaces the URL path.
func (b *URLBuilder) SetPath(path string) *URLBuilder {
	b.path = splitPath(path)
	return b
}

// AddPath appends path segments.
func (b *URLBuilder) AddPath(path string) *URLBuilder {
	b.path = append(b.path, splitPath(path)...)
	return b
}

// AddPathSegment appends one path segment.
func (b *URLBuilder) AddPathSegment(segment string) *URLBuilder {
	b.path = append(b.path, segment)
	return b
}

// PathString returns the encoded path string.
func (b *URLBuilder) PathString() string { return buildPath(b.path, b.endSlash) }

// SetQuery replaces the query from raw query text.
func (b *URLBuilder) SetQuery(rawQuery string) *URLBuilder {
	values, err := url.ParseQuery(strings.TrimPrefix(rawQuery, "?"))
	if err == nil {
		b.query = values
	}
	return b
}

// AddQuery adds a query parameter value.
func (b *URLBuilder) AddQuery(key string, value any) *URLBuilder {
	if b.query == nil {
		b.query = url.Values{}
	}
	b.query.Add(key, valueString(value))
	return b
}

// Query returns query values.
func (b *URLBuilder) Query() url.Values { return b.query }

// QueryString returns the encoded query string.
func (b *URLBuilder) QueryString() string { return b.query.Encode() }

// SetFragment sets the fragment text.
func (b *URLBuilder) SetFragment(fragment string) *URLBuilder { b.fragment = fragment; return b }

// Fragment returns the raw fragment text.
func (b *URLBuilder) Fragment() string { return b.fragment }

// FragmentEncoded returns the encoded fragment text.
func (b *URLBuilder) FragmentEncoded() string { return EncodeFragment(b.fragment) }

// Build returns the URL string.
func (b *URLBuilder) Build() string {
	u := &url.URL{Scheme: b.scheme, Host: b.Authority(), Path: "/" + strings.Join(b.path, "/"), RawQuery: b.query.Encode(), Fragment: b.fragment}
	if len(b.path) == 0 {
		u.Path = ""
	}
	if b.endSlash && !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	}
	return u.String()
}

func (b *URLBuilder) String() string { return b.Build() }

func splitPath(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return nil
	}
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if decoded, err := url.PathUnescape(part); err == nil {
			parts[i] = decoded
		}
	}
	return parts
}

func buildPath(parts []string, endSlash bool) string {
	if len(parts) == 0 {
		if endSlash {
			return "/"
		}
		return ""
	}
	escaped := make([]string, len(parts))
	for i, part := range parts {
		escaped[i] = url.PathEscape(part)
	}
	out := "/" + strings.Join(escaped, "/")
	if endSlash && !strings.HasSuffix(out, "/") {
		out += "/"
	}
	return out
}

func valueString(v any) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}
