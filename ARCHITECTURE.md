# Architecture

This document describes the architecture and design decisions of the Cursor Admin API Exporter.

## Overview

The Cursor Admin API Exporter is designed as a lightweight, scalable Prometheus exporter that collects metrics from the Cursor Admin API and exposes them in a format suitable for monitoring and alerting.

## High-Level Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│                 │     │                 │     │                 │
│   Prometheus    │────▶│   Cursor API    │────▶│   Cursor API    │
│                 │     │   Exporter      │     │   Server        │
│                 │     │                 │     │                 │
└─────────────────┘     └─────────────────┘     └─────────────────┘
                                 │
                                 ▼
                        ┌─────────────────┐
                        │                 │
                        │   Grafana       │
                        │   Dashboard     │
                        │                 │
                        └─────────────────┘
```

## Components

### 1. HTTP Server (`main.go`)

The main HTTP server provides:
- **Metrics Endpoint** (`/metrics`): Exposes Prometheus metrics
- **Health Endpoint** (`/health`): Health check for monitoring
- **Root Endpoint** (`/`): Information page with links
- **Graceful Shutdown**: Proper cleanup on termination
- **Request Logging**: Debug logging for HTTP requests

### 2. Cursor API Client (`pkg/client/cursor.go`)

The API client handles:
- **Authentication**: Bearer token authentication
- **HTTP Client**: Configurable timeout and retry logic
- **Request/Response Handling**: JSON marshaling/unmarshaling
- **Error Handling**: Proper error propagation and logging

#### API Endpoints Supported:
- `GET /admin/team/members` - Team member information
- `GET /admin/usage/daily` - Daily usage statistics
- `GET /admin/spending` - Spending data
- `GET /admin/usage/events` - Usage events

### 3. Exporter System (`pkg/exporters/`)

The exporter system consists of:

#### Main Exporter (`exporter.go`)
- Coordinates all sub-exporters
- Implements Prometheus `Collector` interface
- Handles panic recovery
- Provides scrape metrics

#### Sub-Exporters

**Team Members Exporter** (`team_members.go`)
- Collects team member counts and role distribution
- Metrics: `cursor_team_members_total`, `cursor_team_members_by_role`

**Daily Usage Exporter** (`daily_usage.go`)
- Collects daily usage statistics
- Metrics: lines of code, suggestion acceptance, feature usage

**Spending Exporter** (`spending.go`)
- Collects spending and billing information
- Metrics: total spending, per-member costs, premium requests

**Usage Events Exporter** (`usage_events.go`)
- Collects granular usage events
- Metrics: event counts, token consumption, model usage

### 4. Configuration (`pkg/utils/config.go`)

Simple configuration management:
- Environment variable handling
- Default value support
- Type-safe configuration

## Design Patterns

### 1. Collector Pattern

Each exporter implements the Prometheus `Collector` interface:

```go
type Collector interface {
    Describe(chan<- *Desc)
    Collect(chan<- Metric)
}
```

This allows for:
- Lazy metric collection
- Dynamic metric registration
- Consistent error handling

### 2. Factory Pattern

Exporters are created using factory functions:

```go
func NewTeamMembersExporter(client *client.CursorClient) *TeamMembersExporter
```

Benefits:
- Consistent initialization
- Dependency injection
- Testable components

### 3. Error Handling Strategy

The exporter uses a defensive error handling approach:

1. **Panic Recovery**: All collector methods are wrapped with panic recovery
2. **Graceful Degradation**: Individual exporter failures don't crash the whole system
3. **Error Metrics**: Scrape errors are tracked and exposed as metrics
4. **Logging**: Comprehensive logging for debugging

## Data Flow

### 1. Metrics Collection Flow

```
Prometheus Scrape Request
         │
         ▼
    HTTP Server (/metrics)
         │
         ▼
    Main Exporter (Collect)
         │
         ▼
    Sub-Exporters (Parallel)
         │
         ▼
    Cursor API Client
         │
         ▼
    Cursor API Server
         │
         ▼
    Metric Processing
         │
         ▼
    Prometheus Response
```

### 2. Error Handling Flow

```
API Call Failure
         │
         ▼
    Error Logging
         │
         ▼
    Error Metric Increment
         │
         ▼
    Graceful Continuation
```

## Scalability Considerations

### 1. Horizontal Scaling

The exporter is designed to be stateless and can be horizontally scaled:
- Multiple instances can run simultaneously
- Each instance maintains its own metrics
- Load balancing can be used for high availability

### 2. Vertical Scaling

Resource usage considerations:
- **Memory**: Minimal memory footprint (~10-50MB)
- **CPU**: Low CPU usage except during scrapes
- **Network**: Dependent on API response sizes

### 3. Rate Limiting

To prevent API rate limiting:
- Configurable scrape intervals
- Efficient API usage patterns
- Caching considerations for future versions

## Security Architecture

### 1. Authentication

- **API Token**: Bearer token authentication with Cursor API
- **Token Storage**: Secure environment variable or secret management
- **Token Rotation**: Support for token rotation without restarts

### 2. Container Security

- **Non-root User**: Runs as user ID 65534 (nobody)
- **Minimal Base Image**: Alpine Linux base
- **Read-only Filesystem**: Container filesystem is read-only
- **Security Scanning**: Automated vulnerability scanning

### 3. Network Security

- **TLS**: All API communications use HTTPS
- **Certificate Validation**: Proper certificate validation
- **No Sensitive Data Exposure**: No sensitive data in logs or metrics

## Performance Characteristics

### 1. Latency

- **Scrape Latency**: Typically 1-5 seconds depending on API response time
- **Startup Time**: ~1-2 seconds for full initialization
- **Health Check**: <10ms response time

### 2. Throughput

- **Concurrent Scrapes**: Supports concurrent Prometheus scrapes
- **API Efficiency**: Batched API calls where possible
- **Resource Usage**: Minimal resource consumption

### 3. Reliability

- **Circuit Breaker**: Future consideration for API failures
- **Retry Logic**: Configurable retry mechanisms
- **Graceful Degradation**: Continues operation despite partial failures

## Monitoring and Observability

### 1. Built-in Metrics

The exporter provides metrics about itself:
- `cursor_exporter_scrape_duration_seconds`: Scrape performance
- `cursor_exporter_scrape_errors_total`: Error tracking
- Standard HTTP metrics (via middleware)

### 2. Logging

Structured logging with configurable levels:
- **Debug**: Detailed execution information
- **Info**: General operational information
- **Warn**: Warning conditions
- **Error**: Error conditions requiring attention

### 3. Health Checks

Health endpoint provides:
- Service availability status
- Timestamp information
- JSON response format

## Deployment Patterns

### 1. Standalone Deployment

- Single binary deployment
- Environment variable configuration
- Service management (systemd, etc.)

### 2. Container Deployment

- Docker containers
- Docker Compose for local development
- Container orchestration support

### 3. Kubernetes Deployment

- Helm charts for easy deployment
- ConfigMaps and Secrets for configuration
- ServiceMonitor for Prometheus integration
- HPA for automatic scaling

## Testing Strategy

### 1. Unit Tests

- Component-level testing
- Mock external dependencies
- High test coverage (>80%)

### 2. Integration Tests

- End-to-end API testing
- Real API interaction testing
- Error scenario testing

### 3. Performance Tests

- Load testing
- Memory leak detection
- Resource usage validation

## Future Enhancements

### 1. Caching

- Redis/Memcached support for API response caching
- Configurable cache TTL
- Cache invalidation strategies

### 2. Multi-tenant Support

- Support for multiple Cursor teams
- Tenant isolation
- Aggregated metrics

### 3. Advanced Features

- Historical data aggregation
- Alerting rule templates
- Custom dashboard generation

## Configuration Management

### 1. Environment Variables

Primary configuration method:
```bash
CURSOR_API_TOKEN=xxx
CURSOR_API_URL=https://api.cursor.com
LISTEN_ADDRESS=:8080
METRICS_PATH=/metrics
LOG_LEVEL=info
```

### 2. Kubernetes Configuration

Helm chart values:
```yaml
cursor:
  apiToken: "token"
  apiUrl: "https://api.cursor.com"

config:
  listenAddress: ":8080"
  metricsPath: "/metrics"
  logLevel: "info"
```

### 3. Docker Configuration

Environment variables and volumes:
```yaml
environment:
  - CURSOR_API_TOKEN=${CURSOR_API_TOKEN}
  - LOG_LEVEL=info
```

This architecture provides a robust, scalable, and maintainable solution for monitoring Cursor Admin API metrics.