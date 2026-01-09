.PHONY: help build run test clean dev lint test-coverage

# Binary name
BINARY_NAME=manual-test-project

help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) ./cmd
	@echo "Build complete: $(BINARY_NAME)"

run: ## Run the application
	@echo "Running $(BINARY_NAME)..."
	@go run ./cmd

dev: ## Run the application with air for hot reload
	@echo "Starting development server with hot reload..."
	@air

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run ./...

test: ## Run tests with race detection
	@echo "Running tests..."
	@go test -v -race ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@echo "Clean complete"

docker-build: ## Build docker image
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):latest .

docker-run: ## Run docker container
	@echo "Running Docker container..."
	@docker run -p 8080:8080 $(BINARY_NAME):latest
