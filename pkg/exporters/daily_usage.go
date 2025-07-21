package exporters

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/matanbaruch/cursor-admin-api-exporter/pkg/client"
)

type DailyUsageExporter struct {
	client *client.CursorClient

	linesAdded               *prometheus.Desc
	linesDeleted             *prometheus.Desc
	suggestionAcceptanceRate *prometheus.Desc
	tabsUsed                 *prometheus.Desc
	composerUsed             *prometheus.Desc
	chatRequests             *prometheus.Desc
	modelUsage               *prometheus.Desc
	extensionUsage           *prometheus.Desc
}

func NewDailyUsageExporter(client *client.CursorClient) *DailyUsageExporter {
	return &DailyUsageExporter{
		client: client,

		linesAdded: prometheus.NewDesc(
			"cursor_daily_lines_added_total",
			"Total lines of code added per day",
			[]string{"date"},
			nil,
		),

		linesDeleted: prometheus.NewDesc(
			"cursor_daily_lines_deleted_total",
			"Total lines of code deleted per day",
			[]string{"date"},
			nil,
		),

		suggestionAcceptanceRate: prometheus.NewDesc(
			"cursor_daily_suggestion_acceptance_rate",
			"AI suggestion acceptance rate per day",
			[]string{"date"},
			nil,
		),

		tabsUsed: prometheus.NewDesc(
			"cursor_daily_tabs_used_total",
			"Total tabs used per day",
			[]string{"date"},
			nil,
		),

		composerUsed: prometheus.NewDesc(
			"cursor_daily_composer_used_total",
			"Total composer usage per day",
			[]string{"date"},
			nil,
		),

		chatRequests: prometheus.NewDesc(
			"cursor_daily_chat_requests_total",
			"Total chat requests per day",
			[]string{"date"},
			nil,
		),

		modelUsage: prometheus.NewDesc(
			"cursor_daily_model_usage",
			"Most used model per day",
			[]string{"date", "model"},
			nil,
		),

		extensionUsage: prometheus.NewDesc(
			"cursor_daily_extension_usage",
			"Most used extension per day",
			[]string{"date", "extension"},
			nil,
		),
	}
}

func (e *DailyUsageExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.linesAdded
	ch <- e.linesDeleted
	ch <- e.suggestionAcceptanceRate
	ch <- e.tabsUsed
	ch <- e.composerUsed
	ch <- e.chatRequests
	ch <- e.modelUsage
	ch <- e.extensionUsage
}

func (e *DailyUsageExporter) Collect(ch chan<- prometheus.Metric) {
	endDate := time.Now().Format("2006-01-02")
	startDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")

	usage, err := e.client.GetDailyUsage(startDate, endDate)
	if err != nil {
		logrus.WithError(err).Error("Failed to get daily usage")
		return
	}

	for _, daily := range usage {
		ch <- prometheus.MustNewConstMetric(
			e.linesAdded,
			prometheus.GaugeValue,
			float64(daily.LinesAdded),
			daily.Date,
		)

		ch <- prometheus.MustNewConstMetric(
			e.linesDeleted,
			prometheus.GaugeValue,
			float64(daily.LinesDeleted),
			daily.Date,
		)

		ch <- prometheus.MustNewConstMetric(
			e.suggestionAcceptanceRate,
			prometheus.GaugeValue,
			daily.SuggestionAcceptanceRate,
			daily.Date,
		)

		ch <- prometheus.MustNewConstMetric(
			e.tabsUsed,
			prometheus.GaugeValue,
			float64(daily.TabsUsed),
			daily.Date,
		)

		ch <- prometheus.MustNewConstMetric(
			e.composerUsed,
			prometheus.GaugeValue,
			float64(daily.ComposerUsed),
			daily.Date,
		)

		ch <- prometheus.MustNewConstMetric(
			e.chatRequests,
			prometheus.GaugeValue,
			float64(daily.ChatRequests),
			daily.Date,
		)

		if daily.MostUsedModel != "" {
			ch <- prometheus.MustNewConstMetric(
				e.modelUsage,
				prometheus.GaugeValue,
				1,
				daily.Date,
				daily.MostUsedModel,
			)
		}

		if daily.MostUsedExtension != "" {
			ch <- prometheus.MustNewConstMetric(
				e.extensionUsage,
				prometheus.GaugeValue,
				1,
				daily.Date,
				daily.MostUsedExtension,
			)
		}
	}
}
