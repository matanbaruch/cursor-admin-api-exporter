package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/matanbaruch/cursor-admin-api-exporter/pkg/exporters"
	"github.com/matanbaruch/cursor-admin-api-exporter/pkg/utils"
)

func debugLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		logrus.WithFields(logrus.Fields{
			"method":      r.Method,
			"path":        r.URL.Path,
			"remote_addr": r.RemoteAddr,
			"user_agent":  r.UserAgent(),
		}).Debug("HTTP request received")

		ww := &wrappedWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(ww, r)

		duration := time.Since(start)
		logrus.WithFields(logrus.Fields{
			"method":      r.Method,
			"path":        r.URL.Path,
			"status_code": ww.statusCode,
			"duration":    duration,
		}).Debug("HTTP request completed")
	})
}

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func main() {
	cursorAPIURL := utils.GetEnvWithDefault("CURSOR_API_URL", "https://api.cursor.com")
	cursorAPIToken := os.Getenv("CURSOR_API_TOKEN")
	listenAddr := utils.GetEnvWithDefault("LISTEN_ADDRESS", ":8080")
	metricsPath := utils.GetEnvWithDefault("METRICS_PATH", "/metrics")
	logLevel := utils.GetEnvWithDefault("LOG_LEVEL", "info")

	helpFlag := false
	for _, arg := range os.Args {
		if arg == "--help" || arg == "-h" {
			helpFlag = true
			break
		}
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.WithError(err).Warn("Invalid log level, using info")
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	if helpFlag {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  This application is configured primarily via environment variables.\n")
		fmt.Fprintf(os.Stderr, "  Key environment variables:\n")
		fmt.Fprintf(os.Stderr, "    CURSOR_API_URL: Cursor API endpoint (default: https://api.cursor.com)\n")
		fmt.Fprintf(os.Stderr, "    CURSOR_API_TOKEN: Cursor API token (required)\n")
		fmt.Fprintf(os.Stderr, "    LISTEN_ADDRESS: HTTP server listen address (default: :8080)\n")
		fmt.Fprintf(os.Stderr, "    METRICS_PATH: Metrics endpoint path (default: /metrics)\n")
		fmt.Fprintf(os.Stderr, "    LOG_LEVEL: Logging level (default: info)\n")
		fmt.Fprintf(os.Stderr, "  Use --help or -h to display this message.\n")
		os.Exit(0)
	}

	if cursorAPIToken == "" {
		logrus.Fatal("CURSOR_API_TOKEN environment variable is required")
	}

	logrus.WithFields(logrus.Fields{
		"cursor_api_url": cursorAPIURL,
		"listen_addr":    listenAddr,
		"metrics_path":   metricsPath,
		"log_level":      logLevel,
	}).Info("Starting Cursor Admin API Exporter")

	exporter := exporters.NewCursorExporter(cursorAPIURL, cursorAPIToken)

	prometheus.MustRegister(exporter)

	mux := http.NewServeMux()

	var handler http.Handler = mux
	if logrus.GetLevel() == logrus.DebugLevel {
		handler = debugLoggingMiddleware(mux)
	}

	mux.Handle(metricsPath, promhttp.Handler())

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		logrus.Debug("Health check endpoint accessed")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprintf(w, `{"status":"healthy","timestamp":"%s"}`, time.Now().Format(time.RFC3339)); err != nil {
			logrus.WithError(err).Error("Failed to write health check response")
		}
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logrus.Debug("Root endpoint accessed")
		w.Header().Set("Content-Type", "text/html")
		if _, err := fmt.Fprintf(w, `
		<html>
		<head><title>Cursor Admin API Exporter</title></head>
		<body>
		<h1>Cursor Admin API Exporter</h1>
		<p>This is a Prometheus exporter for Cursor Admin API metrics.</p>
		<ul>
		<li><a href="%s">Metrics</a></li>
		<li><a href="/health">Health Check</a></li>
		</ul>
		<h2>Available Metrics</h2>
		<ul>
		<li><strong>Team Members API:</strong> Team member counts, roles, status</li>
		<li><strong>Daily Usage API:</strong> Lines of code, suggestion acceptance, feature usage</li>
		<li><strong>Spending API:</strong> Per-member spending, premium requests</li>
		<li><strong>Usage Events API:</strong> Token consumption, model usage, granular events</li>
		</ul>
		</body>
		</html>
		`, metricsPath); err != nil {
			logrus.WithError(err).Error("Failed to write root page response")
		}
	})

	server := &http.Server{
		Addr:              listenAddr,
		Handler:           handler,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logrus.Info("Shutting down server...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logrus.WithError(err).Error("Error during server shutdown")
		}
		cancel()
	}()

	logrus.WithField("address", listenAddr).Info("Starting HTTP server")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.WithError(err).Fatal("HTTP server error")
	}

	<-ctx.Done()
	logrus.Info("Server stopped")
}
