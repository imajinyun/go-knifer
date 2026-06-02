package net

import (
	"fmt"
	stdnet "net"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/idna"
)

const (
	// PortRangeMin is the default minimum dynamic helper port.
	PortRangeMin = 1024
	// PortRangeMax is the maximum TCP/UDP port.
	PortRangeMax = 0xFFFF
)

// IsValidPort reports whether port is a valid TCP/UDP port number.
func IsValidPort(port int) bool { return port >= 0 && port <= PortRangeMax }

// IsUsableLocalPort reports whether port can be bound locally on TCP.
func IsUsableLocalPort(port int) bool {
	if !IsValidPort(port) || port == 0 {
		return false
	}
	ln, err := stdnet.Listen("tcp", stdnet.JoinHostPort("127.0.0.1", strconvPort(port)))
	if err != nil {
		return false
	}
	_ = ln.Close()
	return true
}

// GetUsableLocalPort returns an available local TCP port in the default range.
func GetUsableLocalPort() (int, error) { return GetUsableLocalPortInRange(PortRangeMin, PortRangeMax) }

// GetUsableLocalPortFrom returns an available local TCP port from minPort to max.
func GetUsableLocalPortFrom(minPort int) (int, error) {
	return GetUsableLocalPortInRange(minPort, PortRangeMax)
}

// GetUsableLocalPortInRange returns an available local TCP port in [minPort, maxPort].
func GetUsableLocalPortInRange(minPort, maxPort int) (int, error) {
	if minPort < 0 || maxPort > PortRangeMax || minPort > maxPort {
		return 0, fmt.Errorf("invalid port range: %d-%d", minPort, maxPort)
	}
	for port := minPort; port <= maxPort; port++ {
		if IsUsableLocalPort(port) {
			return port, nil
		}
	}
	return 0, fmt.Errorf("no usable local port in range %d-%d", minPort, maxPort)
}

// GetUsableLocalPorts returns up to numRequested available ports in [minPort, maxPort].
func GetUsableLocalPorts(numRequested, minPort, maxPort int) ([]int, error) {
	if numRequested <= 0 {
		return nil, nil
	}
	ports := make([]int, 0, numRequested)
	for port := minPort; port <= maxPort && len(ports) < numRequested; port++ {
		if IsUsableLocalPort(port) {
			ports = append(ports, port)
		}
	}
	if len(ports) < numRequested {
		return ports, fmt.Errorf("only found %d usable local ports", len(ports))
	}
	return ports, nil
}

// LocalPortGenerator generates available local ports from a moving cursor.
type LocalPortGenerator struct{ next int }

// NewLocalPortGenerator creates a local port generator.
func NewLocalPortGenerator(beginPort int) *LocalPortGenerator {
	return &LocalPortGenerator{next: beginPort}
}

// Generate returns the next available local port.
func (g *LocalPortGenerator) Generate() (int, error) {
	if g == nil {
		return 0, fmt.Errorf("nil local port generator")
	}
	port, err := GetUsableLocalPortInRange(g.next, PortRangeMax)
	if err != nil {
		return 0, err
	}
	g.next = port + 1
	return port, nil
}

// HideIPPart hides the last IPv4 segment.
func HideIPPart(ip string) string {
	idx := strings.LastIndex(ip, ".")
	if idx < 0 {
		return ip
	}
	return ip[:idx+1] + "*"
}

// HideIPPartLong hides the last segment of an IPv4 integer.
func HideIPPartLong(ip uint32) string { return HideIPPart(LongToIPv4(ip)) }

// BuildInetSocketAddress builds a TCP address with a default port when host contains none.
func BuildInetSocketAddress(host string, defaultPort int) (*stdnet.TCPAddr, error) {
	if _, _, err := stdnet.SplitHostPort(host); err == nil {
		return stdnet.ResolveTCPAddr("tcp", host)
	}
	return stdnet.ResolveTCPAddr("tcp", stdnet.JoinHostPort(host, strconvPort(defaultPort)))
}

// CreateAddress builds a TCP address from host and port.
func CreateAddress(host string, port int) *stdnet.TCPAddr {
	addr, _ := stdnet.ResolveTCPAddr("tcp", stdnet.JoinHostPort(host, strconvPort(port)))
	return addr
}

// GetIPByHost resolves hostName to the first IP string.
func GetIPByHost(hostName string) string {
	ips, err := stdnet.LookupIP(hostName)
	if err != nil || len(ips) == 0 {
		return hostName
	}
	return ips[0].String()
}

// GetNetworkInterface returns a network interface by name.
func GetNetworkInterface(name string) (*stdnet.Interface, error) { return stdnet.InterfaceByName(name) }

// GetNetworkInterfaces returns all network interfaces.
func GetNetworkInterfaces() ([]stdnet.Interface, error) { return stdnet.Interfaces() }

// LocalIPv4s returns local IPv4 addresses.
func LocalIPv4s() []string { return ToIPList(LocalAddressList(nil)) }

// LocalIPv6s returns local IPv6 addresses.
func LocalIPv6s() []string {
	return ToIPList(LocalAddressList(func(ip stdnet.IP) bool { return ip.To4() == nil && ip.To16() != nil }))
}

// LocalIPs returns all local IP addresses.
func LocalIPs() []string {
	return ToIPList(LocalAddressList(func(ip stdnet.IP) bool { return ip != nil }))
}

// ToIPList converts IP addresses to strings.
func ToIPList(addressList []stdnet.IP) []string {
	out := make([]string, 0, len(addressList))
	seen := map[string]struct{}{}
	for _, ip := range addressList {
		if ip == nil {
			continue
		}
		s := ip.String()
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			out = append(out, s)
		}
	}
	return out
}

// LocalAddressList returns local IP addresses matching addressFilter. nil means non-loopback IPv4.
func LocalAddressList(addressFilter func(stdnet.IP) bool) []stdnet.IP {
	return LocalAddressListByInterface(nil, addressFilter)
}

// LocalAddressListByInterface returns local IP addresses matching interface and address filters.
func LocalAddressListByInterface(interfaceFilter func(stdnet.Interface) bool, addressFilter func(stdnet.IP) bool) []stdnet.IP {
	interfaces, err := stdnet.Interfaces()
	if err != nil {
		return nil
	}
	out := make([]stdnet.IP, 0)
	for _, iface := range interfaces {
		if interfaceFilter != nil && !interfaceFilter(iface) {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ip := addrIP(addr)
			if ip == nil {
				continue
			}
			if addressFilter == nil {
				if ip.IsLoopback() || ip.To4() == nil {
					continue
				}
			} else if !addressFilter(ip) {
				continue
			}
			out = append(out, ip)
		}
	}
	return out
}

// GetLocalhostStr returns a preferred local host IP string.
func GetLocalhostStr() string {
	ips := LocalAddressList(nil)
	if len(ips) > 0 {
		return ips[0].String()
	}
	return LocalIP
}

// GetLocalhost returns a preferred local host IP.
func GetLocalhost() stdnet.IP { return stdnet.ParseIP(GetLocalhostStr()) }

// GetLocalHostName returns the OS host name.
func GetLocalHostName() string {
	host, err := stdnet.LookupAddr(GetLocalhostStr())
	if err == nil && len(host) > 0 {
		return strings.TrimSuffix(host[0], ".")
	}
	name, _ := osHostname()
	return name
}

// GetLocalMACAddress returns the first non-empty local hardware address.
func GetLocalMACAddress(separator ...string) string {
	hw := GetLocalHardwareAddress()
	if hw == nil {
		return ""
	}
	sep := ":"
	if len(separator) > 0 {
		sep = separator[0]
	}
	return formatHardwareAddress(hw, sep)
}

// GetMACAddress returns the hardware address of the interface owning inetAddress.
func GetMACAddress(inetAddress stdnet.IP, separator ...string) string {
	hw := GetHardwareAddress(inetAddress)
	if hw == nil {
		return ""
	}
	sep := ":"
	if len(separator) > 0 {
		sep = separator[0]
	}
	return formatHardwareAddress(hw, sep)
}

// GetHardwareAddress returns the hardware address of the interface owning inetAddress.
func GetHardwareAddress(inetAddress stdnet.IP) stdnet.HardwareAddr {
	interfaces, err := stdnet.Interfaces()
	if err != nil {
		return nil
	}
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if ip := addrIP(addr); ip != nil && ip.Equal(inetAddress) {
				return iface.HardwareAddr
			}
		}
	}
	return nil
}

// GetLocalHardwareAddress returns the first non-empty local hardware address.
func GetLocalHardwareAddress() stdnet.HardwareAddr {
	interfaces, err := stdnet.Interfaces()
	if err != nil {
		return nil
	}
	for _, iface := range interfaces {
		if len(iface.HardwareAddr) > 0 && iface.Flags&stdnet.FlagLoopback == 0 {
			return iface.HardwareAddr
		}
	}
	return nil
}

// NetCat sends data to host:port over TCP.
func NetCat(host string, port int, data []byte, timeout time.Duration) error {
	conn, err := stdnet.DialTimeout("tcp", stdnet.JoinHostPort(host, strconvPort(port)), timeout)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()
	_, err = conn.Write(data)
	return err
}

// Ping checks whether an IP or host is reachable by opening a TCP connection to common ports.
func Ping(ip string, timeout time.Duration) bool {
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	for _, port := range []string{"80", "443"} {
		conn, err := stdnet.DialTimeout("tcp", stdnet.JoinHostPort(ip, port), timeout)
		if err == nil {
			_ = conn.Close()
			return true
		}
	}
	return false
}

// IsOpen reports whether address can be opened within timeout.
func IsOpen(address *stdnet.TCPAddr, timeout time.Duration) bool {
	if address == nil {
		return false
	}
	conn, err := stdnet.DialTimeout("tcp", address.String(), timeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// IDNToASCII converts a Unicode domain name to ASCII.
func IDNToASCII(unicode string) (string, error) { return idna.ToASCII(unicode) }

// GetMultistageReverseProxyIP returns the first non-unknown IP in a comma-separated proxy header.
func GetMultistageReverseProxyIP(ip string) string {
	for _, part := range strings.Split(ip, ",") {
		part = strings.TrimSpace(part)
		if !IsUnknown(part) {
			return part
		}
	}
	return ""
}

// IsUnknown reports whether checkString is empty or equals unknown case-insensitively.
func IsUnknown(checkString string) bool {
	return strings.TrimSpace(checkString) == "" || strings.EqualFold(strings.TrimSpace(checkString), "unknown")
}

// ParseCookies parses a Cookie header value.
func ParseCookies(cookieStr string) []*http.Cookie {
	req := &http.Request{Header: http.Header{"Cookie": []string{cookieStr}}}
	return req.Cookies()
}

// GetDNSInfo looks up DNS records by attribute names such as A, CNAME, MX, NS, or TXT.
func GetDNSInfo(hostName string, attrNames ...string) ([]string, error) {
	if len(attrNames) == 0 {
		attrNames = []string{"A"}
	}
	out := make([]string, 0)
	for _, attr := range attrNames {
		switch strings.ToUpper(attr) {
		case "A", "AAAA":
			ips, err := stdnet.LookupIP(hostName)
			if err != nil {
				return out, err
			}
			for _, ip := range ips {
				if strings.ToUpper(attr) == "A" && ip.To4() == nil {
					continue
				}
				if strings.ToUpper(attr) == "AAAA" && ip.To4() != nil {
					continue
				}
				out = append(out, ip.String())
			}
		case "CNAME":
			v, err := stdnet.LookupCNAME(hostName)
			if err != nil {
				return out, err
			}
			out = append(out, v)
		case "MX":
			mxs, err := stdnet.LookupMX(hostName)
			if err != nil {
				return out, err
			}
			for _, mx := range mxs {
				out = append(out, mx.Host)
			}
		case "NS":
			nss, err := stdnet.LookupNS(hostName)
			if err != nil {
				return out, err
			}
			for _, ns := range nss {
				out = append(out, ns.Host)
			}
		case "TXT":
			txts, err := stdnet.LookupTXT(hostName)
			if err != nil {
				return out, err
			}
			out = append(out, txts...)
		}
	}
	return out, nil
}

// Connect opens a TCP connection to host:port.
func Connect(hostname string, port int, timeout time.Duration) (stdnet.Conn, error) {
	addr := stdnet.JoinHostPort(hostname, strconvPort(port))
	if timeout > 0 {
		return stdnet.DialTimeout("tcp", addr, timeout)
	}
	return stdnet.Dial("tcp", addr)
}

// GetRemoteAddress returns conn's remote address string.
func GetRemoteAddress(conn stdnet.Conn) string {
	if conn == nil || conn.RemoteAddr() == nil {
		return ""
	}
	return conn.RemoteAddr().String()
}

// IsConnected reports whether conn appears open.
func IsConnected(conn stdnet.Conn) bool { return conn != nil && conn.RemoteAddr() != nil }

func addrIP(addr stdnet.Addr) stdnet.IP {
	switch v := addr.(type) {
	case *stdnet.IPNet:
		return v.IP
	case *stdnet.IPAddr:
		return v.IP
	default:
		return nil
	}
}

func formatHardwareAddress(hw stdnet.HardwareAddr, sep string) string {
	parts := make([]string, len(hw))
	for i, b := range hw {
		parts[i] = fmt.Sprintf("%02x", b)
	}
	return strings.Join(parts, sep)
}

func strconvPort(port int) string { return fmt.Sprintf("%d", port) }

func osHostname() (string, error) { return os.Hostname() }
