# Cursor Admin API Exporter Helm Chart

This Helm chart deploys the Cursor Admin API Exporter, a Prometheus exporter that provides comprehensive metrics about your Cursor Admin API including team members, daily usage, spending, and usage events.

## Features

- Comprehensive Cursor Admin API metrics collection
- Secure API token handling via Kubernetes secrets
- External Secret Operator integration for enterprise secret management
- Optional Prometheus Operator integration with ServiceMonitor
- Configurable resource limits and requests
- Health checks and readiness probes
- Optional ingress for external access
- Horizontal Pod Autoscaler support
- Pod Disruption Budget for high availability
- Security hardening with non-root containers and read-only filesystem

## Prerequisites

- Kubernetes 1.19+
- Helm 3.2.0+
- A valid Cursor Admin API token
- **(Optional)** External Secret Operator for external secret management

## Installing the Chart

### Install from OCI Registry

The chart is available as an OCI artifact. Since it's hosted in an OCI registry, there's no need to add a Helm repository.

```bash
# Install with specific version (recommended)
helm install my-cursor-admin-api-exporter \
  oci://ghcr.io/matanbaruch/cursor-admin-api-exporter/charts/cursor-admin-api-exporter \
  --version 0.1.1 \
  --set cursor.apiToken="your-cursor-api-token"

# Install latest version
helm install my-cursor-admin-api-exporter \
  oci://ghcr.io/matanbaruch/cursor-admin-api-exporter/charts/cursor-admin-api-exporter \
  --set cursor.apiToken="your-cursor-api-token"
```

> **Note on OCI Registry**: This chart is distributed via OCI (Open Container Initiative) registry instead of traditional Helm repositories. This provides better security, versioning, and integration with container registries. No `helm repo add` command is needed - you can install directly using the `oci://` URL with specific versions.

### Install from local directory

```bash
# Clone the repository
git clone https://github.com/matanbaruch/cursor-admin-api-exporter.git
cd cursor-admin-api-exporter

# Install the chart
helm install my-cursor-admin-api-exporter ./charts/cursor-admin-api-exporter \
  --set cursor.apiToken="your-cursor-api-token"
```

### Install with custom values

```bash
# From OCI registry
helm install my-cursor-admin-api-exporter \
  oci://ghcr.io/matanbaruch/cursor-admin-api-exporter/charts/cursor-admin-api-exporter \
  --version 0.1.1 \
  --values my-values.yaml

# From local directory
helm install my-cursor-admin-api-exporter ./charts/cursor-admin-api-exporter \
  --values my-values.yaml
```

### Install with ExternalSecret

```bash
# First, create a SecretStore (example with AWS Secrets Manager)
kubectl apply -f - <<EOF
apiVersion: external-secrets.io/v1beta1
kind: SecretStore
metadata:
  name: aws-secrets-manager
spec:
  provider:
    aws:
      service: SecretsManager
      region: us-west-2
      auth:
        jwt:
          serviceAccountRef:
            name: external-secrets-sa
EOF

# Install with ExternalSecret
helm install cursor-admin-api-exporter \
  oci://ghcr.io/matanbaruch/cursor-admin-api-exporter/charts/cursor-admin-api-exporter \
  --version 0.1.1 \
  --set externalSecret.enabled=true \
  --set externalSecret.secretStoreRef.name=aws-secrets-manager \
  --set-json 'externalSecret.data=[{"secretKey":"cursor-api-token","remoteRef":{"key":"cursor/api-token"}}]'

# Alternative: Install from local directory
# helm install cursor-admin-api-exporter ./charts/cursor-admin-api-exporter \
#   --set externalSecret.enabled=true \
#   --set externalSecret.secretStoreRef.name=aws-secrets-manager \
#   --set-json 'externalSecret.data=[{"secretKey":"cursor-api-token","remoteRef":{"key":"cursor/api-token"}}]'
```

## Configuration

The following table lists the configurable parameters and their default values.

### Basic Configuration

| Parameter          | Description        | Default                        |
| ------------------ | ------------------ | ------------------------------ |
| `replicaCount`     | Number of replicas | `1`                            |
| `image.repository` | Image repository   | `cursor-admin-api-exporter` |
| `image.pullPolicy` | Image pull policy  | `IfNotPresent`                 |
| `image.tag`        | Image tag          | `"latest"`                     |

### Cursor Configuration

| Parameter          | Description                            | Default                            |
| ------------------ | -------------------------------------- | ---------------------------------- |
| `cursor.apiUrl`    | Cursor Admin API URL                   | `"https://api.cursor.com"`         |
| `cursor.apiToken`  | Cursor Admin API token (stored in secret) | `""`                           |

### External Secret Configuration

| Parameter                            | Description                                           | Default        |
| ------------------------------------ | ----------------------------------------------------- | -------------- |
| `externalSecret.enabled`             | Enable External Secret Operator integration           | `false`        |
| `externalSecret.secretStoreRef.name` | Secret Store reference name                           | `""`           |
| `externalSecret.secretStoreRef.kind` | Secret Store kind (SecretStore or ClusterSecretStore) | `SecretStore`  |
| `externalSecret.secretName`          | Target secret name (optional)                         | Auto-generated |
| `externalSecret.data`                | Remote references for the external secret             | `[]`           |
| `externalSecret.refreshInterval`     | Refresh interval for the external secret              | `"1h"`         |
| `externalSecret.annotations`         | Additional annotations for ExternalSecret resource    | `{}`           |
| `externalSecret.labels`              | Additional labels for ExternalSecret resource         | `{}`           |

### Application Configuration

| Parameter              | Description             | Default      |
| ---------------------- | ----------------------- | ------------ |
| `config.listenAddress` | Listen address and port | `":8080"`    |
| `config.metricsPath`   | Metrics endpoint path   | `"/metrics"` |
| `config.logLevel`      | Log level               | `"info"`     |

## Usage Examples

### Basic Installation

```bash
# From OCI registry (recommended)
helm install cursor-admin-api-exporter \
  oci://ghcr.io/matanbaruch/cursor-admin-api-exporter/charts/cursor-admin-api-exporter \
  --version 0.1.1 \
  --set cursor.apiToken="cursor_token_xxx"

# From local directory
helm install cursor-admin-api-exporter ./charts/cursor-admin-api-exporter \
  --set cursor.apiToken="cursor_token_xxx"
```

### With Prometheus Operator

```bash
# From OCI registry (recommended)
helm install cursor-admin-api-exporter \
  oci://ghcr.io/matanbaruch/cursor-admin-api-exporter/charts/cursor-admin-api-exporter \
  --version 0.1.1 \
  --set cursor.apiToken="cursor_token_xxx" \
  --set serviceMonitor.enabled=true \
  --set serviceMonitor.additionalLabels.release=prometheus

# From local directory
helm install cursor-admin-api-exporter ./charts/cursor-admin-api-exporter \
  --set cursor.apiToken="cursor_token_xxx" \
  --set serviceMonitor.enabled=true \
  --set serviceMonitor.additionalLabels.release=prometheus
```

### Production Configuration

```bash
# From OCI registry (recommended)
helm install cursor-admin-api-exporter \
  oci://ghcr.io/matanbaruch/cursor-admin-api-exporter/charts/cursor-admin-api-exporter \
  --version 0.1.1 \
  --values ./values-production.yaml

# From local directory
helm install cursor-admin-api-exporter ./charts/cursor-admin-api-exporter \
  --values ./charts/cursor-admin-api-exporter/values-production.yaml
```

## Getting Your Cursor Admin API Token

1. Log into your Cursor Admin dashboard
2. Go to **Settings** → **API Keys**
3. Create a new API key with appropriate permissions
4. Copy the token and use it as `cursor.apiToken`

## Monitoring

Once deployed, the exporter provides metrics at the `/metrics` endpoint. Key metrics include:

- `cursor_team_members` - Total number of team members
- `cursor_daily_usage` - Daily usage metrics
- `cursor_spending` - Spending metrics
- `cursor_usage_events` - Usage events metrics

## Testing

Run the included test:

```bash
helm test cursor-admin-api-exporter
```

## Troubleshooting

### Check pod logs

```bash
kubectl logs -l app.kubernetes.io/name=cursor-admin-api-exporter
```

### Test connectivity

```bash
kubectl port-forward svc/cursor-admin-api-exporter 8080:8080
curl http://localhost:8080/health
curl http://localhost:8080/metrics
```