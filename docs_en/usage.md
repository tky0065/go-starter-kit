# Usage Guide

This guide details how to use `create-go-starter` and provides an in-depth explanation of the structure of generated projects.

## Basic Command

The basic syntax is very simple:

```bash
create-go-starter <nom-du-projet>
```

### Example

```bash
create-go-starter mon-api-backend
```

This command will create a new `mon-api-backend/` directory with the entire project structure using the **full** template by default.

## Available Templates

`create-go-starter` offers **three templates** to meet different project needs. Choose a template with the `--template` flag:

```bash
create-go-starter mon-projet --template minimal    # API REST basique
create-go-starter mon-projet --template full       # API complète avec auth (défaut)
create-go-starter mon-projet --template graphql    # API GraphQL
```

### Templates Overview

| Template | Description | Use Case |
|----------|-------------|----------|
| `minimal` | Basic REST API with Swagger (no authentication) | Quick prototypes, simple public APIs, microservices without auth |
| `full` | Complete API with JWT auth, user management and Swagger (**default**) | Complete backend applications, APIs requiring authentication |
| `graphql` | GraphQL API with gqlgen and GraphQL Playground | Applications requiring GraphQL, modern frontend clients |

### Detailed Feature Comparison

| Feature | minimal | full | graphql |
|---------|---------|------|---------|
| **REST API** | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> | <i class="material-icons error">cancel</i> |
| **GraphQL API** | <i class="material-icons error">cancel</i> | <i class="material-icons error">cancel</i> | <i class="material-icons success">check_circle</i> |
| **JWT Authentication** | <i class="material-icons error">cancel</i> | <i class="material-icons success">check_circle</i> | <i class="material-icons error">cancel</i> |
| **User Management** | <i class="material-icons error">cancel</i> | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> |
| **Swagger Documentation** | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> | <i class="material-icons error">cancel</i> |
| **GraphQL Playground** | <i class="material-icons error">cancel</i> | <i class="material-icons error">cancel</i> | <i class="material-icons success">check_circle</i> |
| **Database (GORM)** | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> |
| **PostgreSQL** | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> |
| **Dependency Injection (fx)** | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> |
| **Structured Logging (zerolog)** | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> |
| **Hexagonal Architecture** | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> |
| **Unit Tests** | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> |
| **Docker** | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> |
| **CI/CD (GitHub Actions)** | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> | <i class="material-icons success">check_circle</i> |

### Major Structural Differences

#### `minimal` Template

**Characteristics**:
- Simple REST API with basic CRUD endpoints
- No authentication or authorization
- Swagger for API documentation
- Perfect for getting started quickly without complexity

**Specific Structure**:
- `internal/adapters/handlers/user_handler.go` - Simple CRUD handlers without auth
- `internal/domain/user/service.go` - Basic business logic
- Automatic Swagger documentation via annotations

**Generated Endpoints**:
```
GET    /health                  # Health check (alias liveness — rétrocompatibilité)
GET    /health/liveness         # Liveness probe K8s (toujours 200 si app tourne)
GET    /health/readiness        # Readiness probe K8s (200 si DB ok, 503 si DB down)
GET    /api/v1/users            # Liste tous les utilisateurs
GET    /api/v1/users/:id        # Récupère un utilisateur
POST   /api/v1/users            # Crée un utilisateur
PUT    /api/v1/users/:id        # Met à jour un utilisateur
DELETE /api/v1/users/:id        # Supprime un utilisateur
GET    /swagger/*               # Documentation Swagger UI
```

**Recommended Use Cases**:
- Quick prototypes and POCs
- Public APIs without sensitive data
- Internal microservices without authentication needs
- Learning hexagonal architecture

---

#### `full` Template (default)

**Characteristics**:
- Complete REST API with JWT authentication
- Auth system with access tokens + refresh tokens
- Complete user management (CRUD + auth)
- Swagger with Bearer token authentication
- Production-ready with built-in security

**Specific Structure**:
- `internal/adapters/handlers/auth_handler.go` - Register, login, refresh endpoints
- `internal/adapters/handlers/user_handler.go` - JWT-protected CRUD
- `internal/adapters/middleware/auth_middleware.go` - JWT verification
- `pkg/auth/` - JWT token generation and validation
- `internal/models/user.go` - User + RefreshToken with bcrypt

**Generated Endpoints**:
```
GET    /health                      # Health check (alias liveness — rétrocompatibilité)
GET    /health/liveness             # Liveness probe K8s
GET    /health/readiness            # Readiness probe K8s (vérifie DB)
POST   /api/v1/auth/register        # Inscription utilisateur
POST   /api/v1/auth/login           # Connexion (retourne access + refresh tokens)
POST   /api/v1/auth/refresh         # Rafraîchir l'access token
GET    /api/v1/users                # Liste utilisateurs (<i class="material-icons warning small">lock</i> JWT requis)
GET    /api/v1/users/:id            # Récupère utilisateur (<i class="material-icons warning small">lock</i> JWT requis)
PUT    /api/v1/users/:id            # Met à jour utilisateur (<i class="material-icons warning small">lock</i> JWT requis)
DELETE /api/v1/users/:id            # Supprime utilisateur (<i class="material-icons warning small">lock</i> JWT requis)
GET    /swagger/*                   # Documentation Swagger UI
```

**Recommended Use Cases**:
- Complete backend applications
- APIs requiring authentication and authorization
- SaaS and multi-user applications
- Publicly exposed APIs with sensitive data

---

#### `graphql` Template

**Characteristics**:
- Complete GraphQL API with gqlgen
- GraphQL Playground to interactively explore the API
- Typed GraphQL schema with resolvers
- User management with mutations and queries
- Hexagonal architecture adapted for GraphQL

**Specific Structure**:
- `graph/schema.graphqls` - GraphQL schema (types, queries, mutations)
- `graph/resolver.go` - Main resolver
- `graph/schema.resolvers.go` - Resolver implementation
- `gqlgen.yml` - gqlgen configuration
- `internal/infrastructure/server/server.go` - GraphQL server with Playground

**Generated GraphQL Schema**:
```graphql
type User {
  id: ID!
  email: String!
  createdAt: Time!
  updatedAt: Time!
}

type Query {
  users(limit: Int, offset: Int): [User!]!
  user(id: ID!): User
}

type Mutation {
  createUser(input: CreateUserInput!): User!
  updateUser(id: ID!, input: UpdateUserInput!): User!
  deleteUser(id: ID!): Boolean!
}
```

**Generated Endpoints**:
```
GET    /health                  # Health check
POST   /graphql                 # Endpoint GraphQL
GET    /                        # GraphQL Playground UI
```

**Recommended Use Cases**:
- Modern frontend applications (React, Vue, Angular)
- APIs requiring flexible queries
- Mobile applications with specific data needs
- Projects favoring GraphQL over REST

---

### How to Choose the Right Template?

**Choose `minimal` if**:
- You want a quick prototype
- Your API is public without sensitive data
- You don't need authentication
- You want to learn hexagonal architecture simply

**Choose `full` if**:
- You're building a complete backend application
- You need JWT authentication
- You want users with login/register
- You prefer REST over GraphQL
- You want a production-ready project immediately

**Choose `graphql` if**:
- You're building a GraphQL API
- Your frontend uses Apollo Client, Relay or urql
- You prefer a strongly typed schema
- You want GraphQL Playground for exploration
- Your clients have varying data needs



## Supported Databases

Go Starter Kit supports **3 database types** to fit your needs:

### Overview

| Database | Ideal For | Configuration | Command |
|----------|-----------|---------------|---------|
| **PostgreSQL** | Production, complex queries | Docker | `create-go-starter mon-app` |
| **MySQL** | Broad compatibility, shared hosting | Docker | `create-go-starter mon-app --database=mysql` |
| **SQLite** | Prototyping, small apps | None | `create-go-starter mon-app --database=sqlite` |

### PostgreSQL (default)

```bash
create-go-starter mon-app
# OU explicitement:
create-go-starter mon-app --database=postgres
```

**Advantages**:
- <i class="material-icons success small">check</i> Advanced SQL features (JSON, arrays, full-text search)
- <i class="material-icons success small">check</i> Excellent performance and reliability
- <i class="material-icons success small">check</i> ACID compliant, strong data integrity
- <i class="material-icons success small">check</i> Ideal for production

**When to use**: Production applications, complex data, reliability requirements

### MySQL/MariaDB

```bash
create-go-starter mon-app --database=mysql
```

**Advantages**:
- <i class="material-icons success small">check</i> Broad compatibility with hosting providers
- <i class="material-icons success small">check</i> Excellent for read-heavy workloads
- <i class="material-icons success small">check</i> Mature and well-documented ecosystem

**When to use**: Shared hosting, MySQL-experienced teams, broad compatibility

### SQLite

```bash
create-go-starter mon-app --database=sqlite
```

**Advantages**:
- <i class="material-icons success small">check</i> Zero configuration (no server needed)
- <i class="material-icons success small">check</i> Perfect for quick prototyping
- <i class="material-icons success small">check</i> Database in a single file
- <i class="material-icons success small">check</i> Very fast for small datasets

**Limitations**:
- <i class="material-icons warning small">warning</i> Limited concurrent writes
- <i class="material-icons warning small">warning</i> Not suited for large-scale production

**When to use**: Prototyping, MVPs, development, embedded apps, small production (<100 users)

### Complete Guide

For a detailed comparison, configuration examples, and migration guides, see:

- **[Database Selection Guide](databases.md)** - Complete comparison and selection help
- **[Migration Guide](database-migration.md)** - Migrating between databases

## Observability (--observability)

<i class="material-icons success">new_releases</i> **New in v1.3.0!** Go Starter Kit supports **3 observability levels** to monitor your generated projects in production.

### Available Levels

| Level | Description | What Is Generated |
|-------|-------------|-------------------|
| `none` | No observability (default) | Standard behavior preserved |
| `basic` | Advanced K8s health checks | `/health/liveness`, `/health/readiness` |
| `advanced` | Complete observability stack | Prometheus + Jaeger + Grafana + K8s Health Checks |

> <i class="material-icons warning">warning</i> The `--observability` flag only works with `--template=full`.

### `advanced` Mode — Complete Observability Stack

```bash
# Générer un projet avec observabilité avancée
create-go-starter mon-app --observability=advanced

# Combiné avec database et template
create-go-starter mon-app --template=full --database=postgres --observability=advanced
```

**Generated Files** (in addition to standard files):

```
mon-app/
├── pkg/
│   ├── metrics/
│   │   └── prometheus.go                # Registry Prometheus + métriques HTTP
│   └── tracing/
│       └── tracer.go                    # Configuration OpenTelemetry + TracerProvider
├── internal/adapters/
│   ├── middleware/
│   │   ├── metrics_middleware.go        # Middleware Prometheus pour métriques HTTP
│   │   └── tracing_middleware.go        # Middleware OpenTelemetry pour spans HTTP
│   ├── handlers/
│   │   └── metrics_handler.go           # Handler GET /metrics
│   └── http/
│       └── health.go                    # Health checks avancés (liveness/readiness)
├── internal/infrastructure/database/
│   └── tracing.go                       # Instrumentation GORM avec spans DB
├── pkg/logger/
│   └── logger_tracing.go               # Logger enrichi avec trace_id/span_id
├── monitoring/
│   ├── grafana/
│   │   ├── provisioning/
│   │   │   ├── datasources/
│   │   │   │   └── prometheus.yml       # Datasource Prometheus auto-configurée
│   │   │   └── dashboards/
│   │   │       └── dashboard.yml        # Auto-provisioning des dashboards
│   │   └── dashboards/
│   │       └── app-dashboard.json       # Dashboard Grafana 7 panneaux
│   └── prometheus/
│       ├── prometheus.yml               # Configuration scraping Prometheus
│       └── alert_rules.yml              # Règles d'alerting
├── deployments/
│   └── kubernetes/
│       └── probes.yaml                  # Configuration probes K8s
└── docker-compose.yml                   # Stack complète (DB + Jaeger + Prometheus + Grafana)
```

#### Prometheus Metrics

**Metrics exposed on `GET /metrics`**:

| Metric | Type | Description |
|--------|------|-------------|
| `http_requests_total` | Counter | Total requests by method, route and status code |
| `http_request_duration_seconds` | Histogram | HTTP latency by route (p50, p90, p95, p99) |
| `http_requests_in_flight` | Gauge | Number of active in-progress requests |

**Library**: [`fiberprometheus/v2`](https://github.com/ansrivas/fiberprometheus) v2.7.0.

```bash
# Tester les métriques
curl http://localhost:8080/metrics
```

#### Distributed Tracing (OpenTelemetry + Jaeger)

**Architecture**: Application → OTLP/gRPC → Jaeger Collector → Jaeger UI

**Features**:
- <i class="material-icons success small">check</i> Automatic spans for each HTTP request
- <i class="material-icons success small">check</i> Automatic spans for each GORM query
- <i class="material-icons success small">check</i> W3C `traceparent` propagation between services
- <i class="material-icons success small">check</i> zerolog logs enriched with `trace_id` and `span_id`

**Environment Variables**:
```bash
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
OTEL_SERVICE_NAME=mon-app
```

**Access Jaeger UI**:
```bash
open http://localhost:16686
```

#### Kubernetes-compatible Health Checks

| Endpoint | K8s Probe | Behavior |
|----------|-----------|----------|
| `GET /health/liveness` | livenessProbe | 200 if the app is running |
| `GET /health/readiness` | readinessProbe | 200 if DB ok, 503 if DB down |
| `GET /health` | — | Backward-compatible alias to liveness |

**Readiness check**: Verifies the DB connection with a 2-second timeout. Returns 503 with details if the DB is unreachable.

**Kubernetes file**: `deployments/kubernetes/probes.yaml` is automatically generated with recommended configurations for `livenessProbe`, `readinessProbe` and `startupProbe`.

#### Grafana Dashboard

**7-panel dashboard** pre-configured and auto-provisioned:

| Panel | Type | Description |
|-------|------|-------------|
| Request Rate | Time series | Requests/second |
| Error Rate | Time series | % errors (4xx, 5xx) |
| Latency P95 | Time series | 95th percentile |
| Latency P99 | Time series | 99th percentile |
| Requests in Flight | Gauge | Active requests |
| Status Codes | Pie chart | Code distribution |
| Top Endpoints | Table | Most used endpoints |

**Alerting rules** included:
- High Error Rate (> 5% for 5 min)
- High Latency (p95 > 1s for 5 min)
- Service Down (no metrics for 1 min)

**Access Grafana**:
```bash
open http://localhost:3000
# Credentials: admin / admin
```

#### Docker Compose — Complete Stack

The generated `docker-compose.yml` includes all services:

| Service | Image | Port | URL |
|---------|-------|------|-----|
| app | Build local | 8080 | `http://localhost:8080` |
| db | postgres:16-alpine | 5432 | — |
| jaeger | jaegertracing/all-in-one:1.56.0 | 16686 | `http://localhost:16686` |
| prometheus | prom/prometheus:v2.51.0 | 9090 | `http://localhost:9090` |
| grafana | grafana/grafana:10.4.0 | 3000 | `http://localhost:3000` |

```bash
# Démarrer toute la stack
docker-compose up -d

# Vérifier
curl http://localhost:8080/health/readiness
curl http://localhost:8080/metrics
```

### `basic` Mode — Enhanced Health Checks

```bash
create-go-starter mon-app --observability=basic
```

Generates the `/health/liveness` and `/health/readiness` endpoints without Prometheus, Jaeger or Grafana.

### `none` Mode (default)

```bash
create-go-starter mon-app
# équivalent à:
create-go-starter mon-app --observability=none
```

Standard behavior preserved — no regression.

---

## Available Options

### Flags and Aliases

All flags have short aliases for faster input.

| Flag | Alias | Default | Description |
|------|-------|---------|-------------|
| `--help` | `-h` | | Display help |
| `--interactive` | `-i` | `false` | Launch guided interactive mode |
| `--dry-run` | `-n` | `false` | Preview files without creating them |
| `--template` | `-t` | `full` | Template: `minimal`, `full`, `graphql` |
| `--database` | `-d` | `postgres` | Database: `postgres`, `mysql`, `sqlite` |
| `--observability` | `-o` | `none` | Observability: `none`, `basic`, `advanced` |

**Subcommands:**

| Command | Description |
|---------|-------------|
| `doctor` | Environment diagnostics (Go, Git, Docker) |
| `add-model` | Generate a CRUD model in an existing project |

**Flexible syntax:** Flags accept `=` or space as separator:
```bash
# Ces syntaxes sont équivalentes:
create-go-starter -t=minimal mon-app
create-go-starter -t minimal mon-app
create-go-starter --template=minimal mon-app
create-go-starter --template minimal mon-app
```

### Examples

```bash
# Template minimal
create-go-starter mon-projet --template minimal
create-go-starter mon-projet -t minimal

# Template full (défaut)
create-go-starter mon-projet --template full
create-go-starter mon-projet  # Même résultat

# Template GraphQL
create-go-starter mon-projet --template graphql
create-go-starter mon-projet -t graphql

# Choisir MySQL comme base de données
create-go-starter mon-projet --database=mysql
create-go-starter mon-projet -d mysql

# Choisir SQLite (idéal pour prototypage)
create-go-starter mon-projet --database=sqlite
create-go-starter mon-projet -d sqlite

# Combiner template et database
create-go-starter mon-projet --template=minimal --database=sqlite
create-go-starter mon-projet -t minimal -d sqlite

# Dry-run : prévisualiser sans créer de fichiers
create-go-starter mon-projet --dry-run
create-go-starter mon-projet -n

# Dry-run avec options spécifiques
create-go-starter mon-projet -t minimal -d sqlite -n

# Mode interactif guidé
create-go-starter --interactive
create-go-starter -i

# Combinaison complète avec observabilité
create-go-starter mon-projet -t full -d postgres -o advanced
```

> **Notes**:
> - The `--template` flag is optional. If not specified, the **full** template is used by default.
> - The `--database` flag is optional. If not specified, **PostgreSQL** is used by default.
> - The `--dry-run` flag is optional. It displays the list of files that would be generated without creating them.
> - The `--interactive` flag launches a guided assistant and ignores other configuration flags.

## Interactive Mode (--interactive) <i class="material-icons success">new_releases</i>

**New in v1.4.0!** Interactive mode guides the user step by step to configure a new project.

### Usage

```bash
create-go-starter --interactive
create-go-starter -i
```

### How It Works

Interactive mode prompts successively for:

1. **Project name** -- With real-time validation
2. **Template** -- Choose between minimal, full (default), graphql
3. **Database** -- Choose between postgres (default), mysql, sqlite
4. **Observability** -- Choose between none (default), basic, advanced (if full template)
5. **Summary** -- Displays the chosen configuration with confirmation

```
create-go-starter - Interactive Mode

Enter project name: mon-app
Select template (minimal/full/graphql) [full]: full
Select database (postgres/mysql/sqlite) [postgres]: postgres
Select observability (none/basic/advanced) [none]: advanced

Configuration Summary:
  Project:       mon-app
  Template:      full
  Database:      postgres
  Observability: advanced

Proceed with generation? (y/n) [y]: y
```

### Notes

- <i class="material-icons info">info</i> Interactive mode requires an interactive terminal (no stdin pipe)
- <i class="material-icons warning">warning</i> `--interactive` and `--dry-run` cannot be used together
- <i class="material-icons info">info</i> Zero external dependencies -- uses only `bufio.NewReader` from the stdlib

## Dry-Run Preview (--dry-run) <i class="material-icons success">new_releases</i>

**New in v1.4.0!** Dry-run mode displays the files that would be generated without creating them.

### Usage

```bash
create-go-starter mon-app --dry-run
create-go-starter -n -t minimal -d sqlite mon-app
```

### How It Works

The dry-run displays:
- The configuration used (template, database, observability)
- The complete list of files that would be created
- The number of files and directories
- A warning if the target directory already exists

```
Dry-run mode: no files will be created

Configuration:
  Project:       mon-app
  Template:      minimal
  Database:      sqlite
  Observability: none

Files that would be generated (23 files):
  mon-app/cmd/main.go
  mon-app/internal/models/user.go
  mon-app/internal/domain/user/service.go
  ...

Summary: 23 files in 12 directories
```

### Compatible with All Flags

```bash
create-go-starter -n -t full -d postgres -o advanced mon-app
```

## Doctor Command <i class="material-icons success">new_releases</i>

**New in v1.4.0!** The `doctor` command checks that your environment is correctly configured.

### Usage

```bash
create-go-starter doctor
```

### How It Works

The command checks:

| Tool | Check | Required |
|------|-------|----------|
| **Go** | Version >= 1.21 installed | <i class="material-icons success small">check</i> Yes |
| **Git** | Binary available | <i class="material-icons success small">check</i> Recommended |
| **Docker** | Binary + active daemon | <i class="material-icons info small">info</i> Optional |

### Example Output

```
create-go-starter doctor v1.5.2

Checking environment...

  [OK] Go 1.25.5 (minimum: 1.21)
  [OK] Git 2.43.0
  [OK] Docker 24.0.7 (daemon running)

All checks passed!
```

### Exit Code

- `0` -- Everything is OK
- `1` -- One or more issues detected

Useful for CI/CD scripts:
```bash
create-go-starter doctor && echo "Ready!" || echo "Fix issues first"
```

## Adding Models (add-model) <i class="material-icons success">new_releases</i>

**New in v1.2.0!** Once your project is created, you can automatically generate new complete CRUD models with the `add-model` command.

### Basic Syntax

```bash
create-go-starter add-model <ModelName> --fields "field1:type1,field2:type2,..."
```

### Simple Example

```bash
cd mon-projet  # Naviguer dans un projet go-starter-kit existant

# Ajouter un modèle Todo
create-go-starter add-model Todo --fields "title:string,completed:bool"
```

**This command automatically generates:**
- <i class="material-icons success small">circle</i> `internal/models/todo.go` - Model with GORM tags
- <i class="material-icons success small">circle</i> `internal/interfaces/todo_repository.go` - Repository interface
- <i class="material-icons success small">circle</i> `internal/adapters/repository/todo_repository.go` - GORM implementation
- <i class="material-icons success small">circle</i> `internal/domain/todo/service.go` - Business logic
- <i class="material-icons success small">circle</i> `internal/domain/todo/module.go` - fx module
- <i class="material-icons success small">circle</i> `internal/adapters/handlers/todo_handler.go` - HTTP handlers
- <i class="material-icons success small">circle</i> `internal/domain/todo/service_test.go` - Service tests
- <i class="material-icons success small">circle</i> `internal/adapters/handlers/todo_handler_test.go` - Handler tests

**And automatically modifies:**
- <i class="material-icons warning small">circle</i> `internal/infrastructure/database/database.go` - Adds AutoMigrate
- <i class="material-icons warning small">circle</i> `internal/adapters/http/routes.go` - Adds CRUD routes
- <i class="material-icons warning small">circle</i> `cmd/main.go` - Adds fx module

### Supported Field Types

| Type | Go Type | Example | Description |
|------|---------|---------|-------------|
| `string` | `string` | `"title:string"` | Character string |
| `int` | `int` | `"age:int"` | Signed integer |
| `uint` | `uint` | `"count:uint"` | Unsigned integer |
| `float64` | `float64` | `"price:float64"` | Decimal number |
| `bool` | `bool` | `"active:bool"` | Boolean (true/false) |
| `time` | `time.Time` | `"birthdate:time"` | Date and time |

### GORM Modifiers

Add modifiers after the type to customize database constraints:

| Modifier | GORM Tag | Description |
|----------|----------|-------------|
| `unique` | `gorm:"unique"` | Uniqueness constraint |
| `not_null` | `gorm:"not null"` | Cannot be NULL |
| `index` | `gorm:"index"` | Create a DB index |

**Syntax:**
```
field:type:modifier1:modifier2:...
```

**Examples:**
```bash
# Email unique et obligatoire
create-go-starter add-model User --fields "email:string:unique:not_null"

# Product avec contraintes
create-go-starter add-model Product --fields "name:string:unique:not_null,price:float64,stock:int:index"
```

### Relationships Between Models

`add-model` supports **BelongsTo** (child → parent) and **HasMany** (parent → children) relationships.

#### BelongsTo (child belongs to a parent)

Adds a **foreign key** and a relationship field to the child model.

```bash
# Le modèle parent DOIT exister d'abord
create-go-starter add-model Todo --fields "title:string,completed:bool"

# Créer Comment qui appartient à Todo
create-go-starter add-model Comment --fields "content:string" --belongs-to Todo
```

**Result in `internal/models/comment.go`:**
```go
type Comment struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Content   string         `gorm:"not null" json:"content"`
    TodoID    uint           `gorm:"not null;index" json:"todo_id"`      // Foreign key
    Todo      Todo           `gorm:"foreignKey:TodoID" json:"todo,omitempty"` // Relation
    CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
```

**Generated Endpoints:**
- `GET /api/v1/comments/:id?include=todo` - Retrieve comment with todo preloaded
- `GET /api/v1/todos/:todoId/comments` - List all comments for a todo
- `POST /api/v1/todos/:todoId/comments` - Create a comment for a todo

#### HasMany (parent has multiple children)

Modifies the **parent model** to add a slice of children.

```bash
# Créer Category avec relation has-many vers Product
create-go-starter add-model Category --fields "name:string:unique" --has-many Product
```

**Result in `internal/models/category.go`:**
```go
type Category struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Name      string         `gorm:"unique;not null" json:"name"`
    Products  []Product      `gorm:"foreignKey:CategoryID" json:"products,omitempty"` // Slice ajouté
    CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
```

**Note:** The `Product` model MUST already exist before using `--has-many`.

#### Complete Example with Relationships

```bash
# 1. Créer le modèle parent
create-go-starter add-model Category --fields "name:string:unique"

# 2. Créer le modèle enfant avec belongs-to
create-go-starter add-model Product --fields "name:string,price:float64" --belongs-to Category

# 3. Optionnel: Ajouter has-many au parent après coup
# (pas nécessaire si fait pendant création de Product)

# 4. Utiliser le preloading dans les requêtes
# GET /api/v1/products/1?include=category
# GET /api/v1/categories/1?include=products
```

### Available Flags

| Flag | Description | Required |
|------|-------------|----------|
| `--fields` | Field definitions (`"name:type:mod,..."`) | Yes <i class="material-icons success small">check</i> |
| `--belongs-to <Model>` | Adds BelongsTo relationship (foreign key) | No |
| `--has-many <Model>` | Adds HasMany relationship (children slice) | No |
| `--public` | Public routes (without auth middleware) | No |
| `--yes`, `-y` | Skip confirmation prompt | No |
| `--help`, `-h` | Display help | No |

### Important Notes

<i class="material-icons warning">warning</i> **Pluralization:** The command uses simple pluralization rules (adds an 's'). For irregular plurals (Person→People, Child→Children), manually edit the generated code.

<i class="material-icons warning">warning</i> **Many-to-many relationships:** Not yet supported. Use `--belongs-to` to create relationships via manual join tables.

<i class="material-icons info">info</i> **Swagger:** After generation, run `make swagger` or `swag init -g cmd/api/main.go -o docs/swagger` to update the API documentation.

### Online Help

```bash
create-go-starter add-model --help
```

### Complete Example: Blog System

Here is a complete example of creating a blog system with 3 nested models: **Category** → **Post** → **Comment**.

#### 1. Create the Initial Project

```bash
create-go-starter blog-api
cd blog-api
./setup.sh  # Configurer la base de données
```

#### 2. Add the Category Model (root)

```bash
create-go-starter add-model Category --fields "name:string:unique:not_null,description:string"
```

**Generated Files:**
- `internal/models/category.go` - Entity
- `internal/interfaces/category_repository.go` - Port
- `internal/adapters/repository/category_repository.go` - GORM impl
- `internal/domain/category/service.go` - Business logic
- `internal/domain/category/module.go` - fx module
- `internal/adapters/handlers/category_handler.go` - HTTP handlers
- Tests: `service_test.go`, `handler_test.go`

**Available Endpoints:**
- `POST /api/v1/categories` - Create category
- `GET /api/v1/categories` - List categories
- `GET /api/v1/categories/:id` - Retrieve category
- `PUT /api/v1/categories/:id` - Update category
- `DELETE /api/v1/categories/:id` - Delete category

#### 3. Add the Post Model (child of Category)

```bash
create-go-starter add-model Post --fields "title:string:not_null,content:string,published:bool" --belongs-to Category
```

**What is automatically added:**
- Fields in `internal/models/post.go`:
  ```go
  CategoryID uint     `gorm:"not null;index" json:"category_id"`
  Category   Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
  ```
- Nested routes in `internal/adapters/http/routes.go`:
  ```go
  categories.Get("/:categoryId/posts", postHandler.GetByParent)
  categories.Post("/:categoryId/posts", postHandler.CreateForParent)
  ```

**Available Endpoints:**
- Standard CRUD: `POST/GET/PUT/DELETE /api/v1/posts`
- **Parent relation**: `GET /api/v1/categories/:categoryId/posts` - Posts of a category
- **Parent relation**: `POST /api/v1/categories/:categoryId/posts` - Create post in category
- **Preloading**: `GET /api/v1/posts/:id?include=category` - Post with its category

#### 4. Add the Comment Model (child of Post)

```bash
create-go-starter add-model Comment --fields "author:string:not_null,content:string:not_null" --belongs-to Post
```

**What is automatically added:**
- Fields in `internal/models/comment.go`:
  ```go
  PostID  uint `gorm:"not null;index" json:"post_id"`
  Post    Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
  ```
- Nested routes:
  ```go
  posts.Get("/:postId/comments", commentHandler.GetByParent)
  posts.Post("/:postId/comments", commentHandler.CreateForParent)
  ```

**Available Endpoints:**
- Standard CRUD: `POST/GET/PUT/DELETE /api/v1/comments`
- **Parent relation**: `GET /api/v1/posts/:postId/comments` - Comments of a post
- **Parent relation**: `POST /api/v1/posts/:postId/comments` - Create comment on post
- **Preloading**: `GET /api/v1/comments/:id?include=post` - Comment with its post

#### 5. Add HasMany to the Parent (optional)

If you want to preload children from the parent:

```bash
# Ajouter Posts à Category
create-go-starter add-model Category --has-many Post

# Ajouter Comments à Post
create-go-starter add-model Post --has-many Comment
```

**What is added:**
- In `internal/models/category.go`:
  ```go
  Posts []Post `gorm:"foreignKey:CategoryID" json:"posts,omitempty"`
  ```
- In `internal/models/post.go`:
  ```go
  Comments []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`
  ```

**New preloading endpoints:**
- `GET /api/v1/categories/:id?include=posts` - Category with all its posts
- `GET /api/v1/posts/:id?include=comments` - Post with all its comments
- `GET /api/v1/posts/:id?include=category,comments` - Post with category AND comments

#### 6. Test the Complete System

```bash
# Rebuild le projet
go mod tidy
go build ./...

# Générer Swagger (optionnel)
make swagger

# Lancer tests
make test

# Démarrer le serveur
make run
```

#### 7. API Usage Examples

**Create a category:**
```bash
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Content-Type: application/json" \
  -d '{"name": "Technology", "description": "Tech articles"}'
```

**Create a post in that category:**
```bash
curl -X POST http://localhost:8080/api/v1/categories/1/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Introduction to Go",
    "content": "Go is a statically typed language...",
    "published": true
  }'
```

**Create a comment on that post:**
```bash
curl -X POST http://localhost:8080/api/v1/posts/1/comments \
  -H "Content-Type: application/json" \
  -d '{
    "author": "John Doe",
    "content": "Great article!"
  }'
```

**Retrieve a post with its category and comments:**
```bash
curl http://localhost:8080/api/v1/posts/1?include=category,comments
```

**Expected Response:**
```json
{
  "data": {
    "id": 1,
    "title": "Introduction to Go",
    "content": "Go is a statically typed language...",
    "published": true,
    "category_id": 1,
    "category": {
      "id": 1,
      "name": "Technology",
      "description": "Tech articles"
    },
    "comments": [
      {
        "id": 1,
        "author": "John Doe",
        "content": "Great article!",
        "post_id": 1
      }
    ],
    "created_at": "2026-02-12T10:00:00Z",
    "updated_at": "2026-02-12T10:00:00Z"
  }
}
```

### Conventions and Limitations

#### Pluralization

The `add-model` command uses simple pluralization rules:

**Applied Rules:**
- Adds 's': `Todo` → `todos`, `Product` → `products`
- Replaces 'y' with 'ies': `Category` → `categories`, `Company` → `companies`
- Adds 'es' for s/x/z: `Class` → `classes`, `Box` → `boxes`, `Quiz` → `quizzes`

**Unsupported Irregular Plurals:**
- `Person` → `People` (will be `persons`)
- `Child` → `Children` (will be `childs`)
- `Mouse` → `Mice` (will be `mouses`)

<i class="material-icons info">info</i> **Solution:** Manually edit the generated files to correct irregular plurals:
- Rename files: `internal/domain/persons/` → `internal/domain/people/`
- Update imports and references in the code
- Update routes: `/api/v1/persons` → `/api/v1/people`

#### Supported Relationships

| Relationship | Status | Description |
|--------------|--------|-------------|
| **BelongsTo** (N:1) | <i class="material-icons success">check_circle</i> Supported | Child belongs to a parent |
| **HasMany** (1:N) | <i class="material-icons success">check_circle</i> Supported | Parent has multiple children |
| **Many-to-Many** | <i class="material-icons warning">error</i> Not yet | Planned for a future version |

**Many-to-Many Workaround:**

To create a many-to-many relationship between `User` and `Role`:

```bash
# 1. Créer la table de jointure
create-go-starter add-model UserRole --fields "user_id:uint:index,role_id:uint:index" --belongs-to User --belongs-to Role

# 2. Éditer manuellement pour ajouter contrainte unique composite
# Dans internal/models/user_role.go:
# type UserRole struct {
#     UserID uint `gorm:"uniqueIndex:user_role_unique;not null"`
#     RoleID uint `gorm:"uniqueIndex:user_role_unique;not null"`
#     ...
# }
```

#### Reserved Fields

These fields are **automatically added** and should NOT be specified in `--fields`:

| Field | Type | Description |
|-------|------|-------------|
| `ID` | `uint` | Primary key |
| `CreatedAt` | `time.Time` | Creation date (auto) |
| `UpdatedAt` | `time.Time` | Modification date (auto) |
| `DeletedAt` | `gorm.DeletedAt` | Soft delete (nullable) |

<i class="material-icons warning">warning</i> If you specify these fields, they will be **ignored** by the generator.

#### Custom Validations

Basic validations are generated automatically, but complex validations must be added manually.

**Automatic Validations:**
- `not_null` → GORM Tag: `gorm:"not null"`
- `unique` → GORM Tag: `gorm:"unique"`

**Validations to Add Manually:**

```go
// Dans internal/models/user.go
type User struct {
    Email string `gorm:"unique;not null" json:"email" validate:"required,email"` // Ajouter validate tag
    Age   int    `gorm:"not null" json:"age" validate:"gte=0,lte=150"`          // Validation range
}
```

```go
// Dans internal/domain/user/service.go
func (s *Service) Create(ctx context.Context, user *models.User) error {
    // Ajouter validation custom
    if user.Age < 18 {
        return domain.ErrValidation("user must be at least 18 years old")
    }
    return s.repo.Create(ctx, user)
}
```

#### Swagger Documentation

After generating a new model, **you must regenerate the Swagger documentation**:

```bash
make swagger
# Ou
swag init -g cmd/api/main.go -o docs/swagger
```

<i class="material-icons warning">warning</i> Without regeneration, the new endpoints will not appear in Swagger UI.

**Add Swagger Annotations** (optional):

```go
// Dans internal/adapters/handlers/todo_handler.go

// Create godoc
// @Summary      Create todo
// @Description  Create a new todo item
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        todo body models.Todo true "Todo object"
// @Success      201  {object}  models.Todo
// @Failure      400  {object}  map[string]interface{}
// @Router       /api/v1/todos [post]
func (h *Handler) Create(c *fiber.Ctx) error {
    // ...
}
```

#### Routes and Middleware

**By default**, all generated routes are **protected by JWT** (authentication required).

To create **public routes** (without authentication):

```bash
create-go-starter add-model Article --fields "title:string,content:string" --public
```

This generates routes **without** the `auth.RequireAuth` middleware:

```go
// routes.go - routes publiques
api.Get("/articles", articleHandler.List)
api.Get("/articles/:id", articleHandler.GetByID)
api.Post("/articles", articleHandler.Create)        // Public!
api.Put("/articles/:id", articleHandler.Update)     // Public!
api.Delete("/articles/:id", articleHandler.Delete)  // Public!
```

<i class="material-icons warning">warning</i> **Use --public with caution** to avoid security vulnerabilities.

## Naming Conventions

The project name must follow certain rules:

### Allowed Characters

- **Letters**: a-z, A-Z
- **Digits**: 0-9
- **Hyphens**: - (hyphen)
- **Underscores**: _ (underscore)

### Restrictions

- No spaces
- No special characters (/, \, @, #, etc.)
- No dots (.)
- Must start with a letter or digit (no leading hyphen)

### Valid Examples

```bash
create-go-starter mon-projet           # ✅ Valide
create-go-starter my-awesome-api       # ✅ Valide
create-go-starter user_service         # ✅ Valide
create-go-starter app2024              # ✅ Valide
create-go-starter MonProjet            # ✅ Valide
```

### Invalid Examples

```bash
create-go-starter mon projet           # contient un espace
create-go-starter mon-projet!          # caractère spécial
create-go-starter my.project           # contient un point
create-go-starter -mon-projet          # commence par un tiret
create-go-starter mon/projet           # contient un slash
```

## Generated Structure

Here is the complete structure created by `create-go-starter`:

```
mon-projet/
├── cmd/
│   └── main.go                              # Point d'entrée de l'application
│
├── internal/
│   ├── models/
│   │   └── user.go                          # Entités partagées: User, RefreshToken, AuthResponse
│   │
│   ├── adapters/
│   │   ├── handlers/
│   │   │   ├── auth_handler.go              # Endpoints auth (register, login, refresh)
│   │   │   ├── auth_handler_test.go         # Tests handlers auth
│   │   │   ├── user_handler.go              # Endpoints CRUD users
│   │   │   └── user_handler_test.go         # Tests handlers users
│   │   ├── middleware/
│   │   │   ├── auth_middleware.go           # Middleware JWT authentication
│   │   │   ├── auth_middleware_test.go      # Tests middleware auth
│   │   │   ├── error_handler.go             # Middleware gestion centralisée erreurs
│   │   │   └── error_handler_test.go        # Tests error handler
│   │   ├── repository/
│   │   │   ├── user_repository.go           # Implémentation GORM du repository
│   │   │   └── user_repository_test.go      # Tests repository
│   │   └── http/
│   │       ├── health.go                    # Handler health check
│   │       ├── health_test.go               # Tests health check
│   │       └── routes.go                    # Routes centralisées de l'API
│   │
│   ├── domain/
│   │   ├── errors.go                        # Erreurs métier personnalisées
│   │   ├── errors_test.go                   # Tests erreurs
│   │   └── user/
│   │       ├── service.go                   # Logique métier (Register, Login, etc.)
│   │       ├── service_test.go              # Tests service
│   │       └── module.go                    # Module fx pour dependency injection
│   │
│   ├── infrastructure/
│   │   ├── database/
│   │   │   ├── database.go                  # Configuration GORM et connexion DB
│   │   │   ├── database_test.go             # Tests database
│   │   │   └── module.go                    # Module fx pour database
│   │   └── server/
│   │       ├── server.go                    # Configuration Fiber app
│   │       ├── server_test.go               # Tests server
│   │       └── module.go                    # Module fx pour server
│   │
│   └── interfaces/
│       ├── auth_service.go                  # Interface AuthService (port)
│       ├── user_service.go                  # Interface UserService (port)
│       └── user_repository.go               # Interface UserRepository (port)
│
├── pkg/
│   ├── auth/
│   │   ├── jwt.go                           # Génération et parsing JWT tokens
│   │   ├── jwt_test.go                      # Tests JWT
│   │   ├── middleware.go                    # Middleware JWT pour Fiber
│   │   ├── middleware_test.go               # Tests middleware
│   │   └── module.go                        # Module fx pour auth
│   ├── config/
│   │   ├── env.go                           # Chargement variables d'environnement
│   │   ├── env_test.go                      # Tests config
│   │   └── module.go                        # Module fx pour config
│   └── logger/
│       ├── logger.go                        # Configuration zerolog
│       ├── logger_test.go                   # Tests logger
│       └── module.go                        # Module fx pour logger
│
├── .github/
│   └── workflows/
│       └── ci.yml                           # Pipeline CI/CD GitHub Actions
│
├── .env                                     # Variables d'environnement (créé automatiquement)
├── .env.example                             # Template de configuration
├── .gitignore                               # Exclusions Git
├── .golangci.yml                            # Configuration golangci-lint
├── Dockerfile                               # Build Docker multi-stage
├── Makefile                                 # Commandes utiles (run, test, lint, etc.)
├── setup.sh                                 # Script de configuration automatique
├── go.mod                                   # Module Go et dépendances
└── README.md                                # Documentation du projet
```

**Total**: ~45+ files created automatically!

## Detailed Explanation of Each Directory

### `/cmd`

**Role**: Application entry point.

**Contents**:
- `main.go`: Application bootstrap with uber-go/fx
  - Initializes all modules (database, server, logger, etc.)
  - Configures the lifecycle (OnStart, OnStop)
  - Launches the application with `fx.New().Run()`

**Pattern**: A single minimal file that orchestrates dependencies.

### `/internal/models`

**Role**: Shared domain entities used throughout the entire application.

**Contents**:
- `user.go`: Defines all user-related entities
  - `User`: Main entity with GORM tags (ID, Email, PasswordHash, timestamps)
  - `RefreshToken`: Entity for JWT refresh token management
  - `AuthResponse`: DTO for authentication responses

**Why a Separate Package?**
- **Avoids circular dependencies**: Previously, `interfaces` → `domain/user` → `interfaces` created a cycle
- **Now**: `interfaces` → `models` ← `domain/user` (no cycle!)
- **Centralization**: Entities are defined in a single place
- **Reusability**: All layers (domain, interfaces, adapters) can import models without conflict

**Import**:
```go
import "mon-projet/internal/models"

user := &models.User{
    Email: "user@example.com",
    PasswordHash: hashedPassword,
}
```

### `/internal/domain`

**Role**: Business layer (business logic), the core of hexagonal architecture.

**Contents**:
- `errors.go`: Business error definitions (DomainError, NotFoundError, ValidationError, etc.)
- `user/`: User domain
  - `service.go`: Business logic (Register, Login, GetUserByID, UpdateUser, DeleteUser)
  - `module.go`: fx module that provides the service

**Principle**: The domain must **never** import other packages (adapters, infrastructure). Dependencies are inverted through interfaces (ports). Entities are now in `internal/models` to avoid dependency cycles.

### `/internal/adapters`

**Role**: Adapters that connect the domain to the outside world (HTTP, DB).

#### `/internal/adapters/handlers`

**Role**: HTTP handlers that expose the REST API.

**Contents**:
- `auth_handler.go`:
  - `Register`: POST /api/v1/auth/register
  - `Login`: POST /api/v1/auth/login
  - `RefreshToken`: POST /api/v1/auth/refresh
- `user_handler.go`:
  - `List`: GET /api/v1/users
  - `GetByID`: GET /api/v1/users/:id
  - `Update`: PUT /api/v1/users/:id
  - `Delete`: DELETE /api/v1/users/:id

**Pattern**: Handlers parse requests, validate, call domain services, return responses.

#### `/internal/adapters/middleware`

**Role**: Fiber middleware for cross-cutting concerns.

**Contents**:
- `auth_middleware.go`: Verifies the JWT token in requests
- `error_handler.go`: Centralized error handling (converts DomainError to HTTP responses)

#### `/internal/adapters/repository`

**Role**: Repository pattern implementation with GORM.

**Contents**:
- `user_repository.go`: Implementation of the UserRepository interface
  - Create, FindByID, FindByEmail, Update, Delete
  - Uses GORM for DB operations

**Pattern**: Repository isolates the domain from the persistence layer.

#### `/internal/adapters/http`

**Role**: HTTP routes and utility handlers.

**Contents**:
- `health.go`: GET /health endpoint for monitoring
- `routes.go`: Centralized configuration of all API routes

**Benefits of Route Centralization**:
- Overview of all routes in a single file
- Facilitates API documentation and versioning
- Clear separation between route definitions and handler logic

### `/internal/infrastructure`

**Role**: Infrastructure configuration (DB, server).

#### `/internal/infrastructure/database`

**Role**: Database configuration and connection.

**Contents**:
- `database.go`:
  - PostgreSQL connection via GORM
  - Connection pool configuration
  - AutoMigrate of entities (`models.User`, `models.RefreshToken`)
  - Lifecycle management (connection closing)

#### `/internal/infrastructure/server`

**Role**: Fiber HTTP server configuration.

**Contents**:
- `server.go`:
  - Fiber app configuration
  - Error handler middleware
  - Server lifecycle (start, graceful shutdown)
  - Note: Routes are registered via `server.Module` which invokes `httpRoutes.RegisterRoutes()` with `fx.Invoke`

### `/internal/interfaces`

**Role**: Port definitions (interfaces) for hexagonal architecture.

**Contents**:
- `user_repository.go`: Interface for user persistence

**Principle**:
- Adapters depend on these interfaces, not on concrete implementations
- Interfaces reference `models.*` for return/parameter types
- Example: `CreateUser(ctx context.Context, user *models.User) error`

### `/pkg`

**Role**: Reusable packages, can be imported by other projects.

#### `/pkg/auth`

**Role**: JWT utilities.

**Contents**:
- `jwt.go`: JWT token generation and validation
  - GenerateAccessToken (15min)
  - GenerateRefreshToken (7 days)
  - ParseToken
- `middleware.go`: Fiber middleware for JWT

#### `/pkg/config`

**Role**: Configuration loading.

**Contents**:
- `env.go`: Loads .env variables (godotenv) and exposes them via Config struct

#### `/pkg/logger`

**Role**: Logger configuration.

**Contents**:
- `logger.go`: Configures zerolog (level, format, output)

### `/.github/workflows`

**Role**: CI/CD with GitHub Actions.

**Contents**:
- `ci.yml`: Pipeline that runs:
  1. golangci-lint (quality checks)
  2. Tests with PostgreSQL (service container)
  3. Build verification

## Detailed Configuration Files

### `.env` and `.env.example`

The `.env.example` file is a template, and `.env` is automatically copied during generation.

**Variables**:

```bash
# Application
APP_NAME=mon-projet                # Nom de l'app (utilisé dans logs)
APP_ENV=development                # Environnement (development, staging, production)
APP_PORT=8080                      # Port HTTP

# Database PostgreSQL
DB_HOST=localhost                  # Hôte de la DB
DB_PORT=5432                       # Port PostgreSQL
DB_USER=postgres                   # Utilisateur DB
DB_PASSWORD=postgres               # Mot de passe DB
DB_NAME=mon-projet                 # Nom de la base de données
DB_SSLMODE=disable                 # SSL mode (require pour production)

# JWT Authentication
JWT_SECRET=                        # SECRET CRITIQUE - À générer!
JWT_EXPIRY=15m                     # Durée des access tokens (15 minutes)
REFRESH_TOKEN_EXPIRY=168h          # Durée des refresh tokens (7 jours)
```

**Important**: Generate a secure JWT_SECRET:

```bash
openssl rand -base64 32
```

Then add it to `.env`:

```bash
JWT_SECRET=votre_secret_genere_ici
```

### `go.mod`

Defines the Go module and dependencies:

```go
module mon-projet

go 1.25

require (
    github.com/gofiber/fiber/v2 v2.x.x
    github.com/golang-jwt/jwt/v5 v5.x.x
    gorm.io/gorm v1.x.x
    gorm.io/driver/postgres v1.x.x
    go.uber.org/fx v1.x.x
    github.com/rs/zerolog v1.x.x
    github.com/go-playground/validator/v10 v10.x.x
    golang.org/x/crypto v0.x.x
    github.com/joho/godotenv v1.x.x
    // ... autres dépendances
)
```

### `Makefile`

Useful commands for development:

```makefile
.PHONY: help run build test test-coverage lint clean docker-build docker-run

help:           # Afficher toutes les commandes
run:            # Lancer l'application (go run)
build:          # Compiler le binaire
test:           # Exécuter les tests avec race detector
test-coverage:  # Tests avec rapport HTML de coverage
lint:           # Linter avec golangci-lint
clean:          # Nettoyer les artifacts
docker-build:   # Build l'image Docker
docker-run:     # Lancer le conteneur Docker
```

### `.golangci.yml`

Linter configuration with strict rules:

- errcheck, gosimple, govet, ineffassign, staticcheck
- gofmt, goimports
- misspell, revive
- Exclusions for tests and generated code

### `Dockerfile`

Optimized multi-stage build:

**Stage 1 (builder)**:
- golang:1.25-alpine image
- Static binary build

**Stage 2 (runtime)**:
- alpine:latest image (lightweight)
- Copies binary only
- Runs as non-root
- EXPOSE 8080

**Final size**: ~15-20MB (vs ~1GB with full golang image)

## Post-Generation Workflow

Once the project is created, you have two options to configure your project:

### Option A: Automatic Configuration with setup.sh (Recommended) <i class="material-icons success">rocket_launch</i>

The CLI automatically generates a `setup.sh` script that automates the entire initial configuration.

**Script Features**:
- <i class="material-icons success small">check_circle</i> Prerequisites check (Go, OpenSSL, Docker)
- <i class="material-icons success small">check_circle</i> Go dependencies installation (`go mod tidy`)
- <i class="material-icons success small">check_circle</i> Automatic JWT secret generation
- <i class="material-icons success small">check_circle</i> PostgreSQL configuration (Docker or local)
- <i class="material-icons success small">check_circle</i> Test execution
- <i class="material-icons success small">check_circle</i> Installation verification

**Usage**:

```bash
cd mon-projet
./setup.sh
```

The script is **interactive** and will guide you through the choices:
- Docker or local PostgreSQL
- JWT secret regeneration if already configured
- Validation at each step

**After running the script**:

```bash
make run
```

That's it! Your application is ready.

---

### Option B: Manual Configuration

If you prefer to configure manually or if the setup.sh script fails, follow these steps:

### Step 1: Navigate to the Project

```bash
cd mon-projet
```

### Step 2: Configure Environment Variables

```bash
# Générer un JWT secret
openssl rand -base64 32

# Éditer .env et ajouter le secret
nano .env  # ou vim, code, etc.
```

Add:
```
JWT_SECRET=le_secret_genere
```

### Step 3: Install Dependencies

```bash
go mod tidy
```

### Step 4: Start PostgreSQL

**Option A: Docker (recommended)**

```bash
docker run -d \
  --name postgres \
  -e POSTGRES_DB=mon-projet \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine
```

**Option B: Local Installation**

```bash
# macOS
brew install postgresql
brew services start postgresql
createdb mon-projet

# Linux
sudo apt install postgresql
sudo systemctl start postgresql
sudo -u postgres createdb mon-projet
```

### Step 5: Launch the Application

```bash
make run
```

Or directly:

```bash
go run cmd/main.go
```

You should see:

```
{"level":"info","time":"...","message":"Starting server on :8080"}
{"level":"info","time":"...","message":"Database connected successfully"}
```

### Step 6: Test the API

```bash
# Health check
curl http://localhost:8080/health
# {"status":"ok"}

# Register un utilisateur
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

## Available Make Commands

The generated Makefile includes these commands:

| Command | Description | Usage |
|---------|-------------|-------|
| `make help` | Display help | See all commands |
| `make run` | Launch the app | Local development |
| `make build` | Compile the binary | Create the executable |
| `make test` | Tests with race detector | Verify the code |
| `make test-coverage` | Tests + HTML report | View coverage |
| `make lint` | golangci-lint | Check quality |
| `make clean` | Clean up | Remove artifacts |
| `make docker-build` | Build Docker image | Containerization |
| `make docker-run` | Launch container | Docker test |

**Examples**:

```bash
# Développement quotidien
make run          # Lance l'app avec hot-reload (si air installé)

# Avant commit
make test         # Vérifie que tous les tests passent
make lint         # Vérifie la qualité du code

# Build pour production
make build        # Crée le binaire
./mon-projet      # Exécute le binaire

# Docker
make docker-build # Build l'image
make docker-run   # Teste le conteneur
```

## Next Steps

Now that you understand the structure, check out:

1. **[Generated Projects Guide](./generated-project-guide.md)** - Complete guide for:
   - Understanding hexagonal architecture
   - Developing new features
   - Using the API (endpoints, authentication)
   - Writing tests
   - Deploying to production

2. **[CLI Architecture](./cli-architecture.md)** - If you want to:
   - Understand how the generator works
   - Contribute to the project
   - Extend the templates

3. **Start developing**:
   ```bash
   # Lancer l'app
   make run

   # Dans un autre terminal, tester l'API
   curl http://localhost:8080/health

   # Lire le code généré
   cat internal/domain/user/service.go
   cat internal/adapters/handlers/auth_handler.go
   ```

## Tips

- **Read the generated code**: Each file is an example of Go best practices
- **Modify according to your needs**: The structure is a starting point, adapt it
- **Follow the patterns**: Repository, dependency injection, centralized error handling
- **Test regularly**: `make test` before each commit
- **Use the linter**: `make lint` to maintain quality

Happy coding! <i class="material-icons success">rocket_launch</i>
