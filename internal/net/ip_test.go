package net

import (
	"errors"
	"math/big"
	stdnet "net"
	"reflect"
	"regexp"
	"strconv"
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
	var compiled string
	if !MatchesWildcardWithOptions("10.0.*.2", "10.0.1.2", WithWildcardCompileFunc(func(pattern string) (*regexp.Regexp, error) {
		compiled = pattern
		return regexp.Compile(pattern)
	})) {
		t.Fatal("MatchesWildcardWithOptions failed")
	}
	if compiled != `^10\.0\.\d{1,3}\.2$` {
		t.Fatalf("compiled wildcard pattern = %q", compiled)
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

func TestIPOptionsUseCustomParsers(t *testing.T) {
	parseIPCalls := 0
	parseIP := func(s string) stdnet.IP {
		parseIPCalls++
		switch s {
		case "aliasv4":
			return stdnet.ParseIP("10.1.2.3")
		case "aliasv6":
			return stdnet.ParseIP("::2")
		case "aliasmask":
			return stdnet.ParseIP("255.255.255.0")
		default:
			return stdnet.ParseIP(s)
		}
	}
	if got, err := IPv4ToLongWithOptions("aliasv4", WithIPParser(parseIP)); err != nil || got != 167838211 {
		t.Fatalf("IPv4ToLongWithOptions = %d %v", got, err)
	}
	if got := IPv4ToLongDefaultWithOptions("bad", 42, WithIPParser(func(string) stdnet.IP { return nil })); got != 42 {
		t.Fatalf("IPv4ToLongDefaultWithOptions = %d", got)
	}
	if got, err := IPv6ToBigIntWithOptions("aliasv6", WithIPParser(parseIP)); err != nil || got.Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("IPv6ToBigIntWithOptions = %v %v", got, err)
	}
	if !IsIPWithOptions("aliasv4", WithIPParser(parseIP)) || !IsIPv4WithOptions("aliasv4", WithIPParser(parseIP)) || !IsIPv6WithOptions("aliasv6", WithIPParser(parseIP)) {
		t.Fatal("IP validators should use custom parser")
	}
	if bit, err := MaskBitByMaskWithOptions("aliasmask", WithIPParser(parseIP)); err != nil || bit != 24 {
		t.Fatalf("MaskBitByMaskWithOptions = %d %v", bit, err)
	}
	if block, err := FormatIPBlockWithOptions("10.1.2.3", "aliasmask", WithIPParser(parseIP)); err != nil || block != "10.1.2.3/24" {
		t.Fatalf("FormatIPBlockWithOptions = %q %v", block, err)
	}
	if ips, err := ListIPsWithOptions("aliasv4/30", false, WithIPParser(parseIP), WithIPIntParser(strconv.Atoi)); err != nil || !reflect.DeepEqual(ips, []string{"10.1.2.1", "10.1.2.2"}) {
		t.Fatalf("ListIPsWithOptions = %#v %v", ips, err)
	}
	if parseIPCalls == 0 {
		t.Fatal("custom IP parser was not called")
	}
}

func TestIPRangeOptionsUseCustomParsers(t *testing.T) {
	parseCIDRCalls := 0
	_, network, err := stdnet.ParseCIDR("192.0.2.0/24")
	if err != nil {
		t.Fatal(err)
	}
	parseCIDR := func(s string) (stdnet.IP, *stdnet.IPNet, error) {
		parseCIDRCalls++
		if s != "alias-cidr" {
			return nil, nil, errors.New("unexpected cidr")
		}
		return stdnet.ParseIP("192.0.2.1"), network, nil
	}
	if !IsInRangeWithOptions("alias-ip", "alias-cidr", WithIPParser(func(s string) stdnet.IP {
		if s == "alias-ip" {
			return stdnet.ParseIP("192.0.2.5")
		}
		return nil
	}), WithCIDRParser(parseCIDR)) {
		t.Fatal("IsInRangeWithOptions should use custom parsers")
	}
	if parseCIDRCalls != 1 {
		t.Fatalf("parseCIDR calls = %d", parseCIDRCalls)
	}

	wildcardParseIntCalls := 0
	if !MatchesWildcardWithOptions("10.x.*.2", "alias", WithWildcardIPParser(func(s string) stdnet.IP {
		if s == "alias" {
			return stdnet.ParseIP("10.9.1.2")
		}
		return nil
	}), WithWildcardIntParser(func(s string) (int, error) {
		wildcardParseIntCalls++
		if s == "x" {
			return 9, nil
		}
		return strconv.Atoi(s)
	})) {
		t.Fatal("MatchesWildcardWithOptions should use custom parsers")
	}
	if wildcardParseIntCalls == 0 {
		t.Fatal("custom wildcard int parser was not called")
	}
}
