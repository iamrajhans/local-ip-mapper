package internal

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
)

type ScanResult struct {
	IP        string `json:"ip"`
	MAC       string `json:"mac,omitempty"`
	OpenPorts []int  `json:"open_ports,omitempty"`
}

// GetLocalNetworkInfo retrieves the local IP address and network range
func GetLocalNetworkInfo() (string, string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				// Local IP and network range
				return ipNet.IP.String(), ipNet.String(), nil
			}
		}
	}

	return "", "", fmt.Errorf("unable to determine local IP address")
}

// GenerateIPList generates a list of IP addresses within the given network range in CIDR notation
func GenerateIPList(networkCIDR string) ([]string, error) {
	_, ipNet, err := net.ParseCIDR(networkCIDR)
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR notation: %v", err)
	}

	var ips []string
	for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); IncrementIP(ip) {
		ips = append(ips, ip.String())
	}

	// Remove the network address and broadcast address (if present)
	if len(ips) > 2 {
		return ips[1 : len(ips)-1], nil
	}

	return ips, nil
}

// incrementIP increments an IP address by 1 (for generating the IP list)
func IncrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// Ping checks if an IP address is reachable within the specified timeout
func Ping(ip string, timeout int, retries int) bool {
	for i := 0; i < retries; i++ {
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("ping", "-n", "1", "-w", fmt.Sprintf("%d", timeout), ip)
		} else {
			cmd = exec.Command("ping", "-c", "1", "-W", fmt.Sprintf("%d", timeout/1000), ip)
		}

		err := cmd.Run()
		if err == nil {
			return true
		}
	}
	return false
}

func ScanNetwork(networkRange string, timeout int, startPort int, endPort int, portTimeout int) []ScanResult {
	fmt.Printf("Starting scan on network: %s with timeout: %d ms\n", networkRange, timeout)

	// Generate the list of IPs to scan
	ipList, err := GenerateIPList(networkRange)
	if err != nil {
		fmt.Printf("Error generating IP list: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %d IPs to scan\n", len(ipList))

	// Initialize progress bar
	bar := pb.StartNew(len(ipList))

	var wg sync.WaitGroup
	maxGoroutines := runtime.NumCPU() * 10
	var results []ScanResult
	var resultsMutex sync.Mutex // Mutex to prevent data races when appending to results slice

	guard := make(chan struct{}, maxGoroutines)

	for _, ip := range ipList {
		wg.Add(1)
		guard <- struct{}{}

		go func(ip string) {
			defer wg.Done()
			defer func() { <-guard }()
			defer bar.Increment()
			retries := 2
			if Ping(ip, timeout, retries) {
				macAddress, err := GetMACAddress(ip)
				if err != nil {
					fmt.Printf("IP %s is active, but MAC address not found\n", ip)
				} else {
					fmt.Printf("IP %s is active, MAC Address: %s\n", ip, macAddress)
				}
				openPorts := ScanPorts(ip, startPort, endPort, portTimeout)
				result := ScanResult{
					IP:        ip,
					MAC:       macAddress,
					OpenPorts: openPorts,
				}

				resultsMutex.Lock()
				results = append(results, result)
				resultsMutex.Unlock()
			}
		}(ip)
	}

	wg.Wait()
	bar.Finish()
	fmt.Println("Network scan complete")

	return results
}

// GetMACAddress returns the MAC address for a given IP address using the `arp` command.
func GetMACAddress(ip string) (string, error) {
	cmd := exec.Command("arp", "-n", ip)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, ip) {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				return fields[3], nil
			}
		}
	}

	return "", fmt.Errorf("MAC address not found for IP: %s", ip)
}

// GetARPTable retrieves the entire ARP table
func GetARPTable() (map[string]string, error) {
	cmd := exec.Command("arp", "-a")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	arpTable := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			ip := strings.Trim(fields[0], "()")
			mac := fields[1]
			if net.ParseIP(ip) != nil {
				arpTable[ip] = mac
			}
		}
	}

	return arpTable, nil
}

// ScanPorts scans a range of ports on the given IP to determine if they are open
func ScanPorts(ip string, startPort, endPort, timeout int) []int {
	var openPorts []int
	for port := startPort; port <= endPort; port++ {
		address := fmt.Sprintf("%s:%d", ip, port)
		conn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Millisecond)
		if err == nil {
			openPorts = append(openPorts, port)
			conn.Close()
		}
	}
	return openPorts
}
