# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.5] - 2025-07-21

## [0.1.4] - 2025-07-21

## [0.1.3] - 2025-07-21

## [0.1.2] - 2025-07-21

## [0.1.1] - 2025-07-21

## [0.1.0] - 2024-01-20

### Added
- Initial release of Cursor Admin API Exporter
- Team members metrics collection and export
- Daily usage metrics including lines of code, suggestion acceptance rates
- Spending analytics with per-member breakdown
- Usage events tracking with token consumption metrics
- Comprehensive test suite with unit, integration, and performance tests
- Docker container support with multi-architecture builds
- Kubernetes Helm chart with full configuration options
- GitHub Actions CI/CD pipeline with automated releases
- Security features including attestations and provenance
- Detailed documentation and examples
- Grafana dashboard template
- Prometheus configuration examples
- Health check and graceful shutdown support

### Features
- **Metrics Collection**: Complete coverage of Cursor Admin API endpoints
- **Multi-platform Support**: Linux, macOS, and Windows binaries
- **Container Security**: Non-root user, minimal base image, security scanning
- **Kubernetes Ready**: Helm chart with ServiceMonitor, HPA, and PDB support
- **Monitoring**: Built-in metrics for exporter health and performance
- **Configuration**: Flexible environment variable configuration
- **Testing**: Comprehensive test coverage with automated coverage reporting

### Security
- All container images signed with build provenance attestations
- Security scanning integrated into CI/CD pipeline
- Non-root container execution
- Minimal attack surface with Alpine base images

[0.1.1]: https://github.com/matanbaruch/cursor-admin-api-exporter/compare/v0.1.0...v0.1.1
[0.1.2]: https://github.com/matanbaruch/cursor-admin-api-exporter/compare/v0.1.1...v0.1.2
[0.1.3]: https://github.com/matanbaruch/cursor-admin-api-exporter/compare/v0.1.2...v0.1.3
[0.1.4]: https://github.com/matanbaruch/cursor-admin-api-exporter/compare/v0.1.3...v0.1.4
[0.1.5]: https://github.com/matanbaruch/cursor-admin-api-exporter/compare/v0.1.4...v0.1.5
[Unreleased]: https://github.com/matanbaruch/cursor-admin-api-exporter/compare/v0.1.5...HEAD
[0.1.0]: https://github.com/matanbaruch/cursor-admin-api-exporter/releases/tag/v0.1.0