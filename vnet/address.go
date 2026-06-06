package vnet

import (
	stdnet "net"

	netimpl "github.com/imajinyun/go-knifer/internal/net"
)

func BuildInetSocketAddress(host string, defaultPort int) (*stdnet.TCPAddr, error) {
	return netimpl.BuildInetSocketAddress(host, defaultPort)
}

func CreateAddress(host string, port int) *stdnet.TCPAddr { return netimpl.CreateAddress(host, port) }

func GetIPByHost(hostName string) string { return netimpl.GetIPByHost(hostName) }

func WithInterfaceByNameFunc(fn func(string) (*stdnet.Interface, error)) InterfaceOption {
	return netimpl.WithInterfaceByNameFunc(fn)
}

func WithInterfacesFunc(fn func() ([]stdnet.Interface, error)) InterfaceOption {
	return netimpl.WithInterfacesFunc(fn)
}

func WithInterfaceAddrsFunc(fn func(stdnet.Interface) ([]stdnet.Addr, error)) InterfaceOption {
	return netimpl.WithInterfaceAddrsFunc(fn)
}

func WithReverseLookupFunc(fn func(string) ([]string, error)) InterfaceOption {
	return netimpl.WithReverseLookupFunc(fn)
}

func WithNetHostnameFunc(fn func() (string, error)) InterfaceOption {
	return netimpl.WithNetHostnameFunc(fn)
}

func GetNetworkInterface(name string) (*stdnet.Interface, error) {
	return netimpl.GetNetworkInterface(name)
}

func GetNetworkInterfaceWithOptions(name string, opts ...InterfaceOption) (*stdnet.Interface, error) {
	return netimpl.GetNetworkInterfaceWithOptions(name, opts...)
}

func GetNetworkInterfaces() ([]stdnet.Interface, error) { return netimpl.GetNetworkInterfaces() }

func GetNetworkInterfacesWithOptions(opts ...InterfaceOption) ([]stdnet.Interface, error) {
	return netimpl.GetNetworkInterfacesWithOptions(opts...)
}

func LocalIPv4s() []string { return netimpl.LocalIPv4s() }

func LocalIPv4sWithOptions(opts ...InterfaceOption) []string {
	return netimpl.LocalIPv4sWithOptions(opts...)
}

func LocalIPv6s() []string { return netimpl.LocalIPv6s() }

func LocalIPv6sWithOptions(opts ...InterfaceOption) []string {
	return netimpl.LocalIPv6sWithOptions(opts...)
}

func LocalIPs() []string { return netimpl.LocalIPs() }

func LocalIPsWithOptions(opts ...InterfaceOption) []string {
	return netimpl.LocalIPsWithOptions(opts...)
}

func ToIPList(addressList []stdnet.IP) []string { return netimpl.ToIPList(addressList) }

func LocalAddressList(addressFilter func(stdnet.IP) bool) []stdnet.IP {
	return netimpl.LocalAddressList(addressFilter)
}

func LocalAddressListWithOptions(addressFilter func(stdnet.IP) bool, opts ...InterfaceOption) []stdnet.IP {
	return netimpl.LocalAddressListWithOptions(addressFilter, opts...)
}

func LocalAddressListByInterface(interfaceFilter func(stdnet.Interface) bool, addressFilter func(stdnet.IP) bool) []stdnet.IP {
	return netimpl.LocalAddressListByInterface(interfaceFilter, addressFilter)
}

func LocalAddressListByInterfaceWithOptions(interfaceFilter func(stdnet.Interface) bool, addressFilter func(stdnet.IP) bool, opts ...InterfaceOption) []stdnet.IP {
	return netimpl.LocalAddressListByInterfaceWithOptions(interfaceFilter, addressFilter, opts...)
}

func GetLocalhostStr() string { return netimpl.GetLocalhostStr() }

func GetLocalhostStrWithOptions(opts ...InterfaceOption) string {
	return netimpl.GetLocalhostStrWithOptions(opts...)
}

func GetLocalhost() stdnet.IP { return netimpl.GetLocalhost() }

func GetLocalhostWithOptions(opts ...InterfaceOption) stdnet.IP {
	return netimpl.GetLocalhostWithOptions(opts...)
}

func GetLocalHostName() string { return netimpl.GetLocalHostName() }

func GetLocalHostNameWithOptions(opts ...InterfaceOption) string {
	return netimpl.GetLocalHostNameWithOptions(opts...)
}

func GetLocalMACAddress(separator ...string) string { return netimpl.GetLocalMACAddress(separator...) }

func GetLocalMACAddressWithOptions(opts []InterfaceOption, separator ...string) string {
	return netimpl.GetLocalMACAddressWithOptions(opts, separator...)
}

func GetMACAddress(inetAddress stdnet.IP, separator ...string) string {
	return netimpl.GetMACAddress(inetAddress, separator...)
}

func GetMACAddressWithOptions(inetAddress stdnet.IP, opts []InterfaceOption, separator ...string) string {
	return netimpl.GetMACAddressWithOptions(inetAddress, opts, separator...)
}

func GetHardwareAddress(inetAddress stdnet.IP) stdnet.HardwareAddr {
	return netimpl.GetHardwareAddress(inetAddress)
}

func GetHardwareAddressWithOptions(inetAddress stdnet.IP, opts ...InterfaceOption) stdnet.HardwareAddr {
	return netimpl.GetHardwareAddressWithOptions(inetAddress, opts...)
}

func GetLocalHardwareAddress() stdnet.HardwareAddr { return netimpl.GetLocalHardwareAddress() }

func GetLocalHardwareAddressWithOptions(opts ...InterfaceOption) stdnet.HardwareAddr {
	return netimpl.GetLocalHardwareAddressWithOptions(opts...)
}

func GetRemoteAddress(conn stdnet.Conn) string { return netimpl.GetRemoteAddress(conn) }

func IsConnected(conn stdnet.Conn) bool { return netimpl.IsConnected(conn) }
