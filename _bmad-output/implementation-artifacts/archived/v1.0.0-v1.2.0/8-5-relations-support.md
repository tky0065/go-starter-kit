# Story 8.5: Relations Support

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

**En tant que** développeur,
**Je veux** définir des relations entre modèles (belongs-to, has-many, many-to-many) lors de `add-model`,
**Afin de** modéliser des données complexes avec des foreign keys et endpoints relationnels.

## Acceptance Criteria

1. **AC1**: Given j'exécute `add-model Comment --fields "content:string" --belongs-to Todo`, When la génération s'exécute, Then la relation GORM BelongsTo est correctement configurée avec foreign key `todo_id`
2. **AC2**: Given une relation belongs-to est définie, When les migrations s'exécutent, Then les foreign keys sont créées dans la base de données
3. **AC3**: Given une relation exists, When les endpoints sont appelés, Then le modèle parent peut être récupéré avec ses relations (preload)
4. **AC4**: Given `--has-many` est utilisé (ex: `add-model User --has-many Todo`), When vérifié, Then le modèle parent inclut un slice du modèle enfant
5. **AC5**: Given toutes les relations supportées, When `--help` est consulté, Then la documentation des flags de relations est claire avec exemples

## Tasks / Subtasks

- [x] Task 1: Supporter le flag `--belongs-to` (AC: 1, 2)
  - [x] 1.1 Parser le flag `--belongs-to <ModelName>` dans add-model
  - [x] 1.2 Ajouter le champ foreign key dans le modèle enfant (`TodoID uint`)
  - [x] 1.3 Ajouter le tag GORM de relation (`gorm:"foreignKey:TodoID"`)
  - [x] 1.4 Ajouter le champ relation dans le struct (`Todo models.Todo`)
  - [x] 1.5 Vérifier que le modèle parent existe dans `internal/models/`

- [x] Task 2: Supporter le flag `--has-many` (AC: 4)
  - [x] 2.1 Parser le flag `--has-many <ModelName>` dans add-model
  - [x] 2.2 Modifier le modèle parent pour ajouter le slice (`Todos []models.Todo`)
  - [x] 2.3 Ajouter le tag GORM `gorm:"foreignKey:ParentID"`
  - [x] 2.4 Modifier le handler parent pour supporter `?include=<children>`

- [x] Task 3: Mettre à jour le repository pour le preloading (AC: 3)
  - [x] 3.1 Ajouter méthode `FindByIDWithRelations(ctx, id, includes []string)` à l'interface
  - [x] 3.2 Implémenter avec `db.Preload("Todos").First(&model, id)`
  - [x] 3.3 Supporter le preloading conditionnel (query param `?include=`)

- [x] Task 4: Mettre à jour les endpoints relationnels (AC: 3)
  - [x] 4.1 GetByID avec preload optionnel via `?include=todos`
  - [x] 4.2 Endpoints imbriqués: `GET /api/v1/todos/:id/comments`
  - [x] 4.3 Création avec relation: `POST /api/v1/todos/:id/comments`

- [x] Task 5: Documentation et aide (AC: 5)
  - [x] 5.1 Documenter les flags de relations dans `--help`
  - [x] 5.2 Exemples dans le message d'aide :
    - `add-model Comment --fields "content:string" --belongs-to Todo`
    - `add-model Todo --has-many Comment`
  - [x] 5.3 Documenter le support `?include=` dans les endpoints

- [x] Task 6: Tests (AC: 1-5)
  - [x] 6.1 Tests du parser de flags de relations
  - [x] 6.2 Tests de génération de modèle avec foreign key
  - [x] 6.3 Tests de modification du modèle parent (has-many)
  - [x] 6.4 Tests E2E: génération avec relations + compilation
  - [x] 6.5 Tests du preloading dans le repository

## Dev Notes

### Relations GORM Supportées

#### BelongsTo (enfant → parent)

```go
// internal/models/comment.go
type Comment struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Content   string         `gorm:"not null" json:"content"`
    TodoID    uint           `gorm:"not null;index" json:"todo_id"`
    Todo      Todo           `gorm:"foreignKey:TodoID" json:"todo,omitempty"`
    CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
```

**Champs générés pour belongs-to :**
- `<Parent>ID uint` avec `gorm:"not null;index"` et `json:"<parent>_id"`
- `<Parent> models.<Parent>` avec `gorm:"foreignKey:<Parent>ID"` et `json:"<parent>,omitempty"`

#### HasMany (parent → enfants)

```go
// internal/models/todo.go (MODIFIÉ)
type Todo struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Title     string         `gorm:"not null" json:"title"`
    Completed bool           `gorm:"default:false" json:"completed"`
    Comments  []Comment      `gorm:"foreignKey:TodoID" json:"comments,omitempty"`
    CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
```

**Champs ajoutés au parent pour has-many :**
- `<Children> []models.<Child>` avec `gorm:"foreignKey:<Parent>ID"` et `json:"<children>,omitempty"`

### Preloading dans le Repository

```go
// Interface étendue
type TodoRepository interface {
    // ... CRUD existant ...
    FindByIDWithRelations(ctx context.Context, id uint, includes []string) (*models.Todo, error)
}

// Implémentation
func (r *todoRepository) FindByIDWithRelations(ctx context.Context, id uint, includes []string) (*models.Todo, error) {
    var todo models.Todo
    query := r.db.WithContext(ctx)

    for _, include := range includes {
        query = query.Preload(include)
    }

    err := query.First(&todo, id).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    return &todo, err
}
```

### Endpoints Relationnels

```go
// Endpoints imbriqués pour Comment sous Todo
// GET /api/v1/todos/:todoId/comments
// POST /api/v1/todos/:todoId/comments

// Routes
todoComments := v1.Group("/todos/:todoId/comments", authMiddleware)
todoComments.Get("", commentHandler.GetByTodo)
todoComments.Post("", commentHandler.CreateForTodo)
```

**Handler pour endpoint imbriqué :**

```go
func (h *CommentHandler) GetByTodo(c *fiber.Ctx) error {
    todoID, err := strconv.ParseUint(c.Params("todoId"), 10, 32)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status": "error", "message": "Invalid todo ID",
        })
    }

    page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 10)

    comments, total, err := h.service.GetByParentID(c.Context(), uint(todoID), page, limit)
    // ... standard response
}
```

### Détection du Modèle Parent

Quand `--belongs-to Todo` est utilisé, le CLI doit :
1. Vérifier que `internal/models/todo.go` existe
2. Parser le fichier pour extraire le nom du type (struct Todo)
3. Si le modèle n'existe pas → erreur claire avec suggestion de le créer d'abord

```go
func validateParentModel(projectPath, parentName string) error {
    modelFile := filepath.Join(projectPath, "internal/models", strings.ToLower(parentName)+".go")
    if _, err := os.Stat(modelFile); os.IsNotExist(err) {
        return fmt.Errorf("parent model '%s' not found at %s. Create it first with: add-model %s --fields ...",
            parentName, modelFile, parentName)
    }
    return nil
}
```

### Modification de Fichiers Existants

Cette story modifie potentiellement des fichiers EXISTANTS dans le projet :
1. Le modèle parent (pour ajouter `[]Children`) quand `--has-many` est utilisé
2. Le repository parent (pour ajouter preloading)
3. Les routes (pour endpoints imbriqués)

**CRITIQUE** : Utiliser des insertions de texte ciblées, PAS de réécriture complète du fichier. Trouver le bon point d'insertion dans le struct existant.

### Flags CLI pour les Relations

```
create-go-starter add-model <Name> [options]

Relation flags:
  --belongs-to <Model>    Add belongs-to relation (adds foreign key)
  --has-many <Model>      Add has-many relation (adds slice to parent)
  --fields "name:type"    Field definitions

Examples:
  add-model Comment --fields "content:string" --belongs-to Todo
  add-model Category --fields "name:string:unique" --has-many Product
```

### Ordre de Migration Important

GORM AutoMigrate gère automatiquement les foreign keys, mais l'ordre compte :
- Le modèle parent DOIT être migré AVANT le modèle enfant
- Donc `&models.Todo{}` avant `&models.Comment{}` dans AutoMigrate

```go
db.AutoMigrate(&models.User{}, &models.RefreshToken{}, &models.Todo{}, &models.Comment{})
```

### Many-to-Many (Scope Limité)

Le support many-to-many est **hors scope** pour cette story. Si demandé, documenter comme enhancement futur :

```go
// Future: Many-to-Many
type Tag struct { ... }
type Todo struct {
    Tags []Tag `gorm:"many2many:todo_tags;"` // Table de jointure auto
}
```

### Anti-Patterns à Éviter

- NE PAS charger les relations par défaut (lazy loading, preload uniquement si `?include=` demandé)
- NE PAS créer de foreign key sans index (`gorm:"not null;index"` obligatoire)
- NE PAS oublier `json:",omitempty"` sur les champs de relation (éviter les null dans la réponse)
- NE PAS permettre de cascader les suppressions sans le documenter
- NE PAS modifier le modèle parent sans backup/vérification
- NE PAS oublier que `DeletedAt` affecte les requêtes de relation (GORM filtre automatiquement)

### References

- [Source: cmd/create-go-starter/templates_user.go:4-46] - User model with GORM tags
- [Source: cmd/create-go-starter/templates_user.go:136-299] - Repository patterns
- [Source: _bmad-output/planning-artifacts/architecture.md#Data Architecture] - GORM, relations
- [Source: _bmad-output/planning-artifacts/epics.md#Story 8.5] - Story specification
- [Source: _bmad-output/project-context.md] - GORM tags rules, naming conventions
- [GORM Docs: https://gorm.io/docs/belongs_to.html] - BelongsTo relations
- [GORM Docs: https://gorm.io/docs/has_many.html] - HasMany relations

## Dev Agent Record

### Agent Model Used

claude-opus-4.6 (via OpenCode / github-copilot)

### Debug Log References

N/A

### Completion Notes List

- All implementation complete: `--belongs-to` and `--has-many` flags fully functional
- `RelationConfig` struct passed through all template generators; `nil` means no relations
- `validateParentModel()` checks file existence at `internal/models/<lowercase>.go`
- Flag parsing follows existing manual arg-parsing pattern (no `flag` package)
- `--belongs-to Parent`: new model gets `ParentID uint` foreign key + `Parent Parent` relation field
- `--has-many Model`: `updateParentModelHasMany()` surgically inserts `[]Child` slice into existing parent model file using multiple anchor points (CreatedAt, UpdatedAt, DeletedAt) with fallback for robustness, includes idempotency check
- `generateHandlerTemplate` uses `result := fmt.Sprintf(...)` pattern so relation handlers can be appended
- Nested routes generated: `v1.Group("/<parentPlural>/:<parentLower>Id/<childPlural>")`
- Repository adds `FindByIDWithRelations` (Preload) and `FindBy<Parent>ID` (WHERE clause) methods
- Generated tests include mock methods and test functions for all relation endpoints
- Many-to-many is out of scope for this story (documented in help text)
- All 10 new relation-specific tests passing; all existing tests unbroken
- `go fmt`, `go vet`, and `go test -short ./cmd/create-go-starter/` all clean
- **Code Review Fixes Applied:**
  - Fixed incomplete File List documentation (added 6 missing files)
  - Improved error handling in `updateParentModelHasMany()` with multiple anchor fallbacks
  - Added pluralization warning in help text for irregular plurals
  - Updated README.md with relation examples and usage
  - Updated docs/usage.md with comprehensive relation documentation
  - All HIGH and MEDIUM issues from code review resolved

### File List

**New Files Created:**
- `cmd/create-go-starter/add_model.go` - Added `RelationConfig` struct, `--belongs-to`/`--has-many` flag parsing, `validateParentModel()`, updated `printModelSummary()` and `printAddModelHelp()`
- `cmd/create-go-starter/add_model_test.go` - Tests for add-model CLI command (help, validation, field parsing)
- `cmd/create-go-starter/model.go` - Model field parsing and validation logic
- `cmd/create-go-starter/model_test.go` - Tests for model field parsing
- `cmd/create-go-starter/database_integration_test.go` - E2E database integration tests

**Modified Files:**
- `cmd/create-go-starter/main.go` - Added `ColorYellow` constant for warning messages, updated ValidDatabases documentation
- `cmd/create-go-starter/templates.go` - Updated to support relation parameters in generated templates
- `cmd/create-go-starter/model_generator.go` - Updated all template generators (`generateModelTemplate`, `generateRepositoryInterfaceTemplate`, `generateRepositoryImplTemplate`, `generateServiceTemplate`, `generateHandlerTemplate`, `generateServiceTestTemplate`, `generateHandlerTestTemplate`, `updateRoutes`) with `*RelationConfig` parameter; added `updateParentModelHasMany()` with improved error handling and fallback anchor points
- `cmd/create-go-starter/model_generator_test.go` - Updated 9 existing call sites to pass `nil` for `*RelationConfig`; added 10 new relation-specific tests
- `README.md` - Updated with relation support examples and usage
- `docs/usage.md` - Added relation documentation with `--belongs-to` and `--has-many` examples
- `_bmad-output/implementation-artifacts/8-5-relations-support.md` - Story status updated to done
- `_bmad-output/implementation-artifacts/sprint-status.yaml` - Story 8-5 and epic-8 status updated to done
