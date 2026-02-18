# Tutorial: Build a Complete Blog API with create-go-starter

Step-by-step guide to building a Blog API with `create-go-starter`, from installation to deployment.

## Table of Contents

1. [Objective](#objective)
2. [Prerequisites](#prerequisites)
3. [Step 1: CLI Installation](#step-1-cli-installation)
4. [Step 2: Project Generation](#step-2-project-generation)
5. [Step 3: Initial Configuration](#step-3-initial-configuration)
6. [Step 4: Test the Base Project](#step-4-test-the-base-project)
7. [Step 5: Add the Post Domain](#step-5-add-the-post-domain)
8. [Step 6: Implement the Post Service](#step-6-implement-the-post-service)
9. [Step 7: Create the Post Repository](#step-7-create-the-post-repository)
10. [Step 8: Create the HTTP Handler](#step-8-create-the-http-handler)
11. [Step 9: Register Routes and the Module](#step-9-register-routes-and-the-module)
12. [Step 10: Test the Posts API](#step-10-test-the-posts-api)
13. [Step 11: Add the Comment Domain](#step-11-add-the-comment-domain)
14. [Step 12: Unit Tests](#step-12-unit-tests)
15. [Step 13: Docker Deployment](#step-13-docker-deployment)
16. [Conclusion](#conclusion)

---

## Objective

Build a complete REST API for a blog with:

- **Posts** with author, title, content, tags
- **Comments** on posts
- **JWT Authentication** (already included in create-go-starter)
- **Comprehensive tests**
- **Docker deployment**

By the end of this tutorial, you will have a production-ready Blog API with all best practices.

## Prerequisites

### Required Software

- **Go 1.25+** - [Download](https://golang.org/dl/)
- **PostgreSQL** or **Docker** - For the database
- **curl** or **Postman** - To test the API
- Code editor (VS Code, GoLand, etc.)

### Recommended Knowledge

- Go basics (structs, interfaces, error handling)
- REST API concepts
- Familiarity with SQL/PostgreSQL (basic)

No need to be an expert! This tutorial explains each step in detail.

---

## Step 1: CLI Installation

### Global Installation (recommended)

The simplest method to install `create-go-starter`:

```bash
go install github.com/tky0065/go-starter-kit/cmd/create-go-starter@latest
```

This command downloads, compiles, and installs the CLI globally.

### Verification

```bash
create-go-starter --help
```

You should see the help output displayed.

**Note**: If the command is not found, add `$GOPATH/bin` to your PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

---

## Step 2: Project Generation

### Create the Project

```bash
create-go-starter blog-api
```

This command generates **~45 files** with all the necessary architecture.

### Generated Structure

```bash
cd blog-api
tree -L 3
```

**Result**:
```
blog-api/
├── cmd/
│   └── main.go                       # Entry point with fx DI
├── internal/
│   ├── models/
│   │   └── user.go                   # Entities: User, RefreshToken, AuthResponse
│   ├── domain/
│   │   ├── user/                     # User domain (pre-generated)
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
├── .env                              # Configuration (auto-copied)
├── .env.example
├── Dockerfile
├── Makefile
├── go.mod
└── README.md
```

<i class="material-icons success">check_circle</i> **Checkpoint 1**: The project has been generated successfully.

---

## Step 3: Initial Configuration

### 3.1 Install Dependencies

```bash
cd blog-api
go mod tidy
```

This command downloads all dependencies (Fiber, GORM, fx, etc.).

### 3.2 Configure PostgreSQL

You have 2 options:

#### Option A: Docker (recommended)

```bash
docker run -d \
  --name blog-postgres \
  -e POSTGRES_DB=blog_api \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine
```

#### Option B: Local PostgreSQL

If PostgreSQL is installed locally:

```bash
createdb blog_api
```

### 3.3 Configure Environment Variables

Generate a secure JWT secret:

```bash
JWT_SECRET=$(openssl rand -base64 32)
echo "JWT_SECRET generated: $JWT_SECRET"
```

Edit the `.env` file:

```bash
nano .env
```

Contents of `.env`:

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
JWT_SECRET=<paste_the_generated_secret_here>
JWT_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=168h
```

**Important**: Replace `<paste_the_generated_secret_here>` with the generated JWT_SECRET.

---

## Step 4: Test the Base Project

### 4.1 Launch the Application

```bash
make run
```

You should see:

```
2024/01/10 10:00:00 INF Starting blog-api server on :8080
```

### 4.2 Test the Health Check

In another terminal:

```bash
curl http://localhost:8080/health
```

**Expected response**:
```json
{"status":"ok"}
```

### 4.3 Test the Default Authentication

#### Create a User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@blog.com",
    "password": "admin123"
  }'
```

**Response**:
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

#### Log In

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@blog.com",
    "password": "admin123"
  }'
```

**Same response** with access_token and refresh_token.

#### Test a Protected Route

```bash
# Replace <ACCESS_TOKEN> with the received token
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

**Response**:
```json
[
  {
    "id": 1,
    "email": "admin@blog.com",
    "created_at": "2024-01-10T10:05:00Z"
  }
]
```

<i class="material-icons success">check_circle</i> **Checkpoint 2**: The base project works perfectly with User and Auth.

---

## Step 5: Add the Post Domain

We will now add our first feature: blog posts.

### 5.1 Create the Post Entity

Create the file `internal/models/post.go`:

```go
package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

// Post represents a blog post
type Post struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Content
	Title   string `gorm:"not null;size:255" json:"title" validate:"required,max=255"`
	Slug    string `gorm:"uniqueIndex;not null;size:255" json:"slug"`
	Content string `gorm:"type:text;not null" json:"content" validate:"required"`

	// Metadata
	Tags      string `gorm:"size:500" json:"tags"`
	Published bool   `gorm:"default:false" json:"published"`

	// Relations
	AuthorID uint `gorm:"not null" json:"author_id"`
}

// BeforeCreate automatically generates a unique slug before insertion
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.Slug == "" {
		p.Slug = slugify(p.Title)
	}
	return nil
}

// slugify converts a title into a URL-friendly slug
// Example: "My Awesome Article!" -> "my-awesome-article"
func slugify(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters
	replacer := strings.NewReplacer(
		"!", "", "?", "", ".", "", ",", "",
		"'", "", "\"", "", ":", "", ";", "",
		"(", "", ")", "", "[", "", "]", "",
	)
	slug = replacer.Replace(slug)

	// Remove multiple dashes
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Remove leading/trailing dashes
	slug = strings.Trim(slug, "-")

	return slug
}
```

**Explanations**:

- **Post struct**: Defines the structure of a blog post
  - `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`: Standard GORM fields
  - `Title`, `Content`: Post content
  - `Slug`: URL-friendly version of the title (e.g., "my-article")
  - `Tags`: Comma-separated tags
  - `Published`: Boolean to publish/unpublish
  - `AuthorID`: Reference to the user (User.ID)

- **BeforeCreate**: GORM hook that executes before insertion into the DB
  - Automatically generates the slug from the title

- **slugify**: Helper function to create a slug
  - "Mon Super Article!" becomes "mon-super-article"

---

## Step 6: Implement the Post Service

### 6.1 Define the PostService Interface

Create `internal/interfaces/post_service.go`:

```go
package interfaces

import (
	"context"

	"blog-api/internal/models"
)

// PostService defines business operations on posts
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

### 6.2 Define the PostRepository Interface

Create `internal/interfaces/post_repository.go`:

```go
package interfaces

import (
	"context"

	"blog-api/internal/models"
)

// PostRepository defines persistence operations for posts
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

### 6.3 Implement the Service

Create `internal/domain/post/service.go`:

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

// NewService creates a new instance of the Post service
func NewService(repo interfaces.PostRepository, logger zerolog.Logger) interfaces.PostService {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

// Create creates a new post
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

// GetByID retrieves a post by its ID
func (s *service) GetByID(ctx context.Context, id uint) (*Post, error) {
	post, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, domain.NewNotFoundError("Post not found", "POST_NOT_FOUND", err)
	}
	return post, nil
}

// GetBySlug retrieves a post by its slug
func (s *service) GetBySlug(ctx context.Context, slug string) (*Post, error) {
	post, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, domain.NewNotFoundError("Post not found", "POST_NOT_FOUND", err)
	}
	return post, nil
}

// List retrieves all posts with pagination
func (s *service) List(ctx context.Context, limit, offset int) ([]*Post, int64, error) {
	return s.repo.FindAll(ctx, limit, offset)
}

// ListByAuthor retrieves posts by an author with pagination
func (s *service) ListByAuthor(ctx context.Context, authorID uint, limit, offset int) ([]*Post, int64, error) {
	return s.repo.FindByAuthorID(ctx, authorID, limit, offset)
}

// Update updates a post
func (s *service) Update(ctx context.Context, id uint, title, content, tags *string) (*Post, error) {
	post, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, domain.NewNotFoundError("Post not found", "POST_NOT_FOUND", err)
	}

	// Update only the provided fields
	if title != nil {
		post.Title = *title
		post.Slug = slugify(*title) // Regenerate the slug
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

// Publish publishes a post
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

// Unpublish unpublishes a post
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

// Delete deletes a post (soft delete)
func (s *service) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error().Err(err).Uint("post_id", id).Msg("Failed to delete post")
		return domain.NewNotFoundError("Post not found", "POST_NOT_FOUND", err)
	}

	s.logger.Info().Uint("post_id", id).Msg("Post deleted successfully")
	return nil
}
```

**Key points**:

- **Dependency Injection**: The service receives the repository and logger via the constructor
- **Error handling**: Uses domain errors (`domain.NewNotFoundError`)
- **Structured logging**: Logs with zerolog for each operation
- **Business logic**: Handles publishing/unpublishing, slug generation, etc.

---

## Step 7: Create the Post Repository

Create `internal/adapters/repository/post_repository.go`:

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

// NewPostRepository creates a new instance of the Post repository
func NewPostRepository(db *gorm.DB) interfaces.PostRepository {
	return &postRepository{db: db}
}

// Create inserts a new post into the database
func (r *postRepository) Create(ctx context.Context, post *models.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

// FindByID retrieves a post by its ID
func (r *postRepository) FindByID(ctx context.Context, id uint) (*models.Post, error) {
	var p post.Post
	err := r.db.WithContext(ctx).First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// FindBySlug retrieves a post by its slug
func (r *postRepository) FindBySlug(ctx context.Context, slug string) (*models.Post, error) {
	var p post.Post
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// FindAll retrieves all posts with pagination
// Returns the posts + total count
func (r *postRepository) FindAll(ctx context.Context, limit, offset int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).Model(&models.Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Retrieve posts
	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&posts).Error

	return posts, total, err
}

// FindByAuthorID retrieves posts by an author with pagination
func (r *postRepository) FindByAuthorID(ctx context.Context, authorID uint, limit, offset int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var total int64

	query := r.db.WithContext(ctx).Where("author_id = ?", authorID)

	// Count total
	if err := query.Model(&models.Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Retrieve posts
	err := query.
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&posts).Error

	return posts, total, err
}

// Update updates a post
func (r *postRepository) Update(ctx context.Context, post *models.Post) error {
	return r.db.WithContext(ctx).Save(post).Error
}

// Delete deletes a post (soft delete with GORM)
func (r *postRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Post{}, id).Error
}
```

**Key points**:

- **GORM**: Uses GORM to interact with PostgreSQL
- **Context**: Each method accepts a context for timeouts/cancellations
- **Pagination**: FindAll and FindByAuthorID return total count + posts
- **Soft Delete**: GORM automatically handles soft delete via DeletedAt

---

## Step 8: Create the HTTP Handler

### 8.1 Create the Handler

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

// CreatePostRequest represents the post creation request
type CreatePostRequest struct {
	Title   string `json:"title" validate:"required,max=255"`
	Content string `json:"content" validate:"required"`
	Tags    string `json:"tags"`
}

// UpdatePostRequest represents the post update request
type UpdatePostRequest struct {
	Title   *string `json:"title,omitempty" validate:"omitempty,max=255"`
	Content *string `json:"content,omitempty"`
	Tags    *string `json:"tags,omitempty"`
}

// Create creates a new post
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

// Get retrieves a post by ID or slug
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

// List retrieves all posts with pagination
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

// ListByAuthor retrieves posts by an author
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

// Update updates a post
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

// Publish publishes a post
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

// Unpublish unpublishes a post
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

// Delete deletes a post
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

## Step 9: Register Routes and the Module

### 9.1 Create the fx Module

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

### 9.2 Register the Routes

Modify `internal/infrastructure/server/routes.go`:

Add after the existing User routes:

```go
// Post routes (protected)
postRoutes := v1.Group("/posts")
postRoutes.Get("/", postHandler.List)                    // List all posts
postRoutes.Get("/:idOrSlug", postHandler.Get)            // Retrieve by ID or slug
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

### 9.3 Add the Module to main

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

### 9.4 Database Migration

Modify `internal/infrastructure/database/migrations.go`:

Add the Post entity to the migrations:

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

### 10.1 Restart the Application

```bash
# Stop the app (Ctrl+C)
# Restart
make run
```

The migrations will create the `posts` table automatically.

### 10.2 Create a Post

First retrieve an access token (see Step 4.3).

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

### 10.3 List Posts

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

### 10.4 Retrieve a Post by Slug

```bash
curl http://localhost:8080/api/v1/posts/mon-premier-article
```

### 10.5 Publish the Post

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

### 10.6 Update the Post

```bash
curl -X PUT http://localhost:8080/api/v1/posts/1 \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Mon Premier Article (Édité)",
    "content": "Contenu mis à jour avec plus d'\''informations!"
  }'
```

### 10.7 Delete the Post

```bash
curl -X DELETE http://localhost:8080/api/v1/posts/1 \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

**Returned code**: 204 No Content

<i class="material-icons success">check_circle</i> **Checkpoint 3**: The Posts API is fully functional!

---

## Step 11: Add the Comment Domain

Now, let's add comments on posts.

### 11.1 Create the Comment Entity

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

// Comment represents a comment on a post
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

### 11.2 Create the Comment Service (simplified)

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

### 11.3 Create the Repository and Handler

I'll let you create these files following the same pattern as Post:

- `internal/interfaces/comment_repository.go`
- `internal/adapters/repository/comment_repository.go`
- `internal/adapters/handlers/comment_handler.go`
- `internal/domain/comment/module.go`

### 11.4 Add the Routes

In `routes.go`:

```go
// Comment routes
commentRoutes := v1.Group("/comments")
commentRoutes.Get("/post/:postID", commentHandler.ListByPost)

commentRoutes.Use(authMiddleware.RequireAuth())
commentRoutes.Post("/", commentHandler.Create)
commentRoutes.Delete("/:id", commentHandler.Delete)
```

### 11.5 Update the Migrations

In `migrations.go`, add `&models.Comment{}`.

<i class="material-icons success">check_circle</i> **Checkpoint 4**: Comments are functional!

---

## Step 12: Unit Tests

### 12.1 Test the Post Service

Create `internal/domain/post/service_test.go`:

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

### 12.2 Run the Tests

```bash
make test
```

---

## Step 13: Docker Deployment

### 13.1 Build the Docker Image

```bash
make docker-build
```

### 13.2 Launch with docker-compose

The `docker-compose.yml` file is already generated:

```bash
docker-compose up -d
```

This launches:
- The application on port 8080
- PostgreSQL on port 5432

### 13.3 Verify the Deployment

```bash
curl http://localhost:8080/health
```

---

## Conclusion

Congratulations! <i class="material-icons success">celebration</i> You have built a complete Blog API with:

<i class="material-icons success">check_circle</i> **JWT Authentication** (User, Login, Register)
<i class="material-icons success">check_circle</i> **Posts** (Full CRUD with slug, tags, publish/unpublish)
<i class="material-icons success">check_circle</i> **Comments** (Create, List, Delete)
<i class="material-icons success">check_circle</i> **Relations** (Post -> Author, Comment -> Post + Author)
<i class="material-icons success">check_circle</i> **Pagination** (Limit/Offset)
<i class="material-icons success">check_circle</i> **Unit tests**
<i class="material-icons success">check_circle</i> **Docker deployment**
<i class="material-icons success">check_circle</i> **Hexagonal architecture**
<i class="material-icons success">check_circle</i> **Structured logging**
<i class="material-icons success">check_circle</i> **Centralized error handling**

### Summary of What You Learned

1. **Installation** of create-go-starter
2. **Generation** of a complete project
3. **Configuration** (.env, PostgreSQL, JWT)
4. **Hexagonal architecture**:
   - Domain (entities, services)
   - Adapters (handlers, repositories)
   - Interfaces (ports)
5. **Dependency Injection** with uber-go/fx
6. **GORM** (migrations, queries, relations)
7. **Fiber** (routes, middleware, handlers)
8. **Tests** with testify and mocks
9. **Docker** and docker-compose

### Next Steps

To go further:

- **Image uploads** for posts
- **Full-text search** in posts
- **Likes/Votes** on posts
- **Categories** to organize posts
- **Swagger** to document the API
- **CI/CD** with GitHub Actions
- **Kubernetes** for production deployment

### Resources

- [Generated project guide](./generated-project-guide.md) - Complete documentation
- [Example repository](https://github.com/tky0065/go-starter-kit/tree/main/examples/blog-api) - Full code
- [Fiber documentation](https://docs.gofiber.io/)
- [GORM documentation](https://gorm.io/docs/)

**Happy coding!** <i class="material-icons info">rocket_launch</i>
