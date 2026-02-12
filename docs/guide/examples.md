# Exemples pratiques

<div class="navigation">
  <a href="index.md"><i class="material-icons">arrow_back</i> Guide Index</a>
</div>

---

## Exemple complet : Entite Product

Nous allons creer une entite `Product` complete avec CRUD. Suivez chaque etape dans l'ordre.

> **Conseil** : Remplacez `mon-projet` par le nom de votre projet dans tous les imports.

---

#### Etape 1 : Creer le Model (Entite)

**Fichier a creer : `internal/models/product.go`**

```go
package models

import (
    "time"

    "gorm.io/gorm"
)

// Product represents a product in the catalog
type Product struct {
    ID          uint           `gorm:"primaryKey" json:"id"`
    Name        string         `gorm:"not null;size:255" json:"name"`
    Description string         `gorm:"type:text" json:"description"`
    Price       float64        `gorm:"not null" json:"price"`
    Stock       int            `gorm:"default:0" json:"stock"`
    SKU         string         `gorm:"uniqueIndex;size:100" json:"sku"`
    Active      bool           `gorm:"default:true" json:"active"`
    CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// ProductResponse is the DTO for API responses (controls what is exposed)
type ProductResponse struct {
    ID          uint    `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    Stock       int     `json:"stock"`
    SKU         string  `json:"sku"`
    Active      bool    `json:"active"`
}

// ToResponse converts Product entity to ProductResponse DTO
func (p *Product) ToResponse() ProductResponse {
    return ProductResponse{
        ID:          p.ID,
        Name:        p.Name,
        Description: p.Description,
        Price:       p.Price,
        Stock:       p.Stock,
        SKU:         p.SKU,
        Active:      p.Active,
    }
}
```

**Pourquoi ?**
- Les entites sont centralisees dans `models/` pour eviter les dependances circulaires
- Tags GORM pour la configuration de la base de donnees
- Tags JSON pour controler la serialisation API
- DTO separe (`ProductResponse`) pour controler ce qui est expose a l'API

---

#### Etape 2 : Definir l'Interface (Port)

**Fichier a creer : `internal/interfaces/product_repository.go`**

```go
package interfaces

import (
    "context"

    "mon-projet/internal/models"
)

// ProductRepository defines the contract for product data access
// This is the "Port" in hexagonal architecture
type ProductRepository interface {
    Create(ctx context.Context, product *models.Product) error
    FindByID(ctx context.Context, id uint) (*models.Product, error)
    FindBySKU(ctx context.Context, sku string) (*models.Product, error)
    FindAll(ctx context.Context, limit, offset int) ([]*models.Product, error)
    FindActive(ctx context.Context) ([]*models.Product, error)
    Update(ctx context.Context, product *models.Product) error
    Delete(ctx context.Context, id uint) error
    Count(ctx context.Context) (int64, error)
}

// ProductService defines the contract for product business logic
type ProductService interface {
    Create(ctx context.Context, name, description, sku string, price float64, stock int) (*models.Product, error)
    GetByID(ctx context.Context, id uint) (*models.Product, error)
    GetAll(ctx context.Context, page, pageSize int) ([]*models.Product, int64, error)
    Update(ctx context.Context, id uint, name, description string, price float64, stock int, active bool) (*models.Product, error)
    Delete(ctx context.Context, id uint) error
    UpdateStock(ctx context.Context, id uint, quantity int) error
}
```

**Pourquoi ?**
- **Abstraction complete** : le domain ne connait pas GORM
- **Contrat clair** : toutes les operations disponibles sont definies
- **Testable** : facile a mocker pour les tests unitaires

---

#### Etape 3 : Implementer le Repository (Adapter)

**Fichier a creer : `internal/adapters/repository/product_repository.go`**

```go
package repository

import (
    "context"

    "gorm.io/gorm"

    "mon-projet/internal/domain"
    "mon-projet/internal/interfaces"
    "mon-projet/internal/models"
)

type productRepositoryGORM struct {
    db *gorm.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *gorm.DB) interfaces.ProductRepository {
    return &productRepositoryGORM{db: db}
}

func (r *productRepositoryGORM) Create(ctx context.Context, product *models.Product) error {
    return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepositoryGORM) FindByID(ctx context.Context, id uint) (*models.Product, error) {
    var product models.Product
    err := r.db.WithContext(ctx).First(&product, id).Error
    if err == gorm.ErrRecordNotFound {
        return nil, domain.NewNotFoundError("Product not found", "PRODUCT_NOT_FOUND", err)
    }
    return &product, err
}

func (r *productRepositoryGORM) FindBySKU(ctx context.Context, sku string) (*models.Product, error) {
    var product models.Product
    err := r.db.WithContext(ctx).Where("sku = ?", sku).First(&product).Error
    if err == gorm.ErrRecordNotFound {
        return nil, domain.NewNotFoundError("Product not found", "PRODUCT_NOT_FOUND", err)
    }
    return &product, err
}

func (r *productRepositoryGORM) FindAll(ctx context.Context, limit, offset int) ([]*models.Product, error) {
    var products []*models.Product
    err := r.db.WithContext(ctx).
        Limit(limit).
        Offset(offset).
        Order("created_at DESC").
        Find(&products).Error
    return products, err
}

func (r *productRepositoryGORM) FindActive(ctx context.Context) ([]*models.Product, error) {
    var products []*models.Product
    err := r.db.WithContext(ctx).Where("active = ?", true).Find(&products).Error
    return products, err
}

func (r *productRepositoryGORM) Update(ctx context.Context, product *models.Product) error {
    return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepositoryGORM) Delete(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}

func (r *productRepositoryGORM) Count(ctx context.Context) (int64, error) {
    var count int64
    err := r.db.WithContext(ctx).Model(&models.Product{}).Count(&count).Error
    return count, err
}
```

**Points cles :**
- Implemente l'interface `ProductRepository`
- Utilise `WithContext(ctx)` pour la propagation du contexte
- Convertit `gorm.ErrRecordNotFound` en `DomainError`

---

#### Etape 4 : Creer le Service (Domain/Business Logic)

**Creer le dossier :**
```bash
mkdir -p internal/domain/product
```

**Fichier a creer : `internal/domain/product/service.go`**

```go
package product

import (
    "context"
    "fmt"

    "github.com/rs/zerolog"

    "mon-projet/internal/domain"
    "mon-projet/internal/interfaces"
    "mon-projet/internal/models"
)

// Service handles product business logic
type Service struct {
    repo   interfaces.ProductRepository
    logger zerolog.Logger
}

// NewService creates a new product service
func NewService(repo interfaces.ProductRepository, logger zerolog.Logger) *Service {
    return &Service{
        repo:   repo,
        logger: logger.With().Str("service", "product").Logger(),
    }
}

// Create creates a new product with business validation
func (s *Service) Create(ctx context.Context, name, description, sku string, price float64, stock int) (*models.Product, error) {
    // Business validation
    if price <= 0 {
        return nil, domain.NewValidationError("Price must be greater than 0", "INVALID_PRICE", nil)
    }
    if stock < 0 {
        return nil, domain.NewValidationError("Stock cannot be negative", "INVALID_STOCK", nil)
    }

    // Check SKU uniqueness (business rule)
    existing, err := s.repo.FindBySKU(ctx, sku)
    if err == nil && existing != nil {
        return nil, domain.NewConflictError("Product with this SKU already exists", "SKU_EXISTS", nil)
    }

    product := &models.Product{
        Name:        name,
        Description: description,
        SKU:         sku,
        Price:       price,
        Stock:       stock,
        Active:      true,
    }

    if err := s.repo.Create(ctx, product); err != nil {
        s.logger.Error().Err(err).Str("sku", sku).Msg("Failed to create product")
        return nil, fmt.Errorf("failed to create product: %w", err)
    }

    s.logger.Info().
        Uint("product_id", product.ID).
        Str("sku", sku).
        Msg("Product created successfully")

    return product, nil
}

// GetByID retrieves a product by its ID
func (s *Service) GetByID(ctx context.Context, id uint) (*models.Product, error) {
    return s.repo.FindByID(ctx, id)
}

// GetAll retrieves all products with pagination
func (s *Service) GetAll(ctx context.Context, page, pageSize int) ([]*models.Product, int64, error) {
    // Validate and set defaults for pagination
    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 20 // Default page size
    }

    offset := (page - 1) * pageSize

    products, err := s.repo.FindAll(ctx, pageSize, offset)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to fetch products: %w", err)
    }

    total, err := s.repo.Count(ctx)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count products: %w", err)
    }

    return products, total, nil
}

// Update updates an existing product
func (s *Service) Update(ctx context.Context, id uint, name, description string, price float64, stock int, active bool) (*models.Product, error) {
    // Fetch existing product
    product, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // Business validation
    if price <= 0 {
        return nil, domain.NewValidationError("Price must be greater than 0", "INVALID_PRICE", nil)
    }
    if stock < 0 {
        return nil, domain.NewValidationError("Stock cannot be negative", "INVALID_STOCK", nil)
    }

    // Update fields
    product.Name = name
    product.Description = description
    product.Price = price
    product.Stock = stock
    product.Active = active

    if err := s.repo.Update(ctx, product); err != nil {
        s.logger.Error().Err(err).Uint("product_id", id).Msg("Failed to update product")
        return nil, fmt.Errorf("failed to update product: %w", err)
    }

    s.logger.Info().Uint("product_id", id).Msg("Product updated successfully")
    return product, nil
}

// Delete soft-deletes a product
func (s *Service) Delete(ctx context.Context, id uint) error {
    // Verify product exists
    _, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return err
    }

    if err := s.repo.Delete(ctx, id); err != nil {
        s.logger.Error().Err(err).Uint("product_id", id).Msg("Failed to delete product")
        return fmt.Errorf("failed to delete product: %w", err)
    }

    s.logger.Info().Uint("product_id", id).Msg("Product deleted successfully")
    return nil
}

// UpdateStock adjusts the stock quantity (positive or negative)
func (s *Service) UpdateStock(ctx context.Context, id uint, quantity int) error {
    product, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return err
    }

    newStock := product.Stock + quantity
    if newStock < 0 {
        return domain.NewValidationError("Insufficient stock", "INSUFFICIENT_STOCK", nil)
    }

    product.Stock = newStock
    if err := s.repo.Update(ctx, product); err != nil {
        return fmt.Errorf("failed to update stock: %w", err)
    }

    s.logger.Info().
        Uint("product_id", id).
        Int("quantity_change", quantity).
        Int("new_stock", newStock).
        Msg("Stock updated")

    return nil
}
```

**Points cles :**
- Toute la **logique metier** est ici (validation, regles business)
- Utilise les `DomainError` pour les erreurs metier
- Logging structure avec contexte
- Le service ne connait que les **interfaces**, pas les implementations

---

#### Etape 5 : Creer le Module fx (Dependency Injection)

**Fichier a creer : `internal/domain/product/module.go`**

```go
package product

import (
    "go.uber.org/fx"

    "mon-projet/internal/adapters/handlers"
    "mon-projet/internal/adapters/repository"
    "mon-projet/internal/interfaces"
)

// Module provides all product-related dependencies
var Module = fx.Module("product",
    fx.Provide(
        // Repository: concrete -> interface
        fx.Annotate(
            repository.NewProductRepository,
            fx.As(new(interfaces.ProductRepository)),
        ),
        // Service: concrete -> interface
        fx.Annotate(
            NewService,
            fx.As(new(interfaces.ProductService)),
        ),
        // Handler
        handlers.NewProductHandler,
    ),
)
```

**Pourquoi fx.Annotate ?**
- Permet de fournir une implementation concrete tout en exposant l'interface
- Facilite le remplacement des implementations (tests, mock, autre DB)

---

#### Etape 6 : Creer le Handler (HTTP Adapter)

**Fichier a creer : `internal/adapters/handlers/product_handler.go`**

```go
package handlers

import (
    "strconv"

    "github.com/go-playground/validator/v10"
    "github.com/gofiber/fiber/v2"

    "mon-projet/internal/domain"
    "mon-projet/internal/interfaces"
)

// ProductHandler handles HTTP requests for products
type ProductHandler struct {
    service  interfaces.ProductService
    validate *validator.Validate
}

// NewProductHandler creates a new product handler
func NewProductHandler(service interfaces.ProductService) *ProductHandler {
    return &ProductHandler{
        service:  service,
        validate: validator.New(),
    }
}

// Request DTOs with validation tags
type CreateProductRequest struct {
    Name        string  `json:"name" validate:"required,max=255"`
    Description string  `json:"description" validate:"max=1000"`
    SKU         string  `json:"sku" validate:"required,max=100"`
    Price       float64 `json:"price" validate:"required,gt=0"`
    Stock       int     `json:"stock" validate:"gte=0"`
}

type UpdateProductRequest struct {
    Name        string  `json:"name" validate:"required,max=255"`
    Description string  `json:"description" validate:"max=1000"`
    Price       float64 `json:"price" validate:"required,gt=0"`
    Stock       int     `json:"stock" validate:"gte=0"`
    Active      bool    `json:"active"`
}

// Create handles POST /api/v1/products
func (h *ProductHandler) Create(c *fiber.Ctx) error {
    var req CreateProductRequest
    if err := c.BodyParser(&req); err != nil {
        return domain.NewValidationError("Invalid request body", "INVALID_BODY", err)
    }

    if err := h.validate.Struct(req); err != nil {
        return domain.NewValidationError("Validation failed", "VALIDATION_ERROR", err)
    }

    product, err := h.service.Create(
        c.Context(),
        req.Name,
        req.Description,
        req.SKU,
        req.Price,
        req.Stock,
    )
    if err != nil {
        return err
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "status": "success",
        "data":   product.ToResponse(),
    })
}

// GetByID handles GET /api/v1/products/:id
func (h *ProductHandler) GetByID(c *fiber.Ctx) error {
    id, err := strconv.ParseUint(c.Params("id"), 10, 32)
    if err != nil {
        return domain.NewValidationError("Invalid product ID", "INVALID_ID", err)
    }

    product, err := h.service.GetByID(c.Context(), uint(id))
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "status": "success",
        "data":   product.ToResponse(),
    })
}

// List handles GET /api/v1/products
func (h *ProductHandler) List(c *fiber.Ctx) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

    products, total, err := h.service.GetAll(c.Context(), page, pageSize)
    if err != nil {
        return err
    }

    // Convert to response DTOs
    responses := make([]interface{}, len(products))
    for i, p := range products {
        responses[i] = p.ToResponse()
    }

    totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

    return c.JSON(fiber.Map{
        "status": "success",
        "data":   responses,
        "meta": fiber.Map{
            "page":        page,
            "page_size":   pageSize,
            "total":       total,
            "total_pages": totalPages,
        },
    })
}

// Update handles PUT /api/v1/products/:id
func (h *ProductHandler) Update(c *fiber.Ctx) error {
    id, err := strconv.ParseUint(c.Params("id"), 10, 32)
    if err != nil {
        return domain.NewValidationError("Invalid product ID", "INVALID_ID", err)
    }

    var req UpdateProductRequest
    if err := c.BodyParser(&req); err != nil {
        return domain.NewValidationError("Invalid request body", "INVALID_BODY", err)
    }

    if err := h.validate.Struct(req); err != nil {
        return domain.NewValidationError("Validation failed", "VALIDATION_ERROR", err)
    }

    product, err := h.service.Update(
        c.Context(),
        uint(id),
        req.Name,
        req.Description,
        req.Price,
        req.Stock,
        req.Active,
    )
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "status": "success",
        "data":   product.ToResponse(),
    })
}

// Delete handles DELETE /api/v1/products/:id
func (h *ProductHandler) Delete(c *fiber.Ctx) error {
    id, err := strconv.ParseUint(c.Params("id"), 10, 32)
    if err != nil {
        return domain.NewValidationError("Invalid product ID", "INVALID_ID", err)
    }

    if err := h.service.Delete(c.Context(), uint(id)); err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "status":  "success",
        "message": "Product deleted successfully",
    })
}
```

**Points cles :**
- Utilise l'interface `ProductService`, pas l'implementation concrete
- Validation avec `go-playground/validator`
- Retourne des DTOs (`ToResponse()`) au lieu des entites directement
- Gestion propre des erreurs avec `DomainError`

---

#### Etape 7 : Ajouter les Routes

**Modifier : `internal/adapters/http/routes.go`**

Ajoutez le parametre `productHandler` et les routes :

```go
func RegisterRoutes(
    app *fiber.App,
    authHandler *handlers.AuthHandler,
    userHandler *handlers.UserHandler,
    productHandler *handlers.ProductHandler,  // <- AJOUTER
    authMiddleware fiber.Handler,
) {
    // Health & Swagger
    RegisterHealthRoutes(app)
    app.Get("/swagger/*", swagger.WrapHandler)

    // API v1
    api := app.Group("/api")
    v1 := api.Group("/v1")

    // Auth routes (public)
    auth := v1.Group("/auth")
    auth.Post("/register", authHandler.Register)
    auth.Post("/login", authHandler.Login)
    auth.Post("/refresh", authHandler.Refresh)

    // User routes (protected)
    users := v1.Group("/users", authMiddleware)
    users.Get("/me", userHandler.GetMe)
    users.Get("", userHandler.GetAllUsers)
    users.Put("/:id", userHandler.UpdateUser)
    users.Delete("/:id", userHandler.DeleteUser)

    // ============================================
    // AJOUTER : Product routes (protected)
    // ============================================
    products := v1.Group("/products", authMiddleware)
    products.Post("", productHandler.Create)
    products.Get("", productHandler.List)
    products.Get("/:id", productHandler.GetByID)
    products.Put("/:id", productHandler.Update)
    products.Delete("/:id", productHandler.Delete)
}
```

**Avantages de cette approche**:
- Toutes les routes sont visibles en un seul fichier
- Facile d'ajouter de nouveaux domaines
- Le versioning de l'API est géré de manière centralisée

---

#### Etape 8 : Ajouter la Migration

**Modifier : `internal/infrastructure/database/database.go`**

```go
func NewDatabase(config *config.Config, logger zerolog.Logger) (*gorm.DB, error) {
    // ... code existant ...

    // AutoMigrate - AJOUTER models.Product
    if err := db.AutoMigrate(
        &models.User{},
        &models.RefreshToken{},
        &models.Product{},  // <- AJOUTER
    ); err != nil {
        return nil, fmt.Errorf("failed to auto-migrate: %w", err)
    }

    // ... reste du code ...
}
```

---

#### Etape 9 : Enregistrer le Module dans le Bootstrap

**Modifier : `cmd/main.go`**

```go
package main

import (
    "go.uber.org/fx"

    "mon-projet/internal/domain/product"  // <- AJOUTER
    "mon-projet/internal/domain/user"
    "mon-projet/internal/infrastructure/database"
    "mon-projet/internal/infrastructure/server"
    "mon-projet/pkg/auth"
    "mon-projet/pkg/config"
    "mon-projet/pkg/logger"
)

func main() {
    fx.New(
        logger.Module,
        config.Module,
        database.Module,
        auth.Module,
        user.Module,
        product.Module,  // <- AJOUTER
        server.Module,
    ).Run()
}
```

---

#### Verification finale

```bash
# 1. Verifier la compilation
go build ./...

# 2. Lancer l'application
make run

# 3. S'authentifier pour obtenir un token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' \
  | jq -r '.data.access_token')

# 4. Creer un produit
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "MacBook Pro 14",
    "description": "Apple laptop with M3 chip",
    "sku": "APPLE-MBP14-M3",
    "price": 1999.99,
    "stock": 50
  }'

# 5. Lister les produits
curl -X GET "http://localhost:8080/api/v1/products?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN"

# 6. Obtenir un produit par ID
curl -X GET http://localhost:8080/api/v1/products/1 \
  -H "Authorization: Bearer $TOKEN"

# 7. Mettre a jour un produit
curl -X PUT http://localhost:8080/api/v1/products/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "MacBook Pro 14 - Updated",
    "description": "Apple laptop with M3 Pro chip",
    "price": 2499.99,
    "stock": 30,
    "active": true
  }'

# 8. Supprimer un produit
curl -X DELETE http://localhost:8080/api/v1/products/1 \
  -H "Authorization: Bearer $TOKEN"
```

---

### Resume : Fichiers crees/modifies

| Action | Fichier | Description |
|--------|---------|-------------|
| Creer | `internal/models/product.go` | Entite GORM + DTO |
| Creer | `internal/interfaces/product_repository.go` | Interfaces (Ports) |
| Creer | `internal/adapters/repository/product_repository.go` | Implementation GORM |
| Creer | `internal/domain/product/service.go` | Business Logic |
| Creer | `internal/domain/product/module.go` | fx.Module |
| Creer | `internal/adapters/handlers/product_handler.go` | HTTP Handler |
| Modifier | `internal/infrastructure/server/server.go` | Ajouter routes |
| Modifier | `internal/infrastructure/database/database.go` | Ajouter AutoMigrate |
| Modifier | `cmd/main.go` | Ajouter product.Module |

### Patterns à suivre

#### 1. Error Handling

Utiliser les DomainError:

```go
// Dans service
if user == nil {
    return domain.NewNotFoundError("User not found", "USER_NOT_FOUND", nil)
}

if exists {
    return domain.NewConflictError("Email already exists", "EMAIL_EXISTS", nil)
}

// Validation
if err := validate.Struct(req); err != nil {
    return domain.NewValidationError("Invalid input", "VALIDATION_ERROR", err)
}
```

Le middleware error_handler convertit automatiquement en réponses HTTP.

#### 2. Repository Pattern

```go
// Interface (port)
type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    FindByEmail(ctx context.Context, email string) (*models.User, error)
}

// Implémentation (adapter)
type userRepositoryGORM struct {
    db *gorm.DB
}
```

#### 3. Dependency Injection avec fx

```go
// Provider
fx.Provide(
    fx.Annotate(
        NewUserService,
        fx.As(new(interfaces.UserService)),  // Interface
    ),
)

// Consumer
type AuthHandler struct {
    userService interfaces.UserService  // Dépend de l'interface
}
```

#### 4. Middleware Chain

```go
protected := api.Group("/users")
protected.Use(authMiddleware.Authenticate())  // JWT required
protected.Get("/", userHandler.List)
```

---


---

## Navigation

**Previous**: [Développement](development.md)  
**Next**: [API Reference](api-reference.md)  
**Index**: [Guide Index](index.md)
