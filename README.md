# Local IP Mapper

Local IP Mapper is a command-line tool written in Go that helps users discover devices on a local network, identify active IP addresses, open ports, and provide detailed information about network activity. The tool includes features such as port scanning.

## Features

- Discover active devices on a local network
- Identify MAC addresses and vendors
- Port scanning to identify open services
- Supports customizable port ranges and network scan parameters

## Requirements

- Go 1.20+


## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/iamrajhans/local-ip-mapper.git
   cd local-ip-mapper
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

## Usage

Run the tool using the following command:

```bash
go run main.go -range 192.168.1.0/24 -timeout 1000 -start-port 20 -end-port 100 -port-timeout 300
```

### Command-Line Flags

- `-range`: The network range to scan (e.g., `192.168.1.0/24`).
- `-timeout`: Timeout in milliseconds for pinging IP addresses (default: 1000ms).
- `-start-port`: Start of the port range to scan (default: 1).
- `-end-port`: End of the port range to scan (default: 1024).
- `-port-timeout`: Timeout in milliseconds for scanning each port (default: 500ms).


## Running Tests

### Unit Tests

Run the unit tests with:

```bash
go test -v ./internal/tests
```

### Integration Tests

Run the integration tests with:

```bash
go test -v ./internal/tests/integration_test.go
```


## Contribution

Contributions are welcome! Please fork the repository and create a pull request with your changes. Ensure all tests pass before submitting.
