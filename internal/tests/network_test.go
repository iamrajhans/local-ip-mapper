package internal_test

import (
	"local-ip-mapper/internal"
	"net"
	"testing"
)

// TestGenerateIPList tests the GenerateIPList function
func TestGenerateIPList(t *testing.T) {
	ipList, err := internal.GenerateIPList("192.168.1.0/30")
	if err != nil {
		t.Fatalf("Failed to generate IP list: %v", err)
	}

	expectedIPs := []string{"192.168.1.1", "192.168.1.2"}
	if len(ipList) != len(expectedIPs) {
		t.Fatalf("Expected %d IPs, got %d", len(expectedIPs), len(ipList))
	}

	for i, ip := range ipList {
		if ip != expectedIPs[i] {
			t.Errorf("Expected IP %s, got %s", expectedIPs[i], ip)
		}
	}
}

// TestIncrementIP tests the incrementIP function
func TestIncrementIP(t *testing.T) {
	ip := net.ParseIP("192.168.1.1")
	internal.IncrementIP(ip)

	expected := "192.168.1.2"
	if ip.String() != expected {
		t.Errorf("Expected %s, got %s", expected, ip.String())
	}
}
