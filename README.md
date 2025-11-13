# GOCR - Professional Go CLI Application

A comprehensive Go CLI application with advanced features including service installation, configuration management, web interface, API endpoints, monitoring, and automatic updates.

## Features

- **Multi-OS/Architecture Support**: Build for Linux, macOS, and Windows across multiple architectures
- **Service Installation**: Install as a systemd service (Linux), launchd service (macOS), or Windows Service
- **Configuration Management**: YAML-based configuration with hot-reload capability
- **Web Interface**: Responsive web UI with multiple theme support (default and dark themes)
- **REST API**: Full-featured API with authentication and rate limiting
- **Monitoring**: Built-in health checks and metrics collection
- **Automatic Updates**: Download and install new versions from GitHub releases
- **Security**: CORS, CSP, rate limiting, and multiple security headers
- **HTML Embedding**: All web assets embedded directly into the binary

## Quick Start

### Prerequisites

- Go 1.21 or later
- Make
- Git

### Building

```bash
# Build for current OS/ARCH
make build

# Build for all supported OS/ARCH combinations
make build-all

# Run the application
make run

# Run tests
make test

# Generate coverage report
make test-coverage
```

### Running the Web Server

```bash
# Start web server with default configuration
gocr web

# Start with custom host and port
gocr web --host 0.0.0.0 --port 8080

# Start with dark theme
gocr web --theme dark

# Start with custom config file
gocr web --config /path/to/config.yaml
```

### Commands

#### Version
Display version, commit, and build date information:
```bash
gocr version
```

#### Web Server
Start the HTTP server with web interface and API:
```bash
gocr web [flags]
```

Flags:
- `--config`: Path to configuration file (default: config.yaml)
- `--host`: Server host (default: localhost)
- `--port`: Server port (default: 8080)
- `--theme`: Theme name (default: dark, dark)

#### Upgrade
Download and install a new version:
```bash
gocr upgrade [--version VERSION] [--token TOKEN]
```

Flags:
- `--version`: Version to upgrade to (default: latest)
- `--token`: GitHub token for private repositories (optional)

#### Install
Install as a system service:
```bash
# Linux (requires root)
sudo gocr install

# macOS (requires admin)
sudo gocr install

# Windows (requires admin)
gocr install
```

Flags:
- `--service`: Service name (default: gocr)
- `--user`: User to run service as (Linux only)

#### Config
Manage configuration:
```bash
gocr config show
```

## Configuration

Copy `config.example.yaml` to `config.yaml` and customize:

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  theme: "default"  # or "dark"

api:
  rate_limit: 100
  rate_limit_window: 60
  auth_enabled: false
  cors_enabled: true
  allowed_origins:
    - "*"

logging:
  level: "info"
  format: "json"
  output_path: "stdout"
```

The configuration file is watched for changes and reloaded automatically.

## API Endpoints

### Public Endpoints

- `GET /health` - Health status check
- `GET /api/heartbeat` - Server heartbeat

### Protected Endpoints (requires authentication)

- `GET /api/metrics` - Application metrics
- `POST /api/auth/token` - Create authentication token

### Creating Tokens

```bash
curl -X POST http://localhost:8080/api/auth/token \
  -H "Content-Type: application/json" \
  -d '{
    "scopes": ["read", "write"],
    "duration": 86400
  }'
```

Response:
```json
{
  "token": "abc123..."
}
```

### Using Tokens

Include the token in the Authorization header:
```bash
curl -H "Authorization: Bearer abc123..." \
  http://localhost:8080/api/metrics
```

## Themes

Currently supported themes:
- **default** - Light theme with blue accents
- **dark** - Dark theme with light blue accents

Themes are organized in the `web/themes/` directory with the following structure:

```
web/themes/theme-name/
├── pages/
│   └── index.html
├── partials/
│   └── (component templates)
└── static/
    ├── css/
    ├── js/
    └── images/
```

## Version Management

The version is managed in the `VERSION` file using Semantic Versioning (MAJOR.MINOR.PATCH).

### Bumping Version

```bash
# Interactive version bump (patch, minor, major, or custom)
make version-bump

# Show current version
make version-show
```

Version information is hard-coded into the binary during compilation via ldflags.

## Release Process

Releases are automatically created via GitHub Actions when pushing to main:

1. Code is pushed to main branch
2. Tests run on all platforms
3. Binaries are built for all supported OS/ARCH combinations
4. A GitHub Release is created with all binaries
5. Release name is auto-generated with timestamp

## Development

### Project Structure

```
gocr/
├── cmd/gocr/
│   ├── main.go
│   └── commands/
│       ├── root.go
│       ├── version.go
│       ├── web.go
│       ├── upgrade.go
│       ├── install.go
│       └── config.go
├── internal/
│   ├── config/
│   ├── web/
│   ├── monitor/
│   ├── upgrade/
│   ├── install/
│   └── api/
├── pkg/
│   ├── version/
│   ├── auth/
│   └── ratelimit/
├── web/
│   └── themes/
│       ├── default/
│       └── dark/
└── .github/
    └── workflows/
```

### Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Run `golangci-lint` for linting
- Write tests for new features

## Security

- **X-Frame-Options**: Set to DENY to prevent clickjacking
- **X-Content-Type-Options**: Set to nosniff to prevent MIME sniffing
- **Content-Security-Policy**: Restricts resource loading
- **Rate Limiting**: Token bucket algorithm prevents abuse
- **Authentication**: Token-based authentication for protected endpoints
- **CORS**: Configurable cross-origin resource sharing

## Monitoring

The application includes built-in monitoring capabilities:

- **Heartbeat API**: Real-time server status
- **Metrics Endpoint**: Application performance metrics
- **Health Checks**: Extensible health check framework
- **Uptime Tracking**: Automatic uptime calculation

## Troubleshooting

### Service fails to start
1. Check logs: `journalctl -u gocr -f` (Linux)
2. Verify permissions: `sudo -l`
3. Check configuration: `gocr config show`

### Port already in use
```bash
# Find process using port
lsof -i :8080

# Use different port
gocr web --port 8081
```

### Hot-reload not working
- Ensure config file exists and is readable
- Check file permissions
- Monitor logs for errors

## License

MIT License - see LICENSE file for details

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Submit a pull request

## Support

For issues and feature requests, please visit the GitHub repository.
