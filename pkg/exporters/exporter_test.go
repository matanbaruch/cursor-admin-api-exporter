package exporters

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/matanbaruch/cursor-admin-api-exporter/pkg/client"
	"github.com/prometheus/client_golang/prometheus"
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
		case "/teams/members":
			response := struct {
				TeamMembers []client.TeamMember `json:"teamMembers"`
			}{
				TeamMembers: []client.TeamMember{
					{Name: "John Doe", Email: "john@example.com", Role: "admin"},
					{Name: "Jane Smith", Email: "jane@example.com", Role: "member"},
				},
			}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Logf("Failed to encode response: %v", err)
			}
		case "/teams/daily-usage-data":
			var reqBody struct {
				StartDate int64 `json:"startDate"`
				EndDate   int64 `json:"endDate"`
			}
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Errorf("Failed to decode request: %v", err)
			}
			response := struct {
				Data []struct {
					Date                 int64  `json:"date"`
					TotalLinesAdded      int    `json:"totalLinesAdded"`
					TotalLinesDeleted    int    `json:"totalLinesDeleted"`
					TotalAccepts         int    `json:"totalAccepts"`
					TotalRejects         int    `json:"totalRejects"`
					TotalTabsAccepted    int    `json:"totalTabsAccepted"`
					ComposerRequests     int    `json:"composerRequests"`
					ChatRequests         int    `json:"chatRequests"`
					MostUsedModel        string `json:"mostUsedModel"`
					TabMostUsedExtension string `json:"tabMostUsedExtension"`
				} `json:"data"`
			}{
				Data: []struct {
					Date                 int64  `json:"date"`
					TotalLinesAdded      int    `json:"totalLinesAdded"`
					TotalLinesDeleted    int    `json:"totalLinesDeleted"`
					TotalAccepts         int    `json:"totalAccepts"`
					TotalRejects         int    `json:"totalRejects"`
					TotalTabsAccepted    int    `json:"totalTabsAccepted"`
					ComposerRequests     int    `json:"composerRequests"`
					ChatRequests         int    `json:"chatRequests"`
					MostUsedModel        string `json:"mostUsedModel"`
					TabMostUsedExtension string `json:"tabMostUsedExtension"`
				}{
					{
						Date:                 time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli(),
						TotalLinesAdded:      100,
						TotalLinesDeleted:    50,
						TotalAccepts:         8,
						TotalRejects:         2,
						TotalTabsAccepted:    20,
						ComposerRequests:     5,
						ChatRequests:         10,
						MostUsedModel:        "gpt-4",
						TabMostUsedExtension: "python",
					},
				},
			}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Logf("Failed to encode response: %v", err)
			}
		case "/teams/spend":
			var reqBody struct {
				Page     int `json:"page"`
				PageSize int `json:"pageSize"`
			}
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Errorf("Failed to decode request: %v", err)
			}
			response := struct {
				TeamMemberSpend []struct {
					SpendCents          int    `json:"spendCents"`
					FastPremiumRequests int    `json:"fastPremiumRequests"`
					Email               string `json:"email"`
				} `json:"teamMemberSpend"`
				SubscriptionCycleStart int64 `json:"subscriptionCycleStart"`
				TotalPages             int   `json:"totalPages"`
			}{
				TeamMemberSpend: []struct {
					SpendCents          int    `json:"spendCents"`
					FastPremiumRequests int    `json:"fastPremiumRequests"`
					Email               string `json:"email"`
				}{
					{
						SpendCents:          1000,
						FastPremiumRequests: 50,
						Email:               "john@example.com",
					},
				},
				SubscriptionCycleStart: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli(),
				TotalPages:             1,
			}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Logf("Failed to encode response: %v", err)
			}
		case "/teams/filtered-usage-events":
			var reqBody struct {
				Page     int `json:"page"`
				PageSize int `json:"pageSize"`
			}
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Errorf("Failed to decode request: %v", err)
			}
			response := struct {
				UsageEvents []struct {
					Timestamp  string `json:"timestamp"`
					Model      string `json:"model"`
					KindLabel  string `json:"kindLabel"`
					TokenUsage *struct {
						InputTokens      int `json:"inputTokens"`
						OutputTokens     int `json:"outputTokens"`
						CacheWriteTokens int `json:"cacheWriteTokens"`
						CacheReadTokens  int `json:"cacheReadTokens"`
					} `json:"tokenUsage"`
					UserEmail  string `json:"userEmail"`
				} `json:"usageEvents"`
				Pagination struct {
					HasNextPage bool `json:"hasNextPage"`
				} `json:"pagination"`
			}{
				UsageEvents: []struct {
					Timestamp  string `json:"timestamp"`
					Model      string `json:"model"`
					KindLabel  string `json:"kindLabel"`
					TokenUsage *struct {
						InputTokens      int `json:"inputTokens"`
						OutputTokens     int `json:"outputTokens"`
						CacheWriteTokens int `json:"cacheWriteTokens"`
						CacheReadTokens  int `json:"cacheReadTokens"`
					} `json:"tokenUsage"`
					UserEmail  string `json:"userEmail"`
				}{
					{
						Timestamp: strconv.FormatInt(time.Now().UnixMilli(), 10),
						Model:     "gpt-4",
						KindLabel: "completion",
						TokenUsage: &struct {
							InputTokens      int `json:"inputTokens"`
							OutputTokens     int `json:"outputTokens"`
							CacheWriteTokens int `json:"cacheWriteTokens"`
							CacheReadTokens  int `json:"cacheReadTokens"`
						}{
							InputTokens:      50,
							OutputTokens:     40,
							CacheWriteTokens: 10,
							CacheReadTokens:  0,
						},
						UserEmail: "john@example.com",
					},
				},
				Pagination: struct {
					HasNextPage bool `json:"hasNextPage"`
				}{
					HasNextPage: false,
				},
			}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Logf("Failed to encode response: %v", err)
			}
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
		client:              nil,
		teamMembersExporter: NewTeamMembersExporter(invalidClient),
		dailyUsageExporter:  NewDailyUsageExporter(invalidClient),
		spendingExporter:    NewSpendingExporter(invalidClient),
		usageEventsExporter: NewUsageEventsExporter(invalidClient),
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
		case "/teams/members":
			if err := json.NewEncoder(w).Encode(struct{ TeamMembers []client.TeamMember }{TeamMembers: []client.TeamMember{}}); err != nil {
				t.Logf("Failed to encode response: %v", err)
			}
		case "/teams/daily-usage-data":
			if err := json.NewEncoder(w).Encode(struct{ Data []struct{} }{Data: []struct{}{}}); err != nil {
				t.Logf("Failed to encode response: %v", err)
			}
		case "/teams/spend":
			if err := json.NewEncoder(w).Encode(struct{ TeamMemberSpend []struct{} }{TeamMemberSpend: []struct{}{}}); err != nil {
				t.Logf("Failed to encode response: %v", err)
			}
		case "/teams/filtered-usage-events":
			if err := json.NewEncoder(w).Encode(struct{ UsageEvents []struct{} }{UsageEvents: []struct{}{}}); err != nil {
				t.Logf("Failed to encode response: %v", err)
			}
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
