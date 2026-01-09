# Story 4.3: Documentation interactive (Swagger)

Status: ready-for-dev

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a consommateur de l'API,
I want accéder à une documentation Swagger auto-générée,
so that je puisse comprendre et tester l'API sans lire le code source.

## Acceptance Criteria

1.  **Swagger UI Accessible**
    *   **Given** Le serveur est démarré (local ou prod).
    *   **When** J'accède à `/swagger` (avec ou sans slash final, ou `/swagger/index.html`).
    *   **Then** L'interface Swagger UI s'affiche correctement.

2.  **Documentation Auto-générée**
    *   **Given** J'ai ajouté des annotations `swag` (@Summary, @Param, etc.) sur mes handlers.
    *   **When** Je lance la commande de génération (ex: `swag init`).
    *   **Then** Le fichier `docs/swagger.json` (ou yaml) est mis à jour.
    *   **And** L'interface UI reflète ces changements après redémarrage.

3.  **Testabilité**
    *   **Given** Je suis sur l'interface Swagger UI.
    *   **When** J'utilise le bouton "Try it out" sur l'endpoint Login.
    *   **Then** Je peux exécuter la requête et voir la réponse réelle du serveur.

4.  **Intégration CLI**
    *   **Given** Je génère un nouveau projet avec le CLI.
    *   **Then** Le Swagger est pré-configuré et fonctionne immédiatement (au moins avec les routes par défaut Health/Auth).

## Tasks / Subtasks

- [ ] **Swagger Setup (`manual-test-project`)**
    - [ ] Ajouter la dépendance `github.com/swaggo/fiber-swagger`.
    - [ ] Ajouter la dépendance `github.com/swaggo/swag/cmd/swag` (tooling).
    - [ ] Initialiser Swagger : `swag init -g internal/infrastructure/server/server.go` (ou point d'entrée principal).
    - [ ] Créer la route dans `internal/infrastructure/server/server.go` : `app.Get("/swagger/*", fiberSwagger.WrapHandler)`.

- [ ] **Annotations (`manual-test-project`)**
    - [ ] Ajouter les annotations générales (@title, @version, @host) dans `cmd/main.go` ou `server.go`.
    - [ ] Ajouter les annotations sur les handlers existants : `auth_handler.go`, `user_handler.go`.
        -   @Summary, @Description, @Tags, @Accept, @Produce, @Param, @Success, @Failure, @Router.

- [ ] **Makefile Update**
    - [ ] Ajouter une commande `make swagger` ou `make docs` qui exécute `swag init`.
    - [ ] Intégrer cette commande dans le build process si nécessaire.

- [ ] **CLI Generator Update**
    - [ ] Mettre à jour `templates.go` pour inclure le dossier `docs/` généré (ou au moins le squelette).
    - [ ] S'assurer que les templates de handlers incluent les commentaires Swagger par défaut.

## Dev Notes

### Swagger Annotations Guide

**General Info (main.go):**
```go
// @title Go Starter Kit API
// @version 1.0
// @description This is a sample server for Go Starter Kit.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /api/v1
```

**Handler Example:**
```go
// Login godoc
// @Summary User Login
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   request body domain.LoginRequest true "Login Credentials"
// @Success 200 {object} domain.LoginResponse
// @Failure 400 {object} domain.AppError
// @Failure 401 {object} domain.AppError
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error { ... }
```

### Libraries
- **Generator:** `github.com/swaggo/swag` (CLI tool)
- **Fiber Middleware:** `github.com/swaggo/fiber-swagger`

### Integration Issues
- **Docker:** Ensure `swag` is installed in the development container if running `make swagger` inside Docker, OR run it on host and mount `docs/`. Ideally, checking in generated docs is standard practice for Go projects to avoid CI dependency on `swag` CLI.

### Architecture Compliance
- **FR17:** Directly addresses the requirement for interactive documentation.
- **NFR9:** Helps document public functions (handlers).

## Dev Agent Record

### Agent Model Used
Gemini 2.0 Flash

### Debug Log References
- Verified `project-context.md` mentions `swaggo/swag`.
- Confirmed Fiber middleware availability.

### Completion Notes List
- [ ] Dependencies added.
- [ ] Swagger init run.
- [ ] Route registered.
- [ ] Handlers annotated.
- [ ] CLI templates updated.

### File List
