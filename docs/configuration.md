---
title: Configuration
nav_order: 3
layout: default
---

# Configuration Guide

This guide covers all configuration options for the Cursor Admin API Exporter.

## Environment Variables

The exporter is configured primarily through environment variables:

### Required Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `CURSOR_API_TOKEN` | Cursor Admin API token | `cur_1234567890abcdef` |

### Optional Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `CURSOR_API_URL` | `https://api.cursor.com` | Cursor API endpoint URL |
| `LISTEN_ADDRESS` | `:8080` | HTTP server listen address |
| `METRICS_PATH` | `/metrics` | Path for metrics endpoint |
| `LOG_LEVEL` | `info` | Logging level (debug, info, warn, error) |

## Configuration Examples

### Basic Configuration

```bash
# Minimal configuration
export CURSOR_API_TOKEN="your_token_here"

# Run the exporter
./cursor-admin-api-exporter
```

### Advanced Configuration

```bash
# Complete configuration
export CURSOR_API_TOKEN="your_token_here"
export CURSOR_API_URL="https://api.cursor.com"
export LISTEN_ADDRESS="0.0.0.0:8080"
export METRICS_PATH="/metrics"
export LOG_LEVEL="debug"

# Run the exporter
./cursor-admin-api-exporter
```

### Docker Configuration

```bash
# Using environment variables
docker run -d \
  --name cursor-exporter \
  -p 8080:8080 \
  -e CURSOR_API_TOKEN=your_token_here \
  -e CURSOR_API_URL=https://api.cursor.com \
  -e LISTEN_ADDRESS=:8080 \
  -e METRICS_PATH=/metrics \
  -e LOG_LEVEL=info \
  ghcr.io/yourusername/cursor-admin-api-exporter:latest
```

### Docker Compose Configuration

```yaml
# docker-compose.yml
version: '3.8'

services:
  cursor-admin-api-exporter:
    image: ghcr.io/yourusername/cursor-admin-api-exporter:latest
    ports:
      - "8080:8080"
    environment:
      - CURSOR_API_TOKEN=${CURSOR_API_TOKEN}
      - CURSOR_API_URL=https://api.cursor.com
      - LISTEN_ADDRESS=:8080
      - METRICS_PATH=/metrics
      - LOG_LEVEL=info
    restart: unless-stopped
```

### Kubernetes Configuration

#### ConfigMap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cursor-exporter-config
data:
  CURSOR_API_URL: "https://api.cursor.com"
  LISTEN_ADDRESS: ":8080"
  METRICS_PATH: "/metrics"
  LOG_LEVEL: "info"
```

#### Secret

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: cursor-exporter-secret
type: Opaque
data:
  CURSOR_API_TOKEN: <base64-encoded-token>
```

#### Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cursor-exporter
spec:
  replicas: 1
  selector:
    matchLabels:
      app: cursor-exporter
  template:
    metadata:
      labels:
        app: cursor-exporter
    spec:
      containers:
      - name: cursor-exporter
        image: ghcr.io/yourusername/cursor-admin-api-exporter:latest
        ports:
        - containerPort: 8080
        env:
        - name: CURSOR_API_TOKEN
          valueFrom:
            secretKeyRef:
              name: cursor-exporter-secret
              key: CURSOR_API_TOKEN
        envFrom:
        - configMapRef:
            name: cursor-exporter-config
```

### Helm Configuration

```yaml
# values.yaml
cursor:
  apiToken: "your_token_here"
  apiUrl: "https://api.cursor.com"

config:
  listenAddress: ":8080"
  metricsPath: "/metrics"
  logLevel: "info"

resources:
  limits:
    cpu: 100m
    memory: 128Mi
  requests:
    cpu: 50m
    memory: 64Mi

serviceMonitor:
  enabled: true
  interval: 30s
```

## Security Configuration

### API Token Management

#### Environment Variables
```bash
# Direct environment variable (not recommended for production)
export CURSOR_API_TOKEN="your_token_here"
```

#### Docker Secrets
```yaml
# docker-compose.yml
services:
  cursor-exporter:
    image: ghcr.io/yourusername/cursor-admin-api-exporter:latest
    secrets:
      - cursor_api_token
    environment:
      - CURSOR_API_TOKEN_FILE=/run/secrets/cursor_api_token

secrets:
  cursor_api_token:
    file: ./cursor_api_token.txt
```

#### Kubernetes Secrets
```yaml
# Using External Secrets Operator
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: cursor-exporter-secret
spec:
  secretStoreRef:
    name: vault-backend
    kind: SecretStore
  target:
    name: cursor-exporter-secret
  data:
  - secretKey: CURSOR_API_TOKEN
    remoteRef:
      key: cursor/api-token
      property: token
```

### Network Security

#### TLS Configuration
```yaml
# For custom TLS certificates
volumes:
  - name: certs
    configMap:
      name: cursor-exporter-certs
      
containers:
- name: cursor-exporter
  volumeMounts:
  - name: certs
    mountPath: /etc/ssl/certs/ca-certificates.crt
    subPath: ca-certificates.crt
```

## Logging Configuration

### Log Levels

| Level | Description | Use Case |
|-------|-------------|----------|
| `debug` | Detailed debugging information | Development, troubleshooting |
| `info` | General operational information | Production (default) |
| `warn` | Warning conditions | Production |
| `error` | Error conditions only | Production (minimal logging) |

### Log Format

The exporter uses structured logging with the following fields:

```json
{
  "level": "info",
  "msg": "Starting HTTP server",
  "address": ":8080",
  "time": "2024-01-20T10:30:00Z"
}
```

### Log Configuration Examples

```bash
# Debug level for troubleshooting
export LOG_LEVEL="debug"

# Error level for minimal logging
export LOG_LEVEL="error"

# Info level for production (default)
export LOG_LEVEL="info"
```

## Performance Configuration

### Resource Limits

#### Docker
```yaml
# docker-compose.yml
services:
  cursor-exporter:
    image: ghcr.io/yourusername/cursor-admin-api-exporter:latest
    deploy:
      resources:
        limits:
          cpus: '0.1'
          memory: 128M
        reservations:
          cpus: '0.05'
          memory: 64M
```

#### Kubernetes
```yaml
resources:
  limits:
    cpu: 100m
    memory: 128Mi
  requests:
    cpu: 50m
    memory: 64Mi
```

### Scraping Configuration

#### Prometheus
```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'cursor-admin-api-exporter'
    static_configs:
      - targets: ['cursor-exporter:8080']
    scrape_interval: 30s
    scrape_timeout: 10s
    metrics_path: /metrics
```

#### ServiceMonitor
```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: cursor-exporter
spec:
  selector:
    matchLabels:
      app: cursor-exporter
  endpoints:
  - port: http
    interval: 30s
    path: /metrics
```

## Advanced Configuration

### Custom API Endpoints

```bash
# For custom Cursor API endpoints
export CURSOR_API_URL="https://custom-api.cursor.com"
```

### Multiple Instances

```yaml
# docker-compose.yml for multiple teams
services:
  cursor-exporter-team1:
    image: ghcr.io/yourusername/cursor-admin-api-exporter:latest
    ports:
      - "8080:8080"
    environment:
      - CURSOR_API_TOKEN=${TEAM1_API_TOKEN}
      - LISTEN_ADDRESS=:8080
      
  cursor-exporter-team2:
    image: ghcr.io/yourusername/cursor-admin-api-exporter:latest
    ports:
      - "8081:8080"
    environment:
      - CURSOR_API_TOKEN=${TEAM2_API_TOKEN}
      - LISTEN_ADDRESS=:8080
```

### Health Check Configuration

```yaml
# docker-compose.yml
services:
  cursor-exporter:
    image: ghcr.io/yourusername/cursor-admin-api-exporter:latest
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

## Validation

### Configuration Validation

```bash
# Test configuration
curl -f http://localhost:8080/health

# Check metrics endpoint
curl -f http://localhost:8080/metrics

# Verify specific metrics
curl -s http://localhost:8080/metrics | grep cursor_team_members_total
```

### Common Validation Commands

```bash
# Check if service is running
docker ps | grep cursor-exporter

# View logs
docker logs cursor-exporter

# Test API connectivity
curl -H "Authorization: Bearer your_token_here" https://api.cursor.com/admin/team/members

# Validate Prometheus scraping
curl -s http://localhost:9090/api/v1/targets | jq '.data.activeTargets[] | select(.job=="cursor-admin-api-exporter")'
```

## Troubleshooting Configuration

### Common Issues

1. **Invalid API Token**
   ```bash
   # Check token validity
   curl -H "Authorization: Bearer your_token_here" https://api.cursor.com/admin/team/members
   ```

2. **Network Connectivity**
   ```bash
   # Test from container
   docker exec cursor-exporter curl -f https://api.cursor.com/admin/team/members
   ```

3. **Permission Issues**
   ```bash
   # Check file permissions
   ls -la /path/to/config/file
   
   # Check container user
   docker exec cursor-exporter id
   ```

### Debug Mode

```bash
# Enable debug logging
export LOG_LEVEL="debug"

# Run with debug output
./cursor-admin-api-exporter 2>&1 | grep -i debug
```

For more troubleshooting information, see the [Troubleshooting Guide](troubleshooting.md).

---

**Next**: [Installation Guide](installation/) for deployment-specific instructions