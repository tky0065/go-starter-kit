.PHONY: help build test test-short install clean lint run smoke-test

# Default target
help:
	@echo "Available targets:"
	@echo "  make build       - Build the CLI binary"
	@echo "  make test        - Run all tests"
	@echo "  make test-short  - Run tests (skip E2E)"
	@echo "  make install     - Install the CLI globally"
	@echo "  make clean       - Remove build artifacts"
	@echo "  make lint        - Run golangci-lint"
	@echo "  make run         - Run the CLI (requires PROJECT_NAME)"
	@echo "  make smoke-test  - Run full smoke test validation"

# Build the CLI binary
build:
	@echo "Building create-go-starter..."
	@go build -o bin/create-go-starter ./cmd/create-go-starter
	@echo "✓ Binary created at bin/create-go-starter"

# Run all tests
test:
	@echo "Running tests..."
	@go test -v ./...
	@echo "✓ All tests passed"

# Install the CLI globally
install:
	@echo "Installing create-go-starter..."
	@go install ./cmd/create-go-starter
	@echo "✓ Installed to $(shell go env GOPATH)/bin/create-go-starter"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@echo "✓ Clean complete"

# Run linting
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run ./...
	@echo "✓ Linting passed"

# Run the CLI (usage: make run PROJECT_NAME=my-app)
run:
	@if [ -z "$(PROJECT_NAME)" ]; then \
		echo "Error: PROJECT_NAME not set. Usage: make run PROJECT_NAME=my-app"; \
		exit 1; \
	fi
	@go run ./cmd/create-go-starter $(PROJECT_NAME)

# Run tests (skip E2E tests for faster feedback)
test-short:
	@echo "Running tests (short mode)..."
	@go test -short -v ./cmd/create-go-starter
	@echo "✓ Tests passed"

# Run full smoke test validation
smoke-test:
	@echo "Running smoke test validation..."
	@./scripts/smoke_test.sh
	@echo "✓ Smoke test complete"

# Run smoke test without runtime (faster, no Docker required)
smoke-test-quick:
	@echo "Running quick smoke test (no runtime)..."
	@./scripts/smoke_test.sh --skip-runtime
	@echo "✓ Quick smoke test complete"
