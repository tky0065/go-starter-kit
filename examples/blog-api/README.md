# Exemple: Blog API

Projet exemple complet généré avec `create-go-starter` et étendu avec des fonctionnalités de blog.

## Vue d'ensemble

Ce projet démontre comment utiliser `create-go-starter` pour créer une API Blog production-ready avec:

- **Authentification JWT** - Register, Login, Refresh tokens
- **Articles (Posts)** - CRUD complet avec slug auto-généré, tags, et statut publish/unpublish
- **Commentaires** - CRUD sur les articles
- **Relations** - Post → Author (User), Comment → Post + Author
- **Pagination** - Limit/Offset pour les listes
- **Tests** - Tests unitaires avec mocks
- **Docker** - Configuration docker-compose prête

## Fonctionnalités par domaine

### User (pré-généré par create-go-starter)
- Register avec email/password
- Login avec JWT access + refresh tokens
- Refresh token rotation
- CRUD utilisateurs

### Post (Article)
- ✅ Créer un article (titre, contenu, tags)
- ✅ Slug auto-généré depuis le titre (ex: "Mon Article" → "mon-article")
- ✅ Lister les articles avec pagination
- ✅ Lister les articles par auteur
- ✅ Récupérer un article par ID ou slug
- ✅ Mettre à jour un article
- ✅ Publier/Dépublier un article
- ✅ Supprimer un article (soft delete)

### Comment (Commentaire)
- ✅ Ajouter un commentaire sur un article
- ✅ Lister les commentaires d'un article
- ✅ Supprimer un commentaire

## Prérequis

- **Go 1.25+**
- **PostgreSQL** ou **Docker**
- **curl** ou **Postman** pour tester

## Installation

### 1. Cloner le repository

```bash
git clone https://github.com/tky0065/go-starter-kit.git
cd go-starter-kit/examples/blog-api
```

### 2. Installer les dépendances

```bash
go mod tidy
```

### 3. Configurer l'environnement

Copier le fichier d'exemple:

```bash
cp .env.example .env
```

Générer un JWT secret:

```bash
openssl rand -base64 32
```

Éditer `.env` et ajouter le JWT_SECRET:

```env
JWT_SECRET=<coller_le_secret_ici>
```

### 4. Lancer la base de données

#### Option A: Docker

```bash
docker run -d \
  --name blog-postgres \
  -e POSTGRES_DB=blog_api \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine
```

#### Option B: PostgreSQL local

```bash
createdb blog_api
```

### 5. Lancer l'application

```bash
make run
```

L'API est disponible sur `http://localhost:8080`

### 6. Vérifier le fonctionnement

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

## Utilisation de l'API

### Authentification

#### Créer un compte

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

Sauvegarder le `access_token` retourné.

#### Se connecter

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Articles (Posts)

**Définir le token**:

```bash
export TOKEN="<access_token_from_login>"
```

#### Créer un article

```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Mon Premier Article",
    "content": "Ceci est le contenu de mon article",
    "tags": "golang,tutorial,blog"
  }'
```

Réponse:
```json
{
  "id": 1,
  "title": "Mon Premier Article",
  "slug": "mon-premier-article",
  "content": "Ceci est le contenu de mon article",
  "tags": "golang,tutorial,blog",
  "published": false,
  "author_id": 1,
  "created_at": "2024-01-10T10:00:00Z"
}
```

#### Lister les articles

```bash
curl http://localhost:8080/api/v1/posts?limit=10&offset=0
```

#### Récupérer un article par slug

```bash
curl http://localhost:8080/api/v1/posts/mon-premier-article
```

#### Publier un article

```bash
curl -X POST http://localhost:8080/api/v1/posts/1/publish \
  -H "Authorization: Bearer $TOKEN"
```

#### Mettre à jour un article

```bash
curl -X PUT http://localhost:8080/api/v1/posts/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Titre Modifié",
    "content": "Contenu mis à jour"
  }'
```

#### Supprimer un article

```bash
curl -X DELETE http://localhost:8080/api/v1/posts/1 \
  -H "Authorization: Bearer $TOKEN"
```

### Commentaires

#### Ajouter un commentaire

```bash
curl -X POST http://localhost:8080/api/v1/comments \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "post_id": 1,
    "content": "Excellent article!"
  }'
```

#### Lister les commentaires d'un article

```bash
curl http://localhost:8080/api/v1/comments/post/1
```

#### Supprimer un commentaire

```bash
curl -X DELETE http://localhost:8080/api/v1/comments/1 \
  -H "Authorization: Bearer $TOKEN"
```

## Structure du projet

```
blog-api/
├── cmd/
│   └── main.go                       # Point d'entrée avec fx DI
├── internal/
│   ├── domain/
│   │   ├── user/                     # Domaine User (authentification)
│   │   ├── post/                     # Domaine Post (articles)
│   │   │   ├── entity.go
│   │   │   ├── service.go
│   │   │   └── module.go
│   │   ├── comment/                  # Domaine Comment
│   │   │   ├── entity.go
│   │   │   ├── service.go
│   │   │   └── module.go
│   │   └── errors.go
│   ├── adapters/
│   │   ├── handlers/
│   │   │   ├── auth_handler.go
│   │   │   ├── user_handler.go
│   │   │   ├── post_handler.go       # Handler HTTP pour Posts
│   │   │   └── comment_handler.go    # Handler HTTP pour Comments
│   │   ├── middleware/
│   │   │   ├── auth_middleware.go
│   │   │   └── error_handler.go
│   │   └── repository/
│   │       ├── user_repository.go
│   │       ├── post_repository.go    # Repository GORM pour Posts
│   │       └── comment_repository.go # Repository GORM pour Comments
│   ├── infrastructure/
│   │   ├── database/
│   │   │   ├── database.go
│   │   │   └── migrations.go         # Migrations Auto (User, Post, Comment)
│   │   └── server/
│   │       ├── server.go
│   │       └── routes.go             # Routes pour tous les domaines
│   └── interfaces/
│       ├── user_service.go
│       ├── user_repository.go
│       ├── post_service.go           # Interface PostService
│       ├── post_repository.go        # Interface PostRepository
│       ├── comment_service.go
│       └── comment_repository.go
├── pkg/
│   ├── auth/                         # JWT utilities
│   ├── config/                       # Configuration
│   └── logger/                       # Zerolog logger
├── .env.example
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── go.mod
└── README.md
```

## Architecture

### Architecture hexagonale (Ports & Adapters)

```
┌─────────────────────────────────────────────────────────┐
│                    HTTP Handlers                        │
│           (adapters/handlers/post_handler.go)           │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ↓
┌─────────────────────────────────────────────────────────┐
│                  Domain Service                         │
│              (domain/post/service.go)                   │
│  - Business Logic (publish/unpublish, slugify)          │
│  - Validation                                           │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ↓
┌─────────────────────────────────────────────────────────┐
│                   Repository                            │
│           (adapters/repository/post_repository.go)      │
│  - GORM queries (FindAll, Create, Update, Delete)       │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ↓
              PostgreSQL Database
```

### Relations entre entités

```
User (1) ────────< (N) Post
                       │
                       │
                       └────< (N) Comment
```

- Un **User** peut avoir plusieurs **Posts** (author_id)
- Un **Post** peut avoir plusieurs **Comments** (post_id)
- Un **Comment** appartient à un **User** (author_id) et un **Post** (post_id)

## Tests

### Lancer tous les tests

```bash
make test
```

### Tests avec coverage

```bash
make test-coverage
```

### Linting

```bash
make lint
```

## Déploiement Docker

### Build l'image

```bash
make docker-build
```

### Lancer avec docker-compose

```bash
docker-compose up -d
```

Cela lance:
- L'application sur `http://localhost:8080`
- PostgreSQL sur le port 5432

### Vérifier

```bash
curl http://localhost:8080/health
docker-compose logs -f app
```

### Arrêter

```bash
docker-compose down
```

## Ce que vous pouvez apprendre de cet exemple

### 1. Architecture hexagonale

- **Domain** (`internal/domain/`): Logique métier pure, indépendante des frameworks
- **Adapters** (`internal/adapters/`): Implémentations concrètes (HTTP handlers, GORM repositories)
- **Interfaces** (`internal/interfaces/`): Contrats entre les couches

### 2. Dependency Injection avec fx

Le fichier `cmd/main.go` montre comment:
- Déclarer les modules fx
- Injecter les dépendances automatiquement
- Gérer le lifecycle (OnStart, OnStop)

### 3. Patterns GORM

- **Migrations**: AutoMigrate dans `migrations.go`
- **Relations**: Foreign keys (author_id, post_id)
- **Soft Delete**: Utilisation de `DeletedAt`
- **Hooks**: BeforeCreate pour générer le slug
- **Pagination**: Limit/Offset avec Count

### 4. Bonnes pratiques

- **Validation** des inputs avec go-playground/validator
- **Error handling** centralisé avec DomainError
- **Logging structuré** avec zerolog
- **JWT** avec access + refresh tokens
- **Middleware** pour l'authentification
- **Makefile** pour automatiser les tâches

## Aller plus loin

Extensions possibles:

- **Upload d'images** pour les articles
- **Recherche full-text** (PostgreSQL FTS)
- **Likes/Votes** sur les articles
- **Catégories** pour organiser les posts
- **Notifications** (emails, webhooks)
- **Rate limiting** pour protéger l'API
- **Swagger** pour documenter l'API
- **Elasticsearch** pour recherche avancée
- **Redis** pour caching
- **Websockets** pour commentaires en temps réel

## Ressources

- [Tutorial complet](../../docs/tutorial-exemple-complet.md) - Guide pas-à-pas pour créer ce projet
- [Guide des projets générés](../../docs/generated-project-guide.md) - Documentation complète
- [create-go-starter](https://github.com/tky0065/go-starter-kit) - CLI generator

## Licence

MIT - Libre d'utilisation pour projets personnels et commerciaux.

---

**Généré avec [create-go-starter](https://github.com/tky0065/go-starter-kit)** 🚀
