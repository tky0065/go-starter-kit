# Story 3.2: Operations CRUD Utilisateur

Status: ready-for-dev

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a administrateur (ou développeur),
I want pouvoir lister, modifier ou supprimer des utilisateurs,
so that je puisse gérer la base d'utilisateurs.

## Acceptance Criteria

1.  **List Users (GET /api/v1/users)**
    *   **Given** Je suis authentifié avec un token valide.
    *   **When** J'envoie une requête GET `/api/v1/users`.
    *   **Then** Je reçois une réponse HTTP 200 OK avec une liste d'utilisateurs.
    *   **And** La liste est paginée (par défaut ou explicitement) ou limitée pour éviter la surcharge.
    *   **And** Les mots de passe ne sont PAS retournés.

2.  **Update User (PUT /api/v1/users/:id)**
    *   **Given** Je suis authentifié.
    *   **When** J'envoie une requête PUT `/api/v1/users/:id` avec des données valides (ex: nom, email).
    *   **Then** L'utilisateur correspondant est mis à jour en base de données.
    *   **And** Je reçois une réponse HTTP 200 OK avec les données mises à jour.
    *   **And** Si l'utilisateur n'existe pas, je reçois une erreur 404.
    *   **And** Si les données sont invalides (ex: email malformé), je reçois une erreur 400.

3.  **Delete User (DELETE /api/v1/users/:id)**
    *   **Given** Je suis authentifié.
    *   **When** J'envoie une requête DELETE `/api/v1/users/:id`.
    *   **Then** L'utilisateur est supprimé (Soft Delete via `deleted_at`).
    *   **And** Je reçois une réponse HTTP 200 OK (ou 204 No Content).
    *   **And** L'utilisateur n'apparaît plus dans les listes ultérieures.

## Tasks / Subtasks

- [ ] **Implementation in manual-test-project (Reference Impl)**
    - [ ] **Interfaces:** Update `internal/interfaces/user.go` to include `GetAll`, `Update`, `Delete`.
    - [ ] **Repository:** Add `FindAll`, `Update`, `Delete` methods in `internal/adapters/repository/user_repository.go`.
        - [ ] Ensure `Update` handles partial updates or struct updates correctly.
        - [ ] Ensure `Delete` performs a Soft Delete (GORM default).
    - [ ] **Service:** Implement business logic in `internal/domain/user/service.go`.
        - [ ] Add input validation call if needed (Service vs Handler responsibility - following Hexagonal Lite, Validation is usually in Adapter, but Business Rules in Service).
    - [ ] **Handlers:** Add `GetAllUsers`, `UpdateUser`, `DeleteUser` in `internal/adapters/handlers/user_handler.go`.
        - [ ] Use `go-playground/validator` for input structs (UpdateDTO).
    - [ ] **Routes:** Register endpoints in `internal/infrastructure/server/server.go` (or where routes are defined).
        - [ ] Ensure they are inside the Protected group (JWT Middleware).

- [ ] **CLI Generator Update**
    - [ ] **Templates:** Update Go templates in `cmd/create-go-starter/` to reflect the changes made in `manual-test-project`.
    - [ ] **Refactoring:** If `templates.go` is too large, split it into `templates_user.go` or similar (as noted in Story 3.1).
    - [ ] **Verification:** Ensure the generated project compiles and includes the full CRUD.

## Dev Notes

- **Role Management:**
    - Currently, the PRD does not specify a complex RBAC system. "Admin" implies an authenticated user for this context, or arguably any authenticated user can CRUD for now (MVP).
    - **Decision:** Protect endpoints with the existing JWT Middleware. Future stories can add `role` checks.

- **Data Privacy:**
    - Ensure `GetAll` returns a "Safe User" struct (DTO) without `PasswordHash`.

- **Validation:**
    - For Update, ensure the Email is unique if changed (Repository check might be needed or DB constraint).

### Project Structure Notes

- **Files to Modify:**
    - `manual-test-project/internal/interfaces/user.go`
    - `manual-test-project/internal/adapters/repository/user_repository.go`
    - `manual-test-project/internal/domain/user/service.go`
    - `manual-test-project/internal/adapters/handlers/user_handler.go`
    - `cmd/create-go-starter/templates.go` (or split files)

### References

- [Epics: Story 3.2](_bmad-output/planning-artifacts/epics.md#story-32-operations-crud-utilisateur)
- [Story 3.1](_bmad-output/implementation-artifacts/3-1-gestion-du-profil-utilisateur.md) (Foundation for User Module)

## Dev Agent Record

### Agent Model Used

Gemini 2.0 Flash

### Debug Log References

- Building upon Story 3.1 structure.
- Assuming `manual-test-project` is the primary prototyping ground before updating CLI templates.

### Completion Notes List

- [ ] CRUD endpoints functional in `manual-test-project`.
- [ ] CLI generates project with full CRUD capabilities.
- [ ] `make test` passes.

### File List
