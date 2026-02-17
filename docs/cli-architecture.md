# Architecture de l'outil CLI

Documentation technique pour contributeurs et développeurs avancés.

## Vue d'ensemble

`create-go-starter` est un générateur de projets Go qui crée des applications complètes avec architecture hexagonale, authentification JWT, API REST, et infrastructure de déploiement.

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

**Statistiques**:
- Lignes de code: ~6,000+
- Fichiers générés par projet: 46+ (60+ avec observability)
- Templates: 45+ fonctions
- Dépendances: Standard library uniquement

## Composants principaux

### 1. main.go - Point d'entrée

**Responsabilités**:
- Parsing des arguments CLI (package `flag`)
- Validation du nom de projet (regex alphanumeric + hyphens/underscores)
- Création de la structure de répertoires
- Orchestration de la génération de fichiers
- Gestion des erreurs et affichage console (avec couleurs)

**Fonctions clés**:

```go
func main()
func validateProjectName(name string) error
func createProjectStructure(projectName string) error
func copyEnvFile(projectPath string) error
func printSuccess(projectName string)
```

**Flux d'exécution**:

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

**Validation du nom**:

```go
// Regex pattern
var validProjectName = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)

// Valid: mon-projet, my-api, user_service, app2024
// Invalid: -projet, mon projet, my.project, @app
```

**Messages colorés**:

```go
func Green(text string) string
func Red(text string) string

// Usage
fmt.Println(Green("✓ Project created successfully!"))
fmt.Println(Red("✗ Error: " + err.Error()))
```

### 2. generator.go - Orchestrateur de génération

**Responsabilités**:
- Validation du répertoire projet
- Validation du module name Go (compatible go.mod)
- Création de tous les fichiers du projet
- Gestion des templates avec injection du nom de projet

**Structure FileGenerator**:

```go
type FileGenerator struct {
    Path    string  // Chemin relatif au projet (ex: cmd/main.go)
    Content string  // Contenu généré du fichier
}
```

**Fonction principale**:

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

### 2.1. Générateur de Modèles (add-model Command)

**Nouveau dans v1.2.0!** Le générateur CRUD scaffolding permet d'ajouter automatiquement de nouveaux modèles à un projet existant.

**Responsabilités**:
- Génération automatique de modèles CRUD complets
- Mise à jour des fichiers existants (database, routes, main)
- Support des relations (BelongsTo, HasMany)
- Génération de tests unitaires

**Structure des fichiers**:

```
cmd/create-go-starter/
├── add_model.go           # CLI orchestration add-model
├── model.go               # Structures: Field, Model, Relation
└── model_generator.go     # 13+ générateurs de templates
```

**Structures de données (model.go)**:

```go
// Field représente un champ du modèle
type Field struct {
    Name      string   // Nom du champ (ex: "Title")
    Type      string   // Type Go (ex: "string", "int", "time.Time")
    GormTags  []string // Tags GORM (ex: ["unique", "not null"])
    JSONTag   string   // Tag JSON (ex: "title")
}

// Model représente le modèle complet
type Model struct {
    Name       string   // Nom du modèle (ex: "Todo")
    NamePlural string   // Forme plurielle (ex: "Todos")
    Fields     []Field  // Champs custom
    Relations  []Relation // Relations avec autres modèles
}

// Relation représente une relation entre modèles
type Relation struct {
    Type        string // "belongs_to" ou "has_many"
    ModelName   string // Nom du modèle lié (ex: "Post")
    ForeignKey  string // Clé étrangère (ex: "PostID")
}
```

**Commande CLI (add_model.go)**:

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
    // - Vérifie que nous sommes dans un projet go-starter-kit
    // - Vérifie que le modèle n'existe pas déjà
    // - Parse et valide les champs
    // - Vérifie que les modèles parents existent (pour relations)

    // Génération
    model := BuildModel(modelName, fields, relations)
    generator := NewModelGenerator(model, projectPath)

    // Génère les 8 fichiers
    generator.GenerateAll()

    // Met à jour les 3 fichiers existants
    generator.UpdateExistingFiles()

    return nil
}
```

**Générateurs de templates (model_generator.go)**:

```go
type ModelGenerator struct {
    model       Model
    projectPath string
    projectName string
}

// 13+ méthodes de génération

// 1-8: Fichiers neufs
func (g *ModelGenerator) GenerateModel() string
func (g *ModelGenerator) GenerateRepositoryInterface() string
func (g *ModelGenerator) GenerateRepositoryImpl() string
func (g *ModelGenerator) GenerateService() string
func (g *ModelGenerator) GenerateServiceModule() string
func (g *ModelGenerator) GenerateHandler() string
func (g *ModelGenerator) GenerateServiceTests() string
func (g *ModelGenerator) GenerateHandlerTests() string

// 9-13: Mise à jour de fichiers existants
func (g *ModelGenerator) UpdateDatabaseMigrations() error
func (g *ModelGenerator) UpdateRoutes() error
func (g *ModelGenerator) UpdateMainModule() error
func (g *ModelGenerator) UpdateParentModel() error
func (g *ModelGenerator) AddBelongsToImport() error
```

**Workflow de génération**:

```
add-model Todo --fields "title:string,completed:bool"
    │
    ├─> 1. Parse et valide arguments
    │   ├─ ModelName: "Todo" → "todos" (pluralisation)
    │   ├─ Fields: [{Name:"Title", Type:"string"}, {Name:"Completed", Type:"bool"}]
    │   └─ Validate: Pas de conflit, types valides
    │
    ├─> 2. Génère 8 fichiers neufs
    │   ├─ internal/models/todo.go (GORM entity)
    │   ├─ internal/interfaces/todo_repository.go (port)
    │   ├─ internal/adapters/repository/todo_repository.go (GORM impl)
    │   ├─ internal/domain/todo/service.go (business logic)
    │   ├─ internal/domain/todo/module.go (fx module)
    │   ├─ internal/adapters/handlers/todo_handler.go (HTTP)
    │   ├─ internal/domain/todo/service_test.go (tests)
    │   └─ internal/adapters/handlers/todo_handler_test.go (tests)
    │
    └─> 3. Met à jour 3 fichiers existants
        ├─ internal/infrastructure/database/database.go (AutoMigrate)
        ├─ internal/adapters/http/routes.go (Routes CRUD)
        └─ cmd/main.go (fx module registration)
```

**Exemple de génération (model)**:

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

**Relations supportées**:

1. **BelongsTo (enfant → parent)**:
   ```bash
   create-go-starter add-model Comment --fields "content:string" --belongs-to Post
   ```
   - Ajoute `PostID uint` + `Post Post` dans Comment
   - Endpoints: `GET /posts/:postId/comments`, `POST /posts/:postId/comments`

2. **HasMany (parent → enfants)**:
   ```bash
   create-go-starter add-model Category --fields "name:string" --has-many Product
   ```
   - Ajoute `Products []Product` dans Category
   - Preloading: `GET /categories/:id?include=products`

**Types supportés**:

| Type CLI | Go Type | GORM Tag | Description |
|----------|---------|----------|-------------|
| `string` | `string` | `gorm:"not null"` | Chaîne de caractères |
| `int` | `int` | `gorm:"not null"` | Entier signé |
| `uint` | `uint` | `gorm:"not null"` | Entier non-signé |
| `float64` | `float64` | `gorm:"not null"` | Nombre décimal |
| `bool` | `bool` | `gorm:"not null"` | Booléen |
| `time` | `time.Time` | `gorm:"not null"` | Date/heure |

**Modificateurs GORM**:

| Modificateur | Tag GORM | Description |
|--------------|----------|-------------|
| `unique` | `gorm:"unique"` | Contrainte d'unicité |
| `not_null` | `gorm:"not null"` | Non nullable |
| `index` | `gorm:"index"` | Index DB |

**Pluralisation**:

Le générateur utilise des règles simples de pluralisation:
- Ajoute 's': Todo → Todos, Product → Products
- Gère 'y': Category → Categories
- Gère 's/x/z': Class → Classes, Box → Boxes

Pour pluriels irréguliers (Person→People), éditer manuellement après génération.

**Validation et sécurité**:

```go
// Validation stricte
- Modèle parent existe (pour --belongs-to)
- Modèle enfant existe (pour --has-many)
- Pas de conflit de noms de fichiers
- Types de champs valides
- Pas de champs réservés (ID, CreatedAt, etc.)

// Sécurité
- Confirmation prompt (sauf --yes)
- Dry-run preview des fichiers générés
- Backup automatique des fichiers modifiés
```

**Intégration tests**:

Chaque modèle généré inclut:
- Tests service: CRUD operations
- Tests handler: HTTP endpoints
- Mock repository (interface)
- Test fixtures

**Liste complète des fichiers générés** (45+ fichiers):

1. **Configuration** (6 fichiers):
   - go.mod, .env.example, .gitignore, .golangci.yml
   - Makefile, docker-compose.yml

2. **Bootstrap** (1 fichier):
   - cmd/main.go

3. **Packages réutilisables** (7 fichiers):
   - pkg/config/env.go, pkg/config/module.go
   - pkg/logger/logger.go, pkg/logger/module.go
   - pkg/auth/jwt.go, pkg/auth/middleware.go, pkg/auth/module.go

4. **Models** (1 fichier):
   - internal/models/user.go (User, RefreshToken, AuthResponse)

5. **Domain** (3 fichiers):
   - internal/domain/errors.go
   - internal/domain/user/service.go
   - internal/domain/user/module.go

6. **Interfaces** (1 fichier):
   - internal/interfaces/user_repository.go

7. **Adapters** (10 fichiers):
   - Handlers: auth_handler.go, user_handler.go, module.go
   - Middleware: error_handler.go
   - Repository: user_repository.go, module.go
   - HTTP: health.go, routes.go

8. **Infrastructure** (2 fichiers):
   - internal/infrastructure/database/database.go
   - internal/infrastructure/server/server.go

9. **Deployment** (3 fichiers):
   - Dockerfile
   - .github/workflows/ci.yml
   - docker-compose.yml

10. **Documentation** (3 fichiers):
   - README.md
   - docs/README.md
   - docs/quick-start.md

**Validation Go module name**:

```go
func validateGoModuleName(name string) error {
    // Must start with letter/number
    // Only alphanumeric, hyphens, underscores
    // No spaces, special chars
    pattern := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)
    return pattern.MatchString(name)
}
```

### 3. templates.go - Templates principaux

**Responsabilités**:
- Définition des templates pour infrastructure et configuration
- Templates: configuration, server, database, middleware, Docker, CI/CD
- Injection dynamique du nom de projet

**Structure ProjectTemplates**:

```go
type ProjectTemplates struct {
    projectName string
}

func NewProjectTemplates(projectName string) *ProjectTemplates {
    return &ProjectTemplates{projectName: projectName}
}
```

**Méthodes principales** (30+ templates):

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
func (t *ProjectTemplates) ServerTemplate() string          // server.go avec Fiber
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

**SetupScriptTemplate** génère un script bash interactif qui:
- Vérifie les prérequis (Go, OpenSSL, Docker, psql)
- Installe les dépendances (`go mod tidy`)
- Génère et configure le JWT secret dans `.env`
- Configure PostgreSQL (Docker ou local)
- Exécute les tests de validation
- Vérifie l'installation complète

Le script est rendu exécutable automatiquement (`chmod 0755`) par `generator.go`.

**Pattern de template**:

Les templates utilisent la concaténation de strings (pas de text/template):

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

**Avantages**:
- Simplicité (pas de parsing)
- Type-safe à la compilation
- Facile à débugger
- Injection directe du projectName

**Inconvénients**:
- Verbeux pour templates complexes
- Échappement manuel des backticks

### 4. templates_user.go - Templates du domaine User

**Responsabilités**:
- Templates spécifiques au domaine User
- Entités (models), services, repositories, handlers
- Tests (si implémenté)

**Méthodes**:

#### Models (Entités partagées)

```go
func (t *ProjectTemplates) ModelsUserTemplate() string  // User, RefreshToken, AuthResponse
```

#### Domain

```go
func (t *ProjectTemplates) DomainErrorsTemplate() string
func (t *ProjectTemplates) UserServiceTemplate() string   // Business logic
func (t *ProjectTemplates) UserModuleTemplate() string    // fx module
```

**Note**: `UserEntityTemplate()` et `UserRefreshTokenTemplate()` sont dépréciés - les entités sont maintenant dans `ModelsUserTemplate()`.

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

**Contenu typique d'un template** (exemple: ModelsUserTemplate):

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

**Pourquoi `models` au lieu de `domain/user`?**
- **Évite les dépendances circulaires**: Les interfaces peuvent référencer les modèles sans créer de cycles
- **Centralisation**: Les entités sont définies en un seul endroit
- **Réutilisabilité**: Tous les layers (domain, interfaces, adapters) peuvent importer models sans conflit

### 5. templates_observability.go - Templates d'Observabilité

**Nouveau dans v1.3.0!** Ce fichier contient tous les templates pour la stack d'observabilité avancée (~1369 lignes).

**Responsabilités**:
- Templates pour Prometheus metrics (endpoint `/metrics`, middleware, registry)
- Templates pour OpenTelemetry distributed tracing (Jaeger export, OTLP/gRPC)
- Templates pour health checks avancés K8s (liveness/readiness)
- Templates pour Grafana dashboard (JSON, provisioning, alerting)
- Templates pour Docker Compose observabilité (Jaeger, Prometheus, Grafana)
- Templates pour configuration Kubernetes probes

**Condition d'activation**: Uniquement généré quand `--observability=advanced` ET `--template=full`.

**Méthodes principales** (15+ templates):

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

#### Health Checks Avancés

```go
func (t *ProjectTemplates) AdvancedHealthHandlerTemplate() string // http/health.go (remplace le basique)
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

#### Docker Compose Observabilité

```go
func (t *ProjectTemplates) DockerComposeObservabilityTemplate() string // docker-compose.yml (version étendue)
```

**Fichiers générés par l'observabilité avancée** (14+ fichiers supplémentaires):

| Catégorie | Fichiers | Description |
|-----------|----------|-------------|
| Metrics | 3 | prometheus.go, metrics_middleware.go, metrics_handler.go |
| Tracing | 4 | tracer.go, tracing_middleware.go, gorm_tracing.go, logger_tracing.go |
| Health | 1 | health.go (version avancée avec liveness/readiness) |
| Grafana | 3 | dashboard JSON, datasource YAML, provisioning YAML |
| Prometheus | 2 | prometheus.yml, alert_rules.yml |
| Kubernetes | 1 | probes.yaml |

**Dépendances ajoutées au go.mod** (uniquement quand observability=advanced):

```
github.com/ansrivas/fiberprometheus/v2 v2.7.0
go.opentelemetry.io/otel v1.24.0
go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.24.0
go.opentelemetry.io/otel/sdk v1.24.0
go.opentelemetry.io/otel/trace v1.24.0
go.opentelemetry.io/contrib/instrumentation/github.com/gofiber/fiber/otelfiber v0.49.0
```

### 6. git.go - Initialisation Git

**Responsabilités**:
- Vérification de la disponibilité de Git sur le système
- Initialisation automatique d'un dépôt Git dans le projet généré
- Création d'un commit initial avec tous les fichiers générés

**Fonctions clés**:

```go
func isGitAvailable() bool           // Vérifie si git est installé
func initGitRepo(projectPath string) error  // Initialise le repo et crée le commit initial
```

**Comportement**:
- Si Git est disponible: initialise le repo et crée un commit "Initial commit from go-starter-kit"
- Si Git n'est pas disponible: affiche un avertissement mais continue (dégradation gracieuse)
- Le `.gitignore` est ajouté automatiquement avant le commit initial

**Intégration**:
- Appelé dans `main.go` après `copyEnvFile()` et avant `printSuccessMessage()`
- Messages: ":material-wrench: Setting up Git repository..." et ":material-check-circle: Git repository initialized"

## Patterns et conventions

### 1. Pattern de templates

**Choix: String concatenation vs text/template**

**Option choisie**: String concatenation

```go
return `package main

import "fmt"

func main() {
    fmt.Println("` + t.projectName + `")
}
`
```

**Pourquoi pas text/template?**
- Simplicité: Pas de parsing, pas de struct de données
- Performance: Pas d'overhead de parsing
- Type-safety: Erreurs à la compilation
- Debugging: Plus facile de voir le template généré

**Challenges**:
- Échappement des backticks: Utiliser `` "`" ``
- Templates longs deviennent verbeux

### 2. Validation en couches

**Layer 1 - CLI level (main.go)**:
- Validation du nom de projet
- Regex: `^[a-zA-Z0-9][a-zA-Z0-9_-]*$`
- Exemples valides: `mon-projet`, `my_app`, `app2024`

**Layer 2 - Generator level (generator.go)**:
- Validation module name Go (même règles que Layer 1)
- Vérification que le répertoire existe

**Layer 3 - Runtime (code généré)**:
- Validation avec go-playground/validator
- Validation métier dans le domain

### 3. Gestion des erreurs

**Convention**:

```go
if err != nil {
    return fmt.Errorf("context: %w", err)  // Wrap error with context
}
```

**Affichage utilisateur**:

```go
fmt.Println(Red("✗ Error: " + err.Error()))
os.Exit(1)
```

**Pas de panic**: Utiliser return error, pas panic()

### 4. Tests

**Organisation**:

```
cmd/create-go-starter/
├── main.go
├── main_test.go           # Tests CLI et fonction run()
├── generator.go
├── generator_test.go      # Tests génération
├── templates.go
├── templates_test.go      # Tests templates
├── templates_user.go
├── git.go
├── git_test.go            # Tests initialisation Git
├── colors_test.go         # Tests utilitaires couleurs
├── env_test.go            # Tests .env copy
├── scaffold_test.go       # Tests création structure
└── smoke_test.go          # Tests E2E smoke tests
scripts/
└── smoke_test.sh          # Script bash E2E validation
```

**Couverture de tests**: 83%+

**Commandes Makefile**:
```bash
make test              # Tous les tests
make test-short        # Tests rapides (skip tests longs)
make smoke-test        # Validation E2E complète avec runtime
make smoke-test-quick  # Validation E2E sans runtime (pas de Docker)
```

**Patterns de tests**:

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

## Extensibilité

### Ajouter un nouveau template

**Étapes**:

1. **Créer la méthode template** dans `templates.go` ou `templates_user.go`:

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

2. **Ajouter FileGenerator** dans `generator.go`:

```go
files := []FileGenerator{
    // ... existing files
    {
        Path:    filepath.Join(projectPath, "pkg", "mypackage", "hello.go"),
        Content: templates.MyNewTemplate(),
    },
}
```

3. **Tester la génération**:

```bash
go run ./cmd/create-go-starter test-project
cat test-project/pkg/mypackage/hello.go
```

4. **Ajouter un test** dans `templates_test.go`:

```go
func TestMyNewTemplate(t *testing.T) {
    templates := NewProjectTemplates("test")
    content := templates.MyNewTemplate()

    assert.Contains(t, content, "package mypackage")
    assert.Contains(t, content, "test")
}
```

### Ajouter une option CLI

**Exemple: Ajouter `--database` flag pour choisir la DB**

1. **Définir le flag** dans `main.go`:

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

2. **Passer à generateProjectFiles**:

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

3. **Adapter les templates**:

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

### Modèles d'extension futurs

**1. Templates multiples**:

```bash
create-go-starter my-project --template=minimal
create-go-starter my-project --template=full
create-go-starter my-project --template=api-only
create-go-starter my-project --template=graphql
```

**Implémentation**:

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

**2. Choix de DB**:

```bash
create-go-starter my-project --db=mysql
create-go-starter my-project --db=mongodb
```

**3. Choix de framework**:

```bash
create-go-starter my-project --framework=gin
create-go-starter my-project --framework=echo
```

**4. CLI interactif**:

```bash
create-go-starter

? Project name: my-awesome-api
? Database: PostgreSQL
? Auth: JWT
? Generate Swagger docs: Yes
? Include Docker: Yes

✓ Generating project...
```

Utiliser: `github.com/manifoldco/promptui`

## Dépendances

**Dépendances du CLI**: **AUCUNE** (seulement stdlib)

**Avantages**:
- Simplicité: Pas de go mod tidy pour le CLI
- Portabilité: Binaire statique sans dépendances
- Installation légère: `go install` très rapide
- Maintenance facile: Pas de breaking changes externes

**Packages stdlib utilisés**:
- `flag` - Parsing CLI
- `fmt` - Formatting et printing
- `os` - File operations, exit codes
- `path/filepath` - Path manipulation
- `regexp` - Validation patterns
- `strings` - String utilities

## Performance

**Métriques**:

- **Temps de génération**: < 1 seconde pour 45+ fichiers
- **Taille binaire**: ~3-4 MB (statique)
- **Mémoire**: < 10 MB pendant génération
- **Disk writes**: 45+ fichiers, ~15,000 lignes de code générées

**Optimisations**:

1. **Pas de parsing de templates**: String concatenation directe
2. **Batch file creation**: Tous les fichiers créés en une passe
3. **MkdirAll une fois**: Créé tous les répertoires parents si nécessaire
4. **Pas de dépendances**: Pas de download ni d'import overhead

## Standards de code

**Conventions suivies**:

1. **gofmt**: Toujours formaté
   ```bash
   go fmt ./...
   ```

2. **golangci-lint**: Quality checks
   ```bash
   golangci-lint run ./...
   ```

3. **Tests coverage**: > 80%
   ```bash
   go test -cover ./...
   ```

4. **Documentation GoDoc**: Pour exports publics
   ```go
   // ProjectTemplates holds all template generation methods
   type ProjectTemplates struct {
       projectName string
   }
   ```

5. **Error handling**: Toujours explicite, jamais ignore
   ```go
   if err != nil {
       return fmt.Errorf("context: %w", err)
   }
   ```

## Débogage

**Activer mode verbose** (à implémenter):

```go
var verbose bool
flag.BoolVar(&verbose, "verbose", false, "Verbose output")

if verbose {
    log.Println("Creating directory:", path)
    log.Println("Writing file:", filepath)
}
```

**Tester génération**:

```bash
# Générer projet test
go run ./cmd/create-go-starter test-project

# Vérifier fichiers
ls -la test-project/
tree test-project/

# Vérifier contenu
cat test-project/go.mod
cat test-project/cmd/main.go

# Test build
cd test-project
go mod tidy
go build ./...

# Nettoyer
cd ..
rm -rf test-project
```

**Debugger avec Delve**:

```bash
dlv debug ./cmd/create-go-starter -- my-project
(dlv) break main.createProjectStructure
(dlv) continue
```

## Contribution

Pour contribuer au CLI:

1. Fork le repository
2. Créer une branche: `git checkout -b feature/my-feature`
3. Faire les changements
4. Tests: `go test ./...`
5. Lint: `golangci-lint run`
6. Commit: `git commit -m "feat: add feature"`
7. Push: `git push origin feature/my-feature`
8. Ouvrir une Pull Request

**Checklist PR**:
- [ ] Tests ajoutés/mis à jour
- [ ] Tests passent (`make test`)
- [ ] Lint passe (`make lint`)
- [ ] Documentation mise à jour
- [ ] Commit messages clairs (conventional commits)

## Roadmap technique

**Court terme**:
- [ ] Version flag (`--version`)
- [ ] Verbose mode (`--verbose`)
- [ ] Dry-run mode (`--dry-run`)
- [ ] Force overwrite (`--force`)

**Moyen terme**:
- [ ] Templates multiples (minimal, full, api-only)
- [ ] Choix de DB (PostgreSQL, MySQL, SQLite, MongoDB)
- [ ] Choix de framework (Fiber, Gin, Echo)
- [ ] CLI interactif (prompts)

**Long terme**:
- [ ] Plugin system pour templates custom
- [ ] Template marketplace
- [ ] Hot-reload des templates
- [ ] GUI pour génération

---

Cette documentation technique devrait permettre aux contributeurs de comprendre l'architecture interne du CLI et de contribuer efficacement au projet.
