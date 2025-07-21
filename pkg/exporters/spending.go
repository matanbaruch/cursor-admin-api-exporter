package exporters

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/matanbaruch/cursor-admin-api-exporter/pkg/client"
)

type SpendingExporter struct {
	client *client.CursorClient

	totalSpending           *prometheus.Desc
	spendingByMember        *prometheus.Desc
	premiumRequestsByMember *prometheus.Desc
	totalPremiumRequests    *prometheus.Desc
}

func NewSpendingExporter(client *client.CursorClient) *SpendingExporter {
	return &SpendingExporter{
		client: client,

		totalSpending: prometheus.NewDesc(
			"cursor_spending_total_cents",
			"Total spending in cents",
			nil,
			nil,
		),

		spendingByMember: prometheus.NewDesc(
			"cursor_spending_by_member_cents",
			"Spending by team member in cents",
			[]string{"member_email", "date"},
			nil,
		),

		premiumRequestsByMember: prometheus.NewDesc(
			"cursor_premium_requests_by_member_total",
			"Premium requests by team member",
			[]string{"member_email", "date"},
			nil,
		),

		totalPremiumRequests: prometheus.NewDesc(
			"cursor_premium_requests_total",
			"Total premium requests",
			nil,
			nil,
		),
	}
}

func (e *SpendingExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.totalSpending
	ch <- e.spendingByMember
	ch <- e.premiumRequestsByMember
	ch <- e.totalPremiumRequests
}

func (e *SpendingExporter) Collect(ch chan<- prometheus.Metric) {
	spending, err := e.client.GetSpending(1000, 0)
	if err != nil {
		logrus.WithError(err).Error("Failed to get spending data")
		return
	}

	var totalSpend int
	var totalPremiumRequests int

	for _, spend := range spending {
		totalSpend += spend.SpendCents
		totalPremiumRequests += spend.PremiumRequests

		ch <- prometheus.MustNewConstMetric(
			e.spendingByMember,
			prometheus.GaugeValue,
			float64(spend.SpendCents),
			spend.MemberEmail,
			spend.Date,
		)

		ch <- prometheus.MustNewConstMetric(
			e.premiumRequestsByMember,
			prometheus.GaugeValue,
			float64(spend.PremiumRequests),
			spend.MemberEmail,
			spend.Date,
		)
	}

	ch <- prometheus.MustNewConstMetric(
		e.totalSpending,
		prometheus.GaugeValue,
		float64(totalSpend),
	)

	ch <- prometheus.MustNewConstMetric(
		e.totalPremiumRequests,
		prometheus.GaugeValue,
		float64(totalPremiumRequests),
	)
}
