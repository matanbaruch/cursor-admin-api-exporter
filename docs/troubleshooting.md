# Troubleshooting Guide

This guide helps you diagnose and resolve common issues with the Cursor Admin API Exporter.

## Quick Diagnostics

### Health Check

First, verify the exporter is running and healthy:

```bash
# Check if the service is running
curl -f http://localhost:8080/health

# Expected response:
# {"status":"healthy","timestamp":"2024-01-20T10:30:00Z"}
```

### Basic Metrics Test

```bash
# Check if metrics are being exported
curl -s http://localhost:8080/metrics | grep cursor_

# Should return multiple metrics starting with "cursor_"
```

### API Connectivity Test

```bash
# Test API connectivity directly
curl -H "Authorization: Bearer YOUR_TOKEN" https://api.cursor.com/admin/team/members

# Should return JSON with team member data
```

## Common Issues

### 1. Authentication Errors

#### Symptoms
- 401 Unauthorized errors in logs
- `cursor_exporter_scrape_errors_total` metric increasing
- No metrics being collected

#### Causes
- Invalid API token
- Expired API token
- Insufficient permissions
- Token not properly set in environment

#### Solutions

```bash
# 1. Verify token is set correctly
echo $CURSOR_API_TOKEN

# 2. Test token directly
curl -H "Authorization: Bearer $CURSOR_API_TOKEN" https://api.cursor.com/admin/team/members

# 3. Check token permissions in Cursor dashboard
# Navigate to Settings → API Keys and verify admin permissions

# 4. Regenerate token if needed
# Create new token in Cursor dashboard with admin permissions
```

#### For Docker:
```bash
# Check environment variables in container
docker exec cursor-exporter env | grep CURSOR

# Recreate container with correct token
docker rm -f cursor-exporter
docker run -d \
  --name cursor-exporter \
  -p 8080:8080 \
  -e CURSOR_API_TOKEN=your_new_token \
  ghcr.io/matanbaruch/cursor-admin-api-exporter:latest
```

### 2. Network Connectivity Issues

#### Symptoms
- Connection timeout errors
- DNS resolution failures
- "No route to host" errors

#### Solutions

```bash
# 1. Test DNS resolution
nslookup api.cursor.com

# 2. Test connectivity from same network
curl -I https://api.cursor.com

# 3. Check firewall rules
# Ensure outbound HTTPS (port 443) is allowed

# 4. Test from inside container (if using Docker)
docker exec cursor-exporter curl -I https://api.cursor.com
```

#### For Corporate Networks:
```bash
# Set proxy if needed
export HTTPS_PROXY=http://proxy.company.com:8080
export HTTP_PROXY=http://proxy.company.com:8080

# For Docker
docker run -d \
  --name cursor-exporter \
  -p 8080:8080 \
  -e CURSOR_API_TOKEN=your_token \
  -e HTTPS_PROXY=http://proxy.company.com:8080 \
  ghcr.io/matanbaruch/cursor-admin-api-exporter:latest
```

### 3. No Metrics Appearing

#### Symptoms
- Exporter running but no metrics in Prometheus
- Empty response from `/metrics` endpoint
- Prometheus showing target as down

#### Solutions

```bash
# 1. Check metrics endpoint directly
curl http://localhost:8080/metrics

# 2. Check if Prometheus is scraping the right endpoint
curl http://localhost:9090/api/v1/targets

# 3. Verify Prometheus configuration
# Check prometheus.yml for correct target configuration

# 4. Check exporter logs
docker logs cursor-exporter
# Or for binary: check application logs
```

#### Example Prometheus Configuration:
```yaml
scrape_configs:
  - job_name: 'cursor-admin-api-exporter'
    static_configs:
      - targets: ['localhost:8080']  # Adjust IP/port as needed
    scrape_interval: 30s
    metrics_path: /metrics
```

### 4. Container Won't Start

#### Symptoms
- Container exits immediately
- Docker shows "Exited (1)" status
- No response from health endpoint

#### Solutions

```bash
# 1. Check container logs
docker logs cursor-exporter

# 2. Check if port is already in use
netstat -tulpn | grep :8080

# 3. Try running interactively
docker run -it --rm \
  -p 8080:8080 \
  -e CURSOR_API_TOKEN=your_token \
  ghcr.io/matanbaruch/cursor-admin-api-exporter:latest

# 4. Check resource limits
docker stats cursor-exporter
```

### 5. High Memory Usage

#### Symptoms
- Container being killed (OOMKilled)
- High memory consumption
- Performance degradation

#### Solutions

```bash
# 1. Check memory usage
docker stats cursor-exporter

# 2. Set memory limits
docker run -d \
  --name cursor-exporter \
  -p 8080:8080 \
  -e CURSOR_API_TOKEN=your_token \
  --memory=128m \
  ghcr.io/matanbaruch/cursor-admin-api-exporter:latest

# 3. Monitor memory patterns
docker exec cursor-exporter cat /proc/meminfo
```

### 6. Rate Limiting Issues

#### Symptoms
- 429 Too Many Requests errors
- Intermittent failures
- High scrape error rates

#### Solutions

```bash
# 1. Increase scrape interval in Prometheus
# Change from 30s to 60s or higher

# 2. Check API rate limits in Cursor documentation

# 3. Monitor request patterns
# Check logs for request timing
```

#### Prometheus Configuration:
```yaml
scrape_configs:
  - job_name: 'cursor-admin-api-exporter'
    static_configs:
      - targets: ['localhost:8080']
    scrape_interval: 60s  # Increased from 30s
    scrape_timeout: 10s
```

### 7. Permission Denied Errors

#### Symptoms
- "Permission denied" in logs
- File access errors
- Container security issues

#### Solutions

```bash
# 1. Check file permissions
ls -la /path/to/config/file

# 2. Run container with correct user
docker run -d \
  --name cursor-exporter \
  -p 8080:8080 \
  -e CURSOR_API_TOKEN=your_token \
  --user $(id -u):$(id -g) \
  ghcr.io/matanbaruch/cursor-admin-api-exporter:latest

# 3. Check SELinux/AppArmor restrictions
# For SELinux: setsebool -P container_manage_cgroup on
# For AppArmor: check /var/log/audit/audit.log
```

### 8. Certificate/TLS Issues

#### Symptoms
- "Certificate verification failed"
- TLS handshake errors
- SSL/TLS connection issues

#### Solutions

```bash
# 1. Update CA certificates
# For container: image already includes updated certs
# For binary: update system CA certificates

# 2. Test TLS connection
openssl s_client -connect api.cursor.com:443 -servername api.cursor.com

# 3. Check system time
date
# Ensure system time is correct (certificate validation depends on it)
```

## Debugging Steps

### 1. Enable Debug Logging

```bash
# For Docker
docker run -d \
  --name cursor-exporter \
  -p 8080:8080 \
  -e CURSOR_API_TOKEN=your_token \
  -e LOG_LEVEL=debug \
  ghcr.io/matanbaruch/cursor-admin-api-exporter:latest

# For binary
export LOG_LEVEL=debug
./cursor-admin-api-exporter
```

### 2. Check Detailed Logs

```bash
# View all logs
docker logs cursor-exporter

# Follow logs in real-time
docker logs -f cursor-exporter

# View recent logs with timestamps
docker logs -t --since 1h cursor-exporter

# Filter for specific log levels
docker logs cursor-exporter 2>&1 | grep -i error
```

### 3. Test Individual Components

```bash
# Test health endpoint
curl -v http://localhost:8080/health

# Test metrics endpoint
curl -v http://localhost:8080/metrics

# Test root endpoint
curl -v http://localhost:8080/

# Test API connectivity
curl -v -H "Authorization: Bearer $CURSOR_API_TOKEN" https://api.cursor.com/admin/team/members
```

### 4. Monitor Resource Usage

```bash
# Check CPU and memory usage
docker stats cursor-exporter

# Check disk usage
docker exec cursor-exporter df -h

# Check network connections
docker exec cursor-exporter netstat -an
```

## Monitoring and Alerting

### Key Metrics to Monitor

```prometheus
# Scrape errors
rate(cursor_exporter_scrape_errors_total[5m])

# Scrape duration
cursor_exporter_scrape_duration_seconds

# Missing metrics (should be > 0)
cursor_team_members_total
```

### Recommended Alerts

```yaml
# High error rate
- alert: CursorExporterHighErrorRate
  expr: rate(cursor_exporter_scrape_errors_total[5m]) > 0.1
  for: 5m

# Long scrape duration
- alert: CursorExporterSlowScrape
  expr: cursor_exporter_scrape_duration_seconds > 10
  for: 5m

# No metrics
- alert: CursorExporterNoMetrics
  expr: cursor_team_members_total == 0
  for: 15m
```

## Getting Help

### Log Analysis

When reporting issues, include:

1. **Error logs** (with timestamps)
2. **Configuration** (environment variables, minus sensitive data)
3. **System information** (OS, Docker version, etc.)
4. **Network setup** (proxy, firewall, etc.)

### Debug Information to Collect

```bash
# System info
uname -a
docker version  # if using Docker

# Network connectivity
curl -I https://api.cursor.com

# Container status
docker ps -a | grep cursor-exporter
docker logs cursor-exporter

# Prometheus configuration
cat /path/to/prometheus.yml

# Exporter health
curl -f http://localhost:8080/health
curl -s http://localhost:8080/metrics | head -20
```

### Community Support

1. **Search existing issues**: [GitHub Issues](https://github.com/matanbaruch/cursor-admin-api-exporter/issues)
2. **Create new issue**: Include debug information above
3. **Community discussions**: [GitHub Discussions](https://github.com/matanbaruch/cursor-admin-api-exporter/discussions)

### Common Log Patterns

```
# Normal operation
INFO Starting HTTP server address=:8080
INFO Starting Cursor metrics collection

# Authentication issues
ERROR Failed to get team members error="API request failed with status 401"

# Network issues
ERROR Failed to make request error="dial tcp: lookup api.cursor.com: no such host"

# Rate limiting
ERROR API request failed with status 429: Too Many Requests
```

## Recovery Procedures

### Service Recovery

```bash
# 1. Stop service
docker stop cursor-exporter

# 2. Check logs for root cause
docker logs cursor-exporter

# 3. Fix configuration issue

# 4. Restart service
docker start cursor-exporter

# 5. Verify health
curl -f http://localhost:8080/health
```

### Data Recovery

```bash
# If metrics are missing, check:
# 1. API connectivity
# 2. Token permissions
# 3. Prometheus scrape configuration
# 4. Time synchronization
```

---

**Related Documentation:**
- [Configuration Guide](configuration.md)
- [Installation Guide](installation/)
- [Metrics Reference](metrics.md)