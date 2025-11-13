# stroganoff - Professional Go CLI Application

A comprehensive Go CLI application with advanced features including service installation, configuration management, web interface, API endpoints, monitoring, and automatic updates.

**Also serves as a production-ready template for creating new Go CLI projects.**

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

### Using as a Template Generator

Stroganoff can generate new Go CLI projects with all the features pre-configured:

```bash
# Build the generator
make build-generator

# Generate a new project
./dist/stroganoff-generate \
  -name myapp \
  -module github.com/username/myapp \
  -app MyApp \
  -output ./myapp

# Set up the new project
cd myapp
go mod tidy
make build
./dist/myapp --help
```

**Parameters:**
- `-name`: Project name (lowercase, e.g., `myapp`)
- `-module`: Go module path (e.g., `github.com/username/myapp`)
- `-app`: Application display name (e.g., `MyApp`)
- `-output`: Output directory path

For detailed documentation on the template generator, see [TEMPLATE_GENERATOR.md](TEMPLATE_GENERATOR.md).

### Running the Web Server

```bash
# Start web server with default configuration
stroganoff web

# Start with custom host and port
stroganoff web --host 0.0.0.0 --port 8080

# Start with dark theme
stroganoff web --theme dark

# Start with custom config file
stroganoff web --config /path/to/config.yaml
```

### Commands

#### Version
Display version, commit, and build date information:
```bash
stroganoff version
```

#### Web Server
Start the HTTP server with web interface and API:
```bash
stroganoff web [flags]
```

Flags:
- `--config`: Path to configuration file (default: config.yaml)
- `--host`: Server host (default: localhost)
- `--port`: Server port (default: 8080)
- `--theme`: Theme name (default: dark, dark)

#### Upgrade
Download and install a new version:
```bash
stroganoff upgrade [--version VERSION] [--token TOKEN]
```

Flags:
- `--version`: Version to upgrade to (default: latest)
- `--token`: GitHub token for private repositories (optional)

#### Install
Install as a system service:
```bash
# Linux (requires root)
sudo stroganoff install

# macOS (requires admin)
sudo stroganoff install

# Windows (requires admin)
stroganoff install
```

Flags:
- `--service`: Service name (default: stroganoff)
- `--user`: User to run service as (Linux only)

#### Config
Manage configuration:
```bash
stroganoff config show
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
stroganoff/
├── cmd/
│   ├── stroganoff/
│   │   ├── main.go
│   │   └── commands/
│   │       ├── root.go
│   │       ├── version.go
│   │       ├── serve.go
│   │       ├── upgrade.go
│   │       ├── install.go
│   │       └── config.go
│   └── generate/
│       └── main.go              # Template generator tool
├── internal/
│   ├── config/                  # Configuration management
│   ├── web/                     # Web server and theme handling
│   ├── monitor/                 # Monitoring and health checks
│   ├── upgrade/                 # Auto-update functionality
│   ├── install/                 # Service installation
│   └── generator/               # Project template generator
├── pkg/
│   ├── version/                 # Version management
│   ├── auth/                    # Authentication and tokens
│   └── ratelimit/               # Rate limiting
├── web/
│   └── themes/
│       ├── default/             # Light theme
│       └── dark/                # Dark theme
└── .github/
    └── workflows/               # GitHub Actions CI/CD
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
1. Check logs: `journalctl -u stroganoff -f` (Linux)
2. Verify permissions: `sudo -l`
3. Check configuration: `stroganoff config show`

### Port already in use
```bash
# Find process using port
lsof -i :8080

# Use different port
stroganoff web --port 8081
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
