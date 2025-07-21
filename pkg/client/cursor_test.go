package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
		if r.URL.Path != "/admin/team/members" {
			t.Errorf("Expected path /admin/team/members, got %s", r.URL.Path)
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got %s", authHeader)
		}

		response := TeamMembersResponse{
			Members: []TeamMember{
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
		if r.URL.Path != "/admin/usage/daily" {
			t.Errorf("Expected path /admin/usage/daily, got %s", r.URL.Path)
		}

		startDate := r.URL.Query().Get("start_date")
		endDate := r.URL.Query().Get("end_date")

		if startDate != "2023-01-01" {
			t.Errorf("Expected start_date '2023-01-01', got %s", startDate)
		}

		if endDate != "2023-01-31" {
			t.Errorf("Expected end_date '2023-01-31', got %s", endDate)
		}

		response := DailyUsageResponse{
			Usage: []DailyUsage{
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
		if r.URL.Path != "/admin/spending" {
			t.Errorf("Expected path /admin/spending, got %s", r.URL.Path)
		}

		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")

		if limit != "100" {
			t.Errorf("Expected limit '100', got %s", limit)
		}

		if offset != "0" {
			t.Errorf("Expected offset '0', got %s", offset)
		}

		response := SpendingResponse{
			Spending: []SpendingData{
				{
					MemberEmail:     "john@example.com",
					SpendCents:      1000,
					PremiumRequests: 50,
					Date:            "2023-01-01",
				},
			},
			Total: 1000,
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
}

func TestCursorClient_GetUsageEvents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/admin/usage/events" {
			t.Errorf("Expected path /admin/usage/events, got %s", r.URL.Path)
		}

		userEmail := r.URL.Query().Get("user_email")
		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")

		if userEmail != "john@example.com" {
			t.Errorf("Expected user_email 'john@example.com', got %s", userEmail)
		}

		if limit != "50" {
			t.Errorf("Expected limit '50', got %s", limit)
		}

		if offset != "0" {
			t.Errorf("Expected offset '0', got %s", offset)
		}

		response := UsageEventsResponse{
			Events: []UsageEvent{
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

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Logf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewCursorClient(server.URL, "test-token")
	events, err := client.GetUsageEvents("john@example.com", 50, 0)

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

	_, err = client.GetUsageEvents("test@example.com", 50, 0)
	if err == nil {
		t.Error("Expected error for unauthorized request")
	}
}
