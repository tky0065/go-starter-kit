# Plan de Corrections - go-starter-kit

## ✅ PROBLÈME RÉSOLU (Commit b5ee53f - 2026-01-10)

**Ce problème a été corrigé par l'introduction du package `internal/models`.**

### Solution implémentée

La dépendance circulaire entre `internal/interfaces` et `internal/domain/user` a été résolue en créant un nouveau package partagé:

**Changements architecturaux**:
- ✅ Création de `internal/models` pour les entités partagées (User, RefreshToken, AuthResponse)
- ✅ Les interfaces dans `internal/interfaces` référencent maintenant `models.*`
- ✅ Le service dans `internal/domain/user` utilise `models.*`
- ✅ Plus de dépendance circulaire: `interfaces` → `models` ← `domain/user`

**Tests de validation réussis**:
- ✅ Application démarre sans erreur fx
- ✅ Health check: `{"status":"ok"}`
- ✅ User registration fonctionne
- ✅ User login avec JWT fonctionne
- ✅ Endpoints protégés fonctionnent
- ✅ Accès non autorisé bloqué correctement

Voir le commit `b5ee53f` pour les détails de l'implémentation.

---

## Rapport Original - Problème Identifié (2026-01-10)

**Date**: 2026-01-10
**Projet testé**: test-api-project (généré avec create-go-starter)
**Statut Original**: ❌ **ÉCHEC - Application ne démarre pas**
**Statut Actuel**: ✅ **RÉSOLU**

## Problèmes Identifiés

### 🔴 Problème Critique #1: Conflit de types dans Dependency Injection

**Erreur rencontrée**:
```
[Fx] ERROR Failed to start: could not build arguments for function "test-api-project/internal/adapters/handlers".RegisterAllRoutes
missing types:
  - user.UserRepository (did you mean to Provide it?)
  - user.TokenService (did you mean to Provide it?)
```

**Cause racine**:
Le fichier `internal/domain/user/service.go` définit des interfaces **locales** au package `user`:
- Ligne 13: `type UserRepository interface { ... }`
- Ligne 27: `type TokenService interface { ... }`

Ces interfaces locales sont **différentes** des interfaces globales définies dans `internal/interfaces/`:
- `internal/interfaces/user_repository.go`: `type UserRepository interface { ... }`
- `internal/interfaces/token_service.go`: `type TokenService interface { ... }`

**Conflit fx**:
- fx fournit des implémentations de type `interfaces.UserRepository` et `interfaces.TokenService`
- Mais `NewServiceWithJWT()` attend `user.UserRepository` et `user.TokenService`
- Go les considère comme des types **incompatibles** car ils appartiennent à des packages différents

**Impact**: L'application ne peut pas démarrer car fx ne peut pas résoudre les dépendances.

---

### Étapes de Reproduction

1. Générer un projet:
   ```bash
   create-go-starter test-api-project
   cd test-api-project
   ```

2. Installer les dépendances:
   ```bash
   go mod tidy
   ```

3. Configurer l'environnement:
   ```bash
   # Ajouter JWT_SECRET dans .env
   JWT_SECRET=Zf5sjJWdsQL//AgFBatdw4gSR0PdTQCUmLK1NEyi0iA=
   ```

4. Démarrer PostgreSQL:
   ```bash
   docker run -d --name test-postgres \
     -e POSTGRES_DB=test-api-project \
     -e POSTGRES_PASSWORD=postgres \
     -p 5432:5432 \
     postgres:16-alpine
   ```

5. Lancer l'application:
   ```bash
   go build ./cmd/main.go
   ./main
   ```

6. **Résultat**: Erreur fx au démarrage (types manquants)

---

## Solutions Proposées

### ✅ Solution #1: Supprimer les interfaces locales et utiliser les interfaces globales

**Fichier à modifier**: `cmd/create-go-starter/templates.go`
**Template concerné**: `UserServiceTemplate()`

**Changements**:

#### Avant (code actuel - INCORRECT):
```go
// Dans internal/domain/user/service.go (généré)

package user

import (
	"context"
	"test-api-project/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository defines the contract for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id uint) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int) ([]*User, int64, error)
}

// TokenService defines the contract for token operations
type TokenService interface {
	GenerateAccessToken(userID uint, email string) (string, error)
	GenerateRefreshToken(userID uint) (*RefreshToken, error)
	ValidateRefreshToken(tokenString string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenID uint) error
}

// Service implements the user business logic
type Service struct {
	repo         UserRepository
	tokenService TokenService
}

// NewServiceWithJWT creates a new user service with JWT support
func NewServiceWithJWT(repo UserRepository, tokenService TokenService) *Service {
	return &Service{
		repo:         repo,
		tokenService: tokenService,
	}
}
```

#### Après (code corrigé - CORRECT):
```go
// Dans internal/domain/user/service.go (généré)

package user

import (
	"context"
	"test-api-project/internal/domain"
	"test-api-project/internal/interfaces"  // AJOUTÉ
	"golang.org/x/crypto/bcrypt"
)

// Service implements the user business logic
type Service struct {
	repo         interfaces.UserRepository     // CHANGÉ
	tokenService interfaces.TokenService       // CHANGÉ
}

// NewServiceWithJWT creates a new user service with JWT support
func NewServiceWithJWT(repo interfaces.UserRepository, tokenService interfaces.TokenService) *Service {  // CHANGÉ
	return &Service{
		repo:         repo,
		tokenService: tokenService,
	}
}
```

**IMPORTANT**: Supprimer complètement les définitions locales de `UserRepository` et `TokenService`.

---

### ✅ Solution #2: Vérifier que les modules fx fournissent les bons types

**Fichier à vérifier**: Modules fx dans `internal/domain/user/module.go`, `pkg/auth/module.go`, etc.

**Vérifier que**:
- `pkg/auth/module.go` fournit `interfaces.TokenService` (pas `auth.TokenService`)
- `internal/adapters/repository/module.go` fournit `interfaces.UserRepository` (pas `repository.UserRepository`)

**Exemple correct**:
```go
// Dans pkg/auth/module.go
var Module = fx.Module("auth",
	fx.Provide(
		fx.Annotate(
			NewTokenService,
			fx.As(new(interfaces.TokenService)),  // ✅ CORRECT: Cast vers l'interface globale
		),
		NewJWTMiddleware,
	),
)
```

---

## Fichiers du Starter à Corriger

### 1. **cmd/create-go-starter/templates.go**

**Fonction**: `UserServiceTemplate()`

**Ligne approximative**: ~1500-1700

**Modifications à apporter**:

1. Supprimer les définitions locales d'interfaces (UserRepository et TokenService)
2. Ajouter l'import `"{{.ProjectName}}/internal/interfaces"`
3. Utiliser `interfaces.UserRepository` et `interfaces.TokenService` partout

**Code à rechercher** (pour localiser la fonction):
```go
func (t *ProjectTemplates) UserServiceTemplate() string {
```

**Remplacement à effectuer dans le template**:

```go
// SUPPRIMER ces lignes:
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	// ...
}

type TokenService interface {
	GenerateAccessToken(userID uint, email string) (string, error)
	// ...
}

// AJOUTER dans les imports:
"` + t.projectName + `/internal/interfaces"

// REMPLACER dans la struct Service:
type Service struct {
	repo         UserRepository        // AVANT
	tokenService TokenService          // AVANT
}

// PAR:
type Service struct {
	repo         interfaces.UserRepository     // APRÈS
	tokenService interfaces.TokenService       // APRÈS
}

// REMPLACER la signature de NewServiceWithJWT:
func NewServiceWithJWT(repo UserRepository, tokenService TokenService) *Service {  // AVANT

// PAR:
func NewServiceWithJWT(repo interfaces.UserRepository, tokenService interfaces.TokenService) *Service {  // APRÈS
```

---

### 2. **Vérification des modules fx (templates.go)**

**Fonctions à vérifier**:
- `AuthModuleTemplate()` - S'assurer que TokenService est fourni avec `fx.As(new(interfaces.TokenService))`
- `RepositoryModuleTemplate()` - S'assurer que UserRepository est fourni avec `fx.As(new(interfaces.UserRepository))`

**Exemple de code correct**:
```go
// Dans AuthModuleTemplate()
var Module = fx.Module("auth",
	fx.Provide(
		fx.Annotate(
			NewTokenService,
			fx.As(new(interfaces.TokenService)),
		),
		NewJWTMiddleware,
	),
)
```

---

## Tests à Effectuer Après Correction

### 1. Regénérer un projet test

```bash
cd /tmp
rm -rf test-correction
mkdir test-correction
cd test-correction

create-go-starter test-fixed-project
cd test-fixed-project
```

### 2. Configuration et build

```bash
# Installer dépendances
go mod tidy

# Configurer .env
echo "JWT_SECRET=$(openssl rand -base64 32)" >> .env

# Build (doit compiler sans erreur)
go build ./cmd/main.go
```

### 3. Démarrer PostgreSQL

```bash
docker run -d --name test-fixed-postgres \
  -e POSTGRES_DB=test-fixed-project \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine

# Attendre que PostgreSQL soit prêt
sleep 5
```

### 4. Démarrer l'application

```bash
./main
```

**Résultat attendu**:
```
[Fx] PROVIDE	...
[Fx] RUN	...
INF Successfully connected to database
INF Database migrations completed successfully
INF Starting test-fixed-project server on :8080
```

**Critère de succès**: ✅ Aucune erreur fx, serveur démarre sur le port 8080

---

### 5. Tester les endpoints

#### Test 1: Health Check

```bash
curl http://localhost:8080/health
```

**Attendu**:
```json
{"status":"ok"}
```

---

#### Test 2: Register (Créer un utilisateur)

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@test.com",
    "password": "password123"
  }'
```

**Attendu** (code 200 ou 201):
```json
{
  "access_token": "eyJhbGci...",
  "refresh_token": "eyJhbGci...",
  "user": {
    "id": 1,
    "email": "user@test.com",
    "created_at": "2026-01-10T..."
  }
}
```

**Sauvegarder le access_token** pour les tests suivants.

---

#### Test 3: Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@test.com",
    "password": "password123"
  }'
```

**Attendu** (même format que register):
```json
{
  "access_token": "...",
  "refresh_token": "...",
  "user": { ... }
}
```

---

#### Test 4: List Users (Protected endpoint)

```bash
TOKEN="<access_token_from_register_or_login>"

curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer $TOKEN"
```

**Attendu**:
```json
{
  "data": [
    {
      "id": 1,
      "email": "user@test.com",
      "created_at": "..."
    }
  ],
  "total": 1,
  "limit": 10,
  "offset": 0
}
```

---

#### Test 5: Get User By ID

```bash
curl -X GET http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer $TOKEN"
```

**Attendu**:
```json
{
  "id": 1,
  "email": "user@test.com",
  "created_at": "..."
}
```

---

#### Test 6: Refresh Token

```bash
REFRESH_TOKEN="<refresh_token_from_register_or_login>"

curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\":\"$REFRESH_TOKEN\"}"
```

**Attendu**:
```json
{
  "access_token": "new_access_token...",
  "refresh_token": "new_refresh_token..."
}
```

---

## Checklist de Validation Complète

Après avoir effectué les corrections:

- [ ] ✅ Le projet se génère sans erreur
- [ ] ✅ `go mod tidy` réussit
- [ ] ✅ `go build` compile sans erreur
- [ ] ✅ L'application démarre avec PostgreSQL connecté
- [ ] ✅ Aucune erreur fx au démarrage
- [ ] ✅ Serveur écoute sur le port 8080
- [ ] ✅ GET `/health` retourne `{"status":"ok"}`
- [ ] ✅ POST `/api/v1/auth/register` crée un utilisateur et retourne des tokens
- [ ] ✅ POST `/api/v1/auth/login` authentifie et retourne des tokens
- [ ] ✅ GET `/api/v1/users` (avec token) retourne la liste des utilisateurs
- [ ] ✅ GET `/api/v1/users/:id` (avec token) retourne un utilisateur
- [ ] ✅ POST `/api/v1/auth/refresh` renouvelle les tokens
- [ ] ✅ Requêtes sans token vers endpoints protégés retournent 401 Unauthorized

---

## Priorité des Corrections

### 🔴 Priorité CRITIQUE (bloque le démarrage)

1. **Corriger le conflit de types UserRepository/TokenService**
   - Fichier: `cmd/create-go-starter/templates.go`
   - Fonction: `UserServiceTemplate()`
   - Impact: Sans cette correction, l'application ne peut pas démarrer

### 🟡 Priorité HAUTE (recommandé)

2. **Vérifier les modules fx**
   - S'assurer que tous les modules fx utilisent `fx.As(new(interfaces.XxxService))`
   - Impact: Prévient d'autres conflits de types similaires

### 🟢 Priorité NORMALE (documentation)

3. **Mettre à jour la documentation**
   - Ajouter une section de troubleshooting dans `docs/generated-project-guide.md`
   - Documenter les erreurs courantes fx et leurs solutions

---

## Notes Supplémentaires

### Pourquoi ce problème existe-t-il?

C'est une **erreur de conception** dans les templates générés. En architecture hexagonale:
- Les **interfaces (ports)** doivent être définies dans `internal/interfaces/`
- Les **implémentations** sont dans `internal/domain/`, `internal/adapters/`, etc.
- Les services du domaine doivent **dépendre des interfaces**, pas définir leurs propres interfaces locales

### Best practice Go

En Go, deux interfaces avec le même nom mais dans des packages différents sont des **types distincts**, même si elles ont les mêmes méthodes:
```go
package user
type UserRepository interface { ... }  // Type: user.UserRepository

package interfaces
type UserRepository interface { ... }  // Type: interfaces.UserRepository

// Ces deux types sont INCOMPATIBLES!
```

### Architecture correcte

```
internal/interfaces/          ← Définitions des interfaces (contrats)
    user_repository.go
    token_service.go

internal/domain/user/         ← Business logic (utilise les interfaces)
    service.go                   → Importe "internal/interfaces"

internal/adapters/repository/ ← Implémentations
    user_repository.go           → Implémente interfaces.UserRepository

pkg/auth/                     ← Implémentations
    token_service.go             → Implémente interfaces.TokenService
```

---

## Conclusion

**Problème principal**: Conflit de types causé par des définitions d'interfaces locales dans le package `user`

**Solution**: Supprimer les interfaces locales et utiliser les interfaces globales définies dans `internal/interfaces/`

**Impact de la correction**: Permettra à l'application de démarrer correctement et aux endpoints d'être testés

**Prochaine étape**: Appliquer les corrections dans `cmd/create-go-starter/templates.go`, regénérer un projet test, et valider tous les endpoints.
