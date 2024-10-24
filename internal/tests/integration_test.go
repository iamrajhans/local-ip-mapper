package internal_test

import (
	"local-ip-mapper/internal"
	"testing"
)

// TestFullNetworkScan tests the full workflow of scanning a network
func TestFullNetworkScan(t *testing.T) {
	// Simulated IP range and parameters
	networkRange := "192.168.1.0/30"
	timeout := 500
	startPort := 20
	endPort := 25
	portTimeout := 100

	// Perform a network scan
	results := internal.ScanNetwork(networkRange, timeout, startPort, endPort, portTimeout)

	// Validate results length (expecting 0 active IPs in this scenario)
	if len(results) > 1 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}

}
