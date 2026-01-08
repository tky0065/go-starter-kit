.PHONY: help build test install clean lint run

# Default target
help:
	@echo "Available targets:"
	@echo "  make build    - Build the CLI binary"
	@echo "  make test     - Run all tests"
	@echo "  make install  - Install the CLI globally"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make lint     - Run golangci-lint"
	@echo "  make run      - Run the CLI (requires PROJECT_NAME)"

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
