# Project Structure

<div class="navigation">
  <a href="index.md"><i class="material-icons">arrow_back</i> Guide Index</a>
</div>

---

## Tech Stack

#### Web Framework: Fiber v2

**Why Fiber?**

- Exceptional performance (built on fasthttp)
- Familiar API (inspired by Express.js)
- Rich middleware
- Excellent documentation

**Configuration**: `internal/infrastructure/server/server.go`

```go
app := fiber.New(fiber.Config{
    ErrorHandler: errorHandler.Handle,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 10 * time.Second,
})
```

**Routes**: Centralized in `internal/adapters/http/routes.go`

```go
// routes.go - All application routes
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

**Advantages of centralized routes**:
- Overview of all API routes in a single file
- Facilitates API documentation and versioning
- Clear separation between route definition and handler logic


#### ORM: GORM

**Why GORM?**

- Most popular ORM in Go
- Automatic migrations
- Hooks and callbacks
- Associations and preloading
- Raw SQL when needed

**Configuration**: `internal/infrastructure/database/database.go`

```go
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})

// Auto-migration
db.AutoMigrate(&models.User{}, &models.RefreshToken{})
```

**Patterns used**:

- Repository pattern for isolation
- Soft deletes (DeletedAt)
- Automatic timestamps
- Indexes on foreign keys

#### Dependency Injection: uber-go/fx

**Why fx?**

- Clean dependency management
- Lifecycle hooks (OnStart, OnStop)
- Startup parallelization
- Clear compile-time errors

**Module Pattern**: Each package exposes an fx module

```go
// domain/user/module.go
var Module = fx.Module("user",
    fx.Provide(
        NewService,       // Provides UserService
        NewUserHandler,   // Provides UserHandler
        NewAuthHandler,   // Provides AuthHandler
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

**Why zerolog?**

- Structured logging (JSON)
- Optimal performance (zero-allocation)
- Log levels (Debug, Info, Warn, Error, Fatal)
- Rich context

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

**HTTP request validation**:

```go
type RegisterRequest struct {
    Email    string `json:"email" validate:"required,email,max=255"`
    Password string `json:"password" validate:"required,min=8,max=72"`
}

// In the handler
if err := validate.Struct(req); err != nil {
    // Return validation error
}
```

**Available tags**: required, email, min, max, uuid, url, alpha, numeric, etc.

#### Authentication: JWT (golang-jwt/jwt)

**Complete flow**:

1. **Register/Login** → Server generates Access Token (15min) + Refresh Token (7d)
2. **Client** → Stores the tokens, uses Access Token for each request
3. **Access Token expires** → Client sends Refresh Token
4. **Server** → Validates Refresh Token, generates new Access Token
5. **Refresh Token expires** → Client must re-login

**Token generation**:

```go
// Access token (short-lived)
accessToken, err := jwt.GenerateAccessToken(userID, jwtSecret, 15*time.Minute)

// Refresh token (long-lived)
refreshToken, err := jwt.GenerateRefreshToken(userID, jwtSecret, 7*24*time.Hour)
```

**Validation**:

```go
claims, err := jwt.ParseToken(tokenString, jwtSecret)
userID := claims.UserID
```

### Detailed Directory Structure

#### `/cmd/main.go`

**Role**: Application bootstrap.

**Contents**:

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

**Principle**: Module composition, no business logic.

#### `/internal/models`

**Models**: Shared domain entities used throughout the application.

**Role**: Centralize data structure definitions (entities) to avoid circular dependencies.

##### `user.go`

Defines the User, RefreshToken and AuthResponse entities:

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

**Principles**:

- **GORM entities**: GORM tags for database configuration
- **JSON serialization**: JSON tags to control the API (e.g.: `json:"-"` hides PasswordHash)
- **Utility methods**: IsExpired(), IsRevoked() for validation logic
- **No dependencies**: No imports from domain or interfaces
- **Usable everywhere**: Imported by interfaces, domain, repository, handlers

**Why a separate package?**

- **Avoids cycles**: Before, `interfaces` → `domain/user` → `interfaces` (<i class="material-icons error">error</i> cycle!)
- **Now**: `interfaces` → `models` ← `domain/user` (<i class="material-icons success">check_circle</i> no cycle)
- **Clarity**: Separation between entities (models) and business logic (domain)

#### `/internal/domain`

**Domain**: Pure business logic, independent of infrastructure.

##### `errors.go`

Defines custom business errors:

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

Business logic:

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

**Responsibilities**:

- Business validation
- Password hashing (Register)
- Password verification (Login)
- Repository call orchestration
- **Uses `models.User`**: Imports the models package for entities

#### `/internal/adapters`

**Adapters**: Connect the domain to the outside world.

##### `handlers/auth_handler.go`

Authentication endpoints:

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

1. Parse JSON body
2. Validate with validator
3. Call service
4. Generate tokens (for Login/Register)
5. Return response

**Register example**:

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

Verifies the JWT token:

```go
type AuthMiddleware struct {
    jwtSecret string
}

func (m *AuthMiddleware) Authenticate() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Extract token from Authorization header
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return fiber.NewError(fiber.StatusUnauthorized, "Missing authorization header")
        }

        // Validate "Bearer <token>" format
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            return fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization format")
        }

        // Parse and validate the token
        claims, err := auth.ParseToken(parts[1], m.jwtSecret)
        if err != nil {
            return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
        }

        // Inject user ID into the context
        c.Locals("user_id", claims.UserID)

        return c.Next()
    }
}
```

##### `middleware/error_handler.go`

Centralized error handling:

```go
func (h *ErrorHandler) Handle(c *fiber.Ctx, err error) error {
    // DomainError → appropriate HTTP status
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

    // Generic error
    return c.Status(fiber.StatusInternalServerError).JSON(...)
}
```

**Advantage**: Handlers don't need to manage HTTP status codes, just return DomainErrors.

##### `repository/user_repository.go`

Repository implementation with GORM:

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

**Infrastructure**: DB and server configuration.

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

The server creates the Fiber application and manages the lifecycle. Routes are registered via `server.Module` which invokes `httpRoutes.RegisterRoutes()` with `fx.Invoke`.

```go
// Module provides the Fiber server dependency via fx
var Module = fx.Module("server",
    fx.Provide(NewServer),
    fx.Invoke(registerHooks),
    fx.Invoke(httpRoutes.RegisterRoutes),  // Centralized routes
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

Centralized file for all application routes:

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

**Reusable packages**: Can be imported by other projects.

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

**Previous**: [Hexagonal Architecture](architecture.md)  
**Next**: [Configuration](configuration.md)  
**Index**: [Guide Index](index.md)
