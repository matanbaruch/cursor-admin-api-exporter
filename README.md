# Cursor Admin API Exporter

<img width="673" height="375" alt="image" src="https://github.com/user-attachments/assets/9bedbbbd-8255-4c34-ae95-4e0813a14f02" />


<!-- Build and Quality Badges -->

[![Lint](https://github.com/matanbaruch/cursor-admin-api-exporter/actions/workflows/lint.yml/badge.svg)](https://github.com/matanbaruch/cursor-admin-api-exporter/actions/workflows/lint.yml)
[![Release](https://github.com/matanbaruch/cursor-admin-api-exporter/actions/workflows/release.yml/badge.svg)](https://github.com/matanbaruch/cursor-admin-api-exporter/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/matanbaruch/cursor-admin-api-exporter)](https://goreportcard.com/report/github.com/matanbaruch/cursor-admin-api-exporter)

<!-- Language and Tech Stack -->

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Prometheus](https://img.shields.io/badge/Prometheus-E6522C?style=for-the-badge&logo=Prometheus&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=for-the-badge&logo=kubernetes&logoColor=white)

<!-- Version and Distribution -->

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/matanbaruch/cursor-admin-api-exporter)](https://github.com/matanbaruch/cursor-admin-api-exporter/releases/latest)
[![Go Version](https://img.shields.io/github/go-mod/go-version/matanbaruch/cursor-admin-api-exporter)](https://github.com/matanbaruch/cursor-admin-api-exporter/blob/main/go.mod)
[![License](https://img.shields.io/github/license/matanbaruch/cursor-admin-api-exporter)](https://github.com/matanbaruch/cursor-admin-api-exporter/blob/main/LICENSE)

<!-- GitHub Stats -->

[![GitHub stars](https://img.shields.io/github/stars/matanbaruch/cursor-admin-api-exporter?style=social)](https://github.com/matanbaruch/cursor-admin-api-exporter/stargazers)

<!-- Distribution Platforms -->

[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/cursor-admin-api-exporter)](https://artifacthub.io/packages/search?repo=cursor-admin-api-exporter)

A Prometheus exporter for Cursor Admin API metrics, providing comprehensive monitoring and observability for your Cursor team's usage, spending, and productivity metrics.

## Features

- **Team Members Monitoring**: Track team size, roles, and member status
- **Daily Usage Metrics**: Monitor lines of code, AI suggestion acceptance rates, and feature usage
- **Spending Analytics**: Track per-member spending and premium request usage
- **Usage Events**: Granular tracking of token consumption and model usage
- **Multi-platform Support**: Available as Docker container, Kubernetes deployment, and native binaries
- **Security First**: Built with security best practices, including attestations and provenance
- **Production Ready**: Comprehensive monitoring, health checks, and graceful shutdown

## Quick Start

### Docker

```bash
# Pull the latest image
docker pull ghcr.io/matanbaruch/cursor-admin-api-exporter:latest

# Run with environment variables
docker run -p 8080:8080 \
  -e CURSOR_API_TOKEN=your_token_here \
  -e CURSOR_API_URL=https://api.cursor.com \
  ghcr.io/matanbaruch/cursor-admin-api-exporter:latest
```

### Docker Compose

```bash
# Copy the example environment file
cp env.example .env

# Edit .env with your Cursor API token
# Then start the services
docker-compose up -d
```

### Kubernetes (Helm)

```bash
# Install the chart
helm install cursor-admin-api-exporter oci://ghcr.io/matanbaruch/cursor-admin-api-exporter/charts/cursor-admin-api-exporter \
  --set cursor.apiToken=your_token_here \
  --set cursor.apiUrl=https://api.cursor.com
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `CURSOR_API_TOKEN` | Cursor Admin API token (required) | - |
| `CURSOR_API_URL` | Cursor API endpoint | `https://api.cursor.com` |
| `LISTEN_ADDRESS` | HTTP server listen address | `:8080` |
| `METRICS_PATH` | Metrics endpoint path | `/metrics` |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` |

### Getting a Cursor API Token

1. Log in to your Cursor team dashboard
2. Navigate to Settings → API Keys
3. Create a new API key with admin permissions
4. Copy the token and use it as `CURSOR_API_TOKEN`

## Metrics

The exporter provides the following metrics:

### Team Members
- `cursor_team_members_total` - Total number of team members
- `cursor_team_members_by_role` - Number of members by role (admin, member, viewer)

### Daily Usage
- `cursor_daily_lines_added_total` - Lines of code added per day
- `cursor_daily_lines_deleted_total` - Lines of code deleted per day
- `cursor_daily_suggestion_acceptance_rate` - AI suggestion acceptance rate
- `cursor_daily_tabs_used_total` - Tab completions used per day
- `cursor_daily_composer_used_total` - Composer usage per day
- `cursor_daily_chat_requests_total` - Chat requests per day
- `cursor_daily_model_usage` - Most used model per day
- `cursor_daily_extension_usage` - Most used extension per day

### Spending
- `cursor_spending_total_cents` - Total spending in cents
- `cursor_spending_by_member_cents` - Spending by team member
- `cursor_premium_requests_by_member_total` - Premium requests by member
- `cursor_premium_requests_total` - Total premium requests

### Usage Events
- `cursor_usage_events_total` - Total usage events
- `cursor_usage_events_by_type_total` - Events by type (completion, chat, etc.)
- `cursor_usage_events_by_user_total` - Events by user
- `cursor_usage_events_by_model_total` - Events by AI model
- `cursor_tokens_consumed_total` - Total tokens consumed
- `cursor_tokens_consumed_by_model_total` - Tokens consumed by model
- `cursor_tokens_consumed_by_user_total` - Tokens consumed by user

### Exporter Metrics
- `cursor_exporter_scrape_duration_seconds` - Time spent scraping the API
- `cursor_exporter_scrape_errors_total` - Total scrape errors

## Development

### Prerequisites

- Go 1.24 or later
- Docker (for container builds)
- Helm 3.x (for Kubernetes deployment)
- Make (for build automation)

### Building from Source

```bash
# Clone the repository
git clone https://github.com/matanbaruch/cursor-admin-api-exporter.git
cd cursor-admin-api-exporter

# Install dependencies
go mod download

# Build the binary
make build

# Run tests
make test

# Build for all platforms
make build-all
```

### Running Tests

```bash
# Run all tests
make test

# Run specific test suites
make test-unit
make test-integration
make test-performance

# Generate coverage report
make coverage

# Run coverage checks
make coverage-check

# Generate coverage badge
make coverage-badge
```

### Code Coverage

The project maintains comprehensive code coverage with the following thresholds:

- **Overall Coverage**: 80% minimum
- **Package Coverage**: 75% minimum  
- **Coverage Decrease**: 0.5% maximum allowed in PRs

```bash
# Generate comprehensive coverage report
make coverage

# Unit tests coverage only
make coverage-unit

# Integration tests coverage only
make coverage-integration

# Check coverage thresholds
make coverage-check

# Generate JSON coverage report
make coverage-check-json

# Generate coverage badge
make coverage-badge

# CI coverage check with strict thresholds
make coverage-ci

# Clean coverage artifacts
make coverage-clean
```

Coverage reports are automatically generated in CI/CD and uploaded to [Codecov](https://codecov.io/gh/matanbaruch/cursor-admin-api-exporter). See [COVERAGE.md](COVERAGE.md) for detailed coverage documentation.

### Local Development

```bash
# Run in development mode
make dev

# Or with environment variables
export CURSOR_API_TOKEN=your_token_here
export LOG_LEVEL=debug
go run main.go
```

## Security

### Authentication

The exporter uses Bearer token authentication with the Cursor Admin API. Ensure your token is:

- Stored securely (use secrets management)
- Rotated regularly
- Has minimal required permissions

### Container Security

Our container images are:

- Built with minimal Alpine base images
- Run as non-root user (uid: 65534)
- Include security scanning
- Signed with build provenance attestations

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- **Issues**: [GitHub Issues](https://github.com/matanbaruch/cursor-admin-api-exporter/issues)
- **Discussions**: [GitHub Discussions](https://github.com/matanbaruch/cursor-admin-api-exporter/discussions)

---

**Note**: This project is not officially affiliated with Cursor. It's a community-driven monitoring solution for Cursor Admin API.
