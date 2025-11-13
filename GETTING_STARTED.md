# stroganoff - Getting Started Guide

Welcome to stroganoff! This guide will help you get up and running quickly.

## Prerequisites

- Go 1.21 or later
- Make (usually pre-installed on macOS/Linux)
- Git

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/stroganoff.git
cd stroganoff
```

### 2. Verify Prerequisites

```bash
go version      # Should be 1.21+
make --version  # Should be available
```

### 3. Download Dependencies

```bash
make install-deps
```

## Your First Build

### Option 1: Quick Build (Current Platform)

```bash
make build
```

This creates a binary in `dist/stroganoff` (or `dist/stroganoff.exe` on Windows).

### Option 2: Build for All Platforms

```bash
make build-all
```

This builds for Linux, macOS, and Windows across multiple architectures.

## Running the Application

### Start the Web Server

```bash
# Copy example configuration
cp config.example.yaml config.yaml

# Run the web server
./dist/stroganoff web

# The server will start on http://localhost:8080
```

### Try Other Commands

```bash
# Show version
./dist/stroganoff version

# Show help
./dist/stroganoff --help

# Show config
./dist/stroganoff config show
```

## Configuration

### Default Configuration

The application uses `config.yaml` for configuration. A template is provided:

```bash
cp config.example.yaml config.yaml
```

### Common Configuration Changes

**Change server port:**
```yaml
server:
  port: 8888  # Change from 8080
```

**Enable authentication:**
```yaml
api:
  auth_enabled: true
```

**Switch theme:**
```bash
./dist/stroganoff web --theme dark
```

## Testing

### Run All Tests

```bash
make test
```

### Run Tests with Coverage

```bash
make test-coverage
```

This generates `coverage.html` which you can open in a browser.

### Run Specific Tests

```bash
go test -v ./pkg/auth
go test -v ./pkg/ratelimit
```

## Development Workflow

### 1. Making Changes

Edit files as needed. Example: adding a new API endpoint.

### 2. Format Code

```bash
make fmt
```

### 3. Lint Code

```bash
make lint
```

Or install golangci-lint first:
```bash
go install github.com/golangci-lint/golangci-lint/cmd/golangci-lint@latest
```

### 4. Run Tests

```bash
make test
```

### 5. Build and Test Locally

```bash
make build
./dist/stroganoff web
```

## Version Management

### Check Current Version

```bash
cat VERSION
make version-show
```

### Bump Version

```bash
make version-bump
```

This will prompt you to:
- Bump patch (0.1.0 → 0.1.1)
- Bump minor (0.1.0 → 0.2.0)
- Bump major (0.1.0 → 1.0.0)
- Enter custom version

## Adding a New Command

### 1. Create Command File

Create `cmd/stroganoff/commands/mycommand.go`:

```go
package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

var myCmd = &cobra.Command{
	Use:   "mycommand",
	Short: "My new command",
	Long:  "A longer description of my command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello from my command!")
	},
}

func init() {
	// Add flags here
	myCmd.Flags().StringP("name", "n", "World", "Name to greet")
}
```

### 2. Register Command

Edit `cmd/stroganoff/commands/root.go` and add:

```go
func init() {
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(myCmd)  // Add this line
	// ... other commands
}
```

### 3. Test Your Command

```bash
make build
./dist/stroganoff mycommand --help
./dist/stroganoff mycommand --name "Alice"
```

## Adding a New API Endpoint

### 1. Add Route to Server

Edit `internal/web/server.go` and add to `setupRoutes()`:

```go
api := s.engine.Group("/api")
{
	api.GET("/myendpoint", s.myEndpointHandler)
}
```

### 2. Add Handler Function

```go
func (s *Server) myEndpointHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello from my endpoint",
	})
}
```

### 3. Test Endpoint

```bash
make build
./dist/stroganoff web &
sleep 2
curl http://localhost:8080/api/myendpoint
```

## Docker Development

### Build Docker Image

```bash
docker build -t stroganoff:dev .
```

### Run in Docker

```bash
docker run -p 8080:8080 stroganoff:dev
```

### Using Docker Compose

```bash
docker-compose up
```

## Deployment

See [DEPLOYMENT.md](DEPLOYMENT.md) for detailed deployment instructions covering:
- Linux (systemd)
- macOS (launchd)
- Windows services
- Docker
- Kubernetes
- Reverse proxy setup

## Troubleshooting

### Build fails

**Issue:** Command not found: "go"

**Solution:** Install Go from https://go.dev/dl

**Issue:** Makefile error on Windows

**Solution:** Install GNU Make from https://gnuwin32.sourceforge.net/packages/make.htm

### Port already in use

**Issue:** Web server won't start, "address already in use"

**Solution:** Use different port
```bash
./dist/stroganoff web --port 8081
```

### Configuration file not found

**Issue:** "failed to load config file"

**Solution:** Create config file first
```bash
cp config.example.yaml config.yaml
```

### Tests fail

**Issue:** One or more tests failing

**Solution:** Check error messages and ensure all dependencies are installed
```bash
make install-deps
make test
```

### Permission denied on macOS/Linux

**Issue:** Cannot run binary: "Permission denied"

**Solution:** Make binary executable
```bash
chmod +x dist/stroganoff
```

## File Structure Explained

```
stroganoff/
├── cmd/              # Commands (CLI entry points)
├── internal/         # Private packages (not importable)
│   ├── config/      # Configuration management
│   ├── web/         # Web server and routes
│   ├── monitor/     # Health & metrics
│   ├── upgrade/     # Update mechanism
│   └── install/     # Service installation
├── pkg/              # Public packages
│   ├── version/     # Version info
│   ├── auth/        # Authentication
│   └── ratelimit/   # Rate limiting
├── web/              # Web assets
│   └── themes/      # Theme files
├── Makefile          # Build commands
├── go.mod            # Dependencies
├── VERSION           # Version file
└── config.example.yaml  # Config template
```

## Key Make Commands

| Command | Purpose |
|---------|---------|
| `make build` | Build for current OS |
| `make build-all` | Build for all platforms |
| `make test` | Run tests |
| `make test-coverage` | Generate coverage report |
| `make lint` | Run linter |
| `make fmt` | Format code |
| `make clean` | Remove build artifacts |
| `make version-bump` | Bump version |
| `make help` | Show all commands |

## Useful Go Commands

```bash
# Download dependencies
go mod download

# Update dependencies
go mod tidy

# Build with custom output
go build -o myapp ./cmd/stroganoff

# Run specific test
go test -run TestName ./...

# Show test coverage
go test -cover ./...

# Generate code
go generate ./...
```

## API Quick Reference

### Health Check
```bash
curl http://localhost:8080/health
```

### Heartbeat
```bash
curl http://localhost:8080/api/heartbeat
```

### Get Metrics (requires auth)
```bash
TOKEN="your-token-here"
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/metrics
```

### Create Token
```bash
curl -X POST http://localhost:8080/api/auth/token \
  -H "Content-Type: application/json" \
  -d '{"scopes": ["read"], "duration": 3600}'
```

## Next Steps

1. **Explore the codebase**: Read through the existing code to understand patterns
2. **Try the commands**: Run all the example commands to see what works
3. **Modify configuration**: Experiment with different config options
4. **Add features**: Create new commands and endpoints
5. **Test thoroughly**: Write tests for your new code
6. **Deploy**: Follow deployment guide for production setup

## Getting Help

- **Read documentation**: Check README.md and DEPLOYMENT.md
- **Check examples**: Look at existing commands and handlers
- **Run tests**: See how existing features are tested
- **GitHub Issues**: Report bugs or request features

## Common Tasks

### Building and Running
```bash
make build && ./dist/stroganoff version
```

### Quick Test Loop
```bash
make test && make build && ./dist/stroganoff web
```

### Check Everything
```bash
make lint && make test && make build-all
```

### Prepare for Release
```bash
make version-bump
make clean
make build-all
git add -A
git commit -m "Release vX.Y.Z"
git push origin main
```

## Tips & Tricks

### Speed up builds
Only build for your current platform during development:
```bash
make build  # Faster than make build-all
```

### Monitor file changes
Use file watchers to rebuild on changes:
```bash
# On macOS/Linux
while inotifywait -e modify -r . ; do make build; done
```

### Quick API testing
Create a test script:
```bash
#!/bin/bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/token \
  -H "Content-Type: application/json" \
  -d '{"scopes": ["read"], "duration": 3600}' | jq -r '.token')
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/metrics
```

## Support

For detailed information about specific features, see:
- [README.md](README.md) - Project overview
- [DEPLOYMENT.md](DEPLOYMENT.md) - Deployment guide
- [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) - Complete feature list

Happy coding! 🚀
