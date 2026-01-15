# Tutorial: Créer une API Blog complète avec create-go-starter

Guide pas-à-pas pour créer une API Blog avec `create-go-starter`, de l'installation au déploiement.

## Table des matières

1. [Objectif](#objectif)
2. [Prérequis](#prérequis)
3. [Étape 1: Installation du CLI](#étape-1-installation-du-cli)
4. [Étape 2: Génération du projet](#étape-2-génération-du-projet)
5. [Étape 3: Configuration initiale](#étape-3-configuration-initiale)
6. [Étape 4: Tester le projet de base](#étape-4-tester-le-projet-de-base)
7. [Étape 5: Ajouter le domaine Post (Article)](#étape-5-ajouter-le-domaine-post-article)
8. [Étape 6: Implémenter le service Post](#étape-6-implémenter-le-service-post)
9. [Étape 7: Créer le repository Post](#étape-7-créer-le-repository-post)
10. [Étape 8: Créer le handler HTTP](#étape-8-créer-le-handler-http)
11. [Étape 9: Enregistrer les routes et le module](#étape-9-enregistrer-les-routes-et-le-module)
12. [Étape 10: Tester l'API Posts](#étape-10-tester-lapi-posts)
13. [Étape 11: Ajouter le domaine Comment](#étape-11-ajouter-le-domaine-comment)
14. [Étape 12: Tests unitaires](#étape-12-tests-unitaires)
15. [Étape 13: Déploiement Docker](#étape-13-déploiement-docker)
16. [Conclusion](#conclusion)

---

## Objectif

Créer une API REST complète pour un blog avec:

- **Articles (Posts)** avec auteur, titre, contenu, tags
- **Commentaires** sur les articles
- **Authentification JWT** (déjà incluse dans create-go-starter)
- **Tests complets**
- **Déploiement Docker**

À la fin de ce tutorial, vous aurez une API Blog production-ready avec toutes les bonnes pratiques.

## Prérequis

### Logiciels requis

- **Go 1.25+** - [Télécharger](https://golang.org/dl/)
- **PostgreSQL** ou **Docker** - Pour la base de données
- **curl** ou **Postman** - Pour tester l'API
- Éditeur de code (VS Code, GoLand, etc.)

### Connaissances recommandées

- Bases de Go (structs, interfaces, error handling)
- Concepts REST API
- Familiarité avec SQL/PostgreSQL (basique)

Pas besoin d'être expert! Ce tutorial explique chaque étape en détail.

---

## Étape 1: Installation du CLI

### Installation globale (recommandée)

La méthode la plus simple pour installer `create-go-starter`:

```bash
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@latest
```

Cette commande télécharge, compile et installe le CLI globalement.

### Vérification

```bash
create-go-starter --help
```

Vous devriez voir l'aide s'afficher.

**Note**: Si la commande n'est pas trouvée, ajoutez `$GOPATH/bin` à votre PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

---

## Étape 2: Génération du projet

### Créer le projet

```bash
create-go-starter blog-api
```

Cette commande génère **~45 fichiers** avec toute l'architecture nécessaire.

### Structure générée

```bash
cd blog-api
tree -L 3
```

**Résultat**:
```
blog-api/
├── cmd/
│   └── main.go                       # Point d'entrée avec fx DI
├── internal/
│   ├── models/
│   │   └── user.go                   # Entités: User, RefreshToken, AuthResponse
│   ├── domain/
│   │   ├── user/                     # Domaine User (pré-généré)
│   │   │   ├── service.go
│   │   │   └── module.go
│   │   └── errors.go
│   ├── adapters/
│   │   ├── handlers/
│   │   │   ├── auth_handler.go
│   │   │   └── user_handler.go
│   │   ├── middleware/
│   │   │   ├── auth_middleware.go
│   │   │   └── error_handler.go
│   │   └── repository/
│   │       └── user_repository.go
│   ├── infrastructure/
│   │   ├── database/
│   │   └── server/
│   └── interfaces/                   # Ports (interfaces)
│       └── user_repository.go
├── pkg/
│   ├── auth/                         # JWT utilities
│   ├── config/                       # Configuration
│   └── logger/                       # Zerolog logger
├── docs/
│   ├── README.md
│   └── quick-start.md
├── .env                              # Configuration (auto-copié)
├── .env.example
├── Dockerfile
├── Makefile
├── go.mod
└── README.md
```

:material-check-circle: **Checkpoint 1**: Le projet est généré avec succès.

---

## Étape 3: Configuration initiale

### 3.1 Installer les dépendances

```bash
cd blog-api
go mod tidy
```

Cette commande télécharge toutes les dépendances (Fiber, GORM, fx, etc.).

### 3.2 Configurer PostgreSQL

Vous avez 2 options:

#### Option A: Docker (recommandé)

```bash
docker run -d \
  --name blog-postgres \
  -e POSTGRES_DB=blog_api \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine
```

#### Option B: PostgreSQL local

Si PostgreSQL est installé localement:

```bash
createdb blog_api
```

### 3.3 Configurer les variables d'environnement

Générer un secret JWT sécurisé:

```bash
JWT_SECRET=$(openssl rand -base64 32)
echo "JWT_SECRET généré: $JWT_SECRET"
```

Éditer le fichier `.env`:

```bash
nano .env
```

Contenu du `.env`:

```env
# Application
APP_NAME=blog-api
APP_ENV=development
APP_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=blog_api
DB_SSLMODE=disable

# JWT
JWT_SECRET=<coller_le_secret_généré_ici>
JWT_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=168h
```

**Important**: Remplacez `<coller_le_secret_généré_ici>` par le JWT_SECRET généré.

---

## Étape 4: Tester le projet de base

### 4.1 Lancer l'application

```bash
make run
```

Vous devriez voir:

```
2024/01/10 10:00:00 INF Starting blog-api server on :8080
```

### 4.2 Tester le health check

Dans un autre terminal:

```bash
curl http://localhost:8080/health
```

**Réponse attendue**:
```json
{"status":"ok"}
```

### 4.3 Tester l'authentification par défaut

#### Créer un utilisateur

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@blog.com",
    "password": "admin123"
  }'
```

**Réponse**:
```json
{
  "access_token": "eyJhbGci...",
  "refresh_token": "eyJhbGci...",
  "user": {
    "id": 1,
    "email": "admin@blog.com",
    "created_at": "2024-01-10T10:05:00Z"
  }
}
```

#### Se connecter

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@blog.com",
    "password": "admin123"
  }'
```

**Même réponse** avec access_token et refresh_token.

#### Tester une route protégée

```bash
# Remplacez <ACCESS_TOKEN> par le token reçu
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

**Réponse**:
```json
[
  {
    "id": 1,
    "email": "admin@blog.com",
    "created_at": "2024-01-10T10:05:00Z"
  }
]
```

:material-check-circle: **Checkpoint 2**: Le projet de base fonctionne parfaitement avec User et Auth.

---

## Étape 5: Ajouter le domaine Post (Article)

Nous allons maintenant ajouter notre première fonctionnalité: les articles de blog.

### 5.1 Créer l'entité Post

Créer le fichier `internal/models/post.go`:

```go
package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

// Post représente un article de blog
type Post struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Contenu
	Title   string `gorm:"not null;size:255" json:"title" validate:"required,max=255"`
	Slug    string `gorm:"uniqueIndex;not null;size:255" json:"slug"`
	Content string `gorm:"type:text;not null" json:"content" validate:"required"`

	// Métadonnées
	Tags      string `gorm:"size:500" json:"tags"`
	Published bool   `gorm:"default:false" json:"published"`

	// Relations
	AuthorID uint `gorm:"not null" json:"author_id"`
}

// BeforeCreate génère automatiquement un slug unique avant l'insertion
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.Slug == "" {
		p.Slug = slugify(p.Title)
	}
	return nil
}

// slugify convertit un titre en slug URL-friendly
// Exemple: "Mon Super Article!" -> "mon-super-article"
func slugify(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")

	// Supprimer les caractères spéciaux
	replacer := strings.NewReplacer(
		"!", "", "?", "", ".", "", ",", "",
		"'", "", "\"", "", ":", "", ";", "",
		"(", "", ")", "", "[", "", "]", "",
	)
	slug = replacer.Replace(slug)

	// Supprimer les tirets multiples
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Supprimer les tirets en début/fin
	slug = strings.Trim(slug, "-")

	return slug
}
```

**Explications**:

- **struct Post**: Définit la structure d'un article
  - `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`: Champs GORM standard
  - `Title`, `Content`: Contenu de l'article
  - `Slug`: URL-friendly version du titre (ex: "mon-article")
  - `Tags`: Tags séparés par virgule
  - `Published`: Boolean pour publier/dépublier
  - `AuthorID`: Référence à l'utilisateur (User.ID)

- **BeforeCreate**: Hook GORM qui s'exécute avant l'insertion en DB
  - Génère automatiquement le slug depuis le titre

- **slugify**: Fonction helper pour créer un slug
  - "Mon Super Article!" devient "mon-super-article"

---

## Étape 6: Implémenter le service Post

### 6.1 Définir l'interface PostService

Créer `internal/interfaces/post_service.go`:

```go
package interfaces

import (
	"context"

	"blog-api/internal/models"
)

// PostService définit les opérations métier sur les articles
type PostService interface {
	Create(ctx context.Context, authorID uint, title, content, tags string) (*models.Post, error)
	GetByID(ctx context.Context, id uint) (*models.Post, error)
	GetBySlug(ctx context.Context, slug string) (*models.Post, error)
	List(ctx context.Context, limit, offset int) ([]*models.Post, int64, error)
	ListByAuthor(ctx context.Context, authorID uint, limit, offset int) ([]*models.Post, int64, error)
	Update(ctx context.Context, id uint, title, content, tags *string) (*models.Post, error)
	Publish(ctx context.Context, id uint) error
	Unpublish(ctx context.Context, id uint) error
	Delete(ctx context.Context, id uint) error
}
```

### 6.2 Définir l'interface PostRepository

Créer `internal/interfaces/post_repository.go`:

```go
package interfaces

import (
	"context"

	"blog-api/internal/models"
)

// PostRepository définit les opérations de persistance pour les articles
type PostRepository interface {
	Create(ctx context.Context, post *models.Post) error
	FindByID(ctx context.Context, id uint) (*models.Post, error)
	FindBySlug(ctx context.Context, slug string) (*models.Post, error)
	FindAll(ctx context.Context, limit, offset int) ([]*models.Post, int64, error)
	FindByAuthorID(ctx context.Context, authorID uint, limit, offset int) ([]*models.Post, int64, error)
	Update(ctx context.Context, post *models.Post) error
	Delete(ctx context.Context, id uint) error
}
```

### 6.3 Implémenter le service

Créer `internal/domain/post/service.go`:

```go
package post

import (
	"context"

	"blog-api/internal/domain"
	"blog-api/internal/interfaces"
	"github.com/rs/zerolog"
)

type service struct {
	repo   interfaces.PostRepository
	logger zerolog.Logger
}

// NewService crée une nouvelle instance du service Post
func NewService(repo interfaces.PostRepository, logger zerolog.Logger) interfaces.PostService {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

// Create crée un nouvel article
func (s *service) Create(ctx context.Context, authorID uint, title, content, tags string) (*Post, error) {
	post := &Post{
		Title:     title,
		Content:   content,
		Tags:      tags,
		AuthorID:  authorID,
		Published: false,
	}

	if err := s.repo.Create(ctx, post); err != nil {
		s.logger.Error().Err(err).Msg("Failed to create post")
		return nil, err
	}

	s.logger.Info().
		Uint("post_id", post.ID).
		Uint("author_id", authorID).
		Str("title", title).
		Msg("Post created successfully")

	return post, nil
}

// GetByID récupère un article par son ID
func (s *service) GetByID(ctx context.Context, id uint) (*Post, error) {
	post, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, domain.NewNotFoundError("Post not found", "POST_NOT_FOUND", err)
	}
	return post, nil
}

// GetBySlug récupère un article par son slug
func (s *service) GetBySlug(ctx context.Context, slug string) (*Post, error) {
	post, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, domain.NewNotFoundError("Post not found", "POST_NOT_FOUND", err)
	}
	return post, nil
}

// List récupère tous les articles avec pagination
func (s *service) List(ctx context.Context, limit, offset int) ([]*Post, int64, error) {
	return s.repo.FindAll(ctx, limit, offset)
}

// ListByAuthor récupère les articles d'un auteur avec pagination
func (s *service) ListByAuthor(ctx context.Context, authorID uint, limit, offset int) ([]*Post, int64, error) {
	return s.repo.FindByAuthorID(ctx, authorID, limit, offset)
}

// Update met à jour un article
func (s *service) Update(ctx context.Context, id uint, title, content, tags *string) (*Post, error) {
	post, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, domain.NewNotFoundError("Post not found", "POST_NOT_FOUND", err)
	}

	// Mettre à jour uniquement les champs fournis
	if title != nil {
		post.Title = *title
		post.Slug = slugify(*title) // Régénérer le slug
	}
	if content != nil {
		post.Content = *content
	}
	if tags != nil {
		post.Tags = *tags
	}

	if err := s.repo.Update(ctx, post); err != nil {
		s.logger.Error().Err(err).Uint("post_id", id).Msg("Failed to update post")
		return nil, err
	}

	s.logger.Info().Uint("post_id", id).Msg("Post updated successfully")
	return post, nil
}

// Publish publie un article
func (s *service) Publish(ctx context.Context, id uint) error {
	post, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return domain.NewNotFoundError("Post not found", "POST_NOT_FOUND", err)
	}

	post.Published = true
	if err := s.repo.Update(ctx, post); err != nil {
		s.logger.Error().Err(err).Uint("post_id", id).Msg("Failed to publish post")
		return err
	}

	s.logger.Info().Uint("post_id", id).Msg("Post published successfully")
	return nil
}

// Unpublish dépublie un article
func (s *service) Unpublish(ctx context.Context, id uint) error {
	post, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return domain.NewNotFoundError("Post not found", "POST_NOT_FOUND", err)
	}

	post.Published = false
	if err := s.repo.Update(ctx, post); err != nil {
		s.logger.Error().Err(err).Uint("post_id", id).Msg("Failed to unpublish post")
		return err
	}

	s.logger.Info().Uint("post_id", id).Msg("Post unpublished successfully")
	return nil
}

// Delete supprime un article (soft delete)
func (s *service) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error().Err(err).Uint("post_id", id).Msg("Failed to delete post")
		return domain.NewNotFoundError("Post not found", "POST_NOT_FOUND", err)
	}

	s.logger.Info().Uint("post_id", id).Msg("Post deleted successfully")
	return nil
}
```

**Points clés**:

- **Dependency Injection**: Le service reçoit le repository et le logger via le constructeur
- **Error handling**: Utilise les erreurs du domaine (`domain.NewNotFoundError`)
- **Logging structuré**: Log avec zerolog pour chaque opération
- **Business logic**: Gère la publication/dépublication, la génération de slug, etc.

---

## Étape 7: Créer le repository Post

Créer `internal/adapters/repository/post_repository.go`:

```go
package repository

import (
	"context"

	"blog-api/internal/models"
	"blog-api/internal/interfaces"
	"gorm.io/gorm"
)

type postRepository struct {
	db *gorm.DB
}

// NewPostRepository crée une nouvelle instance du repository Post
func NewPostRepository(db *gorm.DB) interfaces.PostRepository {
	return &postRepository{db: db}
}

// Create insère un nouvel article dans la base de données
func (r *postRepository) Create(ctx context.Context, post *models.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

// FindByID récupère un article par son ID
func (r *postRepository) FindByID(ctx context.Context, id uint) (*models.Post, error) {
	var p post.Post
	err := r.db.WithContext(ctx).First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// FindBySlug récupère un article par son slug
func (r *postRepository) FindBySlug(ctx context.Context, slug string) (*models.Post, error) {
	var p post.Post
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// FindAll récupère tous les articles avec pagination
// Retourne les posts + le total count
func (r *postRepository) FindAll(ctx context.Context, limit, offset int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).Model(&models.Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Récupérer les posts
	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&posts).Error

	return posts, total, err
}

// FindByAuthorID récupère les articles d'un auteur avec pagination
func (r *postRepository) FindByAuthorID(ctx context.Context, authorID uint, limit, offset int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var total int64

	query := r.db.WithContext(ctx).Where("author_id = ?", authorID)

	// Count total
	if err := query.Model(&models.Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Récupérer les posts
	err := query.
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&posts).Error

	return posts, total, err
}

// Update met à jour un article
func (r *postRepository) Update(ctx context.Context, post *models.Post) error {
	return r.db.WithContext(ctx).Save(post).Error
}

// Delete supprime un article (soft delete avec GORM)
func (r *postRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Post{}, id).Error
}
```

**Points clés**:

- **GORM**: Utilise GORM pour interagir avec PostgreSQL
- **Context**: Chaque méthode accepte un context pour les timeouts/annulations
- **Pagination**: FindAll et FindByAuthorID retournent total count + posts
- **Soft Delete**: GORM gère automatiquement le soft delete via DeletedAt

---

## Étape 8: Créer le handler HTTP

### 8.1 Créer le handler

Créer `internal/adapters/handlers/post_handler.go`:

```go
package handlers

import (
	"strconv"

	"blog-api/internal/interfaces"
	"blog-api/pkg/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type PostHandler struct {
	postService interfaces.PostService
	logger      zerolog.Logger
}

func NewPostHandler(postService interfaces.PostService, logger zerolog.Logger) *PostHandler {
	return &PostHandler{
		postService: postService,
		logger:      logger,
	}
}

// CreatePostRequest représente la requête de création d'article
type CreatePostRequest struct {
	Title   string `json:"title" validate:"required,max=255"`
	Content string `json:"content" validate:"required"`
	Tags    string `json:"tags"`
}

// UpdatePostRequest représente la requête de mise à jour d'article
type UpdatePostRequest struct {
	Title   *string `json:"title,omitempty" validate:"omitempty,max=255"`
	Content *string `json:"content,omitempty"`
	Tags    *string `json:"tags,omitempty"`
}

// Create crée un nouvel article
// POST /api/v1/posts
func (h *PostHandler) Create(c *fiber.Ctx) error {
	var req CreatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validation
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Récupérer l'utilisateur authentifié depuis le context
	userID := c.Locals("userID").(uint)

	// Créer le post
	post, err := h.postService.Create(c.Context(), userID, req.Title, req.Content, req.Tags)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to create post")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create post",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(post)
}

// Get récupère un article par ID ou slug
// GET /api/v1/posts/:idOrSlug
func (h *PostHandler) Get(c *fiber.Ctx) error {
	idOrSlug := c.Params("idOrSlug")

	// Essayer de parser comme ID
	if id, err := strconv.ParseUint(idOrSlug, 10, 32); err == nil {
		post, err := h.postService.GetByID(c.Context(), uint(id))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Post not found",
			})
		}
		return c.JSON(post)
	}

	// Sinon, chercher par slug
	post, err := h.postService.GetBySlug(c.Context(), idOrSlug)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Post not found",
		})
	}

	return c.JSON(post)
}

// List récupère tous les articles avec pagination
// GET /api/v1/posts?limit=10&offset=0
func (h *PostHandler) List(c *fiber.Ctx) error {
	// Paramètres de pagination
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	// Limiter le nombre max de résultats
	if limit > 100 {
		limit = 100
	}

	posts, total, err := h.postService.List(c.Context(), limit, offset)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list posts")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list posts",
		})
	}

	return c.JSON(fiber.Map{
		"data":   posts,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// ListByAuthor récupère les articles d'un auteur
// GET /api/v1/posts/author/:authorID?limit=10&offset=0
func (h *PostHandler) ListByAuthor(c *fiber.Ctx) error {
	authorID, err := strconv.ParseUint(c.Params("authorID"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid author ID",
		})
	}

	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	if limit > 100 {
		limit = 100
	}

	posts, total, err := h.postService.ListByAuthor(c.Context(), uint(authorID), limit, offset)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list posts by author")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list posts",
		})
	}

	return c.JSON(fiber.Map{
		"data":   posts,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// Update met à jour un article
// PUT /api/v1/posts/:id
func (h *PostHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid post ID",
		})
	}

	var req UpdatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validation
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Mettre à jour
	post, err := h.postService.Update(c.Context(), uint(id), req.Title, req.Content, req.Tags)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to update post")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update post",
		})
	}

	return c.JSON(post)
}

// Publish publie un article
// POST /api/v1/posts/:id/publish
func (h *PostHandler) Publish(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid post ID",
		})
	}

	if err := h.postService.Publish(c.Context(), uint(id)); err != nil {
		h.logger.Error().Err(err).Msg("Failed to publish post")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to publish post",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Post published successfully",
	})
}

// Unpublish dépublie un article
// POST /api/v1/posts/:id/unpublish
func (h *PostHandler) Unpublish(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid post ID",
		})
	}

	if err := h.postService.Unpublish(c.Context(), uint(id)); err != nil {
		h.logger.Error().Err(err).Msg("Failed to unpublish post")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to unpublish post",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Post unpublished successfully",
	})
}

// Delete supprime un article
// DELETE /api/v1/posts/:id
func (h *PostHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid post ID",
		})
	}

	if err := h.postService.Delete(c.Context(), uint(id)); err != nil {
		h.logger.Error().Err(err).Msg("Failed to delete post")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Post not found",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
```

**Points clés**:

- **Validation**: Utilise validator pour valider les requêtes
- **Authentication**: Récupère userID depuis le context (middleware auth)
- **Error handling**: Retourne des codes HTTP appropriés
- **Pagination**: Support limit/offset pour les listes

---

## Étape 9: Enregistrer les routes et le module

### 9.1 Créer le module fx

Créer `internal/domain/post/module.go`:

```go
package post

import (
	"blog-api/internal/adapters/handlers"
	"blog-api/internal/adapters/repository"
	"go.uber.org/fx"
)

// Module provides all Post domain dependencies
var Module = fx.Module("post",
	fx.Provide(
		repository.NewPostRepository,
		NewService,
		handlers.NewPostHandler,
	),
)
```

### 9.2 Enregistrer les routes

Modifier `internal/infrastructure/server/routes.go`:

Ajouter après les routes User existantes:

```go
// Post routes (protected)
postRoutes := v1.Group("/posts")
postRoutes.Get("/", postHandler.List)                    // Liste tous les posts
postRoutes.Get("/:idOrSlug", postHandler.Get)            // Récupérer par ID ou slug
postRoutes.Get("/author/:authorID", postHandler.ListByAuthor) // Posts par auteur

postRoutes.Use(authMiddleware.RequireAuth())             // Routes protégées ci-dessous
postRoutes.Post("/", postHandler.Create)                 // Créer un post
postRoutes.Put("/:id", postHandler.Update)               // Mettre à jour
postRoutes.Post("/:id/publish", postHandler.Publish)     // Publier
postRoutes.Post("/:id/unpublish", postHandler.Unpublish) // Dépublier
postRoutes.Delete("/:id", postHandler.Delete)            // Supprimer
```

Le fichier complet `routes.go` devient:

```go
package server

import (
	"blog-api/internal/adapters/handlers"
	"blog-api/internal/adapters/middleware"
	"github.com/gofiber/fiber/v2"
)

type RouteParams struct {
	App            *fiber.App
	AuthHandler    *handlers.AuthHandler
	UserHandler    *handlers.UserHandler
	PostHandler    *handlers.PostHandler  // Ajouté
	AuthMiddleware *middleware.AuthMiddleware
}

func RegisterRoutes(params RouteParams) {
	app := params.App
	authHandler := params.AuthHandler
	userHandler := params.UserHandler
	postHandler := params.PostHandler  // Ajouté
	authMiddleware := params.AuthMiddleware

	// Health check (public)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// API v1
	v1 := app.Group("/api/v1")

	// Auth routes (public)
	auth := v1.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)

	// User routes (protected)
	users := v1.Group("/users")
	users.Use(authMiddleware.RequireAuth())
	users.Get("/", userHandler.List)
	users.Get("/:id", userHandler.GetByID)
	users.Put("/:id", userHandler.Update)
	users.Delete("/:id", userHandler.Delete)

	// Post routes
	postRoutes := v1.Group("/posts")
	postRoutes.Get("/", postHandler.List)
	postRoutes.Get("/:idOrSlug", postHandler.Get)
	postRoutes.Get("/author/:authorID", postHandler.ListByAuthor)

	postRoutes.Use(authMiddleware.RequireAuth())
	postRoutes.Post("/", postHandler.Create)
	postRoutes.Put("/:id", postHandler.Update)
	postRoutes.Post("/:id/publish", postHandler.Publish)
	postRoutes.Post("/:id/unpublish", postHandler.Unpublish)
	postRoutes.Delete("/:id", postHandler.Delete)
}
```

### 9.3 Ajouter le module au main

Modifier `cmd/main.go`:

```go
package main

import (
	"context"

	"blog-api/internal/models"  // Ajouté
	"blog-api/internal/domain/user"
	"blog-api/internal/infrastructure/database"
	"blog-api/internal/infrastructure/server"
	"blog-api/pkg/config"
	"blog-api/pkg/logger"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		// Configuration
		fx.Provide(
			config.Load,
			logger.New,
		),

		// Infrastructure
		database.Module,
		server.Module,

		// Domains
		user.Module,
		post.Module,  // Ajouté

		fx.Invoke(func(lc fx.Lifecycle, srv *server.Server) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go srv.Start()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return srv.Shutdown()
				},
			})
		}),
	).Run()
}
```

### 9.4 Migration de la base de données

Modifier `internal/infrastructure/database/migrations.go`:

Ajouter l'entité Post aux migrations:

```go
package database

import (
	"blog-api/internal/models"  // Ajouté
	"blog-api/internal/domain/user"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// RunMigrations exécute les migrations automatiques pour toutes les entités
func RunMigrations(db *gorm.DB, logger zerolog.Logger) error {
	logger.Info().Msg("Running database migrations...")

	if err := db.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		&models.Post{},  // Ajouté
	); err != nil {
		logger.Error().Err(err).Msg("Failed to run migrations")
		return err
	}

	logger.Info().Msg("Database migrations completed successfully")
	return nil
}
```

---

## Étape 10: Tester l'API Posts

### 10.1 Relancer l'application

```bash
# Arrêter l'app (Ctrl+C)
# Relancer
make run
```

Les migrations vont créer la table `posts` automatiquement.

### 10.2 Créer un article

Récupérez d'abord un access token (voir Étape 4.3).

```bash
# Remplacez <ACCESS_TOKEN> par votre token
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Mon Premier Article",
    "content": "Ceci est le contenu de mon premier article de blog!",
    "tags": "golang,tutorial,blog"
  }'
```

**Réponse**:
```json
{
  "id": 1,
  "created_at": "2024-01-10T11:00:00Z",
  "updated_at": "2024-01-10T11:00:00Z",
  "title": "Mon Premier Article",
  "slug": "mon-premier-article",
  "content": "Ceci est le contenu de mon premier article de blog!",
  "tags": "golang,tutorial,blog",
  "published": false,
  "author_id": 1
}
```

Notez que le **slug** a été généré automatiquement!

### 10.3 Lister les articles

```bash
curl http://localhost:8080/api/v1/posts
```

**Réponse**:
```json
{
  "data": [
    {
      "id": 1,
      "created_at": "2024-01-10T11:00:00Z",
      "updated_at": "2024-01-10T11:00:00Z",
      "title": "Mon Premier Article",
      "slug": "mon-premier-article",
      "content": "Ceci est le contenu de mon premier article de blog!",
      "tags": "golang,tutorial,blog",
      "published": false,
      "author_id": 1
    }
  ],
  "total": 1,
  "limit": 10,
  "offset": 0
}
```

### 10.4 Récupérer un article par slug

```bash
curl http://localhost:8080/api/v1/posts/mon-premier-article
```

### 10.5 Publier l'article

```bash
curl -X POST http://localhost:8080/api/v1/posts/1/publish \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

**Réponse**:
```json
{
  "message": "Post published successfully"
}
```

### 10.6 Mettre à jour l'article

```bash
curl -X PUT http://localhost:8080/api/v1/posts/1 \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Mon Premier Article (Édité)",
    "content": "Contenu mis à jour avec plus d\'informations!"
  }'
```

### 10.7 Supprimer l'article

```bash
curl -X DELETE http://localhost:8080/api/v1/posts/1 \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

**Code retourné**: 204 No Content

:material-check-circle: **Checkpoint 3**: L'API Posts fonctionne complètement!

---

## Étape 11: Ajouter le domaine Comment

Maintenant, ajoutons les commentaires sur les articles.

### 11.1 Créer l'entité Comment

```bash
mkdir -p internal/domain/comment
```

Créer `internal/models/comment.go`:

```go
package models

import (
	"time"

	"gorm.io/gorm"
)

// Comment représente un commentaire sur un article
type Comment struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Contenu
	Content string `gorm:"type:text;not null" json:"content" validate:"required"`

	// Relations
	PostID   uint `gorm:"not null;index" json:"post_id"`
	AuthorID uint `gorm:"not null" json:"author_id"`
}
```

### 11.2 Créer le service Comment (simplifié)

Créer `internal/interfaces/comment_repository.go`:

```go
package interfaces

import (
	"context"

	"blog-api/internal/models"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *models.Comment) error
	FindByPost(ctx context.Context, postID uint) ([]*models.Comment, error)
	Delete(ctx context.Context, id uint) error
}
```

Créer `internal/domain/comment/service.go`:

```go
package comment

import (
	"context"

	"blog-api/internal/models"
	"blog-api/internal/interfaces"
	"github.com/rs/zerolog"
)

type Service struct {
	repo   interfaces.CommentRepository
	logger zerolog.Logger
}

func NewService(repo interfaces.CommentRepository, logger zerolog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

func (s *Service) Create(ctx context.Context, postID, authorID uint, content string) (*models.Comment, error) {
	comment := &models.Comment{
		PostID:   postID,
		AuthorID: authorID,
		Content:  content,
	}

	if err := s.repo.Create(ctx, comment); err != nil {
		s.logger.Error().Err(err).Msg("Failed to create comment")
		return nil, err
	}

	s.logger.Info().Uint("comment_id", comment.ID).Uint("post_id", postID).Msg("Comment created")
	return comment, nil
}

func (s *service) ListByPost(ctx context.Context, postID uint) ([]*Comment, error) {
	return s.repo.FindByPostID(ctx, postID)
}

func (s *service) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error().Err(err).Uint("comment_id", id).Msg("Failed to delete comment")
		return err
	}

	s.logger.Info().Uint("comment_id", id).Msg("Comment deleted")
	return nil
}
```

### 11.3 Créer le repository et handler

Je vais vous laisser créer ces fichiers en suivant le même pattern que Post:

- `internal/interfaces/comment_repository.go`
- `internal/adapters/repository/comment_repository.go`
- `internal/adapters/handlers/comment_handler.go`
- `internal/domain/comment/module.go`

### 11.4 Ajouter les routes

Dans `routes.go`:

```go
// Comment routes
commentRoutes := v1.Group("/comments")
commentRoutes.Get("/post/:postID", commentHandler.ListByPost)

commentRoutes.Use(authMiddleware.RequireAuth())
commentRoutes.Post("/", commentHandler.Create)
commentRoutes.Delete("/:id", commentHandler.Delete)
```

### 11.5 Mettre à jour les migrations

Dans `migrations.go`, ajouter `&models.Comment{}`.

:material-check-circle: **Checkpoint 4**: Les commentaires sont fonctionnels!

---

## Étape 12: Tests unitaires

### 12.1 Tester le service Post

Créer `internal/domain/post/service_test.go`:

```go
package post_test

import (
	"context"
	"testing"

	"blog-api/internal/models"
	"blog-api/internal/interfaces/mocks"
	"blog-api/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostService_Create(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.PostRepository)
	log := logger.New(&config.Config{AppEnv: "test"})
	service := post.NewService(mockRepo, log)

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Post")).
		Return(nil)

	// Act
	result, err := service.Create(context.Background(), 1, "Test Title", "Test Content", "tag1,tag2")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Title", result.Title)
	assert.Equal(t, "test-title", result.Slug)
	mockRepo.AssertExpectations(t)
}
```

### 12.2 Lancer les tests

```bash
make test
```

---

## Étape 13: Déploiement Docker

### 13.1 Build l'image Docker

```bash
make docker-build
```

### 13.2 Lancer avec docker-compose

Le fichier `docker-compose.yml` est déjà généré:

```bash
docker-compose up -d
```

Cela lance:
- L'application sur le port 8080
- PostgreSQL sur le port 5432

### 13.3 Vérifier le déploiement

```bash
curl http://localhost:8080/health
```

---

## Conclusion

Félicitations! :material-party-popper: Vous avez créé une API Blog complète avec:

:material-check-circle: **Authentification JWT** (User, Login, Register)
:material-check-circle: **Articles** (CRUD complet avec slug, tags, publish/unpublish)
:material-check-circle: **Commentaires** (Create, List, Delete)
:material-check-circle: **Relations** (Post → Author, Comment → Post + Author)
:material-check-circle: **Pagination** (Limit/Offset)
:material-check-circle: **Tests unitaires**
:material-check-circle: **Déploiement Docker**
:material-check-circle: **Architecture hexagonale**
:material-check-circle: **Logging structuré**
:material-check-circle: **Error handling centralisé**

### Résumé de ce que vous avez appris

1. **Installation** de create-go-starter
2. **Génération** d'un projet complet
3. **Configuration** (.env, PostgreSQL, JWT)
4. **Architecture hexagonale**:
   - Domain (entities, services)
   - Adapters (handlers, repositories)
   - Interfaces (ports)
5. **Dependency Injection** avec uber-go/fx
6. **GORM** (migrations, queries, relations)
7. **Fiber** (routes, middleware, handlers)
8. **Tests** avec testify et mocks
9. **Docker** et docker-compose

### Prochaines étapes

Pour aller plus loin:

- **Upload d'images** pour les articles
- **Recherche full-text** dans les posts
- **Likes/Votes** sur les articles
- **Catégories** pour organiser les posts
- **Swagger** pour documenter l'API
- **CI/CD** avec GitHub Actions
- **Kubernetes** pour déploiement en production

### Ressources

- [Guide des projets générés](./generated-project-guide.md) - Documentation complète
- [Repository exemple](https://github.com/tky0065/go-starter-kit/tree/main/examples/blog-api) - Code complet
- [Fiber documentation](https://docs.gofiber.io/)
- [GORM documentation](https://gorm.io/docs/)

**Bon coding!** :material-rocket-launch:
