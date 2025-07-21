package exporters

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/matanbaruch/cursor-admin-api-exporter/pkg/client"
)

type TeamMembersExporter struct {
	client *client.CursorClient

	totalMembers *prometheus.Desc
	membersByRole *prometheus.Desc
}

func NewTeamMembersExporter(client *client.CursorClient) *TeamMembersExporter {
	return &TeamMembersExporter{
		client: client,

		totalMembers: prometheus.NewDesc(
			"cursor_team_members_total",
			"Total number of team members",
			nil,
			nil,
		),

		membersByRole: prometheus.NewDesc(
			"cursor_team_members_by_role",
			"Number of team members by role",
			[]string{"role"},
			nil,
		),
	}
}

func (e *TeamMembersExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.totalMembers
	ch <- e.membersByRole
}

func (e *TeamMembersExporter) Collect(ch chan<- prometheus.Metric) {
	members, err := e.client.GetTeamMembers()
	if err != nil {
		logrus.WithError(err).Error("Failed to get team members")
		return
	}

	ch <- prometheus.MustNewConstMetric(
		e.totalMembers,
		prometheus.GaugeValue,
		float64(len(members)),
	)

	roleCount := make(map[string]int)
	for _, member := range members {
		roleCount[member.Role]++
	}

	for role, count := range roleCount {
		ch <- prometheus.MustNewConstMetric(
			e.membersByRole,
			prometheus.GaugeValue,
			float64(count),
			role,
		)
	}
}