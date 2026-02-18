# Generated Project Guide

> Comprehensive guide to developing, testing, and deploying projects created with create-go-starter

## Table of Contents

1. [Architecture](#architecture)
2. [Configuration](#configuration)
3. [Development](#development)
4. [API Reference](#api-reference)
5. [Tests](#tests)
6. [Database](#database)
7. [Security](#security)
8. [Deployment](#deployment)
9. [Monitoring & Logging](#monitoring--logging)
10. [Best Practices](#best-practices)

---

## Architecture

### Hexagonal Architecture (Ports & Adapters)

Generated projects follow hexagonal architecture, also known as "Ports and Adapters".

**Core principle**: The business domain (business logic) is at the center and depends on nothing. All dependencies point toward the domain.

```
┌─────────────────────────────────────────────────────────┐
│                   HTTP Layer (Fiber)                    │
│              adapters/handlers + middleware             │
│  • AuthHandler (register, login, refresh)               │
│  • UserHandler (CRUD operations)                        │
│  • AuthMiddleware (JWT verification)                    │
│  • ErrorHandler (centralized error handling)            │
└───────────────────────┬─────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────────┐
│              Shared Entities Layer                      │
│                   models/                               │
│  • User (entity with GORM tags)                         │
│  • RefreshToken (entity with GORM tags)                 │
│  • AuthResponse (DTO)                                   │
└──────────┬───────────────────────────┬──────────────────┘
           │                           │
           ▼                           ▼
┌──────────────────────┐  ┌──────────────────────────────┐
│   Interfaces Layer   │  │      Domain Layer            │
│   interfaces/        │  │      domain/user             │
│  • UserRepository    │  │  • UserService (logic)       │
│    (port)            │  │  • Business rules            │
└──────────┬───────────┘  └──────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────────────┐
│              Infrastructure Layer                       │
│         database + repository + server                  │
│  • GORM Database Connection                             │
│  • UserRepository (GORM implementation)                 │
│  • Fiber Server Configuration                           │
└─────────────────────────────────────────────────────────┘
```

### Complete Architecture Diagram (Mermaid)

The following diagram shows the complete hexagonal architecture with all components and their interactions:

```mermaid
flowchart TB
    subgraph External["External World"]
        Client["HTTP Client<br/>(Web, Mobile, API)"]
        DB[("PostgreSQL<br/>Database")]
    end

    subgraph Adapters["Adapters Layer"]
        direction TB
        subgraph Inbound["Inbound Adapters (Input)"]
            Handlers["Handlers<br/>AuthHandler<br/>UserHandler"]
            Middleware["Middleware<br/>AuthMiddleware<br/>ErrorHandler"]
        end
        subgraph Outbound["Outbound Adapters (Output)"]
            RepoImpl["GORM Repository<br/>UserRepository"]
        end
    end

    subgraph Core["Core Business (Hexagon)"]
        direction TB
        Models["Models (Entities)<br/>User<br/>RefreshToken"]
        Domain["Domain Services<br/>UserService<br/>Business Logic"]
        Interfaces["Interfaces/Ports<br/>UserRepository<br/>UserService"]
        Errors["Domain Errors<br/>NotFound<br/>Validation<br/>Conflict"]
    end

    subgraph Infrastructure["Infrastructure Layer"]
        Server["Fiber Server<br/>Routes and Config"]
        DBConn["Database Connection<br/>GORM Setup"]
        Config["Configuration<br/>Environment vars"]
    end

    subgraph Packages["Reusable Packages (pkg/)"]
        Auth["Auth Package<br/>JWT Generation<br/>Token Parsing"]
        Logger["Logger Package<br/>Zerolog Config"]
        ConfigPkg["Config Package<br/>Env Loading"]
    end

    Client -->|"HTTP Request"| Server
    Server -->|"Route"| Handlers
    Handlers --> Middleware
    Handlers -->|"Calls"| Domain
    Domain -->|"Uses"| Interfaces
    Domain -->|"Uses"| Models
    Domain -->|"Returns"| Errors
    RepoImpl -.->|"Implements"| Interfaces
    RepoImpl -->|"Uses"| Models
    RepoImpl -->|"Query"| DBConn
    DBConn -->|"SQL"| DB
    Handlers -->|"Uses"| Auth
    Server -->|"Uses"| Config
    Domain -->|"Uses"| Logger
```

### HTTP Request Flow (Sequence Diagram)

This diagram shows the complete path of an HTTP request through the architecture:

```mermaid
sequenceDiagram
    autonumber
    participant C as Client
    participant S as Server (Fiber)
    participant M as Middleware
    participant H as Handler
    participant SVC as Service
    participant P as Port (Interface)
    participant R as Repository
    participant DB as Database

    C->>S: POST /api/v1/auth/register
    S->>M: Route to Handler
    M->>M: Validation (if protected)
    M->>H: Validated request
    
    rect rgb(240, 248, 255)
        Note over H: Handler Layer
        H->>H: Parse JSON Body
        H->>H: Validate Input (validator)
    end
    
    H->>SVC: service.Register(email, password)
    
    rect rgb(255, 250, 240)
        Note over SVC: Domain Layer
        SVC->>SVC: Hash Password (bcrypt)
        SVC->>SVC: Business Validation
    end
    
    SVC->>P: repo.Create(user)
    P->>R: Call implementation
    
    rect rgb(240, 255, 240)
        Note over R: Repository Layer
        R->>DB: INSERT INTO users...
        DB-->>R: User created (ID)
    end
    
    R-->>SVC: User entity
    SVC-->>H: User + nil error
    H->>H: Generate JWT tokens
    H-->>C: HTTP 201 + JSON Response
```

### Dependency Inversion Principle

The core of hexagonal architecture relies on **Dependency Inversion**:

```mermaid
flowchart LR
    subgraph Traditional["Traditional Approach"]
        direction TB
        T_Handler["Handler"] --> T_Service["Service"]
        T_Service --> T_Repo["Repository"]
        T_Repo --> T_DB["Database"]
    end

    subgraph Hexagonal["Hexagonal Architecture"]
        direction TB
        H_Handler["Handler"]
        H_Service["Service"]
        H_Interface["Interface<br/>(Port)"]
        H_Repo["Repository<br/>(Adapter)"]
        H_DB["Database"]
        
        H_Handler --> H_Service
        H_Service --> H_Interface
        H_Repo -.->|"implements"| H_Interface
        H_Repo --> H_DB
    end
```

**Advantages of this approach**:

| Aspect | Without Hexagonal | With Hexagonal |
|--------|-------------------|----------------|
| **Testability** | Difficult (depends on DB) | Easy (mock interfaces) |
| **Changing DB** | Modifications everywhere | Only the repository |
| **Changing framework** | Complete refactoring | Only the handlers |
| **Business logic** | Scattered | Centralized in the domain |

### File Structure and Responsibilities

```mermaid
flowchart TD
    subgraph CMD["cmd/"]
        Main["main.go<br/>Bootstrap fx.New()"]
    end

    subgraph Internal["internal/"]
        subgraph Models["models/"]
            User["user.go<br/>GORM Entities"]
        end
        
        subgraph Domain["domain/"]
            DErrors["errors.go<br/>Business Errors"]
            subgraph UserDomain["user/"]
                Service["service.go<br/>Business Logic"]
                Module["module.go<br/>fx.Module"]
            end
        end
        
        subgraph InterfacesPkg["interfaces/"]
            Repos["*_repository.go<br/>Ports (abstractions)"]
        end
        
        subgraph AdaptersPkg["adapters/"]
            subgraph HandlersPkg["handlers/"]
                AuthH["auth_handler.go"]
                UserH["user_handler.go"]
            end
            subgraph HttpPkg["http/"]
                Health["health.go"]
                Routes["routes.go<br/>Centralized Routes"]
            end
            subgraph MiddlewarePkg["middleware/"]
                AuthM["auth_middleware.go"]
                ErrorM["error_handler.go"]
            end
            subgraph RepoPkg["repository/"]
                UserRepo["user_repository.go<br/>GORM Implementation"]
            end
        end
        
        subgraph Infra["infrastructure/"]
            DBPkg["database/<br/>GORM Connection"]
            ServerPkg["server/<br/>Fiber Config"]
        end
    end

    subgraph Pkg["pkg/"]
        AuthPkg["auth/<br/>JWT utilities"]
        ConfigPkg2["config/<br/>Env loading"]
        LoggerPkg["logger/<br/>Zerolog setup"]
    end

    Main --> Domain
    Main --> Infra
    Main --> Pkg
    HandlersPkg --> Domain
    HttpPkg --> HandlersPkg
    Domain --> InterfacesPkg
    RepoPkg -.-> InterfacesPkg
    RepoPkg --> Models
    Domain --> Models
```

**Data flow**:

1. **HTTP Request** → Handler (adapters/handlers)
2. **Handler** → Calls the Service via the interface (domain)
3. **Service** → Executes business logic, calls the Repository via the interface
4. **Repository** → Persists to DB (infrastructure)
5. **Return** → Bubbles back up to the Handler which returns the HTTP response

**Advantages**:

- **Testability**: The domain can be tested without DB or HTTP
- **Flexibility**: Changing DB (PostgreSQL → MySQL) or framework (Fiber → Gin) is easy
- **Maintainability**: Clear separation of responsibilities
- **Scalability**: Adding new features without breaking existing ones

### Technical Stack

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
- Clear separation between route definitions and handler logic


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
2. **Client** → Stores tokens, uses Access Token for each request
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

**Role**: Centralize data structure (entity) definitions to avoid circular dependencies.

##### `user.go`

Defines the User, RefreshToken, and AuthResponse entities:

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

- **GORM Entities**: GORM tags for database configuration
- **JSON Serialization**: JSON tags to control the API (e.g., `json:"-"` hides PasswordHash)
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

        // Inject user ID into context
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

## Configuration

### Environment Variables

The `.env` file contains all configuration:

```bash
# Application
APP_NAME=mon-projet
APP_ENV=development
APP_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=mon-projet
DB_SSLMODE=disable

# JWT
JWT_SECRET=          # MUST BE FILLED IN!
JWT_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=168h
```

### Generate a Secure JWT_SECRET

**CRITICAL**: Always generate a strong secret:

```bash
openssl rand -base64 32
```

Example result:
```
XqR7nP3vY2kL9wH4sT6mU8jC1bN5aD0f
```

Add it to `.env`:
```
JWT_SECRET=XqR7nP3vY2kL9wH4sT6mU8jC1bN5aD0f
```

### Configuration by Environment

#### Development

```bash
APP_ENV=development
DB_HOST=localhost
DB_SSLMODE=disable
```

#### Staging

```bash
APP_ENV=staging
DB_HOST=staging-db.example.com
DB_SSLMODE=require
JWT_SECRET=<secret-from-vault>
```

#### Production

```bash
APP_ENV=production
DB_HOST=prod-db.example.com
DB_SSLMODE=require
DB_PASSWORD=<secret-from-secrets-manager>
JWT_SECRET=<secret-from-secrets-manager>
```

**Best practice**: Use secrets managers:

- **AWS**: Secrets Manager, Parameter Store
- **GCP**: Secret Manager
- **Kubernetes**: Secrets
- **HashiCorp**: Vault

### PostgreSQL Configuration

#### Option 1: Local PostgreSQL

**macOS (Homebrew)**:

```bash
brew install postgresql@16
brew services start postgresql@16
createdb mon-projet
```

**Linux (apt)**:

```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
sudo -u postgres createdb mon-projet
```

#### Option 2: Docker

```bash
docker run -d \
  --name postgres \
  -e POSTGRES_DB=mon-projet \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine
```

#### Option 3: Docker Compose

If a `docker-compose.yml` is generated:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: mon-projet
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

Start it:

```bash
docker-compose up -d postgres
```

#### Connection Verification

```bash
# With psql
psql -h localhost -U postgres -d mon-projet

# Or test from the app
make run
# Check logs: "Database connected successfully"
```

---

## Development

### Daily Workflow

**1. Start the database**

```bash
# Docker
docker start postgres
# or
docker-compose up -d postgres

# Local
brew services start postgresql  # macOS
sudo systemctl start postgresql # Linux
```

**2. Start the application**

```bash
make run
```

Or with hot-reload (if `air` is installed):

```bash
# Install air
go install github.com/cosmtrek/air@latest

# Start with hot-reload
air
```

**3. Develop**

- Modify the code
- Save (auto-reload with air)
- Check the logs

**4. Test**

```bash
# Unit tests
make test

# Tests with coverage
make test-coverage

# Open the report
open coverage.html
```

**5. Lint**

```bash
make lint
```

### Makefile Commands

| Command | Description |
|---------|-------------|
| `make help` | Show help |
| `make run` | Start the app |
| `make build` | Build binary |
| `make test` | Tests with race detector |
| `make test-coverage` | Tests + HTML report |
| `make lint` | golangci-lint |
| `make clean` | Clean artifacts |
| `make docker-build` | Build Docker image |
| `make docker-run` | Run Docker container |

### Model Management with add-model <i class="material-icons success">new_releases</i>

**New in v1.2.0!** The CRUD scaffolding generator fully automates the creation of new models.

#### Quick Workflow

Instead of manually creating 8 files and modifying 3 existing files (see next section), use:

```bash
create-go-starter add-model <ModelName> --fields "field:type,..."
```

**Example:**
```bash
cd mon-projet  # Navigate to your existing project

# Create a complete Todo model
create-go-starter add-model Todo --fields "title:string,completed:bool,priority:int"
```

**Result:** 8 files generated + 3 files automatically updated in < 2 seconds.

#### Automatically Generated Files

| File | Role | Content |
|------|------|---------|
| `internal/models/todo.go` | Entity | Struct with GORM tags |
| `internal/interfaces/todo_repository.go` | Port | Repository interface |
| `internal/adapters/repository/todo_repository.go` | Adapter | GORM implementation |
| `internal/domain/todo/service.go` | Business Logic | CRUD operations |
| `internal/domain/todo/module.go` | fx Module | Dependency injection |
| `internal/adapters/handlers/todo_handler.go` | HTTP Adapter | REST endpoints |
| `internal/domain/todo/service_test.go` | Tests | Service unit tests |
| `internal/adapters/handlers/todo_handler_test.go` | Tests | HTTP handler tests |

#### Automatically Updated Files

| File | Modification |
|------|--------------|
| `internal/infrastructure/database/database.go` | Adds `&models.Todo{}` in AutoMigrate |
| `internal/adapters/http/routes.go` | Adds CRUD routes `/api/v1/todos/*` |
| `cmd/main.go` | Adds `todo.Module` in fx.New |

#### Types and Modifiers

**Supported field types:**
- `string`, `int`, `uint`, `float64`, `bool`, `time`

**GORM modifiers:**
- `unique` - Uniqueness constraint
- `not_null` - Required field
- `index` - Database index

**Syntax:**
```bash
--fields "field:type:modifier1:modifier2,..."
```

**Examples:**
```bash
# Unique and required email
create-go-starter add-model User --fields "email:string:unique:not_null,age:int"

# Product with price and indexed stock
create-go-starter add-model Product --fields "name:string:unique,price:float64,stock:int:index"

# Article with optional publication
create-go-starter add-model Article --fields "title:string:not_null,content:string,published:bool"
```

#### Relationships Between Models

##### BelongsTo (N:1 - child to parent)

Create a model that **belongs to** an existing parent:

```bash
# The parent MUST exist first
create-go-starter add-model Category --fields "name:string:unique"

# Create child with BelongsTo relationship
create-go-starter add-model Product --fields "name:string,price:float64" --belongs-to Category
```

**What is added in `internal/models/product.go`:**
```go
type Product struct {
    // ... custom fields
    CategoryID uint     `gorm:"not null;index" json:"category_id"`
    Category   Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}
```

**Generated nested routes:**
- `GET /api/v1/categories/:categoryId/products` - List products of a category
- `POST /api/v1/categories/:categoryId/products` - Create product in category

**Preloading:**
- `GET /api/v1/products/:id?include=category` - Product with its category

##### HasMany (1:N - parent to children)

Add a slice of children to an **existing parent model**:

```bash
# Both the parent AND child MUST exist
create-go-starter add-model Category --fields "name:string"
create-go-starter add-model Product --fields "name:string" --belongs-to Category

# Add HasMany to the parent
create-go-starter add-model Category --has-many Product
```

**What is added in `internal/models/category.go`:**
```go
type Category struct {
    // ... existing fields
    Products []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}
```

**Preloading:**
- `GET /api/v1/categories/:id?include=products` - Category with all its products

##### Nested Relationships (3+ levels)

Example: **Category** → **Post** → **Comment**

```bash
# 1. Create the root
create-go-starter add-model Category --fields "name:string:unique"

# 2. Create level 2 (child of Category)
create-go-starter add-model Post \
  --fields "title:string:not_null,content:string,published:bool" \
  --belongs-to Category

# 3. Create level 3 (child of Post)
create-go-starter add-model Comment \
  --fields "author:string:not_null,content:string:not_null" \
  --belongs-to Post

# 4. Optional: Add HasMany to parents
create-go-starter add-model Category --has-many Post
create-go-starter add-model Post --has-many Comment
```

**Result:**
- `Category` has `[]Post`
- `Post` has `CategoryID` + `Category` AND `[]Comment`
- `Comment` has `PostID` + `Post`

**Generated endpoints:**
```bash
# Standard CRUD
GET    /api/v1/categories
GET    /api/v1/posts
GET    /api/v1/comments

# Nested relationships
GET    /api/v1/categories/:categoryId/posts
POST   /api/v1/categories/:categoryId/posts
GET    /api/v1/posts/:postId/comments
POST   /api/v1/posts/:postId/comments

# Preloading
GET    /api/v1/posts/:id?include=category,comments
GET    /api/v1/categories/:id?include=posts
```

#### Public vs Protected Routes

**By default**, all routes are **protected by JWT** (middleware `auth.RequireAuth`).

To create **public routes** (without authentication):

```bash
create-go-starter add-model Article --fields "title:string,content:string" --public
```

This generates:
```go
// routes.go - NO auth middleware
api.Get("/articles", articleHandler.List)
api.Post("/articles", articleHandler.Create)  // Public!
```

<i class="material-icons warning">warning</i> **Warning:** Use `--public` with caution to avoid security vulnerabilities.

#### Customization After Generation

The generated code follows Go best practices and can be **easily extended**:

##### 1. Add custom validations

```go
// internal/domain/todo/service.go
func (s *Service) Create(ctx context.Context, todo *models.Todo) error {
    // Custom business validation
    if todo.Priority < 0 || todo.Priority > 10 {
        return domain.ErrValidation("priority must be between 0 and 10")
    }

    return s.repo.Create(ctx, todo)
}
```

##### 2. Add business methods

```go
// internal/models/todo.go
func (t *Todo) IsOverdue() bool {
    return t.DueDate.Before(time.Now()) && !t.Completed
}

func (t *Todo) MarkComplete() {
    t.Completed = true
    t.CompletedAt = time.Now()
}
```

##### 3. Add custom endpoints

```go
// internal/adapters/handlers/todo_handler.go
func (h *Handler) MarkComplete(c *fiber.Ctx) error {
    id, _ := c.ParamsInt("id")

    todo, err := h.service.GetByID(c.Context(), uint(id))
    if err != nil {
        return err
    }

    todo.MarkComplete()
    return h.service.Update(c.Context(), uint(id), todo)
}

// internal/adapters/http/routes.go
todos.Put("/:id/complete", todoHandler.MarkComplete)
```

##### 4. Add custom queries to the repository

```go
// internal/interfaces/todo_repository.go
type TodoRepository interface {
    // ... generated CRUD methods
    FindOverdue(ctx context.Context) ([]models.Todo, error)
    FindByPriority(ctx context.Context, priority int) ([]models.Todo, error)
}

// internal/adapters/repository/todo_repository.go
func (r *Repository) FindOverdue(ctx context.Context) ([]models.Todo, error) {
    var todos []models.Todo
    err := r.db.WithContext(ctx).
        Where("due_date < ? AND completed = ?", time.Now(), false).
        Find(&todos).Error
    return todos, err
}
```

#### Advanced Relationships

##### Preloading multiple relationships

```bash
# Post with Category AND Comments
GET /api/v1/posts/:id?include=category,comments

# Category with Posts, and each Post with its Comments
GET /api/v1/categories/:id?include=posts.comments
```

##### Avoiding N+1 queries

The generated code automatically uses `Preload()` to avoid N+1:

```go
// internal/adapters/repository/post_repository.go
func (r *Repository) GetByID(ctx context.Context, id uint) (*models.Post, error) {
    var post models.Post
    err := r.db.WithContext(ctx).
        Preload("Category").     // Loads the category in 1 query
        Preload("Comments").     // Loads comments in 1 query
        First(&post, id).Error
    return &post, err
}
```

#### Complete Workflow with add-model

```bash
# 1. Create initial project
create-go-starter blog-api
cd blog-api
./setup.sh

# 2. Generate models
create-go-starter add-model Category --fields "name:string:unique"
create-go-starter add-model Post --fields "title:string,content:string" --belongs-to Category
create-go-starter add-model Comment --fields "author:string,content:string" --belongs-to Post

# 3. Rebuild and test
go mod tidy
go build ./...
make test

# 4. Optional: Regenerate Swagger
make swagger

# 5. Start the server
make run

# 6. Test the API
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Content-Type: application/json" \
  -d '{"name": "Technology"}'

curl -X POST http://localhost:8080/api/v1/categories/1/posts \
  -H "Content-Type: application/json" \
  -d '{"title": "Go is awesome", "content": "..."}'
```

#### Comparison: add-model vs Manual

| Aspect | add-model | Manual |
|--------|-----------|--------|
| **Time** | < 2 seconds | ~30-60 minutes |
| **Files created** | 8 automatically | 8 manually |
| **Files modified** | 3 automatically | 3 manually |
| **Errors** | Minimal (tested generator) | High risk (typos, omissions) |
| **Tests** | Automatically generated | Must be written manually |
| **Best practices** | Always followed | Depends on the developer |
| **Relationships** | Native BelongsTo/HasMany support | Manual configuration |
| **Customization** | Easy after generation | Full control from the start |

**Recommendation:** Use `add-model` for 90%+ of cases, then customize as needed.

#### Limitations and Workarounds

##### Pluralization

Simple rules: `Todo→todos`, `Category→categories`, `Person→persons` (not `people`)

**Workaround:** Manually edit files for irregular plurals.

##### Many-to-many relationships

Not yet natively supported (planned for v1.3.0).

**Workaround:** Create a manual join table:

```bash
create-go-starter add-model UserRole \
  --fields "user_id:uint:index,role_id:uint:index"

# Then edit internal/models/user_role.go to add a unique constraint:
# UserID uint `gorm:"uniqueIndex:user_role_unique"`
# RoleID uint `gorm:"uniqueIndex:user_role_unique"`
```

#### Full Documentation

For more details on `add-model`, see:
- [Usage Guide](./usage.md#ajouter-des-modeles-add-model)
- [CLI Architecture](./cli-architecture.md#21-générateur-de-modèles-add-model-command)
- [Changelog v1.2.0](../CHANGELOG.md#120---2026-02-12)

### Adding a New Feature (manual method)

This section guides you step by step to add a new entity/feature while respecting hexagonal architecture.

#### Overview of the 9 Steps

```mermaid
flowchart LR
    A["1. Model"] --> B["2. Interface"]
    B --> C["3. Repository"]
    B --> D["4. Service"]
    C --> E["5. fx Module"]
    D --> E
    D --> F["6. Handler"]
    F --> G["7. Routes"]
    A --> H["8. Migration"]
    E --> I["9. Bootstrap"]
    G --> I
```

#### Quick Checklist

Use this checklist to make sure you don't miss anything:

| Step | File to create/modify | Depends on | Status |
|------|----------------------|------------|--------|
| 1. Model | `internal/models/<entity>.go` | - | [ ] |
| 2. Interface | `internal/interfaces/<entity>_repository.go` | Step 1 | [ ] |
| 3. Repository | `internal/adapters/repository/<entity>_repository.go` | Steps 1, 2 | [ ] |
| 4. Service | `internal/domain/<entity>/service.go` | Steps 1, 2 | [ ] |
| 5. fx Module | `internal/domain/<entity>/module.go` | Steps 3, 4 | [ ] |
| 6. Handler | `internal/adapters/handlers/<entity>_handler.go` | Steps 1, 4 | [ ] |
| 7. Routes | `internal/infrastructure/server/server.go` (modify) | Step 6 | [ ] |
| 8. Migration | `internal/infrastructure/database/database.go` (modify) | Step 1 | [ ] |
| 9. Bootstrap | `cmd/main.go` (modify) | Step 5 | [ ] |

#### File Dependency Diagram

This diagram shows the file creation order and their dependencies:

```mermaid
flowchart TD
    subgraph Step1["Step 1 - Foundation"]
        Model["models/product.go<br/>GORM Entity"]
    end

    subgraph Step2["Step 2 - Abstraction"]
        Interface["interfaces/product_repository.go<br/>Port (contract)"]
    end

    subgraph Step34["Steps 3 and 4 - Implementation"]
        Repo["repository/product_repository.go<br/>GORM Adapter"]
        Service["domain/product/service.go<br/>Business Logic"]
    end

    subgraph Step5["Step 5 - DI"]
        Module["domain/product/module.go<br/>fx.Module"]
    end

    subgraph Step6["Step 6 - HTTP"]
        Handler["handlers/product_handler.go<br/>REST endpoints"]
    end

    subgraph Step789["Steps 7, 8, 9 - Integration"]
        Routes["server/server.go<br/>Add routes"]
        Migration["database/database.go<br/>AutoMigrate"]
        Main["cmd/main.go<br/>Add module"]
    end

    Model --> Interface
    Model --> Repo
    Model --> Service
    Interface --> Repo
    Interface --> Service
    Repo --> Module
    Service --> Module
    Service --> Handler
    Model --> Handler
    Handler --> Routes
    Module --> Main
    Routes --> Main
    Model --> Migration
```

---

### Complete Example: Product Entity

We will create a complete `Product` entity with CRUD. Follow each step in order.

> **Tip**: Replace `mon-projet` with your project name in all imports.

---

#### Step 1: Create the Model (Entity)

**File to create: `internal/models/product.go`**

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

**Why?**
- Entities are centralized in `models/` to avoid circular dependencies
- GORM tags for database configuration
- JSON tags to control API serialization
- Separate DTO (`ProductResponse`) to control what is exposed to the API

---

#### Step 2: Define the Interface (Port)

**File to create: `internal/interfaces/product_repository.go`**

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

**Why?**
- **Complete abstraction**: the domain doesn't know about GORM
- **Clear contract**: all available operations are defined
- **Testable**: easy to mock for unit tests

---

#### Step 3: Implement the Repository (Adapter)

**File to create: `internal/adapters/repository/product_repository.go`**

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

**Key points:**
- Implements the `ProductRepository` interface
- Uses `WithContext(ctx)` for context propagation
- Converts `gorm.ErrRecordNotFound` to `DomainError`

---

#### Step 4: Create the Service (Domain/Business Logic)

**Create the directory:**
```bash
mkdir -p internal/domain/product
```

**File to create: `internal/domain/product/service.go`**

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

**Key points:**
- All **business logic** is here (validation, business rules)
- Uses `DomainError` for business errors
- Structured logging with context
- The service only knows about **interfaces**, not implementations

---

#### Step 5: Create the fx Module (Dependency Injection)

**File to create: `internal/domain/product/module.go`**

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

**Why fx.Annotate?**
- Allows providing a concrete implementation while exposing the interface
- Facilitates replacing implementations (tests, mocks, different DB)

---

#### Step 6: Create the Handler (HTTP Adapter)

**File to create: `internal/adapters/handlers/product_handler.go`**

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

**Key points:**
- Uses the `ProductService` interface, not the concrete implementation
- Validation with `go-playground/validator`
- Returns DTOs (`ToResponse()`) instead of entities directly
- Clean error handling with `DomainError`

---

#### Step 7: Add Routes

**Modify: `internal/adapters/http/routes.go`**

Add the `productHandler` parameter and routes:

```go
func RegisterRoutes(
    app *fiber.App,
    authHandler *handlers.AuthHandler,
    userHandler *handlers.UserHandler,
    productHandler *handlers.ProductHandler,  // <- ADD
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
    // ADD: Product routes (protected)
    // ============================================
    products := v1.Group("/products", authMiddleware)
    products.Post("", productHandler.Create)
    products.Get("", productHandler.List)
    products.Get("/:id", productHandler.GetByID)
    products.Put("/:id", productHandler.Update)
    products.Delete("/:id", productHandler.Delete)
}
```

**Advantages of this approach**:
- All routes are visible in a single file
- Easy to add new domains
- API versioning is managed centrally

---

#### Step 8: Add the Migration

**Modify: `internal/infrastructure/database/database.go`**

```go
func NewDatabase(config *config.Config, logger zerolog.Logger) (*gorm.DB, error) {
    // ... existing code ...

    // AutoMigrate - ADD models.Product
    if err := db.AutoMigrate(
        &models.User{},
        &models.RefreshToken{},
        &models.Product{},  // <- ADD
    ); err != nil {
        return nil, fmt.Errorf("failed to auto-migrate: %w", err)
    }

    // ... rest of the code ...
}
```

---

#### Step 9: Register the Module in Bootstrap

**Modify: `cmd/main.go`**

```go
package main

import (
    "go.uber.org/fx"

    "mon-projet/internal/domain/product"  // <- ADD
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
        product.Module,  // <- ADD
        server.Module,
    ).Run()
}
```

---

#### Final Verification

```bash
# 1. Verify compilation
go build ./...

# 2. Start the application
make run

# 3. Authenticate to get a token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' \
  | jq -r '.data.access_token')

# 4. Create a product
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

# 5. List products
curl -X GET "http://localhost:8080/api/v1/products?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN"

# 6. Get a product by ID
curl -X GET http://localhost:8080/api/v1/products/1 \
  -H "Authorization: Bearer $TOKEN"

# 7. Update a product
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

# 8. Delete a product
curl -X DELETE http://localhost:8080/api/v1/products/1 \
  -H "Authorization: Bearer $TOKEN"
```

---

### Summary: Files Created/Modified

| Action | File | Description |
|--------|------|-------------|
| Create | `internal/models/product.go` | GORM Entity + DTO |
| Create | `internal/interfaces/product_repository.go` | Interfaces (Ports) |
| Create | `internal/adapters/repository/product_repository.go` | GORM Implementation |
| Create | `internal/domain/product/service.go` | Business Logic |
| Create | `internal/domain/product/module.go` | fx.Module |
| Create | `internal/adapters/handlers/product_handler.go` | HTTP Handler |
| Modify | `internal/infrastructure/server/server.go` | Add routes |
| Modify | `internal/infrastructure/database/database.go` | Add AutoMigrate |
| Modify | `cmd/main.go` | Add product.Module |

### Patterns to Follow

#### 1. Error Handling

Use DomainErrors:

```go
// In service
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

The error_handler middleware automatically converts these to HTTP responses.

#### 2. Repository Pattern

```go
// Interface (port)
type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    FindByEmail(ctx context.Context, email string) (*models.User, error)
}

// Implementation (adapter)
type userRepositoryGORM struct {
    db *gorm.DB
}
```

#### 3. Dependency Injection with fx

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
    userService interfaces.UserService  // Depends on the interface
}
```

#### 4. Middleware Chain

```go
protected := api.Group("/users")
protected.Use(authMiddleware.Authenticate())  // JWT required
protected.Get("/", userHandler.List)
```

---

## API Reference

### Available Endpoints

#### Health Checks (Kubernetes-compatible)

**Liveness probe** — is the application running?
```bash
GET /health/liveness
# Backward-compatible alias:
GET /health
```

**Response** (200 — always if the app is running):
```json
{
  "status": "alive",
  "service": "mon-projet",
  "timestamp": "2026-02-17T10:00:00Z"
}
```

**Readiness probe** — is the application ready to receive traffic?
```bash
GET /health/readiness
```

**Response** (200 — DB accessible):
```json
{
  "status": "ready",
  "service": "mon-projet",
  "timestamp": "2026-02-17T10:00:00Z",
  "checks": {"database": "ok"}
}
```

**Response** (503 — DB inaccessible):
```json
{
  "status": "not_ready",
  "service": "mon-projet",
  "timestamp": "2026-02-17T10:00:00Z",
  "checks": {"database": "error"},
  "error": "database connection failed"
}
```

> **K8s configuration**: The `deployments/kubernetes/probes.yaml` file is automatically generated with recommended configurations for `livenessProbe`, `readinessProbe`, and `startupProbe`.

---

### Authentication

#### Register

```
POST /api/v1/auth/register
Content-Type: application/json
```

**Body**:
```json
{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

**Validation**:
- email: required, valid email, max 255 chars
- password: required, min 8 chars, max 72 chars

**Response** (201):
```json
{
  "status": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 900
  }
}
```

**Errors**:
- `400 Bad Request`: Invalid input
- `409 Conflict`: Email already exists

**Curl example**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

---

#### Login

```
POST /api/v1/auth/login
Content-Type: application/json
```

**Body**:
```json
{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

**Response** (200):
```json
{
  "status": "success",
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "token_type": "Bearer",
    "expires_in": 900
  }
}
```

**Errors**:
- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Invalid credentials

**Curl example**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

---

#### Refresh Token

```
POST /api/v1/auth/refresh
Content-Type: application/json
```

**Body**:
```json
{
  "refresh_token": "eyJhbGc..."
}
```

**Response** (200):
```json
{
  "status": "success",
  "data": {
    "access_token": "eyJhbGc...",
    "token_type": "Bearer",
    "expires_in": 900
  }
}
```

**Errors**:
- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Invalid or expired refresh token

**Curl example**:
```bash
REFRESH_TOKEN="<refresh_token_from_login>"
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\":\"$REFRESH_TOKEN\"}"
```

---

### Users (Protected)

All user endpoints require a valid JWT token.

#### List Users

```
GET /api/v1/users
Authorization: Bearer <access_token>
```

**Response** (200):
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "email": "user1@example.com",
      "created_at": "2026-01-09T10:00:00Z",
      "updated_at": "2026-01-09T10:00:00Z"
    },
    {
      "id": 2,
      "email": "user2@example.com",
      "created_at": "2026-01-09T11:00:00Z",
      "updated_at": "2026-01-09T11:00:00Z"
    }
  ]
}
```

**Errors**:
- `401 Unauthorized`: Missing or invalid token

**Curl example**:
```bash
TOKEN="<access_token>"
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer $TOKEN"
```

---

#### Get User by ID

```
GET /api/v1/users/:id
Authorization: Bearer <access_token>
```

**Response** (200):
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "created_at": "2026-01-09T10:00:00Z",
    "updated_at": "2026-01-09T10:00:00Z"
  }
}
```

**Errors**:
- `401 Unauthorized`: Invalid token
- `404 Not Found`: User not found

**Curl example**:
```bash
TOKEN="<access_token>"
curl -X GET http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer $TOKEN"
```

---

#### Update User

```
PUT /api/v1/users/:id
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Body**:
```json
{
  "email": "newemail@example.com"
}
```

**Response** (200):
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "email": "newemail@example.com",
    "created_at": "2026-01-09T10:00:00Z",
    "updated_at": "2026-01-10T15:30:00Z"
  }
}
```

**Errors**:
- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Invalid token
- `404 Not Found`: User not found
- `409 Conflict`: Email already exists

---

#### Delete User

```
DELETE /api/v1/users/:id
Authorization: Bearer <access_token>
```

**Response** (200):
```json
{
  "status": "success",
  "message": "User deleted successfully"
}
```

**Errors**:
- `401 Unauthorized`: Invalid token
- `404 Not Found`: User not found

**Note**: Uses soft delete (DeletedAt), data remains in the DB.

---

### Complete API Workflow

```bash
# 1. Register
REGISTER_RESP=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}')

# Extract access_token
ACCESS_TOKEN=$(echo $REGISTER_RESP | jq -r '.data.access_token')
REFRESH_TOKEN=$(echo $REGISTER_RESP | jq -r '.data.refresh_token')

# 2. List users (with token)
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer $ACCESS_TOKEN"

# 3. Update user
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"email":"updated@example.com"}'

# 4. When access token expires (15min), use refresh token
NEW_ACCESS=$(curl -s -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\":\"$REFRESH_TOKEN\"}" | jq -r '.data.access_token')

# 5. Continue with new token
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer $NEW_ACCESS"
```

---

## Tests

### Test Organization

Tests are co-located with source code:

```
internal/
├── adapters/
│   ├── handlers/
│   │   ├── auth_handler.go
│   │   ├── auth_handler_test.go
│   │   ├── user_handler.go
│   │   └── user_handler_test.go
│   ├── middleware/
│   │   ├── auth_middleware.go
│   │   └── auth_middleware_test.go
│   └── repository/
│       ├── user_repository.go
│       └── user_repository_test.go
├── domain/
│   ├── user/
│   │   ├── service.go
│   │   └── service_test.go
│   ├── errors.go
│   └── errors_test.go
```

### Running Tests

```bash
# All tests
make test

# Tests with coverage
make test-coverage

# Open the HTML report
open coverage.html  # macOS
xdg-open coverage.html  # Linux

# Tests for a specific package
go test -v ./internal/domain/user

# Specific test
go test -run TestRegister ./internal/adapters/handlers

# Tests with race detector (race condition detection)
go test -race ./...
```

### Types of Tests

#### 1. Unit Tests

Test an isolated function or method, with mocks.

**Example: Service test**

```go
// internal/domain/user/service_test.go
package user

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func TestService_Register(t *testing.T) {
    // Arrange
    mockRepo := new(MockUserRepository)
    logger := zerolog.Nop()
    service := NewService(mockRepo, logger)

    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)

    // Act
    user, err := service.Register(context.Background(), "test@example.com", "password123")

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "test@example.com", user.Email)
    mockRepo.AssertExpectations(t)
}
```

#### 2. Integration Tests

Test multiple components together, with a real DB (SQLite in-memory).

**Example: Handler test with DB**

```go
// internal/adapters/handlers/auth_handler_integration_test.go
func TestAuthHandler_RegisterIntegration(t *testing.T) {
    // Setup in-memory DB
    db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
    require.NoError(t, err)

    db.AutoMigrate(&models.User{})

    // Create real dependencies
    repo := repository.NewUserRepository(db)
    service := user.NewService(repo, zerolog.Nop())
    handler := handlers.NewAuthHandler(service, "test-secret")

    // Create Fiber app
    app := fiber.New()
    app.Post("/register", handler.Register)

    // Test request
    body := `{"email":"test@example.com","password":"password123"}`
    req := httptest.NewRequest("POST", "/register", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req)
    require.NoError(t, err)

    // Assert
    assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)
    assert.Equal(t, "success", result["status"])
}
```

#### 3. Table-Driven Tests

For testing multiple cases with a common structure:

```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid email", "user@example.com", false},
        {"invalid email - no @", "userexample.com", true},
        {"invalid email - no domain", "user@", true},
        {"empty email", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateEmail(tt.email)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Best Practices for Tests

1. **Use testify/assert** for clear assertions:
   ```go
   assert.Equal(t, expected, actual)
   assert.NoError(t, err)
   assert.NotNil(t, user)
   ```

2. **Arrange-Act-Assert pattern**:
   ```go
   // Arrange - Setup
   mockRepo := new(MockUserRepository)

   // Act - Execute
   result, err := service.DoSomething()

   // Assert - Verify
   assert.NoError(t, err)
   assert.Equal(t, expected, result)
   ```

3. **Mock external dependencies**:
   - DB (except for integration tests)
   - External APIs
   - Third-party services

4. **Integration tests with SQLite in-memory**:
   ```go
   db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
   ```

5. **Clean up after each test**:
   ```go
   t.Cleanup(func() {
       db.Exec("DELETE FROM users")
   })
   ```

6. **Descriptive test names**:
   ```go
   func TestUserService_Register_WhenEmailAlreadyExists_ReturnsConflictError(t *testing.T)
   ```

### Coverage

Target: **> 80% coverage**

```bash
# Generate report
make test-coverage

# View coverage per package
go test -cover ./...

# Output:
# ok      mon-projet/internal/domain/user         0.123s  coverage: 85.7% of statements
# ok      mon-projet/internal/adapters/handlers   0.234s  coverage: 92.3% of statements
```

---

## Database

### Migrations

The project uses GORM AutoMigrate to simplify migrations in development:

```go
// internal/infrastructure/database/database.go
db.AutoMigrate(
    &models.User{},
    &models.RefreshToken{},
)
```

**AutoMigrate**:
- Creates tables if they don't exist
- Adds missing columns
- Creates indexes
- **Does NOT delete** columns or tables

**For production**, consider a versioned migration solution:

- **golang-migrate/migrate**: SQL or Go migrations
- **pressly/goose**: Up/down migrations
- **Advanced GORM Migrator**: Programmatic API

**Example with golang-migrate**:

```bash
# Install migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create migration
migrate create -ext sql -dir migrations -seq create_users_table

# Files created:
# migrations/000001_create_users_table.up.sql
# migrations/000001_create_users_table.down.sql

# Run migrations
migrate -path migrations -database "postgresql://user:pass@localhost/dbname?sslmode=disable" up
```

### GORM Models

Conventions and patterns:

```go
type User struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

    Email    string `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
    Password string `gorm:"not null" json:"-"`
}
```

**Important GORM tags**:

| Tag | Description |
|-----|-------------|
| `primarykey` | Primary key |
| `uniqueIndex` | Unique index |
| `index` | Simple index |
| `not null` | NOT NULL column |
| `default:value` | Default value |
| `size:255` | Column size |
| `type:varchar(100)` | Custom SQL type |
| `foreignKey:UserID` | Foreign key |
| `references:ID` | FK reference |

**Conventions**:

- **Soft deletes**: `DeletedAt gorm.DeletedAt`
- **Auto timestamps**: `CreatedAt`, `UpdatedAt`
- **JSON hiding**: `json:"-"` for password
- **Index on FK**: Always index foreign keys

### Advanced Queries

#### Pagination

```go
var users []models.User
limit := 10
offset := 20

db.Limit(limit).Offset(offset).Find(&users)
```

#### Filtering

```go
// Simple Where
db.Where("email = ?", "user@example.com").First(&user)

// Where with multiple conditions
db.Where("created_at > ? AND email LIKE ?", time.Now().Add(-24*time.Hour), "%@example.com").Find(&users)

// Or
db.Where("email = ?", email1).Or("email = ?", email2).Find(&users)
```

#### Sorting

```go
// Order ASC
db.Order("created_at asc").Find(&users)

// Order DESC
db.Order("created_at desc").Find(&users)

// Multiple sorts
db.Order("created_at desc, email asc").Find(&users)
```

#### Joins

```go
// Inner join
db.Joins("LEFT JOIN refresh_tokens ON refresh_tokens.user_id = users.id").
   Where("refresh_tokens.expires_at > ?", time.Now()).
   Find(&users)

// Preload associations
db.Preload("RefreshTokens").Find(&users)
```

#### Aggregations

```go
// Count
var count int64
db.Model(&models.User{}).Count(&count)

// With where
db.Model(&models.User{}).Where("created_at > ?", yesterday).Count(&count)
```

#### Transactions

```go
err := db.Transaction(func(tx *gorm.DB) error {
    // Create user
    if err := tx.Create(&user).Error; err != nil {
        return err  // Rollback
    }

    // Create profile
    if err := tx.Create(&profile).Error; err != nil {
        return err  // Rollback
    }

    return nil  // Commit
})
```

#### Raw SQL

```go
// Raw query
var users []models.User
db.Raw("SELECT * FROM users WHERE email LIKE ?", "%@example.com").Scan(&users)

// Exec
db.Exec("UPDATE users SET email = ? WHERE id = ?", newEmail, userID)
```

### Performance Tips

1. **Index frequently queried columns**:
   ```go
   Email string `gorm:"uniqueIndex"`
   ```

2. **Avoid N+1 queries with Preload**:
   ```go
   // ❌ N+1
   for _, user := range users {
       db.Model(&user).Association("RefreshTokens").Find(&tokens)
   }

   // :material-check-circle: Single query
   db.Preload("RefreshTokens").Find(&users)
   ```

3. **Select only necessary columns**:
   ```go
   db.Select("id, email").Find(&users)
   ```

4. **Use connection pools**:
   ```go
   sqlDB, _ := db.DB()
   sqlDB.SetMaxIdleConns(10)
   sqlDB.SetMaxOpenConns(100)
   sqlDB.SetConnMaxLifetime(time.Hour)
   ```

---

## Security

### JWT Authentication

#### Complete Flow

```
┌─────────┐                                  ┌─────────┐
│ Client  │                                  │ Server  │
└────┬────┘                                  └────┬────┘
     │                                            │
     │ 1. POST /auth/register or /login          │
     │───────────────────────────────────────────>│
     │                                            │
     │ 2. Access Token (15min) + Refresh (7d)    │
     │<───────────────────────────────────────────│
     │                                            │
     │ 3. GET /users (Authorization: Bearer AT)  │
     │───────────────────────────────────────────>│
     │                                            │
     │ 4. Response                                │
     │<───────────────────────────────────────────│
     │                                            │
     │ [15 minutes later - Access Token expires] │
     │                                            │
     │ 5. GET /users (Authorization: Bearer AT)  │
     │───────────────────────────────────────────>│
     │                                            │
     │ 6. 401 Unauthorized (token expired)        │
     │<───────────────────────────────────────────│
     │                                            │
     │ 7. POST /auth/refresh (Refresh Token)     │
     │───────────────────────────────────────────>│
     │                                            │
     │ 8. New Access Token (15min)                │
     │<───────────────────────────────────────────│
     │                                            │
     │ 9. Continue with new Access Token          │
     │───────────────────────────────────────────>│
     │                                            │
```

#### Token Storage (client-side)

**Access Token** (short-lived: 15min):
- **Recommended**: In memory (JavaScript variable)
- **Not in localStorage** (vulnerable to XSS)
- Lost on page refresh → Use refresh token

**Refresh Token** (long-lived: 7d):
- **Option 1**: httpOnly cookie (most secure)
- **Option 2**: localStorage (if no XSS risk)

**React example**:

```javascript
// Store access token in memory
let accessToken = null;

// Login
const login = async (email, password) => {
    const response = await fetch('/api/v1/auth/login', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({email, password})
    });
    const data = await response.json();

    accessToken = data.data.access_token;  // Memory
    localStorage.setItem('refresh_token', data.data.refresh_token);
};

// API call
const fetchUsers = async () => {
    const response = await fetch('/api/v1/users', {
        headers: {'Authorization': `Bearer ${accessToken}`}
    });

    if (response.status === 401) {
        // Token expired, refresh
        await refreshAccessToken();
        // Retry request
    }
};

// Refresh
const refreshAccessToken = async () => {
    const refreshToken = localStorage.getItem('refresh_token');
    const response = await fetch('/api/v1/auth/refresh', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({refresh_token: refreshToken})
    });
    const data = await response.json();
    accessToken = data.data.access_token;
};
```

### Route Protection

Authentication middleware:

```go
// internal/adapters/middleware/auth_middleware.go
func (m *AuthMiddleware) Authenticate() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // 1. Extract token
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return fiber.NewError(fiber.StatusUnauthorized, "Missing authorization header")
        }

        // 2. Parse "Bearer <token>"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            return fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization format")
        }

        // 3. Validate JWT
        claims, err := auth.ParseToken(parts[1], m.jwtSecret)
        if err != nil {
            return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
        }

        // 4. Inject user ID in context
        c.Locals("user_id", claims.UserID)

        return c.Next()
    }
}
```

**Usage**:

```go
// Protected routes
users := api.Group("/users")
users.Use(authMiddleware.Authenticate())  // Middleware applied
users.Get("/", userHandler.List)
```

**In the handler, retrieve user ID**:

```go
func (h *UserHandler) List(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    // Use userID to check permissions, etc.
}
```

### Input Validation

go-playground/validator v10:

```go
type RegisterRequest struct {
    Email    string `json:"email" validate:"required,email,max=255"`
    Password string `json:"password" validate:"required,min=8,max=72"`
}

// In handler
validate := validator.New()
if err := validate.Struct(req); err != nil {
    return domain.NewValidationError("Invalid input", "VALIDATION_ERROR", err)
}
```

**Common validation tags**:

| Tag | Description |
|-----|-------------|
| `required` | Required field |
| `email` | Valid email format |
| `min=N` | Minimum length/value |
| `max=N` | Maximum length/value |
| `len=N` | Exact length |
| `gte=N` | Greater than or equal |
| `lte=N` | Less than or equal |
| `alpha` | Letters only |
| `alphanum` | Letters + digits |
| `numeric` | Digits only |
| `uuid` | UUID format |
| `url` | URL format |

**Custom validators**:

```go
validate := validator.New()

// Register custom validator
validate.RegisterValidation("strong_password", func(fl validator.FieldLevel) bool {
    password := fl.Field().String()
    // Custom logic: must contain uppercase, lowercase, number, special char
    return hasUppercase(password) && hasLowercase(password) && hasNumber(password)
})

// Usage
type Request struct {
    Password string `validate:"required,strong_password"`
}
```

### Password Hashing

bcrypt (golang.org/x/crypto/bcrypt):

```go
// internal/domain/user/service.go

import "golang.org/x/crypto/bcrypt"

func (s *Service) Register(ctx context.Context, email, password string) (*models.User, error) {
    // Hash the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    // Create user with hashed password
    user := &models.User{
        Email:        email,
        PasswordHash: string(hashedPassword),
    }

    if err := s.repo.CreateUser(ctx, user); err != nil {
        return nil, err
    }

    return user, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*models.User, error) {
    user, err := s.repo.GetUserByEmail(ctx, email)
    if err != nil {
        return nil, err
    }

    // Compare password
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
        return nil, domain.NewUnauthorizedError("Invalid credentials", "INVALID_CREDENTIALS", err)
    }

    return user, nil
}
```

**DefaultCost** = 10 (2^10 iterations) - Good balance between security and performance

**Note**: Password hashing is handled in the service (business logic), not in the entity. The `models.User` entity stores the `PasswordHash` which is always hashed.

**Usage**:

```go
// Register (in the service)
user, err := userService.Register(ctx, email, plainPassword)
if err != nil {
    return err
}
db.Create(user)

// Login
user, _ := repo.FindByEmail(email)
if err := user.ComparePassword(plainPassword); err != nil {
    return domain.NewUnauthorizedError("Invalid credentials", "INVALID_CREDENTIALS", err)
}
```

### Security Checklist

Production checklist:

- [ ] **Strong JWT_SECRET**: Generated with `openssl rand -base64 32`
- [ ] **HTTPS in production**: Always use TLS
- [ ] **Rate limiting**: Implement with fiber/limiter
  ```go
  import "github.com/gofiber/fiber/v2/middleware/limiter"

  app.Use(limiter.New(limiter.Config{
      Max:        100,
      Expiration: 1 * time.Minute,
  }))
  ```
- [ ] **CORS configured**: Restrict origins
  ```go
  import "github.com/gofiber/fiber/v2/middleware/cors"

  app.Use(cors.New(cors.Config{
      AllowOrigins: "https://your-frontend.com",
      AllowHeaders: "Origin, Content-Type, Accept, Authorization",
  }))
  ```
- [ ] **Strict validation**: All inputs validated
- [ ] **SQL Injection**: GORM prevents it automatically
- [ ] **XSS**: Escape HTML outputs (if using templates)
- [ ] **Logs without secrets**: Never log passwords, tokens
  ```go
  // ❌ BAD
  logger.Info().Str("password", password).Msg("User login")

  // :material-check-circle: GOOD
  logger.Info().Str("email", email).Msg("User login")
  ```
- [ ] **Environment variables**: Secrets in .env (not in code)
- [ ] **DB SSL**: `DB_SSLMODE=require` in production
- [ ] **Helmet headers**: Implement security headers
  ```go
  import "github.com/gofiber/fiber/v2/middleware/helmet"

  app.Use(helmet.New())
  ```
- [ ] **Request timeouts**: Prevent DoS
  ```go
  app := fiber.New(fiber.Config{
      ReadTimeout:  10 * time.Second,
      WriteTimeout: 10 * time.Second,
      IdleTimeout:  30 * time.Second,
  })
  ```

---

## Deployment

### Docker

#### Building the Image

```bash
make docker-build
```

Or manually:

```bash
docker build -t mon-projet:latest .
```

The generated Dockerfile uses a multi-stage build:

```dockerfile
# Stage 1: Build
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# Stage 2: Runtime
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY .env.example .env

EXPOSE 8080
CMD ["./main"]
```

**Advantages**:
- Lightweight final image (~15-20MB vs ~1GB)
- Security (minimal alpine image)
- Static binary (no dependencies)

#### Running with Docker

```bash
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PASSWORD=postgres \
  -e JWT_SECRET=<your_secret> \
  mon-projet:latest
```

**Note**: `host.docker.internal` allows accessing localhost from Docker.

#### Docker Compose

If `docker-compose.yml` is generated:

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      DB_HOST: postgres
      DB_USER: postgres
      DB_PASSWORD: postgres
      DB_NAME: mon-projet
      JWT_SECRET: ${JWT_SECRET}
    depends_on:
      - postgres

  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: mon-projet
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

**Start**:

```bash
# Set JWT_SECRET
export JWT_SECRET=$(openssl rand -base64 32)

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop
docker-compose down
```

### Kubernetes

Basic manifests for K8s deployment:

#### Secret

```yaml
# k8s/secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: mon-projet-secret
type: Opaque
stringData:
  jwt-secret: "<your_secret_base64>"
  db-password: "postgres"
```

```bash
kubectl apply -f k8s/secret.yaml
```

#### Deployment

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mon-projet
  labels:
    app: mon-projet
spec:
  replicas: 3
  selector:
    matchLabels:
      app: mon-projet
  template:
    metadata:
      labels:
        app: mon-projet
    spec:
      containers:
      - name: mon-projet
        image: mon-projet:latest
        ports:
        - containerPort: 8080
        env:
        - name: APP_PORT
          value: "8080"
        - name: DB_HOST
          value: postgres-service
        - name: DB_USER
          value: postgres
        - name: DB_NAME
          value: mon-projet
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: mon-projet-secret
              key: db-password
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: mon-projet-secret
              key: jwt-secret
        livenessProbe:
          httpGet:
            path: /health/liveness
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /health/readiness
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "500m"
```

```bash
kubectl apply -f k8s/deployment.yaml
```

#### Service

```yaml
# k8s/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: mon-projet-service
spec:
  selector:
    app: mon-projet
  ports:
  - port: 80
    targetPort: 8080
    protocol: TCP
  type: LoadBalancer
```

```bash
kubectl apply -f k8s/service.yaml
```

#### Deploy PostgreSQL (StatefulSet)

```yaml
# k8s/postgres.yaml
apiVersion: v1
kind: Service
metadata:
  name: postgres-service
spec:
  selector:
    app: postgres
  ports:
  - port: 5432
  clusterIP: None
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgres
spec:
  serviceName: postgres-service
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:16-alpine
        ports:
        - containerPort: 5432
        env:
        - name: POSTGRES_DB
          value: mon-projet
        - name: POSTGRES_USER
          value: postgres
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: mon-projet-secret
              key: db-password
        volumeMounts:
        - name: postgres-storage
          mountPath: /var/lib/postgresql/data
  volumeClaimTemplates:
  - metadata:
      name: postgres-storage
    spec:
      accessModes: ["ReadWriteOnce"]
      resources:
        requests:
          storage: 10Gi
```

```bash
kubectl apply -f k8s/postgres.yaml
```

### CI/CD with GitHub Actions

The generated workflow (`.github/workflows/ci.yml`):

```yaml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.25'

    - name: golangci-lint
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest

  test:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:16-alpine
        env:
          POSTGRES_DB: test_db
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: postgres
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432

    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.25'

    - name: Run tests
      env:
        DB_HOST: localhost
        DB_PORT: 5432
        DB_USER: postgres
        DB_PASSWORD: postgres
        DB_NAME: test_db
        JWT_SECRET: test-secret
      run: go test -v -race -coverprofile=coverage.out ./...

    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        files: ./coverage.out

  build:
    runs-on: ubuntu-latest
    needs: [quality, test]
    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.25'

    - name: Build
      run: go build -v -o mon-projet cmd/main.go
```

**Pipeline**:
1. **Quality**: golangci-lint
2. **Test**: Tests with PostgreSQL (service container)
3. **Build**: Build verification

### Production Deployment

#### Pre-deployment Checklist

- [ ] All tests pass (`make test`)
- [ ] Lint passes (`make lint`)
- [ ] Environment variables configured
- [ ] JWT_SECRET generated (strong, random)
- [ ] DB_SSLMODE=require
- [ ] DB migrations executed
- [ ] Health check works
- [ ] Logs configured
- [ ] Monitoring in place

#### Recommended Platforms

**1. Google Cloud Run** (simplest):

```bash
# Build and push image
gcloud builds submit --tag gcr.io/PROJECT_ID/mon-projet

# Deploy
gcloud run deploy mon-projet \
  --image gcr.io/PROJECT_ID/mon-projet \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars JWT_SECRET=$JWT_SECRET,DB_HOST=$DB_HOST
```

**2. AWS ECS/Fargate**:

- Build image → Push to ECR
- Create Task Definition
- Create ECS Service
- Configure ALB

**3. Heroku**:

```bash
# Login
heroku login

# Create app
heroku create mon-projet

# Add PostgreSQL
heroku addons:create heroku-postgresql:hobby-dev

# Set env vars
heroku config:set JWT_SECRET=$(openssl rand -base64 32)

# Deploy
git push heroku main
```

**4. Kubernetes** (most flexible):

```bash
# Apply all manifests
kubectl apply -f k8s/

# Check status
kubectl get pods
kubectl get services

# View logs
kubectl logs -f deployment/mon-projet
```

---

## Monitoring & Logging

### Logging with zerolog

#### Configuration

The logger is configured in `pkg/logger/logger.go`:

```go
func NewLogger(config *config.Config) zerolog.Logger {
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

    var logger zerolog.Logger

    if config.AppEnv == "production" {
        logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
    } else {
        logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
            With().
            Timestamp().
            Logger()
    }

    // Set level based on env
    switch config.AppEnv {
    case "production":
        zerolog.SetGlobalLevel(zerolog.InfoLevel)
    case "development":
        zerolog.SetGlobalLevel(zerolog.DebugLevel)
    default:
        zerolog.SetGlobalLevel(zerolog.InfoLevel)
    }

    return logger
}
```

#### Usage

Injection via fx:

```go
type UserService struct {
    logger zerolog.Logger
}

func NewUserService(logger zerolog.Logger) *UserService {
    return &UserService{logger: logger}
}
```

**Structured logging**:

```go
// Info
logger.Info().
    Str("email", user.Email).
    Uint("user_id", user.ID).
    Msg("User registered successfully")

// Error
logger.Error().
    Err(err).
    Str("operation", "create_user").
    Str("email", email).
    Msg("Failed to create user")

// Debug
logger.Debug().
    Interface("request", req).
    Msg("Received request")

// Warn
logger.Warn().
    Dur("duration", elapsed).
    Msg("Slow query detected")

// Fatal (exits)
logger.Fatal().
    Err(err).
    Msg("Cannot connect to database")
```

#### Log Levels

| Level | Usage |
|-------|-------|
| **Debug** | Detailed information for debugging |
| **Info** | Important events (user login, etc.) |
| **Warn** | Non-critical abnormal behaviors |
| **Error** | Errors requiring attention |
| **Fatal** | Critical errors (app exit) |

#### Best Practices

**<i class="material-icons success">check_circle</i> GOOD - Structured logging**:

```go
logger.Info().
    Str("user_id", userID).
    Str("action", "login").
    Dur("duration", elapsed).
    Msg("User logged in")
```

**❌ BAD - String formatting**:

```go
logger.Info().Msgf("User %s logged in after %v", userID, elapsed)
```

**<i class="material-icons success">check_circle</i> GOOD - No secrets**:

```go
logger.Info().Str("email", email).Msg("User login attempt")
```

**❌ BAD - Logging secrets**:

```go
logger.Info().Str("password", password).Msg("Login")  // NEVER!
```

### Full Observability (`--observability=advanced`)

When a project is generated with `--observability=advanced`, a complete observability stack is included:

#### Included Stack

| Service | URL | Description |
|---------|-----|-------------|
| **Grafana** | `http://localhost:3000` | Dashboards & alerts (admin/admin) |
| **Prometheus** | `http://localhost:9090` | Metrics & alert rules |
| **Jaeger** | `http://localhost:16686` | Distributed traces |

#### Starting the Stack

```bash
docker-compose up -d
```

Grafana is **auto-provisioned**: the API dashboard is available immediately without manual configuration.

#### Default Grafana Credentials

```
URL:       http://localhost:3000
Login:     admin
Password:  admin  (or value of GF_SECURITY_ADMIN_PASSWORD in .env)
```

> **Security**: Change `GF_SECURITY_ADMIN_PASSWORD` in `.env` before deploying to production.

#### Exposed Metrics (`/metrics`)

The application exposes Prometheus metrics on `/metrics`:

| Metric | Type | Description |
|--------|------|-------------|
| `http_requests_total` | Counter | Number of HTTP requests by method/status |
| `http_request_duration_seconds` | Histogram | Request latency |
| `health_check_status` | Gauge | Health check status (1=OK, 0=KO) |
| `health_check_duration_seconds` | Histogram | Health check duration |

#### Grafana Dashboard — Panels

The `api-dashboard.json` dashboard contains 7 pre-configured panels:

- **Request Rate** — Requests/second (timeseries)
- **Error Rate %** — 5xx error rate with thresholds (stat: green/yellow/red)
- **P95 Latency** — 95th percentile latency in ms (gauge)
- **DB Query P95** — 95th percentile DB query duration (timeseries)
- **Active DB Connections** — Active DB connections (stat)
- **Health Status** — Database status (UP/DOWN)
- **HTTP Requests by Status** — Request distribution by status code (timeseries)

#### Pre-configured Prometheus Alerts

Alert rules in `deployments/prometheus/rules/api_alerts.yml`:

| Alert | Condition | Delay | Severity |
|-------|-----------|-------|----------|
| `HighErrorRate` | Error rate > 5% | 2 min | warning |
| `HighP95Latency` | P95 latency > 1s | 5 min | warning |
| `DatabaseDown` | health_check_status == 0 | 1 min | critical |

#### Generated File Structure

```
<project>/
├── deployments/
│   ├── prometheus/
│   │   ├── prometheus.yml              # Scraping config + rules
│   │   └── rules/
│   │       └── api_alerts.yml          # Alerts (ErrorRate, Latency, DB)
│   └── grafana/
│       ├── provisioning/
│       │   ├── datasources/
│       │   │   └── prometheus.yaml     # Auto-datasource Prometheus
│       │   └── dashboards/
│       │       └── default.yaml        # Auto-provisioning dashboards
│       └── dashboards/
│           └── api-dashboard.json      # Main dashboard (7 panels)
└── docker-compose.yml                  # Includes prometheus, grafana, jaeger
```

### Monitoring (additional recommendations)

For production, complement with:

#### 1. Sentry (Error Tracking)

```go
import "github.com/getsentry/sentry-go"

sentry.Init(sentry.ClientOptions{
    Dsn: os.Getenv("SENTRY_DSN"),
})

// Capture errors
sentry.CaptureException(err)
```

#### 2. APM (Application Performance Monitoring)

- **New Relic**: Complete APM
- **Datadog**: Monitoring + logs
- **Elastic APM**: Open source

### Health Checks

The `/health` endpoint is crucial for:

- Load balancers
- Kubernetes probes
- Monitoring tools

**Enhanced**:

```go
type HealthResponse struct {
    Status   string            `json:"status"`
    Version  string            `json:"version"`
    Services map[string]string `json:"services"`
}

func (h *HealthHandler) Check(c *fiber.Ctx) error {
    // Check database
    dbStatus := "ok"
    if err := h.db.Exec("SELECT 1").Error; err != nil {
        dbStatus = "error"
    }

    response := HealthResponse{
        Status:  "ok",
        Version: "1.0.0",
        Services: map[string]string{
            "database": dbStatus,
        },
    }

    if dbStatus != "ok" {
        return c.Status(fiber.StatusServiceUnavailable).JSON(response)
    }

    return c.JSON(response)
}
```

### Advanced Observability (v1.3.0)

If your project was generated with `--observability=advanced`, a complete observability stack is integrated.

#### Monitoring Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /health/liveness` | Liveness probe — is the app running? |
| `GET /health/readiness` | Readiness probe — is the DB accessible? |
| `GET /health` | Backward-compatible alias to liveness |
| `GET /metrics` | Prometheus metrics (advanced mode only) |

#### Docker Compose Stack

```bash
# Start all monitoring services
docker-compose up -d

# Available services:
# - Jaeger UI:     http://localhost:16686  (distributed traces)
# - Prometheus UI: http://localhost:9090   (metrics)
# - Grafana UI:    http://localhost:3000   (dashboards, credentials: admin/admin)
```

#### Prometheus Metrics

HTTP metrics are automatically collected by the middleware:

```bash
curl http://localhost:8080/metrics
# http_requests_total{method="GET",path="/health",status="200"} 42
# http_request_duration_seconds_bucket{method="GET",path="/health",le="0.1"} 42
```

#### Distributed Traces

OpenTelemetry traces are exported to Jaeger via OTLP/gRPC. Each HTTP request and DB query automatically generates spans:

```bash
# Configure the Jaeger endpoint in .env
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317

# Access Jaeger UI
open http://localhost:16686
```

Zerolog logs are enriched with `trace_id` and `span_id` to correlate logs and traces.

#### Grafana Dashboard

A pre-configured dashboard with 7 panels (request rate, error rate, latency percentiles, etc.) is automatically provisioned when Grafana starts.

For more details, see the [Monitoring & Observability Guide](./guide/monitoring.md).

---

## Best Practices

### Architecture

**1. Domain isolation**

The domain must **never** import other packages:

```go
// ❌ BAD - Domain importing adapter
package user

import "mon-projet/internal/adapters/repository"  // NO!

// :material-check-circle: GOOD - Domain only imports interfaces
package user

import "mon-projet/internal/interfaces"
```

**2. Single Responsibility Principle**

Each component has a single responsibility:

- **Handlers**: Parse + validate + call service
- **Services**: Business logic only
- **Repositories**: Data access only

**3. Dependency Injection**

Always via fx.Provide, no global variables:

```go
// ❌ BAD - Global variable
var db *gorm.DB

// :material-check-circle: GOOD - Injection
type UserService struct {
    db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
    return &UserService{db: db}
}
```

### Code Style

**1. gofmt**

Always format:

```bash
go fmt ./...
```

Or configure your IDE to format on save.

**2. golangci-lint**

Follow the rules:

```bash
make lint
```

**3. GoDoc Documentation**

For public exports:

```go
// UserService handles user-related business logic.
// It provides methods for user registration, authentication, and CRUD operations.
type UserService struct {
    repo   interfaces.UserRepository
    logger zerolog.Logger
}

// Register creates a new user with the provided email and password.
// The password is automatically hashed before storage.
// Returns an error if the email already exists or if validation fails.
func (s *UserService) Register(ctx context.Context, email, password string) (*User, error) {
    // ...
}
```

**4. Explicit error handling**

Always handle errors, don't use `panic`:

```go
// ❌ BAD
user := getUserByID(id)  // What if error?

// :material-check-circle: GOOD
user, err := getUserByID(id)
if err != nil {
    return nil, fmt.Errorf("failed to get user: %w", err)
}
```

### Naming Conventions

**Interfaces**:
- Suffix `-er` or `-Service`
- Examples: `UserRepository`, `AuthService`, `Logger`

**Repositories**:
- Suffix `-Repository`
- Examples: `UserRepository`, `ProductRepository`

**Handlers**:
- Suffix `-Handler`
- Examples: `AuthHandler`, `UserHandler`

**Constructors**:
- Prefix `New`
- Examples: `NewUserService`, `NewAuthHandler`

**Private methods**:
- lowerCamelCase
- Examples: `hashPassword`, `validateEmail`

### Error Handling Patterns

**Wrap errors with context**:

```go
// :material-check-circle: GOOD
if err != nil {
    return fmt.Errorf("failed to create user %s: %w", email, err)
}
```

**Domain errors for business logic**:

```go
if user == nil {
    return domain.NewNotFoundError("User not found", "USER_NOT_FOUND", nil)
}
```

**Don't handle HTTP status in the service**:

```go
// ❌ BAD - Service returning HTTP status
func (s *UserService) GetByID(id uint) (int, *User, error) {
    return 404, nil, errors.New("not found")
}

// :material-check-circle: GOOD - Service returning domain error
func (s *UserService) GetByID(id uint) (*User, error) {
    return nil, domain.NewNotFoundError("User not found", "USER_NOT_FOUND", nil)
}
```

### Testing Best Practices

**1. Coverage > 80%**

```bash
go test -cover ./...
```

**2. Table-driven tests**

```go
tests := []struct {
    name    string
    input   string
    want    string
    wantErr bool
}{
    {"valid", "test", "TEST", false},
    {"empty", "", "", true},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        got, err := ToUpper(tt.input)
        if tt.wantErr {
            assert.Error(t, err)
        } else {
            assert.Equal(t, tt.want, got)
        }
    })
}
```

**3. Descriptive names**

```go
func TestUserService_Register_WhenEmailAlreadyExists_ReturnsConflictError(t *testing.T)
```

**4. Setup/teardown with t.Cleanup()**

```go
func TestSomething(t *testing.T) {
    db := setupTestDB(t)
    t.Cleanup(func() {
        db.Exec("DELETE FROM users")
        db.Close()
    })

    // Test code
}
```

### Performance

**1. GORM - Avoid N+1 queries**

```go
// ❌ N+1 problem
for _, user := range users {
    db.Model(&user).Association("Posts").Find(&posts)
}

// :material-check-circle: Single query with Preload
db.Preload("Posts").Find(&users)
```

**2. Context - Always pass context.Context**

```go
func (s *UserService) GetByID(ctx context.Context, id uint) (*User, error) {
    return s.repo.FindByID(ctx, id)
}
```

**3. Database indexes**

```go
Email string `gorm:"uniqueIndex"`  // Index on frequently queried columns
```

**4. Connection pooling**

```go
sqlDB, _ := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### Security Recap

- [ ] Validate all user inputs
- [ ] Never log passwords or tokens
- [ ] Rate limiting on public endpoints
- [ ] HTTPS in production
- [ ] Strong JWT secret (32+ characters)
- [ ] Bcrypt for passwords
- [ ] Update dependencies regularly

```bash
# Check vulnerabilities
go list -json -m all | nancy sleuth
```

---

## Conclusion

This guide covers all aspects of development with projects generated by `create-go-starter`. To go further:

- **Code examples**: All patterns are in the generated code
- **Tests**: Look at `*_test.go` files for examples
- **Official documentation**:
  - [Fiber](https://docs.gofiber.io/)
  - [GORM](https://gorm.io/docs/)
  - [fx](https://uber-go.github.io/fx/)
  - [zerolog](https://github.com/rs/zerolog)

**Happy coding!** <i class="material-icons info">rocket_launch</i>

If you encounter issues or have questions, check:
- [GitHub Issues](https://github.com/tky0065/go-starter-kit/issues)
- [GitHub Discussions](https://github.com/tky0065/go-starter-kit/discussions)
