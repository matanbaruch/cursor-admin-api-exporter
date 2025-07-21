---
layout: default
title: Home
nav_order: 1
description: "Cursor Admin API Exporter - A Prometheus exporter for Cursor Admin API"
permalink: /
---

# Cursor Admin API Exporter

A Prometheus exporter for Cursor Admin API that provides comprehensive metrics about your Cursor usage, team members, and spending.

## Features

- **Team Metrics**: Monitor team member activity and usage
- **Usage Tracking**: Track daily usage events and patterns
- **Spending Insights**: Monitor costs and spending trends
- **Prometheus Integration**: Export metrics in Prometheus format
- **Docker Support**: Easy deployment with Docker containers
- **Kubernetes Ready**: Helm charts for Kubernetes deployment

## Quick Start

```bash
# Using Docker
docker run -p 8080:8080 cursor-admin-api-exporter

# Using Go
go run main.go
```

## Documentation

- [Getting Started](getting-started.md)
- [Configuration](configuration.md)
- [Installation](installation/)
- [Metrics](metrics.md)
- [Troubleshooting](troubleshooting.md)

## Contributing

We welcome contributions! Please see our [Contributing Guide](../CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.