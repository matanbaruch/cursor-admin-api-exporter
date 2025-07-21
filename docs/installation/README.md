# Installation Guide

This guide provides detailed installation instructions for all supported platforms and deployment methods.

## Installation Methods

Choose the installation method that best fits your environment:

| Method | Best For | Difficulty | Features |
|--------|----------|------------|----------|
| [Docker](docker.md) | Development, Testing | Easy | Quick setup, isolated |
| [Docker Compose](docker-compose.md) | Local Development | Easy | Full stack with monitoring |
| [Kubernetes/Helm](helm.md) | Production | Medium | Scalable, enterprise-ready |
| [Binary](binary.md) | Simple deployments | Easy | Minimal dependencies |
| [Systemd](systemd.md) | Linux servers | Medium | System service integration |

## Quick Start

For the fastest setup, use Docker:

```bash
docker run -d \
  --name cursor-exporter \
  -p 8080:8080 \
  -e CURSOR_API_TOKEN=your_token_here \
  ghcr.io/yourusername/cursor-admin-api-exporter:latest
```

## Platform Support

### Operating Systems

| OS | Binary | Docker | Kubernetes |
|----|--------|--------|------------|
| Linux (AMD64) | ✅ | ✅ | ✅ |
| Linux (ARM64) | ✅ | ✅ | ✅ |
| macOS (Intel) | ✅ | ✅ | ✅ |
| macOS (ARM64) | ✅ | ✅ | ✅ |
| Windows | ✅ | ✅ | ✅ |

### Container Platforms

| Platform | Status | Notes |
|----------|--------|-------|
| Docker | ✅ | Recommended |
| Podman | ✅ | Compatible |
| containerd | ✅ | Via Kubernetes |
| CRI-O | ✅ | Via Kubernetes |

### Kubernetes Distributions

| Distribution | Status | Notes |
|-------------|--------|-------|
| Vanilla Kubernetes | ✅ | Fully supported |
| Amazon EKS | ✅ | Tested |
| Google GKE | ✅ | Tested |
| Azure AKS | ✅ | Tested |
| Red Hat OpenShift | ✅ | Compatible |
| Rancher | ✅ | Compatible |
| k3s | ✅ | Compatible |
| MicroK8s | ✅ | Compatible |

## Prerequisites

### General Requirements

- **API Token**: Cursor Admin API token with appropriate permissions
- **Network Access**: Outbound HTTPS access to `api.cursor.com`
- **Resources**: Minimal (50MB RAM, negligible CPU)

### Platform-Specific Requirements

#### Docker
- Docker Engine 20.10+ or Docker Desktop
- 50MB available disk space

#### Kubernetes
- Kubernetes 1.19+
- Helm 3.0+ (for Helm installation)
- RBAC enabled cluster

#### Binary
- No additional dependencies
- Appropriate architecture binary

## Installation Steps

### 1. Get API Token

Before installation, obtain your Cursor Admin API token:

1. Log in to your Cursor team dashboard
2. Navigate to **Settings** → **API Keys**
3. Create a new API key with **Admin** permissions
4. Copy the token for use in configuration

### 2. Choose Installation Method

Select your preferred installation method:

- **New to containers?** → [Docker](docker.md)
- **Want full monitoring stack?** → [Docker Compose](docker-compose.md)
- **Production deployment?** → [Kubernetes/Helm](helm.md)
- **Simple server deployment?** → [Binary](binary.md)
- **Linux system service?** → [Systemd](systemd.md)

### 3. Configure and Deploy

Follow the specific guide for your chosen method:

- [Docker Installation](docker.md)
- [Docker Compose Installation](docker-compose.md)
- [Helm Installation](helm.md)
- [Binary Installation](binary.md)
- [Systemd Installation](systemd.md)

### 4. Verify Installation

After installation, verify the exporter is working:

```bash
# Check health
curl http://localhost:8080/health

# Check metrics
curl http://localhost:8080/metrics
```

## Next Steps

After installation:

1. **Configure Monitoring**: Set up Prometheus and Grafana
2. **Set Up Alerts**: Configure alerting rules
3. **Security Review**: Implement security best practices
4. **Performance Tuning**: Optimize for your environment

## Need Help?

If you encounter issues during installation:

1. Check the [Troubleshooting Guide](../troubleshooting.md)
2. Review the [Configuration Guide](../configuration.md)
3. Search [GitHub Issues](https://github.com/yourusername/cursor-admin-api-exporter/issues)
4. Ask in [GitHub Discussions](https://github.com/yourusername/cursor-admin-api-exporter/discussions)

---

**Choose your installation method:**
- [Docker](docker.md) - Quick and easy
- [Docker Compose](docker-compose.md) - Full development stack
- [Kubernetes/Helm](helm.md) - Production deployment
- [Binary](binary.md) - Simple and lightweight
- [Systemd](systemd.md) - Linux system service