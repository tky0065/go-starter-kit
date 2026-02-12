# Story 8.2: Model & Repository Generation

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

**En tant que** développeur,
**Je veux** que le modèle GORM et le repository soient générés automatiquement après `add-model`,
**Afin de** ne pas écrire le code boilerplate de persistance.

## Acceptance Criteria

1. **AC1**: Given la commande `add-model Todo --fields "title:string,completed:bool"` est validée, When la génération s'exécute, Then `internal/models/todo.go` est créé avec la struct Go, les tags GORM et JSON corrects
2. **AC2**: Given le modèle est généré, When le code est compilé, Then `internal/interfaces/todo_repository.go` définit l'interface CRUD complète
3. **AC3**: Given l'interface existe, When la génération continue, Then `internal/adapters/repository/todo_repository.go` implémente le CRUD avec GORM
4. **AC4**: Given tous les fichiers sont générés, When `go build ./...` est exécuté dans le projet, Then la compilation réussit sans erreur
5. **AC5**: Given le modèle est généré, When `database.go` AutoMigrate est vérifié, Then le nouveau modèle est ajouté à la liste de migration

## Tasks / Subtasks

- [x] Task 1: Générer le fichier modèle (AC: 1)
  - [x] 1.1 Créer template pour `internal/models/<name>.go`
  - [x] 1.2 Générer struct avec champs ID, CreatedAt, UpdatedAt, DeletedAt (base GORM)
  - [x] 1.3 Ajouter champs utilisateur avec tags `gorm:""` et `json:""` corrects
  - [x] 1.4 Gérer les imports nécessaires (`time`, `gorm.io/gorm`)

- [x] Task 2: Générer l'interface repository (AC: 2)
  - [x] 2.1 Créer template pour `internal/interfaces/<name>_repository.go`
  - [x] 2.2 Définir méthodes CRUD: Create, FindByID, FindAll (paginé), Update, Delete
  - [x] 2.3 Utiliser les bons types (context.Context, pointeurs de modèle)

- [x] Task 3: Générer l'implémentation repository (AC: 3)
  - [x] 3.1 Créer template pour `internal/adapters/repository/<name>_repository.go`
  - [x] 3.2 Implémenter CRUD complet avec GORM et contexte
  - [x] 3.3 Pagination avec offset/limit pattern existant
  - [x] 3.4 Soft delete via GORM (automatique avec DeletedAt)
  - [x] 3.5 Gestion d'erreurs: `gorm.ErrRecordNotFound` → `nil, nil`

- [x] Task 4: Mettre à jour les modules fx (AC: 4)
  - [x] 4.1 Ajouter le repository au module fx dans `internal/adapters/repository/module.go`
  - [x] 4.2 S'assurer que l'interface est satisfaite (provide interface, not concrete)

- [x] Task 5: Mettre à jour AutoMigrate (AC: 5)
  - [x] 5.1 Ajouter `&models.Todo{}` dans l'appel `db.AutoMigrate()` de `internal/infrastructure/database/database.go`

- [x] Task 6: Tests unitaires (AC: 1-5)
  - [x] 6.1 Test de génération du modèle (struct correcte, tags)
  - [x] 6.2 Test de génération de l'interface (méthodes CRUD)
  - [x] 6.3 Test de génération du repository (implémente interface)
  - [x] 6.4 Test E2E: génération complète + `go build ./...`

## Dev Notes

### Pattern du Modèle User Existant (Référence)

Le modèle User dans `templates_user.go` est la référence pour le pattern :

```go
type User struct {
    ID           uint           `gorm:"primaryKey" json:"id"`
    Email        string         `gorm:"uniqueIndex;not null" json:"email"`
    PasswordHash string         `gorm:"not null" json:"-"`
    CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
```

**Champs de base OBLIGATOIRES pour tout modèle généré :**
- `ID uint` avec `gorm:"primaryKey" json:"id"`
- `CreatedAt time.Time` avec `gorm:"autoCreateTime" json:"created_at"`
- `UpdatedAt time.Time` avec `gorm:"autoUpdateTime" json:"updated_at"`
- `DeletedAt gorm.DeletedAt` avec `gorm:"index" json:"deleted_at,omitempty"`

### Template Modèle Généré (Exemple : Todo)

```go
package models

import (
    "time"
    "gorm.io/gorm"
)

// Todo represents a todo item in the system.
type Todo struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Title     string         `gorm:"not null" json:"title"`
    Completed bool           `gorm:"default:false" json:"completed"`
    CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
```

### Template Interface Repository (Référence : UserRepository)

```go
package interfaces

import (
    "context"
    "<module>/internal/models"
)

// TodoRepository defines the contract for todo persistence operations.
type TodoRepository interface {
    Create(ctx context.Context, todo *models.Todo) error
    FindByID(ctx context.Context, id uint) (*models.Todo, error)
    FindAll(ctx context.Context, page, limit int) ([]*models.Todo, int64, error)
    Update(ctx context.Context, todo *models.Todo) error
    Delete(ctx context.Context, id uint) error
}
```

### Template Repository GORM (Référence : user_repository.go)

```go
package repository

import (
    "context"
    "errors"
    "<module>/internal/interfaces"
    "<module>/internal/models"
    "gorm.io/gorm"
)

type todoRepository struct {
    db *gorm.DB
}

// NewTodoRepository creates a new GORM-backed TodoRepository.
func NewTodoRepository(db *gorm.DB) interfaces.TodoRepository {
    return &todoRepository{db: db}
}

func (r *todoRepository) Create(ctx context.Context, todo *models.Todo) error {
    return r.db.WithContext(ctx).Create(todo).Error
}

func (r *todoRepository) FindByID(ctx context.Context, id uint) (*models.Todo, error) {
    var todo models.Todo
    err := r.db.WithContext(ctx).First(&todo, id).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    return &todo, err
}

func (r *todoRepository) FindAll(ctx context.Context, page, limit int) ([]*models.Todo, int64, error) {
    var todos []*models.Todo
    var total int64

    r.db.WithContext(ctx).Model(&models.Todo{}).Count(&total)

    offset := (page - 1) * limit
    err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&todos).Error
    return todos, total, err
}

func (r *todoRepository) Update(ctx context.Context, todo *models.Todo) error {
    return r.db.WithContext(ctx).Save(todo).Error
}

func (r *todoRepository) Delete(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).Delete(&models.Todo{}, id).Error
}
```

### Mapping Types → Tags GORM

| Field Type | GORM Tag | Notes |
|------------|----------|-------|
| `string` | `gorm:"type:varchar(255)"` | Défaut, `not null` si spécifié |
| `string:unique` | `gorm:"uniqueIndex;not null"` | Index unique |
| `int/uint` | (aucun tag spécial) | GORM infère le type |
| `float64` | (aucun tag spécial) | GORM infère |
| `bool` | `gorm:"default:false"` | Toujours un default |
| `time.Time` | `gorm:"type:timestamp"` | Pour dates personnalisées |

### Modification de database.go (AutoMigrate)

Le fichier `internal/infrastructure/database/database.go` contient l'appel `AutoMigrate`. Il faut **ajouter** le nouveau modèle sans casser les existants :

```go
// AVANT
if err := db.AutoMigrate(&models.User{}, &models.RefreshToken{}); err != nil {

// APRÈS
if err := db.AutoMigrate(&models.User{}, &models.RefreshToken{}, &models.Todo{}); err != nil {
```

**ATTENTION** : Cette modification est dans le projet GÉNÉRÉ, pas dans le CLI. Le `add-model` doit lire le fichier database.go existant, trouver la ligne AutoMigrate, et ajouter le nouveau modèle.

### Modification de repository/module.go (fx)

```go
// AJOUTER dans le module existant :
fx.Provide(func(db *gorm.DB) interfaces.TodoRepository {
    return NewTodoRepository(db)
}),
```

### Pattern FileGenerator pour la Génération

```go
// Dans add_model.go ou model_generator.go
func generateModelFiles(projectPath, modulePath, modelName string, fields []FieldDefinition) []FileGenerator {
    return []FileGenerator{
        {
            Path:    filepath.Join(projectPath, "internal/models", strings.ToLower(modelName)+".go"),
            Content: generateModelTemplate(modulePath, modelName, fields),
        },
        {
            Path:    filepath.Join(projectPath, "internal/interfaces", strings.ToLower(modelName)+"_repository.go"),
            Content: generateRepositoryInterfaceTemplate(modulePath, modelName),
        },
        {
            Path:    filepath.Join(projectPath, "internal/adapters/repository", strings.ToLower(modelName)+"_repository.go"),
            Content: generateRepositoryImplTemplate(modulePath, modelName, fields),
        },
    }
}
```

### Project Structure Notes

- Les templates sont des strings Go dans le CLI (`cmd/create-go-starter/`)
- Le CLI écrit dans le projet cible (working directory), PAS dans son propre répertoire
- Le `modulePath` (ex: `github.com/user/my-app`) est lu depuis `go.mod` du projet cible
- Utiliser `strings.ToLower()` pour les noms de fichiers, `strings.Title()` pour les types Go

### Anti-Patterns à Éviter

- NE PAS utiliser `text/template` Go - les templates existants sont des string concatenations
- NE PAS modifier le go.mod du projet cible manuellement
- NE PAS créer de dossier `internal/domain/<model>/` dans cette story (c'est Story 8.3)
- NE PAS oublier `gorm.DeletedAt` (soft delete obligatoire)
- NE PAS utiliser `Id` au lieu de `ID` (Go idiomatique)

### References

- [Source: cmd/create-go-starter/templates_user.go:4-46] - User model struct pattern
- [Source: cmd/create-go-starter/templates_user.go:136-299] - UserRepository interface + implementation
- [Source: cmd/create-go-starter/templates.go:856-859] - AutoMigrate pattern
- [Source: cmd/create-go-starter/generator.go] - FileGenerator struct, generateProjectFiles()
- [Source: cmd/create-go-starter/templates_user.go:1402-1422] - Repository fx module pattern
- [Source: _bmad-output/planning-artifacts/architecture.md#Data Architecture] - GORM, DB decisions
- [Source: _bmad-output/project-context.md] - Naming conventions, tags rules

## Dev Agent Record

### Agent Model Used

claude-sonnet-4.5 (github-copilot/claude-sonnet-4.5)

### Debug Log References

N/A - Story implemented via direct code generation

### Completion Notes List

- **Model Generation Pattern**: Followed existing User model pattern from `templates_user.go` for struct layout and GORM tags
- **Template Strategy**: Used string concatenation approach (not `text/template`) consistent with existing CLI codebase
- **File Modification Approach**: Implemented surgical updates to `database.go` and `repository/module.go` using string parsing to preserve existing code
- **Idempotency**: All file modification functions (`updateAutoMigrate`, `updateRepositoryModule`) check for existing entries before adding
- **Test Coverage**: Implemented unit tests for each template function + E2E test with actual `go build` validation
- **Naming Conventions**: Followed Go idioms - `ID` (not `Id`), PascalCase for types, camelCase for unexported, snake_case for JSON/DB
- **GORM Tags**: Applied consistent tagging - `primaryKey`, `autoCreateTime`, `autoUpdateTime`, soft delete with `DeletedAt`
- **Repository Pattern**: Used exact pattern from UserRepository - context propagation, `ErrRecordNotFound` handling, offset/limit pagination
- **fx Module Integration**: New repositories automatically registered as interface providers (not concrete types)

### File List

**Created Files:**
- `cmd/create-go-starter/model.go` - Field parsing logic, validation, type mapping
- `cmd/create-go-starter/model_test.go` - Tests for field parsing and validation
- `cmd/create-go-starter/model_generator.go` - Template generation and file modification logic
- `cmd/create-go-starter/model_generator_test.go` - Tests for generation (unit + E2E)

**Modified Files:**
- `cmd/create-go-starter/add_model.go` - Already existed from Story 8-1 (CLI command entry point)
- `cmd/create-go-starter/add_model_test.go` - Already existed from Story 8-1 (CLI integration tests)
