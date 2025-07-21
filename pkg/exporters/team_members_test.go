package exporters

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matanbaruch/cursor-admin-api-exporter/pkg/client"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func dtoMetric(m prometheus.Metric) *dto.Metric {
	var dm dto.Metric
	if err := m.Write(&dm); err != nil {
		return nil
	}
	return &dm
}

func findMetric(metrics []prometheus.Metric, name string) prometheus.Metric {
	for _, m := range metrics {
		desc := m.Desc()
		if desc.String() != "" && strings.Contains(desc.String(), name) {
			return m
		}
	}
	return nil
}

func findMetricWithLabel(metrics []prometheus.Metric, name, labelName, labelValue string) prometheus.Metric {
	for _, m := range metrics {
		desc := m.Desc()
		if !strings.Contains(desc.String(), name) {
			continue
		}
		dm := dtoMetric(m)
		for _, lp := range dm.GetLabel() {
			if lp.GetName() == labelName && lp.GetValue() == labelValue {
				return m
			}
		}
	}
	return nil
}

func TestNewTeamMembersExporter(t *testing.T) {
	mockClient := client.NewCursorClient("https://api.cursor.com", "test-token")
	exporter := NewTeamMembersExporter(mockClient)

	if exporter == nil {
		t.Fatal("Expected exporter to be non-nil")
	}

	if exporter.client == nil {
		t.Error("Expected client to be non-nil")
	}

	if exporter.totalMembers == nil {
		t.Error("Expected totalMembers metric to be non-nil")
	}

	if exporter.membersByRole == nil {
		t.Error("Expected membersByRole metric to be non-nil")
	}
}

func TestTeamMembersExporter_Describe(t *testing.T) {
	mockClient := client.NewCursorClient("https://api.cursor.com", "test-token")
	exporter := NewTeamMembersExporter(mockClient)

	ch := make(chan *prometheus.Desc, 10)
	go func() {
		exporter.Describe(ch)
		close(ch)
	}()

	count := 0
	for desc := range ch {
		if desc == nil {
			t.Error("Expected metric description to be non-nil")
		}
		count++
	}

	if count != 2 {
		t.Errorf("Expected 2 metric descriptions, got %d", count)
	}
}

func TestTeamMembersExporter_Collect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/teams/members" {
			t.Errorf("Expected path /teams/members, got %s", r.URL.Path)
		}

		response := struct {
			TeamMembers []client.TeamMember `json:"teamMembers"`
		}{
			TeamMembers: []client.TeamMember{
				{Name: "Admin User", Email: "admin@example.com", Role: "admin"},
				{Name: "Member One", Email: "member1@example.com", Role: "member"},
				{Name: "Member Two", Email: "member2@example.com", Role: "member"},
				{Name: "Viewer User", Email: "viewer@example.com", Role: "viewer"},
			},
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Logf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := client.NewCursorClient(server.URL, "test-token")
	exporter := NewTeamMembersExporter(client)

	ch := make(chan prometheus.Metric)
	go func() {
		exporter.Collect(ch)
		close(ch)
	}()

	metrics := []prometheus.Metric{}
	for m := range ch {
		metrics = append(metrics, m)
	}

	if len(metrics) != 4 {
		t.Errorf("Expected 4 metrics (1 total + 3 role counts), got %d", len(metrics))
	}

	// Check total members
	totalMetric := findMetric(metrics, "cursor_team_members_total")
	if totalMetric == nil {
		t.Error("Expected cursor_team_members_total metric")
	}
	if dtoMetric(totalMetric).GetGauge().GetValue() != 4 {
		t.Errorf("Expected total members 4, got %f", dtoMetric(totalMetric).GetGauge().GetValue())
	}

	// Check by role
	adminMetric := findMetricWithLabel(metrics, "cursor_team_members_by_role", "role", "admin")
	if adminMetric == nil {
		t.Error("Expected cursor_team_members_by_role{role=\"admin\"} metric")
	}
	if dtoMetric(adminMetric).GetGauge().GetValue() != 1 {
		t.Errorf("Expected 1 admin, got %f", dtoMetric(adminMetric).GetGauge().GetValue())
	}

	memberMetric := findMetricWithLabel(metrics, "cursor_team_members_by_role", "role", "member")
	if dtoMetric(memberMetric).GetGauge().GetValue() != 2 {
		t.Errorf("Expected 2 members, got %f", dtoMetric(memberMetric).GetGauge().GetValue())
	}

	viewerMetric := findMetricWithLabel(metrics, "cursor_team_members_by_role", "role", "viewer")
	if dtoMetric(viewerMetric).GetGauge().GetValue() != 1 {
		t.Errorf("Expected 1 viewer, got %f", dtoMetric(viewerMetric).GetGauge().GetValue())
	}
}

func TestTeamMembersExporter_Collect_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := client.TeamMembersResponse{
			Members: []client.TeamMember{},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Logf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	mockClient := client.NewCursorClient(server.URL, "test-token")
	exporter := NewTeamMembersExporter(mockClient)

	ch := make(chan prometheus.Metric, 10)
	go func() {
		exporter.Collect(ch)
		close(ch)
	}()

	metricCount := 0
	for metric := range ch {
		if metric == nil {
			t.Error("Expected metric to be non-nil")
		}
		metricCount++
	}

	if metricCount != 1 {
		t.Errorf("Expected 1 metric (total members = 0), got %d", metricCount)
	}
}

func TestTeamMembersExporter_Collect_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer server.Close()

	mockClient := client.NewCursorClient(server.URL, "test-token")
	exporter := NewTeamMembersExporter(mockClient)

	ch := make(chan prometheus.Metric, 10)
	go func() {
		exporter.Collect(ch)
		close(ch)
	}()

	metricCount := 0
	for metric := range ch {
		if metric == nil {
			t.Error("Expected metric to be non-nil")
		}
		metricCount++
	}

	if metricCount != 0 {
		t.Errorf("Expected 0 metrics on API error, got %d", metricCount)
	}
}
