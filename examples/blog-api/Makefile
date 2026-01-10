.PHONY: help run build test test-coverage lint docker-build docker-run clean

help: ## Afficher l'aide
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

run: ## Lancer l'application
	@echo "🚀 Starting blog-api..."
	go run ./cmd/main.go

build: ## Build le binaire
	@echo "🔨 Building binary..."
	go build -o bin/blog-api ./cmd/main.go
	@echo "✅ Build complete: bin/blog-api"

test: ## Exécuter les tests
	@echo "🧪 Running tests..."
	go test -v ./...

test-coverage: ## Tests avec coverage
	@echo "📊 Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"

lint: ## Linter le code
	@echo "🔍 Running linter..."
	golangci-lint run ./...

docker-build: ## Build l'image Docker
	@echo "🐳 Building Docker image..."
	docker build -t blog-api:latest .
	@echo "✅ Docker image built: blog-api:latest"

docker-run: ## Lancer le conteneur Docker
	@echo "🐳 Running Docker container..."
	docker run -d --name blog-api -p 8080:8080 --env-file .env blog-api:latest

docker-compose: ## Lancer avec docker-compose
	@echo "🐳 Starting with docker-compose..."
	docker-compose up -d
	@echo "✅ Services started. API: http://localhost:8080"

docker-compose-down: ## Arrêter docker-compose
	@echo "🛑 Stopping docker-compose..."
	docker-compose down

clean: ## Nettoyer les artifacts
	@echo "🧹 Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	@echo "✅ Cleaned"

.DEFAULT_GOAL := help
