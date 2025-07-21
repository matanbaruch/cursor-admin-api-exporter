# Getting Started

This guide will help you get the Cursor Admin API Exporter up and running quickly.

## Prerequisites

Before you begin, ensure you have:

- A Cursor team with admin access
- A Cursor Admin API token
- Docker, Kubernetes, or Go installed (depending on your deployment method)

## Step 1: Get Your API Token

1. Log in to your Cursor team dashboard
2. Navigate to **Settings** → **API Keys**
3. Click **Create New API Key**
4. Select **Admin** permissions
5. Copy the generated token (you'll need this for configuration)

⚠️ **Important**: Store this token securely and never commit it to version control.

## Step 2: Choose Your Deployment Method

### Option A: Docker (Recommended for Testing)

```bash
# Pull the latest image
docker pull ghcr.io/matanbaruch/cursor-admin-api-exporter:latest

# Run the exporter
docker run -d \
  --name cursor-exporter \
  -p 8080:8080 \
  -e CURSOR_API_TOKEN=your_token_here \
  -e CURSOR_API_URL=https://api.cursor.com \
  -e LOG_LEVEL=info \
  ghcr.io/matanbaruch/cursor-admin-api-exporter:latest
```

### Option B: Docker Compose (Recommended for Development)

```bash
# Clone the repository
git clone https://github.com/matanbaruch/cursor-admin-api-exporter.git
cd cursor-admin-api-exporter

# Copy and edit environment file
cp env.example .env
# Edit .env file with your API token

# Start services (includes Prometheus and Grafana)
docker-compose up -d
```

### Option C: Kubernetes with Helm (Recommended for Production)

```bash
# Install the chart
helm install cursor-admin-api-exporter \
  oci://ghcr.io/matanbaruch/cursor-admin-api-exporter/charts/cursor-admin-api-exporter \
  --set cursor.apiToken=your_token_here \
  --set cursor.apiUrl=https://api.cursor.com
```

### Option D: Binary Installation

```bash
# Download the latest binary (Linux AMD64 example)
wget https://github.com/matanbaruch/cursor-admin-api-exporter/releases/latest/download/cursor-admin-api-exporter-linux-amd64

# Make it executable
chmod +x cursor-admin-api-exporter-linux-amd64

# Run the exporter
export CURSOR_API_TOKEN=your_token_here
./cursor-admin-api-exporter-linux-amd64
```

## Step 3: Verify the Installation

1. **Check the health endpoint**:
   ```bash
   curl http://localhost:8080/health
   ```

2. **View available metrics**:
   ```bash
   curl http://localhost:8080/metrics
   ```

3. **Check the web interface**:
   Open http://localhost:8080 in your browser

## Step 4: Configure Monitoring

### Prometheus Configuration

Add this to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'cursor-admin-api-exporter'
    static_configs:
      - targets: ['localhost:8080']
    scrape_interval: 30s
    metrics_path: /metrics
```

### Grafana Dashboard

1. Import the dashboard from `grafana-dashboard.json`
2. Or access the pre-built dashboard at http://localhost:3000 (if using docker-compose)
3. Default credentials: admin/admin

## Step 5: Explore the Metrics

The exporter provides several categories of metrics:

### Team Metrics
- `cursor_team_members_total` - Total team members
- `cursor_team_members_by_role` - Members by role

### Usage Metrics
- `cursor_daily_lines_added_total` - Lines of code added
- `cursor_daily_suggestion_acceptance_rate` - AI suggestion acceptance
- `cursor_daily_tabs_used_total` - Tab completions used
- `cursor_daily_chat_requests_total` - Chat requests made

### Spending Metrics
- `cursor_spending_total_cents` - Total spending
- `cursor_spending_by_member_cents` - Per-member spending
- `cursor_premium_requests_total` - Premium requests

### Usage Events
- `cursor_usage_events_total` - Total usage events
- `cursor_tokens_consumed_total` - Total tokens consumed
- `cursor_tokens_consumed_by_model_total` - Tokens by model

## Common Configuration Options

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `CURSOR_API_TOKEN` | - | **Required** - Your Cursor API token |
| `CURSOR_API_URL` | `https://api.cursor.com` | Cursor API endpoint |
| `LISTEN_ADDRESS` | `:8080` | HTTP server address |
| `METRICS_PATH` | `/metrics` | Metrics endpoint path |
| `LOG_LEVEL` | `info` | Logging level (debug, info, warn, error) |

### Example Configuration

```bash
# Basic configuration
export CURSOR_API_TOKEN="your_token_here"
export CURSOR_API_URL="https://api.cursor.com"
export LISTEN_ADDRESS=":8080"
export METRICS_PATH="/metrics"
export LOG_LEVEL="info"

# Run the exporter
./cursor-admin-api-exporter
```

## Next Steps

Now that you have the exporter running:

1. **Configure Monitoring**: Set up Prometheus and Grafana dashboards
2. **Set Up Alerts**: Configure alerts for important metrics
3. **Customize Configuration**: Adjust settings for your environment
4. **Scale for Production**: Deploy with high availability if needed

## Troubleshooting

### Common Issues

**Connection refused or timeout errors**:
- Check if the Cursor API is accessible from your network
- Verify your API token is correct and has admin permissions
- Check firewall settings

**Authentication errors**:
- Ensure your API token is valid and not expired
- Verify you have admin permissions on the Cursor team
- Check the token is properly set in environment variables

**No metrics appearing**:
- Check the logs for error messages: `docker logs cursor-exporter`
- Verify the metrics endpoint is accessible: `curl http://localhost:8080/metrics`
- Ensure Prometheus is scraping the correct endpoint

**Permission denied errors**:
- Check if running as non-root user (security best practice)
- Verify file permissions for configuration files
- Ensure proper container security context

### Getting Help

If you encounter issues:

1. Check the [Troubleshooting Guide](troubleshooting.md)
2. Review the [Configuration Guide](configuration.md)
3. Search [GitHub Issues](https://github.com/matanbaruch/cursor-admin-api-exporter/issues)
4. Ask questions in [GitHub Discussions](https://github.com/matanbaruch/cursor-admin-api-exporter/discussions)

## Security Best Practices

- **Never commit API tokens** to version control
- **Use environment variables** or secret management systems
- **Rotate tokens regularly** 
- **Run with minimal permissions**
- **Use HTTPS** for all communications
- **Monitor for unusual activity**

## Performance Considerations

- **Scrape interval**: Default 30s is recommended
- **API rate limits**: Monitor for rate limiting
- **Resource usage**: ~50MB RAM, minimal CPU
- **Network**: Dependent on API response sizes

---

**Next**: [Installation Guide](installation/) for detailed deployment instructions