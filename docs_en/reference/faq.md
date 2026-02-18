# Frequently Asked Questions (FAQ)

Frequently asked questions about `create-go-starter` and generated projects.

---

## <i class="material-icons">download</i> Installation & Configuration

### How do I install `create-go-starter`?

The recommended method is installation via `go install`:

```bash
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@latest
```

This command automatically downloads and installs the binary into `$GOPATH/bin`.

**Verify the installation:**

```bash
create-go-starter --help
```

<i class="material-icons info">arrow_forward</i> **Details:** See the [complete installation guide](../installation.md)

---

### <i class="material-icons warning">warning</i> Error "command not found: create-go-starter"

**Cause:** The `$GOPATH/bin` directory is not in your `PATH`, or the shell cache has not been reloaded.

**Solutions:**

**1. Reload the shell cache** (try this first):

```bash
hash -r
which create-go-starter
```

**2. Restart the terminal** (often the simplest solution)

Close and reopen your terminal.

**3. Add `$GOPATH/bin` to `PATH`:**

```bash
# Check GOPATH
go env GOPATH

# Add to PATH (in ~/.zshrc or ~/.bashrc)
export PATH=$PATH:$(go env GOPATH)/bin

# Reload
source ~/.zshrc  # or ~/.bashrc
```

**4. Use the full path temporarily:**

```bash
$(go env GOPATH)/bin/create-go-starter mon-projet
```

<i class="material-icons info">arrow_forward</i> **Details:** [Installation Guide - Troubleshooting](../installation.md#résolution-de-problèmes)

---

### What version of Go is required?

**Go 1.25 or higher** is required.

**Check your version:**

```bash
go version
# Should display: go version go1.25.x ...
```

**Update Go:** Download from [golang.org/dl](https://golang.org/dl/)

---

### How do I update `create-go-starter`?

Simply re-run the installation command:

```bash
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@latest
```

Go will automatically download the latest version.

---

## <i class="material-icons">storage</i> Database & Configuration

### Why is PostgreSQL the default database?

**PostgreSQL** is chosen as the default for several reasons:

- <i class="material-icons success">check</i> **Production reliability** - ACID compliant, strong data integrity
- <i class="material-icons success">check</i> **Advanced features** - JSON, arrays, full-text search
- <i class="material-icons success">check</i> **Performance** - Excellent for complex queries
- <i class="material-icons success">check</i> **Community** - Rich ecosystem, comprehensive documentation

PostgreSQL is the standard choice for modern backend applications.

<i class="material-icons info">arrow_forward</i> **Comparison:** [Database Selection Guide](../databases.md)

---

### How do I change the database?

**When creating a new project:**

```bash
# MySQL
create-go-starter mon-app --database=mysql

# SQLite
create-go-starter mon-app --database=sqlite

# PostgreSQL (default)
create-go-starter mon-app --database=postgres
```

**For an existing project:**

You need to regenerate the project with the new `--database` flag and migrate the data.

<i class="material-icons warning">warning</i> **Important:** Back up your data before migrating.

<i class="material-icons info">arrow_forward</i> **Complete guide:** [Database Migration Guide](../database-migration.md)

---

### Which database should I choose for my project?

**PostgreSQL** (default):
- Production applications
- Complex relational data
- Advanced queries (JSON, analytics)
- Reliability and scalability

**MySQL**:
- Shared hosting
- Broad compatibility
- Read-intensive workloads
- Teams familiar with MySQL

**SQLite**:
- Rapid prototyping / MVP
- Small applications (<100 concurrent users)
- Development and testing
- Desktop/embedded applications

<i class="material-icons info">arrow_forward</i> **Detailed comparison:** [databases.md - Decision Matrix](../databases.md#matrice-de-décision)

---

### <i class="material-icons error">error</i> Error "JWT_SECRET not set" - what should I do?

**Cause:** The JWT secret is not configured in the `.env` file.

**Solution:**

**1. Generate a JWT secret:**

```bash
openssl rand -base64 32
# Example output: xK7vZ9mN2pQ8rT4wL6jH3sB5cA1dF0eG=
```

**2. Add it to the `.env` file:**

```bash
# Edit .env
JWT_SECRET=xK7vZ9mN2pQ8rT4wL6jH3sB5cA1dF0eG=
```

**Or use the automatic script:**

```bash
./setup.sh
# The script automatically generates JWT_SECRET
```

!!! warning "Security"
    - Never commit the `.env` file to Git
    - Use a different secret for each environment (dev, staging, prod)
    - The secret must be at least 32 characters

<i class="material-icons info">arrow_forward</i> **Documentation:** [Generated Project Guide - Configuration](../generated-project-guide.md#configuration)

---

### Do I need Docker?

It depends on the chosen database:

| Database | Docker required? | Alternative |
|----------|------------------|-------------|
| **PostgreSQL** | <i class="material-icons success">check</i> Yes (recommended) | Local installation possible |
| **MySQL** | <i class="material-icons success">check</i> Yes (recommended) | Local installation possible |
| **SQLite** | <i class="material-icons error">close</i> No | Embedded database (file) |

**Launch PostgreSQL with Docker:**

```bash
docker run -d --name postgres \
  -e POSTGRES_DB=mon-app \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine
```

**Or use docker-compose** (included in generated projects):

```bash
docker-compose up -d
```

---

### How do I resolve database connection errors?

**Checks to perform:**

**1. Verify that the database is running:**

```bash
# For Docker
docker ps
# Should show the postgres container as "Up"

# For local PostgreSQL
pg_isready -h localhost -p 5432
```

**2. Verify the `.env` configuration:**

```bash
# PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=mon-app
DB_SSLMODE=disable

# MySQL
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=mon-app

# SQLite
DB_NAME=mon-app.db
```

**3. Test the connection manually:**

```bash
# PostgreSQL
psql -h localhost -U postgres -d mon-app

# MySQL
mysql -h localhost -u root -p mon-app
```

**4. Check the application logs:**

```bash
make run
# The logs often indicate the exact cause
```

<i class="material-icons warning">warning</i> **Common mistake:** Forgetting to start Docker or PostgreSQL before launching the app.

---

### How do I resolve database migration errors?

**Symptoms:**
- Tables not created
- Missing columns
- "relation does not exist" errors

**Solutions:**

**1. Force table recreation (development only):**

```go
// In infrastructure/database/database.go
db.AutoMigrate(&models.User{}, &models.RefreshToken{})
// Temporarily add:
// db.Migrator().DropTable(&models.User{}, &models.RefreshToken{})
// db.AutoMigrate(&models.User{}, &models.RefreshToken{})
```

**2. Check the GORM logs:**

The logs indicate the SQL queries executed and the errors.

**3. Manual migration:**

```bash
# Connect to the database
psql -h localhost -U postgres -d mon-app

# Check the tables
\dt

# Recreate if necessary
DROP TABLE users, refresh_tokens CASCADE;
```

**4. Restart the application:**

```bash
make run
# GORM automatically recreates the tables at startup
```

!!! danger "Production"
    Never use `DropTable` or `AutoMigrate` in production. Use versioned migrations (e.g., golang-migrate).

---

## <i class="material-icons">architecture</i> Architecture & Design

### Why hexagonal architecture?

**Hexagonal architecture** (Ports & Adapters) offers several advantages:

**<i class="material-icons success">check</i> Separation of concerns**
- Isolated business domain
- Interchangeable infrastructure
- Easy to test code

**<i class="material-icons success">check</i> Testability**
- Easy to create mocks (interfaces)
- Unit tests without a database
- Isolated integration tests

**<i class="material-icons success">check</i> Maintainability**
- Localized changes
- Easy to add features
- Safe refactoring

**<i class="material-icons success">check</i> Scalability**
- Easy to add new layers
- Easier migration to microservices

<i class="material-icons info">arrow_forward</i> **Details:** [Generated Project Guide - Architecture](../generated-project-guide.md#architecture)

---

### What is the difference between `models/`, `domain/` and `interfaces/`?

**`internal/models/`** - **Shared entities**
- GORM data structures (User, RefreshToken)
- DTOs (AuthResponse)
- Shared by all modules
- Avoids circular dependencies

**`internal/domain/`** - **Business logic**
- Services (UserService)
- Business rules
- Use cases
- Depends on nothing (except models and interfaces)

**`internal/interfaces/`** - **Ports (contracts)**
- Interfaces defining contracts
- UserRepository (interface)
- Enables dependency inversion

**Example flow:**

```go
Handler → Service (domain/) → Repository (interfaces/) ← Implementation (adapters/repository/)
```

<i class="material-icons info">arrow_forward</i> **Diagrams:** [Generated Project Guide - Hexagonal Architecture](../generated-project-guide.md#architecture-hexagonale-ports--adapters)

---

### When should I use the `minimal` vs `full` template?

**`minimal` template** - Basic REST API:
- <i class="material-icons small">circle</i> Rapid prototyping
- <i class="material-icons small">circle</i> Simple public APIs
- <i class="material-icons small">circle</i> No authentication needed
- <i class="material-icons small">circle</i> Learning / demonstration

**`full` template** (default) - Complete API with JWT:
- <i class="material-icons small">circle</i> Complete backend applications
- <i class="material-icons small">circle</i> Authentication needed
- <i class="material-icons small">circle</i> User management
- <i class="material-icons small">circle</i> Production-ready

**`graphql` template** - GraphQL API:
- <i class="material-icons small">circle</i> Applications requiring GraphQL
- <i class="material-icons small">circle</i> Flexible queries
- <i class="material-icons small">circle</i> React/Vue/Angular frontend

**Example:**

```bash
# Quick prototype without auth
create-go-starter demo-api --template=minimal --database=sqlite

# Production application with auth
create-go-starter mon-app --template=full --database=postgres

# Modern GraphQL API
create-go-starter graphql-api --template=graphql --database=mysql
```

<i class="material-icons info">arrow_forward</i> **Comparison:** [Usage Guide - Available Templates](../usage.md#templates-disponibles)

---

### How do I add a new model (entity)?

**`add-model` command (v1.2.0+):**

```bash
# Simple model
create-go-starter add-model Product --fields "name:string,price:float64"

# With BelongsTo relationship
create-go-starter add-model Comment --fields "content:string" --belongs-to Post

# With HasMany relationship
create-go-starter add-model Category --fields "name:string:unique" --has-many Product
```

**What is automatically generated:**
- <i class="material-icons success">check</i> Model with GORM tags (`internal/models/`)
- <i class="material-icons success">check</i> Repository interface (`internal/interfaces/`)
- <i class="material-icons success">check</i> Repository implementation (`internal/adapters/repository/`)
- <i class="material-icons success">check</i> Service with business logic (`internal/domain/`)
- <i class="material-icons success">check</i> HTTP CRUD handlers (`internal/adapters/handlers/`)
- <i class="material-icons success">check</i> Unit tests
- <i class="material-icons success">check</i> Routes added automatically

<i class="material-icons info">arrow_forward</i> **Details:** [README - Add Models with Relationships](../index.md#ajouter-des-modèles-avec-relations-nouveau-v120)

---

## <i class="material-icons">bug_report</i> Common Errors

### <i class="material-icons error">error</i> Compilation error "cannot find package"

**Cause:** Go dependencies not installed.

**Solution:**

```bash
# Install dependencies
go mod tidy

# Check go.mod and go.sum
cat go.mod

# Clean and reinstall
go clean -modcache
go mod tidy
```

---

### <i class="material-icons error">error</i> Port 8080 already in use

**Cause:** Another process is using port 8080.

**Solution 1 - Change the port in `.env`:**

```bash
# .env
PORT=3000
```

**Solution 2 - Stop the existing process:**

```bash
# Find the process
lsof -i :8080

# Kill the process
kill -9 <PID>
```

**Solution 3 - Use a different port temporarily:**

```bash
PORT=3000 make run
```

---

### <i class="material-icons error">error</i> Tests fail with "database connection failed"

**Cause:** Tests are trying to connect to a real database.

**Solutions:**

**1. Use unit tests with mocks:**

The generated tests use mocks to avoid DB connections:

```go
// Example mock repository
type mockUserRepository struct{}

func (m *mockUserRepository) Create(user *models.User) error {
    return nil
}
```

**2. For integration tests, configure a test DB:**

```bash
# .env.test
DB_NAME=myapp_test
```

```go
// In the test
os.Setenv("DB_NAME", "myapp_test")
```

**3. Use in-memory SQLite for tests:**

```go
dsn := ":memory:"
db, _ := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
```

<i class="material-icons info">arrow_forward</i> **Examples:** [Generated Project Guide - Tests](../generated-project-guide.md#tests)

---

### <i class="material-icons error">error</i> Docker build fails

**Common errors:**

**1. "go.sum not found":**

```bash
# Generate go.sum
go mod tidy

# Then rebuild Docker
docker build -t mon-app .
```

**2. "cannot find package":**

Verify that `go.mod` and `go.sum` are up to date:

```bash
go mod download
go mod verify
docker build -t mon-app .
```

**3. Build too slow:**

Use Docker cache:

```dockerfile
# The generated Dockerfile already uses caching
COPY go.mod go.sum ./
RUN go mod download
# Then copy source code
COPY . .
```

---

## <i class="material-icons">verified_user</i> Best Practices

### How should I organize my code?

**Follow the hexagonal architecture:**

```
internal/
├── models/          # Shared entities (User, Product, etc.)
├── domain/          # Business logic by domain
│   ├── user/        # User domain
│   │   ├── service.go
│   │   └── module.go
│   └── product/     # New domain
│       ├── service.go
│       └── module.go
├── interfaces/      # Ports (interfaces)
│   ├── user_repository.go
│   └── product_repository.go
├── adapters/        # Implementations
│   ├── handlers/    # HTTP handlers
│   ├── repository/  # DB repositories
│   └── middleware/  # Middleware
└── infrastructure/  # Infrastructure
    ├── database/
    └── server/
```

**Rules:**
- <i class="material-icons success">check</i> Domain depends on nothing (except models and interfaces)
- <i class="material-icons success">check</i> One package per business domain
- <i class="material-icons success">check</i> No business logic in handlers
- <i class="material-icons success">check</i> Interfaces for all contracts

<i class="material-icons info">arrow_forward</i> **Details:** [Generated Project Guide - Best Practices](../generated-project-guide.md#bonnes-pratiques)

---

### How do I test my code effectively?

**Unit tests** - Isolated business logic:

```go
// Test the service with mock repository
func TestUserService_Register(t *testing.T) {
    mockRepo := &mockUserRepository{}
    service := NewUserService(mockRepo, logger)

    user, err := service.Register("test@example.com", "password123")

    assert.NoError(t, err)
    assert.NotNil(t, user)
}
```

**Integration tests** - With a real DB:

```go
// Use in-memory SQLite
dsn := ":memory:"
db, _ := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
repo := NewUserRepository(db)

user, err := repo.Create(&models.User{Email: "test@example.com"})
assert.NoError(t, err)
```

**E2E tests** - Complete HTTP requests:

```go
app := fiber.New()
app.Post("/api/v1/auth/register", handler.Register)

req := httptest.NewRequest("POST", "/api/v1/auth/register", body)
resp, _ := app.Test(req)

assert.Equal(t, 201, resp.StatusCode)
```

**Run the tests:**

```bash
make test              # All tests
make test-coverage     # With coverage report
go test ./... -v       # Verbose mode
go test -run TestName  # Specific test
```

---

### What environment variables are needed?

**Essential variables** (`.env`):

```bash
# Application
PORT=8080
ENV=development

# Database (PostgreSQL)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=mon-app
DB_SSLMODE=disable

# Security
JWT_SECRET=<generated with: openssl rand -base64 32>
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d

# Logging
LOG_LEVEL=info
```

**Optional variables:**

```bash
# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://monapp.com

# Rate limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_MAX=100
```

!!! warning "Security"
    - Always use `.env.example` as a template
    - Never commit `.env` to Git
    - Use different secrets per environment

<i class="material-icons info">arrow_forward</i> **Complete configuration:** [Generated Project Guide - Configuration](../generated-project-guide.md#configuration)

---

## <i class="material-icons">cloud_upload</i> Deployment

### How do I deploy to production?

**Option 1 - Docker (Recommended):**

```bash
# Build the image
docker build -t mon-app:latest .

# Run the container
docker run -d \
  -p 8080:8080 \
  --env-file .env.production \
  mon-app:latest
```

**Option 2 - Native binary:**

```bash
# Build for production
make build

# Copy to server
scp mon-app user@server:/opt/mon-app/

# Launch with systemd
sudo systemctl start mon-app
```

**Option 3 - Docker Compose:**

```bash
# On the server
docker-compose -f docker-compose.prod.yml up -d
```

**Deployment checklist:**
- <i class="material-icons small">circle</i> Environment variables configured (`.env.production`)
- <i class="material-icons small">circle</i> JWT secret generated (different from dev)
- <i class="material-icons small">circle</i> Database configured and accessible
- <i class="material-icons small">circle</i> Logs configured (level `info` or `warn`)
- <i class="material-icons small">circle</i> CORS configured for production domain
- <i class="material-icons small">circle</i> Reverse proxy (nginx) configured
- <i class="material-icons small">circle</i> SSL/TLS enabled (Let's Encrypt)
- <i class="material-icons small">circle</i> Monitoring configured

<i class="material-icons info">arrow_forward</i> **Complete guide:** [Generated Project Guide - Deployment](../generated-project-guide.md#déploiement)

---

### How do I manage migrations in production?

**Recommended approach - Versioned migrations:**

Install `golang-migrate`:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

**Create migrations:**

```bash
migrate create -ext sql -dir migrations -seq create_users_table
```

**Apply migrations:**

```bash
migrate -path migrations -database "postgresql://user:pass@localhost:5432/dbname?sslmode=disable" up
```

**Alternative - Manual SQL script:**

```sql
-- migrations/001_initial.sql
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

```bash
psql -h localhost -U postgres -d mon-app < migrations/001_initial.sql
```

!!! danger "Production"
    - Always test migrations in staging first
    - Back up the database before migrating
    - Never use `AutoMigrate` in production

---

### How do I monitor my application?

**Structured logs** (zerolog included):

```go
logger.Info().
    Str("user_id", userID).
    Str("action", "login").
    Msg("User logged in successfully")
```

**Health check endpoint:**

```bash
# Already included in generated projects
curl http://localhost:8080/health
# {"status":"ok"}
```

**External monitoring:**

Integrate with services such as:
- Prometheus + Grafana (metrics)
- Sentry (error tracking)
- DataDog (complete monitoring)
- New Relic (APM)

**Example with Prometheus:**

```go
// Add metrics
import "github.com/prometheus/client_golang/prometheus"

var requestCount = prometheus.NewCounterVec(
    prometheus.CounterOpts{Name: "http_requests_total"},
    []string{"method", "endpoint"},
)
```

<i class="material-icons info">arrow_forward</i> **Details:** [Generated Project Guide - Monitoring](../generated-project-guide.md#monitoring--logging)

---

## <i class="material-icons">terminal</i> CLI - v1.4.0 Features

### How do I use interactive mode?

Interactive mode launches a step-by-step assistant that guides you through configuring your project:

```bash
create-go-starter --interactive
# or with the short alias
create-go-starter -i
```

The assistant asks you questions for each option (project name, template, database, observability) and validates your answers in real time.

<i class="material-icons info">arrow_forward</i> **Details:** [Usage Guide - Interactive Mode](../usage.md#mode-interactif---interactive)

---

### How do I preview files without creating them?

Use **dry-run** mode to see the complete list of files that would be generated, without writing anything to disk:

```bash
create-go-starter mon-app --dry-run
# or with the short alias
create-go-starter mon-app -n

# Combine with other options
create-go-starter mon-app -t minimal -d sqlite -n
```

The dry-run displays the number of files and directories that would be created, with a configuration summary.

<i class="material-icons info">arrow_forward</i> **Details:** [Usage Guide - Dry-Run](../usage.md#prévisualisation-dry-run---dry-run)

---

### How do I diagnose my environment with `doctor`?

The `doctor` command checks that your environment has all the necessary tools:

```bash
create-go-starter doctor
```

It checks:
- <i class="material-icons success">check</i> **Go** (version >= 1.21 required)
- <i class="material-icons success">check</i> **Git** (presence detection)
- <i class="material-icons success">check</i> **Docker** (optional, for PostgreSQL/MySQL)

Each check displays a clear status (OK / WARNING / FAIL) with recommendations if a problem is detected.

<i class="material-icons info">arrow_forward</i> **Details:** [Usage Guide - Doctor Command](../usage.md#commande-doctor)

---

### What short aliases are available?

Since v1.4.0, all CLI options have single-letter short aliases:

| Alias | Long option | Description |
|-------|-------------|-------------|
| `-h` | `--help` | Show help |
| `-i` | `--interactive` | Interactive mode |
| `-n` | `--dry-run` | Preview without writing |
| `-t` | `--template` | Template (minimal, full, graphql) |
| `-d` | `--database` | Database (postgres, mysql, sqlite) |
| `-o` | `--observability` | Observability (none, basic, advanced) |

**Flexible syntax examples:**

```bash
# All these syntaxes are valid
create-go-starter mon-app -t minimal -d sqlite
create-go-starter mon-app -t=minimal -d=sqlite
create-go-starter mon-app --template minimal --database sqlite
```

---

### The progress bar doesn't show up, is that normal?

The visual progress bar is automatically disabled in certain contexts:

- **Non-interactive terminals** (CI/CD, pipes) - automatic detection via `os.Stdout.Fd()`
- **`NO_COLOR` variable set** - respects the [no-color.org](https://no-color.org) convention
- **Output redirection** (`> file`, `| grep`)

In these cases, generation happens silently and the final statistics are always displayed as plain text.

---

## <i class="material-icons">help</i> Miscellaneous Questions

### Can I use `create-go-starter` for commercial projects?

**Yes, absolutely!** The project is under the MIT license, which allows:
- <i class="material-icons success">check</i> Commercial use
- <i class="material-icons success">check</i> Modification
- <i class="material-icons success">check</i> Distribution
- <i class="material-icons success">check</i> Private use

No attribution required (but appreciated!).

---

### How do I contribute to the project?

**Process:**

1. **Fork** the repository
2. **Clone** your fork
3. **Create** a feature branch
4. **Develop** with tests
5. **Commit** following conventions
6. **Push** and create a Pull Request

```bash
git clone https://github.com/VOTRE-USERNAME/go-starter-kit.git
cd go-starter-kit
git checkout -b feature/ma-fonctionnalite

# Develop, test
make test
make lint

# Commit
git commit -m "feat: ajouter fonctionnalité X"

# Push
git push origin feature/ma-fonctionnalite
```

<i class="material-icons info">arrow_forward</i> **Complete guide:** [Contribution Guide](../contributing.md)

---

### Where can I get help?

**Resources:**

- <i class="material-icons">menu_book</i> **Documentation:** [docs/](../)
- <i class="material-icons">bug_report</i> **Issues:** [GitHub Issues](https://github.com/tky0065/go-starter-kit/issues)
- <i class="material-icons">forum</i> **Discussions:** [GitHub Discussions](https://github.com/tky0065/go-starter-kit/discussions)
- <i class="material-icons">chat</i> **Questions:** Open a discussion

**Before asking a question:**
1. Check this FAQ
2. Read the relevant documentation
3. Search existing Issues/Discussions

---

### Is the documentation up to date?

Yes! The documentation is **systematically updated** with every code change.

**Process:**
- Code modified → Documentation updated **immediately**
- Tests documented with examples
- Changes tested locally (`mkdocs serve`)

**If you find an error:**
Open an [Issue](https://github.com/tky0065/go-starter-kit/issues) or [Pull Request](https://github.com/tky0065/go-starter-kit/pulls).

---

## <i class="material-icons">lightbulb</i> Need more help?

**This FAQ didn't answer your question?**

- <i class="material-icons">arrow_forward</i> **Complete documentation:** [Index](../index.md)
- <i class="material-icons">arrow_forward</i> **Installation guide:** [installation.md](../installation.md)
- <i class="material-icons">arrow_forward</i> **Usage guide:** [usage.md](../usage.md)
- <i class="material-icons">arrow_forward</i> **Generated project guide:** [generated-project-guide.md](../generated-project-guide.md)
- <i class="material-icons">arrow_forward</i> **Open a Discussion:** [GitHub Discussions](https://github.com/tky0065/go-starter-kit/discussions)

---

**Last updated:** 2026-02-18
