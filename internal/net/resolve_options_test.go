package net

import (
	"testing"
	"time"
)

func TestResolveWithOptions(t *testing.T) {
	ips, err := GetIPByHostWithOptions("localhost", WithResolveNetwork("ip4"), WithResolveTimeout(time.Second))
	if err != nil {
		t.Fatalf("GetIPByHostWithOptions: %v", err)
	}
	if len(ips) == 0 {
		t.Fatal("GetIPByHostWithOptions returned no IPs")
	}
	dns, err := GetDNSInfoWithOptions("localhost", WithDNSTypes("A"), WithResolveTimeout(time.Second))
	if err != nil {
		t.Fatalf("GetDNSInfoWithOptions: %v", err)
	}
	if len(dns) == 0 {
		t.Fatal("GetDNSInfoWithOptions returned no A records")
	}
}
