# Docker Installation

This guide covers installing the Cursor Admin API Exporter using Docker.

## Quick Start

```bash
docker run -d \
  --name cursor-exporter \
  -p 8080:8080 \
  -e CURSOR_API_TOKEN=your_token_here \
  ghcr.io/matanbaruch/cursor-admin-api-exporter:latest
```

## Prerequisites

- Docker Engine 20.10+ or Docker Desktop
- Cursor Admin API token
- 50MB available disk space

## Installation Steps

### 1. Pull the Image

```bash
# Pull the latest version
docker pull ghcr.io/matanbaruch/cursor-admin-api-exporter:latest

# Or pull a specific version
docker pull ghcr.io/matanbaruch/cursor-admin-api-exporter:v0.1.0
```

### 2. Run the Container

#### Basic Configuration

```bash
docker run -d \
  --name cursor-exporter \
  -p 8080:8080 \
  -e CURSOR_API_TOKEN=your_token_here \
  ghcr.io/matanbaruch/cursor-admin-api-exporter:latest
```

#### Advanced Configuration

```bash
docker run -d \
  --name cursor-exporter \
  -p 8080:8080 \
  -e CURSOR_API_TOKEN=your_token_here \
  -e CURSOR_API_URL=https://api.cursor.com \
  -e LISTEN_ADDRESS=:8080 \
  -e METRICS_PATH=/metrics \
  -e LOG_LEVEL=info \
  --restart=unless-stopped \
  ghcr.io/matanbaruch/cursor-admin-api-exporter:latest
```

### 3. Verify Installation

```bash
# Check container status
docker ps | grep cursor-exporter

# Check logs
docker logs cursor-exporter

# Test health endpoint
curl http://localhost:8080/health

# Test metrics endpoint
curl http://localhost:8080/metrics
```

## Configuration Options

### Environment Variables

```bash
# Required
-e CURSOR_API_TOKEN=your_token_here

# Optional
-e CURSOR_API_URL=https://api.cursor.com
-e LISTEN_ADDRESS=:8080
-e METRICS_PATH=/metrics
-e LOG_LEVEL=info
```

### Port Mapping

```bash
# Default port
-p 8080:8080

# Custom port
-p 9090:8080

# Bind to specific interface
-p 127.0.0.1:8080:8080
```

### Volume Mounting

```bash
# Mount configuration file
-v /path/to/config:/app/config

# Mount logs directory
-v /path/to/logs:/app/logs
```

## Security Configuration

### Non-Root User

The container runs as a non-root user (UID 65534) by default:

```bash
docker run -d \
  --name cursor-exporter \
  -p 8080:8080 \
  -e CURSOR_API_TOKEN=your_token_here \
  --user 65534:65534 \
  ghcr.io/matanbaruch/cursor-admin-api-exporter:latest
```

### Read-Only Filesystem

```bash
docker run -d \
  --name cursor-exporter \
  -p 8080:8080 \
  -e CURSOR_API_TOKEN=your_token_here \
  --read-only \
  --tmpfs /tmp:rw,noexec,nosuid,size=10m \
  ghcr.io/matanbaruch/cursor-admin-api-exporter:latest
```

### Resource Limits

```bash
docker run -d \
  --name cursor-exporter \
  -p 8080:8080 \
  -e CURSOR_API_TOKEN=your_token_here \
  --memory=128m \
  --cpus=0.1 \
  ghcr.io/matanbaruch/cursor-admin-api-exporter:latest
```

## Health Checks

### Built-in Health Check

```bash
# The image includes a built-in health check
docker run -d \
  --name cursor-exporter \
  -p 8080:8080 \
  -e CURSOR_API_TOKEN=your_token_here \
  --health-cmd="wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1" \
  --health-interval=30s \
  --health-timeout=3s \
  --health-retries=3 \
  ghcr.io/matanbaruch/cursor-admin-api-exporter:latest
```

### Check Health Status

```bash
# Check container health
docker inspect cursor-exporter | jq '.[0].State.Health'

# View health check logs
docker logs cursor-exporter --since 5m | grep health
```

## Networking

### Bridge Network (Default)

```bash
docker run -d \
  --name cursor-exporter \
  -p 8080:8080 \
  -e CURSOR_API_TOKEN=your_token_here \
  ghcr.io/matanbaruch/cursor-admin-api-exporter:latest
```

### Custom Network

```bash
# Create custom network
docker network create monitoring

# Run container on custom network
docker run -d \
  --name cursor-exporter \
  --network monitoring \
  -p 8080:8080 \
  -e CURSOR_API_TOKEN=your_token_here \
  ghcr.io/matanbaruch/cursor-admin-api-exporter:latest
```

### Host Network

```bash
docker run -d \
  --name cursor-exporter \
  --network host \
  -e CURSOR_API_TOKEN=your_token_here \
  -e LISTEN_ADDRESS=:8080 \
  ghcr.io/matanbaruch/cursor-admin-api-exporter:latest
```

## Monitoring Integration

### With Prometheus

```bash
# Run Prometheus alongside
docker run -d \
  --name prometheus \
  --network monitoring \
  -p 9090:9090 \
  -v ./prometheus.yml:/etc/prometheus/prometheus.yml \
  prom/prometheus:latest
```

### With Grafana

```bash
# Run Grafana alongside
docker run -d \
  --name grafana \
  --network monitoring \
  -p 3000:3000 \
  -e GF_SECURITY_ADMIN_PASSWORD=admin \
  grafana/grafana:latest
```

## Logging

### View Logs

```bash
# View all logs
docker logs cursor-exporter

# Follow logs
docker logs -f cursor-exporter

# View recent logs
docker logs --since 1h cursor-exporter

# View logs with timestamps
docker logs -t cursor-exporter
```

### Log Configuration

```bash
# Set log level
docker run -d \
  --name cursor-exporter \
  -p 8080:8080 \
  -e CURSOR_API_TOKEN=your_token_here \
  -e LOG_LEVEL=debug \
  ghcr.io/matanbaruch/cursor-admin-api-exporter:latest
```

## Troubleshooting

### Common Issues

1. **Container won't start**
   ```bash
   # Check logs for errors
   docker logs cursor-exporter
   
   # Check container status
   docker ps -a | grep cursor-exporter
   ```

2. **Permission denied**
   ```bash
   # Run as specific user
   docker run -d \
     --name cursor-exporter \
     -p 8080:8080 \
     -e CURSOR_API_TOKEN=your_token_here \
     --user $(id -u):$(id -g) \
     ghcr.io/matanbaruch/cursor-admin-api-exporter:latest
   ```

3. **Port already in use**
   ```bash
   # Use different port
   docker run -d \
     --name cursor-exporter \
     -p 8081:8080 \
     -e CURSOR_API_TOKEN=your_token_here \
     ghcr.io/matanbaruch/cursor-admin-api-exporter:latest
   ```

4. **API connection issues**
   ```bash
   # Test from inside container
   docker exec cursor-exporter curl -f https://api.cursor.com/admin/team/members
   ```

### Debug Mode

```bash
# Run in debug mode
docker run -d \
  --name cursor-exporter \
  -p 8080:8080 \
  -e CURSOR_API_TOKEN=your_token_here \
  -e LOG_LEVEL=debug \
  ghcr.io/matanbaruch/cursor-admin-api-exporter:latest

# View debug logs
docker logs -f cursor-exporter | grep -i debug
```

## Maintenance

### Update Container

```bash
# Pull latest image
docker pull ghcr.io/matanbaruch/cursor-admin-api-exporter:latest

# Stop and remove old container
docker stop cursor-exporter
docker rm cursor-exporter

# Run new container
docker run -d \
  --name cursor-exporter \
  -p 8080:8080 \
  -e CURSOR_API_TOKEN=your_token_here \
  ghcr.io/matanbaruch/cursor-admin-api-exporter:latest
```

### Backup Configuration

```bash
# Export environment variables
docker inspect cursor-exporter | jq '.[0].Config.Env'

# Save run command
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
  assaflavie/runlike cursor-exporter
```

### Cleanup

```bash
# Stop and remove container
docker stop cursor-exporter
docker rm cursor-exporter

# Remove image
docker rmi ghcr.io/matanbaruch/cursor-admin-api-exporter:latest

# Clean up unused resources
docker system prune -f
```

## Next Steps

After successful installation:

1. **Configure Monitoring**: Set up Prometheus and Grafana
2. **Set Up Alerts**: Configure alerting rules
3. **Production Deployment**: Consider Docker Compose or Kubernetes
4. **Security Review**: Implement security best practices

---

**Related Documentation:**
- [Docker Compose Installation](docker-compose.md)
- [Configuration Guide](../configuration.md)
- [Troubleshooting Guide](../troubleshooting.md)