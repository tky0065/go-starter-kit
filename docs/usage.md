# Guide d'utilisation

Ce guide détaille l'utilisation de `create-go-starter` et explique en profondeur la structure des projets générés.

## Commande de base

La syntaxe de base est très simple:

```bash
create-go-starter <nom-du-projet>
```

### Exemple

```bash
create-go-starter mon-api-backend
```

Cette commande va créer un nouveau répertoire `mon-api-backend/` avec toute la structure du projet.

## Options disponibles

### Flags actuels

```bash
create-go-starter --help    # Afficher l'aide
create-go-starter -h        # Alias pour --help
```

> **Note**: D'autres options seront ajoutées dans les futures versions (choix de templates, bases de données, frameworks, etc.)

## Conventions de nommage

Le nom du projet doit respecter certaines règles:

### Caractères autorisés

- **Lettres**: a-z, A-Z
- **Chiffres**: 0-9
- **Tirets**: - (hyphen)
- **Underscores**: _ (trait de soulignement)

### Restrictions

- Pas d'espaces
- Pas de caractères spéciaux (/, \, @, #, etc.)
- Pas de points (.)
- Doit commencer par une lettre ou un chiffre (pas de tiret au début)

### Exemples valides

```bash
create-go-starter mon-projet           ✅
create-go-starter my-awesome-api       ✅
create-go-starter user_service         ✅
create-go-starter app2024              ✅
create-go-starter MonProjet            ✅
```

### Exemples invalides

```bash
create-go-starter mon projet           ❌ (contient un espace)
create-go-starter mon-projet!          ❌ (caractère spécial)
create-go-starter my.project           ❌ (contient un point)
create-go-starter -mon-projet          ❌ (commence par un tiret)
create-go-starter mon/projet           ❌ (contient un slash)
```

## Structure générée

Voici la structure complète créée par `create-go-starter`:

```
mon-projet/
├── cmd/
│   └── main.go                              # Point d'entrée de l'application
│
├── internal/
│   ├── models/
│   │   └── user.go                          # Entités partagées: User, RefreshToken, AuthResponse
│   │
│   ├── adapters/
│   │   ├── handlers/
│   │   │   ├── auth_handler.go              # Endpoints auth (register, login, refresh)
│   │   │   ├── auth_handler_test.go         # Tests handlers auth
│   │   │   ├── user_handler.go              # Endpoints CRUD users
│   │   │   └── user_handler_test.go         # Tests handlers users
│   │   ├── middleware/
│   │   │   ├── auth_middleware.go           # Middleware JWT authentication
│   │   │   ├── auth_middleware_test.go      # Tests middleware auth
│   │   │   ├── error_handler.go             # Middleware gestion centralisée erreurs
│   │   │   └── error_handler_test.go        # Tests error handler
│   │   ├── repository/
│   │   │   ├── user_repository.go           # Implémentation GORM du repository
│   │   │   └── user_repository_test.go      # Tests repository
│   │   └── http/
│   │       ├── health.go                    # Handler health check
│   │       └── health_test.go               # Tests health check
│   │
│   ├── domain/
│   │   ├── errors.go                        # Erreurs métier personnalisées
│   │   ├── errors_test.go                   # Tests erreurs
│   │   └── user/
│   │       ├── service.go                   # Logique métier (Register, Login, etc.)
│   │       ├── service_test.go              # Tests service
│   │       └── module.go                    # Module fx pour dependency injection
│   │
│   ├── infrastructure/
│   │   ├── database/
│   │   │   ├── database.go                  # Configuration GORM et connexion DB
│   │   │   ├── database_test.go             # Tests database
│   │   │   └── module.go                    # Module fx pour database
│   │   └── server/
│   │       ├── server.go                    # Configuration Fiber et routes
│   │       ├── server_test.go               # Tests server
│   │       └── module.go                    # Module fx pour server
│   │
│   └── interfaces/
│       ├── auth_service.go                  # Interface AuthService (port)
│       ├── user_service.go                  # Interface UserService (port)
│       └── user_repository.go               # Interface UserRepository (port)
│
├── pkg/
│   ├── auth/
│   │   ├── jwt.go                           # Génération et parsing JWT tokens
│   │   ├── jwt_test.go                      # Tests JWT
│   │   ├── middleware.go                    # Middleware JWT pour Fiber
│   │   ├── middleware_test.go               # Tests middleware
│   │   └── module.go                        # Module fx pour auth
│   ├── config/
│   │   ├── env.go                           # Chargement variables d'environnement
│   │   ├── env_test.go                      # Tests config
│   │   └── module.go                        # Module fx pour config
│   └── logger/
│       ├── logger.go                        # Configuration zerolog
│       ├── logger_test.go                   # Tests logger
│       └── module.go                        # Module fx pour logger
│
├── .github/
│   └── workflows/
│       └── ci.yml                           # Pipeline CI/CD GitHub Actions
│
├── .env                                     # Variables d'environnement (créé automatiquement)
├── .env.example                             # Template de configuration
├── .gitignore                               # Exclusions Git
├── .golangci.yml                            # Configuration golangci-lint
├── Dockerfile                               # Build Docker multi-stage
├── Makefile                                 # Commandes utiles (run, test, lint, etc.)
├── setup.sh                                 # Script de configuration automatique
├── go.mod                                   # Module Go et dépendances
└── README.md                                # Documentation du projet
```

**Total**: ~45+ fichiers créés automatiquement!

## Explication détaillée de chaque répertoire

### `/cmd`

**Rôle**: Point d'entrée de l'application.

**Contenu**:
- `main.go`: Bootstrap de l'application avec uber-go/fx
  - Initialise tous les modules (database, server, logger, etc.)
  - Configure le lifecycle (OnStart, OnStop)
  - Lance l'application avec `fx.New().Run()`

**Pattern**: Un seul fichier minimal qui orchestre les dépendances.

### `/internal/models`

**Rôle**: Entités de domaine partagées (domain entities) utilisées à travers toute l'application.

**Contenu**:
- `user.go`: Définit toutes les entités liées aux utilisateurs
  - `User`: Entité principale avec tags GORM (ID, Email, PasswordHash, timestamps)
  - `RefreshToken`: Entité pour la gestion des refresh tokens JWT
  - `AuthResponse`: DTO pour les réponses d'authentification

**Pourquoi un package séparé?**
- **Évite les dépendances circulaires**: Avant, `interfaces` → `domain/user` → `interfaces` créait un cycle
- **Maintenant**: `interfaces` → `models` ← `domain/user` (pas de cycle!)
- **Centralisation**: Les entités sont définies en un seul endroit
- **Réutilisabilité**: Tous les layers (domain, interfaces, adapters) peuvent importer models sans conflit

**Import**:
```go
import "mon-projet/internal/models"

user := &models.User{
    Email: "user@example.com",
    PasswordHash: hashedPassword,
}
```

### `/internal/domain`

**Rôle**: Couche métier (logique business), cœur de l'architecture hexagonale.

**Contenu**:
- `errors.go`: Définition des erreurs métier (DomainError, NotFoundError, ValidationError, etc.)
- `user/`: Domaine User
  - `service.go`: Logique métier (Register, Login, GetUserByID, UpdateUser, DeleteUser)
  - `module.go`: Module fx qui fournit le service

**Principe**: Le domaine ne doit **jamais** importer d'autres packages (adapters, infrastructure). Les dépendances sont inversées via interfaces (ports). Les entités sont maintenant dans `internal/models` pour éviter les cycles de dépendances.

### `/internal/adapters`

**Rôle**: Adaptateurs qui connectent le domaine au monde extérieur (HTTP, DB).

#### `/internal/adapters/handlers`

**Rôle**: Handlers HTTP qui exposent l'API REST.

**Contenu**:
- `auth_handler.go`:
  - `Register`: POST /api/v1/auth/register
  - `Login`: POST /api/v1/auth/login
  - `RefreshToken`: POST /api/v1/auth/refresh
- `user_handler.go`:
  - `List`: GET /api/v1/users
  - `GetByID`: GET /api/v1/users/:id
  - `Update`: PUT /api/v1/users/:id
  - `Delete`: DELETE /api/v1/users/:id

**Pattern**: Handlers parsent les requêtes, valident, appellent les services du domaine, retournent les réponses.

#### `/internal/adapters/middleware`

**Rôle**: Middleware Fiber pour cross-cutting concerns.

**Contenu**:
- `auth_middleware.go`: Vérifie le JWT token dans les requêtes
- `error_handler.go`: Gestion centralisée des erreurs (convertit DomainError en réponses HTTP)

#### `/internal/adapters/repository`

**Rôle**: Implémentation du pattern Repository avec GORM.

**Contenu**:
- `user_repository.go`: Implémentation de l'interface UserRepository
  - Create, FindByID, FindByEmail, Update, Delete
  - Utilise GORM pour les opérations DB

**Pattern**: Repository isole le domaine de la couche de persistance.

#### `/internal/adapters/http`

**Rôle**: Handlers HTTP simples (non liés au domaine).

**Contenu**:
- `health.go`: Endpoint GET /health pour monitoring

### `/internal/infrastructure`

**Rôle**: Configuration de l'infrastructure (DB, serveur).

#### `/internal/infrastructure/database`

**Rôle**: Configuration et connexion à la base de données.

**Contenu**:
- `database.go`:
  - Connexion PostgreSQL via GORM
  - Configuration du pool de connexions
  - AutoMigrate des entités (`models.User`, `models.RefreshToken`)
  - Gestion du lifecycle (fermeture connexion)

#### `/internal/infrastructure/server`

**Rôle**: Configuration du serveur HTTP Fiber.

**Contenu**:
- `server.go`:
  - Configuration Fiber app
  - Enregistrement des routes
  - Middleware error handler
  - Groupes de routes (/api/v1/auth, /api/v1/users)
  - Lifecycle du serveur (start, graceful shutdown)

### `/internal/interfaces`

**Rôle**: Définition des ports (interfaces) pour l'architecture hexagonale.

**Contenu**:
- `user_repository.go`: Interface pour la persistance utilisateur

**Principe**:
- Les adapters dépendent de ces interfaces, pas des implémentations concrètes
- Les interfaces référencent `models.*` pour les types de retour/paramètres
- Exemple: `CreateUser(ctx context.Context, user *models.User) error`

### `/pkg`

**Rôle**: Packages réutilisables, peuvent être importés par d'autres projets.

#### `/pkg/auth`

**Rôle**: Utilitaires JWT.

**Contenu**:
- `jwt.go`: Génération et validation de tokens JWT
  - GenerateAccessToken (15min)
  - GenerateRefreshToken (7 jours)
  - ParseToken
- `middleware.go`: Middleware Fiber pour JWT

#### `/pkg/config`

**Rôle**: Chargement de la configuration.

**Contenu**:
- `env.go`: Charge les variables .env (godotenv) et les expose via struct Config

#### `/pkg/logger`

**Rôle**: Configuration du logger.

**Contenu**:
- `logger.go`: Configure zerolog (niveau, format, output)

### `/.github/workflows`

**Rôle**: CI/CD avec GitHub Actions.

**Contenu**:
- `ci.yml`: Pipeline qui exécute:
  1. golangci-lint (quality checks)
  2. Tests avec PostgreSQL (service container)
  3. Build verification

## Fichiers de configuration détaillés

### `.env` et `.env.example`

Le fichier `.env.example` est un template, et `.env` est copié automatiquement lors de la génération.

**Variables**:

```bash
# Application
APP_NAME=mon-projet                # Nom de l'app (utilisé dans logs)
APP_ENV=development                # Environnement (development, staging, production)
APP_PORT=8080                      # Port HTTP

# Database PostgreSQL
DB_HOST=localhost                  # Hôte de la DB
DB_PORT=5432                       # Port PostgreSQL
DB_USER=postgres                   # Utilisateur DB
DB_PASSWORD=postgres               # Mot de passe DB
DB_NAME=mon-projet                 # Nom de la base de données
DB_SSLMODE=disable                 # SSL mode (require pour production)

# JWT Authentication
JWT_SECRET=                        # SECRET CRITIQUE - À générer!
JWT_EXPIRY=15m                     # Durée des access tokens (15 minutes)
REFRESH_TOKEN_EXPIRY=168h          # Durée des refresh tokens (7 jours)
```

**Important**: Générez un JWT_SECRET sécurisé:

```bash
openssl rand -base64 32
```

Puis ajoutez-le dans `.env`:

```bash
JWT_SECRET=votre_secret_genere_ici
```

### `go.mod`

Définit le module Go et les dépendances:

```go
module mon-projet

go 1.25

require (
    github.com/gofiber/fiber/v2 v2.x.x
    github.com/golang-jwt/jwt/v5 v5.x.x
    gorm.io/gorm v1.x.x
    gorm.io/driver/postgres v1.x.x
    go.uber.org/fx v1.x.x
    github.com/rs/zerolog v1.x.x
    github.com/go-playground/validator/v10 v10.x.x
    golang.org/x/crypto v0.x.x
    github.com/joho/godotenv v1.x.x
    // ... autres dépendances
)
```

### `Makefile`

Commandes utiles pour le développement:

```makefile
.PHONY: help run build test test-coverage lint clean docker-build docker-run

help:           # Afficher toutes les commandes
run:            # Lancer l'application (go run)
build:          # Compiler le binaire
test:           # Exécuter les tests avec race detector
test-coverage:  # Tests avec rapport HTML de coverage
lint:           # Linter avec golangci-lint
clean:          # Nettoyer les artifacts
docker-build:   # Build l'image Docker
docker-run:     # Lancer le conteneur Docker
```

### `.golangci.yml`

Configuration du linter avec règles strictes:

- errcheck, gosimple, govet, ineffassign, staticcheck
- gofmt, goimports
- misspell, revive
- Exclusions pour tests et code généré

### `Dockerfile`

Build multi-stage optimisé:

**Stage 1 (builder)**:
- Image golang:1.25-alpine
- Build du binaire statique

**Stage 2 (runtime)**:
- Image alpine:latest (légère)
- Copie du binaire seulement
- Exécution en tant que non-root
- EXPOSE 8080

**Taille finale**: ~15-20MB (vs ~1GB avec image golang complète)

## Workflow après génération

Une fois le projet créé, vous avez deux options pour configurer votre projet:

### Option A: Configuration automatique avec setup.sh (Recommandé) 🚀

Le CLI génère automatiquement un script `setup.sh` qui automatise toute la configuration initiale.

**Fonctionnalités du script**:
- ✅ Vérification des prérequis (Go, OpenSSL, Docker)
- ✅ Installation des dépendances Go (`go mod tidy`)
- ✅ Génération automatique du JWT secret
- ✅ Configuration de PostgreSQL (Docker ou local)
- ✅ Exécution des tests
- ✅ Vérification de l'installation

**Utilisation**:

```bash
cd mon-projet
./setup.sh
```

Le script est **interactif** et vous guidera à travers les choix:
- Docker ou PostgreSQL local
- Régénération du JWT secret si déjà configuré
- Validation à chaque étape

**Après l'exécution du script**:

```bash
make run
```

C'est tout! Votre application est prête.

---

### Option B: Configuration manuelle

Si vous préférez configurer manuellement ou si le script setup.sh échoue, suivez ces étapes:

### Étape 1: Naviguer dans le projet

```bash
cd mon-projet
```

### Étape 2: Configurer les variables d'environnement

```bash
# Générer un JWT secret
openssl rand -base64 32

# Éditer .env et ajouter le secret
nano .env  # ou vim, code, etc.
```

Ajoutez:
```
JWT_SECRET=le_secret_genere
```

### Étape 3: Installer les dépendances

```bash
go mod tidy
```

### Étape 4: Lancer PostgreSQL

**Option A: Docker (recommandé)**

```bash
docker run -d \
  --name postgres \
  -e POSTGRES_DB=mon-projet \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine
```

**Option B: Installation locale**

```bash
# macOS
brew install postgresql
brew services start postgresql
createdb mon-projet

# Linux
sudo apt install postgresql
sudo systemctl start postgresql
sudo -u postgres createdb mon-projet
```

### Étape 5: Lancer l'application

```bash
make run
```

Ou directement:

```bash
go run cmd/main.go
```

Vous devriez voir:

```
{"level":"info","time":"...","message":"Starting server on :8080"}
{"level":"info","time":"...","message":"Database connected successfully"}
```

### Étape 6: Tester l'API

```bash
# Health check
curl http://localhost:8080/health
# {"status":"ok"}

# Register un utilisateur
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

## Commandes Make disponibles

Le Makefile généré inclut ces commandes:

| Commande | Description | Usage |
|----------|-------------|-------|
| `make help` | Afficher l'aide | Voir toutes les commandes |
| `make run` | Lancer l'app | Développement local |
| `make build` | Compiler le binaire | Créer l'exécutable |
| `make test` | Tests avec race detector | Vérifier le code |
| `make test-coverage` | Tests + rapport HTML | Voir le coverage |
| `make lint` | golangci-lint | Vérifier la qualité |
| `make clean` | Nettoyer | Supprimer artifacts |
| `make docker-build` | Build image Docker | Containerisation |
| `make docker-run` | Lancer conteneur | Test Docker |

**Exemples**:

```bash
# Développement quotidien
make run          # Lance l'app avec hot-reload (si air installé)

# Avant commit
make test         # Vérifie que tous les tests passent
make lint         # Vérifie la qualité du code

# Build pour production
make build        # Crée le binaire
./mon-projet      # Exécute le binaire

# Docker
make docker-build # Build l'image
make docker-run   # Teste le conteneur
```

## Prochaines étapes

Maintenant que vous comprenez la structure, consultez:

1. **[Guide des projets générés](./generated-project-guide.md)** - Guide complet pour:
   - Comprendre l'architecture hexagonale
   - Développer de nouvelles fonctionnalités
   - Utiliser l'API (endpoints, authentification)
   - Écrire des tests
   - Déployer en production

2. **[Architecture du CLI](./cli-architecture.md)** - Si vous voulez:
   - Comprendre comment fonctionne le générateur
   - Contribuer au projet
   - Étendre les templates

3. **Commencer à développer**:
   ```bash
   # Lancer l'app
   make run

   # Dans un autre terminal, tester l'API
   curl http://localhost:8080/health

   # Lire le code généré
   cat internal/domain/user/service.go
   cat internal/adapters/handlers/auth_handler.go
   ```

## Conseils

- **Lisez le code généré**: Chaque fichier est un exemple de best practice Go
- **Modifiez selon vos besoins**: La structure est un point de départ, adaptez-la
- **Suivez les patterns**: Repository, dependency injection, error handling centralisé
- **Testez régulièrement**: `make test` avant chaque commit
- **Utilisez le linter**: `make lint` pour maintenir la qualité

Bon développement! 🚀
