# Makefile for cursor-admin-api-exporter

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=cursor-admin-api-exporter
BINARY_DIR=bin

# Build info
VERSION?=$(shell git describe --tags --always --dirty)
COMMIT=$(shell git rev-parse HEAD)
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Linker flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)"

# Platform-specific settings
PLATFORMS=linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

.PHONY: all build build-all clean test test-unit test-integration test-performance test-benchmark coverage coverage-unit coverage-integration coverage-check coverage-check-json coverage-badge coverage-ci coverage-clean lint fmt vet security docker docker-build docker-push helm-lint helm-test deps tidy help

all: clean fmt vet lint test build

build: deps
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME) -v

build-all: deps
	@echo "Building for multiple platforms..."
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*} GOARCH=$${platform#*/} $(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-$${platform%/*}-$${platform#*/} -v; \
		if [ "$${platform%/*}" = "windows" ]; then \
			mv $(BINARY_DIR)/$(BINARY_NAME)-$${platform%/*}-$${platform#*/} $(BINARY_DIR)/$(BINARY_NAME)-$${platform%/*}-$${platform#*/}.exe; \
		fi; \
	done

clean:
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)
	rm -rf coverage-reports coverage*.out coverage*.html coverage*.json coverage*.xml coverage*.lcov coverage*.svg

test: test-unit test-integration test-performance test-benchmark

test-unit:
	./scripts/run-tests.sh unit

test-integration:
	./scripts/run-tests.sh integration

test-performance:
	./scripts/run-tests.sh performance

test-benchmark:
	./scripts/run-tests.sh benchmark

coverage:
	@echo "Generating comprehensive coverage report..."
	@./scripts/coverage.sh --all-tests --format all --threshold 80

coverage-unit:
	@echo "Generating unit test coverage..."
	@./scripts/coverage.sh --unit-only --format all --threshold 80

coverage-integration:
	@echo "Generating integration test coverage..."
	@./scripts/coverage.sh --integration-only --format all --threshold 80

coverage-check:
	@echo "Checking coverage thresholds..."
	@./scripts/check-coverage.sh --threshold 80 --package-threshold 75

coverage-check-json:
	@echo "Generating JSON coverage report..."
	@./scripts/check-coverage.sh --output json --output-dir coverage-reports

coverage-badge:
	@echo "Generating coverage badge..."
	@./scripts/check-coverage.sh --output badge --output-dir coverage-reports

coverage-ci:
	@echo "Running CI coverage check..."
	@./scripts/coverage.sh --all-tests --format all --threshold 80 --clean
	@./scripts/check-coverage.sh --threshold 80 --package-threshold 75 --fail-on-decrease

coverage-clean:
	@echo "Cleaning coverage artifacts..."
	@rm -rf coverage-reports coverage*.out coverage*.html coverage*.json coverage*.xml coverage*.lcov coverage*.svg

lint:
	@echo "Running linters..."
	@golangci-lint run --timeout=5m

fmt:
	@echo "Formatting code..."
	@gofmt -s -w .

vet:
	@echo "Running go vet..."
	@go vet ./...

security:
	@echo "Running security scanner..."
	@gosec ./...

docker: docker-build

docker-build:
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):latest .

docker-push:
	@echo "Pushing Docker image..."
	@docker push $(BINARY_NAME):latest

helm-lint:
	@echo "Linting Helm chart..."
	@helm lint charts/cursor-admin-api-exporter

helm-test:
	@echo "Testing Helm chart..."
	@helm template cursor-admin-api-exporter charts/cursor-admin-api-exporter --debug

deps:
	@echo "Downloading dependencies..."
	@$(GOMOD) download

tidy:
	@echo "Tidying dependencies..."
	@$(GOMOD) tidy

install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BINARY_DIR)/$(BINARY_NAME) $(GOPATH)/bin/

uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(GOPATH)/bin/$(BINARY_NAME)

run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BINARY_DIR)/$(BINARY_NAME)

dev:
	@echo "Running in development mode..."
	@go run main.go

check-env:
	@echo "Checking environment variables..."
	@echo "CURSOR_API_TOKEN: $${CURSOR_API_TOKEN:-not set}"
	@echo "CURSOR_API_URL: $${CURSOR_API_URL:-not set}"
	@echo "LISTEN_ADDRESS: $${LISTEN_ADDRESS:-not set}"
	@echo "METRICS_PATH: $${METRICS_PATH:-not set}"
	@echo "LOG_LEVEL: $${LOG_LEVEL:-not set}"

help:
	@echo "Available targets:"
	@echo "  all             - Run fmt, vet, lint, test, and build"
	@echo "  build           - Build the binary"
	@echo "  build-all       - Build binaries for all platforms"
	@echo "  clean           - Clean build artifacts"
	@echo "  test            - Run all tests"
	@echo "  test-unit       - Run unit tests"
	@echo "  test-integration- Run integration tests"
	@echo "  test-performance- Run performance tests"
	@echo "  test-benchmark  - Run benchmark tests"
	@echo "  coverage        - Generate comprehensive test coverage report"
	@echo "  coverage-unit   - Generate unit test coverage only"
	@echo "  coverage-integration - Generate integration test coverage only"
	@echo "  coverage-check  - Check coverage thresholds"
	@echo "  coverage-check-json - Generate JSON coverage report"
	@echo "  coverage-badge  - Generate coverage badge"
	@echo "  coverage-ci     - Run CI coverage check with strict thresholds"
	@echo "  coverage-clean  - Clean coverage artifacts"
	@echo "  lint            - Run linters"
	@echo "  fmt             - Format code"
	@echo "  vet             - Run go vet"
	@echo "  security        - Run security scanner"
	@echo "  docker          - Build Docker image"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-push     - Push Docker image"
	@echo "  helm-lint       - Lint Helm chart"
	@echo "  helm-test       - Test Helm chart"
	@echo "  deps            - Download dependencies"
	@echo "  tidy            - Tidy dependencies"
	@echo "  install         - Install binary to GOPATH/bin"
	@echo "  uninstall       - Remove binary from GOPATH/bin"
	@echo "  run             - Build and run the binary"
	@echo "  dev             - Run in development mode"
	@echo "  check-env       - Check environment variables"
	@echo "  help            - Show this help message"