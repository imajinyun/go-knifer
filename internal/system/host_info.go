package system

import (
	"net"
	"strings"
)

// HostInfo 对应 hutool HostInfo，代表当前主机的信息。
type HostInfo struct {
	Name    string
	Address string
}

// NewHostInfo 收集当前主机信息。
func NewHostInfo() *HostInfo {
	h := &HostInfo{}
	if hostname, err := osHostname(); err == nil {
		h.Name = hostname
	}
	h.Address = firstNonLoopbackIPv4()
	return h
}

// GetName 取得主机名。
func (h *HostInfo) GetName() string { return h.Name }

// GetAddress 取得主机 IP 地址。
func (h *HostInfo) GetAddress() string { return h.Address }

// String 实现 fmt.Stringer。
func (h *HostInfo) String() string {
	var b strings.Builder
	appendLine(&b, "Host Name:    ", h.Name)
	appendLine(&b, "Host Address: ", h.Address)
	return b.String()
}

// firstNonLoopbackIPv4 返回第一个非回环的 IPv4 地址。
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
