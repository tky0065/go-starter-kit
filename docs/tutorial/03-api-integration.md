# Partie 3: Exposer l'API HTTP

<i class="material-icons warning small">circle</i> **Partie 3/4** - Temps estimé: 30 minutes

[<i class="material-icons">arrow_back</i> Retour à l'index](index.md)

---

## Objectif

Dans cette partie, vous allez exposer le domaine Posts via une API HTTP REST avec Fiber, enregistrer les routes et tester l'API complète.

**Ce que vous allez créer**:
- Handler HTTP pour Posts (CRUD complet)
- Module fx pour l'injection de dépendances
- Enregistrement des routes
- Migration de base de données
- Tests de l'API avec curl

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

Récupérez d'abord un access token (voir Partie 1, Étape 4.3).

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
    "content": "Contenu mis à jour avec plus d'\''informations!"
  }'
```

### 10.7 Supprimer l'article

```bash
curl -X DELETE http://localhost:8080/api/v1/posts/1 \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

**Code retourné**: 204 No Content

<i class="material-icons success">check_circle</i> **Checkpoint 3**: L'API Posts fonctionne complètement!

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

<i class="material-icons success">check_circle</i> **Checkpoint 4**: Les commentaires sont fonctionnels!

---

## Résumé de la Partie 3

<i class="material-icons success">check</i> Handler HTTP créé pour Posts (CRUD complet)
<i class="material-icons success">check</i> Module fx configuré avec injection de dépendances
<i class="material-icons success">check</i> Routes enregistrées (publiques et protégées)
<i class="material-icons success">check</i> Migration de base de données
<i class="material-icons success">check</i> API testée avec curl
<i class="material-icons success">check</i> Domaine Comment ajouté (exercice)

L'API Blog est maintenant complètement fonctionnelle avec Posts et Commentaires!

---

## Navigation

[<i class="material-icons">arrow_back</i> Partie 2: Créer votre premier domaine](02-first-domain.md) | [Partie 4: Tests et Déploiement <i class="material-icons">arrow_forward</i>](04-testing-deployment.md)
