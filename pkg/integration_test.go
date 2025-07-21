//go:build integration
// +build integration

package pkg

import (
	"os"
	"testing"
	"time"

	"github.com/matanbaruch/cursor-admin-api-exporter/pkg/client"
	"github.com/matanbaruch/cursor-admin-api-exporter/pkg/exporters"
)

func TestIntegration_CursorAPIExporter(t *testing.T) {
	apiToken := os.Getenv("CURSOR_API_TOKEN")
	if apiToken == "" {
		t.Skip("CURSOR_API_TOKEN not set, skipping integration tests")
	}

	apiURL := os.Getenv("CURSOR_API_URL")
	if apiURL == "" {
		apiURL = "https://api.cursor.com"
	}

	// Create client
	cursorClient := client.NewCursorClient(apiURL, apiToken)

	// Test API client
	t.Run("GetTeamMembers", func(t *testing.T) {
		members, err := cursorClient.GetTeamMembers()
		if err != nil {
			t.Fatalf("Failed to get team members: %v", err)
		}
		t.Logf("Retrieved %d team members", len(members))
	})

	t.Run("GetDailyUsage", func(t *testing.T) {
		endDate := time.Now().Format("2006-01-02")
		startDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02")

		usage, err := cursorClient.GetDailyUsage(startDate, endDate)
		if err != nil {
			t.Fatalf("Failed to get daily usage: %v", err)
		}
		t.Logf("Retrieved %d daily usage records", len(usage))
	})

	t.Run("GetSpending", func(t *testing.T) {
		spending, err := cursorClient.GetSpending(100, 0)
		if err != nil {
			t.Fatalf("Failed to get spending: %v", err)
		}
		t.Logf("Retrieved %d spending records", len(spending))
	})

	t.Run("GetUsageEvents", func(t *testing.T) {
		events, err := cursorClient.GetUsageEvents("", 100, 0)
		if err != nil {
			t.Fatalf("Failed to get usage events: %v", err)
		}
		t.Logf("Retrieved %d usage events", len(events))
	})

	// Test exporters
	t.Run("ExporterMetrics", func(t *testing.T) {
		exporter := exporters.NewCursorExporter(apiURL, apiToken)

		// Test Describe
		descCh := make(chan *prometheus.Desc, 100)
		go func() {
			exporter.Describe(descCh)
			close(descCh)
		}()

		descCount := 0
		for desc := range descCh {
			if desc == nil {
				t.Error("Received nil metric description")
			}
			descCount++
		}

		if descCount == 0 {
			t.Error("No metric descriptions received")
		}

		// Test Collect
		metricCh := make(chan prometheus.Metric, 100)
		go func() {
			exporter.Collect(metricCh)
			close(metricCh)
		}()

		metricCount := 0
		for metric := range metricCh {
			if metric == nil {
				t.Error("Received nil metric")
			}
			metricCount++
		}

		if metricCount == 0 {
			t.Error("No metrics collected")
		}

		t.Logf("Collected %d metrics", metricCount)
	})
}
