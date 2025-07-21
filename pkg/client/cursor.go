package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
)

type CursorClient struct {
	BaseURL    string
	APIToken   string
	HTTPClient *http.Client
}

type TeamMember struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type DailyUsage struct {
	Date                     string  `json:"date"`
	LinesAdded               int     `json:"lines_added"`
	LinesDeleted             int     `json:"lines_deleted"`
	SuggestionAcceptanceRate float64 `json:"suggestion_acceptance_rate"`
	TabsUsed                 int     `json:"tabs_used"`
	ComposerUsed             int     `json:"composer_used"`
	ChatRequests             int     `json:"chat_requests"`
	MostUsedModel            string  `json:"most_used_model"`
	MostUsedExtension        string  `json:"most_used_extension"`
}

type SpendingData struct {
	MemberEmail     string `json:"member_email"`
	SpendCents      int    `json:"spend_cents"`
	PremiumRequests int    `json:"premium_requests"`
	Date            string `json:"date"`
}

type UsageEvent struct {
	EventType      string    `json:"event_type"`
	UserEmail      string    `json:"user_email"`
	TokensConsumed int       `json:"tokens_consumed"`
	Model          string    `json:"model"`
	Timestamp      time.Time `json:"timestamp"`
}

type TeamMembersResponse struct {
	Members []TeamMember `json:"members"`
}

type DailyUsageResponse struct {
	Usage []DailyUsage `json:"usage"`
}

type SpendingResponse struct {
	Spending []SpendingData `json:"spending"`
	Total    int            `json:"total"`
}

type UsageEventsResponse struct {
	Events []UsageEvent `json:"events"`
	Total  int          `json:"total"`
}

func NewCursorClient(baseURL, apiToken string) *CursorClient {
	return &CursorClient{
		BaseURL:  baseURL,
		APIToken: apiToken,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *CursorClient) makeRequest(endpoint string, params url.Values) ([]byte, error) {
	fullURL := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	if params != nil {
		fullURL = fmt.Sprintf("%s?%s", fullURL, params.Encode())
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIToken))
	req.Header.Set("Content-Type", "application/json")

	logrus.WithField("url", fullURL).Debug("Making API request")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logrus.WithError(err).Debug("Failed to close response body")
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func (c *CursorClient) GetTeamMembers() ([]TeamMember, error) {
	body, err := c.makeRequest("/admin/team/members", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}

	var response TeamMembersResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal team members response: %w", err)
	}

	return response.Members, nil
}

func (c *CursorClient) GetDailyUsage(startDate, endDate string) ([]DailyUsage, error) {
	params := url.Values{}
	if startDate != "" {
		params.Set("start_date", startDate)
	}
	if endDate != "" {
		params.Set("end_date", endDate)
	}

	body, err := c.makeRequest("/admin/usage/daily", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily usage: %w", err)
	}

	var response DailyUsageResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal daily usage response: %w", err)
	}

	return response.Usage, nil
}

func (c *CursorClient) GetSpending(limit int, offset int) ([]SpendingData, error) {
	params := url.Values{}
	if limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}
	if offset >= 0 {
		params.Set("offset", fmt.Sprintf("%d", offset))
	}

	body, err := c.makeRequest("/admin/spending", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get spending data: %w", err)
	}

	var response SpendingResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal spending response: %w", err)
	}

	return response.Spending, nil
}

func (c *CursorClient) GetUsageEvents(userEmail string, limit int, offset int) ([]UsageEvent, error) {
	params := url.Values{}
	if userEmail != "" {
		params.Set("user_email", userEmail)
	}
	if limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}
	if offset >= 0 {
		params.Set("offset", fmt.Sprintf("%d", offset))
	}

	body, err := c.makeRequest("/admin/usage/events", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage events: %w", err)
	}

	var response UsageEventsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal usage events response: %w", err)
	}

	return response.Events, nil
}
