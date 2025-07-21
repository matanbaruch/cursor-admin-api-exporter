# Metrics Reference

This document provides a comprehensive reference for all metrics exported by the Cursor Admin API Exporter.

## Overview

The exporter provides metrics across several categories:
- **Team Members**: Team composition and roles
- **Daily Usage**: Code activity and AI assistance
- **Spending**: Cost and billing information
- **Usage Events**: Detailed usage tracking
- **Exporter Health**: System performance metrics

## Team Members Metrics

### `cursor_team_members_total`
- **Type**: Gauge
- **Description**: Total number of team members
- **Labels**: None

```prometheus
# HELP cursor_team_members_total Total number of team members
# TYPE cursor_team_members_total gauge
cursor_team_members_total 15
```

### `cursor_team_members_by_role`
- **Type**: Gauge
- **Description**: Number of team members by role
- **Labels**: `role`

```prometheus
# HELP cursor_team_members_by_role Number of team members by role
# TYPE cursor_team_members_by_role gauge
cursor_team_members_by_role{role="admin"} 3
cursor_team_members_by_role{role="member"} 10
cursor_team_members_by_role{role="viewer"} 2
```

## Daily Usage Metrics

### `cursor_daily_lines_added_total`
- **Type**: Gauge
- **Description**: Total lines of code added per day
- **Labels**: `date`

```prometheus
# HELP cursor_daily_lines_added_total Total lines of code added per day
# TYPE cursor_daily_lines_added_total gauge
cursor_daily_lines_added_total{date="2024-01-20"} 1250
```

### `cursor_daily_lines_deleted_total`
- **Type**: Gauge
- **Description**: Total lines of code deleted per day
- **Labels**: `date`

```prometheus
# HELP cursor_daily_lines_deleted_total Total lines of code deleted per day
# TYPE cursor_daily_lines_deleted_total gauge
cursor_daily_lines_deleted_total{date="2024-01-20"} 430
```

### `cursor_daily_suggestion_acceptance_rate`
- **Type**: Gauge
- **Description**: AI suggestion acceptance rate per day (0-1)
- **Labels**: `date`

```prometheus
# HELP cursor_daily_suggestion_acceptance_rate AI suggestion acceptance rate per day
# TYPE cursor_daily_suggestion_acceptance_rate gauge
cursor_daily_suggestion_acceptance_rate{date="2024-01-20"} 0.75
```

### `cursor_daily_tabs_used_total`
- **Type**: Gauge
- **Description**: Total tab completions used per day
- **Labels**: `date`

```prometheus
# HELP cursor_daily_tabs_used_total Total tab completions used per day
# TYPE cursor_daily_tabs_used_total gauge
cursor_daily_tabs_used_total{date="2024-01-20"} 1850
```

### `cursor_daily_composer_used_total`
- **Type**: Gauge
- **Description**: Total composer usage per day
- **Labels**: `date`

```prometheus
# HELP cursor_daily_composer_used_total Total composer usage per day
# TYPE cursor_daily_composer_used_total gauge
cursor_daily_composer_used_total{date="2024-01-20"} 45
```

### `cursor_daily_chat_requests_total`
- **Type**: Gauge
- **Description**: Total chat requests per day
- **Labels**: `date`

```prometheus
# HELP cursor_daily_chat_requests_total Total chat requests per day
# TYPE cursor_daily_chat_requests_total gauge
cursor_daily_chat_requests_total{date="2024-01-20"} 120
```

### `cursor_daily_model_usage`
- **Type**: Gauge
- **Description**: Most used model per day (1 if this model was most used, 0 otherwise)
- **Labels**: `date`, `model`

```prometheus
# HELP cursor_daily_model_usage Most used model per day
# TYPE cursor_daily_model_usage gauge
cursor_daily_model_usage{date="2024-01-20",model="gpt-4"} 1
cursor_daily_model_usage{date="2024-01-20",model="claude-3"} 0
```

### `cursor_daily_extension_usage`
- **Type**: Gauge
- **Description**: Most used extension per day (1 if this extension was most used, 0 otherwise)
- **Labels**: `date`, `extension`

```prometheus
# HELP cursor_daily_extension_usage Most used extension per day
# TYPE cursor_daily_extension_usage gauge
cursor_daily_extension_usage{date="2024-01-20",extension="python"} 1
cursor_daily_extension_usage{date="2024-01-20",extension="javascript"} 0
```

## Spending Metrics

### `cursor_spending_total_cents`
- **Type**: Gauge
- **Description**: Total spending in cents
- **Labels**: None

```prometheus
# HELP cursor_spending_total_cents Total spending in cents
# TYPE cursor_spending_total_cents gauge
cursor_spending_total_cents 125000
```

### `cursor_spending_by_member_cents`
- **Type**: Gauge
- **Description**: Spending by team member in cents
- **Labels**: `member_email`, `date`

```prometheus
# HELP cursor_spending_by_member_cents Spending by team member in cents
# TYPE cursor_spending_by_member_cents gauge
cursor_spending_by_member_cents{member_email="john@example.com",date="2024-01-20"} 8500
cursor_spending_by_member_cents{member_email="jane@example.com",date="2024-01-20"} 12000
```

### `cursor_premium_requests_by_member_total`
- **Type**: Gauge
- **Description**: Premium requests by team member
- **Labels**: `member_email`, `date`

```prometheus
# HELP cursor_premium_requests_by_member_total Premium requests by team member
# TYPE cursor_premium_requests_by_member_total gauge
cursor_premium_requests_by_member_total{member_email="john@example.com",date="2024-01-20"} 45
cursor_premium_requests_by_member_total{member_email="jane@example.com",date="2024-01-20"} 67
```

### `cursor_premium_requests_total`
- **Type**: Gauge
- **Description**: Total premium requests
- **Labels**: None

```prometheus
# HELP cursor_premium_requests_total Total premium requests
# TYPE cursor_premium_requests_total gauge
cursor_premium_requests_total 450
```

## Usage Events Metrics

### `cursor_usage_events_total`
- **Type**: Gauge
- **Description**: Total number of usage events
- **Labels**: None

```prometheus
# HELP cursor_usage_events_total Total number of usage events
# TYPE cursor_usage_events_total gauge
cursor_usage_events_total 2850
```

### `cursor_usage_events_by_type_total`
- **Type**: Gauge
- **Description**: Number of usage events by type
- **Labels**: `event_type`

```prometheus
# HELP cursor_usage_events_by_type_total Number of usage events by type
# TYPE cursor_usage_events_by_type_total gauge
cursor_usage_events_by_type_total{event_type="completion"} 1200
cursor_usage_events_by_type_total{event_type="chat"} 450
cursor_usage_events_by_type_total{event_type="edit"} 800
```

### `cursor_usage_events_by_user_total`
- **Type**: Gauge
- **Description**: Number of usage events by user
- **Labels**: `user_email`

```prometheus
# HELP cursor_usage_events_by_user_total Number of usage events by user
# TYPE cursor_usage_events_by_user_total gauge
cursor_usage_events_by_user_total{user_email="john@example.com"} 350
cursor_usage_events_by_user_total{user_email="jane@example.com"} 425
```

### `cursor_usage_events_by_model_total`
- **Type**: Gauge
- **Description**: Number of usage events by model
- **Labels**: `model`

```prometheus
# HELP cursor_usage_events_by_model_total Number of usage events by model
# TYPE cursor_usage_events_by_model_total gauge
cursor_usage_events_by_model_total{model="gpt-4"} 1850
cursor_usage_events_by_model_total{model="claude-3"} 650
cursor_usage_events_by_model_total{model="gpt-3.5-turbo"} 350
```

### `cursor_tokens_consumed_total`
- **Type**: Gauge
- **Description**: Total tokens consumed
- **Labels**: None

```prometheus
# HELP cursor_tokens_consumed_total Total tokens consumed
# TYPE cursor_tokens_consumed_total gauge
cursor_tokens_consumed_total 1250000
```

### `cursor_tokens_consumed_by_model_total`
- **Type**: Gauge
- **Description**: Tokens consumed by model
- **Labels**: `model`

```prometheus
# HELP cursor_tokens_consumed_by_model_total Tokens consumed by model
# TYPE cursor_tokens_consumed_by_model_total gauge
cursor_tokens_consumed_by_model_total{model="gpt-4"} 850000
cursor_tokens_consumed_by_model_total{model="claude-3"} 300000
cursor_tokens_consumed_by_model_total{model="gpt-3.5-turbo"} 100000
```

### `cursor_tokens_consumed_by_user_total`
- **Type**: Gauge
- **Description**: Tokens consumed by user
- **Labels**: `user_email`

```prometheus
# HELP cursor_tokens_consumed_by_user_total Tokens consumed by user
# TYPE cursor_tokens_consumed_by_user_total gauge
cursor_tokens_consumed_by_user_total{user_email="john@example.com"} 125000
cursor_tokens_consumed_by_user_total{user_email="jane@example.com"} 180000
```

## Exporter Health Metrics

### `cursor_exporter_scrape_duration_seconds`
- **Type**: Histogram
- **Description**: Time spent scraping the Cursor Admin API
- **Labels**: None

```prometheus
# HELP cursor_exporter_scrape_duration_seconds Time spent scraping Cursor Admin API
# TYPE cursor_exporter_scrape_duration_seconds histogram
cursor_exporter_scrape_duration_seconds_bucket{le="0.1"} 0
cursor_exporter_scrape_duration_seconds_bucket{le="0.5"} 12
cursor_exporter_scrape_duration_seconds_bucket{le="1"} 45
cursor_exporter_scrape_duration_seconds_bucket{le="2.5"} 48
cursor_exporter_scrape_duration_seconds_bucket{le="5"} 48
cursor_exporter_scrape_duration_seconds_bucket{le="10"} 48
cursor_exporter_scrape_duration_seconds_bucket{le="+Inf"} 48
cursor_exporter_scrape_duration_seconds_sum 18.5
cursor_exporter_scrape_duration_seconds_count 48
```

### `cursor_exporter_scrape_errors_total`
- **Type**: Counter
- **Description**: Total number of scrape errors
- **Labels**: None

```prometheus
# HELP cursor_exporter_scrape_errors_total Total number of scrape errors
# TYPE cursor_exporter_scrape_errors_total counter
cursor_exporter_scrape_errors_total 2
```

## Metric Labels

### Common Labels

| Label | Description | Example Values |
|-------|-------------|----------------|
| `date` | Date in YYYY-MM-DD format | `2024-01-20` |
| `member_email` | Team member email address | `john@example.com` |
| `user_email` | User email address | `jane@example.com` |
| `role` | Team member role | `admin`, `member`, `viewer` |
| `model` | AI model name | `gpt-4`, `claude-3`, `gpt-3.5-turbo` |
| `extension` | File extension | `python`, `javascript`, `go` |
| `event_type` | Type of usage event | `completion`, `chat`, `edit` |

## Metric Collection

### Collection Frequency

- **Default**: Every 30 seconds (configurable via Prometheus)
- **Recommended**: 30-60 seconds for production
- **Minimum**: 15 seconds (to avoid API rate limits)

### Data Freshness

- **Team Members**: Real-time
- **Daily Usage**: Updated daily at midnight UTC
- **Spending**: Updated hourly
- **Usage Events**: Updated every 5 minutes

## Querying Examples

### PromQL Queries

#### Team Productivity
```prometheus
# Average suggestion acceptance rate over 7 days
avg_over_time(cursor_daily_suggestion_acceptance_rate[7d])

# Total lines of code added per week
sum(increase(cursor_daily_lines_added_total[7d]))

# Most active users by token consumption
topk(5, cursor_tokens_consumed_by_user_total)
```

#### Cost Analysis
```prometheus
# Daily spending trend
cursor_spending_total_cents / 100

# Top spenders by team member
topk(10, cursor_spending_by_member_cents / 100)

# Premium request utilization
cursor_premium_requests_total / cursor_team_members_total
```

#### Model Usage
```prometheus
# Most used AI models
topk(5, cursor_tokens_consumed_by_model_total)

# Model usage distribution
(cursor_tokens_consumed_by_model_total / cursor_tokens_consumed_total) * 100
```

### Grafana Queries

#### Dashboard Variables
```prometheus
# Team members for filtering
label_values(cursor_usage_events_by_user_total, user_email)

# Available models
label_values(cursor_tokens_consumed_by_model_total, model)

# Date range for daily metrics
label_values(cursor_daily_lines_added_total, date)
```

#### Panels
```prometheus
# Lines of code trend
cursor_daily_lines_added_total - cursor_daily_lines_deleted_total

# Suggestion acceptance rate trend
cursor_daily_suggestion_acceptance_rate * 100

# Cost per team member
cursor_spending_total_cents / cursor_team_members_total / 100
```

## Alerting Rules

### Suggested Alerts

```yaml
# Low suggestion acceptance rate
- alert: CursorLowSuggestionAcceptance
  expr: cursor_daily_suggestion_acceptance_rate < 0.3
  for: 1h
  labels:
    severity: warning
  annotations:
    summary: "Low AI suggestion acceptance rate"

# High spending
- alert: CursorHighSpending
  expr: increase(cursor_spending_total_cents[24h]) > 10000
  for: 15m
  labels:
    severity: critical
  annotations:
    summary: "High daily spending detected"

# API scrape errors
- alert: CursorExporterErrors
  expr: rate(cursor_exporter_scrape_errors_total[5m]) > 0.1
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: "Cursor exporter experiencing errors"
```

## Troubleshooting Metrics

### Missing Metrics

1. **Check API connectivity**: `cursor_exporter_scrape_errors_total`
2. **Verify authentication**: Look for 401 errors in logs
3. **Check rate limits**: High scrape error rates
4. **Validate permissions**: Ensure API token has admin access

### Incorrect Values

1. **Time zone issues**: Daily metrics use UTC
2. **API delays**: Some metrics may have delays
3. **Caching**: API responses may be cached
4. **Permissions**: Limited data based on token permissions

---

**Related Documentation:**
- [Configuration Guide](configuration.md)
- [Grafana Dashboard](grafana-dashboard.md)
- [Troubleshooting Guide](troubleshooting.md)