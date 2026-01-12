# Plan de Corrections - go-starter-kit

## ✅ PROBLÈME RÉSOLU (2026-01-12)

### Problème: Routes non enregistrées

**Symptôme**: Le projet compilait mais les routes API ne fonctionnaient pas car `RegisterRoutes()` n'était jamais invoqué par fx.

**Solution appliquée**:
1. Modifié `ServerTemplate()` dans `templates.go` pour ajouter:
   - Import `httpRoutes "<project>/internal/adapters/http"`
   - `fx.Invoke(httpRoutes.RegisterRoutes)` dans le module server

**Fichiers modifiés**:
- `cmd/create-go-starter/templates.go` - ServerTemplate avec fx.Invoke
- `cmd/create-go-starter/templates_test.go` - Tests mis à jour

**Tests de validation**:
- ✅ Tous les tests passent (`go test ./cmd/create-go-starter/...`)
- ✅ Test E2E passe (projet généré compile)
- ✅ `server.go` généré contient `fx.Invoke(httpRoutes.RegisterRoutes)`

---

## Historique des Corrections

### ✅ RÉSOLU - Conflit de types DI (Commit b5ee53f - 2026-01-10)

La dépendance circulaire entre `internal/interfaces` et `internal/domain/user` a été résolue en créant le package `internal/models` pour les entités partagées.

### ✅ RÉSOLU - Routes centralisées (Commit 5905c7b - 2026-01-12)

Refactoring des routes dans `http/routes.go` au lieu de `handlers/module.go`.

### ✅ RÉSOLU - RegisterRoutes non invoqué (2026-01-12)

Ajout de `fx.Invoke(httpRoutes.RegisterRoutes)` dans `ServerTemplate()`.

---

## Architecture Finale

```
internal/
├── adapters/
│   ├── handlers/
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   └── module.go          # fx.Provide handlers (PAS d'invoke routes)
│   ├── http/
│   │   ├── health.go
│   │   └── routes.go          # RegisterRoutes() - enregistre toutes les routes
│   ├── middleware/
│   │   └── error_handler.go
│   └── repository/
│       ├── user_repository.go
│       └── module.go
├── domain/
│   ├── user/
│   │   ├── service.go
│   │   └── module.go
│   └── errors.go
├── infrastructure/
│   ├── database/
│   │   └── database.go
│   └── server/
│       └── server.go          # fx.Invoke(httpRoutes.RegisterRoutes)
├── interfaces/
│   ├── services.go            # TokenService interface
│   └── user_repository.go     # UserRepository interface
└── models/
    └── user.go                # User, RefreshToken, AuthResponse
```

## Flux d'Injection de Dépendances (fx)

```
main.go
  └── fx.New(...)
       ├── logger.Module
       ├── database.Module     → *gorm.DB
       ├── auth.Module         → interfaces.TokenService, fiber.Handler (JWT middleware)
       ├── user.Module         → *user.Service
       ├── repository.Module   → interfaces.UserRepository
       ├── handlers.Module     → *AuthHandler, *UserHandler
       └── server.Module
            ├── fx.Provide(NewServer)       → *fiber.App
            ├── fx.Invoke(registerHooks)    → Start/Stop server
            └── fx.Invoke(RegisterRoutes)   → Enregistre les routes
```

**Clé**: `RegisterRoutes` est invoqué depuis `server.Module`, pas depuis `handlers.Module`, pour éviter les imports cycliques.
