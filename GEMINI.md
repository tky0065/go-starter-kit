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
    - `models/`: Shared domain entities (User, RefreshToken, AuthResponse). Prevents circular dependencies.
    - `domain/`: Core business logic (services, not entities).
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

**Option 1: Automated Setup (Recommended)**
```bash
cd <project-name>
./setup.sh
make run
```

The `setup.sh` script automatically handles:
- Go dependency installation (`go mod tidy`)
- JWT secret generation and configuration
- PostgreSQL setup (Docker or local)
- Installation verification

**Option 2: Manual Setup**
1.  **Initialize:**
    ```bash
    cd <project-name>
    go mod tidy
    ```
2.  **Configure Environment:**
    ```bash
    # Generate JWT secret
    openssl rand -base64 32
    # Edit .env and add: JWT_SECRET=<generated-secret>
    ```
3.  **Start PostgreSQL:**
    ```bash
    docker run -d --name postgres \
      -e POSTGRES_DB=<project-name> \
      -e POSTGRES_PASSWORD=postgres \
      -p 5432:5432 \
      postgres:16-alpine
    ```
4.  **Run in Development:**
    ```bash
    make run
    # OR with hot-reload (requires air)
    make dev
    ```
5.  **Run Tests:**
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

---

## Documentation Maintenance Policy

**⚠️ MANDATORY REQUIREMENT**: Documentation MUST be updated whenever code changes are made.

### Update Documentation When:

1. **Adding Features**:
   - New templates → Update `docs/cli-architecture.md`
   - New generated files → Update `docs/usage.md` and `docs/generated-project-guide.md`
   - New entities/models → Update all architecture diagrams and code examples

2. **Modifying Architecture**:
   - Package restructuring → Update ALL documentation files
   - Dependency graph changes → Update architecture diagrams
   - Design pattern changes → Update examples and explanations

3. **Template Changes**:
   - Any template modification → Update corresponding documentation
   - Import changes → Update all code examples throughout docs

4. **Bug Fixes**:
   - If fix changes generated project behavior → Documentation update required

### Critical Documentation Files

Always review these after code changes:

| File | Purpose | Update When |
|------|---------|-------------|
| `README.md` | Project overview | Structure or features change |
| `docs/usage.md` | Usage guide | Generated structure changes |
| `docs/generated-project-guide.md` | Complete guide | Any template or architecture change |
| `docs/cli-architecture.md` | CLI internals | Generator or template logic changes |
| `CLAUDE.md` | AI context (this repo) | Project structure or conventions change |
| `GEMINI.md` | AI context (this file) | Project structure or conventions change |

### Standard Workflow

```bash
# Step 1: Code changes
vim cmd/create-go-starter/templates_user.go

# Step 2: Test thoroughly
go build -o create-go-starter ./cmd/create-go-starter
./create-go-starter test-validation-project
cd test-validation-project && go mod tidy && go build ./...

# Step 3: Update ALL affected docs
vim docs/cli-architecture.md
vim docs/generated-project-guide.md
vim README.md

# Step 4: Verify documentation accuracy
grep -r "old_pattern" docs/  # Find outdated references
# Update all found instances

# Step 5: Commit code + docs together
git add cmd/ docs/ README.md CLAUDE.md GEMINI.md
git commit -m "feat: implement feature X

- Add feature implementation in templates
- Update all documentation to reflect changes
- Add comprehensive examples in docs"
```

### Why This Is Critical

- **User Trust**: Outdated docs erode confidence in the project
- **Development Velocity**: Accurate docs enable faster onboarding and contributions
- **AI Assistance**: LLMs rely on current documentation for context
- **Maintenance**: Prevents technical debt accumulation

### Enforcement

- Pull requests with code changes but no doc updates will be questioned
- Documentation updates are NOT optional—they are part of the feature
- When in doubt, over-document rather than under-document

**Remember**: Undocumented code changes are considered incomplete work.
