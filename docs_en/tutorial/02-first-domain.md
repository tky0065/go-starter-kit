# Part 2: Create Your First Domain (Posts)

<i class="material-icons success small">circle</i> **Part 2/4** - Estimated time: 30 minutes

[<i class="material-icons">arrow_back</i> Back to index](index.md)

---

## Goal

In this part, you will implement the Posts domain (blog articles) following hexagonal architecture.

**What you will create**:
- Post entity with GORM
- PostService interface (port)
- Post service implementation
- PostRepository interface (port)
- Post repository implementation

---

## Step 5: Add the Post (Article) Domain

We will now add our first feature: blog articles.

### 5.1 Create the Post entity

Create the file `internal/models/post.go`:

```go
package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

// Post represents a blog article
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

	// Remove multiple hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Remove leading/trailing hyphens
	slug = strings.Trim(slug, "-")

	return slug
}
```

**Explanations**:

- **Post struct**: Defines the structure of an article
  - `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`: Standard GORM fields
  - `Title`, `Content`: Article content
  - `Slug`: URL-friendly version of the title (e.g. "my-article")
  - `Tags`: Comma-separated tags
  - `Published`: Boolean to publish/unpublish
  - `AuthorID`: Reference to the user (User.ID)

- **BeforeCreate**: GORM hook that runs before insertion in the DB
  - Automatically generates the slug from the title

- **slugify**: Helper function to create a slug
  - "My Awesome Article!" becomes "my-awesome-article"

---

## Step 6: Implement the Post Service

### 6.1 Define the PostService interface

Create `internal/interfaces/post_service.go`:

```go
package interfaces

import (
	"context"

	"blog-api/internal/models"
)

// PostService defines the business operations on articles
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

### 6.2 Define the PostRepository interface

Create `internal/interfaces/post_repository.go`:

```go
package interfaces

import (
	"context"

	"blog-api/internal/models"
)

// PostRepository defines the persistence operations for articles
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

### 6.3 Implement the service

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

// Create creates a new article
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

// GetByID retrieves an article by its ID
func (s *service) GetByID(ctx context.Context, id uint) (*Post, error) {
	post, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, domain.NewNotFoundError("Post not found", "POST_NOT_FOUND", err)
	}
	return post, nil
}

// GetBySlug retrieves an article by its slug
func (s *service) GetBySlug(ctx context.Context, slug string) (*Post, error) {
	post, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, domain.NewNotFoundError("Post not found", "POST_NOT_FOUND", err)
	}
	return post, nil
}

// List retrieves all articles with pagination
func (s *service) List(ctx context.Context, limit, offset int) ([]*Post, int64, error) {
	return s.repo.FindAll(ctx, limit, offset)
}

// ListByAuthor retrieves an author's articles with pagination
func (s *service) ListByAuthor(ctx context.Context, authorID uint, limit, offset int) ([]*Post, int64, error) {
	return s.repo.FindByAuthorID(ctx, authorID, limit, offset)
}

// Update updates an article
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

// Publish publishes an article
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

// Unpublish unpublishes an article
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

// Delete deletes an article (soft delete)
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

// Create inserts a new article into the database
func (r *postRepository) Create(ctx context.Context, post *models.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

// FindByID retrieves an article by its ID
func (r *postRepository) FindByID(ctx context.Context, id uint) (*models.Post, error) {
	var p post.Post
	err := r.db.WithContext(ctx).First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// FindBySlug retrieves an article by its slug
func (r *postRepository) FindBySlug(ctx context.Context, slug string) (*models.Post, error) {
	var p post.Post
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// FindAll retrieves all articles with pagination
// Returns the posts + the total count
func (r *postRepository) FindAll(ctx context.Context, limit, offset int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).Model(&models.Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Retrieve the posts
	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&posts).Error

	return posts, total, err
}

// FindByAuthorID retrieves an author's articles with pagination
func (r *postRepository) FindByAuthorID(ctx context.Context, authorID uint, limit, offset int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var total int64

	query := r.db.WithContext(ctx).Where("author_id = ?", authorID)

	// Count total
	if err := query.Model(&models.Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Retrieve the posts
	err := query.
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&posts).Error

	return posts, total, err
}

// Update updates an article
func (r *postRepository) Update(ctx context.Context, post *models.Post) error {
	return r.db.WithContext(ctx).Save(post).Error
}

// Delete deletes an article (soft delete with GORM)
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

## Part 2 Summary

<i class="material-icons success">check</i> Creation of the Post entity with GORM
<i class="material-icons success">check</i> Definition of interfaces (PostService, PostRepository)
<i class="material-icons success">check</i> Implementation of the Post service with business logic
<i class="material-icons success">check</i> Implementation of the Post repository with GORM

<i class="material-icons success">check_circle</i> **Checkpoint**: The Post domain is implemented (service + repository)

In the next part, we will expose this domain via the HTTP API.

---

## Navigation

[<i class="material-icons">arrow_back</i> Part 1: Installation and Configuration](01-setup.md) | [Part 3: Expose the HTTP API <i class="material-icons">arrow_forward</i>](03-api-integration.md)
