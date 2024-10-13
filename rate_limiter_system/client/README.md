
# Client

## Overview

The **Client** package simulates API requests to the server to test the rate-limiting behavior. It can make both sequential and concurrent requests to observe how the server handles load and applies rate limits based on IP addresses.

## Usage

### Running the Client

To start the client and make requests:

```bash
cd client/
go run main.go
```

### Configuration

The client reads IP addresses from `config/ip_address.json`. You can modify this file to test requests from different clients:

```json
{
  "ipAddresses": [
    "10.0.0.1",
    "10.0.0.2"
  ]
}
```

### Request Modes

- **Sequential Requests**: Send one request at a time.
- **Concurrent Requests**: Send multiple requests in parallel.

### Design Overview

- The client sends requests and logs the responses from the server to determine if the requests are accepted or rate-limited.
- The IP addresses are read dynamically from the configuration file, allowing easy adjustments for testing.

