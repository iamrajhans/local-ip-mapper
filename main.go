package main

import (
	"flag"
	"fmt"
	"local-ip-mapper/internal"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
)

func main() {
	setupLogging()
	// Parse command line arguments
	networkRange := flag.String("range", "", "The network range to scan (e.g., 192.168.1.0/24)")
	timeout := flag.Int("timeout", 1000, "Timeout in milliseconds for pinging IP addresses")
	export := flag.String("export", "", "Export the results to a file (csv or json format)")
	startPort := flag.Int("start-port", 1, "Start of the port range to scan (default: 1)")
	endPort := flag.Int("end-port", 1024, "End of the port range to scan (default: 1024)")
	portTimeout := flag.Int("port-timeout", 500, "Timeout in milliseconds for scanning each port (default: 500ms)")

	flag.Parse()

	if *networkRange == "" {
		localIP, networkRangeCIDR, err := internal.GetLocalNetworkInfo()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("No network range provided. Using the local network range: %s (%s)\n", networkRangeCIDR, localIP)
		*networkRange = networkRangeCIDR
	}
	results := internal.ScanNetwork(*networkRange, *timeout, *startPort, *endPort, *portTimeout)
	displaySummaryTable(results)

	if *export != "" {
		fmt.Printf("Exporting results to: %s\n", *export)
		// Placeholder for export functionality
	}
}

func setupLogging() {
	logFile, err := os.OpenFile("scan.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		os.Exit(1)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func displaySummaryTable(results []internal.ScanResult) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"IP Address", "MAC Address", "Open Ports"})

	for _, result := range results {
		openPortsStr := "None"
		if len(result.OpenPorts) > 0 {
			openPortsStr = fmt.Sprintf("%v", result.OpenPorts)
		}
		table.Append([]string{result.IP, result.MAC, openPortsStr})
	}

	table.Render()
}
