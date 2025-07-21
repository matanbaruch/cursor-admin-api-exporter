#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_color() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Function to run unit tests
run_unit_tests() {
    print_color $GREEN "Running unit tests..."
    
    # Create coverage directory
    mkdir -p coverage
    
    # Run tests with coverage
    go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
    
    # Generate coverage report
    if [ -f coverage.out ]; then
        go tool cover -html=coverage.out -o coverage.html
        go tool cover -func=coverage.out
        
        # Calculate coverage percentage
        COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
        print_color $GREEN "Total coverage: ${COVERAGE}%"
        
        # Copy to coverage directory
        cp coverage.out coverage/coverage.out
        cp coverage.html coverage/coverage.html
    fi
    
    print_color $GREEN "Unit tests completed successfully!"
}

# Function to run integration tests
run_integration_tests() {
    print_color $GREEN "Running integration tests..."
    
    # Check if API token is set
    if [ -z "$CURSOR_API_TOKEN" ]; then
        print_color $YELLOW "CURSOR_API_TOKEN not set, skipping integration tests"
        return 0
    fi
    
    # Run integration tests
    go test -v -tags=integration ./pkg/integration_test.go
    
    print_color $GREEN "Integration tests completed successfully!"
}

# Function to run performance tests
run_performance_tests() {
    print_color $GREEN "Running performance tests..."
    
    # Run performance tests
    go test -v -tags=performance ./pkg/exporters/performance_test.go
    
    print_color $GREEN "Performance tests completed successfully!"
}

# Function to run benchmark tests
run_benchmark_tests() {
    print_color $GREEN "Running benchmark tests..."
    
    # Run benchmark tests
    go test -v -bench=. -benchmem ./...
    
    print_color $GREEN "Benchmark tests completed successfully!"
}

# Function to run all tests
run_all_tests() {
    print_color $GREEN "Running all tests..."
    
    run_unit_tests
    run_integration_tests
    run_performance_tests
    run_benchmark_tests
    
    print_color $GREEN "All tests completed successfully!"
}

# Function to show help
show_help() {
    echo "Usage: $0 [OPTION]"
    echo ""
    echo "Options:"
    echo "  unit          Run unit tests with coverage"
    echo "  integration   Run integration tests"
    echo "  performance   Run performance tests"
    echo "  benchmark     Run benchmark tests"
    echo "  all           Run all tests"
    echo "  help          Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 unit"
    echo "  $0 integration"
    echo "  $0 all"
}

# Main script logic
case "${1:-all}" in
    "unit")
        run_unit_tests
        ;;
    "integration")
        run_integration_tests
        ;;
    "performance")
        run_performance_tests
        ;;
    "benchmark")
        run_benchmark_tests
        ;;
    "all")
        run_all_tests
        ;;
    "help"|"-h"|"--help")
        show_help
        ;;
    *)
        print_color $RED "Invalid option: $1"
        show_help
        exit 1
        ;;
esac