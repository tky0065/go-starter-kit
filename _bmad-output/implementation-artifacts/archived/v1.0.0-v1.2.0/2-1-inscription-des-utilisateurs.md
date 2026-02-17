# Story 2.1: Inscription des utilisateurs (Register)

Status: ready-for-dev

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

- [ ] Définir l'entité User dans `internal/domain/user` (AC: 3)
  - [ ] Ajouter les tags GORM et JSON (snake_case)
  - [ ] Inclure `ID`, `Email`, `PasswordHash`, `CreatedAt`, `UpdatedAt`
- [ ] Créer l'interface Repository dans `internal/interfaces` (AC: 1, 5)
  - [ ] Méthode `CreateUser(user *User) error`
  - [ ] Méthode `GetUserByEmail(email string) (*User, error)`
- [ ] Implémenter le Repository GORM dans `internal/adapters/repository` (AC: 1)
- [ ] Implémenter la logique métier dans `internal/domain/user/service.go` (AC: 3, 5)
  - [ ] Gérer le hachage du mot de passe avec `bcrypt`
  - [ ] Vérifier l'unicité de l'email
- [ ] Créer le Handler HTTP dans `internal/adapters/handlers` (AC: 1, 2, 4)
  - [ ] Définir le DTO `RegisterRequest` avec tags de validation
  - [ ] Valider la requête avec le validator centralisé
  - [ ] Mapper le résultat vers une réponse JSON standardisée
- [ ] Enregistrer les composants dans le système de DI (fx) (AC: 1)
- [ ] Ajouter les annotations Swagger pour la documentation (AC: 1)

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
Gemini 2.0 Flash

### Debug Log References
None

### Completion Notes List
- First story of Epic 2 contexted.
- Detailed implementation steps for the Hexagonal Lite pattern.
- Security requirements for bcrypt cost and sensitive data exposure integrated.

### File List
- internal/domain/user/entity.go
- internal/domain/user/service.go
- internal/interfaces/user_repository.go
- internal/adapters/repository/user_repository.go
- internal/adapters/handlers/auth_handler.go
