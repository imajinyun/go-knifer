package vnet_test

import (
	"testing"

	"github.com/imajinyun/go-knifer/vnet"
)

func TestVNetFacade(t *testing.T) {
	v, err := vnet.IPv4ToLong("127.0.0.1")
	if err != nil || vnet.LongToIPv4(v) != "127.0.0.1" {
		t.Fatalf("IPv4 facade failed: %d %v", v, err)
	}
	if !vnet.IsIPv4("192.168.1.1") || !vnet.IsIPv6("::1") || !vnet.IsInnerIP("10.0.0.1") {
		t.Fatal("IP validators failed")
	}
	if !vnet.IsValidPort(80) || vnet.HideIPPart("192.168.1.2") != "192.168.1.*" {
		t.Fatal("port or hide helper failed")
	}
	if vnet.CreateTLSConfig(false) == nil || vnet.NewUploadSetting().MemoryThreshold == 0 {
		t.Fatal("TLS/upload helpers failed")
	}
}
