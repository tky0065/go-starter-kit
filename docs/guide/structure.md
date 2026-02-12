# Structure du projet

<div class="navigation">
  <a href="index.md"><i class="material-icons">arrow_back</i> Guide Index</a>
</div>

---

## Stack technique

#### Web Framework: Fiber v2

**Pourquoi Fiber?**

- Performance exceptionnelle (built on fasthttp)
- API familière (inspirée d'Express.js)
- Middleware riche
- Documentation excellente

**Configuration**: `internal/infrastructure/server/server.go`

```go
app := fiber.New(fiber.Config{
    ErrorHandler: errorHandler.Handle,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 10 * time.Second,
})
```

**Routes**: Centralisées dans `internal/adapters/http/routes.go`

```go
// routes.go - Toutes les routes de l'application
func RegisterRoutes(
    app *fiber.App,
    authHandler *handlers.AuthHandler,
    userHandler *handlers.UserHandler,
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
}
```

**Avantages de la centralisation des routes**:
- Vue d'ensemble de toutes les routes API en un seul fichier
- Facilite la documentation et le versioning de l'API
- Séparation claire entre la définition des routes et la logique des handlers


#### ORM: GORM

**Pourquoi GORM?**

- ORM le plus populaire en Go
- Migrations automatiques
- Hooks et callbacks
- Associations et preloading
- Raw SQL quand nécessaire

**Configuration**: `internal/infrastructure/database/database.go`

```go
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})

// Auto-migration
db.AutoMigrate(&models.User{}, &models.RefreshToken{})
```

**Patterns utilisés**:

- Repository pattern pour isolation
- Soft deletes (DeletedAt)
- Timestamps automatiques
- Indexes sur clés étrangères

#### Dependency Injection: uber-go/fx

**Pourquoi fx?**

- Gestion propre des dépendances
- Lifecycle hooks (OnStart, OnStop)
- Parallélisation du démarrage
- Erreurs claires à la compilation

**Pattern Module**: Chaque package expose un module fx

```go
// domain/user/module.go
var Module = fx.Module("user",
    fx.Provide(
        NewService,       // Fournit UserService
        NewUserHandler,   // Fournit UserHandler
        NewAuthHandler,   // Fournit AuthHandler
    ),
)
```

**Bootstrap**: `cmd/main.go`

```go
fx.New(
    logger.Module,      // Logger
    config.Module,      // Configuration
    database.Module,    // Database
    auth.Module,        // JWT utilities
    user.Module,        // User domain
    server.Module,      // Fiber server
).Run()
```

#### Logging: zerolog

**Pourquoi zerolog?**

- Logging structuré (JSON)
- Performance optimale (zero-allocation)
- Niveaux de log (Debug, Info, Warn, Error, Fatal)
- Contexte riche

**Usage**:

```go
logger.Info().
    Str("email", user.Email).
    Uint("user_id", user.ID).
    Msg("User registered successfully")

logger.Error().
    Err(err).
    Str("operation", "create_user").
    Msg("Failed to create user")
```

#### Validation: go-playground/validator v10

**Validation des requêtes HTTP**:

```go
type RegisterRequest struct {
    Email    string `json:"email" validate:"required,email,max=255"`
    Password string `json:"password" validate:"required,min=8,max=72"`
}

// Dans le handler
if err := validate.Struct(req); err != nil {
    // Retourner erreur de validation
}
```

**Tags disponibles**: required, email, min, max, uuid, url, alpha, numeric, etc.

#### Authentication: JWT (golang-jwt/jwt)

**Flow complet**:

1. **Register/Login** → Serveur génère Access Token (15min) + Refresh Token (7j)
2. **Client** → Stocke les tokens, utilise Access Token pour chaque requête
3. **Access Token expire** → Client envoie Refresh Token
4. **Serveur** → Valide Refresh Token, génère nouveau Access Token
5. **Refresh Token expire** → Client doit se re-login

**Génération de tokens**:

```go
// Access token (courte durée)
accessToken, err := jwt.GenerateAccessToken(userID, jwtSecret, 15*time.Minute)

// Refresh token (longue durée)
refreshToken, err := jwt.GenerateRefreshToken(userID, jwtSecret, 7*24*time.Hour)
```

**Validation**:

```go
claims, err := jwt.ParseToken(tokenString, jwtSecret)
userID := claims.UserID
```

### Structure des répertoires détaillée

#### `/cmd/main.go`

**Rôle**: Bootstrap de l'application.

**Contenu**:

```go
package main

import (
    "go.uber.org/fx"
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
        server.Module,
    ).Run()
}
```

**Principe**: Composition de modules, pas de logique métier.

#### `/internal/models`

**Models**: Entités de domaine partagées utilisées à travers toute l'application.

**Rôle**: Centraliser les définitions des structures de données (entities) pour éviter les dépendances circulaires.

##### `user.go`

Définit les entités User, RefreshToken et AuthResponse:

```go
package models

import (
    "time"
    "gorm.io/gorm"
)

// User represents the domain entity for a user
type User struct {
    ID           uint           `gorm:"primaryKey" json:"id"`
    Email        string         `gorm:"uniqueIndex;not null" json:"email"`
    PasswordHash string         `gorm:"not null" json:"-"`
    CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// RefreshToken represents a refresh token for session management
type RefreshToken struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    UserID    uint      `gorm:"not null;index" json:"user_id"`
    Token     string    `gorm:"uniqueIndex;not null" json:"token"`
    ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
    Revoked   bool      `gorm:"not null;default:false" json:"revoked"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (rt *RefreshToken) IsExpired() bool {
    return time.Now().After(rt.ExpiresAt)
}

func (rt *RefreshToken) IsRevoked() bool {
    return rt.Revoked
}

// AuthResponse represents the authentication response with tokens
type AuthResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    int64  `json:"expires_in"`
}
```

**Principes**:

- **Entités GORM**: Tags GORM pour configuration base de données
- **Serialization JSON**: Tags json pour contrôler l'API (ex: `json:"-"` cache PasswordHash)
- **Méthodes utilitaires**: IsExpired(), IsRevoked() pour la logique de validation
- **Pas de dépendances**: Aucun import de domain ou interfaces
- **Utilisable partout**: Importé par interfaces, domain, repository, handlers

**Pourquoi un package séparé?**

- **Évite les cycles**: Avant, `interfaces` → `domain/user` → `interfaces` (❌ cycle!)
- **Maintenant**: `interfaces` → `models` ← `domain/user` (:material-check-circle: pas de cycle)
- **Clarté**: Séparation entre entities (models) et business logic (domain)

#### `/internal/domain`

**Domaine**: Logique métier pure, indépendante de l'infrastructure.

##### `errors.go`

Définit les erreurs métier personnalisées:

```go
type DomainError struct {
    Type    string
    Message string
    Code    string
    Err     error
}

func NewNotFoundError(message, code string, err error) *DomainError
func NewValidationError(message, code string, err error) *DomainError
func NewConflictError(message, code string, err error) *DomainError
func NewUnauthorizedError(message, code string, err error) *DomainError
```

**Usage**:

```go
if user == nil {
    return domain.NewNotFoundError("User not found", "USER_NOT_FOUND", nil)
}
```

##### `user/service.go`

Logique métier:

```go
package user

import (
    "context"
    "mon-projet/internal/models"
    "mon-projet/internal/interfaces"
)

type Service struct {
    repo   interfaces.UserRepository
    logger zerolog.Logger
}

func (s *Service) Register(ctx context.Context, email, password string) (*models.User, error)
func (s *Service) Login(ctx context.Context, email, password string) (*models.User, error)
func (s *Service) GetByID(ctx context.Context, id uint) (*models.User, error)
func (s *Service) Update(ctx context.Context, id uint, email string) (*models.User, error)
func (s *Service) Delete(ctx context.Context, id uint) error
```

**Responsabilités**:

- Validation métier
- Hashage de password (Register)
- Vérification de password (Login)
- Orchestration d'appels repository
- **Utilise `models.User`**: Importe le package models pour les entités

#### `/internal/adapters`

**Adapters**: Connectent le domaine au monde extérieur.

##### `handlers/auth_handler.go`

Endpoints d'authentification:

```go
type AuthHandler struct {
    authService interfaces.AuthService
    userService interfaces.UserService
    jwtSecret   string
    validate    *validator.Validate
}

func (h *AuthHandler) Register(c *fiber.Ctx) error
func (h *AuthHandler) Login(c *fiber.Ctx) error
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error
```

**Pattern**:

1. Parse body JSON
2. Validate avec validator
3. Appeler service
4. Générer tokens (pour Login/Register)
5. Retourner réponse

**Exemple Register**:

```go
func (h *AuthHandler) Register(c *fiber.Ctx) error {
    var req RegisterRequest
    if err := c.BodyParser(&req); err != nil {
        return err
    }

    if err := h.validate.Struct(req); err != nil {
        return domain.NewValidationError("Invalid input", "VALIDATION_ERROR", err)
    }

    user, err := h.userService.Register(c.Context(), req.Email, req.Password)
    if err != nil {
        return err
    }

    accessToken, _ := auth.GenerateAccessToken(user.ID, h.jwtSecret, 15*time.Minute)
    refreshToken, _ := auth.GenerateRefreshToken(user.ID, h.jwtSecret, 168*time.Hour)

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "status": "success",
        "data": fiber.Map{
            "access_token":  accessToken,
            "refresh_token": refreshToken,
            "token_type":    "Bearer",
            "expires_in":    900,
        },
    })
}
```

##### `middleware/auth_middleware.go`

Vérifie le JWT token:

```go
type AuthMiddleware struct {
    jwtSecret string
}

func (m *AuthMiddleware) Authenticate() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Extraire token du header Authorization
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return fiber.NewError(fiber.StatusUnauthorized, "Missing authorization header")
        }

        // Valider format "Bearer <token>"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            return fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization format")
        }

        // Parser et valider le token
        claims, err := auth.ParseToken(parts[1], m.jwtSecret)
        if err != nil {
            return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
        }

        // Injecter user ID dans le contexte
        c.Locals("user_id", claims.UserID)

        return c.Next()
    }
}
```

##### `middleware/error_handler.go`

Gestion centralisée des erreurs:

```go
func (h *ErrorHandler) Handle(c *fiber.Ctx, err error) error {
    // DomainError → HTTP status approprié
    if domainErr, ok := err.(*domain.DomainError); ok {
        switch domainErr.Type {
        case "not_found":
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "status": "error",
                "error":  domainErr.Message,
                "code":   domainErr.Code,
            })
        case "validation":
            return c.Status(fiber.StatusBadRequest).JSON(...)
        case "unauthorized":
            return c.Status(fiber.StatusUnauthorized).JSON(...)
        case "conflict":
            return c.Status(fiber.StatusConflict).JSON(...)
        }
    }

    // Erreur générique
    return c.Status(fiber.StatusInternalServerError).JSON(...)
}
```

**Avantage**: Les handlers n'ont pas besoin de gérer les status HTTP, juste retourner des DomainError.

##### `repository/user_repository.go`

Implémentation du repository avec GORM:

```go
type userRepositoryGORM struct {
    db *gorm.DB
}

func (r *userRepositoryGORM) Create(ctx context.Context, user *models.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepositoryGORM) FindByEmail(ctx context.Context, email string) (*models.User, error) {
    var user models.User
    err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
    if err == gorm.ErrRecordNotFound {
        return nil, domain.NewNotFoundError("User not found", "USER_NOT_FOUND", err)
    }
    return &user, err
}
```

#### `/internal/infrastructure`

**Infrastructure**: Configuration DB et serveur.

##### `database/database.go`

```go
func NewDatabase(config *config.Config, logger zerolog.Logger) (*gorm.DB, error) {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
        config.DBHost, config.DBUser, config.DBPassword,
        config.DBName, config.DBPort, config.DBSSLMode)

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

    // AutoMigrate
    db.AutoMigrate(&models.User{}, &models.RefreshToken{})

    return db, nil
}
```

##### `server/server.go`

Le serveur crée l'application Fiber et gère le lifecycle. Les routes sont enregistrées via `server.Module` qui invoque `httpRoutes.RegisterRoutes()` avec `fx.Invoke`.

```go
// Module provides the Fiber server dependency via fx
var Module = fx.Module("server",
    fx.Provide(NewServer),
    fx.Invoke(registerHooks),
    fx.Invoke(httpRoutes.RegisterRoutes),  // Routes centralisées
)

func NewServer(logger zerolog.Logger, db *gorm.DB) *fiber.App {
    app := fiber.New(fiber.Config{
        AppName:      "mon-projet",
        ErrorHandler: middleware.ErrorHandler,
    })

    logger.Info().Msg("Fiber server initialized with centralized error handler")

    return app
}

// registerHooks registers lifecycle hooks for server startup and shutdown
func registerHooks(lifecycle fx.Lifecycle, app *fiber.App, logger zerolog.Logger) {
    lifecycle.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            port := config.GetEnv("APP_PORT", "8080")
            logger.Info().Str("port", port).Msg("Starting Fiber server")

            go func() {
                if err := app.Listen(":" + port); err != nil {
                    logger.Error().Err(err).Msg("Server stopped unexpectedly")
                }
            }()

            return nil
        },
        OnStop: func(ctx context.Context) error {
            logger.Info().Msg("Shutting down Fiber server gracefully")
            return app.ShutdownWithContext(ctx)
        },
    })
}
```

##### `http/routes.go`

Fichier centralisé pour toutes les routes de l'application :

```go
func RegisterRoutes(
    app *fiber.App,
    authHandler *handlers.AuthHandler,
    userHandler *handlers.UserHandler,
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
}
```

#### `/pkg`

**Packages réutilisables**: Peuvent être importés par d'autres projets.

##### `auth/jwt.go`

```go
func GenerateAccessToken(userID uint, secret string, expiry time.Duration) (string, error) {
    claims := &Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}
```

---



---

## Navigation

**Previous**: [Architecture hexagonale](architecture.md)  
**Next**: [Configuration](configuration.md)  
**Index**: [Guide Index](index.md)
