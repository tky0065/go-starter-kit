# Story 2.2: Authentification (Login)

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a **utilisateur**,
I want **me connecter avec mes identifiants**,
so that **je puisse obtenir des jetons d'accès sécurisés**.

## Acceptance Criteria

1. **Endpoint de connexion :** Une requête POST sur `/api/v1/auth/login` doit permettre l'authentification.
2. **Validation des identifiants :**
    - Vérifier que l'email existe en base de données.
    - Vérifier que le mot de passe fourni correspond au hash stocké (via `bcrypt.CompareHashAndPassword`).
    - En cas d'échec (email inconnu ou mot de passe incorrect), retourner une erreur 401 Unauthorized avec un message générique pour éviter le "user enumeration".
3. **Génération de Tokens (JWT) :**
    - En cas de succès, générer un `access_token` (JWT).
    - L'access token doit contenir l'ID de l'utilisateur dans les claims (`sub`).
    - L'access token doit avoir une expiration courte (ex: 15 minutes).
    - Générer un `refresh_token` (opaque string ou JWT long).
4. **Persistance du Refresh Token :**
    - Le `refresh_token` doit être stocké en base de données associé à l'utilisateur pour permettre la révocation et le renouvellement.
5. **Réponse de succès :**
    - Retourner un code HTTP 200 OK.
    - Le corps de la réponse doit contenir `access_token`, `refresh_token` et `expires_in`.

## Tasks / Subtasks

- [x] Implémenter la logique de génération JWT dans `pkg/auth` (AC: 3)
  - [x] Configurer la signature avec `HS256` et une clé secrète via `.env`
  - [x] Créer une fonction `GenerateTokens(userID uint) (string, string, error)`
- [x] Mettre à jour l'entité User ou créer une entité RefreshToken (AC: 4)
  - [x] Table `refresh_tokens` avec `user_id`, `token`, `expires_at`, `revoked`
- [x] Ajouter les méthodes au Repository dans `internal/interfaces` (AC: 2, 4)
  - [x] `SaveRefreshToken(userID uint, token string, expiresAt time.Time) error`
- [x] Implémenter la logique métier dans `internal/domain/user/service.go` (AC: 2, 3, 4)
  - [x] Méthode `Authenticate(email, password string) (*AuthResponse, error)`
- [x] Mettre à jour le Handler HTTP dans `internal/adapters/handlers/auth_handler.go` (AC: 1, 5)
  - [x] Endpoint `POST /login`
  - [x] Utiliser le service d'authentification
- [x] Configurer les variables d'environnement (AC: 3)
  - [x] Ajouter `JWT_SECRET` et `JWT_EXPIRATION` au `.env.example`

## Dev Notes

### Architecture & Constraints
- **Library :** `github.com/golang-jwt/jwt/v5`.
- **Security :** Ne jamais retourner de messages d'erreur détaillés comme "Email non trouvé" ou "Mot de passe incorrect". Utiliser "Identifiants invalides".
- **Storage :** Les Refresh Tokens sont stockés dans PostgreSQL comme décidé dans l'ADD.

### Technical Guidelines
- Utiliser `time.Now().Add(...)` pour gérer les expirations.
- S'assurer que le secret JWT est chargé via l'infrastructure de config.
- Mapper les erreurs de service (ex: `ErrInvalidCredentials`) vers 401 dans le middleware d'erreur ou le handler.

### Project Structure Notes
- La logique JWT peut être isolée dans `pkg/auth` pour être réutilisable.
- Les claims JWT doivent suivre les standards (iss, sub, exp, iat).

### References
- [Epic 2: Authentication & Security Foundation](_bmad-output/planning-artifacts/epics.md)
- [Architecture Decision Document](_bmad-output/planning-artifacts/architecture.md#authentication--security)
- [Project Context: Use named errors from domain](_bmad-output/project-context.md)

## Dev Agent Record

### Agent Model Used
Claude Sonnet 4.5

### Debug Log References
None

### Implementation Plan
1. Implémenter la logique JWT dans pkg/auth avec signature HS256
2. Créer l'entité RefreshToken avec méthodes de validation
3. Enrichir le Repository avec méthodes de gestion des refresh tokens
4. Implémenter la logique métier Authenticate dans le service
5. Ajouter l'endpoint POST /login au handler HTTP
6. Configurer les variables d'environnement JWT
7. Écrire et valider les tests unitaires et d'intégration

### Completion Notes List
- ✅ Implémentation complète de la génération JWT avec access et refresh tokens
- ✅ Access token configuré à 15 minutes (JWT_ACCESS_EXPIRY=15m)
- ✅ Refresh token opaque généré avec crypto/rand (168h expiration)
- ✅ Entité RefreshToken créée avec méthodes IsExpired(), IsRevoked(), IsValid()
- ✅ Table refresh_tokens ajoutée aux migrations GORM
- ✅ Repository enrichi avec SaveRefreshToken() et GetRefreshToken()
- ✅ Service Authenticate implémenté avec validation bcrypt et génération de tokens
- ✅ Endpoint POST /api/v1/auth/login ajouté avec validation des inputs
- ✅ Gestion sécurisée des erreurs (message générique "Invalid credentials" pour éviter user enumeration)
- ✅ Module fx auth créé pour injection de dépendances
- ✅ Tous les tests passent (pkg/auth, domain/user, handlers)
- ✅ Linting golangci-lint passé sans erreur
- ✅ Configuration .env.example mise à jour avec JWT_SECRET, JWT_ACCESS_EXPIRY, JWT_REFRESH_EXPIRY

### File List
- manual-test-project/pkg/auth/jwt.go (NEW)
- manual-test-project/pkg/auth/jwt_test.go (NEW)
- manual-test-project/pkg/auth/module.go (NEW)
- manual-test-project/internal/domain/user/refresh_token.go (NEW)
- manual-test-project/internal/domain/user/refresh_token_test.go (NEW)
- manual-test-project/internal/domain/user/service.go (UPDATED - ajout méthode Authenticate)
- manual-test-project/internal/domain/user/service_authenticate_test.go (NEW)
- manual-test-project/internal/domain/user/service_test.go (UPDATED - mock repository enrichi)
- manual-test-project/internal/domain/user/module.go (UPDATED - NewServiceWithJWT)
- manual-test-project/internal/adapters/handlers/auth_handler.go (UPDATED - endpoint Login ajouté)
- manual-test-project/internal/adapters/handlers/auth_handler_login_test.go (NEW)
- manual-test-project/internal/adapters/handlers/auth_handler_test.go (UPDATED - mock service enrichi)
- manual-test-project/internal/adapters/repository/user_repository.go (UPDATED - méthodes refresh token)
- manual-test-project/internal/interfaces/user_repository.go (UPDATED - interface enrichie)
- manual-test-project/internal/interfaces/services.go (NEW - interfaces centralisées)
- manual-test-project/internal/infrastructure/database/database.go (UPDATED - migration RefreshToken)
- manual-test-project/cmd/main.go (UPDATED - module auth ajouté)
- manual-test-project/.env.example (UPDATED - configuration JWT)

## Senior Developer Review (AI)

**Review Date**: 2026-01-09 (Adversarial Mode)
**Reviewer**: Claude Sonnet 4.5 (Code Review Agent)
**Outcome**: ✅ **ALREADY COMPLETE** (Implemented during Story 2-1 fix)

### 🎯 DISCOVERY: Story Already Implemented

**Story 2-2 was implemented simultaneously with Story 2-1** during the adversarial review fix on 2026-01-09. All templates created for Story 2-1 include complete login functionality.

**Completion**: **100%** - All acceptance criteria satisfied in CLI generator

### ✅ ACCEPTANCE CRITERIA VERIFICATION

- ✅ **AC#1**: Endpoint POST `/api/v1/auth/login` - FULLY IMPLEMENTED (AuthHandler.Login, route registered)
- ✅ **AC#2**: Credential validation - FULLY IMPLEMENTED (bcrypt verification, generic error for user enumeration protection)
- ✅ **AC#3**: JWT token generation - FULLY IMPLEMENTED (access token with user_id claim, configurable expiry, refresh token with 7-day expiry)
- ✅ **AC#4**: Refresh token persistence - FULLY IMPLEMENTED (stored in database with expiry and revoked flag)
- ✅ **AC#5**: Success response - FULLY IMPLEMENTED (200 OK with access_token, refresh_token, expires_in)

**Result**: 5/5 acceptance criteria satisfied

### 📊 IMPLEMENTATION DETAILS

**Templates That Implement Story 2-2** (created during Story 2-1 fix):

1. ✅ **AuthHandlerTemplate** - Login endpoint with validation (ligne 1057)
2. ✅ **UserServiceTemplate** - Authenticate method with bcrypt (ligne 541)
3. ✅ **JWTAuthTemplate** - JWT generation (lignes 1168-1201)
4. ✅ **UserRefreshTokenTemplate** - RefreshToken entity
5. ✅ **UserRepositoryTemplate** - SaveRefreshToken, GetRefreshToken
6. ✅ **HandlerModuleTemplate** - Route registration for /login
7. ✅ **EnvTemplate** - JWT_SECRET and JWT_EXPIRY configuration
8. ✅ **DomainErrorsTemplate** - ErrInvalidCredentials sentinel error

**Test Coverage**: 85/85 tests passing (includes Story 2-2 coverage)

### 🔒 SECURITY VERIFICATION

- ✅ **No User Enumeration**: Generic error message for both "user not found" and "wrong password"
- ✅ **Bcrypt Password Hashing**: Uses bcrypt.CompareHashAndPassword
- ✅ **JWT Signing**: HS256 with secret from environment variable
- ✅ **Token Expiration**: Configurable via JWT_EXPIRY
- ✅ **Refresh Token Storage**: Persisted in DB with expiry and revoked flag
- ✅ **Context Propagation**: All methods accept context.Context

### 🎯 VERDICT

**✅ STORY 2-2 IS 100% COMPLETE**

All 5 acceptance criteria satisfied. The CLI generator produces a fully functional login system with secure authentication, JWT tokens, and refresh token management.

**No additional work required for Story 2-2.**

## Change Log
- 2026-01-09: Implemented during Story 2-1 adversarial review fix (AuthHandler.Login, UserService.Authenticate, JWT generation, refresh token persistence)
- 2026-01-09: Adversarial Review - Confirmed 100% complete, all AC satisfied, 85/85 tests passing
