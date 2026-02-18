# CLI Tool Architecture

Technical documentation for contributors and advanced developers.

## Overview

`create-go-starter` is a Go project generator that creates complete applications with hexagonal architecture, JWT authentication, REST API, and deployment infrastructure.

```
create-go-starter (CLI)
├── main.go              # Entry point, validation, orchestration
├── generator.go         # File generation orchestrator
├── templates.go              # Core templates (config, server, domain, setup.sh)
├── templates_user.go         # User domain specific templates
├── templates_observability.go # Observability templates (Prometheus, Jaeger, Grafana, Health)
├── git.go               # Git repository initialization
├── smoke_test.go        # E2E smoke tests
└── scripts/
    └── smoke_test.sh    # Bash E2E validation script
```

**Statistics**:
- Lines of code: ~6,000+
- Files generated per project: 46+ (60+ with observability)
- Templates: 45+ functions
- Dependencies: Standard library only

## Main Components

### 1. main.go - Entry Point

**Responsibilities**:
- CLI argument parsing (`flag` package)
- Project name validation (regex alphanumeric + hyphens/underscores)
- Directory structure creation
- File generation orchestration
- Error handling and console output (with colors)

**Key Functions**:

```go
func main()
func validateProjectName(name string) error
func createProjectStructure(projectName string) error
func copyEnvFile(projectPath string) error
func printSuccess(projectName string)
```

**Execution Flow**:

```
1. Parse command-line arguments (flag.Parse)
2. Validate project name (alphanumeric + - _)
3. Check if directory already exists
4. Create project directory
5. Create subdirectory structure (cmd/, internal/, pkg/, etc.)
6. Generate all files via generateProjectFiles()
7. Copy .env.example → .env
8. Print success message with next steps
```

**Name Validation**:

```go
// Regex pattern
var validProjectName = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)

// Valid: mon-projet, my-api, user_service, app2024
// Invalid: -projet, mon projet, my.project, @app
```

**Colored Messages**:

```go
func Green(text string) string
func Red(text string) string

// Usage
fmt.Println(Green("✓ Project created successfully!"))
fmt.Println(Red("✗ Error: " + err.Error()))
```

### 2. generator.go - Generation Orchestrator

**Responsibilities**:
- Project directory validation
- Go module name validation (compatible with go.mod)
- Creation of all project files
- Template management with project name injection

**FileGenerator Structure**:

```go
type FileGenerator struct {
    Path    string  // Relative path to project (e.g.: cmd/main.go)
    Content string  // Generated file content
}
```

**Main Function**:

```go
func generateProjectFiles(projectPath, projectName string) error {
    // Validate project directory exists
    // Validate Go module name
    // Create templates instance
    templates := NewProjectTemplates(projectName)

    // Define all files
    files := []FileGenerator{
        {Path: "go.mod", Content: templates.GoModTemplate()},
        {Path: "cmd/main.go", Content: templates.MainGoTemplate()},
        // ... 40+ more files
    }

    // Write all files
    for _, file := range files {
        os.MkdirAll(filepath.Dir(file.Path), 0755)
        os.WriteFile(file.Path, []byte(file.Content), 0644)
    }

    return nil
}
```

### 2.1. Model Generator (add-model Command)

**New in v1.2.0!** The CRUD scaffolding generator allows automatically adding new models to an existing project.

**Responsibilities**:
- Automatic generation of complete CRUD models
- Updating existing files (database, routes, main)
- Support for relations (BelongsTo, HasMany)
- Unit test generation

**File Structure**:

```
cmd/create-go-starter/
├── add_model.go           # CLI orchestration add-model
├── model.go               # Structures: Field, Model, Relation
└── model_generator.go     # 13+ template generators
```

**Data Structures (model.go)**:

```go
// Field represents a model field
type Field struct {
    Name      string   // Field name (e.g.: "Title")
    Type      string   // Go type (e.g.: "string", "int", "time.Time")
    GormTags  []string // GORM tags (e.g.: ["unique", "not null"])
    JSONTag   string   // JSON tag (e.g.: "title")
}

// Model represents the complete model
type Model struct {
    Name       string   // Model name (e.g.: "Todo")
    NamePlural string   // Plural form (e.g.: "Todos")
    Fields     []Field  // Custom fields
    Relations  []Relation // Relations with other models
}

// Relation represents a relation between models
type Relation struct {
    Type        string // "belongs_to" or "has_many"
    ModelName   string // Related model name (e.g.: "Post")
    ForeignKey  string // Foreign key (e.g.: "PostID")
}
```

**CLI Command (add_model.go)**:

```go
func handleAddModel(args []string) error {
    // Parse flags
    var (
        fieldsFlag   string
        belongsTo    string
        hasMany      string
        public       bool
        skipConfirm  bool
    )

    // Validation
    // - Checks that we are in a go-starter-kit project
    // - Checks that the model doesn't already exist
    // - Parses and validates fields
    // - Checks that parent models exist (for relations)

    // Generation
    model := BuildModel(modelName, fields, relations)
    generator := NewModelGenerator(model, projectPath)

    // Generates the 8 files
    generator.GenerateAll()

    // Updates the 3 existing files
    generator.UpdateExistingFiles()

    return nil
}
```

**Template Generators (model_generator.go)**:

```go
type ModelGenerator struct {
    model       Model
    projectPath string
    projectName string
}

// 13+ generation methods

// 1-8: New files
func (g *ModelGenerator) GenerateModel() string
func (g *ModelGenerator) GenerateRepositoryInterface() string
func (g *ModelGenerator) GenerateRepositoryImpl() string
func (g *ModelGenerator) GenerateService() string
func (g *ModelGenerator) GenerateServiceModule() string
func (g *ModelGenerator) GenerateHandler() string
func (g *ModelGenerator) GenerateServiceTests() string
func (g *ModelGenerator) GenerateHandlerTests() string

// 9-13: Updating existing files
func (g *ModelGenerator) UpdateDatabaseMigrations() error
func (g *ModelGenerator) UpdateRoutes() error
func (g *ModelGenerator) UpdateMainModule() error
func (g *ModelGenerator) UpdateParentModel() error
func (g *ModelGenerator) AddBelongsToImport() error
```

**Generation Workflow**:

```
add-model Todo --fields "title:string,completed:bool"
    │
    ├─> 1. Parse and validate arguments
    │   ├─ ModelName: "Todo" → "todos" (pluralization)
    │   ├─ Fields: [{Name:"Title", Type:"string"}, {Name:"Completed", Type:"bool"}]
    │   └─ Validate: No conflicts, valid types
    │
    ├─> 2. Generate 8 new files
    │   ├─ internal/models/todo.go (GORM entity)
    │   ├─ internal/interfaces/todo_repository.go (port)
    │   ├─ internal/adapters/repository/todo_repository.go (GORM impl)
    │   ├─ internal/domain/todo/service.go (business logic)
    │   ├─ internal/domain/todo/module.go (fx module)
    │   ├─ internal/adapters/handlers/todo_handler.go (HTTP)
    │   ├─ internal/domain/todo/service_test.go (tests)
    │   └─ internal/adapters/handlers/todo_handler_test.go (tests)
    │
    └─> 3. Update 3 existing files
        ├─ internal/infrastructure/database/database.go (AutoMigrate)
        ├─ internal/adapters/http/routes.go (CRUD Routes)
        └─ cmd/main.go (fx module registration)
```

**Generation Example (model)**:

Input:
```bash
create-go-starter add-model Product --fields "name:string:unique:not_null,price:float64,stock:int:index"
```

Output (`internal/models/product.go`):
```go
package models

import (
    "time"
    "gorm.io/gorm"
)

type Product struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Name      string         `gorm:"unique;not null" json:"name"`
    Price     float64        `gorm:"not null" json:"price"`
    Stock     int            `gorm:"index;not null" json:"stock"`
    CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
```

**Supported Relations**:

1. **BelongsTo (child → parent)**:
   ```bash
   create-go-starter add-model Comment --fields "content:string" --belongs-to Post
   ```
   - Adds `PostID uint` + `Post Post` to Comment
   - Endpoints: `GET /posts/:postId/comments`, `POST /posts/:postId/comments`

2. **HasMany (parent → children)**:
   ```bash
   create-go-starter add-model Category --fields "name:string" --has-many Product
   ```
   - Adds `Products []Product` to Category
   - Preloading: `GET /categories/:id?include=products`

**Supported Types**:

| CLI Type | Go Type | GORM Tag | Description |
|----------|---------|----------|-------------|
| `string` | `string` | `gorm:"not null"` | Character string |
| `int` | `int` | `gorm:"not null"` | Signed integer |
| `uint` | `uint` | `gorm:"not null"` | Unsigned integer |
| `float64` | `float64` | `gorm:"not null"` | Decimal number |
| `bool` | `bool` | `gorm:"not null"` | Boolean |
| `time` | `time.Time` | `gorm:"not null"` | Date/time |

**GORM Modifiers**:

| Modifier | GORM Tag | Description |
|----------|----------|-------------|
| `unique` | `gorm:"unique"` | Uniqueness constraint |
| `not_null` | `gorm:"not null"` | Non-nullable |
| `index` | `gorm:"index"` | DB index |

**Pluralization**:

The generator uses simple pluralization rules:
- Adds 's': Todo → Todos, Product → Products
- Handles 'y': Category → Categories
- Handles 's/x/z': Class → Classes, Box → Boxes

For irregular plurals (Person→People), edit manually after generation.

**Validation and Safety**:

```go
// Strict validation
- Parent model exists (for --belongs-to)
- Child model exists (for --has-many)
- No file name conflicts
- Valid field types
- No reserved fields (ID, CreatedAt, etc.)

// Safety
- Confirmation prompt (except --yes)
- Dry-run preview of generated files
- Automatic backup of modified files
```

**Test Integration**:

Each generated model includes:
- Service tests: CRUD operations
- Handler tests: HTTP endpoints
- Mock repository (interface)
- Test fixtures

**Complete List of Generated Files** (45+ files):

1. **Configuration** (6 files):
   - go.mod, .env.example, .gitignore, .golangci.yml
   - Makefile, docker-compose.yml

2. **Bootstrap** (1 file):
   - cmd/main.go

3. **Reusable Packages** (7 files):
   - pkg/config/env.go, pkg/config/module.go
   - pkg/logger/logger.go, pkg/logger/module.go
   - pkg/auth/jwt.go, pkg/auth/middleware.go, pkg/auth/module.go

4. **Models** (1 file):
   - internal/models/user.go (User, RefreshToken, AuthResponse)

5. **Domain** (3 files):
   - internal/domain/errors.go
   - internal/domain/user/service.go
   - internal/domain/user/module.go

6. **Interfaces** (1 file):
   - internal/interfaces/user_repository.go

7. **Adapters** (10 files):
   - Handlers: auth_handler.go, user_handler.go, module.go
   - Middleware: error_handler.go
   - Repository: user_repository.go, module.go
   - HTTP: health.go, routes.go

8. **Infrastructure** (2 files):
   - internal/infrastructure/database/database.go
   - internal/infrastructure/server/server.go

9. **Deployment** (3 files):
   - Dockerfile
   - .github/workflows/ci.yml
   - docker-compose.yml

10. **Documentation** (3 files):
    - README.md
    - docs/README.md
    - docs/quick-start.md

**Go module name validation**:

```go
func validateGoModuleName(name string) error {
    // Must start with letter/number
    // Only alphanumeric, hyphens, underscores
    // No spaces, special chars
    pattern := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)
    return pattern.MatchString(name)
}
```

### 3. templates.go - Main Templates

**Responsibilities**:
- Definition of templates for infrastructure and configuration
- Templates: configuration, server, database, middleware, Docker, CI/CD
- Dynamic injection of the project name

**ProjectTemplates Structure**:

```go
type ProjectTemplates struct {
    projectName string
}

func NewProjectTemplates(projectName string) *ProjectTemplates {
    return &ProjectTemplates{projectName: projectName}
}
```

**Main Methods** (30+ templates):

#### Configuration & Build

```go
func (t *ProjectTemplates) GoModTemplate() string
func (t *ProjectTemplates) MakefileTemplate() string
func (t *ProjectTemplates) DockerfileTemplate() string
func (t *ProjectTemplates) DockerComposeTemplate() string
func (t *ProjectTemplates) GolangCILintTemplate() string
func (t *ProjectTemplates) GitHubActionsWorkflowTemplate() string
func (t *ProjectTemplates) EnvTemplate() string
func (t *ProjectTemplates) GitignoreTemplate() string
```

#### Bootstrap

```go
func (t *ProjectTemplates) UpdatedMainGoTemplate() string // fx.New bootstrap
```

#### Packages (pkg/)

```go
func (t *ProjectTemplates) ConfigTemplate() string          // pkg/config/env.go
func (t *ProjectTemplates) LoggerTemplate() string          // pkg/logger/logger.go
func (t *ProjectTemplates) JWTAuthTemplate() string         // pkg/auth/jwt.go
func (t *ProjectTemplates) JWTMiddlewareTemplate() string   // pkg/auth/middleware.go
func (t *ProjectTemplates) AuthModuleTemplate() string      // pkg/auth/module.go
```

#### Infrastructure

```go
func (t *ProjectTemplates) DatabaseTemplate() string        // database.go
func (t *ProjectTemplates) ServerTemplate() string          // server.go with Fiber
```

#### Middleware

```go
func (t *ProjectTemplates) ErrorHandlerMiddlewareTemplate() string
```

#### Health Check & Routes

```go
func (t *ProjectTemplates) HealthHandlerTemplate() string    // health.go
func (t *ProjectTemplates) RoutesTemplate() string           // routes.go - Centralized routes
```

#### Documentation

```go
func (t *ProjectTemplates) ReadmeTemplate() string
func (t *ProjectTemplates) DocsReadmeTemplate() string      // docs/README.md
func (t *ProjectTemplates) QuickStartTemplate() string      // docs/quick-start.md
```

#### Setup & Automation

```go
func (t *ProjectTemplates) SetupScriptTemplate() string     // setup.sh - Automated setup
```

**SetupScriptTemplate** generates an interactive bash script that:
- Checks prerequisites (Go, OpenSSL, Docker, psql)
- Installs dependencies (`go mod tidy`)
- Generates and configures the JWT secret in `.env`
- Configures PostgreSQL (Docker or local)
- Runs validation tests
- Verifies the complete installation

The script is automatically made executable (`chmod 0755`) by `generator.go`.

**Template Pattern**:

Templates use string concatenation (not text/template):

```go
func (t *ProjectTemplates) GoModTemplate() string {
    return `module ` + t.projectName + `

go 1.25.5

require (
    github.com/gofiber/fiber/v2 v2.52.10
    // ...
)
`
}
```

**Advantages**:
- Simplicity (no parsing)
- Type-safe at compile time
- Easy to debug
- Direct projectName injection

**Disadvantages**:
- Verbose for complex templates
- Manual backtick escaping

### 4. templates_user.go - User Domain Templates

**Responsibilities**:
- Templates specific to the User domain
- Entities (models), services, repositories, handlers
- Tests (if implemented)

**Methods**:

#### Models (Shared Entities)

```go
func (t *ProjectTemplates) ModelsUserTemplate() string  // User, RefreshToken, AuthResponse
```

#### Domain

```go
func (t *ProjectTemplates) DomainErrorsTemplate() string
func (t *ProjectTemplates) UserServiceTemplate() string   // Business logic
func (t *ProjectTemplates) UserModuleTemplate() string    // fx module
```

**Note**: `UserEntityTemplate()` and `UserRefreshTokenTemplate()` are deprecated - entities are now in `ModelsUserTemplate()`.

#### Interfaces (Ports)

```go
func (t *ProjectTemplates) UserInterfacesTemplate() string
func (t *ProjectTemplates) UserRepositoryInterfaceTemplate() string
```

#### Adapters

```go
func (t *ProjectTemplates) UserRepositoryTemplate() string  // GORM implementation
func (t *ProjectTemplates) RepositoryModuleTemplate() string
func (t *ProjectTemplates) AuthHandlerTemplate() string     // Register, Login, Refresh
func (t *ProjectTemplates) UserHandlerTemplate() string     // CRUD endpoints
func (t *ProjectTemplates) HandlerModuleTemplate() string
```

**Typical Template Content** (example: ModelsUserTemplate):

```go
func (t *ProjectTemplates) ModelsUserTemplate() string {
    return `package models

import (
    "time"
    "gorm.io/gorm"
)

// User represents the domain entity for a user
type User struct {
    ID           uint           ` + "`gorm:\"primaryKey\" json:\"id\"`" + `
    Email        string         ` + "`gorm:\"uniqueIndex;not null\" json:\"email\"`" + `
    PasswordHash string         ` + "`gorm:\"not null\" json:\"-\"`" + `
    CreatedAt    time.Time      ` + "`gorm:\"autoCreateTime\" json:\"created_at\"`" + `
    UpdatedAt    time.Time      ` + "`gorm:\"autoUpdateTime\" json:\"updated_at\"`" + `
    DeletedAt    gorm.DeletedAt ` + "`gorm:\"index\" json:\"deleted_at,omitempty\"`" + `
}

// RefreshToken represents a refresh token for session management
type RefreshToken struct {
    ID        uint      ` + "`gorm:\"primaryKey\" json:\"id\"`" + `
    UserID    uint      ` + "`gorm:\"not null;index\" json:\"user_id\"`" + `
    Token     string    ` + "`gorm:\"uniqueIndex;not null\" json:\"token\"`" + `
    ExpiresAt time.Time ` + "`gorm:\"not null\" json:\"expires_at\"`" + `
    Revoked   bool      ` + "`gorm:\"not null;default:false\" json:\"revoked\"`" + `
    CreatedAt time.Time ` + "`gorm:\"autoCreateTime\" json:\"created_at\"`" + `
    UpdatedAt time.Time ` + "`gorm:\"autoUpdateTime\" json:\"updated_at\"`" + `
}

func (rt *RefreshToken) IsExpired() bool {
    return time.Now().After(rt.ExpiresAt)
}

func (rt *RefreshToken) IsRevoked() bool {
    return rt.Revoked
}

// AuthResponse represents the authentication response with tokens
type AuthResponse struct {
    AccessToken  string ` + "`json:\"access_token\"`" + `
    RefreshToken string ` + "`json:\"refresh_token\"`" + `
    ExpiresIn    int64  ` + "`json:\"expires_in\"`" + `
}
`
}
```

**Why `models` instead of `domain/user`?**
- **Avoids circular dependencies**: Interfaces can reference models without creating cycles
- **Centralization**: Entities are defined in a single location
- **Reusability**: All layers (domain, interfaces, adapters) can import models without conflict

### 5. templates_observability.go - Observability Templates

**New in v1.3.0!** This file contains all templates for the advanced observability stack (~1369 lines).

**Responsibilities**:
- Templates for Prometheus metrics (`/metrics` endpoint, middleware, registry)
- Templates for OpenTelemetry distributed tracing (Jaeger export, OTLP/gRPC)
- Templates for advanced K8s health checks (liveness/readiness)
- Templates for Grafana dashboard (JSON, provisioning, alerting)
- Templates for observability Docker Compose (Jaeger, Prometheus, Grafana)
- Templates for Kubernetes probes configuration

**Activation Condition**: Only generated when `--observability=advanced` AND `--template=full`.

**Main Methods** (15+ templates):

#### Prometheus Metrics

```go
func (t *ProjectTemplates) PrometheusTemplate() string           // pkg/metrics/prometheus.go
func (t *ProjectTemplates) MetricsMiddlewareTemplate() string    // middleware/metrics_middleware.go
func (t *ProjectTemplates) MetricsHandlerTemplate() string       // handlers/metrics_handler.go
```

#### Distributed Tracing (OpenTelemetry)

```go
func (t *ProjectTemplates) TracerTemplate() string               // pkg/tracing/tracer.go
func (t *ProjectTemplates) TracingMiddlewareTemplate() string    // middleware/tracing_middleware.go
func (t *ProjectTemplates) GORMTracingTemplate() string          // database/tracing.go
func (t *ProjectTemplates) LoggerWithTracingTemplate() string    // logger/logger_tracing.go
```

#### Advanced Health Checks

```go
func (t *ProjectTemplates) AdvancedHealthHandlerTemplate() string // http/health.go (replaces the basic one)
```

#### Grafana Dashboard & Configuration

```go
func (t *ProjectTemplates) GrafanaDashboardJSONTemplate() string       // monitoring/grafana/dashboards/app-dashboard.json
func (t *ProjectTemplates) GrafanaDatasourceTemplate() string          // monitoring/grafana/provisioning/datasources/prometheus.yml
func (t *ProjectTemplates) GrafanaDashboardProvisioningTemplate() string // monitoring/grafana/provisioning/dashboards/dashboard.yml
```

#### Prometheus & Alerting Configuration

```go
func (t *ProjectTemplates) PrometheusConfigTemplate() string     // monitoring/prometheus/prometheus.yml
func (t *ProjectTemplates) PrometheusAlertRulesTemplate() string // monitoring/prometheus/alert_rules.yml
```

#### Kubernetes Probes

```go
func (t *ProjectTemplates) KubernetesProbesTemplate() string     // deployments/kubernetes/probes.yaml
```

#### Observability Docker Compose

```go
func (t *ProjectTemplates) DockerComposeObservabilityTemplate() string // docker-compose.yml (extended version)
```

**Files Generated by Advanced Observability** (14+ additional files):

| Category | Files | Description |
|----------|-------|-------------|
| Metrics | 3 | prometheus.go, metrics_middleware.go, metrics_handler.go |
| Tracing | 4 | tracer.go, tracing_middleware.go, gorm_tracing.go, logger_tracing.go |
| Health | 1 | health.go (advanced version with liveness/readiness) |
| Grafana | 3 | dashboard JSON, datasource YAML, provisioning YAML |
| Prometheus | 2 | prometheus.yml, alert_rules.yml |
| Kubernetes | 1 | probes.yaml |

**Dependencies Added to go.mod** (only when observability=advanced):

```
github.com/ansrivas/fiberprometheus/v2 v2.7.0
go.opentelemetry.io/otel v1.24.0
go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.24.0
go.opentelemetry.io/otel/sdk v1.24.0
go.opentelemetry.io/otel/trace v1.24.0
go.opentelemetry.io/contrib/instrumentation/github.com/gofiber/fiber/otelfiber v0.49.0
```

### 6. git.go - Git Initialization

**Responsibilities**:
- Checking Git availability on the system
- Automatic initialization of a Git repository in the generated project
- Creating an initial commit with all generated files

**Key Functions**:

```go
func isGitAvailable() bool           // Checks if git is installed
func initGitRepo(projectPath string) error  // Initializes the repo and creates the initial commit
```

**Behavior**:
- If Git is available: initializes the repo and creates a commit "Initial commit from go-starter-kit"
- If Git is not available: displays a warning but continues (graceful degradation)
- The `.gitignore` is automatically added before the initial commit

**Integration**:
- Called in `main.go` after `copyEnvFile()` and before `printSuccessMessage()`
- Messages: "<i class="material-icons">build</i> Setting up Git repository..." and "<i class="material-icons success">check_circle</i> Git repository initialized"

## Patterns and Conventions

### 1. Template Pattern

**Choice: String concatenation vs text/template**

**Chosen option**: String concatenation

```go
return `package main

import "fmt"

func main() {
    fmt.Println("` + t.projectName + `")
}
`
```

**Why not text/template?**
- Simplicity: No parsing, no data struct
- Performance: No parsing overhead
- Type-safety: Compile-time errors
- Debugging: Easier to see the generated template

**Challenges**:
- Backtick escaping: Use `` "`" ``
- Long templates become verbose

### 2. Layered Validation

**Layer 1 - CLI level (main.go)**:
- Project name validation
- Regex: `^[a-zA-Z0-9][a-zA-Z0-9_-]*$`
- Valid examples: `mon-projet`, `my_app`, `app2024`

**Layer 2 - Generator level (generator.go)**:
- Go module name validation (same rules as Layer 1)
- Directory existence verification

**Layer 3 - Runtime (generated code)**:
- Validation with go-playground/validator
- Business validation in the domain

### 3. Error Handling

**Convention**:

```go
if err != nil {
    return fmt.Errorf("context: %w", err)  // Wrap error with context
}
```

**User Display**:

```go
fmt.Println(Red("✗ Error: " + err.Error()))
os.Exit(1)
```

**No panic**: Use return error, not panic()

### 4. Tests

**Organization**:

```
cmd/create-go-starter/
├── main.go
├── main_test.go           # CLI tests and run() function
├── generator.go
├── generator_test.go      # Generation tests
├── templates.go
├── templates_test.go      # Template tests
├── templates_user.go
├── git.go
├── git_test.go            # Git initialization tests
├── colors_test.go         # Color utilities tests
├── env_test.go            # .env copy tests
├── scaffold_test.go       # Structure creation tests
└── smoke_test.go          # E2E smoke tests
scripts/
└── smoke_test.sh          # Bash E2E validation script
```

**Test Coverage**: 83%+

**Makefile Commands**:
```bash
make test              # All tests
make test-short        # Quick tests (skip long tests)
make smoke-test        # Full E2E validation with runtime
make smoke-test-quick  # E2E validation without runtime (no Docker)
```

**Test Patterns**:

1. **Table-driven tests**:

```go
func TestValidateProjectName(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid simple", "myproject", false},
        {"valid with hyphen", "my-project", false},
        {"invalid space", "my project", true},
        {"invalid special char", "my@project", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateProjectName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("validateProjectName() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

2. **Template validation**:

```go
func TestReadmeTemplate(t *testing.T) {
    templates := NewProjectTemplates("test-project")
    content := templates.ReadmeTemplate()

    if !strings.Contains(content, "# test-project") {
        t.Error("README should contain project name as title")
    }
}
```

## Extensibility

### Adding a New Template

**Steps**:

1. **Create the template method** in `templates.go` or `templates_user.go`:

```go
func (t *ProjectTemplates) MyNewTemplate() string {
    return `package mypackage

// Generated code for ` + t.projectName + `

func Hello() {
    fmt.Println("Hello from ` + t.projectName + `")
}
`
}
```

2. **Add FileGenerator** in `generator.go`:

```go
files := []FileGenerator{
    // ... existing files
    {
        Path:    filepath.Join(projectPath, "pkg", "mypackage", "hello.go"),
        Content: templates.MyNewTemplate(),
    },
}
```

3. **Test the generation**:

```bash
go run ./cmd/create-go-starter test-project
cat test-project/pkg/mypackage/hello.go
```

4. **Add a test** in `templates_test.go`:

```go
func TestMyNewTemplate(t *testing.T) {
    templates := NewProjectTemplates("test")
    content := templates.MyNewTemplate()

    assert.Contains(t, content, "package mypackage")
    assert.Contains(t, content, "test")
}
```

### Adding a CLI Option

**Example: Adding a `--database` flag to choose the DB**

1. **Define the flag** in `main.go`:

```go
var database string

func main() {
    flag.StringVar(&database, "database", "postgres", "Database type (postgres, mysql, sqlite)")
    flag.Parse()

    // Validate
    if database != "postgres" && database != "mysql" && database != "sqlite" {
        fmt.Println(Red("✗ Invalid database type"))
        os.Exit(1)
    }
}
```

2. **Pass to generateProjectFiles**:

```go
func generateProjectFiles(projectPath, projectName, database string) error {
    templates := NewProjectTemplates(projectName)
    templates.database = database  // Add field to struct

    // Conditional template generation
    switch database {
    case "postgres":
        // Generate PostgreSQL specific files
    case "mysql":
        // Generate MySQL specific files
    }
}
```

3. **Adapt the templates**:

```go
func (t *ProjectTemplates) GoModTemplate() string {
    driver := "gorm.io/driver/postgres"
    if t.database == "mysql" {
        driver = "gorm.io/driver/mysql"
    }

    return `module ` + t.projectName + `

require (
    gorm.io/gorm v1.31.1
    ` + driver + ` v1.5.11
)
`
}
```

### Future Extension Models

**1. Multiple Templates**:

```bash
create-go-starter my-project --template=minimal
create-go-starter my-project --template=full
create-go-starter my-project --template=api-only
create-go-starter my-project --template=graphql
```

**Implementation**:

```go
// templates/minimal.go
type MinimalTemplates struct { ... }

// templates/full.go
type FullTemplates struct { ... }

// Factory pattern
func NewTemplates(projectName, templateType string) TemplateGenerator {
    switch templateType {
    case "minimal":
        return NewMinimalTemplates(projectName)
    case "full":
        return NewFullTemplates(projectName)
    }
}
```

**2. DB Choice**:

```bash
create-go-starter my-project --db=mysql
create-go-starter my-project --db=mongodb
```

**3. Framework Choice**:

```bash
create-go-starter my-project --framework=gin
create-go-starter my-project --framework=echo
```

**4. Interactive CLI**:

```bash
create-go-starter

? Project name: my-awesome-api
? Database: PostgreSQL
? Auth: JWT
? Generate Swagger docs: Yes
? Include Docker: Yes

✓ Generating project...
```

Use: `github.com/manifoldco/promptui`

## Dependencies

**CLI Dependencies**: **NONE** (stdlib only)

**Advantages**:
- Simplicity: No go mod tidy for the CLI
- Portability: Static binary with no dependencies
- Lightweight installation: Very fast `go install`
- Easy maintenance: No external breaking changes

**Stdlib Packages Used**:
- `flag` - CLI parsing
- `fmt` - Formatting and printing
- `os` - File operations, exit codes
- `path/filepath` - Path manipulation
- `regexp` - Validation patterns
- `strings` - String utilities

## Performance

**Metrics**:

- **Generation time**: < 1 second for 45+ files
- **Binary size**: ~3-4 MB (static)
- **Memory**: < 10 MB during generation
- **Disk writes**: 45+ files, ~15,000 lines of generated code

**Optimizations**:

1. **No template parsing**: Direct string concatenation
2. **Batch file creation**: All files created in a single pass
3. **MkdirAll once**: Creates all parent directories if needed
4. **No dependencies**: No download or import overhead

## Code Standards

**Conventions Followed**:

1. **gofmt**: Always formatted
   ```bash
   go fmt ./...
   ```

2. **golangci-lint**: Quality checks
   ```bash
   golangci-lint run ./...
   ```

3. **Test coverage**: > 80%
   ```bash
   go test -cover ./...
   ```

4. **GoDoc documentation**: For public exports
   ```go
   // ProjectTemplates holds all template generation methods
   type ProjectTemplates struct {
       projectName string
   }
   ```

5. **Error handling**: Always explicit, never ignored
   ```go
   if err != nil {
       return fmt.Errorf("context: %w", err)
   }
   ```

## Debugging

**Enable verbose mode** (to be implemented):

```go
var verbose bool
flag.BoolVar(&verbose, "verbose", false, "Verbose output")

if verbose {
    log.Println("Creating directory:", path)
    log.Println("Writing file:", filepath)
}
```

**Test generation**:

```bash
# Generate test project
go run ./cmd/create-go-starter test-project

# Check files
ls -la test-project/
tree test-project/

# Check content
cat test-project/go.mod
cat test-project/cmd/main.go

# Test build
cd test-project
go mod tidy
go build ./...

# Clean up
cd ..
rm -rf test-project
```

**Debug with Delve**:

```bash
dlv debug ./cmd/create-go-starter -- my-project
(dlv) break main.createProjectStructure
(dlv) continue
```

## Contributing

To contribute to the CLI:

1. Fork the repository
2. Create a branch: `git checkout -b feature/my-feature`
3. Make changes
4. Tests: `go test ./...`
5. Lint: `golangci-lint run`
6. Commit: `git commit -m "feat: add feature"`
7. Push: `git push origin feature/my-feature`
8. Open a Pull Request

**PR Checklist**:
- [ ] Tests added/updated
- [ ] Tests pass (`make test`)
- [ ] Lint passes (`make lint`)
- [ ] Documentation updated
- [ ] Clear commit messages (conventional commits)

## Technical Roadmap

**Completed**:
- [x] Dry-run mode (`--dry-run`)
- [x] Multiple templates (minimal, full, graphql)
- [x] DB choice (PostgreSQL, MySQL, SQLite)
- [x] Interactive CLI (`--interactive`)
- [x] Short aliases (`-t`, `-d`, `-o`, `-i`, `-n`, `-h`)
- [x] `doctor` command for diagnostics
- [x] Progress bar and statistics
- [x] Advanced observability (Prometheus, Jaeger, Grafana)

**Short term**:
- [ ] Version flag (`--version`)
- [ ] Verbose mode (`--verbose`)
- [ ] Force overwrite (`--force`)

**Medium term**:
- [ ] Framework choice (Fiber, Gin, Echo)
- [ ] Colored diff for `--dry-run`

**Long term**:
- [ ] Plugin system for custom templates
- [ ] Template marketplace
- [ ] Hot-reload for templates
- [ ] GUI for generation

---

This technical documentation should allow contributors to understand the internal architecture of the CLI and contribute effectively to the project.
