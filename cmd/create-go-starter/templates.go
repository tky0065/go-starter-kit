package main

// ProjectTemplates holds all the templates for project file generation
type ProjectTemplates struct {
	projectName string
	database    string // Database type: mysql, postgres (default), sqlite, mongodb
}

// NewProjectTemplates creates a new templates instance with the given project name
// Defaults to postgres for backward compatibility
func NewProjectTemplates(projectName string) *ProjectTemplates {
	return &ProjectTemplates{
		projectName: projectName,
		database:    "postgres", // Default to postgres for backward compatibility
	}
}

// NewProjectTemplatesWithDatabase creates a new templates instance with project name and database type
func NewProjectTemplatesWithDatabase(projectName, database string) *ProjectTemplates {
	if database == "" {
		database = "postgres" // Default to postgres if empty
	}
	return &ProjectTemplates{
		projectName: projectName,
		database:    database,
	}
}

// Helper methods for database-specific configuration

// getDatabaseDescription returns a description for the database type
func (t *ProjectTemplates) getDatabaseDescription() string {
	switch t.database {
	case "mysql":
		return "MySQL"
	case "sqlite":
		return "SQLite"
	case "postgres", "":
		return "PostgreSQL"
	default:
		return "PostgreSQL"
	}
}

// getDefaultPort returns the default port for the database type
func (t *ProjectTemplates) getDefaultPort() string {
	switch t.database {
	case "mysql":
		return "3306"
	case "sqlite":
		return ""
	case "postgres", "":
		return "5432"
	default:
		return "5432"
	}
}

// getDefaultUser returns the default user for the database type
func (t *ProjectTemplates) getDefaultUser() string {
	switch t.database {
	case "mysql":
		return "root"
	case "sqlite":
		return ""
	case "postgres", "":
		return "postgres"
	default:
		return "postgres"
	}
}

// getDefaultPassword returns the default password for the database type
func (t *ProjectTemplates) getDefaultPassword() string {
	switch t.database {
	case "mysql":
		return "secret"
	case "sqlite":
		return ""
	case "postgres", "":
		return "postgres"
	default:
		return "postgres"
	}
}

// getDatabaseConnectionCode returns the GORM connection code for the database type
func (t *ProjectTemplates) getDatabaseConnectionCode() string {
	switch t.database {
	case "mysql":
		return "mysql.Open(dsn)"
	case "sqlite":
		return "sqlite.Open(dsn)"
	case "postgres", "":
		return "postgres.Open(dsn)"
	default:
		return "postgres.Open(dsn)"
	}
}

// getDatabaseDSN returns the DSN format string for the database type
func (t *ProjectTemplates) getDatabaseDSN() string {
	switch t.database {
	case "mysql":
		return `"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"`
	case "sqlite":
		return `"./%s.db"`
	case "postgres", "":
		return `"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s"`
	default:
		return `"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s"`
	}
}

func (t *ProjectTemplates) GoModTemplate() string {
	// Get the appropriate database driver dependency
	databaseDriver := GoModDependencies(t.database)

	return `module ` + t.projectName + `

go 1.25.5

require (
	github.com/go-playground/validator/v10 v10.30.1
	github.com/gofiber/contrib/jwt v1.1.2
	github.com/gofiber/fiber/v2 v2.52.10
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/joho/godotenv v1.5.1
	github.com/rs/zerolog v1.33.0
	github.com/swaggo/fiber-swagger v1.3.0
	github.com/swaggo/swag v1.16.4
	go.uber.org/fx v1.24.0
	golang.org/x/crypto v0.32.0
	` + databaseDriver + `
	gorm.io/gorm v1.31.1
)
`
}

// MainGoTemplate returns the main.go file content
func (t *ProjectTemplates) MainGoTemplate() string {
	return `package main

import (
	"fmt"
)

func main() {
	// TODO: This is a placeholder main.go
	// Infrastructure components (Fiber, GORM, fx) will be added in Story 1.4
	fmt.Println("` + t.projectName + ` - Project structure initialized")
	fmt.Println("Next steps:")
	fmt.Println("  1. Run 'go mod tidy' to fetch dependencies")
	fmt.Println("  2. Implement your application logic")
	fmt.Println("  3. Run 'make build' to build the binary")
}
`
}

// DockerfileTemplate returns the Dockerfile content
// Optimized for minimal image size (<50MB) with security best practices
func (t *ProjectTemplates) DockerfileTemplate() string {
	return `# =============================================================================
# Build stage - Compile the Go application
# =============================================================================
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install ca-certificates for HTTPS and git for private modules (if needed)
RUN apk --no-cache add ca-certificates

# Copy go mod files first for better layer caching
COPY go.mod ./

# Download dependencies and generate go.sum
RUN go mod download

# Copy source code
COPY . .

# Run go mod tidy to ensure all dependencies are resolved
RUN go mod tidy

# Build a statically linked binary with optimized flags
# -s: Omit the symbol table and debug information
# -w: Omit the DWARF symbol table
# CGO_ENABLED=0: Disable cgo for a fully static binary (required for scratch/alpine)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o ` + t.projectName + ` ./cmd

# =============================================================================
# Runtime stage - Minimal production image
# =============================================================================
FROM alpine:3.21

# Add ca-certificates for HTTPS requests and wget for healthcheck
RUN apk --no-cache add ca-certificates wget

# Create non-root user for security (AC #2)
# -D: No password, -g: GECOS, -s: Shell, -H: No home directory
RUN addgroup -g 1000 -S appgroup && \
    adduser -u 1000 -S appuser -G appgroup -s /sbin/nologin -H

# Set working directory
WORKDIR /app

# Copy the binary from builder with proper ownership
COPY --from=builder --chown=appuser:appgroup /app/` + t.projectName + ` .

# Copy ca-certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Switch to non-root user
USER appuser

# Expose application port
EXPOSE 8080

# Healthcheck to monitor application status (AC #4)
# Check /health endpoint every 30s, timeout 3s, start after 5s, fail after 3 retries
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the binary
CMD ["./` + t.projectName + `"]
`
}

// GolangCILintTemplate returns the .golangci.yml file content
// Compatible with golangci-lint v1.x (widely deployed)
func (t *ProjectTemplates) GolangCILintTemplate() string {
	return `run:
  timeout: 5m
  tests: false # Don't lint test files strictly

linters:
  disable-all: true
  enable:
    - errcheck      # Check for unchecked errors
    - gosimple      # Simplify code
    - govet         # Vet examines Go source code
    - ineffassign   # Detect ineffectual assignments
    - staticcheck   # Advanced Go linter
    - typecheck     # Type-check Go code
    - unused        # Check for unused constants, variables, functions and types
    - gocyclo       # Compute cyclomatic complexities
    - gofmt         # Check formatting
    - gosec         # Security-focused linter (basic)

linters-settings:
  gocyclo:
    min-complexity: 15
  gosec:
    excludes:
      - G404 # Allow weak random number generator in non-crypto contexts
`
}

// MakefileTemplate returns the Makefile content
func (t *ProjectTemplates) MakefileTemplate() string {
	return `.PHONY: help build run test clean dev lint test-coverage

# Binary name
BINARY_NAME=` + t.projectName + `

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

swagger: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/main.go --output docs
	@echo "Swagger documentation generated in docs/ directory"
	@echo "Run the application and visit http://localhost:8080/swagger/index.html"

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
`
}

// EnvTemplate returns the .env.example file content
func (t *ProjectTemplates) EnvTemplate() string {
	return `# Application Configuration
APP_NAME=` + t.projectName + `
APP_ENV=development
APP_PORT=8080

` + DatabaseEnvVars(t.database, t.projectName) + `

# JWT Configuration
# IMPORTANT: Generate a secure random secret for production!
# Example: openssl rand -base64 32
JWT_SECRET=
JWT_EXPIRY=24h
`
}

// GitignoreTemplate returns the .gitignore file content
func (t *ProjectTemplates) GitignoreTemplate() string {
	return `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib
` + t.projectName + `

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool
*.out

# Dependency directories
vendor/

# Go workspace file
go.work

# Environment files
.env
.env.local

# IDE files
.vscode/
.idea/
*.swp
*.swo
*~

# OS files
.DS_Store
Thumbs.db

# Temporary files
tmp/
temp/

# SQLite database files
*.db
*.db-shm
*.db-wal
`
}

// DockerComposeTemplate returns the docker-compose.yml file content
func (t *ProjectTemplates) DockerComposeTemplate() string {
	// Get database-specific environment variables for API service
	dbHost := "db"
	dbPort := "5432"
	dbUser := "postgres"
	dbPassword := "postgres"
	dbName := t.projectName

	// Adjust for MySQL
	if t.database == "mysql" {
		dbPort = "3306"
		dbUser = "root"
		dbPassword = "secret"
	}

	// For SQLite, no db service needed
	if t.database == "sqlite" {
		return `version: '3.8'

services:
  # Application API (SQLite - no separate DB service needed)
  api:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: ` + t.projectName + `_api
    environment:
      APP_NAME: ` + t.projectName + `
      APP_ENV: development
      APP_PORT: 8080
      DB_NAME: ` + dbName + `.db
      JWT_SECRET: dev-secret-change-in-production
      JWT_EXPIRY: 24h
    ports:
      - "8080:8080"
    networks:
      - ` + t.projectName + `_network
    volumes:
      - .:/app
      - sqlite_data:/app/data
    command: /app/` + t.projectName + `

volumes:
  sqlite_data:

networks:
  ` + t.projectName + `_network:
    driver: bridge
`
	}

	// For MySQL and PostgreSQL, include database service
	dbServiceComment := "# PostgreSQL Database"
	if t.database == "mysql" {
		dbServiceComment = "# MySQL Database"
	}

	// Get database service configuration from helper function
	dbServiceConfig := DatabaseDockerService(t.database, t.projectName)

	return `version: '3.8'

services:
  ` + dbServiceComment + `
` + dbServiceConfig + `
    networks:
      - ` + t.projectName + `_network

  # Application API
  api:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: ` + t.projectName + `_api
    environment:
      APP_NAME: ` + t.projectName + `
      APP_ENV: development
      APP_PORT: 8080
      DB_HOST: ` + dbHost + `
      DB_PORT: ` + dbPort + `
      DB_USER: ` + dbUser + `
      DB_PASSWORD: ` + dbPassword + `
      DB_NAME: ` + dbName + `
      DB_SSLMODE: disable
      JWT_SECRET: dev-secret-change-in-production
      JWT_EXPIRY: 24h
    ports:
      - "8080:8080"
    depends_on:
      db:
        condition: service_healthy
    networks:
      - ` + t.projectName + `_network
    volumes:
      - .:/app
    command: /app/` + t.projectName + `

networks:
  ` + t.projectName + `_network:
    driver: bridge
`
}

// ReadmeTemplate returns the README.md file content
func (t *ProjectTemplates) ReadmeTemplate() string {
	dbDesc := t.readmeDatabaseDescription()
	dbPrereqs := t.readmeDatabasePrerequisites()
	return `# ` + t.projectName + `

Application backend Go générée avec create-go-starter. Architecture hexagonale complète avec authentification JWT, API REST, et intégration ` + dbDesc + `.

## Fonctionnalités

- **Architecture hexagonale** (Ports & Adapters) - Séparation claire des responsabilités
- **Authentification JWT** - Access tokens + Refresh tokens avec rotation sécurisée
- **API REST** avec Fiber v2 - Framework web haute performance
- **Base de données** - GORM avec ` + dbDesc + ` et migrations automatiques
- **Injection de dépendances** - uber-go/fx pour architecture modulaire
- **Tests complets** - Tests unitaires et d'intégration
- **Documentation Swagger** - API documentée automatiquement avec OpenAPI
- **Docker** - Build multi-stage optimisé
- **CI/CD** - Pipeline GitHub Actions pré-configuré
- **Logging structuré** - rs/zerolog pour logs professionnels

## Prérequis

- ` + dbPrereqs + `
  ` + "```bash" + `
  go install github.com/swaggo/swag/cmd/swag@latest
  ` + "```" + `

## Installation rapide

### 1. Installer les dépendances

` + "```bash" + `
go mod tidy
` + "```" + `

### 2. Configurer l'environnement

Le fichier ` + "`.env`" + ` a déjà été créé depuis ` + "`.env.example`" + `. Éditez-le pour ajouter votre JWT secret:

` + "```bash" + `
# Générer un JWT secret sécurisé
openssl rand -base64 32

# Éditer .env et ajouter le secret
nano .env
` + "```" + `

Ajoutez dans ` + "`.env`" + `:
` + "```" + `
JWT_SECRET=<votre_secret_généré>
` + "```" + `

` + t.readmeDatabaseSetup() + `

L'API sera disponible sur ` + "`http://localhost:8080`" + `

### 5. Tester

` + "```bash" + `
# Health check
curl http://localhost:8080/health

# Register un utilisateur
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
` + "```" + `

## Documentation

Pour plus de détails, consultez la documentation complète dans le dossier ` + "`docs/`" + `:

- **[Quick Start](./docs/quick-start.md)** - Démarrage en 5 minutes
- **[Documentation complète](./docs/)** - Guides complets

## Architecture

Ce projet suit l'architecture hexagonale (Ports and Adapters):

` + "```" + `
` + t.projectName + `/
├── cmd/                     # Point d'entrée
│   └── main.go              # Bootstrap avec fx
├── internal/
│   ├── domain/              # Logique métier (cœur)
│   │   ├── user/            # Domaine User
│   │   │   ├── entity.go    # Entités
│   │   │   └── service.go   # Logique métier
│   │   └── errors.go        # Erreurs métier
│   ├── adapters/            # Adapters (HTTP, DB)
│   │   ├── handlers/        # HTTP handlers
│   │   ├── middleware/      # Middleware Fiber
│   │   └── repository/      # Implémentation GORM
│   ├── infrastructure/      # Infrastructure
│   │   ├── database/        # Configuration DB
│   │   └── server/          # Configuration Fiber
│   └── interfaces/          # Ports (interfaces)
├── pkg/                     # Packages réutilisables
│   ├── auth/                # JWT utilities
│   ├── config/              # Configuration
│   └── logger/              # Logger
├── .env                     # Configuration (créé automatiquement)
├── .env.example             # Template
├── Dockerfile               # Build Docker
├── Makefile                 # Commandes
└── go.mod                   # Dépendances
` + "```" + `

**Principe**: Le domaine (` + "`internal/domain`" + `) ne dépend de rien. Toutes les dépendances pointent vers le domaine via des interfaces (` + "`internal/interfaces`" + `).

## API Endpoints

### Authentication (Public)

- ` + "`POST /api/v1/auth/register`" + ` - Créer un compte
- ` + "`POST /api/v1/auth/login`" + ` - Se connecter
- ` + "`POST /api/v1/auth/refresh`" + ` - Rafraîchir le token

### Users (Protected - JWT required)

- ` + "`GET /api/v1/users`" + ` - Liste des utilisateurs
- ` + "`GET /api/v1/users/:id`" + ` - Détails d'un utilisateur
- ` + "`PUT /api/v1/users/:id`" + ` - Mettre à jour
- ` + "`DELETE /api/v1/users/:id`" + ` - Supprimer (soft delete)

### Health

- ` + "`GET /health`" + ` - Health check

## Développement

### Commandes Make

| Commande | Description |
|----------|-------------|
| ` + "`make help`" + ` | Afficher l'aide |
| ` + "`make run`" + ` | Lancer l'application |
| ` + "`make build`" + ` | Compiler le binaire |
| ` + "`make test`" + ` | Tests avec race detector |
| ` + "`make test-coverage`" + ` | Tests + rapport HTML |
| ` + "`make lint`" + ` | golangci-lint |
| ` + "`make clean`" + ` | Nettoyer artifacts |
| ` + "`make docker-build`" + ` | Build image Docker |
| ` + "`make docker-run`" + ` | Run conteneur Docker |

### Tests

` + "```bash" + `
# Tous les tests
make test

# Tests avec coverage
make test-coverage

# Ouvrir le rapport
open coverage.html  # macOS
xdg-open coverage.html  # Linux
` + "```" + `

### Linting

` + "```bash" + `
make lint
` + "```" + `

## Stack technique

| Composant | Bibliothèque | Description |
|-----------|-------------|-------------|
| Web Framework | [Fiber](https://gofiber.io/) v2 | Framework HTTP rapide |
| ORM | [GORM](https://gorm.io/) | ` + t.readmeDatabaseStackDescription() + ` |
| DI | [fx](https://uber-go.github.io/fx/) | Dependency injection |
| Logging | [zerolog](https://github.com/rs/zerolog) | Logger structuré |
| JWT | [golang-jwt](https://github.com/golang-jwt/jwt) v5 | Authentification |
| Validation | [validator](https://github.com/go-playground/validator) v10 | Validation |
| Swagger | [swaggo](https://github.com/swaggo/swag) | Documentation API |

## Variables d'environnement

Fichier ` + "`.env`" + `:

` + "```bash" + `
` + t.readmeDatabaseConfig() + `
` + "```" + `

## Déploiement

### Docker

` + "```bash" + `
` + t.readmeDatabaseDockerRun() + `
` + "```" + `

### Docker Compose

Si disponible:

` + "```bash" + `
docker-compose up -d
` + "```" + `

## Contribuer

1. Fork le projet
2. Créer une branche (` + "`git checkout -b feature/ma-fonctionnalite`" + `)
3. Commit (` + "`git commit -m 'feat: ajouter fonctionnalité'`" + `)
4. Push (` + "`git push origin feature/ma-fonctionnalite`" + `)
5. Ouvrir une Pull Request

## Sécurité

- ✅ JWT avec secrets forts
- ✅ Passwords hashés avec bcrypt
- ✅ Validation des entrées
- ✅ Soft deletes
- ✅ GORM prévient SQL injection
- ✅ Error handling centralisé

**Production checklist**:
- [ ] Générer JWT_SECRET fort (` + "`openssl rand -base64 32`" + `)
- [ ] HTTPS/TLS activé
- [ ] DB_SSLMODE=require
- [ ] Rate limiting configuré
- [ ] CORS configuré
- [ ] Secrets dans gestionnaire de secrets

## Licence

MIT

---

**Généré avec [create-go-starter](https://github.com/tky0065/go-starter-kit)** 🚀
`
}

// LoggerTemplate returns the pkg/logger/logger.go file content
func (t *ProjectTemplates) LoggerTemplate() string {
	return `// Package logger provides structured logging utilities using zerolog.
// It configures the logger based on the application environment, using JSON format
// in production for log aggregation systems and console format in development
// for human readability. The logger is provided via fx for dependency injection.
package logger

import (
	"os"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

// Module provides the logger dependency via fx for application-wide logging.
var Module = fx.Module("logger",
	fx.Provide(NewLogger),
)

// NewLogger creates a new zerolog logger instance configured for the current environment.
// In production (APP_ENV=production), it outputs JSON format for log aggregation.
// In other environments, it uses a human-readable console format with colors.
func NewLogger() zerolog.Logger {
	// Use JSON format in production, console format in development
	env := os.Getenv("APP_ENV")

	var logger zerolog.Logger
	if env == "production" {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	}

	return logger
}
`
}

// DatabaseTemplate returns the internal/infrastructure/database/database.go file content
func (t *ProjectTemplates) DatabaseTemplate() string {
	// Get database-specific import path and description
	databaseImport := DatabaseImportPath(t.database)
	databaseDesc := t.getDatabaseDescription()
	connectionCode := DatabaseConnectionCode(t.database)

	// Generate database-specific DSN code
	var dsnCode string
	switch t.database {
	case "mysql":
		dsnCode = `dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.GetEnv("DB_USER", "root"),
		config.GetEnv("DB_PASSWORD", "secret"),
		config.GetEnv("DB_HOST", "localhost"),
		config.GetEnv("DB_PORT", "3306"),
		config.GetEnv("DB_NAME", "` + t.projectName + `"),
	)`
	case "sqlite":
		dsnCode = `dsn := fmt.Sprintf(
		"./%s.db",
		config.GetEnv("DB_NAME", "` + t.projectName + `"),
	)`
	default: // postgres
		dsnCode = `dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.GetEnv("DB_HOST", "localhost"),
		config.GetEnv("DB_PORT", "5432"),
		config.GetEnv("DB_USER", "postgres"),
		config.GetEnv("DB_PASSWORD", "postgres"),
		config.GetEnv("DB_NAME", "` + t.projectName + `"),
		config.GetEnv("DB_SSLMODE", "disable"),
	)`
	}

	return `// Package database provides ` + databaseDesc + ` database connectivity and management.
// It configures GORM for database operations, handles connection pooling,
// runs automatic migrations, and manages graceful shutdown through fx lifecycle hooks.
// This package is part of the infrastructure layer in the hexagonal architecture.
package database

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"` + databaseImport + `"
	"gorm.io/gorm"

	"` + t.projectName + `/internal/models"
	"` + t.projectName + `/pkg/config"
)

// Module provides the database dependency via fx with automatic lifecycle management.
var Module = fx.Module("database",
	fx.Provide(NewDatabase),
	fx.Invoke(registerHooks),
)

// NewDatabase creates a new GORM database connection configured from environment variables.
// It establishes a database connection, configures connection pooling, and runs
// automatic migrations for all domain models. Returns an error if connection fails.
func NewDatabase(logger zerolog.Logger) (*gorm.DB, error) {
	// Build DSN from environment variables
	` + dsnCode + `

	db, err := gorm.Open(` + connectionCode + `, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	logger.Info().Msg("Successfully connected to database")

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool parameters
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * 60) // 5 minutes

	// AutoMigrate database schemas
	if err := db.AutoMigrate(&models.User{}, &models.RefreshToken{}); err != nil {
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}

	logger.Info().Msg("Database migrations completed successfully")
	logger.Info().Msg("Database connection pool configured and ready")

	return db, nil
}

// registerHooks registers fx lifecycle hooks for graceful database shutdown.
// It ensures the database connection is properly closed when the application stops.
func registerHooks(lifecycle fx.Lifecycle, db *gorm.DB, logger zerolog.Logger) {
	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info().Msg("Closing database connection")
			sqlDB, err := db.DB()
			if err != nil {
				return err
			}
			return sqlDB.Close()
		},
	})
}

`
}

// ServerTemplate returns the internal/infrastructure/server/server.go file content
func (t *ProjectTemplates) ServerTemplate() string {
	return `// Package server provides HTTP server configuration and lifecycle management.
// It creates and configures a Fiber application with middleware, error handling,
// and graceful shutdown support through fx lifecycle hooks. This package is part
// of the infrastructure layer and coordinates all HTTP-related concerns.
package server

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"gorm.io/gorm"

	"` + t.projectName + `/pkg/config"
	httpRoutes "` + t.projectName + `/internal/adapters/http"
	"` + t.projectName + `/internal/adapters/middleware"

	// Swagger docs - generated by swag init
	_ "` + t.projectName + `/docs"
)

// Module provides the Fiber server dependency via fx with automatic lifecycle management.
var Module = fx.Module("server",
	fx.Provide(NewServer),
	fx.Invoke(registerHooks),
	fx.Invoke(httpRoutes.RegisterRoutes),
)

// NewServer creates and configures a new Fiber application with centralized error handling.
// It sets up the application name, error handler, and common routes like favicon handling.
// The server is ready to accept route registrations after creation.
func NewServer(logger zerolog.Logger, db *gorm.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "` + t.projectName + `",
		ErrorHandler: middleware.ErrorHandler,
		// Increase buffer sizes to prevent "Request Header Fields Too Large" errors
		ReadBufferSize:  16384, // 16KB (default is 4KB)
		WriteBufferSize: 16384,
	})

	// Ignore common browser requests (favicon, apple-touch-icon)
	// These would otherwise pollute error logs
	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})
	app.Get("/apple-touch-icon*.png", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	logger.Info().Msg("Fiber server initialized with centralized error handler")

	return app
}

// registerHooks registers fx lifecycle hooks for server startup and graceful shutdown.
// It starts the server in a background goroutine on startup and properly shuts it down
// when the application receives a termination signal.
func registerHooks(lifecycle fx.Lifecycle, app *fiber.App, logger zerolog.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			port := config.GetEnv("APP_PORT", "8080")
			logger.Info().Str("port", port).Msg("Starting Fiber server")

			// Start server in background goroutine
			go func() {
				if err := app.Listen(":" + port); err != nil {
					logger.Error().Err(err).Msg("Server stopped unexpectedly")
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info().Msg("Shutting down Fiber server gracefully")
			return app.ShutdownWithContext(ctx)
		},
	})
}

`
}

// HealthHandlerTemplate returns the internal/adapters/http/health.go file content
func (t *ProjectTemplates) HealthHandlerTemplate() string {
	return `// Package http provides HTTP route registration and health check endpoints.
// It coordinates route setup for the Fiber application and provides essential
// endpoints like health checks for container orchestration and load balancers.
package http

import (
	"github.com/gofiber/fiber/v2"
)

// HealthResponse represents the health check response structure.
// It provides a simple status field for health monitoring systems.
type HealthResponse struct {
	Status string ` + "`json:\"status\"`" + `
}

// RegisterHealthRoutes registers health check routes on the Fiber application.
// The health endpoint is used by container orchestrators and load balancers
// to verify the application is running and ready to accept requests.
func RegisterHealthRoutes(app *fiber.App) {
	app.Get("/health", healthHandler)
}

// healthHandler handles health check requests and returns the application status.
// It returns a simple JSON response indicating the service is operational.
func healthHandler(c *fiber.Ctx) error {
	return c.JSON(HealthResponse{
		Status: "ok",
	})
}
`
}

// ConfigTemplate returns the pkg/config/env.go file content
func (t *ProjectTemplates) ConfigTemplate() string {
	return `// Package config provides configuration management utilities for the application.
// It offers a simple interface for accessing environment variables with sensible
// defaults, enabling twelve-factor app configuration patterns.
package config

import "os"

// GetEnv retrieves an environment variable with a fallback default value.
// If the environment variable is not set or is empty, the default value is returned.
// This function is the primary way to access configuration throughout the application.
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
`
}

// UpdatedMainGoTemplate returns the updated cmd/main.go file content with fx integration
func (t *ProjectTemplates) UpdatedMainGoTemplate() string {
	return `package main

import (
	"log"

	"github.com/joho/godotenv"
	"go.uber.org/fx"

	"` + t.projectName + `/internal/adapters/handlers"
	"` + t.projectName + `/internal/adapters/repository"
	"` + t.projectName + `/internal/domain/user"
	"` + t.projectName + `/internal/infrastructure/database"
	"` + t.projectName + `/internal/infrastructure/server"
	"` + t.projectName + `/pkg/auth"
	"` + t.projectName + `/pkg/logger"
)

// @title ` + t.projectName + ` API
// @version 1.0
// @description A Go starter kit with authentication, user management, and CRUD operations
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load environment variables from .env file
	// This is primarily for local development; in production, use system environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found or couldn't be loaded")
	}

	fx.New(
		// Core infrastructure
		logger.Module,
		database.Module,

		// Authentication & authorization
		auth.Module,

		// Domain services
		user.Module,

		// Data persistence
		repository.Module,

		// HTTP handlers
		handlers.Module,

		// HTTP server (must be last as it depends on handlers)
		server.Module,
	).Run()
}
`
}

// GitHubActionsWorkflowTemplate returns the .github/workflows/ci.yml file content
func (t *ProjectTemplates) GitHubActionsWorkflowTemplate() string {
	return `name: CI

on:
  push:
    branches: [ "main" ]
  pull_request:
    branches: [ "main" ]

jobs:
  quality:
    name: Quality & Security
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
          cache: false # golangci-lint-action handles its own caching

      - name: Run Linter
        uses: golangci/golangci-lint-action@v6
        with:
          version: v1.60
          args: --timeout=5m

  test:
    name: Test & Build
    runs-on: ubuntu-latest
    needs: quality # Run tests only if lint passes
    services:
      postgres:
        image: postgres:16-alpine
        env:
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: ` + t.projectName + `
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'

      - name: Run Tests
        run: make test
        env:
          DB_HOST: localhost
          DB_PORT: 5432
          DB_USER: postgres
          DB_PASSWORD: postgres
          DB_NAME: ` + t.projectName + `
          DB_SSLMODE: disable

      - name: Build Check
        run: go build -v ./...
`
}

// DocsReadmeTemplate returns the docs/README.md file content (navigation hub)
func (t *ProjectTemplates) DocsReadmeTemplate() string {
	return `# Documentation ` + t.projectName + `

Documentation complète pour le projet ` + t.projectName + `.

## Table des matières

1. [Démarrage rapide](./quick-start.md)

## Aide rapide

- **Lancer le projet**: ` + "`make run`" + `
- **Tests**: ` + "`make test`" + `
- **API Health**: ` + "`http://localhost:8080/health`" + `

## Ressources

- [create-go-starter Documentation](https://github.com/tky0065/go-starter-kit)
- [Fiber Documentation](https://docs.gofiber.io/)
- [GORM Documentation](https://gorm.io/docs/)
`
}

// SwaggerDocsTemplate returns a placeholder docs/docs.go file
// This file will be overwritten when running 'make swagger' (swag init)
// but is needed for the project to compile before swagger generation
func (t *ProjectTemplates) SwaggerDocsTemplate() string {
	return `// Package docs provides Swagger documentation for the API.
// This is a placeholder file that allows the project to compile
// before running 'make swagger' to generate the actual documentation.
//
// Run 'make swagger' to generate the complete Swagger documentation.
// This will overwrite this file with the generated content.
package docs

import "github.com/swaggo/swag"

const docTemplate = ` + "`" + `{
    "swagger": "2.0",
    "info": {
        "description": "` + t.projectName + ` API - Run 'make swagger' to generate complete documentation",
        "title": "` + t.projectName + ` API",
        "contact": {},
        "version": "1.0"
    },
    "host": "localhost:8080",
    "basePath": "/api/v1",
    "paths": {}
}` + "`" + `

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "` + t.projectName + ` API",
	Description:      "` + t.projectName + ` API - Run 'make swagger' to generate complete documentation",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
`
}

// QuickStartTemplate returns the docs/quick-start.md file content
func (t *ProjectTemplates) QuickStartTemplate() string {
	return `# Démarrage rapide

Guide pour lancer ` + t.projectName + ` en 5 minutes.

## Prérequis

- Go 1.25+
- PostgreSQL (ou Docker)

## Installation

### 1. Installer les dépendances

` + "```bash" + `
go mod tidy
` + "```" + `

### 2. Configurer la base de données

**Option A: PostgreSQL local**

` + "```bash" + `
# macOS
brew install postgresql
brew services start postgresql
createdb ` + t.projectName + `

# Linux
sudo apt install postgresql
sudo systemctl start postgresql
sudo -u postgres createdb ` + t.projectName + `
` + "```" + `

**Option B: Docker (recommandé)**

` + "```bash" + `
docker run -d \
  --name postgres \
  -e POSTGRES_DB=` + t.projectName + ` \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine
` + "```" + `

### 3. Configurer l'environnement

Le fichier ` + "`.env`" + ` a déjà été créé. Générez un JWT secret:

` + "```bash" + `
# Générer un secret fort
openssl rand -base64 32

# Éditer .env
nano .env
` + "```" + `

Ajoutez dans ` + "`.env`" + `:
` + "```bash" + `
JWT_SECRET=<secret_généré_ci-dessus>
` + "```" + `

### 4. Lancer l'application

` + "```bash" + `
make run
` + "```" + `

L'API sera disponible sur ` + "`http://localhost:8080`" + `

### 5. Tester

` + "```bash" + `
# Health check
curl http://localhost:8080/health
# {"status":"ok"}
` + "```" + `

## Premier utilisateur

### Register

` + "```bash" + `
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password123"}'
` + "```" + `

Réponse (exemple):
` + "```json" + `
{
  "status": "success",
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "token_type": "Bearer",
    "expires_in": 900
  }
}
` + "```" + `

### Login

` + "```bash" + `
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password123"}'
` + "```" + `

### Utiliser l'access token

` + "```bash" + `
# Sauvegarder le token (remplacez par votre token)
TOKEN="eyJhbGc..."

# Lister les utilisateurs
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer $TOKEN"
` + "```" + `

## Endpoints disponibles

### Public (sans auth)

- ` + "`GET /health`" + ` - Health check
- ` + "`POST /api/v1/auth/register`" + ` - Créer un compte
- ` + "`POST /api/v1/auth/login`" + ` - Se connecter
- ` + "`POST /api/v1/auth/refresh`" + ` - Rafraîchir le token

### Protected (JWT required)

- ` + "`GET /api/v1/users`" + ` - Liste des utilisateurs
- ` + "`GET /api/v1/users/:id`" + ` - Détails d'un utilisateur
- ` + "`PUT /api/v1/users/:id`" + ` - Mettre à jour
- ` + "`DELETE /api/v1/users/:id`" + ` - Supprimer (soft delete)

## Développement

### Commandes utiles

` + "```bash" + `
# Lancer l'app
make run

# Tests
make test

# Tests avec coverage
make test-coverage

# Linting
make lint

# Build
make build

# Docker
make docker-build
make docker-run
` + "```" + `

### Structure du projet

` + "```" + `
` + t.projectName + `/
├── cmd/main.go                  # Point d'entrée (fx bootstrap)
├── internal/
│   ├── domain/                  # Logique métier
│   │   ├── user/                # Domaine User
│   │   └── errors.go            # Erreurs métier
│   ├── adapters/                # HTTP handlers, middleware, repository
│   ├── infrastructure/          # DB, server config
│   └── interfaces/              # Ports (interfaces)
├── pkg/                         # Packages réutilisables (auth, config, logger)
├── .env                         # Configuration
└── Makefile                     # Commandes
` + "```" + `

## Dépannage

### Erreur: "connection refused" sur DB

Vérifiez que PostgreSQL est démarré:

` + "```bash" + `
# Docker
docker ps | grep postgres

# Local
brew services list  # macOS
systemctl status postgresql  # Linux
` + "```" + `

### Erreur: "Invalid JWT secret"

Assurez-vous que ` + "`JWT_SECRET`" + ` est défini dans ` + "`.env`" + `:

` + "```bash" + `
cat .env | grep JWT_SECRET
` + "```" + `

Si vide, générez-en un:

` + "```bash" + `
echo "JWT_SECRET=$(openssl rand -base64 32)" >> .env
` + "```" + `

### Port 8080 déjà utilisé

Changez ` + "`APP_PORT`" + ` dans ` + "`.env`" + `:

` + "```bash" + `
APP_PORT=3000
` + "```" + `

## Prochaines étapes

- Lisez le README principal pour plus de détails
- Consultez le code dans ` + "`internal/domain/user/`" + ` pour comprendre la structure
- Ajoutez vos propres domaines en suivant le pattern User
- Déployez avec Docker: ` + "`make docker-build && make docker-run`" + `

Bon développement! 🚀
`
}

// readmeDatabaseDescription returns a database-specific description for the README
// This ensures generated projects have accurate documentation for their chosen database
func (t *ProjectTemplates) readmeDatabaseDescription() string {
	switch t.database {
	case "sqlite":
		return `**SQLite** - Base de données locale (fichier .db), parfait pour prototypage et tests`
	case "mysql":
		return `**MySQL/MariaDB** - Base de données relationnelle, compatible partout, facile à déployer`
	case "mongodb":
		return `**MongoDB** - Base de données NoSQL orientée documents, flexible et scalable`
	case "postgres":
		fallthrough
	default:
		return `**PostgreSQL** - Base de données relationnelle avancée, robuste et performante`
	}
}

// readmeDatabasePrerequisites returns database-specific prerequisites section for README
func (t *ProjectTemplates) readmeDatabasePrerequisites() string {
	switch t.database {
	case "sqlite":
		return `**Go 1.25+** - [Télécharger](https://golang.org/dl/)
- **SQLite** - Inclus dans Go (aucune installation externe requise)
- **Docker** (optionnel) - Pour containerisation
- **Make** - Pour les commandes de build
- **swag** (optionnel) - Pour régénérer la documentation Swagger`
	case "mysql":
		return `**Go 1.25+** - [Télécharger](https://golang.org/dl/)
- **MySQL 8.0+** ou **MariaDB 10.5+** - Base de données (peut être lancée via Docker)
- **Docker** (optionnel) - Pour containerisation
- **Make** - Pour les commandes de build
- **swag** (optionnel) - Pour régénérer la documentation Swagger`
	case "mongodb":
		return `**Go 1.25+** - [Télécharger](https://golang.org/dl/)
- **MongoDB 5.0+** - Base de données (peut être lancée via Docker)
- **Docker** (optionnel) - Pour containerisation
- **Make** - Pour les commandes de build
- **swag** (optionnel) - Pour régénérer la documentation Swagger`
	case "postgres":
		fallthrough
	default:
		return `**Go 1.25+** - [Télécharger](https://golang.org/dl/)
- **PostgreSQL 14+** - Base de données (peut être lancée via Docker)
- **Docker** (optionnel) - Pour containerisation
- **Make** - Pour les commandes de build
- **swag** (optionnel) - Pour régénérer la documentation Swagger`
	}
}

// readmeDatabaseConfig returns database-specific configuration example for README
func (t *ProjectTemplates) readmeDatabaseConfig() string {
	switch t.database {
	case "sqlite":
		return `# Application
APP_NAME=` + t.projectName + `
APP_ENV=development
APP_PORT=8080

# Database (SQLite - fichier local)
DB_NAME=` + t.projectName + `.db

# JWT
JWT_SECRET=                  # À REMPLIR!
JWT_EXPIRY=15m               # 15 minutes
REFRESH_TOKEN_EXPIRY=168h    # 7 jours
`
	case "mysql":
		return `# Application
APP_NAME=` + t.projectName + `
APP_ENV=development
APP_PORT=8080

# Database (MySQL)
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=secret
DB_NAME=` + t.projectName + `

# JWT
JWT_SECRET=                  # À REMPLIR!
JWT_EXPIRY=15m               # 15 minutes
REFRESH_TOKEN_EXPIRY=168h    # 7 jours
`
	case "mongodb":
		return `# Application
APP_NAME=` + t.projectName + `
APP_ENV=development
APP_PORT=8080

# Database (MongoDB)
DB_HOST=localhost
DB_PORT=27017
DB_NAME=` + t.projectName + `

# JWT
JWT_SECRET=                  # À REMPLIR!
JWT_EXPIRY=15m               # 15 minutes
REFRESH_TOKEN_EXPIRY=168h    # 7 jours
`
	case "postgres":
		fallthrough
	default:
		return `# Application
APP_NAME=` + t.projectName + `
APP_ENV=development
APP_PORT=8080

# Database (PostgreSQL)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=` + t.projectName + `
DB_SSLMODE=disable

# JWT
JWT_SECRET=                  # À REMPLIR!
JWT_EXPIRY=15m               # 15 minutes
REFRESH_TOKEN_EXPIRY=168h    # 7 jours
`
	}
}

// readmeDatabaseSetup returns database-specific setup instructions for README
// For SQLite, this returns an empty string (no setup needed)
// For other databases, returns the setup instructions
func (t *ProjectTemplates) readmeDatabaseSetup() string {
	switch t.database {
	case "sqlite":
		// SQLite needs no external setup
		return ``
	case "mysql":
		return `### 3. Lancer MySQL

**Option A: Docker (recommandé)**

` + "```bash" + `
docker run -d \
  --name mysql \
  -e MYSQL_ROOT_PASSWORD=secret \
  -e MYSQL_DATABASE=` + t.projectName + ` \
  -p 3306:3306 \
  mysql:8.0
` + "```" + `

**Option B: MySQL local**

` + "```bash" + `
# macOS
brew install mysql
brew services start mysql
mysql -u root -p"secret" -e "CREATE DATABASE ` + t.projectName + `;"

# Linux
sudo apt install mysql-server
sudo systemctl start mysql
mysql -u root -e "CREATE DATABASE ` + t.projectName + `;"
` + "```" + `

### 4. Lancer l'application

` + "```bash" + `
make run
` + "```" + ``
	case "mongodb":
		return `### 3. Lancer MongoDB

**Option A: Docker (recommandé)**

` + "```bash" + `
docker run -d \
  --name mongodb \
  -e MONGO_INITDB_DATABASE=` + t.projectName + ` \
  -p 27017:27017 \
  mongo:latest
` + "```" + `

**Option B: MongoDB local**

` + "```bash" + `
# macOS
brew install mongodb-community
brew services start mongodb-community

# Linux
sudo apt install -y mongodb
sudo systemctl start mongodb
` + "```" + `

### 4. Lancer l'application

` + "```bash" + `
make run
` + "```" + ``
	case "postgres":
		fallthrough
	default:
		return `### 3. Lancer PostgreSQL

**Option A: Docker (recommandé)**

` + "```bash" + `
docker run -d \
  --name postgres \
  -e POSTGRES_DB=` + t.projectName + ` \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine
` + "```" + `

**Option B: PostgreSQL local**

` + "```bash" + `
# macOS
brew install postgresql
brew services start postgresql
createdb ` + t.projectName + `

# Linux
sudo apt install postgresql
sudo systemctl start postgresql
sudo -u postgres createdb ` + t.projectName + `
` + "```" + `

### 4. Lancer l'application

` + "```bash" + `
make run
` + "```" + ``
	}
}

// readmeDatabaseStackDescription returns database-specific description for tech stack table
func (t *ProjectTemplates) readmeDatabaseStackDescription() string {
	switch t.database {
	case "sqlite":
		return `ORM avec SQLite (fichier local)`
	case "mysql":
		return `ORM avec MySQL/MariaDB`
	case "mongodb":
		return `ODM avec MongoDB`
	case "postgres":
		fallthrough
	default:
		return `ORM avec PostgreSQL`
	}
}

// readmeDatabaseDockerRun returns database-specific Docker run command
func (t *ProjectTemplates) readmeDatabaseDockerRun() string {
	if t.database == "sqlite" {
		// For SQLite, just need the JWT_SECRET and basic config
		return `# Build
make docker-build

# Run with SQLite (database file mounted as volume)
docker run -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  -e JWT_SECRET=<secret> \
  ` + t.projectName + `:latest
`
	}
	// For other databases, include DB_HOST
	return `# Build
make docker-build

# Run
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e JWT_SECRET=<secret> \
  ` + t.projectName + `:latest
`
}

// SetupScriptTemplate returns the setup.sh file content for automated project setup
func (t *ProjectTemplates) SetupScriptTemplate() string {
	return `#!/bin/bash

# setup.sh - Automated setup script for ` + t.projectName + `
# This script configures your development environment with all required dependencies

set -e  # Exit on error

# Color codes for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Helper functions
print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_step() {
    echo -e "\n${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}$1${NC}"
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Welcome message
echo -e "\n${GREEN}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║  Configuration automatique de ` + t.projectName + `${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════════╝${NC}\n"

# ============================================================================
# STEP 1: Check Prerequisites
# ============================================================================
print_step "Étape 1/6: Vérification des prérequis"

MISSING_DEPS=0

# Check Go
if command_exists go; then
    GO_VERSION=$(go version | awk '{print $3}')
    print_success "Go est installé: $GO_VERSION"
else
    print_error "Go n'est pas installé. Installez Go 1.25+ depuis https://golang.org/dl/"
    MISSING_DEPS=1
fi

# Check openssl
if command_exists openssl; then
    print_success "OpenSSL est installé"
else
    print_error "OpenSSL n'est pas installé. Installez avec: brew install openssl (macOS) ou apt install openssl (Linux)"
    MISSING_DEPS=1
fi

# Check Docker (optional but recommended)
if command_exists docker; then
    print_success "Docker est installé"
    DOCKER_AVAILABLE=1
else
    print_info "Docker n'est pas installé (optionnel). PostgreSQL devra être installé localement."
    DOCKER_AVAILABLE=0
fi

# Check psql (PostgreSQL client)
if command_exists psql; then
    print_success "Client PostgreSQL (psql) est installé"
    PSQL_AVAILABLE=1
else
    print_info "Client PostgreSQL (psql) n'est pas installé (optionnel)"
    PSQL_AVAILABLE=0
fi

if [ $MISSING_DEPS -eq 1 ]; then
    print_error "Des dépendances obligatoires sont manquantes. Installez-les et relancez ce script."
    exit 1
fi

# ============================================================================
# STEP 2: Install Go Dependencies
# ============================================================================
print_step "Étape 2/6: Installation des dépendances Go"

print_info "Exécution de 'go mod tidy'..."
if go mod tidy; then
    print_success "Dépendances Go installées avec succès"
else
    print_error "Échec de l'installation des dépendances Go"
    exit 1
fi

# Install swag CLI silently for Swagger documentation
if ! command_exists swag; then
    print_info "Installation de swag (générateur Swagger)..."
    go install github.com/swaggo/swag/cmd/swag@latest 2>/dev/null
    if command_exists swag; then
        print_success "swag installé avec succès"
    fi
fi

# ============================================================================
# STEP 3: Generate JWT Secret
# ============================================================================
print_step "Étape 3/6: Génération du JWT secret"

if [ -f .env ]; then
    JWT_CURRENT=$(grep "^JWT_SECRET=" .env | cut -d '=' -f2)
    if [ -n "$JWT_CURRENT" ] && [ "$JWT_CURRENT" != "" ]; then
        print_info "JWT_SECRET existe déjà dans .env"
        echo -n "Voulez-vous le régénérer? (y/N): "
        read -r REGEN_JWT
        if [[ ! $REGEN_JWT =~ ^[Yy]$ ]]; then
            print_info "JWT_SECRET conservé"
            SKIP_JWT=1
        else
            SKIP_JWT=0
        fi
    else
        SKIP_JWT=0
    fi
else
    print_error "Fichier .env introuvable. Création depuis .env.example..."
    if [ -f .env.example ]; then
        cp .env.example .env
        print_success "Fichier .env créé"
    else
        print_error ".env.example introuvable. Impossible de continuer."
        exit 1
    fi
    SKIP_JWT=0
fi

if [ $SKIP_JWT -eq 0 ]; then
    print_info "Génération d'un JWT secret sécurisé..."
    JWT_SECRET=$(openssl rand -base64 32)

    # Update .env file with JWT secret
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        sed -i '' "s|^JWT_SECRET=.*|JWT_SECRET=$JWT_SECRET|" .env
    else
        # Linux
        sed -i "s|^JWT_SECRET=.*|JWT_SECRET=$JWT_SECRET|" .env
    fi

    print_success "JWT_SECRET généré et ajouté à .env (chargé automatiquement au démarrage)"
fi

# ============================================================================
# STEP 4: Configure PostgreSQL
# ============================================================================
print_step "Étape 4/6: Configuration de PostgreSQL"

if [ $DOCKER_AVAILABLE -eq 1 ]; then
    echo -n "Voulez-vous démarrer PostgreSQL avec Docker? (Y/n): "
    read -r USE_DOCKER
    if [[ ! $USE_DOCKER =~ ^[Nn]$ ]]; then
        # Check if postgres container already exists
        if docker ps -a --format '{{.Names}}' | grep -q "^postgres$"; then
            print_info "Conteneur PostgreSQL 'postgres' existe déjà"

            # Check if it's running
            if docker ps --format '{{.Names}}' | grep -q "^postgres$"; then
                print_success "PostgreSQL est déjà en cours d'exécution"
            else
                print_info "Démarrage du conteneur existant..."
                docker start postgres
                sleep 2
                print_success "PostgreSQL démarré"
            fi
        else
            print_info "Création et démarrage d'un nouveau conteneur PostgreSQL..."
            docker run -d \
                --name postgres \
                -e POSTGRES_DB=` + t.projectName + ` \
                -e POSTGRES_PASSWORD=postgres \
                -p 5432:5432 \
                postgres:16-alpine

            # Wait for PostgreSQL to be ready
            print_info "Attente du démarrage de PostgreSQL (10 secondes)..."
            sleep 10
            print_success "PostgreSQL démarré avec Docker"
        fi

        POSTGRES_STARTED=1
    else
        print_info "Configuration Docker PostgreSQL ignorée"
        POSTGRES_STARTED=0
    fi
else
    print_info "Docker non disponible. Vérification de PostgreSQL local..."
    POSTGRES_STARTED=0
fi

# Try to connect to PostgreSQL to verify it's running
print_info "Vérification de la connexion PostgreSQL..."
if [ $PSQL_AVAILABLE -eq 1 ]; then
    if PGPASSWORD=postgres psql -h localhost -U postgres -d ` + t.projectName + ` -c '\q' 2>/dev/null; then
        print_success "Connexion PostgreSQL réussie"
        POSTGRES_STARTED=1
    else
        if [ $POSTGRES_STARTED -eq 0 ]; then
            print_error "Impossible de se connecter à PostgreSQL"
            print_info "Assurez-vous que PostgreSQL est installé et démarré:"
            print_info "  macOS: brew install postgresql && brew services start postgresql"
            print_info "  Linux: sudo apt install postgresql && sudo systemctl start postgresql"
            print_info "\nPuis créez la base de données:"
            print_info "  createdb ` + t.projectName + `"
            exit 1
        fi
    fi
else
    print_info "Client psql non disponible, impossible de vérifier la connexion"
    if [ $POSTGRES_STARTED -eq 0 ]; then
        print_info "Assurez-vous que PostgreSQL est installé et démarré manuellement"
    fi
fi

# ============================================================================
# STEP 5: Generate Swagger & Run Tests
# ============================================================================
print_step "Étape 5/6: Génération Swagger & Tests"

# Generate Swagger documentation
if command_exists swag; then
    print_info "Génération de la documentation Swagger..."
    if swag init -g cmd/main.go --output docs 2>/dev/null; then
        print_success "Documentation Swagger générée"
    else
        print_info "Génération Swagger ignorée (exécutez 'make swagger' manuellement)"
    fi
else
    print_info "swag non installé, génération Swagger ignorée"
fi

print_info "Lancement des tests unitaires..."
if go test ./... 2>/dev/null; then
    print_success "Tous les tests passent"
else
    print_info "Certains tests ont échoué (normal si la base n'est pas encore configurée)"
fi

# ============================================================================
# STEP 6: Verify Installation
# ============================================================================
print_step "Étape 6/6: Vérification de l'installation"

print_info "Vérification de la configuration..."

# Check .env file
if [ -f .env ]; then
    if grep -q "^JWT_SECRET=..*" .env; then
        print_success ".env configuré avec JWT_SECRET"
    else
        print_error ".env manque JWT_SECRET"
    fi
else
    print_error "Fichier .env manquant"
fi

# Check go.mod
if [ -f go.mod ]; then
    print_success "go.mod présent"
else
    print_error "go.mod manquant"
fi

# ============================================================================
# Summary and Next Steps
# ============================================================================
echo -e "\n${GREEN}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║  ✅ Configuration terminée avec succès!${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════════╝${NC}\n"

print_info "Prochaines étapes:"
echo "  1. Lancer l'application:    make run"
echo "  2. Vérifier la santé:       curl http://localhost:8080/health"
echo "  3. Documentation Swagger:   http://localhost:8080/swagger/index.html"
echo ""
print_info "Documentation:"
echo "  - Guide rapide: docs/quick-start.md"
echo "  - README:       README.md"
echo ""
print_success "Bon développement! 🚀"
`
}
