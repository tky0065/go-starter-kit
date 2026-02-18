# Part 3: Expose the HTTP API

<i class="material-icons warning small">circle</i> **Part 3/4** - Estimated time: 30 minutes

[<i class="material-icons">arrow_back</i> Back to index](index.md)

---

## Goal

In this part, you will expose the Posts domain via an HTTP REST API with Fiber, register the routes and test the complete API.

**What you will create**:
- HTTP handler for Posts (full CRUD)
- fx module for dependency injection
- Route registration
- Database migration
- API tests with curl

---

## Step 8: Create the HTTP Handler

### 8.1 Create the handler

Create `internal/adapters/handlers/post_handler.go`:

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

// CreatePostRequest represents the article creation request
type CreatePostRequest struct {
	Title   string `json:"title" validate:"required,max=255"`
	Content string `json:"content" validate:"required"`
	Tags    string `json:"tags"`
}

// UpdatePostRequest represents the article update request
type UpdatePostRequest struct {
	Title   *string `json:"title,omitempty" validate:"omitempty,max=255"`
	Content *string `json:"content,omitempty"`
	Tags    *string `json:"tags,omitempty"`
}

// Create creates a new article
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

	// Retrieve the authenticated user from the context
	userID := c.Locals("userID").(uint)

	// Create the post
	post, err := h.postService.Create(c.Context(), userID, req.Title, req.Content, req.Tags)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to create post")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create post",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(post)
}

// Get retrieves an article by ID or slug
// GET /api/v1/posts/:idOrSlug
func (h *PostHandler) Get(c *fiber.Ctx) error {
	idOrSlug := c.Params("idOrSlug")

	// Try to parse as ID
	if id, err := strconv.ParseUint(idOrSlug, 10, 32); err == nil {
		post, err := h.postService.GetByID(c.Context(), uint(id))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Post not found",
			})
		}
		return c.JSON(post)
	}

	// Otherwise, search by slug
	post, err := h.postService.GetBySlug(c.Context(), idOrSlug)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Post not found",
		})
	}

	return c.JSON(post)
}

// List retrieves all articles with pagination
// GET /api/v1/posts?limit=10&offset=0
func (h *PostHandler) List(c *fiber.Ctx) error {
	// Pagination parameters
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	// Limit the maximum number of results
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

// ListByAuthor retrieves an author's articles
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

// Update updates an article
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

	// Update
	post, err := h.postService.Update(c.Context(), uint(id), req.Title, req.Content, req.Tags)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to update post")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update post",
		})
	}

	return c.JSON(post)
}

// Publish publishes an article
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

// Unpublish unpublishes an article
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

// Delete deletes an article
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

**Key points**:

- **Validation**: Uses validator to validate requests
- **Authentication**: Retrieves userID from the context (auth middleware)
- **Error handling**: Returns appropriate HTTP status codes
- **Pagination**: Supports limit/offset for lists

---

## Step 9: Register the Routes and Module

### 9.1 Create the fx module

Create `internal/domain/post/module.go`:

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

### 9.2 Register the routes

Modify `internal/infrastructure/server/routes.go`:

Add after the existing User routes:

```go
// Post routes (protected)
postRoutes := v1.Group("/posts")
postRoutes.Get("/", postHandler.List)                    // List all posts
postRoutes.Get("/:idOrSlug", postHandler.Get)            // Get by ID or slug
postRoutes.Get("/author/:authorID", postHandler.ListByAuthor) // Posts by author

postRoutes.Use(authMiddleware.RequireAuth())             // Protected routes below
postRoutes.Post("/", postHandler.Create)                 // Create a post
postRoutes.Put("/:id", postHandler.Update)               // Update
postRoutes.Post("/:id/publish", postHandler.Publish)     // Publish
postRoutes.Post("/:id/unpublish", postHandler.Unpublish) // Unpublish
postRoutes.Delete("/:id", postHandler.Delete)            // Delete
```

The complete `routes.go` file becomes:

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
	PostHandler    *handlers.PostHandler  // Added
	AuthMiddleware *middleware.AuthMiddleware
}

func RegisterRoutes(params RouteParams) {
	app := params.App
	authHandler := params.AuthHandler
	userHandler := params.UserHandler
	postHandler := params.PostHandler  // Added
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

### 9.3 Add the module to main

Modify `cmd/main.go`:

```go
package main

import (
	"context"

	"blog-api/internal/models"  // Added
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
		post.Module,  // Added

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

### 9.4 Database migration

Modify `internal/infrastructure/database/migrations.go`:

Add the Post entity to migrations:

```go
package database

import (
	"blog-api/internal/models"  // Added
	"blog-api/internal/domain/user"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// RunMigrations runs automatic migrations for all entities
func RunMigrations(db *gorm.DB, logger zerolog.Logger) error {
	logger.Info().Msg("Running database migrations...")

	if err := db.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		&models.Post{},  // Added
	); err != nil {
		logger.Error().Err(err).Msg("Failed to run migrations")
		return err
	}

	logger.Info().Msg("Database migrations completed successfully")
	return nil
}
```

---

## Step 10: Test the Posts API

### 10.1 Restart the application

```bash
# Stop the app (Ctrl+C)
# Restart
make run
```

The migrations will automatically create the `posts` table.

### 10.2 Create an article

First retrieve an access token (see Part 1, Step 4.3).

```bash
# Replace <ACCESS_TOKEN> with your token
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Mon Premier Article",
    "content": "Ceci est le contenu de mon premier article de blog!",
    "tags": "golang,tutorial,blog"
  }'
```

**Response**:
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

Note that the **slug** was generated automatically!

### 10.3 List the articles

```bash
curl http://localhost:8080/api/v1/posts
```

**Response**:
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

### 10.4 Retrieve an article by slug

```bash
curl http://localhost:8080/api/v1/posts/mon-premier-article
```

### 10.5 Publish the article

```bash
curl -X POST http://localhost:8080/api/v1/posts/1/publish \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

**Response**:
```json
{
  "message": "Post published successfully"
}
```

### 10.6 Update the article

```bash
curl -X PUT http://localhost:8080/api/v1/posts/1 \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Mon Premier Article (Édité)",
    "content": "Contenu mis à jour avec plus d'\''informations!"
  }'
```

### 10.7 Delete the article

```bash
curl -X DELETE http://localhost:8080/api/v1/posts/1 \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

**Returned code**: 204 No Content

<i class="material-icons success">check_circle</i> **Checkpoint 3**: The Posts API is fully working!

---

## Step 11: Add the Comment Domain

Now, let's add comments on articles.

### 11.1 Create the Comment entity

```bash
mkdir -p internal/domain/comment
```

Create `internal/models/comment.go`:

```go
package models

import (
	"time"

	"gorm.io/gorm"
)

// Comment represents a comment on an article
type Comment struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Content
	Content string `gorm:"type:text;not null" json:"content" validate:"required"`

	// Relations
	PostID   uint `gorm:"not null;index" json:"post_id"`
	AuthorID uint `gorm:"not null" json:"author_id"`
}
```

### 11.2 Create the Comment service (simplified)

Create `internal/interfaces/comment_repository.go`:

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

Create `internal/domain/comment/service.go`:

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

### 11.3 Create the repository and handler

We'll leave you to create these files following the same pattern as Post:

- `internal/interfaces/comment_repository.go`
- `internal/adapters/repository/comment_repository.go`
- `internal/adapters/handlers/comment_handler.go`
- `internal/domain/comment/module.go`

### 11.4 Add the routes

In `routes.go`:

```go
// Comment routes
commentRoutes := v1.Group("/comments")
commentRoutes.Get("/post/:postID", commentHandler.ListByPost)

commentRoutes.Use(authMiddleware.RequireAuth())
commentRoutes.Post("/", commentHandler.Create)
commentRoutes.Delete("/:id", commentHandler.Delete)
```

### 11.5 Update the migrations

In `migrations.go`, add `&models.Comment{}`.

<i class="material-icons success">check_circle</i> **Checkpoint 4**: Comments are functional!

---

## Part 3 Summary

<i class="material-icons success">check</i> HTTP handler created for Posts (full CRUD)
<i class="material-icons success">check</i> fx module configured with dependency injection
<i class="material-icons success">check</i> Routes registered (public and protected)
<i class="material-icons success">check</i> Database migration
<i class="material-icons success">check</i> API tested with curl
<i class="material-icons success">check</i> Comment domain added (exercise)

The Blog API is now fully functional with Posts and Comments!

---

## Navigation

[<i class="material-icons">arrow_back</i> Part 2: Create Your First Domain](02-first-domain.md) | [Part 4: Testing and Deployment <i class="material-icons">arrow_forward</i>](04-testing-deployment.md)
