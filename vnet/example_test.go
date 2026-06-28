package vnet_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/imajinyun/knifer-go/vnet"
)

const exampleCertPEM = `-----BEGIN CERTIFICATE-----
MIIBhTCCASugAwIBAgIRAPWQSq0Qr7yZD5twH61BxFIwCgYIKoZIzj0EAwIwEjEQ
MA4GA1UEChMHZ28tdGVzdDAeFw0yNjA2MDYwMDAwMDBaFw0yNzA2MDYwMDAwMDBa
MBIxEDAOBgNVBAoTB2dvLXRlc3QwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASm
1YPqMC7UTw4R7ovbHYgk4+LALoU6hr61VnsBiKCdsMCMScpLob8ldIl+6o4f/ntM
5kmXvEFd9Mp6FfaHkgnbo0IwQDAOBgNVHQ8BAf8EBAMCAqQwDwYDVR0TAQH/BAUw
AwEB/zAdBgNVHQ4EFgQUX90U1OkOXbGUzD2JNoWlqQtk3/0wCgYIKoZIzj0EAwID
SQAwRgIhANw7UzN0vtxOfygWqANg00uGOo7y98q1/Ac3N1wQxVBkAiEA7QjQRHtH
LA6wKo8yoCnW36b+nvxlhHvzrIxwWCgwCWM=
-----END CERTIFICATE-----`

type exampleDialer struct {
	network string
	address string
	data    chan []byte
}

func (d *exampleDialer) DialContext(_ context.Context, network, address string) (net.Conn, error) {
	d.network = network
	d.address = address
	client, server := net.Pipe()
	if d.data != nil {
		go func() {
			defer func() { _ = server.Close() }()
			payload, _ := io.ReadAll(server)
			d.data <- payload
		}()
	} else {
		_ = server.Close()
	}
	return client, nil
}

type exampleListener struct{}

func (exampleListener) Accept() (net.Conn, error) { return nil, io.EOF }
func (exampleListener) Close() error              { return nil }
func (exampleListener) Addr() net.Addr            { return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 1000} }

func ExampleCreateAddress() {
	addr := vnet.CreateAddress("127.0.0.1", 8080)
	fmt.Println(addr.String())
	// Output: 127.0.0.1:8080
}

func ExampleIPv4ToLong() {
	n, err := vnet.IPv4ToLong("192.0.2.1")

	fmt.Println(n)
	fmt.Println(vnet.LongToIPv4(n))
	fmt.Println(err)
	// Output:
	// 3221225985
	// 192.0.2.1
	// <nil>
}

func ExampleIPv4ToLongWithOptions() {
	parseIP := vnet.WithIPParser(func(string) net.IP {
		return net.ParseIP("198.51.100.9")
	})
	n, err := vnet.IPv4ToLongWithOptions("ignored", parseIP)
	fmt.Println(vnet.LongToIPv4(n))
	fmt.Println(err)
	// Output:
	// 198.51.100.9
	// <nil>
}

func ExampleLongToIPv4() {
	fmt.Println(vnet.LongToIPv4(3221225985))
	// Output: 192.0.2.1
}

func ExampleIPv4ToLongDefault() {
	fmt.Println(vnet.IPv4ToLongDefault("not-an-ip", 42))
	// Output: 42
}

func ExampleIPv4ToLongDefaultWithOptions() {
	parseIP := vnet.WithIPParser(func(string) net.IP {
		return nil
	})
	fmt.Println(vnet.IPv4ToLongDefaultWithOptions("ignored", 42, parseIP))
	// Output: 42
}

func ExampleIPv6ToBigInt() {
	n, err := vnet.IPv6ToBigInt("2001:db8::1")
	fmt.Println(n.String())
	fmt.Println(err)
	// Output:
	// 42540766411282592856903984951653826561
	// <nil>
}

func ExampleIPv6ToBigIntWithOptions() {
	parseIP := vnet.WithIPParser(func(string) net.IP {
		return net.ParseIP("2001:db8::2")
	})
	n, err := vnet.IPv6ToBigIntWithOptions("ignored", parseIP)
	fmt.Println(n.String())
	fmt.Println(err)
	// Output:
	// 42540766411282592856903984951653826562
	// <nil>
}

func ExampleBigIntToIPv6() {
	n, _ := vnet.IPv6ToBigInt("2001:db8::1")
	ip, err := vnet.BigIntToIPv6(n)
	fmt.Println(ip)
	fmt.Println(err)
	// Output:
	// 2001:db8::1
	// <nil>
}

func ExampleIsIP() {
	fmt.Println(vnet.IsIP("192.0.2.1"))
	fmt.Println(vnet.IsIP("example"))
	// Output:
	// true
	// false
}

func ExampleIsIPWithOptions() {
	parseIP := vnet.WithIPParser(func(string) net.IP {
		return net.ParseIP("2001:db8::1")
	})
	fmt.Println(vnet.IsIPWithOptions("ignored", parseIP))
	// Output: true
}

func ExampleIsIPv4() {
	fmt.Println(vnet.IsIPv4("192.0.2.1"))
	fmt.Println(vnet.IsIPv4("2001:db8::1"))
	// Output:
	// true
	// false
}

func ExampleIsIPv4WithOptions() {
	parseIP := vnet.WithIPParser(func(string) net.IP {
		return net.ParseIP("192.0.2.1")
	})
	fmt.Println(vnet.IsIPv4WithOptions("ignored", parseIP))
	// Output: true
}

func ExampleIsIPv6() {
	fmt.Println(vnet.IsIPv6("2001:db8::1"))
	fmt.Println(vnet.IsIPv6("192.0.2.1"))
	// Output:
	// true
	// false
}

func ExampleIsIPv6WithOptions() {
	parseIP := vnet.WithIPParser(func(string) net.IP {
		return net.ParseIP("2001:db8::1")
	})
	fmt.Println(vnet.IsIPv6WithOptions("ignored", parseIP))
	// Output: true
}

func ExampleIsInnerIP() {
	fmt.Println(vnet.IsInnerIP("10.0.0.1"))
	fmt.Println(vnet.IsInnerIP("203.0.113.1"))
	// Output:
	// true
	// false
}

func ExampleIsInnerIPWithOptions() {
	parseIP := vnet.WithIPParser(func(string) net.IP {
		return net.ParseIP("172.16.0.1")
	})
	fmt.Println(vnet.IsInnerIPWithOptions("ignored", parseIP))
	// Output: true
}

func ExampleFormatIPBlock() {
	block, err := vnet.FormatIPBlock("192.0.2.9", "255.255.255.0")
	fmt.Println(block)
	fmt.Println(err)
	// Output:
	// 192.0.2.9/24
	// <nil>
}

func ExampleFormatIPBlockWithOptions() {
	parseIP := vnet.WithIPParser(net.ParseIP)
	maskBit, _ := vnet.FormatIPBlockWithOptions("192.0.2.9", "255.255.255.0", parseIP)
	fmt.Println(maskBit)
	// Output: 192.0.2.9/24
}

func ExampleBeginIP() {
	begin, beginErr := vnet.BeginIP("192.0.2.9", 24)
	end, endErr := vnet.EndIP("192.0.2.9", 24)
	count, countErr := vnet.CountByMaskBit(24, true)

	fmt.Println(begin, beginErr)
	fmt.Println(end, endErr)
	fmt.Println(count, countErr)
	// Output:
	// 192.0.2.0 <nil>
	// 192.0.2.255 <nil>
	// 256 <nil>
}

func ExampleBeginIPWithOptions() {
	parseIP := vnet.WithIPParser(func(string) net.IP {
		return net.ParseIP("192.0.2.9")
	})
	begin, err := vnet.BeginIPWithOptions("ignored", 24, parseIP)
	fmt.Println(begin)
	fmt.Println(err)
	// Output:
	// 192.0.2.0
	// <nil>
}

func ExampleBeginIPLong() {
	n, err := vnet.BeginIPLong("192.0.2.9", 24)
	fmt.Println(n)
	fmt.Println(vnet.LongToIPv4(n))
	fmt.Println(err)
	// Output:
	// 3221225984
	// 192.0.2.0
	// <nil>
}

func ExampleBeginIPLongWithOptions() {
	parseIP := vnet.WithIPParser(func(string) net.IP {
		return net.ParseIP("192.0.2.9")
	})
	n, err := vnet.BeginIPLongWithOptions("ignored", 24, parseIP)
	fmt.Println(vnet.LongToIPv4(n))
	fmt.Println(err)
	// Output:
	// 192.0.2.0
	// <nil>
}

func ExampleEndIP() {
	ip, err := vnet.EndIP("192.0.2.9", 24)
	fmt.Println(ip)
	fmt.Println(err)
	// Output:
	// 192.0.2.255
	// <nil>
}

func ExampleEndIPWithOptions() {
	parseIP := vnet.WithIPParser(func(string) net.IP {
		return net.ParseIP("192.0.2.9")
	})
	ip, err := vnet.EndIPWithOptions("ignored", 24, parseIP)
	fmt.Println(ip)
	fmt.Println(err)
	// Output:
	// 192.0.2.255
	// <nil>
}

func ExampleEndIPLong() {
	n, err := vnet.EndIPLong("192.0.2.9", 24)
	fmt.Println(n)
	fmt.Println(vnet.LongToIPv4(n))
	fmt.Println(err)
	// Output:
	// 3221226239
	// 192.0.2.255
	// <nil>
}

func ExampleEndIPLongWithOptions() {
	parseIP := vnet.WithIPParser(func(string) net.IP {
		return net.ParseIP("192.0.2.9")
	})
	n, err := vnet.EndIPLongWithOptions("ignored", 24, parseIP)
	fmt.Println(vnet.LongToIPv4(n))
	fmt.Println(err)
	// Output:
	// 192.0.2.255
	// <nil>
}

func ExampleMaskBitByMask() {
	maskBit, err := vnet.MaskBitByMask("255.255.255.0")
	fmt.Println(maskBit, err)
	// Output: 24 <nil>
}

func ExampleMaskBitByMaskWithOptions() {
	parseInt := vnet.WithIPIntParser(func(part string) (int, error) {
		return map[string]int{"255": 255, "0": 0}[part], nil
	})
	maskBit, err := vnet.MaskBitByMaskWithOptions("255.255.255.0", parseInt)
	fmt.Println(maskBit, err)
	// Output: 24 <nil>
}

func ExampleCountByMaskBit() {
	all, _ := vnet.CountByMaskBit(30, true)
	usable, _ := vnet.CountByMaskBit(30, false)
	fmt.Println(all)
	fmt.Println(usable)
	// Output:
	// 4
	// 2
}

func ExampleMaskByMaskBit() {
	mask, err := vnet.MaskByMaskBit(24)
	fmt.Println(mask)
	fmt.Println(err)
	// Output:
	// 255.255.255.0
	// <nil>
}

func ExampleMaskByIPRange() {
	mask, err := vnet.MaskByIPRange("192.0.2.0", "192.0.2.255")
	fmt.Println(mask)
	fmt.Println(err)
	// Output:
	// 255.255.255.0
	// <nil>
}

func ExampleMaskByIPRangeWithOptions() {
	mask, err := vnet.MaskByIPRangeWithOptions("192.0.2.0", "192.0.2.255")
	fmt.Println(mask)
	fmt.Println(err)
	// Output:
	// 255.255.255.0
	// <nil>
}

func ExampleCountByIPRange() {
	count, err := vnet.CountByIPRange("192.0.2.1", "192.0.2.4")
	fmt.Println(count)
	fmt.Println(err)
	// Output:
	// 4
	// <nil>
}

func ExampleCountByIPRangeWithOptions() {
	count, err := vnet.CountByIPRangeWithOptions("192.0.2.1", "192.0.2.4")
	fmt.Println(count)
	fmt.Println(err)
	// Output:
	// 4
	// <nil>
}

func ExampleIsMaskValid() {
	fmt.Println(vnet.IsMaskValid("255.255.255.0"))
	fmt.Println(vnet.IsMaskValid("255.0.255.0"))
	// Output:
	// true
	// false
}

func ExampleIsMaskValidWithOptions() {
	parseInt := vnet.WithIPIntParser(func(part string) (int, error) {
		return map[string]int{"255": 255, "0": 0}[part], nil
	})
	fmt.Println(vnet.IsMaskValidWithOptions("255.255.255.0", parseInt))
	// Output: true
}

func ExampleIsMaskBitValid() {
	fmt.Println(vnet.IsMaskBitValid(24))
	fmt.Println(vnet.IsMaskBitValid(40))
	// Output:
	// true
	// false
}

func ExampleListIPs() {
	ips, err := vnet.ListIPs("192.0.2.1-192.0.2.3", true)
	fmt.Println(ips)
	fmt.Println(err)
	// Output:
	// [192.0.2.1 192.0.2.2 192.0.2.3]
	// <nil>
}

func ExampleListIPsWithOptions() {
	ips, err := vnet.ListIPsWithOptions("192.0.2.1-192.0.2.2", true)
	fmt.Println(ips)
	fmt.Println(err)
	// Output:
	// [192.0.2.1 192.0.2.2]
	// <nil>
}

func ExampleListIPCIDR() {
	ips, err := vnet.ListIPCIDR("192.0.2.0", 30, false)
	fmt.Println(ips)
	fmt.Println(err)
	// Output:
	// [192.0.2.1 192.0.2.2]
	// <nil>
}

func ExampleListIPCIDRWithOptions() {
	ips, err := vnet.ListIPCIDRWithOptions("192.0.2.0", 30, false)
	fmt.Println(ips)
	fmt.Println(err)
	// Output:
	// [192.0.2.1 192.0.2.2]
	// <nil>
}

func ExampleListIPRange() {
	ips, err := vnet.ListIPRange("192.0.2.1", "192.0.2.2")
	fmt.Println(ips)
	fmt.Println(err)
	// Output:
	// [192.0.2.1 192.0.2.2]
	// <nil>
}

func ExampleListIPRangeWithOptions() {
	ips, err := vnet.ListIPRangeWithOptions("192.0.2.1", "192.0.2.2")
	fmt.Println(ips)
	fmt.Println(err)
	// Output:
	// [192.0.2.1 192.0.2.2]
	// <nil>
}

func ExampleParseCookies() {
	cookies := vnet.ParseCookies("sid=abc; theme=dark")
	for _, cookie := range cookies {
		fmt.Println(cookie.Name, cookie.Value)
	}
	// Output:
	// sid abc
	// theme dark
}

func ExampleMatchesWildcard() {
	fmt.Println(vnet.MatchesWildcard("192.168.*.*", "192.168.1.2"))
	fmt.Println(vnet.MatchesWildcard("10.0.*.*", "192.168.1.2"))
	// Output:
	// true
	// false
}

func ExampleMatchesWildcardWithOptions() {
	parseIP := vnet.WithWildcardIPParser(func(string) net.IP {
		return net.ParseIP("203.0.113.9")
	})

	fmt.Println(vnet.MatchesWildcardWithOptions("203.0.113.*", "ignored", parseIP))
	// Output: true
}

func ExampleIsInRange() {
	fmt.Println(vnet.IsInRange("192.0.2.10", "192.0.2.0/24"))
	fmt.Println(vnet.IsInRange("198.51.100.10", "192.0.2.0/24"))
	// Output:
	// true
	// false
}

func ExampleIsInRangeWithOptions() {
	parseIP := vnet.WithIPParser(func(string) net.IP {
		return net.ParseIP("192.0.2.10")
	})

	fmt.Println(vnet.IsInRangeWithOptions("ignored", "192.0.2.0/24", parseIP))
	// Output: true
}

func ExampleHideIPPart() {
	fmt.Println(vnet.HideIPPart("192.0.2.99"))
	// Output: 192.0.2.*
}

func ExampleHideIPPartLong() {
	n, _ := vnet.IPv4ToLong("192.0.2.99")
	fmt.Println(vnet.HideIPPartLong(n))
	// Output: 192.0.2.*
}

func ExampleIDNToASCII() {
	domain, err := vnet.IDNToASCII("例子.测试")
	fmt.Println(domain)
	fmt.Println(err)
	// Output:
	// xn--fsqu00a.xn--0zwm56d
	// <nil>
}

func ExampleGetMultistageReverseProxyIP() {
	fmt.Println(vnet.GetMultistageReverseProxyIP("unknown, 198.51.100.7, 203.0.113.9"))
	// Output: 198.51.100.7
}

func ExampleIsUnknown() {
	fmt.Println(vnet.IsUnknown("unknown"))
	fmt.Println(vnet.IsUnknown("198.51.100.7"))
	// Output:
	// true
	// false
}

func ExampleIsValidPort() {
	fmt.Println(vnet.IsValidPort(443))
	fmt.Println(vnet.IsValidPort(70000))
	// Output:
	// true
	// false
}

func ExampleToIPList() {
	ips := vnet.ToIPList([]net.IP{
		net.ParseIP("192.0.2.1"),
		net.ParseIP("192.0.2.1"),
		net.ParseIP("2001:db8::1"),
	})

	fmt.Println(ips)
	// Output: [192.0.2.1 2001:db8::1]
}

func ExampleBuildInetSocketAddressWithOptions() {
	resolver := vnet.WithTCPAddrResolver(func(network, address string) (*net.TCPAddr, error) {
		fmt.Println(network, address)
		return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080}, nil
	})

	addr, err := vnet.BuildInetSocketAddressWithOptions("example.test", 8080, resolver)
	fmt.Println(addr.String())
	fmt.Println(err)
	// Output:
	// tcp example.test:8080
	// 127.0.0.1:8080
	// <nil>
}

func ExampleBuildInetSocketAddress() {
	addr, err := vnet.BuildInetSocketAddress("127.0.0.1", 8080)
	fmt.Println(addr.String())
	fmt.Println(err)
	// Output:
	// 127.0.0.1:8080
	// <nil>
}

func ExampleCreateAddressWithOptions() {
	resolver := vnet.WithTCPAddrResolver(func(network, address string) (*net.TCPAddr, error) {
		fmt.Println(network, address)
		return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9000}, nil
	})

	addr, err := vnet.CreateAddressWithOptions("localhost", 9000, resolver)
	fmt.Println(addr.String())
	fmt.Println(err)
	// Output:
	// tcp localhost:9000
	// 127.0.0.1:9000
	// <nil>
}

func ExampleGetNetworkInterfacesWithOptions() {
	interfaces, err := vnet.GetNetworkInterfacesWithOptions(vnet.WithInterfacesFunc(func() ([]net.Interface, error) {
		return []net.Interface{{Name: "eth0", Index: 1}}, nil
	}))

	fmt.Println(interfaces[0].Name)
	fmt.Println(err)
	// Output:
	// eth0
	// <nil>
}

func ExampleGetNetworkInterfaceWithOptions() {
	iface, err := vnet.GetNetworkInterfaceWithOptions("eth0", vnet.WithInterfaceByNameFunc(func(name string) (*net.Interface, error) {
		return &net.Interface{Name: name, Index: 1}, nil
	}))
	fmt.Println(iface.Name)
	fmt.Println(err)
	// Output:
	// eth0
	// <nil>
}

func ExampleLocalIPsWithOptions() {
	interfaces := vnet.WithInterfacesFunc(func() ([]net.Interface, error) {
		return []net.Interface{{Name: "eth0", Index: 1}}, nil
	})
	addrs := vnet.WithInterfaceAddrsFunc(func(net.Interface) ([]net.Addr, error) {
		_, ipNet, _ := net.ParseCIDR("192.0.2.10/24")
		ipNet.IP = net.ParseIP("192.0.2.10")
		return []net.Addr{ipNet}, nil
	})

	fmt.Println(vnet.LocalIPsWithOptions(interfaces, addrs))
	// Output: [192.0.2.10]
}

func ExampleLocalIPv4sWithOptions() {
	interfaces := vnet.WithInterfacesFunc(func() ([]net.Interface, error) {
		return []net.Interface{{Name: "eth0", Index: 1}}, nil
	})
	addrs := vnet.WithInterfaceAddrsFunc(func(net.Interface) ([]net.Addr, error) {
		_, ipNet, _ := net.ParseCIDR("192.0.2.10/24")
		ipNet.IP = net.ParseIP("192.0.2.10")
		return []net.Addr{ipNet}, nil
	})
	fmt.Println(vnet.LocalIPv4sWithOptions(interfaces, addrs))
	// Output: [192.0.2.10]
}

func ExampleLocalIPv6sWithOptions() {
	interfaces := vnet.WithInterfacesFunc(func() ([]net.Interface, error) {
		return []net.Interface{{Name: "lo0", Index: 1}}, nil
	})
	addrs := vnet.WithInterfaceAddrsFunc(func(net.Interface) ([]net.Addr, error) {
		_, ipNet, _ := net.ParseCIDR("2001:db8::1/64")
		ipNet.IP = net.ParseIP("2001:db8::1")
		return []net.Addr{ipNet}, nil
	})
	fmt.Println(vnet.LocalIPv6sWithOptions(interfaces, addrs))
	// Output: [2001:db8::1]
}

func ExampleLocalAddressListWithOptions() {
	interfaces := vnet.WithInterfacesFunc(func() ([]net.Interface, error) {
		return []net.Interface{{Name: "eth0", Index: 1}}, nil
	})
	addrs := vnet.WithInterfaceAddrsFunc(func(net.Interface) ([]net.Addr, error) {
		_, ipNet, _ := net.ParseCIDR("192.0.2.10/24")
		ipNet.IP = net.ParseIP("192.0.2.10")
		return []net.Addr{ipNet}, nil
	})
	ips := vnet.LocalAddressListWithOptions(func(ip net.IP) bool {
		return ip.To4() != nil
	}, interfaces, addrs)
	fmt.Println(vnet.ToIPList(ips))
	// Output: [192.0.2.10]
}

func ExampleLocalAddressListByInterfaceWithOptions() {
	interfaces := vnet.WithInterfacesFunc(func() ([]net.Interface, error) {
		return []net.Interface{{Name: "eth0", Index: 1}, {Name: "lo0", Flags: net.FlagLoopback, Index: 2}}, nil
	})
	addrs := vnet.WithInterfaceAddrsFunc(func(iface net.Interface) ([]net.Addr, error) {
		_, ipNet, _ := net.ParseCIDR("192.0.2.10/24")
		ipNet.IP = net.ParseIP("192.0.2.10")
		return []net.Addr{ipNet}, nil
	})
	ips := vnet.LocalAddressListByInterfaceWithOptions(
		func(iface net.Interface) bool { return iface.Flags&net.FlagLoopback == 0 },
		func(ip net.IP) bool { return ip.To4() != nil },
		interfaces,
		addrs,
	)
	fmt.Println(vnet.ToIPList(ips))
	// Output: [192.0.2.10]
}

func ExampleGetLocalhostStrWithOptions() {
	interfaces := vnet.WithInterfacesFunc(func() ([]net.Interface, error) {
		return []net.Interface{{Name: "eth0", Index: 1}}, nil
	})
	addrs := vnet.WithInterfaceAddrsFunc(func(net.Interface) ([]net.Addr, error) {
		_, ipNet, _ := net.ParseCIDR("192.0.2.10/24")
		ipNet.IP = net.ParseIP("192.0.2.10")
		return []net.Addr{ipNet}, nil
	})
	fmt.Println(vnet.GetLocalhostStrWithOptions(interfaces, addrs))
	// Output: 192.0.2.10
}

func ExampleGetLocalhostWithOptions() {
	interfaces := vnet.WithInterfacesFunc(func() ([]net.Interface, error) {
		return []net.Interface{{Name: "eth0", Index: 1}}, nil
	})
	addrs := vnet.WithInterfaceAddrsFunc(func(net.Interface) ([]net.Addr, error) {
		_, ipNet, _ := net.ParseCIDR("192.0.2.10/24")
		ipNet.IP = net.ParseIP("192.0.2.10")
		return []net.Addr{ipNet}, nil
	})
	fmt.Println(vnet.GetLocalhostWithOptions(interfaces, addrs))
	// Output: 192.0.2.10
}

func ExampleGetLocalHostNameWithOptions() {
	reverse := vnet.WithReverseLookupFunc(func(string) ([]string, error) {
		return []string{"example.local."}, nil
	})
	interfaces := vnet.WithInterfacesFunc(func() ([]net.Interface, error) {
		return []net.Interface{{Name: "eth0", Index: 1}}, nil
	})
	addrs := vnet.WithInterfaceAddrsFunc(func(net.Interface) ([]net.Addr, error) {
		_, ipNet, _ := net.ParseCIDR("192.0.2.10/24")
		ipNet.IP = net.ParseIP("192.0.2.10")
		return []net.Addr{ipNet}, nil
	})

	fmt.Println(vnet.GetLocalHostNameWithOptions(reverse, interfaces, addrs))
	// Output: example.local
}

func ExampleGetLocalMACAddressWithOptions() {
	iface := net.Interface{Name: "eth0", HardwareAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}}
	opts := []vnet.InterfaceOption{
		vnet.WithInterfacesFunc(func() ([]net.Interface, error) {
			return []net.Interface{iface}, nil
		}),
	}
	fmt.Println(vnet.GetLocalMACAddressWithOptions(opts, "-"))
	// Output: aa-bb-cc-dd-ee-ff
}

func ExampleGetMACAddressWithOptions() {
	iface := net.Interface{Name: "eth0", HardwareAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}}
	interfaces := vnet.WithInterfacesFunc(func() ([]net.Interface, error) {
		return []net.Interface{iface}, nil
	})
	addrs := vnet.WithInterfaceAddrsFunc(func(net.Interface) ([]net.Addr, error) {
		_, ipNet, _ := net.ParseCIDR("192.0.2.10/24")
		ipNet.IP = net.ParseIP("192.0.2.10")
		return []net.Addr{ipNet}, nil
	})
	fmt.Println(vnet.GetMACAddressWithOptions(net.ParseIP("192.0.2.10"), []vnet.InterfaceOption{interfaces, addrs}, "-"))
	// Output: aa-bb-cc-dd-ee-ff
}

func ExampleGetHardwareAddressWithOptions() {
	iface := net.Interface{Name: "eth0", HardwareAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}}
	interfaces := vnet.WithInterfacesFunc(func() ([]net.Interface, error) {
		return []net.Interface{iface}, nil
	})
	addrs := vnet.WithInterfaceAddrsFunc(func(net.Interface) ([]net.Addr, error) {
		_, ipNet, _ := net.ParseCIDR("192.0.2.10/24")
		ipNet.IP = net.ParseIP("192.0.2.10")
		return []net.Addr{ipNet}, nil
	})
	fmt.Println(vnet.GetHardwareAddressWithOptions(net.ParseIP("192.0.2.10"), interfaces, addrs))
	// Output: aa:bb:cc:dd:ee:ff
}

func ExampleGetLocalHardwareAddressWithOptions() {
	iface := net.Interface{Name: "eth0", HardwareAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}}
	interfaces := vnet.WithInterfacesFunc(func() ([]net.Interface, error) {
		return []net.Interface{iface}, nil
	})
	fmt.Println(vnet.GetLocalHardwareAddressWithOptions(interfaces))
	// Output: aa:bb:cc:dd:ee:ff
}

func ExampleIsUsableLocalPortWithOptions() {
	listenerFactory := vnet.WithPortListenerFactory(func(network, address string) (net.Listener, error) {
		fmt.Println(network, address)
		return exampleListener{}, nil
	})
	fmt.Println(vnet.IsUsableLocalPortWithOptions(12345, vnet.WithPortNetwork("tcp4"), vnet.WithPortHost("127.0.0.1"), listenerFactory))
	// Output:
	// tcp4 127.0.0.1:12345
	// true
}

func ExampleGetUsableLocalPortWithOptions() {
	listenerFactory := vnet.WithPortListenerFactory(func(_ string, address string) (net.Listener, error) {
		_, port, _ := net.SplitHostPort(address)
		if port == "1026" {
			return exampleListener{}, nil
		}
		return nil, os.ErrPermission
	})
	port, err := vnet.GetUsableLocalPortWithOptions(listenerFactory)
	fmt.Println(port)
	fmt.Println(err)
	// Output:
	// 1026
	// <nil>
}

func ExampleGetUsableLocalPortFromWithOptions() {
	listenerFactory := vnet.WithPortListenerFactory(func(_ string, address string) (net.Listener, error) {
		_, port, _ := net.SplitHostPort(address)
		if port == "2002" {
			return exampleListener{}, nil
		}
		return nil, os.ErrPermission
	})
	port, err := vnet.GetUsableLocalPortFromWithOptions(2000, listenerFactory)
	fmt.Println(port)
	fmt.Println(err)
	// Output:
	// 2002
	// <nil>
}

func ExampleGetUsableLocalPortInRangeWithOptions() {
	listenerFactory := vnet.WithPortListenerFactory(func(_ string, address string) (net.Listener, error) {
		_, port, _ := net.SplitHostPort(address)
		if port == "3001" {
			return exampleListener{}, nil
		}
		return nil, os.ErrPermission
	})
	port, err := vnet.GetUsableLocalPortInRangeWithOptions(3000, 3002, listenerFactory)
	fmt.Println(port)
	fmt.Println(err)
	// Output:
	// 3001
	// <nil>
}

func ExampleGetUsableLocalPortsWithOptions() {
	listenerFactory := vnet.WithPortListenerFactory(func(_ string, address string) (net.Listener, error) {
		_, port, _ := net.SplitHostPort(address)
		if port == "4000" || port == "4002" {
			return exampleListener{}, nil
		}
		return nil, os.ErrPermission
	})
	ports, err := vnet.GetUsableLocalPortsWithOptions(2, 4000, 4003, listenerFactory)
	fmt.Println(ports)
	fmt.Println(err)
	// Output:
	// [4000 4002]
	// <nil>
}

func ExampleNewLocalPortGeneratorWithOptions() {
	listenerFactory := vnet.WithPortListenerFactory(func(_ string, _ string) (net.Listener, error) {
		return exampleListener{}, nil
	})
	generator := vnet.NewLocalPortGeneratorWithOptions(5000, listenerFactory)
	first, _ := generator.Gen()
	second, _ := generator.Gen()
	fmt.Println(first, second)
	// Output: 5000 5001
}

func ExampleConnectWithOptions() {
	dialer := &exampleDialer{}
	conn, err := vnet.ConnectWithOptions("example.test", 443, vnet.WithConnectNetwork("tcp4"), vnet.WithConnectDialer(dialer))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	fmt.Println(dialer.network)
	fmt.Println(dialer.address)
	// Output:
	// tcp4
	// example.test:443
}

func ExampleNetCatWithOptions() {
	dialer := &exampleDialer{data: make(chan []byte, 1)}
	err := vnet.NetCatWithOptions("example.test", 9000, []byte("hello"), vnet.WithConnectDialer(dialer))
	fmt.Println(err)
	fmt.Println(string(<-dialer.data))
	// Output:
	// <nil>
	// hello
}

func ExamplePingWithOptions() {
	dialer := &exampleDialer{}
	ok := vnet.PingWithOptions("example.test", vnet.WithPingDialer(dialer), vnet.WithPingPorts(8443), vnet.WithPingNetwork("tcp4"))
	fmt.Println(ok)
	fmt.Println(dialer.network)
	fmt.Println(dialer.address)
	// Output:
	// true
	// tcp4
	// example.test:8443
}

func ExampleIsOpenWithOptions() {
	dialer := &exampleDialer{}
	addr := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8443}
	fmt.Println(vnet.IsOpenWithOptions(addr, vnet.WithConnectDialer(dialer)))
	fmt.Println(dialer.address)
	// Output:
	// true
	// 127.0.0.1:8443
}

func ExampleGetRemoteAddress() {
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()
	fmt.Println(vnet.GetRemoteAddress(client) != "")
	// Output: true
}

func ExampleIsConnected() {
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()
	fmt.Println(vnet.IsConnected(client))
	// Output: true
}

func ExampleTLSVersion() {
	fmt.Println(vnet.TLSVersion("TLSv1.2") == tls.VersionTLS12)
	fmt.Println(vnet.TLSVersion("TLSv1.3") == tls.VersionTLS13)
	// Output:
	// true
	// true
}

func ExampleCreateTLSConfig() {
	config := vnet.CreateTLSConfig()
	fmt.Println(config.MinVersion == tls.VersionTLS12)
	// Output: true
}

func ExampleNewCertPool() {
	pool := vnet.NewCertPool()
	fmt.Println(pool != nil)
	// Output: true
}

func ExampleNewTLSConfigBuilder() {
	config := vnet.NewTLSConfigBuilder().
		SetMinVersion(tls.VersionTLS12).
		SetServerName("example.test").
		Build()
	fmt.Println(config.MinVersion == tls.VersionTLS12)
	fmt.Println(config.ServerName)
	// Output:
	// true
	// example.test
}

func ExampleAddRootCABytes() {
	builder := vnet.NewTLSConfigBuilder()
	err := vnet.AddRootCABytes(builder, []byte(exampleCertPEM))
	fmt.Println(err)
	fmt.Println(builder.Build().RootCAs != nil)
	// Output:
	// <nil>
	// true
}

func ExampleAddRootCAFileWithOptions() {
	builder := vnet.NewTLSConfigBuilder()
	err := vnet.AddRootCAFileWithOptions(builder, "ca.pem", vnet.WithTLSReadFile(func(path string) ([]byte, error) {
		fmt.Println(path)
		return []byte(exampleCertPEM), nil
	}))
	fmt.Println(err)
	fmt.Println(builder.Build().RootCAs != nil)
	// Output:
	// ca.pem
	// <nil>
	// true
}

func ExampleAddRootCAReader() {
	builder := vnet.NewTLSConfigBuilder()
	err := vnet.AddRootCAReader(builder, strings.NewReader(exampleCertPEM))
	fmt.Println(err)
	fmt.Println(builder.Build().RootCAs != nil)
	// Output:
	// <nil>
	// true
}

func ExampleAddRootCAReaderWithOptions() {
	builder := vnet.NewTLSConfigBuilder()
	err := vnet.AddRootCAReaderWithOptions(builder, strings.NewReader("ignored"), vnet.WithTLSReadAll(func(io.Reader) ([]byte, error) {
		return []byte(exampleCertPEM), nil
	}))
	fmt.Println(err)
	fmt.Println(builder.Build().RootCAs != nil)
	// Output:
	// <nil>
	// true
}

func ExampleNewUploadSetting() {
	setting := vnet.NewUploadSetting()
	fmt.Println(setting.MaxFileSize > 0)
	fmt.Println(setting.AllowFileExts)
	// Output:
	// true
	// true
}

func ExampleUploadFileContentType() {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, _ := writer.CreateFormFile("avatar", "avatar.txt")
	_, _ = part.Write([]byte("hello"))
	_ = writer.Close()

	request, _ := http.NewRequest(http.MethodPost, "/upload", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	form, _ := vnet.ParseMultipartForm(request, vnet.NewUploadSetting())
	file := form.Form.File["avatar"][0]
	fmt.Println(vnet.UploadFileName(file))
	fmt.Println(vnet.UploadFileSize(file))
	fmt.Println(vnet.UploadFileContentType(file))
	// Output:
	// avatar.txt
	// 5
	// application/octet-stream
}

func ExampleSaveUploadedFile() {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, _ := writer.CreateFormFile("avatar", "avatar.txt")
	_, _ = part.Write([]byte("hello"))
	_ = writer.Close()

	request, _ := http.NewRequest(http.MethodPost, "/upload", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	form, _ := vnet.ParseMultipartForm(request, vnet.NewUploadSetting())
	dest := filepath.Join(os.TempDir(), "knifer-go-vnet-upload-example.txt")
	defer os.Remove(dest)

	err := vnet.SaveUploadedFile(form.Form.File["avatar"][0], dest, vnet.WithUploadOverwrite(true))
	data, _ := os.ReadFile(dest)
	fmt.Println(err)
	fmt.Println(string(data))
	// Output:
	// <nil>
	// hello
}

func ExampleParseMultipartForm() {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	_ = writer.WriteField("name", "gopher")
	part, _ := writer.CreateFormFile("avatar", "avatar.txt")
	_, _ = part.Write([]byte("hello"))
	_ = writer.Close()

	request, _ := http.NewRequest(http.MethodPost, "/upload", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	form, err := vnet.ParseMultipartForm(request, vnet.NewUploadSetting())

	fmt.Println(form.GetParam("name"))
	fmt.Println(vnet.UploadFileName(form.Form.File["avatar"][0]))
	fmt.Println(vnet.UploadFileSize(form.Form.File["avatar"][0]))
	fmt.Println(err)
	// Output:
	// gopher
	// avatar.txt
	// 5
	// <nil>
}
