package exporters

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/matanbaruch/cursor-admin-api-exporter/pkg/client"
)

type CursorExporter struct {
	client               *client.CursorClient
	teamMembersExporter  *TeamMembersExporter
	dailyUsageExporter   *DailyUsageExporter
	spendingExporter     *SpendingExporter
	usageEventsExporter  *UsageEventsExporter

	scrapeDuration prometheus.Histogram
	scrapeErrors   prometheus.Counter
}

func NewCursorExporter(baseURL, token string) *CursorExporter {
	cursorClient := client.NewCursorClient(baseURL, token)

	return &CursorExporter{
		client:               cursorClient,
		teamMembersExporter:  NewTeamMembersExporter(cursorClient),
		dailyUsageExporter:   NewDailyUsageExporter(cursorClient),
		spendingExporter:     NewSpendingExporter(cursorClient),
		usageEventsExporter:  NewUsageEventsExporter(cursorClient),

		scrapeDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name: "cursor_exporter_scrape_duration_seconds",
				Help: "Time spent scraping Cursor Admin API",
			},
		),

		scrapeErrors: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "cursor_exporter_scrape_errors_total",
				Help: "Total number of scrape errors",
			},
		),
	}
}

func (e *CursorExporter) Describe(ch chan<- *prometheus.Desc) {
	e.teamMembersExporter.Describe(ch)
	e.dailyUsageExporter.Describe(ch)
	e.spendingExporter.Describe(ch)
	e.usageEventsExporter.Describe(ch)
	e.scrapeDuration.Describe(ch)
	e.scrapeErrors.Describe(ch)
}

func (e *CursorExporter) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		e.scrapeDuration.Observe(duration.Seconds())
		e.scrapeDuration.Collect(ch)
		e.scrapeErrors.Collect(ch)
		logrus.WithField("total_duration", duration).Debug("Completed Cursor metrics collection")
	}()

	logrus.Debug("Starting Cursor metrics collection")

	func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.WithField("panic", r).Error("Panic during team members collection")
				e.scrapeErrors.Inc()
			}
		}()
		logrus.Debug("Starting team members collection")
		e.teamMembersExporter.Collect(ch)
		logrus.Debug("Completed team members collection")
	}()

	func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.WithField("panic", r).Error("Panic during daily usage collection")
				e.scrapeErrors.Inc()
			}
		}()
		logrus.Debug("Starting daily usage collection")
		e.dailyUsageExporter.Collect(ch)
		logrus.Debug("Completed daily usage collection")
	}()

	func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.WithField("panic", r).Error("Panic during spending collection")
				e.scrapeErrors.Inc()
			}
		}()
		logrus.Debug("Starting spending collection")
		e.spendingExporter.Collect(ch)
		logrus.Debug("Completed spending collection")
	}()

	func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.WithField("panic", r).Error("Panic during usage events collection")
				e.scrapeErrors.Inc()
			}
		}()
		logrus.Debug("Starting usage events collection")
		e.usageEventsExporter.Collect(ch)
		logrus.Debug("Completed usage events collection")
	}()
}