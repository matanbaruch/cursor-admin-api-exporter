package exporters

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/matanbaruch/cursor-admin-api-exporter/pkg/client"
)

type UsageEventsExporter struct {
	client *client.CursorClient

	totalEvents           *prometheus.Desc
	eventsByType          *prometheus.Desc
	eventsByUser          *prometheus.Desc
	eventsByModel         *prometheus.Desc
	tokensConsumed        *prometheus.Desc
	tokensConsumedByModel *prometheus.Desc
	tokensConsumedByUser  *prometheus.Desc
}

func NewUsageEventsExporter(client *client.CursorClient) *UsageEventsExporter {
	return &UsageEventsExporter{
		client: client,

		totalEvents: prometheus.NewDesc(
			"cursor_usage_events_total",
			"Total number of usage events",
			nil,
			nil,
		),

		eventsByType: prometheus.NewDesc(
			"cursor_usage_events_by_type_total",
			"Number of usage events by type",
			[]string{"event_type"},
			nil,
		),

		eventsByUser: prometheus.NewDesc(
			"cursor_usage_events_by_user_total",
			"Number of usage events by user",
			[]string{"user_email"},
			nil,
		),

		eventsByModel: prometheus.NewDesc(
			"cursor_usage_events_by_model_total",
			"Number of usage events by model",
			[]string{"model"},
			nil,
		),

		tokensConsumed: prometheus.NewDesc(
			"cursor_tokens_consumed_total",
			"Total tokens consumed",
			nil,
			nil,
		),

		tokensConsumedByModel: prometheus.NewDesc(
			"cursor_tokens_consumed_by_model_total",
			"Tokens consumed by model",
			[]string{"model"},
			nil,
		),

		tokensConsumedByUser: prometheus.NewDesc(
			"cursor_tokens_consumed_by_user_total",
			"Tokens consumed by user",
			[]string{"user_email"},
			nil,
		),
	}
}

func (e *UsageEventsExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.totalEvents
	ch <- e.eventsByType
	ch <- e.eventsByUser
	ch <- e.eventsByModel
	ch <- e.tokensConsumed
	ch <- e.tokensConsumedByModel
	ch <- e.tokensConsumedByUser
}

func (e *UsageEventsExporter) Collect(ch chan<- prometheus.Metric) {
	endDate := time.Now().Format("2006-01-02")
	startDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")

	events, err := e.client.GetUsageEvents("", 5000, 0, startDate, endDate)
	if err != nil {
		logrus.WithError(err).Error("Failed to get usage events")
		return
	}

	eventTypeCount := make(map[string]int)
	userEventCount := make(map[string]int)
	modelEventCount := make(map[string]int)
	modelTokenCount := make(map[string]int)
	userTokenCount := make(map[string]int)
	totalTokens := 0

	for _, event := range events {
		eventTypeCount[event.EventType]++
		userEventCount[event.UserEmail]++
		modelEventCount[event.Model]++
		modelTokenCount[event.Model] += event.TokensConsumed
		userTokenCount[event.UserEmail] += event.TokensConsumed
		totalTokens += event.TokensConsumed
	}

	ch <- prometheus.MustNewConstMetric(
		e.totalEvents,
		prometheus.GaugeValue,
		float64(len(events)),
	)

	for eventType, count := range eventTypeCount {
		ch <- prometheus.MustNewConstMetric(
			e.eventsByType,
			prometheus.GaugeValue,
			float64(count),
			eventType,
		)
	}

	for userEmail, count := range userEventCount {
		ch <- prometheus.MustNewConstMetric(
			e.eventsByUser,
			prometheus.GaugeValue,
			float64(count),
			userEmail,
		)
	}

	for model, count := range modelEventCount {
		ch <- prometheus.MustNewConstMetric(
			e.eventsByModel,
			prometheus.GaugeValue,
			float64(count),
			model,
		)
	}

	ch <- prometheus.MustNewConstMetric(
		e.tokensConsumed,
		prometheus.GaugeValue,
		float64(totalTokens),
	)

	for model, tokens := range modelTokenCount {
		ch <- prometheus.MustNewConstMetric(
			e.tokensConsumedByModel,
			prometheus.GaugeValue,
			float64(tokens),
			model,
		)
	}

	for userEmail, tokens := range userTokenCount {
		ch <- prometheus.MustNewConstMetric(
			e.tokensConsumedByUser,
			prometheus.GaugeValue,
			float64(tokens),
			userEmail,
		)
	}
}
