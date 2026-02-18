# create-go-starter

A powerful CLI tool to generate production-ready Go projects in seconds.

## <i class="material-icons info">info</i> About

**create-go-starter** is a stable, production-ready CLI generator for creating modern Go applications.

- <i class="material-icons success">check</i> **Stable version** - v1.4.0 available
- <i class="material-icons success">check</i> **3 templates** - minimal, full, GraphQL
- <i class="material-icons success">check</i> **Tested and validated** - Used in production
- <i class="material-icons success">check</i> **Open source** - MIT License
- <i class="material-icons success">check</i> **Actively maintained** - Regular updates

## Overview

`create-go-starter` is a Go project generator that creates a complete hexagonal architecture with all the essential features of a modern backend application. With a single command, get a structured project with JWT authentication, REST API, database, tests, and Docker configuration ready for deployment.

### Included Features

- **Hexagonal architecture** (Ports & Adapters) - Clear separation of concerns
- **JWT Authentication** - Access tokens + Refresh tokens with secure rotation
- **REST API** with Fiber v2 - High-performance web framework
- **Database** - GORM with PostgreSQL and automatic migrations
- **Dependency injection** - uber-go/fx for a modular architecture
- **Comprehensive tests** - Unit and integration tests
- **Swagger documentation** - API automatically documented with OpenAPI
- **Docker** - Optimized multi-stage build and docker-compose
- **CI/CD** - Pre-configured GitHub Actions pipeline
- **Structured logging** - rs/zerolog for professional logs
- **Validation** - go-playground/validator for input validation
- **Makefile** - Useful commands for dev, test, build, and deployment

### <i class="material-icons success">new_releases</i> New in v1.2.0: CRUD Generator

**Add new models in 2 seconds!** The `add-model` command automatically generates all the necessary CRUD code:

```bash
cd mon-projet
create-go-starter add-model Todo --fields "title:string,completed:bool"
```

**What is automatically generated:**
- <i class="material-icons success small">check</i> Model with GORM tags
- <i class="material-icons success small">check</i> Repository (interface + implementation)
- <i class="material-icons success small">check</i> Service (business logic)
- <i class="material-icons success small">check</i> HTTP Handlers (REST endpoints)
- <i class="material-icons success small">check</i> Tests (unit + integration)
- <i class="material-icons success small">check</i> Routes, migrations, fx modules

**Supported relations:**
- `--belongs-to Parent` - N:1 relation (child → parent)
- `--has-many Child` - 1:N relation (parent → children)

**Example with relations:**
```bash
create-go-starter add-model Category --fields "name:string:unique"
create-go-starter add-model Post --fields "title:string,content:string" --belongs-to Category
create-go-starter add-model Comment --fields "author:string,content:string" --belongs-to Post
```

**Result:** Complete blog (Category → Post → Comment) with nested endpoints and preloading.

<i class="material-icons">arrow_forward</i> [Complete add-model guide](./usage.md#ajouter-des-modeles-add-model)

### <i class="material-icons success">new_releases</i> New in v1.3.0: Advanced Observability

**Production-ready monitoring in one command!** The `--observability=advanced` flag generates a complete observability stack:

```bash
create-go-starter mon-app --template=full --observability=advanced
```

**What is automatically generated:**
- <i class="material-icons success small">check</i> **Prometheus Metrics** — `/metrics` endpoint with HTTP metrics (latency, throughput, errors)
- <i class="material-icons success small">check</i> **Distributed Tracing** — OpenTelemetry + Jaeger with W3C traceparent propagation
- <i class="material-icons success small">check</i> **K8s Health Checks** — `/health/liveness` and `/health/readiness` with DB verification
- <i class="material-icons success small">check</i> **Grafana Dashboard** — Pre-configured 7-panel dashboard with alerting
- <i class="material-icons success small">check</i> **Docker Compose** — Complete stack (Jaeger + Prometheus + Grafana)
- <i class="material-icons success small">check</i> **Kubernetes Probes** — Automatically generated `probes.yaml` file

<i class="material-icons">arrow_forward</i> [Complete observability guide](./usage.md#observabilité---observability)

### <i class="material-icons success">new_releases</i> New in v1.4.0: CLI Improvements

**Redesigned developer experience!** Guided interactive mode, dry-run preview, environment diagnostics, and short aliases for all flags.

```bash
# Guided interactive mode
create-go-starter -i

# Preview without creating
create-go-starter mon-app --dry-run

# Environment diagnostics
create-go-starter doctor

# Short aliases for all flags
create-go-starter -t minimal -d sqlite -n mon-app
```

**What's new in v1.4.0:**
- <i class="material-icons success small">check</i> **Interactive Mode** (`--interactive` / `-i`) — Step-by-step guided assistant
- <i class="material-icons success small">check</i> **Dry-Run** (`--dry-run` / `-n`) — File preview without writing
- <i class="material-icons success small">check</i> **Doctor** (`create-go-starter doctor`) — Go, Git, Docker diagnostics
- <i class="material-icons success small">check</i> **Progress bar** — Visual feedback during generation
- <i class="material-icons success small">check</i> **Short aliases** — `-t`, `-d`, `-o`, `-i`, `-n`, `-h`

<i class="material-icons">arrow_forward</i> [Complete CLI improvements guide](./usage.md#mode-interactif---interactive)

## Quick Installation

### Method 1: Direct Installation (Recommended)

Global installation with a single command, without cloning the repository:

```bash
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@latest
```

The binary will be installed in `$GOPATH/bin` (usually `~/go/bin`). Make sure this directory is in your PATH.

**Note**: If `create-go-starter` is not recognized after installation, add `$GOPATH/bin` to your PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Method 2: Build from Source

For contributors or customization:

```bash
git clone https://github.com/tky0065/go-starter-kit.git
cd go-starter-kit
go build -o create-go-starter ./cmd/create-go-starter
# The binary is now available: ./create-go-starter
```

### Method 3: Build with Makefile

```bash
git clone https://github.com/tky0065/go-starter-kit.git
cd go-starter-kit
make build
# The binary is available: ./create-go-starter
```

For more details, see the [complete installation guide](./installation.md).

## Basic Usage

### Create a New Project

```bash
create-go-starter mon-super-projet
```

This command will:
1. Create the complete project structure
2. Generate ~30+ files (handlers, services, repositories, tests, etc.)
3. Configure all necessary files (.env, Dockerfile, Makefile, etc.)
4. Copy the `.env.example` file to `.env`
5. Initialize a Git repository with an initial commit (if Git is available)

### Choose a Template

By default, `create-go-starter` generates a **full** project with JWT authentication and user management. You can choose a different template with the `--template` flag:

```bash
create-go-starter mon-projet --template minimal    # API REST basique avec Swagger
create-go-starter mon-projet --template full       # API complète avec JWT auth (défaut)
create-go-starter mon-projet --template graphql    # API GraphQL avec gqlgen
```

**Available templates**:

| Template | Description | Use Case |
|----------|-------------|----------|
| `minimal` | Basic REST API with Swagger (no authentication) | Quick prototypes, simple public APIs |
| `full` | Complete API with JWT auth, user management, and Swagger | Complete backend applications (default) |
| `graphql` | GraphQL API with gqlgen and GraphQL Playground | Applications requiring GraphQL |

For more details on the differences between templates, see the [usage guide](./usage.md#templates-disponibles).

### Run the Generated Project

#### Option 1: Automatic Setup (Recommended) <i class="material-icons">rocket_launch</i>

```bash
cd mon-super-projet
./setup.sh
make run
```

The `setup.sh` script automates:
- Go dependency installation
- JWT secret generation
- PostgreSQL configuration (Docker or local)
- Installation verification

#### Option 2: Manual Setup

```bash
cd mon-super-projet

# Install dependencies and generate go.sum
go mod tidy

# Configure the JWT secret in .env
# JWT_SECRET=<generate with: openssl rand -base64 32>

# Start the database (PostgreSQL)
docker run -d --name postgres \
  -e POSTGRES_DB=mon-super-projet \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine

# Start the application
make run
```

The API will be available at `http://localhost:8080`

```bash
# Test the health check
curl http://localhost:8080/health
# {"status":"ok"}
```

## Generated Structure

Here is what `create-go-starter` generates for you:

```
mon-super-projet/
├── cmd/
│   └── main.go                    # Entry point with fx dependency injection
├── internal/
│   ├── models/                    # Shared domain entities
│   │   └── user.go                # User, RefreshToken, AuthResponse
│   ├── domain/                    # Domain layer (business logic)
│   │   ├── user/                  # User domain
│   │   │   ├── service.go         # Business logic (Register, Login, etc.)
│   │   │   └── module.go          # fx module
│   │   └── errors.go              # Custom business errors
│   ├── adapters/                  # Adapters (HTTP, DB)
│   │   ├── handlers/              # HTTP handlers
│   │   │   ├── auth_handler.go    # Auth endpoints (register, login, refresh)
│   │   │   └── user_handler.go    # User CRUD endpoints
│   │   ├── middleware/            # Fiber middleware
│   │   │   ├── auth_middleware.go # JWT verification
│   │   │   └── error_handler.go   # Centralized error handling
│   │   ├── repository/            # Repository implementation
│   │   │   └── user_repository.go # GORM implementation
│   │   └── http/                  # HTTP utilities
│   │       ├── health.go          # Health check handler
│   │       └── routes.go          # Centralized routes
│   ├── infrastructure/            # Infrastructure
│   │   ├── database/              # DB configuration (GORM, migrations)
│   │   └── server/                # Fiber app configuration
│   └── interfaces/                # Ports (interfaces)
│       └── user_repository.go     # UserRepository interface
├── pkg/                           # Reusable packages
│   ├── auth/                      # JWT utilities
│   ├── config/                    # Configuration loading (.env)
│   └── logger/                    # zerolog configuration
├── .github/workflows/
│   └── ci.yml                     # CI/CD pipeline (lint, test, build)
├── .env                           # Configuration (copied from .env.example)
├── .env.example                   # Configuration template
├── .gitignore                     # Git exclusions
├── .golangci.yml                  # Linter configuration
├── Dockerfile                     # Multi-stage build for production
├── Makefile                       # Useful commands (run, test, lint, docker, etc.)
├── setup.sh                       # Automatic setup script
├── go.mod                         # Go module with dependencies
└── README.md                      # Project documentation
```

For a detailed explanation of each component, see the [usage guide](./usage.md).

## Tech Stack

Generated projects use the best libraries from the Go ecosystem:

| Component | Library | Version | Description |
|-----------|---------|---------|-------------|
| Web Framework | [Fiber](https://gofiber.io/) | v2 | Fast HTTP framework, inspired by Express |
| ORM | [GORM](https://gorm.io/) | v1 | Go ORM with PostgreSQL support |
| Dependency Injection | [fx](https://uber-go.github.io/fx/) | latest | DI framework by Uber |
| Logging | [zerolog](https://github.com/rs/zerolog) | latest | High-performance structured logger |
| JWT | [golang-jwt](https://github.com/golang-jwt/jwt) | v5 | JWT tokens for authentication |
| Validation | [validator](https://github.com/go-playground/validator) | v10 | Struct validation |
| Swagger | [swaggo](https://github.com/swaggo/swag) | latest | OpenAPI API documentation |
| Crypto | golang.org/x/crypto | latest | bcrypt hashing for passwords |

## Documentation

### Essential Guides

- **[Installation Guide](./installation.md)** - Detailed installation with all methods
- **[Usage Guide](./usage.md)** - CLI usage and complete generated structure
- **[Generated Projects Guide](./generated-project-guide.md)** - Complete guide for developing with created projects (architecture, API, tests, deployment)

### Advanced Documentation

- **[CLI Architecture](./cli-architecture.md)** - Technical documentation for contributors
- **[Contributing Guide](./contributing.md)** - How to contribute to the project

## Quick Start in 30 Seconds

```bash
# 1. Install the tool
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@latest

# 2. Create a project
create-go-starter mon-projet

# 3. Automatic setup
cd mon-projet
./setup.sh

# 4. Run
make run

# 5. Test
curl http://localhost:8080/health
```

Or manual setup:
```bash
cd mon-projet
echo "JWT_SECRET=$(openssl rand -base64 32)" >> .env
docker run -d --name postgres -e POSTGRES_DB=mon-projet -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres:16-alpine
make run
```

## API Usage Examples

Once your project is running, you can test the API:

```bash
# Create a user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securePass123"}'

# Log in
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securePass123"}'

# Use the returned token to access protected endpoints
TOKEN="<access_token_from_login_response>"
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer $TOKEN"
```

For more examples and the complete API documentation, see the [generated projects guide](./generated-project-guide.md#api-reference).

## Available Makefile Commands

Generated projects include a Makefile with useful commands:

```bash
make help           # Display all available commands
make run            # Run the application
make build          # Build the binary
make test           # Run the tests
make test-coverage  # Tests with coverage report
make lint           # Lint the code (golangci-lint)
make docker-build   # Build the Docker image
make docker-run     # Run the Docker container
make clean          # Clean artifacts
```

## Prerequisites

- **Go 1.25 or higher** - [Download Go](https://golang.org/dl/)
- **PostgreSQL** - For generated projects (can be run via Docker)
- **Git** - For cloning and contributing
- **Docker** (optional) - For running PostgreSQL and containerizing the application
- **golangci-lint** (optional) - For linting

## Why create-go-starter?

### Time Savings

Instead of spending hours configuring:
- The project architecture
- JWT authentication
- The database connection
- Tests
- Docker and CI/CD
- Swagger documentation

Get all of this with **a single command** and start developing your business features immediately.

### Built-in Best Practices

- **Hexagonal architecture** - Clear separation between domain, adapters, and infrastructure
- **Dependency injection** - Testable and modular code
- **Centralized error handling** - Consistent error management
- **Security-first** - JWT, bcrypt, validation, CORS
- **Tests** - Unit and integration test examples
- **Clean code** - Adherence to Go conventions and strict linting

### Production-Ready

Generated projects are ready for production:
- Optimized multi-stage Docker build
- CI/CD with automated tests
- Structured logging for monitoring
- Environment-based configuration
- Health checks
- Graceful shutdown

## Contributing

Contributions are welcome! See the [contributing guide](./docs/contributing.md) to get started.

### Contribution Process

1. Fork the project
2. Create a branch (`git checkout -b feature/ma-fonctionnalite`)
3. Commit the changes (`git commit -m 'feat: ajouter une fonctionnalité'`)
4. Push to the branch (`git push origin feature/ma-fonctionnalite`)
5. Open a Pull Request

## Roadmap

**Completed features**:

- [x] **Multiple templates** - Three templates available (minimal, full, graphql)
- [x] **Multi-Database Support** - PostgreSQL, MySQL, SQLite (v1.1.0)
- [x] **CRUD Scaffolding Generator** - add-model command to generate complete models (v1.2.0)
- [x] **Advanced Observability** - Prometheus, Jaeger, Grafana, K8s Health Checks (v1.3.0)
- [x] **CLI Improvements** - Interactive mode, dry-run, doctor, short aliases (v1.4.0)

**Planned features**:

- [ ] Web framework choice (Gin, Echo, Chi)
- [ ] Microservices generation
- [ ] E2E test templates

## License

[MIT License](LICENSE) - Free to use for personal and commercial projects.

## Support

- **Issues**: [GitHub Issues](https://github.com/tky0065/go-starter-kit/issues)
- **Discussions**: [GitHub Discussions](https://github.com/tky0065/go-starter-kit/discussions)
- **Documentation**: [docs/](./docs/)

## Acknowledgments

Built with the excellent libraries from the Go community. Thanks to the maintainers of Fiber, GORM, fx, zerolog, and all other dependencies.

---

**Made with <i class="material-icons small" style="color:#e25555">favorite</i> for the Go community**

Start building your next backend application in seconds, not days!
