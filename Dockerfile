# Multi-stage Dockerfile for stroganoff

# Stage 1: Build
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Copy VERSION file
COPY VERSION .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-X github.com/yourusername/stroganoff/pkg/version.Version=$(cat VERSION) \
              -X github.com/yourusername/stroganoff/pkg/version.Commit=$(git rev-parse --short HEAD 2>/dev/null || echo unknown) \
              -X github.com/yourusername/stroganoff/pkg/version.BuildDate=$(date -u +'%Y-%m-%dT%H:%M:%SZ')" \
    -o stroganoff ./cmd/stroganoff

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /app/stroganoff /usr/local/bin/stroganoff

# Copy configuration template
COPY config.example.yaml ./config.yaml

# Expose default port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD stroganoff version || exit 1

# Run the application
ENTRYPOINT ["stroganoff"]
CMD ["web", "--host", "0.0.0.0"]
