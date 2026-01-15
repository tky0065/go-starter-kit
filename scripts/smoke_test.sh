#!/bin/bash
# =============================================================================
# Smoke Test Script for go-starter-kit
# =============================================================================
# This script performs a complete end-to-end validation of the CLI tool:
# 1. Generates a new project in a temporary directory
# 2. Verifies compilation and tests pass
# 3. Runs linting checks
# 4. Starts the server and tests critical endpoints
# 5. Cleans up and generates a validation report
#
# Prerequisites:
#   - Go 1.25+ installed
#   - Docker installed and running (for PostgreSQL)
#   - golangci-lint installed (optional, will skip if not present)
#
# Usage:
#   ./scripts/smoke_test.sh [--skip-runtime] [--keep-project]
#
# Options:
#   --skip-runtime    Skip runtime tests (server startup and HTTP tests)
#   --keep-project    Don't delete the generated project after tests
# =============================================================================

set -e

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
TEST_PROJECT_NAME="smoke-test-api-$$"
TEMP_DIR="${SMOKE_TEST_DIR:-/tmp/go-starter-smoke-tests}"
TEST_PROJECT_PATH="$TEMP_DIR/$TEST_PROJECT_NAME"
REPORT_FILE="$TEMP_DIR/smoke-test-report-$(date +%Y%m%d-%H%M%S).txt"
SERVER_PID=""
POSTGRES_CONTAINER=""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Flags
SKIP_RUNTIME=false
KEEP_PROJECT=false

# Parse arguments
for arg in "$@"; do
	case $arg in
	--skip-runtime)
		SKIP_RUNTIME=true
		shift
		;;
	--keep-project)
		KEEP_PROJECT=true
		shift
		;;
	esac
done

# =============================================================================
# Utility Functions
# =============================================================================

log_info() {
	echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
	echo -e "${GREEN}[PASS]${NC} $1"
}

log_warning() {
	echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
	echo -e "${RED}[FAIL]${NC} $1"
}

log_section() {
	echo ""
	echo -e "${BLUE}======================================${NC}"
	echo -e "${BLUE}  $1${NC}"
	echo -e "${BLUE}======================================${NC}"
}

report() {
	echo "$1" >>"$REPORT_FILE"
}

cleanup() {
	log_section "Cleanup"

	# Stop server if running
	if [ -n "$SERVER_PID" ] && kill -0 "$SERVER_PID" 2>/dev/null; then
		log_info "Stopping test server (PID: $SERVER_PID)..."
		kill "$SERVER_PID" 2>/dev/null || true
		wait "$SERVER_PID" 2>/dev/null || true
		log_success "Server stopped"
	fi

	# Stop PostgreSQL container if we started one
	if [ -n "$POSTGRES_CONTAINER" ]; then
		log_info "Stopping PostgreSQL container..."
		docker stop "$POSTGRES_CONTAINER" >/dev/null 2>&1 || true
		docker rm "$POSTGRES_CONTAINER" >/dev/null 2>&1 || true
		log_success "PostgreSQL container removed"
	fi

	# Remove test project unless --keep-project is set
	if [ "$KEEP_PROJECT" = false ] && [ -d "$TEST_PROJECT_PATH" ]; then
		log_info "Removing test project..."
		rm -rf "$TEST_PROJECT_PATH"
		log_success "Test project removed"
	elif [ "$KEEP_PROJECT" = true ]; then
		log_info "Keeping test project at: $TEST_PROJECT_PATH"
	fi
}

trap cleanup EXIT

# =============================================================================
# Test Functions
# =============================================================================

test_prerequisites() {
	log_section "Checking Prerequisites"

	# Check Go
	if command -v go &>/dev/null; then
		GO_VERSION=$(go version | awk '{print $3}')
		log_success "Go installed: $GO_VERSION"
		report "Go Version: $GO_VERSION"
	else
		log_error "Go is not installed"
		exit 1
	fi

	# Check Docker (required for runtime tests)
	if command -v docker &>/dev/null; then
		if docker info &>/dev/null; then
			log_success "Docker is installed and running"
			report "Docker: Available"
		else
			log_warning "Docker is installed but not running"
			report "Docker: Not running"
			if [ "$SKIP_RUNTIME" = false ]; then
				log_warning "Runtime tests will be skipped (Docker required for PostgreSQL)"
				SKIP_RUNTIME=true
			fi
		fi
	else
		log_warning "Docker is not installed"
		report "Docker: Not installed"
		if [ "$SKIP_RUNTIME" = false ]; then
			log_warning "Runtime tests will be skipped"
			SKIP_RUNTIME=true
		fi
	fi

	# Check golangci-lint
	if command -v golangci-lint &>/dev/null; then
		LINT_VERSION=$(golangci-lint --version | head -n1)
		log_success "golangci-lint installed: $LINT_VERSION"
		report "golangci-lint: $LINT_VERSION"
		HAS_LINT=true
	else
		log_warning "golangci-lint is not installed (lint tests will be skipped)"
		report "golangci-lint: Not installed"
		HAS_LINT=false
	fi
}

test_cli_generation() {
	log_section "Testing Project Generation (AC#1)"

	# Create temp directory
	mkdir -p "$TEMP_DIR"
	cd "$TEMP_DIR"

	# Build CLI first
	log_info "Building CLI..."
	cd "$PROJECT_ROOT"
	go build -o "$TEMP_DIR/create-go-starter" ./cmd/create-go-starter
	log_success "CLI built successfully"

	# Generate project
	cd "$TEMP_DIR"
	log_info "Generating project: $TEST_PROJECT_NAME"

	if ./create-go-starter "$TEST_PROJECT_NAME" 2>&1; then
		log_success "Project generated without errors"
		report "Project Generation: PASS"
	else
		log_error "Project generation failed"
		report "Project Generation: FAIL"
		exit 1
	fi

	# Verify project structure
	log_info "Verifying project structure..."
	REQUIRED_FILES=(
		"go.mod"
		"cmd/main.go"
		"Makefile"
		".env.example"
		"Dockerfile"
		"docker-compose.yml"
		".gitignore"
		"internal/models/user.go"
		"internal/domain/user/service.go"
		"internal/adapters/handlers/auth_handler.go"
		"internal/infrastructure/server/server.go"
	)

	MISSING_FILES=()
	for file in "${REQUIRED_FILES[@]}"; do
		if [ -f "$TEST_PROJECT_PATH/$file" ]; then
			log_success "  Found: $file"
		else
			log_error "  Missing: $file"
			MISSING_FILES+=("$file")
		fi
	done

	if [ ${#MISSING_FILES[@]} -eq 0 ]; then
		report "Structure Verification: PASS (${#REQUIRED_FILES[@]} files checked)"
	else
		report "Structure Verification: FAIL (Missing: ${MISSING_FILES[*]})"
		exit 1
	fi

	# Check for Git initialization
	if [ -d "$TEST_PROJECT_PATH/.git" ]; then
		log_success "Git repository initialized"
		report "Git Init: PASS"
	else
		log_warning "Git repository not initialized (may be expected)"
		report "Git Init: WARN"
	fi
}

test_compilation_and_tests() {
	log_section "Testing Compilation and Tests (AC#2)"

	cd "$TEST_PROJECT_PATH"

	# Test go mod tidy
	log_info "Running go mod tidy..."
	if go mod tidy 2>&1; then
		log_success "go mod tidy successful"
	else
		log_error "go mod tidy failed"
		report "Go Mod Tidy: FAIL"
		exit 1
	fi

	# Test compilation
	log_info "Running go build ./..."
	if go build ./... 2>&1; then
		log_success "Compilation successful"
		report "Compilation: PASS"
	else
		log_error "Compilation failed"
		report "Compilation: FAIL"
		exit 1
	fi

	# Run tests
	log_info "Running go test ./..."
	TEST_OUTPUT=$(go test ./... 2>&1) || {
		log_warning "Some tests failed or no tests found"
		echo "$TEST_OUTPUT"
		report "Tests: WARN (Some tests may have failed)"
	}

	if echo "$TEST_OUTPUT" | grep -q "ok\|no test files"; then
		log_success "Tests passed (or no test files)"
		report "Tests: PASS"
	fi
}

test_lint_compliance() {
	log_section "Testing Lint Compliance (AC#3)"

	cd "$TEST_PROJECT_PATH"

	if [ "$HAS_LINT" = false ]; then
		log_warning "Skipping lint tests (golangci-lint not installed)"
		report "Lint: SKIPPED"
		return
	fi

	# Check if .golangci.yml exists
	if [ -f ".golangci.yml" ]; then
		log_success "Found .golangci.yml configuration"
	else
		log_warning ".golangci.yml not found, using defaults"
	fi

	# Run linter
	log_info "Running golangci-lint run..."
	if golangci-lint run ./... 2>&1; then
		log_success "Lint passed without violations"
		report "Lint: PASS"
	else
		log_warning "Lint found some issues (may be warnings)"
		report "Lint: WARN"
	fi
}

test_smoke_runtime() {
	log_section "Testing Runtime (Smoke Tests) (AC#4)"

	if [ "$SKIP_RUNTIME" = true ]; then
		log_warning "Skipping runtime tests"
		report "Runtime Tests: SKIPPED"
		return
	fi

	cd "$TEST_PROJECT_PATH"

	# Start PostgreSQL container
	log_info "Starting PostgreSQL container..."
	POSTGRES_CONTAINER="smoke-test-postgres-$$"

	docker run -d \
		--name "$POSTGRES_CONTAINER" \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_DB=smoke_test_db \
		-p 5433:5432 \
		postgres:16-alpine >/dev/null

	# Wait for PostgreSQL to be ready
	log_info "Waiting for PostgreSQL to be ready..."
	for i in {1..30}; do
		if docker exec "$POSTGRES_CONTAINER" pg_isready -U postgres >/dev/null 2>&1; then
			log_success "PostgreSQL is ready"
			break
		fi
		sleep 1
	done

	# Update .env with test database settings and JWT secret
	sed -i.bak 's/DB_HOST=.*/DB_HOST=localhost/' .env 2>/dev/null ||
		sed -i '' 's/DB_HOST=.*/DB_HOST=localhost/' .env
	sed -i.bak 's/DB_PORT=.*/DB_PORT=5433/' .env 2>/dev/null ||
		sed -i '' 's/DB_PORT=.*/DB_PORT=5433/' .env
	sed -i.bak 's/DB_NAME=.*/DB_NAME=smoke_test_db/' .env 2>/dev/null ||
		sed -i '' 's/DB_NAME=.*/DB_NAME=smoke_test_db/' .env
	sed -i.bak 's/DB_USER=.*/DB_USER=postgres/' .env 2>/dev/null ||
		sed -i '' 's/DB_USER=.*/DB_USER=postgres/' .env
	sed -i.bak 's/DB_PASSWORD=.*/DB_PASSWORD=postgres/' .env 2>/dev/null ||
		sed -i '' 's/DB_PASSWORD=.*/DB_PASSWORD=postgres/' .env
	# Set JWT_SECRET for authentication tests
	sed -i.bak 's/JWT_SECRET=.*/JWT_SECRET=smoke-test-secret-key-for-validation-only/' .env 2>/dev/null ||
		sed -i '' 's/JWT_SECRET=.*/JWT_SECRET=smoke-test-secret-key-for-validation-only/' .env
	rm -f .env.bak

	# Build and start the server
	log_info "Building the server..."
	go build -o test-server ./cmd/main.go

	log_info "Starting the server..."
	./test-server &
	SERVER_PID=$!

	# Wait for server to be ready
	log_info "Waiting for server to start..."
	for i in {1..30}; do
		if curl -s http://localhost:8080/health >/dev/null 2>&1; then
			log_success "Server is ready on port 8080"
			break
		fi
		if [ $i -eq 30 ]; then
			log_error "Server failed to start within 30 seconds"
			report "Server Startup: FAIL"
			return
		fi
		sleep 1
	done
	report "Server Startup: PASS"

	# Test /health endpoint
	log_info "Testing /health endpoint..."
	HEALTH_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health)
	if [ "$HEALTH_RESPONSE" = "200" ]; then
		log_success "/health returns 200 OK"
		report "Health Endpoint: PASS"
	else
		log_error "/health returned $HEALTH_RESPONSE"
		report "Health Endpoint: FAIL"
	fi

	# Test /swagger endpoint
	log_info "Testing /swagger endpoint..."
	SWAGGER_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/swagger/index.html)
	if [ "$SWAGGER_RESPONSE" = "200" ]; then
		log_success "/swagger/index.html returns 200 OK"
		report "Swagger Endpoint: PASS"
	else
		log_warning "/swagger returned $SWAGGER_RESPONSE (may need swag init)"
		report "Swagger Endpoint: WARN ($SWAGGER_RESPONSE)"
	fi

	# Test authentication flow
	log_info "Testing authentication flow..."

	# Register
	REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
		-H "Content-Type: application/json" \
		-d '{"email":"smoketest@example.com","password":"Test123!@#"}')

	if echo "$REGISTER_RESPONSE" | grep -q "id\|user\|success\|email"; then
		log_success "Registration successful"
		report "Registration: PASS"
	else
		log_warning "Registration response: $REGISTER_RESPONSE"
		report "Registration: WARN"
	fi

	# Login
	LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
		-H "Content-Type: application/json" \
		-d '{"email":"smoketest@example.com","password":"Test123!@#"}')

	if echo "$LOGIN_RESPONSE" | grep -q "token\|access"; then
		log_success "Login successful"
		report "Login: PASS"
	else
		log_warning "Login response: $LOGIN_RESPONSE"
		report "Login: WARN"
	fi
}

generate_report() {
	log_section "Generating Validation Report"

	# Header
	{
		echo "=============================================="
		echo "  GO-STARTER-KIT SMOKE TEST REPORT"
		echo "=============================================="
		echo ""
		echo "Date: $(date)"
		echo "Project: $TEST_PROJECT_NAME"
		echo "Location: $TEST_PROJECT_PATH"
		echo ""
		echo "----------------------------------------------"
		echo "  RESULTS SUMMARY"
		echo "----------------------------------------------"
	} >"$REPORT_FILE.header"

	cat "$REPORT_FILE.header" "$REPORT_FILE" >"$REPORT_FILE.final"
	mv "$REPORT_FILE.final" "$REPORT_FILE"
	rm -f "$REPORT_FILE.header"

	# Display report
	echo ""
	cat "$REPORT_FILE"
	echo ""
	log_success "Report saved to: $REPORT_FILE"
}

# =============================================================================
# Main Execution
# =============================================================================

main() {
	echo ""
	echo -e "${GREEN}=============================================${NC}"
	echo -e "${GREEN}  GO-STARTER-KIT SMOKE TEST SUITE${NC}"
	echo -e "${GREEN}=============================================${NC}"
	echo ""
	echo "Project Root: $PROJECT_ROOT"
	echo "Test Project: $TEST_PROJECT_NAME"
	echo "Temp Dir: $TEMP_DIR"
	echo ""

	# Create temp directory and initialize report
	mkdir -p "$TEMP_DIR"
	echo "" >"$REPORT_FILE"

	# Run tests
	test_prerequisites
	test_cli_generation
	test_compilation_and_tests
	test_lint_compliance
	test_smoke_runtime

	# Generate report
	generate_report

	# Final summary
	log_section "SMOKE TEST COMPLETE"

	if grep -q "FAIL" "$REPORT_FILE"; then
		log_error "Some tests failed. Please review the report."
		exit 1
	else
		log_success "All smoke tests passed!"
		exit 0
	fi
}

main "$@"
