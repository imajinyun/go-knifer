package system

import (
	"net"
	"strings"
)

// HostInfo describes current host information.
type HostInfo struct {
	Name    string
	Address string
}

// NewHostInfo collects current host information.
func NewHostInfo() *HostInfo {
	h := &HostInfo{}
	if hostname, err := osHostname(); err == nil {
		h.Name = hostname
	}
	h.Address = firstNonLoopbackIPv4()
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
func firstNonLoopbackIPv4() string {
	addrs, err := net.InterfaceAddrs()
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
