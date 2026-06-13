package system

import (
	"net"
	"strings"
	"testing"
)

func TestHostInfo(t *testing.T) {
	h := GetHostInfo()
	if h == nil {
		t.Fatal("HostInfo 不应为 nil")
	}
	if h.GetName() == "" {
		t.Errorf("Host Name 不应为空")
	}
	if !strings.Contains(h.String(), "Host Name:") {
		t.Errorf("HostInfo.String 缺少 caption: %s", h.String())
	}
}

func TestHostInfoWithOptions(t *testing.T) {
	_, ipNet, err := net.ParseCIDR("10.0.0.2/24")
	if err != nil {
		t.Fatal(err)
	}
	ipNet.IP = net.ParseIP("10.0.0.2")
	h := NewHostInfoWithOptions(
		WithHostNameFunc(func() (string, error) { return "option-host", nil }),
		WithHostInterfaceAddrsFunc(func() ([]net.Addr, error) { return []net.Addr{ipNet}, nil }),
	)
	if h.GetName() != "option-host" || h.GetAddress() != "10.0.0.2" {
		t.Fatalf("NewHostInfoWithOptions = %#v", h)
	}

	h = GetHostInfoWithOptions(WithHostAddressFunc(func() string { return "192.0.2.10" }))
	if h.GetAddress() != "192.0.2.10" {
		t.Fatalf("GetHostInfoWithOptions address = %q", h.GetAddress())
	}
}
