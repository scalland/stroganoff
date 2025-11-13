.PHONY: help build run serve test clean version-bump install-deps lint docker-build docker-run fmt

# Variables
VERSION_FILE := VERSION
BINARY_NAME := gocr
DIST_DIR := dist
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
VERSION := $(shell cat $(VERSION_FILE))
LDFLAGS := -ldflags "-X github.com/yourusername/gocr/pkg/version.Version=$(VERSION) -X github.com/yourusername/gocr/pkg/version.Commit=$(GIT_COMMIT) -X github.com/yourusername/gocr/pkg/version.BuildDate=$(BUILD_DATE)"

# OS and ARCH combinations
OS_ARCH := linux/amd64 linux/arm64 linux/arm darwin/amd64 darwin/arm64 windows/amd64 windows/arm64

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build              Build binary for current OS/ARCH"
	@echo "  build-all          Build binaries for all OS/ARCH combinations"
	@echo "  run                Run the application"
	@echo "  serve              Build and start the server (port 8080)"
	@echo "  test               Run tests"
	@echo "  test-coverage      Run tests with coverage"
	@echo "  clean              Clean build artifacts"
	@echo "  lint               Run linter"
	@echo "  fmt                Format code"
	@echo "  install-deps       Install dependencies"
	@echo "  version-bump       Interactive version bumping (major/minor/patch)"
	@echo "  version-show       Show current version"
	@echo "  docker-build       Build Docker image"
	@echo "  docker-run         Run Docker container"

install-deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

build: install-deps
	@echo "Building $(BINARY_NAME) v$(VERSION) for $(GOOS)/$(GOARCH)..."
	go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME) ./cmd/gocr

build-all: clean install-deps
	@echo "Building $(BINARY_NAME) v$(VERSION) for all OS/ARCH combinations..."
	@mkdir -p $(DIST_DIR)
	@for os_arch in $(OS_ARCH); do \
		os=$${os_arch%/*}; \
		arch=$${os_arch#*/}; \
		binary=$(BINARY_NAME); \
		if [ "$$os" = "windows" ]; then \
			binary=$(BINARY_NAME).exe; \
		fi; \
		echo "Building for $$os/$$arch..."; \
		GOOS=$$os GOARCH=$$arch go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-$(VERSION)-$$os-$$arch/$(BINARY_NAME) ./cmd/gocr; \
		if [ "$$os" = "windows" ]; then \
			mv $(DIST_DIR)/$(BINARY_NAME)-$(VERSION)-$$os-$$arch/$(BINARY_NAME) $(DIST_DIR)/$(BINARY_NAME)-$(VERSION)-$$os-$$arch/$(BINARY_NAME).exe; \
		fi; \
		cd $(DIST_DIR)/$(BINARY_NAME)-$(VERSION)-$$os-$$arch && tar czf ../$(BINARY_NAME)-$(VERSION)-$$os-$$arch.tar.gz . && cd ../..; \
	done
	@echo "Build complete. Artifacts in $(DIST_DIR)/"

run: build
	@./$(DIST_DIR)/$(BINARY_NAME) --help

serve: build
	@cp config.example.yaml config.yaml 2>/dev/null || true
	@echo "Starting GOCR server..."
	@./$(DIST_DIR)/$(BINARY_NAME) serve

test: install-deps
	@echo "Running tests..."
	go test -v -race ./...

test-coverage: install-deps
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint: install-deps
	@echo "Running linter..."
	golangci-lint run ./... 2>/dev/null || echo "Install golangci-lint: https://golangci-lint.run/usage/install/"

fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w . 2>/dev/null || echo "Install goimports: go install golang.org/x/tools/cmd/goimports@latest"

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(DIST_DIR)
	rm -f coverage.out coverage.html

version-show:
	@echo "Current version: $(VERSION)"
	@echo "Git commit: $(GIT_COMMIT)"
	@echo "Build date: $(BUILD_DATE)"

version-bump:
	@echo "Current version: $(VERSION)"
	@echo "Select new version:"
	@echo "  1. Patch (0.1.0 -> 0.1.1)"
	@echo "  2. Minor (0.1.0 -> 0.2.0)"
	@echo "  3. Major (0.1.0 -> 1.0.0)"
	@echo "  4. Manual (enter custom version)"
	@read -p "Enter choice (1-4): " choice; \
	case $$choice in \
		1) \
			major=$$(echo $(VERSION) | cut -d. -f1); \
			minor=$$(echo $(VERSION) | cut -d. -f2); \
			patch=$$(echo $(VERSION) | cut -d. -f3 | cut -d- -f1); \
			new_patch=$$((patch + 1)); \
			new_version="$$major.$$minor.$$new_patch"; \
			echo "$$new_version" > $(VERSION_FILE); \
			echo "Version bumped to $$new_version"; \
			;; \
		2) \
			major=$$(echo $(VERSION) | cut -d. -f1); \
			minor=$$(echo $(VERSION) | cut -d. -f2); \
			minor=$$((minor + 1)); \
			new_version="$$major.$$minor.0"; \
			echo "$$new_version" > $(VERSION_FILE); \
			echo "Version bumped to $$new_version"; \
			;; \
		3) \
			major=$$(echo $(VERSION) | cut -d. -f1); \
			major=$$((major + 1)); \
			new_version="$$major.0.0"; \
			echo "$$new_version" > $(VERSION_FILE); \
			echo "Version bumped to $$new_version"; \
			;; \
		4) \
			read -p "Enter new version: " new_version; \
			echo "$$new_version" > $(VERSION_FILE); \
			echo "Version set to $$new_version"; \
			;; \
		*) echo "Invalid choice"; ;; \
	esac

docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):$(VERSION) .

docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(BINARY_NAME):$(VERSION)

.DEFAULT_GOAL := help
