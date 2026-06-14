package vnet_test

import (
	"errors"
	stdnet "net"
	"testing"

	"github.com/imajinyun/go-knifer/vnet"
)

type stubListener struct{}

func (stubListener) Accept() (stdnet.Conn, error) { return nil, errors.New("stub listener") }
func (stubListener) Close() error                 { return nil }
func (stubListener) Addr() stdnet.Addr {
	return &stdnet.TCPAddr{IP: stdnet.ParseIP("127.0.0.1"), Port: 12345}
}

func TestVNetProviderOptionsFacade(t *testing.T) {
	var network, address string
	addr, err := vnet.BuildInetSocketAddressWithOptions("example.com", 8080, vnet.WithAddressNetwork("tcp4"), vnet.WithTCPAddrResolver(func(n, a string) (*stdnet.TCPAddr, error) {
		network, address = n, a
		return &stdnet.TCPAddr{IP: stdnet.ParseIP("10.0.0.2"), Port: 8080}, nil
	}))
	if err != nil || addr.Port != 8080 {
		t.Fatalf("BuildInetSocketAddressWithOptions = %#v %v", addr, err)
	}
	if network != "tcp4" || address != "example.com:8080" {
		t.Fatalf("address resolver target = %s %s", network, address)
	}

	if !vnet.IsUsableLocalPortWithOptions(23456, vnet.WithPortNetwork("tcp4"), vnet.WithPortHost("127.0.0.2"), vnet.WithPortListenerFactory(func(n, a string) (stdnet.Listener, error) {
		network, address = n, a
		return stubListener{}, nil
	})) {
		t.Fatal("IsUsableLocalPortWithOptions should use listener factory")
	}
	if network != "tcp4" || address != "127.0.0.2:23456" {
		t.Fatalf("listener target = %s %s", network, address)
	}

	iface := stdnet.Interface{Name: "vnet0", HardwareAddr: stdnet.HardwareAddr{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}}
	_, ipNet, err := stdnet.ParseCIDR("10.9.8.7/24")
	if err != nil {
		t.Fatal(err)
	}
	ipNet.IP = stdnet.ParseIP("10.9.8.7")
	opts := []vnet.InterfaceOption{
		vnet.WithInterfaceByNameFunc(func(name string) (*stdnet.Interface, error) { return &iface, nil }),
		vnet.WithInterfacesFunc(func() ([]stdnet.Interface, error) { return []stdnet.Interface{iface}, nil }),
		vnet.WithInterfaceAddrsFunc(func(stdnet.Interface) ([]stdnet.Addr, error) { return []stdnet.Addr{ipNet}, nil }),
		vnet.WithReverseLookupFunc(func(string) ([]string, error) { return []string{"vnet.local."}, nil }),
		vnet.WithNetHostnameFunc(func() (string, error) { return "fallback", nil }),
	}
	gotIface, err := vnet.GetNetworkInterfaceWithOptions("vnet0", opts...)
	if err != nil || gotIface.Name != "vnet0" {
		t.Fatalf("GetNetworkInterfaceWithOptions = %#v %v", gotIface, err)
	}
	if got := vnet.LocalIPv4sWithOptions(opts...); len(got) != 1 || got[0] != "10.9.8.7" {
		t.Fatalf("LocalIPv4sWithOptions = %#v", got)
	}
	if got := vnet.GetLocalHostNameWithOptions(opts...); got != "vnet.local" {
		t.Fatalf("GetLocalHostNameWithOptions = %q", got)
	}
	if got := vnet.GetLocalMACAddressWithOptions(opts, "-"); got != "01-02-03-04-05-06" {
		t.Fatalf("GetLocalMACAddressWithOptions = %q", got)
	}
}
