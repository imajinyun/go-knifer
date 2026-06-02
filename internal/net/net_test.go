package net

import (
	"crypto/tls"
	"math/big"
	"reflect"
	"testing"
)

func TestIPv4Helpers(t *testing.T) {
	v, err := IPv4ToLong("127.0.0.1")
	if err != nil || v != 2130706433 {
		t.Fatalf("IPv4ToLong = %d %v", v, err)
	}
	if got := LongToIPv4(v); got != "127.0.0.1" {
		t.Fatalf("LongToIPv4 = %q", got)
	}
	if !IsIPv4("192.168.1.1") || IsIPv4("999.1.1.1") || !IsIPv6("::1") || !IsIP("::1") {
		t.Fatal("IP validators failed")
	}
	if !IsInnerIP("192.168.1.1") || IsInnerIP("8.8.8.8") {
		t.Fatal("IsInnerIP failed")
	}
	if got, _ := BeginIP("192.168.1.9", 24); got != "192.168.1.0" {
		t.Fatalf("BeginIP = %q", got)
	}
	if got, _ := EndIP("192.168.1.9", 24); got != "192.168.1.255" {
		t.Fatalf("EndIP = %q", got)
	}
	if bit, _ := MaskBitByMask("255.255.255.0"); bit != 24 {
		t.Fatalf("MaskBitByMask = %d", bit)
	}
	if mask, _ := MaskByMaskBit(24); mask != "255.255.255.0" {
		t.Fatalf("MaskByMaskBit = %q", mask)
	}
	if count, _ := CountByMaskBit(30, false); count != 2 {
		t.Fatalf("CountByMaskBit = %d", count)
	}
	if ips, _ := ListIPCIDR("192.168.1.0", 30, false); !reflect.DeepEqual(ips, []string{"192.168.1.1", "192.168.1.2"}) {
		t.Fatalf("ListIPCIDR = %#v", ips)
	}
	if !MatchesWildcard("192.168.*.*", "192.168.1.2") || !IsInRange("192.168.1.2", "192.168.1.0/24") {
		t.Fatal("range matching failed")
	}
}

func TestIPv6BigInt(t *testing.T) {
	v, err := IPv6ToBigInt("::1")
	if err != nil || v.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("IPv6ToBigInt = %v %v", v, err)
	}
	if got, err := BigIntToIPv6(big.NewInt(1)); err != nil || got != "::1" {
		t.Fatalf("BigIntToIPv6 = %q %v", got, err)
	}
}

func TestPortAndMiscHelpers(t *testing.T) {
	if !IsValidPort(65535) || IsValidPort(70000) {
		t.Fatal("IsValidPort failed")
	}
	if HideIPPart("192.168.1.2") != "192.168.1.*" {
		t.Fatal("HideIPPart failed")
	}
	if got := GetMultistageReverseProxyIP("unknown, 10.0.0.1, 8.8.8.8"); got != "10.0.0.1" {
		t.Fatalf("GetMultistageReverseProxyIP = %q", got)
	}
	if IsUnknown("10.0.0.1") || !IsUnknown("unknown") {
		t.Fatal("IsUnknown failed")
	}
	if ascii, err := IDNToASCII("中国.cn"); err != nil || ascii == "" {
		t.Fatalf("IDNToASCII = %q %v", ascii, err)
	}
	if len(ParseCookies("a=1; b=2")) != 2 {
		t.Fatal("ParseCookies failed")
	}
}

func TestTLSHelpers(t *testing.T) {
	cfg := NewTLSConfigBuilder().SetMinVersion(tls.VersionTLS12).SetServerName("example.com").Build()
	if cfg.MinVersion != tls.VersionTLS12 || cfg.ServerName != "example.com" {
		t.Fatalf("TLS builder failed: %#v", cfg)
	}
	if TLSVersion(TLSv13) != tls.VersionTLS13 {
		t.Fatal("TLSVersion failed")
	}
}
