# Reverse Proxy

A simple, lightweight HTTP reverse proxy server written in Go that routes requests based on the Host header to different backend services.

### Prerequisites

- Go 1.19 or later

### Installation

1. Clone the repository:
```bash
git clone https://github.com/skyefactory/reverseproxy.git
cd reverseproxy
```

2. Build the application:
```bash
go build -o rproxy
```

3. Configure your routes (see [Configuration](#configuration))

4. Run the application:
```bash
sudo ./rproxy
```

The proxy will start on port 80 and begin routing requests according to the config file
## Configuration

### Routes Configuration

Edit the `config.cfg` file to define your routing rules. Each line follows the format:

```
hostname -> target_url
```

**Example:**
```
domain.com -> localhost:10000
subdomain.domain.com -> localhost:10001
api.domain.com -> http://192.168.1.100:8080
```

- Lines starting with `#` are treated as comments
- Empty lines are ignored
- The proxy will match the incoming `Host` header against the configured hostnames

### 404 Page

Place your custom 404 HTML template in `404.html`. If this file doesn't exist, a basic 404 response will be returned for unmatched routes.

## Logging

The proxy automatically logs all requests to `access.log` with the following information:

- Timestamp
- Host header
- Request URL path
- Client IP address (supports X-Forwarded-For header)
- HTTP method
- Protocol version
- User agent

**Log format:**
```
Time: 2025-08-01T10:30:45Z, Host: domain.com, URL Path: /api/users, Client IP: 192.168.1.10, Method: GET, Protocol: HTTP/1.1, User-Agent: Mozilla/5.0...
```