package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
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

type aggregatedData struct {
	linesAdded       int
	linesDeleted     int
	totalAccepts     int
	totalRejects     int
	tabsUsed         int
	composerUsed     int
	chatRequests     int
	modelCounts      map[string]int
	extensionCounts map[string]int
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

func (c *CursorClient) makeRequest(method string, endpoint string, params url.Values, body io.Reader) ([]byte, error) {
	fullURL := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	if params != nil {
		fullURL = fmt.Sprintf("%s?%s", fullURL, params.Encode())
	}

	req, err := http.NewRequest(method, fullURL, body)
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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return bodyBytes, nil
}

func (c *CursorClient) GetTeamMembers() ([]TeamMember, error) {
	body, err := c.makeRequest("GET", "/teams/members", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}

	var response struct {
		TeamMembers []TeamMember `json:"teamMembers"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal team members response: %w", err)
	}

	return response.TeamMembers, nil
}

func (c *CursorClient) GetDailyUsage(startDate, endDate string) ([]DailyUsage, error) {
	startT, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}
	endT, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}
	startMs := startT.UnixMilli()
	endMs := endT.Add(24*time.Hour - time.Millisecond).UnixMilli()

	reqBody := struct {
		StartDate int64 `json:"startDate"`
		EndDate   int64 `json:"endDate"`
	}{
		StartDate: startMs,
		EndDate:   endMs,
	}
	reqJson, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	body, err := c.makeRequest("POST", "/teams/daily-usage-data", nil, bytes.NewReader(reqJson))
	if err != nil {
		return nil, fmt.Errorf("failed to get daily usage: %w", err)
	}

	var response struct {
		Data []struct {
			Date                     int64  `json:"date"`
			TotalLinesAdded          int    `json:"totalLinesAdded"`
			TotalLinesDeleted        int    `json:"totalLinesDeleted"`
			TotalAccepts             int    `json:"totalAccepts"`
			TotalRejects             int    `json:"totalRejects"`
			TotalTabsAccepted        int    `json:"totalTabsAccepted"`
			ComposerRequests         int    `json:"composerRequests"`
			ChatRequests             int    `json:"chatRequests"`
			MostUsedModel            string `json:"mostUsedModel"`
			TabMostUsedExtension     string `json:"tabMostUsedExtension"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal daily usage response: %w", err)
	}

	aggMap := make(map[string]*aggregatedData)
	for _, d := range response.Data {
		dateStr := time.UnixMilli(d.Date).Format("2006-01-02")
		if _, ok := aggMap[dateStr]; !ok {
			aggMap[dateStr] = &aggregatedData{
				modelCounts:      make(map[string]int),
				extensionCounts: make(map[string]int),
			}
		}
		agg := aggMap[dateStr]
		agg.linesAdded += d.TotalLinesAdded
		agg.linesDeleted += d.TotalLinesDeleted
		agg.totalAccepts += d.TotalAccepts
		agg.totalRejects += d.TotalRejects
		agg.tabsUsed += d.TotalTabsAccepted
		agg.composerUsed += d.ComposerRequests
		agg.chatRequests += d.ChatRequests
		if d.MostUsedModel != "" {
			agg.modelCounts[d.MostUsedModel]++
		}
		if d.TabMostUsedExtension != "" {
			agg.extensionCounts[d.TabMostUsedExtension]++
		}
	}

	var usage []DailyUsage
	var dates []string
	for date := range aggMap {
		dates = append(dates, date)
	}
	sort.Strings(dates)
	for _, date := range dates {
		agg := aggMap[date]
		rate := 0.0
		totalSuggestions := agg.totalAccepts + agg.totalRejects
		if totalSuggestions > 0 {
			rate = float64(agg.totalAccepts) / float64(totalSuggestions)
		}

		var mostModel string
		maxModel := 0
		for m, count := range agg.modelCounts {
			if count > maxModel {
				maxModel = count
				mostModel = m
			}
		}

		var mostExt string
		maxExt := 0
		for e, count := range agg.extensionCounts {
			if count > maxExt {
				maxExt = count
				mostExt = e
			}
		}

		usage = append(usage, DailyUsage{
			Date:                     date,
			LinesAdded:               agg.linesAdded,
			LinesDeleted:             agg.linesDeleted,
			SuggestionAcceptanceRate: rate,
			TabsUsed:                 agg.tabsUsed,
			ComposerUsed:             agg.composerUsed,
			ChatRequests:             agg.chatRequests,
			MostUsedModel:            mostModel,
			MostUsedExtension:        mostExt,
		})
	}

	return usage, nil
}

func (c *CursorClient) GetSpending(limit int, offset int) ([]SpendingData, error) {
	var allSpending []SpendingData
	page := 1
	for {
		reqBody := struct {
			Page     int `json:"page"`
			PageSize int `json:"pageSize"`
		}{
			Page:     page,
			PageSize: limit,
		}
		reqJson, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}

		body, err := c.makeRequest("POST", "/teams/spend", nil, bytes.NewReader(reqJson))
		if err != nil {
			return nil, fmt.Errorf("failed to get spending data: %w", err)
		}

		var response struct {
			TeamMemberSpend []struct {
				SpendCents          int    `json:"spendCents"`
				FastPremiumRequests int    `json:"fastPremiumRequests"`
				Email               string `json:"email"`
			} `json:"teamMemberSpend"`
			SubscriptionCycleStart int64 `json:"subscriptionCycleStart"`
			TotalPages             int   `json:"totalPages"`
		}
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, fmt.Errorf("failed to unmarshal spending response: %w", err)
		}

		dateStr := time.UnixMilli(response.SubscriptionCycleStart).Format("2006-01-02")

		for _, s := range response.TeamMemberSpend {
			allSpending = append(allSpending, SpendingData{
				MemberEmail:     s.Email,
				SpendCents:      s.SpendCents,
				PremiumRequests: s.FastPremiumRequests,
				Date:            dateStr,
			})
		}

		if page >= response.TotalPages {
			break
		}
		page++
	}

	return allSpending, nil
}

func (c *CursorClient) GetUsageEvents(userEmail string, limit int, offset int, startDate, endDate string) ([]UsageEvent, error) {
	var allEvents []UsageEvent
	page := 1
	for {
		reqBody := struct {
			Email     *string `json:"email,omitempty"`
			StartDate *int64  `json:"startDate,omitempty"`
			EndDate   *int64  `json:"endDate,omitempty"`
			Page      int     `json:"page"`
			PageSize  int     `json:"pageSize"`
		}{
			Page:     page,
			PageSize: limit,
		}
		if userEmail != "" {
			reqBody.Email = &userEmail
		}
		if startDate != "" {
			sT, err := time.Parse("2006-01-02", startDate)
			if err == nil {
				ms := sT.UnixMilli()
				reqBody.StartDate = &ms
			}
		}
		if endDate != "" {
			eT, err := time.Parse("2006-01-02", endDate)
			if err == nil {
				ms := eT.Add(24*time.Hour - time.Millisecond).UnixMilli()
				reqBody.EndDate = &ms
			}
		}
		reqJson, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}

		body, err := c.makeRequest("POST", "/teams/filtered-usage-events", nil, bytes.NewReader(reqJson))
		if err != nil {
			return nil, fmt.Errorf("failed to get usage events: %w", err)
		}

		var response struct {
			UsageEvents []struct {
				Timestamp        string `json:"timestamp"`
				Model            string `json:"model"`
				KindLabel        string `json:"kindLabel"`
				TokenUsage       *struct {
					InputTokens      int `json:"inputTokens"`
					OutputTokens     int `json:"outputTokens"`
					CacheWriteTokens int `json:"cacheWriteTokens"`
					CacheReadTokens  int `json:"cacheReadTokens"`
				} `json:"tokenUsage"`
				UserEmail        string `json:"userEmail"`
			} `json:"usageEvents"`
			Pagination struct {
				HasNextPage bool `json:"hasNextPage"`
			} `json:"pagination"`
		}
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, fmt.Errorf("failed to unmarshal usage events response: %w", err)
		}

		for _, e := range response.UsageEvents {
			tsMs, perr := strconv.ParseInt(e.Timestamp, 10, 64)
			if perr != nil {
				logrus.WithError(perr).Warn("Failed to parse timestamp")
				continue
			}
			ts := time.UnixMilli(tsMs)

			tokens := 0
			if e.TokenUsage != nil {
				tokens = e.TokenUsage.InputTokens + e.TokenUsage.OutputTokens + e.TokenUsage.CacheReadTokens + e.TokenUsage.CacheWriteTokens
			}

			allEvents = append(allEvents, UsageEvent{
				EventType:      e.KindLabel,
				UserEmail:      e.UserEmail,
				TokensConsumed: tokens,
				Model:          e.Model,
				Timestamp:      ts,
			})
		}

		if !response.Pagination.HasNextPage {
			break
		}
		page++
	}

	return allEvents, nil
}
