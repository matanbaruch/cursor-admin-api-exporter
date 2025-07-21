package exporters

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/matanbaruch/cursor-admin-api-exporter/pkg/client"
)

func TestNewCursorExporter(t *testing.T) {
	baseURL := "https://api.cursor.com"
	token := "test-token"

	exporter := NewCursorExporter(baseURL, token)

	if exporter == nil {
		t.Fatal("Expected exporter to be non-nil")
	}

	if exporter.client == nil {
		t.Error("Expected client to be non-nil")
	}

	if exporter.teamMembersExporter == nil {
		t.Error("Expected teamMembersExporter to be non-nil")
	}

	if exporter.dailyUsageExporter == nil {
		t.Error("Expected dailyUsageExporter to be non-nil")
	}

	if exporter.spendingExporter == nil {
		t.Error("Expected spendingExporter to be non-nil")
	}

	if exporter.usageEventsExporter == nil {
		t.Error("Expected usageEventsExporter to be non-nil")
	}

	if exporter.scrapeDuration == nil {
		t.Error("Expected scrapeDuration metric to be non-nil")
	}

	if exporter.scrapeErrors == nil {
		t.Error("Expected scrapeErrors metric to be non-nil")
	}
}

func TestCursorExporter_Describe(t *testing.T) {
	exporter := NewCursorExporter("https://api.cursor.com", "test-token")

	ch := make(chan *prometheus.Desc, 100)
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

	if count == 0 {
		t.Error("Expected at least one metric description")
	}
}

func TestCursorExporter_Collect_WithMockServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		switch r.URL.Path {
		case "/admin/team/members":
			response := client.TeamMembersResponse{
				Members: []client.TeamMember{
					{Name: "John Doe", Email: "john@example.com", Role: "admin"},
					{Name: "Jane Smith", Email: "jane@example.com", Role: "member"},
				},
			}
			json.NewEncoder(w).Encode(response)
		case "/admin/usage/daily":
			response := client.DailyUsageResponse{
				Usage: []client.DailyUsage{
					{
						Date:                     "2023-01-01",
						LinesAdded:               100,
						LinesDeleted:             50,
						SuggestionAcceptanceRate: 0.8,
						TabsUsed:                 20,
						ComposerUsed:             5,
						ChatRequests:             10,
						MostUsedModel:            "gpt-4",
						MostUsedExtension:        "python",
					},
				},
			}
			json.NewEncoder(w).Encode(response)
		case "/admin/spending":
			response := client.SpendingResponse{
				Spending: []client.SpendingData{
					{
						MemberEmail:     "john@example.com",
						SpendCents:      1000,
						PremiumRequests: 50,
						Date:            "2023-01-01",
					},
				},
				Total: 1000,
			}
			json.NewEncoder(w).Encode(response)
		case "/admin/usage/events":
			response := client.UsageEventsResponse{
				Events: []client.UsageEvent{
					{
						EventType:      "completion",
						UserEmail:      "john@example.com",
						TokensConsumed: 100,
						Model:          "gpt-4",
						Timestamp:      time.Now(),
					},
				},
				Total: 1,
			}
			json.NewEncoder(w).Encode(response)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	exporter := NewCursorExporter(server.URL, "test-token")

	ch := make(chan prometheus.Metric, 100)
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

	if metricCount == 0 {
		t.Error("Expected at least one metric to be collected")
	}
}

func TestCursorExporter_Collect_HandlesErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer server.Close()

	exporter := NewCursorExporter(server.URL, "test-token")

	ch := make(chan prometheus.Metric, 100)
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

	if metricCount == 0 {
		t.Error("Expected at least error metrics to be collected")
	}
}

func TestCursorExporter_Collect_HandlesPanics(t *testing.T) {
	invalidClient := client.NewCursorClient("http://invalid", "token")
	exporter := &CursorExporter{
		client:               nil,
		teamMembersExporter:  NewTeamMembersExporter(invalidClient),
		dailyUsageExporter:   NewDailyUsageExporter(invalidClient),
		spendingExporter:     NewSpendingExporter(invalidClient),
		usageEventsExporter:  NewUsageEventsExporter(invalidClient),
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

	ch := make(chan prometheus.Metric, 100)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Collect should not panic, but got: %v", r)
			}
			close(ch)
		}()
		exporter.Collect(ch)
	}()

	for range ch {
	}
}

func TestCursorExporter_MetricsRegistration(t *testing.T) {
	registry := prometheus.NewRegistry()
	exporter := NewCursorExporter("https://api.cursor.com", "test-token")

	err := registry.Register(exporter)
	if err != nil {
		t.Fatalf("Failed to register exporter: %v", err)
	}

	families, err := registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	if len(families) == 0 {
		t.Error("Expected at least one metric family")
	}
}

func TestCursorExporter_ScrapeMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		switch r.URL.Path {
		case "/admin/team/members":
			json.NewEncoder(w).Encode(client.TeamMembersResponse{Members: []client.TeamMember{}})
		case "/admin/usage/daily":
			json.NewEncoder(w).Encode(client.DailyUsageResponse{Usage: []client.DailyUsage{}})
		case "/admin/spending":
			json.NewEncoder(w).Encode(client.SpendingResponse{Spending: []client.SpendingData{}})
		case "/admin/usage/events":
			json.NewEncoder(w).Encode(client.UsageEventsResponse{Events: []client.UsageEvent{}})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	exporter := NewCursorExporter(server.URL, "test-token")
	registry := prometheus.NewRegistry()
	registry.MustRegister(exporter)

	families, err := registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	found := false
	for _, family := range families {
		if family.GetName() == "cursor_exporter_scrape_duration_seconds" {
			found = true
			if len(family.GetMetric()) == 0 {
				t.Error("Expected scrape duration metric to have values")
			}
			break
		}
	}

	if !found {
		t.Error("Expected to find scrape duration metric")
	}
}