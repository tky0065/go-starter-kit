# Story 3.1: Gestion du Profil Utilisateur (Me)

Status: ready-for-dev

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a utilisateur connecté,
I want consulter mon propre profil,
so that je puisse vérifier mes informations de compte.

## Acceptance Criteria

1. **Given** Je suis authentifié avec un token valide (JWT)
   **When** J'envoie une requête GET `/api/v1/users/me`
   **Then** Je reçois une réponse HTTP 200 OK avec mes informations (ID, email, nom, created_at)
   **And** Les informations retournées correspondent à l'utilisateur identifié par le token (ID extrait du contexte)
   **And** Aucune donnée sensible (mot de passe haché) n'est retournée.

2. **Given** Je ne suis PAS authentifié (pas de token ou token invalide)
   **When** J'envoie une requête GET `/api/v1/users/me`
   **Then** Je reçois une réponse HTTP 401 Unauthorized (géré par le middleware Auth existant).

## Tasks / Subtasks

- [ ] **Analysis & Prototype (manual-test-project)**
    - [ ] Verify existing Auth Middleware in `manual-test-project`.
    - [ ] Define `GetUser` interface in `internal/interfaces`.
    - [ ] Implement `GetProfile` logic in `internal/domain/user/service.go`.
    - [ ] Implement `FindByID` in `internal/adapters/repository/user_repository.go` (if not exists).
    - [ ] Implement `GetMe` handler in `internal/adapters/handlers/user_handler.go`.
    - [ ] Register route `GET /api/v1/users/me` in Fiber app (protected by Auth Middleware).
    - [ ] Add Integration Test for `GetMe`.

- [ ] **Generator Implementation (CLI)**
    - [ ] **CRITICAL:** Port existing Auth logic (Epic 2) from `manual-test-project` to `cmd/create-go-starter/templates.go` (or new `templates_user.go` / `templates_auth.go`) if not already present.
    - [ ] Create templates for User Domain, Repository, Service, and Handlers.
    - [ ] Update `generator.go` to write these new files during project generation.
    - [ ] Update `MainGoTemplate` or `ServerTemplate` to wire up the new User module via `fx`.

## Dev Notes

- **Existing Codebase State:**
    - Epic 2 (Auth) is marked "Done", and code exists in `manual-test-project/internal/adapters/handlers/auth_handler.go`.
    - **WARNING:** The CLI templates (`cmd/create-go-starter/templates.go`) DO NOT appear to contain the Auth code yet. You must bridge this gap.
    - The `manual-test-project` is the source of truth for the implementation patterns.

- **Relevant Architecture Patterns:**
    - **Hexagonal Lite:** Keep logic in Service, persistence in Repo, HTTP transport in Handler.
    - **Interfaces:** Define interactions in `internal/interfaces`.
    - **Response Format:** Use the standard JSON envelope (`status`, `data`).
    - **Dependency Injection:** Use `uber-go/fx` modules.

- **Security:**
    - Extract UserID from the Fiber Local Context (set by JWT Middleware).
    - NEVER accept UserID as a URL parameter for "My Profile" (e.g. `/users/:id` for "me" is bad practice, use `/users/me` or `/me`).

### Project Structure Notes

- **Target Files (Generated Project):**
    - `internal/domain/user/entity.go` (User struct)
    - `internal/domain/user/service.go`
    - `internal/adapters/handlers/user_handler.go`
    - `internal/adapters/repository/user_repository.go`
    - `internal/interfaces/user.go`

- **CLI Files:**
    - `cmd/create-go-starter/templates.go` (Consider splitting if too large, e.g., `templates_user.go`).
    - `cmd/create-go-starter/generator.go`

### References

- [Epics: Story 3.1](_bmad-output/planning-artifacts/epics.md#story-31-gestion-du-profil-utilisateur-me)
- [Architecture: API Patterns](_bmad-output/planning-artifacts/architecture.md#api-naming-conventions)
- [Project Context: Tech Stack](_bmad-output/project-context.md#technology-stack--versions)

## Dev Agent Record

### Agent Model Used

Gemini 2.0 Flash

### Debug Log References

- Validated `manual-test-project` has Auth implementation.
- Detected missing templates in `create-go-starter`.

### Completion Notes List

- [ ] Verified `GetMe` works in manual-test-project.
- [ ] Verified `create-go-starter` generates a project containing the new User/Auth modules.
- [ ] Validated generated project passes `make test`.

### File List
- `cmd/create-go-starter/templates.go`
- `cmd/create-go-starter/generator.go`
- `manual-test-project/internal/adapters/handlers/user_handler.go`
- `manual-test-project/internal/domain/user/service.go`
