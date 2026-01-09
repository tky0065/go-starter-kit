# Story 2.4: Sécurisation des routes (Auth Middleware)

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a **développeur**,
I want **protéger certaines routes API via un middleware**,
so that **seuls les utilisateurs authentifiés puissent y accéder**.

## Acceptance Criteria

1. **Middleware d'authentification :**
    - Créer un middleware réutilisable basé sur `github.com/gofiber/contrib/jwt`.
    - Le middleware doit valider la signature JWT (`HS256`) en utilisant la clé secrète définie dans `JWT_SECRET` (env).
    - Le middleware doit extraire le token depuis le header `Authorization: Bearer <token>`.
2. **Gestion des Erreurs (Unauthorized) :**
    - Si le token est absent ou malformé -> Retourner `401 Unauthorized`.
    - Si le token est expiré -> Retourner `401 Unauthorized` (avec message "Token expired" si possible, mais standardisé JSON).
    - Si la signature est invalide -> Retourner `401 Unauthorized`.
    - **Format de réponse :** Doit respecter le format d'erreur standard du projet (`{"status": "error", ...}`).
3. **Contexte Utilisateur :**
    - En cas de succès, le middleware doit injecter les claims du token (notamment `user_id` / `sub`) dans le contexte local de Fiber (`c.Locals("user")` ou équivalent).
    - Permettre aux handlers suivants d'accéder facilement à l'ID de l'utilisateur connecté.
4. **Protection des Routes :**
    - Appliquer ce middleware au groupe de routes `/api/v1` globalement ou sélectivement (ex: `/api/v1/users/me` doit être protégé).
    - Les routes publiques (Login, Register) doivent rester accessibles sans token.

## Tasks / Subtasks

- [x] Créer le package middleware dans `internal/adapters/middleware/auth_middleware.go` (AC: 1, 2, 3)
    - [x] Fonction `NewAuthMiddleware(secret string) fiber.Handler`
    - [x] Configurer `jwtware.Config` avec `SigningKey`, `SigningMethod` ("HS256")
    - [x] Implémenter `ErrorHandler` pour retourner le JSON standard 401
    - [x] Configurer `SuccessHandler` (optionnel) pour mapper les claims vers un struct contextuel si nécessaire
- [x] Mettre à jour `internal/infrastructure/server/server.go` (AC: 4)
    - [x] Injecter le middleware via `fx`
    - [x] Appliquer le middleware aux routes protégées (stratégie : soit middleware global avec `Filter` pour routes publiques, soit application spécifique sur les groupes `protected`)
    - [x] *Recommandation :* Créer deux groupes dans `RegisterHandlers` : `public` et `protected`
- [x] Tester la protection (AC: 2, 3)
    - [x] Test d'intégration : Appel sans token -> 401
    - [x] Test d'intégration : Appel avec token invalide -> 401
    - [x] Test d'intégration : Appel avec token valide -> 200 + Accès UserID

## Dev Notes

### Architecture & Constraints
- **Library :** Utiliser `github.com/gofiber/contrib/jwt` (le wrapper officiel Fiber pour `golang-jwt/v5`).
- **Configuration :** Le secret doit provenir de la config (chargée via `dotenv`). Ne jamais hardcoder.
- **Dependency Injection :** Le middleware doit être fourni comme un composant `fx` ou instancié dans le module `server`.

### Technical Guidelines
- **Context Key :** Par défaut, `contrib/jwt` stocke le token dans `c.Locals("user")`. Assurez-vous de documenter ou d'helper pour extraire l'ID (ex: `func GetUserID(c *fiber.Ctx) (uint, error)`).
- **Security :** Utiliser `jwtware.Config{ ... }` avec des gestionnaires d'erreurs stricts.

### Project Structure Notes
- Placer le middleware dans `internal/adapters/middleware`.
- Si vous créez un helper `GetUserID`, placez-le dans le même package ou dans `pkg/auth`.

### References
- [Epic 2: Authentication & Security Foundation](_bmad-output/planning-artifacts/epics.md)
- [Architecture Decision Document](_bmad-output/planning-artifacts/architecture.md#authentication--security)
- [Official Fiber JWT Middleware Docs](https://docs.gofiber.io/contrib/jwt)

## Dev Agent Record

### Agent Model Used
Claude Sonnet 4.5

### Debug Log References
None

### Completion Notes List
- ✅ Ajouté la dépendance `github.com/gofiber/contrib/jwt` v1.1.2
- ✅ Créé le middleware d'authentification dans `internal/adapters/middleware/auth_middleware.go`
- ✅ Implémenté le helper `GetUserID` dans `pkg/auth/context.go` pour extraire l'ID utilisateur du contexte
- ✅ Mis à jour `internal/adapters/handlers/module.go` pour créer des groupes de routes publiques et protégées
- ✅ Ajouté la route protégée `/api/v1/users/me` pour tester le middleware
- ✅ Créé des tests d'intégration complets couvrant tous les scénarios (6 tests dans auth_middleware_test.go, 4 tests dans protected_routes_test.go)
- ✅ Tous les tests passent (100% success rate)
- ✅ Le linting passe sans erreur
- ✅ Les routes publiques (/auth/register, /auth/login, /auth/refresh) restent accessibles sans token
- ✅ Le middleware retourne le format d'erreur standard `{"status": "error", ...}` pour les 401

### Implementation Plan
**Approche retenue :** Création de deux groupes de routes distincts (public et protected) dans le module handlers, avec injection du middleware via fx.

**Détails techniques :**
- Le middleware utilise `jwtware.Config` avec `SigningKey` pour valider les tokens JWT HS256
- L'ErrorHandler personnalisé retourne le format JSON standard du projet
- Le token est automatiquement stocké dans `c.Locals("user")` par le middleware
- Le helper `GetUserID` extrait le `sub` claim et le parse en uint

### File List
- manual-test-project/internal/adapters/middleware/auth_middleware.go (NEW)
- manual-test-project/internal/adapters/middleware/auth_middleware_test.go (NEW)
- manual-test-project/internal/adapters/handlers/module.go (MODIFIED)
- manual-test-project/internal/adapters/handlers/auth_handler.go (MODIFIED - added GetCurrentUser)
- manual-test-project/internal/adapters/handlers/protected_routes_test.go (NEW)
- manual-test-project/internal/adapters/handlers/auth_handler_login_test.go (MODIFIED - removed RegisterRoutes calls)
- manual-test-project/internal/adapters/handlers/auth_handler_integration_test.go (MODIFIED - removed RegisterRoutes calls)
- manual-test-project/internal/infrastructure/server/server.go (MODIFIED - added middleware import)
- manual-test-project/pkg/auth/context.go (NEW)
- manual-test-project/go.mod (MODIFIED - added gofiber/contrib/jwt dependency)
- manual-test-project/go.sum (MODIFIED)

## Senior Developer Review (AI)

**Date: 2026-01-09**
- **Issue Fixed:** Explicitly added `SigningMethod: "HS256"` to `NewAuthMiddleware` configuration to enforce secure signing algorithm, mitigating potential "none" algorithm attacks.
- **Documentation Update:** Corrected file paths in the "File List" to accurately reflect the `manual-test-project` directory structure.
- **Verification:** Middleware logic verified correct, and tests provide excellent coverage for both success and error scenarios.

## Change Log

**Date: 2026-01-09**
- Implémenté le middleware d'authentification JWT utilisant `gofiber/contrib/jwt`
- Créé le système de groupes de routes (public/protected) dans le module handlers
- Ajouté le helper `GetUserID` pour extraire l'ID utilisateur du contexte
- Implémenté la route protégée `/api/v1/users/me` comme exemple
- Créé 10 tests d'intégration couvrant tous les cas d'usage du middleware
- Tous les acceptance criteria sont satisfaits et validés par les tests
- Review & Fix: Ajout explicite de HS256 et correction des chemins de fichiers

