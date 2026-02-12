# Partie 2: Créer votre premier domaine (Posts)

<i class="material-icons success small">circle</i> **Partie 2/4** - Temps estimé: 30 minutes

[<i class="material-icons">arrow_back</i> Retour à l'index](index.md)

---

## Objectif

Dans cette partie, vous allez implémenter le domaine Posts (articles de blog) en suivant l'architecture hexagonale.

**Ce que vous allez créer**:
- Entité Post avec GORM
- Interface PostService (port)
- Implémentation du service Post
- Interface PostRepository (port)
- Implémentation du repository Post

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

## Résumé de la Partie 2

<i class="material-icons success">check</i> Création de l'entité Post avec GORM
<i class="material-icons success">check</i> Définition des interfaces (PostService, PostRepository)
<i class="material-icons success">check</i> Implémentation du service Post avec business logic
<i class="material-icons success">check</i> Implémentation du repository Post avec GORM

<i class="material-icons success">check_circle</i> **Checkpoint**: Le domaine Post est implémenté (service + repository)

Dans la prochaine partie, nous allons exposer ce domaine via l'API HTTP.

---

## Navigation

[<i class="material-icons">arrow_back</i> Partie 1: Installation et Configuration](01-setup.md) | [Partie 3: Exposer l'API HTTP <i class="material-icons">arrow_forward</i>](03-api-integration.md)
