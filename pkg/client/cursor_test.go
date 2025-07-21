package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestNewCursorClient(t *testing.T) {
	baseURL := "https://api.cursor.com"
	apiToken := "test-token"

	client := NewCursorClient(baseURL, apiToken)

	if client == nil {
		t.Fatal("Expected client to be non-nil")
	}

	if client.BaseURL != baseURL {
		t.Errorf("Expected BaseURL to be %s, got %s", baseURL, client.BaseURL)
	}

	if client.APIToken != apiToken {
		t.Errorf("Expected APIToken to be %s, got %s", apiToken, client.APIToken)
	}

	if client.HTTPClient == nil {
		t.Error("Expected HTTPClient to be non-nil")
	}

	if client.HTTPClient.Timeout != 30*time.Second {
		t.Errorf("Expected timeout to be 30s, got %v", client.HTTPClient.Timeout)
	}
}

func TestCursorClient_GetTeamMembers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/teams/members" {
			t.Errorf("Expected path /teams/members, got %s", r.URL.Path)
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got %s", authHeader)
		}

		response := struct {
			TeamMembers []TeamMember `json:"teamMembers"`
		}{
			TeamMembers: []TeamMember{
				{Name: "John Doe", Email: "john@example.com", Role: "admin"},
				{Name: "Jane Smith", Email: "jane@example.com", Role: "member"},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Logf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewCursorClient(server.URL, "test-token")
	members, err := client.GetTeamMembers()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(members) != 2 {
		t.Errorf("Expected 2 members, got %d", len(members))
	}

	if members[0].Name != "John Doe" {
		t.Errorf("Expected first member name to be 'John Doe', got %s", members[0].Name)
	}

	if members[1].Role != "member" {
		t.Errorf("Expected second member role to be 'member', got %s", members[1].Role)
	}
}

func TestCursorClient_GetDailyUsage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/teams/daily-usage-data" {
			t.Errorf("Expected path /teams/daily-usage-data, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		var reqBody struct {
			StartDate int64 `json:"startDate"`
			EndDate   int64 `json:"endDate"`
		}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if reqBody.StartDate != time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli() {
			t.Errorf("Expected startDate %d, got %d", time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli(), reqBody.StartDate)
		}
		if reqBody.EndDate != time.Date(2023, 1, 31, 23, 59, 59, 999000000, time.UTC).UnixMilli() {
			t.Errorf("Expected endDate %d, got %d", time.Date(2023, 1, 31, 23, 59, 59, 999000000, time.UTC).UnixMilli(), reqBody.EndDate)
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

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Logf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewCursorClient(server.URL, "test-token")
	usage, err := client.GetDailyUsage("2023-01-01", "2023-01-31")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(usage) != 1 {
		t.Errorf("Expected 1 usage record, got %d", len(usage))
	}

	if usage[0].LinesAdded != 100 {
		t.Errorf("Expected 100 lines added, got %d", usage[0].LinesAdded)
	}

	if usage[0].SuggestionAcceptanceRate != 0.8 {
		t.Errorf("Expected 0.8 acceptance rate, got %f", usage[0].SuggestionAcceptanceRate)
	}
}

func TestCursorClient_GetSpending(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/teams/spend" {
			t.Errorf("Expected path /teams/spend, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		var reqBody struct {
			Page     int `json:"page"`
			PageSize int `json:"pageSize"`
		}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if reqBody.Page != 1 {
			t.Errorf("Expected page 1, got %d", reqBody.Page)
		}
		if reqBody.PageSize != 100 {
			t.Errorf("Expected pageSize 100, got %d", reqBody.PageSize)
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

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Logf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewCursorClient(server.URL, "test-token")
	spending, err := client.GetSpending(100, 0)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(spending) != 1 {
		t.Errorf("Expected 1 spending record, got %d", len(spending))
	}

	if spending[0].SpendCents != 1000 {
		t.Errorf("Expected 1000 cents, got %d", spending[0].SpendCents)
	}

	if spending[0].PremiumRequests != 50 {
		t.Errorf("Expected 50 premium requests, got %d", spending[0].PremiumRequests)
	}
	if spending[0].Date != "2023-01-01" {
		t.Errorf("Expected date '2023-01-01', got %s", spending[0].Date)
	}
}

func TestCursorClient_GetUsageEvents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/teams/filtered-usage-events" {
			t.Errorf("Expected path /teams/filtered-usage-events, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		var reqBody struct {
			Email    string `json:"email"`
			Page     int    `json:"page"`
			PageSize int    `json:"pageSize"`
		}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if reqBody.Email != "john@example.com" {
			t.Errorf("Expected email 'john@example.com', got %s", reqBody.Email)
		}
		if reqBody.Page != 1 {
			t.Errorf("Expected page 1, got %d", reqBody.Page)
		}
		if reqBody.PageSize != 50 {
			t.Errorf("Expected pageSize 50, got %d", reqBody.PageSize)
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
				UserEmail string `json:"userEmail"`
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
				UserEmail string `json:"userEmail"`
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

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Logf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewCursorClient(server.URL, "test-token")
	events, err := client.GetUsageEvents("john@example.com", 50, 0, "", "")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	if events[0].EventType != "completion" {
		t.Errorf("Expected event type 'completion', got %s", events[0].EventType)
	}

	if events[0].TokensConsumed != 100 {
		t.Errorf("Expected 100 tokens consumed, got %d", events[0].TokensConsumed)
	}
}

func TestCursorClient_ErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewCursorClient(server.URL, "invalid-token")

	_, err := client.GetTeamMembers()
	if err == nil {
		t.Error("Expected error for unauthorized request")
	}

	_, err = client.GetDailyUsage("2023-01-01", "2023-01-31")
	if err == nil {
		t.Error("Expected error for unauthorized request")
	}

	_, err = client.GetSpending(100, 0)
	if err == nil {
		t.Error("Expected error for unauthorized request")
	}

	_, err = client.GetUsageEvents("test@example.com", 50, 0, "", "")
	if err == nil {
		t.Error("Expected error for unauthorized request")
	}
}
