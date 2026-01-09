# Story 2.1: Inscription des utilisateurs (Register)

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a **visiteur**,
I want **créer un compte avec mon email et mot de passe**,
so that **je puisse accéder aux fonctionnalités protégées**.

## Acceptance Criteria

1. **Endpoint d'inscription :** Une requête POST sur `/api/v1/auth/register` doit permettre l'inscription.
2. **Validation des entrées :**
    - L'email doit être valide et unique en base de données.
    - Le mot de passe doit être présent (une validation de force peut être ajoutée).
    - Retourner une erreur 400 Bad Request avec des détails clairs si la validation échoue.
3. **Sécurité des données :**
    - Le mot de passe **DOIT** être haché avec `bcrypt` (coût >= 10) avant d'être enregistré.
    - Ne jamais stocker le mot de passe en clair.
4. **Réponse de succès :**
    - Retourner un code HTTP 201 Created en cas de succès.
    - Le corps de la réponse ne doit **JAMAIS** contenir le mot de passe (haché ou non).
    - Retourner l'objet utilisateur créé (ID, Email, CreatedAt).
5. **Gestion des doublons :** Si l'email existe déjà, retourner une erreur 409 Conflict ou 400 Bad Request selon les standards (le PRD suggère une gestion centralisée des erreurs).

## Tasks / Subtasks

- [x] Définir l'entité User dans `internal/domain/user` (AC: 3)
  - [x] Ajouter les tags GORM et JSON (snake_case)
  - [x] Inclure `ID`, `Email`, `PasswordHash`, `CreatedAt`, `UpdatedAt`
- [x] Créer l'interface Repository dans `internal/interfaces` (AC: 1, 5)
  - [x] Méthode `CreateUser(user *User) error`
  - [x] Méthode `GetUserByEmail(email string) (*User, error)`
- [x] Implémenter le Repository GORM dans `internal/adapters/repository` (AC: 1)
- [x] Implémenter la logique métier dans `internal/domain/user/service.go` (AC: 3, 5)
  - [x] Gérer le hachage du mot de passe avec `bcrypt`
  - [x] Vérifier l'unicité de l'email
- [x] Créer le Handler HTTP dans `internal/adapters/handlers` (AC: 1, 2, 4)
  - [x] Définir le DTO `RegisterRequest` avec tags de validation
  - [x] Valider la requête avec le validator centralisé
  - [x] Mapper le résultat vers une réponse JSON standardisée
- [x] Enregistrer les composants dans le système de DI (fx) (AC: 1)
- [x] Ajouter les annotations Swagger pour la documentation (AC: 1)

## Dev Notes

### Architecture & Constraints
- **Pattern :** Hexagonale Lite. Séparation stricte entre Handler (Adapter), Service (Domain) et Repository (Adapter).
- **Security :** Utiliser `golang.org/x/crypto/bcrypt`.
- **Validation :** Utiliser `go-playground/validator/v10`.
- **Naming :** L'entité doit être `User`, la table `users`.

### Technical Guidelines
- Injecter le Repository dans le Service, et le Service dans le Handler via `fx`.
- Utiliser le middleware d'erreur centralisé pour gérer les erreurs de duplication de clé (PostgreSQL unique constraint).
- **Standard Response :** Utiliser l'enveloppe `{"status": "success", "data": {...}}`.

### Project Structure Notes
- Ce module inaugure la structure métier sous `/internal/domain/user`.
- Les interfaces de repository doivent être dans `/internal/interfaces/user_repository.go`.

### References
- [Epic 2: Authentication & Security Foundation](_bmad-output/planning-artifacts/epics.md)
- [Architecture Decision Document](_bmad-output/planning-artifacts/architecture.md)
- [Project Context: Password Hashing with bcrypt](_bmad-output/project-context.md)

## Dev Agent Record

### Agent Model Used
Claude Sonnet 4.5

### Debug Log References
None

### Implementation Plan
- Implemented user registration following Hexagonal Lite architecture
- Used bcrypt with DefaultCost (10) for password hashing as per AC requirement (>= 10)
- Validated email uniqueness before user creation
- Implemented proper error handling for duplicate emails (409 Conflict)
- Added validator/v10 for request validation
- Ensured JSON tags use snake_case convention
- All components registered with fx dependency injection

### Completion Notes List
- All acceptance criteria validated and implemented
- Unit tests pass for all layers (entity, service, repository, handler)
- Integration tests verify the complete registration flow
- Linting passes with golangci-lint v2.7.2
- Code follows project standards from project-context.md
- PasswordHash properly excluded from JSON responses using `json:"-"` tag
- Response follows standard format: `{"status": "success", "data": {...}}`

### File List
- manual-test-project/internal/domain/user/entity.go (modified - added GORM auto timestamps)
- manual-test-project/internal/domain/user/errors.go (new - added sentinel errors)
- manual-test-project/internal/domain/user/service.go (modified - added context.Context, sentinel errors, bcrypt validation)
- manual-test-project/internal/domain/user/service_test.go (modified - added context.Context, enhanced tests)
- manual-test-project/internal/domain/user/module.go (existing - validated)
- manual-test-project/internal/interfaces/user_repository.go (modified - added context.Context)
- manual-test-project/internal/interfaces/user_repository_test.go (modified - fixed linting)
- manual-test-project/internal/adapters/repository/user_repository.go (modified - added context.Context propagation)
- manual-test-project/internal/adapters/repository/module.go (existing - validated)
- manual-test-project/internal/adapters/handlers/auth_handler.go (modified - context, RFC3339, max validation, error handling)
- manual-test-project/internal/adapters/handlers/auth_handler_test.go (modified - added comprehensive validation tests)
- manual-test-project/internal/adapters/handlers/module.go (existing - validated)
- manual-test-project/cmd/main.go (modified - registered auth modules)
- manual-test-project/internal/infrastructure/database/database.go (modified - added User migration)
- manual-test-project/go.mod (modified - added validator/v10 dependency)

## Senior Developer Review (AI)

**Review Date:** 2026-01-09
**Reviewer:** Claude Sonnet 4.5
**Outcome:** Changes Requested (Auto-Fixed)

### Issues Found: 4 High, 6 Medium, 2 Low

### Action Items
- [x] [HIGH] Create sentinel errors for type-safe error handling (errors.go)
- [x] [HIGH] Add context.Context to all functions for timeout/cancellation support
- [x] [HIGH] Fix fragile string-based error matching in handler
- [x] [HIGH] Add bcrypt hash validation (ensure non-empty result)
- [x] [MEDIUM] Fix date format to RFC3339 standard
- [x] [MEDIUM] Add max length validations (email: 255, password: 72)
- [x] [MEDIUM] Improve validation error messages for user-friendly output
- [x] [MEDIUM] Inject validator properly instead of global variable
- [x] [MEDIUM] Add comprehensive validation tests (invalid email, empty fields, etc.)
- [x] [MEDIUM] Update File List with missed files (cmd/main.go, database.go)
- [x] [LOW] Remove TODO comment admitting non-ideal implementation
- [x] [LOW] Increase min password length from 6 to 8

### Review Summary
All HIGH and MEDIUM issues were automatically fixed. The implementation now follows Go best practices with:
- Type-safe error handling using sentinel errors
- Proper context propagation for request lifecycle management
- RFC3339 date formatting for API responses
- Comprehensive input validation with user-friendly error messages
- 7 test scenarios covering all validation edge cases

## Change Log
- **Date:** 2026-01-08
  - Initial implementation: Validated and adjusted existing user registration
  - Fixed entity.go GORM tags to use autoCreateTime/autoUpdateTime
  - Added validator/v10 dependency for request validation
  - Fixed linting issue in user_repository_test.go
  - All tests passing, linting clean, build successful

- **Date:** 2026-01-09
  - Code review fixes: Applied 12 corrections (4 HIGH, 6 MEDIUM, 2 LOW)
  - Added sentinel errors (user/errors.go) for type-safe error handling
  - Added context.Context propagation across all layers
  - Fixed date format to RFC3339, added max length validations
  - Enhanced test coverage: 7 validation scenarios
  - Improved error messages for user-friendly validation feedback
  - All tests passing (100%), linting clean
