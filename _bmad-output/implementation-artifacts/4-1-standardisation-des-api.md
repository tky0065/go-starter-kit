# Story 4.1: Standardisation des API (Grouping & V1)

Status: ready-for-dev

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a développeur,
I want que toutes les routes soient préfixées par `/api/v1`,
so that je puisse versionner mon API facilement à l'avenir et maintenir une structure cohérente.

## Acceptance Criteria

1.  **Global Prefix /api/v1**
    *   **Given** Le serveur est démarré.
    *   **When** J'accède à n'importe quelle route métier (ex: Auth, Users).
    *   **Then** L'URL doit OBLIGATOIREMENT commencer par `/api/v1` (ex: `/api/v1/auth/login`, `/api/v1/users`).
    *   **And** Les anciennes routes (si elles existaient sans préfixe) ne sont plus accessibles ou sont redirigées.

2.  **Route Grouping**
    *   **Given** Je regarde le code de définition des routes (`server.go` ou `router.go`).
    *   **Then** Je vois l'utilisation explicite des `Group` de Fiber.
    *   **And** Il y a un groupe principal `v1`.
    *   **And** Il y a des sous-groupes pour chaque domaine (`auth`, `users`, etc.).

3.  **Consistency Check**
    *   **Given** L'application tourne.
    *   **When** Je teste tous les endpoints définis dans les Epics précédents (Register, Login, Me, CRUD Users).
    *   **Then** Ils répondent tous correctement sous le nouveau préfixe `/api/v1`.

## Tasks / Subtasks

- [ ] **Implementation in manual-test-project (Reference Impl)**
    - [ ] **Server Setup:** Refactor `internal/infrastructure/server/server.go` to create a root group `/api` and a version group `/v1`.
    - [ ] **Route Refactoring:** Move existing route registrations into specific functions or methods that attach to the `v1` group (e.g., `RegisterAuthRoutes(v1)`, `RegisterUserRoutes(v1)`).
    - [ ] **Verification:** Update any integration tests or manual requests (e.g., `.http` files or Postman collections) to use the new URL structure.

- [ ] **CLI Generator Update**
    - [ ] **Templates:** Update the Server Template in `cmd/create-go-starter/templates.go` (or `templates_server.go`) to generate this grouping structure by default.
    - [ ] **New Projects:** Ensure that when a user generates a NEW project, it comes with `/api/v1` pre-configured.

## Dev Notes

- **Fiber Groups:**
    ```go
    api := app.Group("/api")
    v1 := api.Group("/v1")
    
    authGroup := v1.Group("/auth")
    userGroup := v1.Group("/users")
    ```

- **Impact on Existing Tests:**
    - This is a breaking change for API consumers (tests). ALL integration tests must be updated to include `/api/v1` in the request paths.

- **Architecture Compliance:**
    - This fulfills **FR14** and **FR13**.
    - Keeps the `RegisterRoutes` logic clean and modular.

### Project Structure Notes

- **Files to Modify:**
    - `manual-test-project/internal/infrastructure/server/server.go`
    - `manual-test-project/tests/integration/...` (all test files)
    - `cmd/create-go-starter/templates.go`

### References

- [Epics: Story 4.1](_bmad-output/planning-artifacts/epics.md#story-41-standardisation-des-api-grouping--v1)
- [Architecture: API Patterns](_bmad-output/planning-artifacts/architecture.md#api-communication-patterns)
- [Project Context: Routing Rules](_bmad-output/project-context.md#framework-specific-rules-fiber--gorm)

## Dev Agent Record

### Agent Model Used

Gemini 2.0 Flash

### Debug Log References

- Verified `project-context.md` explicitly mandates `/api/v1` prefix.
- Validated logic against Fiber documentation patterns for Groups.

### Completion Notes List

- [ ] All routes in `manual-test-project` moved to `/api/v1`.
- [ ] Integration tests updated and passing.
- [ ] CLI generator produces project with `/api/v1` structure.

### File List
