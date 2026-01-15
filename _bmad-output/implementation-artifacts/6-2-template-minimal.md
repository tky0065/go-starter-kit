# Story 6.2: Template Minimal

Status: done

## Story

**En tant que** développeur,
**Je veux** générer un projet Go minimal avec API REST et Swagger,
**Afin de** démarrer rapidement sans la complexité de l'authentification.

## Acceptance Criteria

1. **AC1**: Given l'utilisateur exécute `create-go-starter mon-projet --template=minimal`, When le projet est généré, Then la structure inclut: Fiber, GORM, Swagger, Health check, Logger ✅
2. **AC2**: Given le projet minimal est généré, When on vérifie les fichiers, Then il n'y a PAS d'authentification JWT ni de gestion utilisateurs ✅
3. **AC3**: Given le projet minimal est généré, When on exécute `go build ./...`, Then le projet compile sans erreur ✅
4. **AC4**: Given le projet minimal est lancé, When on accède à `/health` et `/swagger/*`, Then les endpoints répondent correctement ✅

## Tasks / Subtasks

- [x] Task 1: Créer la structure conditionnelle dans generator.go (AC: 1, 2)
  - [x] 1.1 Modifier `generateProjectFiles(projectPath, projectName, template string)`
  - [x] 1.2 Créer `getDirectoriesForTemplate(template string) []string`
  - [x] 1.3 Créer `generateMinimalTemplateFiles(projectPath, projectName string) error`

- [x] Task 2: Créer les templates minimaux (AC: 1)
  - [x] 2.1 `MinimalMainGoTemplate()` - main.go sans auth, avec Swagger
  - [x] 2.2 `MinimalRoutesTemplate()` - routes.go sans auth middleware
  - [x] 2.3 `MinimalServerTemplate()` - server.go simplifié
  - [x] 2.4 `MinimalGoModTemplate()` - go.mod sans dépendances JWT
  - [x] 2.5 `MinimalReadmeTemplate()` - README adapté
  - [x] 2.6 `MinimalEnvTemplate()` - .env.example sans JWT_SECRET
  - [x] 2.7 `MinimalDatabaseTemplate()` - database.go sans User migrations
  - [x] 2.8 `MinimalSetupScriptTemplate()` - setup.sh sans JWT setup
  - [x] 2.9 `MinimalDockerComposeTemplate()` - docker-compose sans JWT_SECRET
  - [x] 2.10 `MinimalDocsReadmeTemplate()` - docs/README.md
  - [x] 2.11 `MinimalQuickStartTemplate()` - docs/quick-start.md

- [x] Task 3: Adapter les fichiers partagés (AC: 1)
  - [x] 3.1 Réutiliser: ConfigTemplate, LoggerTemplate
  - [x] 3.2 Réutiliser: HealthHandlerTemplate, DockerfileTemplate, MakefileTemplate
  - [x] 3.3 Réutiliser: GitignoreTemplate, GolangCILintTemplate, GitHubActionsWorkflowTemplate
  - [x] 3.4 Réutiliser: SwaggerDocsTemplate

- [x] Task 4: Mettre à jour la structure de répertoires (AC: 1)
  - [x] 4.1 Créer `getDirectoriesForTemplate(template string) []string`
  - [x] 4.2 Modifier `createProjectStructure(projectPath, template string) error`
  - [x] 4.3 Minimal n'a pas besoin de: `pkg/auth/`, `internal/domain/user/`, `internal/adapters/handlers/`

- [x] Task 5: Tests (AC: 3, 4)
  - [x] 5.1 Test génération template minimal (TestGenerateMinimalProjectFiles)
  - [x] 5.2 Test compilation du projet généré (TestE2EMinimalProjectBuilds)
  - [x] 5.3 Test absence des fichiers auth (TestMinimalTemplateNoAuthFiles)
  - [x] 5.4 Tests unitaires templates (TestMinimal*Template)
  - [x] 5.5 Test directories pour minimal (TestGetDirectoriesForMinimalTemplate)

## Dev Notes

### Dépendance Story 6.1

Cette story dépend de la story 6.1 (flag `--template`). Le flag doit être implémenté et le paramètre `template` doit être passé à `generateProjectFiles()`.

### Fichiers générés pour template "minimal"

```
mon-projet/
├── cmd/
│   └── main.go                    # Simplifié, sans auth
├── internal/
│   ├── adapters/
│   │   └── http/
│   │       ├── health.go          # Identique
│   │       └── routes.go          # Sans auth routes
│   └── infrastructure/
│       ├── database/
│       │   └── database.go        # Identique
│       └── server/
│           └── server.go          # Sans auth middleware
├── pkg/
│   ├── config/
│   │   └── env.go                 # Identique
│   └── logger/
│       └── logger.go              # Identique
├── docs/
│   ├── README.md
│   ├── docs.go                    # Swagger
│   └── quick-start.md
├── .env.example                   # Sans JWT_SECRET
├── .gitignore
├── .golangci.yml
├── Dockerfile
├── Makefile
├── README.md
├── go.mod                         # Sans deps JWT
└── setup.sh
```

### Fichiers EXCLUS du template minimal

- `pkg/auth/jwt.go`
- `pkg/auth/middleware.go`
- `pkg/auth/module.go`
- `internal/models/user.go`
- `internal/domain/errors.go`
- `internal/domain/user/service.go`
- `internal/domain/user/module.go`
- `internal/interfaces/services.go`
- `internal/interfaces/user_repository.go`
- `internal/adapters/middleware/error_handler.go`
- `internal/adapters/repository/user_repository.go`
- `internal/adapters/repository/module.go`
- `internal/adapters/handlers/auth_handler.go`
- `internal/adapters/handlers/user_handler.go`
- `internal/adapters/handlers/module.go`

### Pattern de code suggéré

```go
// generator.go

func generateProjectFiles(projectPath, projectName, template string) error {
    // ... validation existante ...

    templates := NewProjectTemplates(projectName)

    // Obtenir les fichiers selon le template
    var files []FileGenerator
    switch template {
    case "minimal":
        files = getMinimalFiles(templates, projectPath)
    case "full":
        files = getFullFiles(templates, projectPath) // actuel
    case "graphql":
        files = getGraphQLFiles(templates, projectPath) // story 6.4
    default:
        files = getFullFiles(templates, projectPath)
    }

    // ... écriture des fichiers ...
}

func getMinimalFiles(t *ProjectTemplates, projectPath string) []FileGenerator {
    return []FileGenerator{
        {Path: filepath.Join(projectPath, "go.mod"), Content: t.MinimalGoModTemplate()},
        {Path: filepath.Join(projectPath, "cmd", "main.go"), Content: t.MinimalMainGoTemplate()},
        // ... autres fichiers minimal ...
    }
}
```

### Templates à créer dans templates.go

1. **MinimalGoModTemplate()** - Sans dépendances JWT:
   - Retirer: `github.com/gofiber/contrib/jwt`, `github.com/golang-jwt/jwt/v5`

2. **MinimalMainGoTemplate()** - main.go simplifié:
   - Sans import auth
   - Sans fx.Module pour auth
   - Swagger activé

3. **MinimalRoutesTemplate()** - routes.go:
   - Uniquement `/health` et `/swagger/*`
   - Pas de groupes `/api/v1/auth` ou `/api/v1/users`

4. **MinimalServerTemplate()** - server.go:
   - Sans auth middleware
   - Configuration Fiber basique

5. **MinimalEnvTemplate()** - .env.example:
   - Sans `JWT_SECRET`
   - `DATABASE_URL`, `PORT`, `ENV` uniquement

### Structure de répertoires pour minimal

```go
func getMinimalDirectories() []string {
    return []string{
        "cmd",
        "internal/adapters/http",
        "internal/infrastructure/database",
        "internal/infrastructure/server",
        "pkg/config",
        "pkg/logger",
        "docs",
    }
}
```

### Project Structure Notes

- Réutiliser le maximum de code entre templates
- Les fichiers communs (health.go, database.go, logger.go) restent identiques
- Seuls les fichiers avec auth-specific code ont besoin de variantes

### References

- [Source: cmd/create-go-starter/generator.go] - Structure actuelle
- [Source: cmd/create-go-starter/templates.go] - Templates existants
- [Source: _bmad-output/planning-artifacts/epic-6.md#Story 6.2] - Spécifications

## Dev Agent Record

### Agent Model Used

Claude Sonnet 4 (Anthropic)

### Completion Notes List

- **Phase RED**: Created failing tests in `template_minimal_test.go` covering all acceptance criteria
- **Phase GREEN**: Implemented all minimal templates and generator logic
- **Key changes**:
  - Added `getDirectoriesForTemplate()` function returning template-specific directories
  - Added `generateMinimalTemplateFiles()` function to generate minimal project files
  - Updated `createProjectStructure()` to accept template parameter
  - Created 11 new Minimal*Template() methods in `templates_minimal.go`
  - Updated switch statement in `generateProjectFiles()` to handle "minimal" template
  - Updated existing tests to pass the template parameter to `createProjectStructure()`
- **All tests pass**: Unit tests, integration tests, and E2E tests including build verification

### File List

**New files:**
- `cmd/create-go-starter/templates_minimal.go` - All minimal template methods (815 lines)
- `cmd/create-go-starter/template_minimal_test.go` - All minimal template tests (510 lines)

**Modified files:**
- `cmd/create-go-starter/generator.go` - Added `getDirectoriesForTemplate()`, `generateMinimalTemplateFiles()`, updated switch
- `cmd/create-go-starter/main.go` - Updated `createProjectStructure()` to accept template parameter
- `cmd/create-go-starter/generator_test.go` - Updated tests to pass TemplateFull
- `cmd/create-go-starter/scaffold_test.go` - Updated tests to pass TemplateFull
- `cmd/create-go-starter/smoke_test.go` - Updated tests to pass TemplateFull
- `cmd/create-go-starter/templates_test.go` - Updated tests to pass TemplateFull
- `cmd/create-go-starter/main_test.go` - Updated minimal template test to expect success
- `cmd/create-go-starter/git_test.go` - Updated test messages from French to English

## Senior Developer Review (AI)

**Review Date:** 2026-01-15
**Reviewer:** Claude Sonnet 4 (Code Review Agent)
**Outcome:** ✅ APPROVED (with fixes applied)

### Issues Found & Fixed

| Severity | Issue | Resolution |
|----------|-------|------------|
| 🔴 HIGH | `SetConnMaxLifetime(5 * 60)` used int instead of time.Duration - caused 300ns lifetime instead of 5 minutes | Fixed: Changed to `5 * time.Minute` |
| 🔴 HIGH | Missing `time` import in MinimalDatabaseTemplate | Fixed: Added `"time"` to imports |
| 🟡 MEDIUM | git_test.go modified but not documented in File List | Fixed: Added to File List |

### Issues Noted (Low Priority - Not Blocking)

- Code duplication between `generateFullTemplateFiles` and `generateMinimalTemplateFiles` (DRY violation) - recommend refactoring in future
- Missing dedicated test for `MinimalDockerComposeTemplate()`
- Hardcoded GitHub URL in documentation templates

### Acceptance Criteria Verification

- ✅ AC1: Fiber, GORM, Swagger, Health, Logger included
- ✅ AC2: No JWT auth, no user management
- ✅ AC3: Project compiles with `go build ./...`
- ✅ AC4: /health and /swagger/* endpoints registered

### Test Results

All tests pass including E2E build verification:
- `TestE2EMinimalProjectBuilds` ✅
- `TestMinimalTemplateNoAuthFiles` ✅
- `TestGetDirectoriesForMinimalTemplate` ✅

## Change Log

| Date | Author | Changes |
|------|--------|---------|
| 2026-01-15 | Dev Agent | Initial implementation - all tasks completed |
| 2026-01-15 | Review Agent | Fixed HIGH issues: time.Duration bug and missing import; Updated File List |
