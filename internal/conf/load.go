package conf

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// DecryptFunc decrypts encrypted configuration values.
type DecryptFunc func(cipherText string) (string, error)

// LoadOptions controls file/remote loading behavior.
type LoadOptions struct {
	// AllowInclude enables include/import keys in loaded configs.
	AllowInclude bool
	// IncludeKeys are keys used to discover included files. Defaults to include/import.
	IncludeKeys []string
	// Decrypt decrypts ENC(...) values after loading and merging.
	Decrypt DecryptFunc
	// RemoteClient is used by LoadRemote. Defaults to http.DefaultClient.
	RemoteClient *http.Client
	// Timeout bounds remote loading when RemoteClient has no timeout.
	Timeout time.Duration
}

// LoadWithOptions reads and parses a configuration file with advanced options.
func LoadWithOptions(path string, opts LoadOptions) (*Conf, error) {
	return loadFile(path, opts, map[string]bool{})
}

// LoadFiles loads multiple configuration files and merges them in order.
func LoadFiles(paths ...string) (*Conf, error) { return LoadFilesWithOptions(LoadOptions{}, paths...) }

// LoadFilesWithOptions loads multiple configuration files and merges them in order.
func LoadFilesWithOptions(opts LoadOptions, paths ...string) (*Conf, error) {
	configs := make([]*Conf, 0, len(paths))
	for _, path := range paths {
		c, err := LoadWithOptions(path, opts)
		if err != nil {
			return nil, err
		}
		configs = append(configs, c)
	}
	return Merge(configs...), nil
}

// LoadRemote loads configuration from an HTTP(S) URL.
func LoadRemote(rawURL string) (*Conf, error) { return LoadRemoteWithOptions(rawURL, LoadOptions{}) }

// LoadRemoteWithOptions loads configuration from an HTTP(S) URL with options.
func LoadRemoteWithOptions(rawURL string, opts LoadOptions) (*Conf, error) {
	return loadRemote(rawURL, opts)
}

// Merge merges configurations in order. Later configurations override earlier ones.
func Merge(configs ...*Conf) *Conf {
	out := New()
	for _, c := range configs {
		out.Merge(c)
	}
	return out
}

// Merge merges other into s. Existing keys are overwritten by other.
func (s *Conf) Merge(other *Conf) *Conf {
	if s == nil {
		return Merge(other)
	}
	if other == nil || other.data == nil {
		return s
	}
	for group, m := range other.data {
		for key, value := range m {
			s.SetByGroup(group, key, value)
		}
	}
	return s
}

func loadFile(path string, opts LoadOptions, seen map[string]bool) (*Conf, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, wrapConfigIO("resolve config file "+path, err)
	}
	if seen[abs] {
		return nil, invalidInputf("circular config include: %s", path)
	}
	seen[abs] = true
	defer delete(seen, abs)

	b, err := os.ReadFile(path) // #nosec G304 G703 -- configuration loader intentionally reads caller-provided paths.
	if err != nil {
		return nil, wrapConfigIO("read config file "+path, err)
	}
	current, err := ParseByExt(path, b)
	if err != nil {
		return nil, err
	}
	if !opts.AllowInclude {
		return current.DecryptValues(opts.Decrypt)
	}

	includes := current.includePaths(includeKeys(opts))
	current.removeIncludeKeys(includeKeys(opts))
	if len(includes) == 0 {
		return current.DecryptValues(opts.Decrypt)
	}
	baseDir := filepath.Dir(path)
	merged := New()
	for _, include := range includes {
		include = strings.TrimSpace(include)
		if include == "" {
			continue
		}
		if !filepath.IsAbs(include) {
			include = filepath.Join(baseDir, include)
		}
		c, err := loadFile(include, opts, seen)
		if err != nil {
			return nil, err
		}
		merged.Merge(c)
	}
	merged.Merge(current)
	return merged.DecryptValues(opts.Decrypt)
}

func loadRemote(rawURL string, opts LoadOptions) (*Conf, error) {
	client := opts.RemoteClient
	if client == nil {
		client = http.DefaultClient
	}
	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, invalidInputf("invalid remote config url %s: %s", rawURL, err.Error())
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, wrapConfigIO("fetch remote config "+rawURL, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, invalidInputf("fetch remote config %s: unexpected status %d", rawURL, resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, wrapConfigIO("read remote config "+rawURL, err)
	}
	parsePath := rawURL
	if u, err := url.Parse(rawURL); err == nil && u.Path != "" {
		parsePath = u.Path
	}
	c, err := ParseByExt(parsePath, b)
	if err != nil {
		return nil, err
	}
	return c.DecryptValues(opts.Decrypt)
}

func includeKeys(opts LoadOptions) []string {
	if len(opts.IncludeKeys) > 0 {
		return opts.IncludeKeys
	}
	return []string{"include", "import"}
}

func (s *Conf) includePaths(keys []string) []string {
	var out []string
	for _, key := range keys {
		if value, ok := s.Lookup(defaultGroup, key); ok {
			out = append(out, splitList(value)...)
		}
	}
	return out
}

func (s *Conf) removeIncludeKeys(keys []string) {
	for _, key := range keys {
		s.Delete(key)
	}
}
