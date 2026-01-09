# Gemini Project Context: go-starter-kit

## Project Overview

`go-starter-kit` is a CLI productivity tool designed to scaffold production-ready Go API projects in minutes. It aims to eliminate the friction of starting new projects by providing an opinionated yet pragmatic foundation that balances development speed with industrial standards.

### Core Philosophy
- **Zero-to-Hero Experience:** A CLI tool (`create-go-starter`) that handles scaffolding, dependency management, and initial configuration.
- **Best Practices by Default:** Includes JWT authentication, validation, Swagger documentation, Dockerization, and CI/CD out of the box.
- **Lite Hexagonal Architecture:** A simplified hexagonal structure that ensures maintainability without excessive complexity.

### Tech Stack
- **Web Framework:** [Fiber](https://gofiber.io/) (Fast, Express-inspired)
- **Dependency Injection:** [uber-go/fx](https://github.com/uber-go/fx)
- **ORM:** [GORM](https://gorm.io/) with PostgreSQL support.
- **Logging:** [zerolog](https://github.com/rs/zerolog)
- **Documentation:** Swagger (auto-generated)
- **Containerization:** Docker (multi-stage builds)

---

## Directory Structure

### Main Tool (`/cmd/create-go-starter`)
- `main.go`: CLI entry point, handles project name validation and directory creation.
- `generator.go`: Orchestrates the file generation process using templates.
- `templates.go`: Contains all boilerplate code templates for the generated project.

### Generated Project Structure (Boilerplate)
- `cmd/main.go`: Application entry point using `fx` for dependency injection.
- `internal/`:
    - `domain/`: Core business logic and entities.
    - `adapters/`: External interfaces (HTTP handlers, Middleware).
    - `infrastructure/`: Infrastructure concerns (Database connection, Server configuration).
    - `interfaces/`: Port definitions (Ports in Hexagonal Architecture).
- `pkg/`: Public/Shared libraries (Config loader, Logger).
- `deployments/`: Deployment configurations (Docker, K8s).

---

## Building and Running

### The Scaffolding Tool (Main Repository)

- **Build the CLI:**
  ```bash
  make build
  ```
- **Run the CLI (Generates a new project):**
  ```bash
  make run PROJECT_NAME=my-new-api
  ```
- **Run Tests:**
  ```bash
  make test
  ```
- **Install CLI Globally:**
  ```bash
  make install
  ```

### The Generated Project

1.  **Initialize:**
    ```bash
    cd <project-name>
    go mod download
    cp .env.example .env
    ```
2.  **Run in Development:**
    ```bash
    make run
    # OR with hot-reload (requires air)
    make dev
    ```
3.  **Run Tests:**
    ```bash
    make test
    ```

---

## Development Conventions

- **Hexagonal Architecture:** Maintain a strict separation between the domain logic and external adapters. Domain logic should not depend on infrastructure or adapters.
- **Dependency Injection:** Always use `fx.Module` and `fx.Provide` for defining components. Avoid global variables.
- **Error Handling:** Use centralized error handling (middleware) and define domain-specific errors in `internal/domain/errors.go`.
- **Environment Variables:** All configuration must be driven by environment variables using the `pkg/config` utility.
- **Linting:** Ensure code passes `golangci-lint` (config provided in templates).
- **Graceful Shutdown:** The server and database connections must handle OS signals for clean termination.

## Key Files for Reference
- `cmd/create-go-starter/main.go`: Logic for the CLI tool.
- `cmd/create-go-starter/templates.go`: Source of truth for the generated boilerplate.
- `go.mod`: Project dependencies.
- `_bmad-output/planning-artifacts/prd.md`: High-level product requirements and vision.
