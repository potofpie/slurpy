# 🍹 Slurpy

A vibe coded POC to try to get a network tab for CLI app 

![Slurpy Demo](https://img.shields.io/badge/demo-available-brightgreen) ![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg) ![License](https://img.shields.io/badge/license-MIT-green.svg)

## ✨ Features

### 🔧 Slurpy SDK
- **Drop-in replacement** for Go's standard HTTP client
- **Automatic logging** of requests and responses
- **Namespace support** for organizing logs by project/service
- **Zero-overhead** when disabled
- **Request/response body capture** with automatic restoration
- **Duration tracking** for performance analysis
- **Error handling** and logging

### 🎨 Slurpy CLI
- **Beautiful TUI** built with Bubble Tea
- **Two-panel layout** similar to browser dev tools
- **Request filtering** and search capabilities
- **Detailed request/response inspection**
- **Namespace-based filtering**
- **Keyboard navigation** for efficiency
- **Color-coded status indicators**

## 🚀 Quick Start

### Installation

```bash
go get github.com/bobby/slurpy
```

### Basic Usage

```go
package main

import (
    "log"
    "github.com/bobby/slurpy/sdk"
)

func main() {
    // Create client with logging enabled
    client, err := slurpy.New(slurpy.Config{
        Namespace: "my-app",
        Enabled:   true,
    })
    if err != nil {
        log.Fatal(err)
    }

    // Use like any HTTP client - requests are automatically logged
    resp, err := client.Get("https://api.example.com/users")
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    // Request is now logged to ~/.config/slurpy/logs/
}
```

### Launch the CLI

```bash
# Build and run
go build -o slurpy ./cli
./slurpy

# Or use the Makefile
make cli
```

## 📖 Documentation

### SDK Configuration

```go
type Config struct {
    Namespace string // Unique identifier for this project
    Enabled   bool   // Enable/disable logging
}
```

**Default namespace:** `"default"`  
**Storage location:** `~/.config/slurpy/logs/`

### Supported HTTP Methods

The Slurpy client implements all standard HTTP methods:

```go
// Standard methods
resp, err := client.Get(url)
resp, err := client.Post(url, contentType, body)
resp, err := client.Put(url, contentType, body)
resp, err := client.Delete(url)

// Custom requests
resp, err := client.Do(req)
```

### Runtime Configuration

```go
// Change namespace dynamically
client.SetNamespace("new-service")

// Toggle logging on/off
client.SetEnabled(false) // Disable logging
client.SetEnabled(true)  // Re-enable logging

// Check current state
namespace := client.GetNamespace()
enabled := client.IsEnabled()
```

### Advanced Usage

```go
// Multiple clients with different namespaces
userClient, _ := slurpy.New(slurpy.Config{
    Namespace: "user-service",
    Enabled:   true,
})

paymentClient, _ := slurpy.New(slurpy.Config{
    Namespace: "payment-service", 
    Enabled:   true,
})

// Custom requests with headers
req, _ := http.NewRequest("GET", "https://api.example.com/data", nil)
req.Header.Set("Authorization", "Bearer token")
req.Header.Set("X-API-Key", "secret")

resp, err := client.Do(req)
```

## 🖥️ CLI Usage

### Key Bindings

| Key | Action |
|-----|--------|
| `↑/k`, `↓/j` | Navigate request list |
| `tab` | Switch between panels |
| `r` | Refresh requests |
| `c` | Clear current namespace |
| `?` | Toggle help |
| `q`/`esc` | Quit |

### Interface

```
┌─────────────────────────────┬─────────────────────────────┐
│        Request List         │      Request Details        │
│                             │                             │
│ ● GET /api/users [200]      │ REQUEST                     │
│   10:30:45 • 45ms • my-app  │ Method: GET                 │
│                             │ URL: https://api.../users   │
│ ● POST /api/users [201]     │ Duration: 45ms              │
│   10:30:50 • 123ms • my-app │                             │
│                             │ Request Headers:            │
│ ● PUT /api/users/1 [200]    │   Authorization: Bearer ... │
│   10:31:15 • 67ms • my-app  │   Content-Type: application │
│                             │                             │
│                             │ RESPONSE                    │
│                             │ Status: 200                 │
│                             │ Size: 1234 bytes            │
│                             │                             │
│                             │ Response Body:              │
│                             │ {"users": [...]}            │
└─────────────────────────────┴─────────────────────────────┘
```

### Features

- **Request List** (Left Panel): Shows all HTTP requests with method, URL, status, and timing
- **Request Details** (Right Panel): Detailed view of headers, body, and response data
- **Filtering**: Built-in search and filter capabilities
- **Color Coding**: Visual status indicators (green=success, red=error, yellow=pending)
- **Namespace Organization**: Requests grouped by service/project namespace

## 🏗️ Architecture

```
slurpy/
├── pkg/           # Shared data structures and utilities
│   ├── models/    # Request/response models
│   └── storage/   # File system storage management
├── sdk/           # Slurpy SDK for Go applications
├── cli/           # Bubble Tea TUI application
│   └── ui/        # UI components and styling
├── examples/      # Usage examples
│   ├── basic/     # Simple usage example
│   └── advanced/  # Advanced features demo
├── Makefile       # Build and development commands
└── test_slurpy.go # Comprehensive test suite
```

## 💾 Storage Format

Requests are stored as JSON files in `~/.config/slurpy/logs/` with the naming convention:
```
{namespace}_{request_id}.json
```

### Data Structure

```json
{
  "id": "abc123def456",
  "timestamp": "2024-01-15T10:30:00Z",
  "method": "GET",
  "url": "https://api.example.com/users",
  "headers": {"Authorization": "Bearer token"},
  "body": "",
  "response": {
    "status_code": 200,
    "headers": {"Content-Type": "application/json"},
    "body": "{\"users\": []}",
    "size": 123
  },
  "duration": "45ms",
  "namespace": "my-app",
  "error": ""
}
```

## 🛠️ Development

### Available Commands

```bash
# Development
make build       # Build the CLI
make cli         # Build and run CLI
make example     # Run basic example to generate test data
make run-example # Run example then start CLI
make deps        # Install dependencies

# Testing & Cleanup
make test        # Run tests
make clean       # Clean built artifacts and logs
make help        # Show available commands
```

### Testing

```bash
# Run comprehensive test suite
go run test_slurpy.go

# Generate test data with examples
go run ./examples/basic
go run ./examples/advanced
```

### Building

```bash
# Build CLI
go build -o slurpy ./cli

# Build for different platforms
GOOS=linux GOARCH=amd64 go build -o slurpy-linux ./cli
GOOS=windows GOARCH=amd64 go build -o slurpy.exe ./cli
```

## 📚 Examples

### Basic Integration

See [`examples/basic/main.go`](examples/basic/main.go) for a simple integration example.

### Advanced Features

See [`examples/advanced/main.go`](examples/advanced/main.go) for:
- Multiple namespaces
- Runtime configuration changes  
- Custom headers and requests
- CRUD operations
- Error handling

### Real-world Integration

```go
// In your existing application
func NewHTTPClient() *http.Client {
    if debug {
        client, _ := slurpy.New(slurpy.Config{
            Namespace: "my-service",
            Enabled:   true,
        })
        return client.Client
    }
    return &http.Client{}
}
```

## 🔧 Configuration

### Environment Variables

While Slurpy doesn't use environment variables directly, you can integrate them:

```go
client, err := slurpy.New(slurpy.Config{
    Namespace: os.Getenv("SERVICE_NAME"),
    Enabled:   os.Getenv("DEBUG") == "true",
})
```

### Storage Location

Default: `~/.config/slurpy/logs/`

To use a custom location, modify the storage package or implement your own storage backend.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Development Guidelines

- Follow Go best practices and conventions
- Add tests for new features
- Update documentation as needed
- Keep the CLI responsive and user-friendly
- Maintain backward compatibility

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Amazing TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Beautiful styling
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- Inspired by browser developer tools and network debugging needs

## 📞 Support

- 🐛 **Issues**: [GitHub Issues](https://github.com/bobby/slurpy/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/bobby/slurpy/discussions)
- 📖 **Documentation**: This README and inline code docs

---

**Happy debugging! 🍹✨** 
