# Story 8.1: Add-Model CLI Command

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

**En tant que** développeur,
**Je veux** exécuter `create-go-starter add-model <name>` dans un projet existant,
**Afin de** générer automatiquement le code CRUD pour ce modèle sans écrire le boilerplate moi-même.

## Acceptance Criteria

1. **AC1**: Given je suis dans un projet go-starter-kit existant, When j'exécute `create-go-starter add-model Todo --fields "title:string,completed:bool"`, Then la commande parse les fields correctement et affiche un résumé de ce qui sera généré *(Note: Actual file generation is implemented in Story 8.2 - this story focuses on validation and summary display)*
2. **AC2**: Given la commande parse les fields, When l'utilisateur confirme, Then la commande détecte automatiquement le type de database du projet (postgres/mysql/sqlite) depuis le go.mod
3. **AC3**: Given un nom de modèle invalide ou un projet non-go-starter-kit, When la commande est exécutée, Then un message d'erreur clair est affiché
4. **AC4**: Given `--help` est passé avec add-model, When l'aide s'affiche, Then la syntaxe des fields et les types supportés sont documentés

## Tasks / Subtasks

- [x] Task 1: Refactorer main.go pour supporter les sous-commandes (AC: 1, 3)
  - [x] 1.1 Détecter si le premier argument est `add-model` (sous-commande) ou un nom de projet (comportement actuel)
  - [x] 1.2 Router vers `runAddModel()` ou `runCreateProject()` selon le cas
  - [x] 1.3 Conserver 100% de compatibilité avec le CLI existant (pas de régression)

- [x] Task 2: Parser les fields du modèle (AC: 1)
  - [x] 2.1 Parser le flag `--fields "name:type,name:type"` en `[]FieldDefinition`
  - [x] 2.2 Supporter les types Go courants: `string`, `int`, `uint`, `float64`, `bool`, `time.Time`
  - [x] 2.3 Supporter les modificateurs GORM: `unique`, `not_null`, `index` (ex: `email:string:unique:not_null`)
  - [x] 2.4 Valider les noms de champs (PascalCase automatique) et types

- [x] Task 3: Détecter le contexte projet (AC: 2)
  - [x] 3.1 Vérifier qu'on est dans un projet go-starter-kit (présence de `internal/models/`, `go.mod`)
  - [x] 3.2 Lire `go.mod` pour détecter le database driver (postgres/mysql/sqlite)
  - [x] 3.3 Extraire le module path du projet depuis `go.mod`

- [x] Task 4: Afficher le résumé et demander confirmation (AC: 1)
  - [x] 4.1 Afficher les fichiers qui seront générés avec leur chemin
  - [x] 4.2 Afficher les fichiers existants qui seront modifiés
  - [x] 4.3 Demander confirmation (y/n) avant de procéder
  - [x] 4.4 Supporter `--yes` pour skip la confirmation (CI/CD)

- [x] Task 5: Gestion d'erreurs et aide (AC: 3, 4)
  - [x] 5.1 Erreur si pas dans un projet go-starter-kit
  - [x] 5.2 Erreur si modèle déjà existant (`internal/models/<name>.go` existe)
  - [x] 5.3 Erreur si fields invalides ou vides
  - [x] 5.4 `--help` pour add-model avec exemples et types supportés

- [x] Task 6: Tests unitaires (AC: 1-4)
  - [x] 6.1 Tests du parser de fields (types valides/invalides, modificateurs)
  - [x] 6.2 Tests de détection projet (go.mod parsing, database detection)
  - [x] 6.3 Tests de validation (modèle existant, noms invalides)
  - [x] 6.4 Tests de la sous-commande (routing, help)

## Dev Notes

### Architecture du CLI Actuel

Le CLI actuel dans `cmd/create-go-starter/main.go` utilise un parsing manuel des arguments :
- Arguments positionnels : nom du projet
- Flags : `--template`, `--database`
- Fonction `run()` encapsule toute la logique (testable)
- Fonctions de validation : `validateTemplate()`, `validateDatabase()`

**IMPORTANT** : Ne PAS utiliser cobra ou une lib externe. Le CLI actuel utilise le package `flag` standard. Rester cohérent.

### Pattern de Sous-commande Recommandé

```go
// Dans main.go - run() function
func run() error {
    if len(os.Args) > 1 && os.Args[1] == "add-model" {
        return runAddModel(os.Args[2:])
    }
    // Comportement existant inchangé
    return runCreateProject()
}
```

### Structure FieldDefinition

```go
// Nouveau fichier: cmd/create-go-starter/model.go
type FieldDefinition struct {
    Name       string   // PascalCase (ex: "Title")
    Type       string   // Go type (ex: "string")
    GORMTags   []string // Modificateurs (ex: ["unique", "not_null"])
    JSONName   string   // snake_case (ex: "title")
}

// Types supportés et leur mapping GORM
var supportedTypes = map[string]string{
    "string":    "string",
    "int":       "int",
    "uint":      "uint",
    "float64":   "float64",
    "bool":      "bool",
    "time":      "time.Time",
}
```

### Détection Database depuis go.mod

```go
func detectDatabase(goModPath string) (string, error) {
    content, err := os.ReadFile(goModPath)
    if err != nil { return "", err }

    switch {
    case strings.Contains(string(content), "gorm.io/driver/postgres"):
        return "postgres", nil
    case strings.Contains(string(content), "gorm.io/driver/mysql"):
        return "mysql", nil
    case strings.Contains(string(content), "gorm.io/driver/sqlite"):
        return "sqlite", nil
    default:
        return "", fmt.Errorf("database driver not found in go.mod")
    }
}
```

### Fichiers à Créer/Modifier

**Nouveaux fichiers dans `cmd/create-go-starter/` :**
- `model.go` - Types `FieldDefinition`, parser, validation
- `model_test.go` - Tests unitaires du parser
- `add_model.go` - Fonction `runAddModel()`, détection projet, résumé
- `add_model_test.go` - Tests de la sous-commande

**Fichiers existants à modifier :**
- `main.go` - Router add-model vs create-project dans `run()`

### Conventions de Nommage

- Nom de modèle en argument : `Todo`, `Product`, `BlogPost` (PascalCase)
- Fichiers générés : `todo.go`, `product.go` (snake_case singulier)
- Tables DB : `todos`, `products` (snake_case pluriel - GORM convention)
- JSON fields : `snake_case` (ex: `blog_post_id`)
- Endpoints API : `/api/v1/todos`, `/api/v1/products` (pluriel snake_case)

### Project Structure Notes

- Le CLI est dans `cmd/create-go-starter/`
- Les templates Go string sont dans `templates.go` et `templates_user.go`
- Le pattern `FileGenerator{Path, Content}` est utilisé pour la génération
- Les fichiers de test sont colocalisés (ex: `main.go` ↔ `main_test.go`)

### References

- [Source: cmd/create-go-starter/main.go] - CLI entry point, flag parsing, run() function
- [Source: cmd/create-go-starter/generator.go] - FileGenerator pattern, generateProjectFiles()
- [Source: cmd/create-go-starter/templates.go] - ProjectTemplates struct, template methods
- [Source: cmd/create-go-starter/templates_user.go] - User model/service/handler templates
- [Source: _bmad-output/planning-artifacts/epics.md#Epic 8] - Epic specification
- [Source: _bmad-output/planning-artifacts/architecture.md] - Architecture decisions

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Debug Log References

- Tests failed initially due to relative binary path when using `cmd.Dir` - fixed by using `absBinaryPath()` helper

### Completion Notes List

- **Task 1**: Refactored `main()` in main.go to detect `add-model` as first argument and route to `runAddModel()`. All existing CLI behavior preserved — zero regressions on existing test suite.
- **Task 2**: Implemented `parseFields()` in model.go with support for all 6 Go types (string, int, uint, float64, bool, time.Time) and 3 GORM modifiers (unique, not_null, index). Automatic PascalCase conversion for field names, snake_case for JSON names.
- **Task 3**: Implemented `detectProject()` in add_model.go — validates go-starter-kit project structure (go.mod + internal/models/), detects database driver from go.mod dependencies, extracts module path.
- **Task 4**: Implemented `printModelSummary()` showing model name, database type, module path, fields with types/modifiers, and files to generate/modify. Confirmation prompt with `--yes` flag to skip.
- **Task 5**: Comprehensive error handling for: missing model name, missing --fields, invalid model names (non-PascalCase), invalid field types, invalid GORM modifiers, not in go-starter-kit project, model already exists. Help system with `--help`/`-h` showing syntax, types, modifiers, and examples.
- **Task 6**: 30+ test cases covering field parsing (valid/invalid), model name validation, PascalCase/snake_case conversion, database detection, module path extraction, subcommand routing, help display, error scenarios, and regression testing.

### Change Log

- 2026-02-11: Implemented Story 8.1 — add-model CLI subcommand with field parsing, project detection, summary display, error handling, help system, and comprehensive tests
- 2026-02-11: Code review fixes applied:
  - **CRITICAL #1**: Fixed `toSnakeCase()` to properly handle acronyms (ID→id, UserID→user_id, APIKey→api_key) instead of breaking them (ID→i_d)
  - **CRITICAL #2**: Clarified AC#1 to note that actual file generation is deferred to Story 8.2
  - **HIGH #3**: Strengthened project validation to verify go-starter-kit dependencies (fiber, fx, gorm) in go.mod
  - **HIGH #4**: Added test cases for field names with numbers (V2API, Product2Name, OAuth2Token)
  - **HIGH #5**: Added E2E test `TestE2EAddModelWithRealProject` that creates a real project and runs add-model in it
  - **HIGH #6**: Help text already complete (verified all types and modifiers documented)

### File List

- `cmd/create-go-starter/main.go` (modified) — Added subcommand routing for `add-model`
- `cmd/create-go-starter/model.go` (new) — FieldDefinition type, parseFields(), validateModelName(), toPascalCase(), toSnakeCase()
- `cmd/create-go-starter/model_test.go` (new) — Unit tests for field parser, validators, and converters
- `cmd/create-go-starter/add_model.go` (new) — runAddModel(), detectProject(), detectDatabase(), printModelSummary(), printAddModelHelp()
- `cmd/create-go-starter/add_model_test.go` (new) — Integration tests for add-model subcommand (help, errors, routing, summary)
