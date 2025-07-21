#!/bin/bash

# Coverage validation and threshold checking script
# Validates coverage against thresholds and provides detailed reporting

set -e

# Configuration
COVERAGE_FILE="coverage.out"
COVERAGE_THRESHOLD=80.0
PACKAGE_THRESHOLD=75.0
OUTPUT_FORMAT="text"
OUTPUT_DIR="coverage-reports"
COMPARE_FILE=""
VERBOSE=false
FAIL_ON_DECREASE=false
DECREASE_THRESHOLD=0.5

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help              Show this help message"
    echo "  -v, --verbose           Enable verbose output"
    echo "  -f, --file FILE         Coverage file to check (default: coverage.out)"
    echo "  -t, --threshold NUM     Overall coverage threshold (default: 80.0)"
    echo "  -p, --package-threshold NUM  Package coverage threshold (default: 75.0)"
    echo "  -o, --output FORMAT     Output format: text,json,html,badge (default: text)"
    echo "  -d, --output-dir DIR    Output directory (default: coverage-reports)"
    echo "  -c, --compare FILE      Compare with previous coverage file"
    echo "  --fail-on-decrease      Fail if coverage decreases by more than threshold"
    echo "  --decrease-threshold NUM  Decrease threshold percentage (default: 0.5)"
    echo ""
    echo "Examples:"
    echo "  $0                      # Check coverage with defaults"
    echo "  $0 -t 85 -p 80          # Check with custom thresholds"
    echo "  $0 -o json -d reports   # Generate JSON report in reports directory"
    echo "  $0 -c old_coverage.out  # Compare with previous coverage"
}

log() {
    echo -e "${BLUE}[$(date '+%Y-%m-%d %H:%M:%S')]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[$(date '+%Y-%m-%d %H:%M:%S')]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[$(date '+%Y-%m-%d %H:%M:%S')]${NC} $1"
}

log_error() {
    echo -e "${RED}[$(date '+%Y-%m-%d %H:%M:%S')]${NC} $1"
}

check_dependencies() {
    if ! command -v go >/dev/null 2>&1; then
        log_error "Go is not installed or not in PATH"
        exit 1
    fi
    
    if ! command -v bc >/dev/null 2>&1; then
        log_error "bc is not installed (required for calculations)"
        exit 1
    fi
}

validate_coverage_file() {
    if [ ! -f "$COVERAGE_FILE" ]; then
        log_error "Coverage file not found: $COVERAGE_FILE"
        exit 1
    fi
    
    if [ ! -s "$COVERAGE_FILE" ]; then
        log_error "Coverage file is empty: $COVERAGE_FILE"
        exit 1
    fi
}

get_total_coverage() {
    local file="$1"
    go tool cover -func="$file" | grep total | awk '{print $3}' | sed 's/%//'
}

get_package_coverage() {
    local file="$1"
    # Get package-level coverage by grouping functions by package
    go tool cover -func="$file" | grep -v total | awk '{
        split($1, parts, "/")
        package = parts[length(parts)-1]
        split(package, fileparts, ":")
        packagename = fileparts[1]
        coverage = $3
        gsub(/%/, "", coverage)
        packages[packagename] = coverage
    } END {
        for (pkg in packages) {
            print pkg":"packages[pkg]
        }
    }'
}

check_overall_threshold() {
    local coverage="$1"
    local threshold="$2"
    
    if (( $(echo "$coverage >= $threshold" | bc -l) )); then
        log_success "Overall coverage: $coverage% (threshold: $threshold%)"
        return 0
    else
        log_error "Overall coverage: $coverage% is below threshold: $threshold%"
        return 1
    fi
}

check_package_thresholds() {
    local threshold="$1"
    local failed_packages=()
    local package_count=0
    local passed_count=0
    
    while IFS=: read -r package coverage; do
        if [[ -n "$package" && -n "$coverage" ]]; then
            package_count=$((package_count + 1))
            
            if (( $(echo "$coverage >= $threshold" | bc -l) )); then
                passed_count=$((passed_count + 1))
                [ "$VERBOSE" = true ] && log_success "Package $package: $coverage% (threshold: $threshold%)"
            else
                failed_packages+=("$package ($coverage%)")
                [ "$VERBOSE" = true ] && log_warning "Package $package: $coverage% is below threshold: $threshold%"
            fi
        fi
    done < <(get_package_coverage "$COVERAGE_FILE")
    
    echo ""
    echo "Package Coverage Summary:"
    echo "========================="
    echo "Total packages: $package_count"
    echo "Passed threshold: $passed_count"
    echo "Failed threshold: $((package_count - passed_count))"
    echo ""
    
    if [ ${#failed_packages[@]} -gt 0 ]; then
        log_warning "Packages below threshold ($threshold%):"
        for package in "${failed_packages[@]}"; do
            echo "  - $package"
        done
        echo ""
        return 1
    else
        log_success "All packages meet the threshold requirement"
        return 0
    fi
}

compare_coverage() {
    local current_file="$1"
    local previous_file="$2"
    
    if [ ! -f "$previous_file" ]; then
        log_error "Previous coverage file not found: $previous_file"
        return 1
    fi
    
    local current_coverage=$(get_total_coverage "$current_file")
    local previous_coverage=$(get_total_coverage "$previous_file")
    
    local difference=$(echo "$current_coverage - $previous_coverage" | bc -l)
    
    echo ""
    echo "Coverage Comparison:"
    echo "==================="
    echo "Previous coverage: $previous_coverage%"
    echo "Current coverage: $current_coverage%"
    echo "Difference: $difference%"
    echo ""
    
    if (( $(echo "$difference >= 0" | bc -l) )); then
        log_success "Coverage improved by $difference%"
        return 0
    else
        local abs_difference=$(echo "$difference * -1" | bc -l)
        
        if (( $(echo "$abs_difference <= $DECREASE_THRESHOLD" | bc -l) )); then
            log_warning "Coverage decreased by $abs_difference% (within threshold: $DECREASE_THRESHOLD%)"
            return 0
        else
            log_error "Coverage decreased by $abs_difference% (exceeds threshold: $DECREASE_THRESHOLD%)"
            return 1
        fi
    fi
}

generate_text_report() {
    echo ""
    echo "Coverage Report:"
    echo "==============="
    go tool cover -func="$COVERAGE_FILE"
    echo ""
}

generate_json_report() {
    mkdir -p "$OUTPUT_DIR"
    
    local total_coverage=$(get_total_coverage "$COVERAGE_FILE")
    local timestamp=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    
    local packages_json=$(get_package_coverage "$COVERAGE_FILE" | awk -F: '{print "{\"package\":\""$1"\",\"coverage\":"$2"}"}' | jq -s '.')
    
    cat > "$OUTPUT_DIR/coverage-check.json" << EOF
{
  "timestamp": "$timestamp",
  "coverage": {
    "total": $total_coverage,
    "threshold": $COVERAGE_THRESHOLD,
    "package_threshold": $PACKAGE_THRESHOLD,
    "status": "$(echo "$total_coverage >= $COVERAGE_THRESHOLD" | bc -l | sed 's/1/pass/;s/0/fail/')"
  },
  "packages": $packages_json
}
EOF
    
    log_success "JSON report generated: $OUTPUT_DIR/coverage-check.json"
}

generate_html_report() {
    mkdir -p "$OUTPUT_DIR"
    
    local total_coverage=$(get_total_coverage "$COVERAGE_FILE")
    local timestamp=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    
    cat > "$OUTPUT_DIR/coverage-check.html" << EOF
<!DOCTYPE html>
<html>
<head>
    <title>Coverage Check Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background: #f4f4f4; padding: 20px; border-radius: 5px; }
        .success { color: #28a745; }
        .warning { color: #ffc107; }
        .error { color: #dc3545; }
        .package { margin: 10px 0; padding: 10px; background: #f8f9fa; border-radius: 3px; }
        .threshold { font-weight: bold; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Coverage Check Report</h1>
        <p>Generated: $timestamp</p>
        <p>Total Coverage: <span class="$(echo "$total_coverage >= $COVERAGE_THRESHOLD" | bc -l | sed 's/1/success/;s/0/error/')">$total_coverage%</span></p>
        <p>Threshold: <span class="threshold">$COVERAGE_THRESHOLD%</span></p>
    </div>
    
    <h2>Package Coverage</h2>
EOF
    
    while IFS=: read -r package coverage; do
        if [[ -n "$package" && -n "$coverage" ]]; then
            local status="success"
            if (( $(echo "$coverage < $PACKAGE_THRESHOLD" | bc -l) )); then
                status="warning"
            fi
            
            echo "    <div class=\"package\">" >> "$OUTPUT_DIR/coverage-check.html"
            echo "        <strong>$package:</strong> <span class=\"$status\">$coverage%</span>" >> "$OUTPUT_DIR/coverage-check.html"
            echo "    </div>" >> "$OUTPUT_DIR/coverage-check.html"
        fi
    done < <(get_package_coverage "$COVERAGE_FILE")
    
    echo "</body></html>" >> "$OUTPUT_DIR/coverage-check.html"
    
    log_success "HTML report generated: $OUTPUT_DIR/coverage-check.html"
}

generate_badge() {
    mkdir -p "$OUTPUT_DIR"
    
    local coverage=$(get_total_coverage "$COVERAGE_FILE")
    local color="red"
    
    if (( $(echo "$coverage >= 80" | bc -l) )); then
        color="brightgreen"
    elif (( $(echo "$coverage >= 60" | bc -l) )); then
        color="yellow"
    fi
    
    local badge_url="https://img.shields.io/badge/coverage-${coverage}%25-${color}"
    
    if command -v curl >/dev/null 2>&1; then
        curl -s "$badge_url" > "$OUTPUT_DIR/coverage-badge.svg"
        log_success "Coverage badge generated: $OUTPUT_DIR/coverage-badge.svg"
    else
        log_warning "curl not found, skipping badge generation"
    fi
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            print_usage
            exit 0
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -f|--file)
            COVERAGE_FILE="$2"
            shift 2
            ;;
        -t|--threshold)
            COVERAGE_THRESHOLD="$2"
            shift 2
            ;;
        -p|--package-threshold)
            PACKAGE_THRESHOLD="$2"
            shift 2
            ;;
        -o|--output)
            OUTPUT_FORMAT="$2"
            shift 2
            ;;
        -d|--output-dir)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        -c|--compare)
            COMPARE_FILE="$2"
            shift 2
            ;;
        --fail-on-decrease)
            FAIL_ON_DECREASE=true
            shift
            ;;
        --decrease-threshold)
            DECREASE_THRESHOLD="$2"
            shift 2
            ;;
        *)
            log_error "Unknown option: $1"
            print_usage
            exit 1
            ;;
    esac
done

# Main execution
main() {
    log "Starting coverage check for cursor-admin-api-exporter"
    
    # Check dependencies
    check_dependencies
    
    # Validate coverage file
    validate_coverage_file
    
    # Get total coverage
    local total_coverage=$(get_total_coverage "$COVERAGE_FILE")
    
    # Generate reports based on output format
    case "$OUTPUT_FORMAT" in
        text)
            generate_text_report
            ;;
        json)
            generate_json_report
            ;;
        html)
            generate_html_report
            ;;
        badge)
            generate_badge
            ;;
        *)
            generate_text_report
            ;;
    esac
    
    # Check thresholds
    local overall_passed=true
    local packages_passed=true
    local comparison_passed=true
    
    if ! check_overall_threshold "$total_coverage" "$COVERAGE_THRESHOLD"; then
        overall_passed=false
    fi
    
    if ! check_package_thresholds "$PACKAGE_THRESHOLD"; then
        packages_passed=false
    fi
    
    # Compare with previous coverage if specified
    if [ -n "$COMPARE_FILE" ]; then
        if ! compare_coverage "$COVERAGE_FILE" "$COMPARE_FILE"; then
            comparison_passed=false
        fi
    fi
    
    # Final result
    if [ "$overall_passed" = true ] && [ "$packages_passed" = true ] && [ "$comparison_passed" = true ]; then
        log_success "Coverage check passed!"
        exit 0
    else
        log_error "Coverage check failed!"
        exit 1
    fi
}

# Run main function
main "$@"