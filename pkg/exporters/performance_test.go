//go:build performance
// +build performance

package exporters

import (
	"testing"
	"time"

	"github.com/matanbaruch/cursor-admin-api-exporter/pkg/client"
)

func TestPerformance_CursorExporter(t *testing.T) {
	// Create mock client for performance testing
	mockClient := client.NewCursorClient("http://mock-api.com", "test-token")

	// Test individual exporter performance
	t.Run("TeamMembersExporter", func(t *testing.T) {
		exporter := NewTeamMembersExporter(mockClient)

		start := time.Now()

		// Simulate multiple collections
		for i := 0; i < 100; i++ {
			ch := make(chan prometheus.Metric, 10)
			go func() {
				exporter.Collect(ch)
				close(ch)
			}()

			// Drain channel
			for range ch {
			}
		}

		duration := time.Since(start)
		t.Logf("100 collections took %v (avg: %v per collection)", duration, duration/100)

		// Performance assertion
		if duration > 10*time.Second {
			t.Errorf("Performance too slow: %v", duration)
		}
	})

	t.Run("DailyUsageExporter", func(t *testing.T) {
		exporter := NewDailyUsageExporter(mockClient)

		start := time.Now()

		for i := 0; i < 100; i++ {
			ch := make(chan prometheus.Metric, 50)
			go func() {
				exporter.Collect(ch)
				close(ch)
			}()

			for range ch {
			}
		}

		duration := time.Since(start)
		t.Logf("100 daily usage collections took %v", duration)

		if duration > 10*time.Second {
			t.Errorf("Performance too slow: %v", duration)
		}
	})

	t.Run("SpendingExporter", func(t *testing.T) {
		exporter := NewSpendingExporter(mockClient)

		start := time.Now()

		for i := 0; i < 100; i++ {
			ch := make(chan prometheus.Metric, 20)
			go func() {
				exporter.Collect(ch)
				close(ch)
			}()

			for range ch {
			}
		}

		duration := time.Since(start)
		t.Logf("100 spending collections took %v", duration)

		if duration > 10*time.Second {
			t.Errorf("Performance too slow: %v", duration)
		}
	})

	t.Run("UsageEventsExporter", func(t *testing.T) {
		exporter := NewUsageEventsExporter(mockClient)

		start := time.Now()

		for i := 0; i < 100; i++ {
			ch := make(chan prometheus.Metric, 30)
			go func() {
				exporter.Collect(ch)
				close(ch)
			}()

			for range ch {
			}
		}

		duration := time.Since(start)
		t.Logf("100 usage events collections took %v", duration)

		if duration > 10*time.Second {
			t.Errorf("Performance too slow: %v", duration)
		}
	})
}

func BenchmarkCursorExporter_Collect(b *testing.B) {
	mockClient := client.NewCursorClient("http://mock-api.com", "test-token")
	exporter := NewCursorExporter("http://mock-api.com", "test-token")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ch := make(chan prometheus.Metric, 100)
		go func() {
			exporter.Collect(ch)
			close(ch)
		}()

		for range ch {
		}
	}
}

func BenchmarkTeamMembersExporter_Collect(b *testing.B) {
	mockClient := client.NewCursorClient("http://mock-api.com", "test-token")
	exporter := NewTeamMembersExporter(mockClient)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ch := make(chan prometheus.Metric, 10)
		go func() {
			exporter.Collect(ch)
			close(ch)
		}()

		for range ch {
		}
	}
}

func BenchmarkDailyUsageExporter_Collect(b *testing.B) {
	mockClient := client.NewCursorClient("http://mock-api.com", "test-token")
	exporter := NewDailyUsageExporter(mockClient)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ch := make(chan prometheus.Metric, 50)
		go func() {
			exporter.Collect(ch)
			close(ch)
		}()

		for range ch {
		}
	}
}

func BenchmarkSpendingExporter_Collect(b *testing.B) {
	mockClient := client.NewCursorClient("http://mock-api.com", "test-token")
	exporter := NewSpendingExporter(mockClient)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ch := make(chan prometheus.Metric, 20)
		go func() {
			exporter.Collect(ch)
			close(ch)
		}()

		for range ch {
		}
	}
}

func BenchmarkUsageEventsExporter_Collect(b *testing.B) {
	mockClient := client.NewCursorClient("http://mock-api.com", "test-token")
	exporter := NewUsageEventsExporter(mockClient)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ch := make(chan prometheus.Metric, 30)
		go func() {
			exporter.Collect(ch)
			close(ch)
		}()

		for range ch {
		}
	}
}
