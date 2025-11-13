# stroganoff Deployment Guide

This guide covers various deployment scenarios for the stroganoff application.

## Standalone Binary Deployment

### Prerequisites
- Go 1.21+ (for building) or pre-built binary

### Building

```bash
# Clone or download the repository
git clone https://github.com/yourusername/stroganoff.git
cd stroganoff

# Build for your platform
make build

# Or build for multiple platforms
make build-all
```

### Running Standalone

```bash
# Create configuration
cp config.example.yaml config.yaml
# Edit config.yaml as needed

# Run the application
./dist/stroganoff web --config config.yaml
```

## Linux Service Deployment (systemd)

### Installation

```bash
# Build the binary
make build

# Get the binary path
BINARY_PATH=$(pwd)/dist/stroganoff

# Install as service
sudo $BINARY_PATH install

# Or specify custom service name and user
sudo $BINARY_PATH install --service myapp --user appuser
```

### Service Management

```bash
# Start the service
sudo systemctl start stroganoff

# Stop the service
sudo systemctl stop stroganoff

# View logs
sudo journalctl -u stroganoff -f

# Enable auto-start
sudo systemctl enable stroganoff

# Disable auto-start
sudo systemctl disable stroganoff
```

### Service Unit File

The install command creates `/etc/systemd/system/stroganoff.service`:

```ini
[Unit]
Description=stroganoff Service
After=network.target

[Service]
Type=simple
User=root
ExecStart=/path/to/stroganoff web
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
```

## macOS Service Deployment (launchd)

### Installation

```bash
# Build the binary
make build

# Get the binary path
BINARY_PATH=$(pwd)/dist/gocr

# Install as service
sudo $BINARY_PATH install
```

### Service Management

```bash
# Start the service
launchctl start gocr

# Stop the service
launchctl stop gocr

# View logs
tail -f /var/log/gocr.log

# Remove service
launchctl unload ~/Library/LaunchAgents/com.gocr.plist
```

## Windows Service Deployment

### Installation

Open Command Prompt as Administrator:

```cmd
# Build the binary
go build -o gocr.exe ./cmd/stroganoff

# Install as service
gocr.exe install

# Or specify custom service name
gocr.exe install --service MyApp
```

### Service Management

```cmd
# Start the service
net start gocr

# Stop the service
net stop gocr

# View logs - Use Windows Event Viewer
# Services: Look for "gocr"
```

## Docker Deployment

### Building Docker Image

```bash
# Build the Docker image
docker build -t gocr:latest .

# Or with a tag
docker build -t myregistry/gocr:1.0.0 .
```

### Running with Docker

```bash
# Run with default configuration
docker run -p 8080:8080 gocr:latest

# Run with custom config
docker run -p 8080:8080 \
  -v $(pwd)/config.yaml:/app/config.yaml \
  gocr:latest

# Run with environment variables
docker run -p 8080:8080 \
  -e GOCR_THEME=dark \
  -e GOCR_LOG_LEVEL=debug \
  gocr:latest
```

### Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f gocr

# Stop services
docker-compose down
```

### Kubernetes Deployment

Example deployment manifest:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gocr
  labels:
    app: gocr
spec:
  replicas: 2
  selector:
    matchLabels:
      app: gocr
  template:
    metadata:
      labels:
        app: gocr
    spec:
      containers:
      - name: gocr
        image: gocr:latest
        ports:
        - containerPort: 8080
        env:
        - name: GOCR_THEME
          value: "dark"
        volumeMounts:
        - name: config
          mountPath: /app/config.yaml
          subPath: config.yaml
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
      volumes:
      - name: config
        configMap:
          name: gocr-config
---
apiVersion: v1
kind: Service
metadata:
  name: gocr
spec:
  selector:
    app: gocr
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: LoadBalancer
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: gocr-config
data:
  config.yaml: |
    server:
      host: "0.0.0.0"
      port: 8080
      theme: "dark"
    api:
      rate_limit: 100
      rate_limit_window: 60
      auth_enabled: true
      cors_enabled: true
```

### Kubernetes Deployment Steps

```bash
# Create namespace
kubectl create namespace gocr

# Create the deployment
kubectl apply -f deployment.yaml -n gocr

# Check deployment status
kubectl get deployment -n gocr

# View logs
kubectl logs -f deployment/gocr -n gocr

# Port forward
kubectl port-forward svc/gocr 8080:80 -n gocr
```

## Production Deployment Checklist

### Pre-Deployment
- [ ] Update VERSION file with new version
- [ ] Run full test suite: `make test`
- [ ] Generate coverage report: `make test-coverage`
- [ ] Review all changes in git log
- [ ] Create release notes
- [ ] Build binaries: `make build-all`

### Configuration
- [ ] Copy and customize `config.example.yaml` to `config.yaml`
- [ ] Set appropriate log level (warn or error for production)
- [ ] Enable authentication if needed
- [ ] Configure rate limiting based on expected load
- [ ] Set proper CORS origins
- [ ] Configure TLS if needed

### Deployment
- [ ] Backup existing configuration
- [ ] Deploy new binary
- [ ] Verify service starts successfully
- [ ] Check logs for errors
- [ ] Test health endpoint: `curl http://localhost:8080/health`
- [ ] Run smoke tests

### Post-Deployment
- [ ] Monitor logs for errors
- [ ] Check metrics and uptime
- [ ] Verify all endpoints are responding
- [ ] Test upgrade functionality
- [ ] Document deployment details
- [ ] Create rollback plan

## Reverse Proxy Configuration

### Nginx

```nginx
upstream gocr {
    server localhost:8080;
}

server {
    listen 80;
    server_name yourdomain.com;

    # Redirect to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name yourdomain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    # Security headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;

    location / {
        proxy_pass http://gocr;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
}
```

### Apache

```apache
<VirtualHost *:80>
    ServerName yourdomain.com
    Redirect permanent / https://yourdomain.com/
</VirtualHost>

<VirtualHost *:443>
    ServerName yourdomain.com

    SSLEngine on
    SSLCertificateFile /path/to/cert.pem
    SSLCertificateKeyFile /path/to/key.pem

    # Security headers
    Header always set Strict-Transport-Security "max-age=31536000; includeSubDomains"
    Header always set X-Frame-Options "SAMEORIGIN"
    Header always set X-Content-Type-Options "nosniff"

    ProxyPreserveHost On
    ProxyPass / http://localhost:8080/
    ProxyPassReverse / http://localhost:8080/

    # Logging
    ErrorLog ${APACHE_LOG_DIR}/gocr-error.log
    CustomLog ${APACHE_LOG_DIR}/gocr-access.log combined
</VirtualHost>
```

## Monitoring and Logging

### Log Rotation (Logrotate)

Create `/etc/logrotate.d/gocr`:

```
/var/log/gocr*.log {
    daily
    missingok
    rotate 14
    compress
    delaycompress
    notifempty
    create 0640 root root
    sharedscripts
    postrotate
        systemctl reload gocr > /dev/null 2>&1 || true
    endscript
}
```

### Health Monitoring

```bash
#!/bin/bash
# health-check.sh

ENDPOINT="http://localhost:8080/health"
TIMEOUT=5

response=$(curl -s -m $TIMEOUT $ENDPOINT)
status=$?

if [ $status -eq 0 ]; then
    echo "Health check passed: $response"
    exit 0
else
    echo "Health check failed with status: $status"
    exit 1
fi
```

## Troubleshooting Deployment

### Service Won't Start

1. Check logs:
   - Linux: `sudo journalctl -u gocr -n 50`
   - macOS: `cat /var/log/gocr.log`
   - Windows: Windows Event Viewer

2. Verify configuration:
   ```bash
   gocr config show
   ```

3. Test manually:
   ```bash
   ./dist/gocr web
   ```

### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill process if needed
kill -9 <PID>

# Or use a different port
gocr web --port 8081
```

### High CPU/Memory Usage

1. Check metrics:
   ```bash
   curl http://localhost:8080/api/metrics
   ```

2. Reduce rate limit window
3. Enable authentication to control access
4. Check for memory leaks in logs

## Backup and Recovery

### Configuration Backup

```bash
# Backup config
cp config.yaml config.yaml.backup.$(date +%s)

# Backup entire deployment
tar czf gocr-backup.tar.gz config.yaml logs/
```

### Recovery

```bash
# Restore from backup
cp config.yaml.backup.<timestamp> config.yaml

# Restart service
sudo systemctl restart gocr
```

## Updates and Upgrades

### Using the Upgrade Command

```bash
# Upgrade to latest version
gocr upgrade

# Upgrade to specific version
gocr upgrade --version v1.2.3

# For private repositories
gocr upgrade --token <github-token>
```

### Manual Upgrade

1. Backup current binary and config
2. Download new binary
3. Replace old binary
4. Restart service
5. Verify health check passes

## Performance Tuning

### Configuration

```yaml
server:
  read_timeout: 30
  write_timeout: 30

api:
  rate_limit: 1000
  rate_limit_window: 60
```

### System Parameters

```bash
# Increase file descriptors (Linux)
ulimit -n 65536

# Increase open connections
sysctl -w net.core.somaxconn=65535
```

## Security Hardening

1. Run service as non-root user
2. Use TLS/SSL for all connections
3. Enable authentication for all APIs
4. Restrict CORS origins to known hosts
5. Use firewall rules to limit access
6. Regularly update dependencies
7. Monitor security logs
8. Use reverse proxy with WAF (Web Application Firewall)
