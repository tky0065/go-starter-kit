# Story 2.4: Sécurisation des routes (Auth Middleware)

Status: ready-for-dev

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

- [ ] Créer le package middleware dans `internal/adapters/middleware/auth_middleware.go` (AC: 1, 2, 3)
    - [ ] Fonction `NewAuthMiddleware(secret string) fiber.Handler`
    - [ ] Configurer `jwtware.Config` avec `SigningKey`, `SigningMethod` ("HS256")
    - [ ] Implémenter `ErrorHandler` pour retourner le JSON standard 401
    - [ ] Configurer `SuccessHandler` (optionnel) pour mapper les claims vers un struct contextuel si nécessaire
- [ ] Mettre à jour `internal/infrastructure/server/server.go` (AC: 4)
    - [ ] Injecter le middleware via `fx`
    - [ ] Appliquer le middleware aux routes protégées (stratégie : soit middleware global avec `Filter` pour routes publiques, soit application spécifique sur les groupes `protected`)
    - [ ] *Recommandation :* Créer deux groupes dans `RegisterHandlers` : `public` et `protected`
- [ ] Tester la protection (AC: 2, 3)
    - [ ] Test d'intégration : Appel sans token -> 401
    - [ ] Test d'intégration : Appel avec token invalide -> 401
    - [ ] Test d'intégration : Appel avec token valide -> 200 + Accès UserID

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
Gemini 2.0 Flash

### Debug Log References
None

### Completion Notes List
- Confirmed usage of `gofiber/contrib/jwt`.
- Specified error handling requirements to match project standard.
- Defined integration strategy via `fx` and route grouping.

### File List
- internal/adapters/middleware/auth_middleware.go
- internal/infrastructure/server/server.go
- pkg/auth/context.go (Optional helper)
