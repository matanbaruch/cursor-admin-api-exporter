package exporters

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matanbaruch/cursor-admin-api-exporter/pkg/client"
	"github.com/prometheus/client_golang/prometheus"
)

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
		if r.URL.Path != "/admin/team/members" {
			http.NotFound(w, r)
			return
		}

		response := client.TeamMembersResponse{
			Members: []client.TeamMember{
				{Name: "John Doe", Email: "john@example.com", Role: "admin"},
				{Name: "Jane Smith", Email: "jane@example.com", Role: "member"},
				{Name: "Bob Johnson", Email: "bob@example.com", Role: "member"},
				{Name: "Alice Brown", Email: "alice@example.com", Role: "viewer"},
			},
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

	if metricCount != 4 {
		t.Errorf("Expected 4 metrics (1 total + 3 role counts), got %d", metricCount)
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
