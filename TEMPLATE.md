# 🎯 GOCR Template Repository

Welcome! This is a **production-ready GitHub template repository** for building professional Go CLI applications.

## What is a GitHub Template Repository?

A template repository is a repository that can be used as a base for new projects. Instead of forking, you can click **"Use this template"** button to create a new repository with the same directory structure and files.

## Features Included

✅ **Complete Go CLI Application** with all 12 professional features:
- Cobra CLI framework
- Gin HTTP web server
- YAML configuration with hot-reload
- API authentication & rate limiting
- Service installation (systemd, launchd, Windows)
- Automatic updates from GitHub releases
- Health monitoring and metrics
- Multiple themes support
- Docker & Kubernetes ready
- GitHub Actions CI/CD
- Comprehensive documentation

## How to Use This Template

### Option 1: Use GitHub's Template Feature (Recommended)

1. Go to the [GOCR repository](https://github.com/yourusername/gocr)
2. Click the green **"Use this template"** button
3. Choose a new repository name
4. Create your new repository

### Option 2: Clone and Customize

```bash
git clone https://github.com/yourusername/gocr.git my-app
cd my-app

# Update module name
sed -i '' 's/yourusername\/gocr/yourusername\/my-app/g' go.mod cmd/gocr/commands/*.go

# Remove template files if desired
rm TEMPLATE.md template.properties

git add -A
git commit -m "Initialize from GOCR template"
```

## Quick Start After Using Template

```bash
# Update your module name in go.mod
vi go.mod  # Change module path

# Update imports in code files
sed -i '' 's/yourusername\/gocr/your-username\/your-app/g' **/*.go

# Build
make build

# Run
./dist/your-app version
./dist/your-app serve

# Create first commit
git add -A
git commit -m "Initial project setup"
```

## What to Customize

### 1. **Project Name**
- Update `go.mod` with your module path
- Update imports throughout the codebase
- Update command name in `cmd/gocr/main.go`

### 2. **Documentation**
- Update `README.md` with your project details
- Update `config.example.yaml` with your config schema
- Update `DEPLOYMENT.md` if needed

### 3. **Application Logic**
- Add your commands in `cmd/gocr/commands/`
- Add API endpoints in `internal/web/server.go`
- Customize configuration in `internal/config/`

### 4. **Branding**
- Update web UI in `web/themes/`
- Update logo/styling
- Update version in `VERSION` file

### 5. **CI/CD**
- Update `.github/workflows/build.yml` for your needs
- Set repository secrets if needed
- Configure deployment settings

## Project Structure

```
your-app/
├── cmd/your-app/              # CLI application entry point
│   └── commands/              # Cobra commands
├── internal/                  # Private packages
│   ├── config/               # Configuration system
│   ├── web/                  # Web server & routes
│   ├── monitor/              # Health & metrics
│   ├── upgrade/              # Auto-update
│   └── install/              # Service installation
├── pkg/                       # Public packages
│   ├── version/              # Version management
│   ├── auth/                 # Authentication
│   └── ratelimit/            # Rate limiting
├── web/                       # Web assets
│   └── themes/               # UI themes
├── Makefile                   # Build automation
├── Dockerfile                 # Docker image
├── docker-compose.yml         # Docker Compose
├── go.mod                     # Go module
├── VERSION                    # Version file
├── config.example.yaml        # Config template
└── Documentation files        # README, DEPLOYMENT, etc.
```

## Key Files to Understand

| File | Purpose |
|------|---------|
| `go.mod` | Go module definition - **UPDATE THIS FIRST** |
| `cmd/gocr/commands/*.go` | CLI command implementations |
| `internal/config/` | Configuration management |
| `internal/web/server.go` | HTTP API routes and handlers |
| `Makefile` | Build and development tasks |
| `VERSION` | Semantic version (update with make version-bump) |
| `.github/workflows/build.yml` | CI/CD pipeline |

## Customization Checklist

- [ ] Update `go.mod` with your module path
- [ ] Update imports in all `.go` files
- [ ] Update README.md with your project description
- [ ] Update VERSION file (or leave as 0.1.0)
- [ ] Customize config.example.yaml for your needs
- [ ] Add your commands in `cmd/your-app/commands/`
- [ ] Add your API endpoints in `internal/web/server.go`
- [ ] Update GitHub Actions workflow if needed
- [ ] Update Docker image name and settings
- [ ] Customize web UI in `web/themes/`
- [ ] Remove TEMPLATE.md and template.properties (optional)
- [ ] Create your first git commit

## Common Tasks

### Build
```bash
make build              # Current platform
make build-all         # All platforms
```

### Run
```bash
make serve             # Start server (port 8080)
./dist/your-app serve --port 9000
```

### Test
```bash
make test              # Run tests
make test-coverage     # With coverage report
```

### Version Management
```bash
make version-show      # Show current version
make version-bump      # Interactive version bumping
```

### Deploy
```bash
# Linux (systemd)
sudo ./dist/your-app install
sudo systemctl start your-app

# macOS (launchd)
sudo ./dist/your-app install
launchctl start your-app

# Docker
docker build -t your-app:latest .
docker run -p 8080:8080 your-app:latest
```

## Documentation

- **README.md** - Project overview and features
- **GETTING_STARTED.md** - Step-by-step setup guide
- **DEPLOYMENT.md** - Production deployment options
- **PROJECT_SUMMARY.md** - Complete feature breakdown

## Included Technologies

- **Language**: Go 1.21+
- **CLI**: Cobra framework
- **HTTP**: Gin web framework
- **Config**: YAML with hot-reload (fsnotify)
- **Database**: Ready for any SQL or NoSQL
- **Auth**: Token-based API authentication
- **Rate Limiting**: Token bucket algorithm
- **Docker**: Multi-stage Dockerfile
- **CI/CD**: GitHub Actions workflow
- **Monitoring**: Health checks & metrics APIs

## Support

For issues or questions about the template:
1. Check the documentation in the repository
2. Review GETTING_STARTED.md for setup help
3. Check DEPLOYMENT.md for deployment help
4. See PROJECT_SUMMARY.md for feature details

## License

This template is provided as-is. When you create a new repository from this template, you can choose your own license.

---

## Quick Reference

| Task | Command |
|------|---------|
| Build | `make build` |
| Start server | `make serve` |
| Run tests | `make test` |
| Show version | `./dist/your-app version` |
| View help | `./dist/your-app --help` |
| Bump version | `make version-bump` |
| Build Docker | `make docker-build` |
| Install as service | `sudo ./dist/your-app install` |

---

**Ready to build your professional Go application!** 🚀

For detailed information, see the included documentation files in the repository.
