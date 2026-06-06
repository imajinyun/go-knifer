package system

import (
	"net"
	"strings"
)

type hostInfoConfig struct {
	hostname       func() (string, error)
	interfaceAddrs func() ([]net.Addr, error)
	address        func() string
}

// HostInfoOption customizes host information collection per call.
type HostInfoOption func(*hostInfoConfig)

// WithHostNameFunc sets the function used to collect the host name.
func WithHostNameFunc(fn func() (string, error)) HostInfoOption {
	return func(c *hostInfoConfig) {
		if fn != nil {
			c.hostname = fn
		}
	}
}

// WithHostInterfaceAddrsFunc sets the function used to collect local interface addresses.
func WithHostInterfaceAddrsFunc(fn func() ([]net.Addr, error)) HostInfoOption {
	return func(c *hostInfoConfig) {
		if fn != nil {
			c.interfaceAddrs = fn
		}
	}
}

// WithHostAddressFunc sets the function used to collect the host address directly.
func WithHostAddressFunc(fn func() string) HostInfoOption {
	return func(c *hostInfoConfig) {
		if fn != nil {
			c.address = fn
		}
	}
}

func applyHostInfoOptions(opts []HostInfoOption) hostInfoConfig {
	cfg := hostInfoConfig{
		hostname:       osHostname,
		interfaceAddrs: net.InterfaceAddrs,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	if cfg.hostname == nil {
		cfg.hostname = osHostname
	}
	if cfg.interfaceAddrs == nil {
		cfg.interfaceAddrs = net.InterfaceAddrs
	}
	return cfg
}

// HostInfo describes current host information.
type HostInfo struct {
	Name    string
	Address string
}

// NewHostInfo collects current host information.
func NewHostInfo() *HostInfo {
	return NewHostInfoWithOptions()
}

// NewHostInfoWithOptions collects host information using custom providers.
func NewHostInfoWithOptions(opts ...HostInfoOption) *HostInfo {
	cfg := applyHostInfoOptions(opts)
	h := &HostInfo{}
	if hostname, err := cfg.hostname(); err == nil {
		h.Name = hostname
	}
	if cfg.address != nil {
		h.Address = cfg.address()
	} else {
		h.Address = firstNonLoopbackIPv4(cfg.interfaceAddrs)
	}
	return h
}

// GetName returns the host name.
func (h *HostInfo) GetName() string { return h.Name }

// GetAddress returns the host IP address.
func (h *HostInfo) GetAddress() string { return h.Address }

// String implements fmt.Stringer.
func (h *HostInfo) String() string {
	var b strings.Builder
	appendLine(&b, "Host Name:    ", h.Name)
	appendLine(&b, "Host Address: ", h.Address)
	return b.String()
}

// firstNonLoopbackIPv4 returns the first non-loopback IPv4 address.
func firstNonLoopbackIPv4(interfaceAddrs func() ([]net.Addr, error)) string {
	addrs, err := interfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}
		ip4 := ipnet.IP.To4()
		if ip4 != nil {
			return ip4.String()
		}
	}
	return ""
}
