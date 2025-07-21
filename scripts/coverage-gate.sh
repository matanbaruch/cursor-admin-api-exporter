#!/bin/bash

# Coverage gate script for local testing
# Tests coverage gate logic locally similar to CI/CD

set -e

# Configuration
THRESHOLD=-0.5
BASE_BRANCH=${1:-main}
STASH_CHANGES=false
VERBOSE=false

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_usage() {
    echo "Usage: $0 [BASE_BRANCH]"
    echo ""
    echo "Options:"
    echo "  BASE_BRANCH         Base branch to compare against (default: main)"
    echo "  -h, --help          Show this help message"
    echo "  -v, --verbose       Enable verbose output"
    echo "  -s, --stash         Stash changes before testing"
    echo ""
    echo "Examples:"
    echo "  $0                  # Compare against main branch"
    echo "  $0 develop          # Compare against develop branch"
    echo "  $0 -s               # Stash changes and compare against main"
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
    if ! command -v git >/dev/null 2>&1; then
        log_error "Git is not installed or not in PATH"
        exit 1
    fi
    
    if ! command -v bc >/dev/null 2>&1; then
        log_error "bc is not installed (required for calculations)"
        exit 1
    fi
    
    if ! command -v go >/dev/null 2>&1; then
        log_error "Go is not installed or not in PATH"
        exit 1
    fi
}

check_git_status() {
    if [ "$(git rev-parse --is-inside-work-tree 2>/dev/null)" != "true" ]; then
        log_error "Not in a Git repository"
        exit 1
    fi
    
    if [ "$(git status --porcelain | wc -l)" -gt 0 ]; then
        if [ "$STASH_CHANGES" = true ]; then
            log_warning "Working directory is not clean, stashing changes..."
            git stash push -m "coverage-gate-test-$(date +%s)"
        else
            log_error "Working directory is not clean. Use -s to stash changes or commit/stash them manually."
            exit 1
        fi
    fi
}

get_current_coverage() {
    log "Getting current branch coverage..."
    
    # Make scripts executable
    chmod +x scripts/coverage.sh scripts/check-coverage.sh scripts/run-tests.sh
    
    # Run coverage on current branch
    ./scripts/coverage.sh --unit-only --format all --threshold 80 --clean
    
    if [ -f coverage.out ]; then
        local coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
        echo "$coverage"
    else
        echo "0"
    fi
}

get_base_coverage() {
    local base_branch="$1"
    
    log "Getting base branch ($base_branch) coverage..."
    
    # Check if base branch exists
    if ! git show-ref --verify --quiet "refs/heads/$base_branch" && \
       ! git show-ref --verify --quiet "refs/remotes/origin/$base_branch"; then
        log_error "Base branch '$base_branch' does not exist"
        exit 1
    fi
    
    # Get current branch name
    local current_branch=$(git rev-parse --abbrev-ref HEAD)
    
    # Switch to base branch
    git checkout "$base_branch" >/dev/null 2>&1
    
    # Run coverage on base branch
    ./scripts/coverage.sh --unit-only --format all --threshold 80 --clean
    
    local coverage
    if [ -f coverage.out ]; then
        coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
    else
        coverage="0"
    fi
    
    # Switch back to current branch
    git checkout "$current_branch" >/dev/null 2>&1
    
    echo "$coverage"
}

calculate_coverage_diff() {
    local current_coverage="$1"
    local base_coverage="$2"
    
    local diff=$(echo "$current_coverage - $base_coverage" | bc -l)
    echo "$diff"
}

check_coverage_gate() {
    local current_coverage="$1"
    local base_coverage="$2"
    local diff="$3"
    
    echo ""
    echo "Coverage Gate Results:"
    echo "======================"
    echo "Current Branch Coverage: $current_coverage%"
    echo "Base Branch Coverage: $base_coverage%"
    echo "Difference: $diff%"
    echo "Threshold: $THRESHOLD%"
    echo ""
    
    if (( $(echo "$diff < $THRESHOLD" | bc -l) )); then
        log_error "Coverage decreased by more than threshold ($THRESHOLD%)"
        echo ""
        echo "Recommendations:"
        echo "1. Add tests to cover new or modified code"
        echo "2. Ensure existing tests still pass"
        echo "3. Consider refactoring complex code"
        echo "4. Review test coverage for edge cases"
        return 1
    else
        log_success "Coverage gate passed!"
        echo ""
        if (( $(echo "$diff >= 0" | bc -l) )); then
            echo "✅ Coverage improved by $diff%"
        else
            echo "✅ Coverage decrease ($diff%) is within acceptable threshold"
        fi
        return 0
    fi
}

cleanup() {
    if [ "$STASH_CHANGES" = true ]; then
        log "Restoring stashed changes..."
        git stash pop >/dev/null 2>&1 || true
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
        -s|--stash)
            STASH_CHANGES=true
            shift
            ;;
        -*)
            log_error "Unknown option: $1"
            print_usage
            exit 1
            ;;
        *)
            BASE_BRANCH="$1"
            shift
            ;;
    esac
done

# Main execution
main() {
    log "Starting coverage gate test"
    
    # Check dependencies
    check_dependencies
    
    # Check git status
    check_git_status
    
    # Set up cleanup trap
    trap cleanup EXIT
    
    # Get current coverage
    local current_coverage=$(get_current_coverage)
    
    # Get base coverage
    local base_coverage=$(get_base_coverage "$BASE_BRANCH")
    
    # Calculate difference
    local diff=$(calculate_coverage_diff "$current_coverage" "$base_coverage")
    
    # Check coverage gate
    if check_coverage_gate "$current_coverage" "$base_coverage" "$diff"; then
        log_success "Coverage gate test completed successfully!"
        exit 0
    else
        log_error "Coverage gate test failed!"
        exit 1
    fi
}

# Run main function
main "$@"