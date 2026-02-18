# Technical Glossary

This glossary defines the technical terms used in go-starter-kit and generated projects.

---

## Architecture

### Adapter

**Definition**: Component that implements an interface (port) to connect the business domain to an external technology (database, API, HTTP, etc.).

**Context**: In go-starter-kit, adapters are located in `internal/adapters/` and include HTTP handlers and database repositories.

**Example**:
```go
// HTTP adapter (handler)
type UserHandler struct {
    userService interfaces.UserService
}

// Database adapter (repository)
type UserRepository struct {
    db *gorm.DB
}
```

**See also**: Port, Hexagonal Architecture, Repository

### Domain

**Definition**: Central layer containing pure business logic, independent of infrastructure technical details.

**Context**: The domain is located in `internal/domain/` and contains business services (UserService, AuthService). Business entities are in `internal/models/`.

**Example**:
```go
// Domain service
type userService struct {
    userRepo interfaces.UserRepository
}

func (s *userService) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error) {
    // Pure business logic
}
```

**See also**: Hexagonal Architecture, Service, Models

### Hexagonal Architecture

**Definition**: Architectural pattern (also called "Ports and Adapters") that isolates the business domain from technical details by using interfaces (ports) and implementations (adapters).

**Context**: go-starter-kit uses hexagonal architecture to ensure separation of concerns and facilitate testing.

**Structure**:
```
internal/
├── models/       # Domain entities
├── domain/       # Business logic
├── interfaces/   # Ports (interfaces)
├── adapters/     # Adapters (HTTP, DB)
└── infrastructure/ # Configuration
```

**See also**: Port, Adapter, Domain, Dependency Injection

### Models

**Definition**: Business domain entities shared across all layers of the application.

**Context**: Models are defined in `internal/models/` to avoid circular dependencies. They represent the fundamental business objects (User, RefreshToken, etc.).

**Example**:
```go
// internal/models/user.go
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Email     string    `gorm:"uniqueIndex;not null"`
    Password  string    `gorm:"not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

**See also**: Domain, Entity

### Port

**Definition**: Interface that defines a contract for interacting with the business domain or an external resource.

**Context**: Ports are defined in `internal/interfaces/` and include UserService, UserRepository, AuthService, etc.

**Example**:
```go
// Service port
type UserService interface {
    CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error)
    GetUserByID(ctx context.Context, id uint) (*models.User, error)
}

// Repository port
type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    FindByEmail(ctx context.Context, email string) (*models.User, error)
}
```

**See also**: Adapter, Interface, Hexagonal Architecture

### Repository

**Definition**: Pattern that encapsulates data access logic and provides a collection-like interface for manipulating entities.

**Context**: Repositories are implemented in `internal/adapters/repositories/` and use GORM for PostgreSQL access.

**Example**:
```go
type UserRepository struct {
    db *gorm.DB
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}
```

**See also**: Adapter, GORM, Port

### Service

**Definition**: Component that contains business logic and orchestrates domain operations.

**Context**: Services are implemented in `internal/domain/` and use repositories to persist data.

**Example**:
```go
type userService struct {
    userRepo interfaces.UserRepository
    logger   *zap.Logger
}

func (s *userService) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error) {
    // Business validation
    if req.Email == "" {
        return nil, errors.New("email required")
    }

    // Password hashing
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

    // Creation via repository
    user := &models.User{Email: req.Email, Password: string(hashedPassword)}
    err = s.userRepo.Create(ctx, user)

    return user, err
}
```

**See also**: Domain, Repository, Business Logic

---

## Go Concepts

### Channel

**Definition**: Native Go type enabling communication between goroutines in a thread-safe manner.

**Context**: Used for asynchronous communication and synchronization between goroutines in generated projects.

**Example**:
```go
// Channel to signal graceful shutdown
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

<-quit // Blocks until a signal is received
```

**See also**: Goroutine, Concurrency

### Context

**Definition**: Go type that carries deadlines, cancellation signals, and request-scoped values across API boundaries.

**Context**: Used in all service and repository functions to manage timeouts and cancellations.

**Example**:
```go
func (s *userService) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
    return s.userRepo.FindByID(ctx, id)
}

// Usage with timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
user, err := userService.GetUserByID(ctx, 1)
```

**See also**: Goroutine, Timeout

### Goroutine

**Definition**: Lightweight thread managed by the Go runtime enabling concurrent execution.

**Context**: Used to start the HTTP server in a non-blocking manner during graceful shutdown.

**Example**:
```go
// Starting the server in a goroutine
go func() {
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        logger.Fatal("Server failed", zap.Error(err))
    }
}()

// Waiting for shutdown signal
<-quit
```

**See also**: Channel, Concurrency

### Interface

**Definition**: Go type that defines a set of method signatures. A value satisfies an interface if it implements all of its methods.

**Context**: Fundamental in the hexagonal architecture of go-starter-kit. Interfaces define the ports.

**Example**:
```go
// Interface definition (port)
type UserService interface {
    CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error)
}

// Implicit implementation
type userService struct {}

func (s *userService) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error) {
    // Implementation
}
```

**See also**: Port, Adapter, Duck Typing

### Pointer

**Definition**: Variable that stores the memory address of another variable rather than its value.

**Context**: Widely used to avoid copies of large structures and to allow in-place modifications.

**Example**:
```go
// Return a pointer to avoid copying
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
    var user models.User
    err := r.db.WithContext(ctx).First(&user, id).Error
    return &user, err // Returns a pointer
}

// Pointer receiver to modify the structure
func (u *User) SetEmail(email string) {
    u.Email = email // Modifies the original
}
```

**See also**: Struct, Method Receiver

### Struct

**Definition**: Go composite type that groups fields of different types under a single name.

**Context**: Used to define entities (models), services, repositories, handlers, and configurations.

**Example**:
```go
// Entity struct
type User struct {
    ID        uint   `gorm:"primaryKey" json:"id"`
    Email     string `gorm:"uniqueIndex;not null" json:"email"`
    Password  string `gorm:"not null" json:"-"`
}

// Service struct
type userService struct {
    userRepo interfaces.UserRepository
    logger   *zap.Logger
}
```

**See also**: Pointer, Interface, Tag

### Struct Tag

**Definition**: Metadata attached to struct fields, used by reflection to configure behavior (serialization, validation, ORM).

**Context**: Used for GORM (DB mapping), JSON (serialization), and validation.

**Example**:
```go
type User struct {
    ID        uint   `gorm:"primaryKey" json:"id"`
    Email     string `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
    Password  string `gorm:"not null" json:"-"` // Excluded from JSON
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
```

**See also**: Struct, GORM, JSON Marshaling

---

## Database

### GORM

**Definition**: Go ORM (Object-Relational Mapping) that facilitates database operations by mapping Go structs to SQL tables.

**Context**: GORM is the default ORM in projects generated by go-starter-kit. It handles migrations, queries, and relationships.

**Example**:
```go
// Initialization
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

// Automatic migration
db.AutoMigrate(&models.User{}, &models.RefreshToken{})

// Query
var user models.User
db.Where("email = ?", "user@example.com").First(&user)
```

**See also**: ORM, Migration, PostgreSQL

### Migration

**Definition**: Process of versioning and applying database schema changes in a controlled manner.

**Context**: GORM AutoMigrate is used at application startup to synchronize the DB schema with Go models.

**Example**:
```go
// internal/infrastructure/database.go
func InitDB(cfg *config.Config) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})

    // Automatic migration
    if err := db.AutoMigrate(
        &models.User{},
        &models.RefreshToken{},
    ); err != nil {
        return nil, err
    }

    return db, nil
}
```

**See also**: GORM, Schema, AutoMigrate

### ORM

**Definition**: Object-Relational Mapping - programming technique that converts data between incompatible type systems (objects vs relational tables).

**Context**: GORM is the ORM used in go-starter-kit to abstract SQL access and map Go structs to PostgreSQL tables.

**See also**: GORM, Repository, PostgreSQL

### PostgreSQL

**Definition**: Open-source relational database management system, robust and compliant with SQL standards.

**Context**: Default database for generated projects. Configured via Docker or local installation.

**Configuration**:
```bash
# Via Docker
docker run -d --name postgres \
  -e POSTGRES_DB=myproject \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine

# Connection string in .env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/myproject?sslmode=disable
```

**See also**: GORM, Migration, Database URL

### Transaction

**Definition**: Set of database operations that execute atomically (all or nothing).

**Context**: Used to ensure data consistency during multiple operations.

**Example**:
```go
func (s *userService) CreateUserWithProfile(ctx context.Context, req CreateUserRequest) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Create the user
        user := &models.User{Email: req.Email}
        if err := tx.Create(user).Error; err != nil {
            return err // Automatic rollback
        }

        // Create the profile
        profile := &models.Profile{UserID: user.ID}
        if err := tx.Create(profile).Error; err != nil {
            return err // Automatic rollback
        }

        return nil // Automatic commit
    })
}
```

**See also**: GORM, ACID, Rollback

---

## Security

### Authentication

**Definition**: Process of verifying a user's identity (who are you?).

**Context**: Implemented via JWT in go-starter-kit. The user provides email/password and receives a JWT access token.

**Example**:
```go
func (s *authService) Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error) {
    // Verify email
    user, err := s.userRepo.FindByEmail(ctx, req.Email)

    // Verify password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        return nil, errors.New("invalid credentials")
    }

    // Generate JWT
    token, err := s.tokenService.GenerateToken(user.ID)
    return &models.AuthResponse{AccessToken: token}, nil
}
```

**See also**: Authorization, JWT, bcrypt

### Authorization

**Definition**: Process of verifying the permissions of an authenticated user (what can you do?).

**Context**: Implemented via JWT middleware that verifies the token and extracts the user ID to check permissions.

**Example**:
```go
// Authorization middleware
func AuthMiddleware(tokenService interfaces.TokenService) fiber.Handler {
    return func(c *fiber.Ctx) error {
        token := extractBearerToken(c)
        userID, err := tokenService.ValidateToken(token)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
        }
        c.Locals("userID", userID)
        return c.Next()
    }
}
```

**See also**: Authentication, Middleware, JWT

### bcrypt

**Definition**: Adaptive password hashing algorithm resistant to brute-force attacks thanks to a configurable computation cost.

**Context**: Used to hash passwords before storing them in the database.

**Example**:
```go
import "golang.org/x/crypto/bcrypt"

// Hash a password
hashedPassword, err := bcrypt.GenerateFromPassword(
    []byte(plainPassword),
    bcrypt.DefaultCost, // Cost 10
)

// Verify a password
err := bcrypt.CompareHashAndPassword(
    []byte(hashedPassword),
    []byte(plainPassword),
)
if err != nil {
    // Invalid password
}
```

**See also**: Hash, Authentication, Password Security

### Bearer Token

**Definition**: Type of authentication token transmitted in the HTTP header `Authorization: Bearer <token>`.

**Context**: JWT is transmitted as a Bearer token in HTTP requests to authenticate the user.

**Example**:
```http
GET /api/users/me HTTP/1.1
Host: localhost:8080
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

```go
// Bearer token extraction
func extractBearerToken(c *fiber.Ctx) string {
    authHeader := c.Get("Authorization")
    if strings.HasPrefix(authHeader, "Bearer ") {
        return strings.TrimPrefix(authHeader, "Bearer ")
    }
    return ""
}
```

**See also**: JWT, Authorization Header, HTTP

### Hash

**Definition**: One-way cryptographic function that transforms data into a fixed-size fingerprint.

**Context**: Used via bcrypt to hash passwords. Hashes are one-way (you cannot recover the original password).

**Example**:
```go
// Password hash (bcrypt)
hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
// Result: $2a$10$N9qo8uLOickgx2ZMRZoMye...

// Impossible to recover "password123" from the hash
// You can only verify:
bcrypt.CompareHashAndPassword(hash, []byte("password123")) // nil (match)
bcrypt.CompareHashAndPassword(hash, []byte("wrongpass"))   // error (no match)
```

**See also**: bcrypt, Encryption, One-way Function

### JWT

**Definition**: JSON Web Token - open standard (RFC 7519) for creating signed access tokens that contain JSON claims.

**Context**: Used for stateless authentication. The token contains the user ID and is signed with JWT_SECRET.

**Structure**: `header.payload.signature`

**Example**:
```go
import "github.com/golang-jwt/jwt/v5"

// Generation
type Claims struct {
    UserID uint `json:"user_id"`
    jwt.RegisteredClaims
}

token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
    UserID: 123,
    RegisteredClaims: jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
    },
})

signedToken, err := token.SignedString([]byte(jwtSecret))

// Validation
parsedToken, err := jwt.ParseWithClaims(signedToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
    return []byte(jwtSecret), nil
})
```

**See also**: Bearer Token, Authentication, Claims

### Refresh Token

**Definition**: Long-lived token used to obtain new access tokens without re-authentication.

**Context**: Stored in the database (refresh_tokens table) and associated with a user. Allows regenerating short-lived JWTs.

**Example**:
```go
type RefreshToken struct {
    ID        uint      `gorm:"primaryKey"`
    Token     string    `gorm:"uniqueIndex;not null"`
    UserID    uint      `gorm:"not null"`
    ExpiresAt time.Time `gorm:"not null"`
    CreatedAt time.Time
}

// Usage
func (s *authService) RefreshAccessToken(ctx context.Context, refreshToken string) (*models.AuthResponse, error) {
    // Verify the refresh token
    rt, err := s.tokenRepo.FindByToken(ctx, refreshToken)
    if err != nil || rt.ExpiresAt.Before(time.Now()) {
        return nil, errors.New("invalid refresh token")
    }

    // Generate new access token
    accessToken, err := s.tokenService.GenerateToken(rt.UserID)
    return &models.AuthResponse{AccessToken: accessToken}, nil
}
```

**See also**: JWT, Token Rotation, Authentication

---

## HTTP/API

### Endpoint

**Definition**: Specific URL exposed by an API to perform a particular operation.

**Context**: Endpoints are defined in `internal/adapters/http/routes.go` and mapped to handlers.

**Example**:
```go
// Endpoint definitions
func SetupRoutes(app *fiber.App, userHandler *UserHandler, authHandler *AuthHandler) {
    api := app.Group("/api")

    // Public endpoints
    api.Post("/auth/register", authHandler.Register)
    api.Post("/auth/login", authHandler.Login)

    // Protected endpoints
    protected := api.Group("", AuthMiddleware(tokenService))
    protected.Get("/users/me", userHandler.GetCurrentUser)
    protected.Put("/users/me", userHandler.UpdateCurrentUser)
}
```

**See also**: Handler, Router, REST

### Fiber

**Definition**: High-performance HTTP web framework for Go, built on fasthttp, offering a fast router and middleware features with an Express.js-inspired API.

**Context**: Default HTTP framework in go-starter-kit. Used to define routes, handlers, and middlewares.

**Example**:
```go
import "github.com/gofiber/fiber/v2"

app := fiber.New()

// Simple route
app.Get("/ping", func(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{"message": "pong"})
})

// Route with parameter
app.Get("/users/:id", func(c *fiber.Ctx) error {
    id := c.Params("id")
    return c.JSON(fiber.Map{"user_id": id})
})

app.Listen(":8080")
```

**See also**: Handler, Router, Middleware

### Handler

**Definition**: Function that processes an HTTP request for a specific endpoint.

**Context**: Handlers are defined in `internal/adapters/http/` and call domain services.

**Example**:
```go
type UserHandler struct {
    userService interfaces.UserService
    logger      *zap.Logger
}

func (h *UserHandler) GetCurrentUser(c *fiber.Ctx) error {
    userID := c.Locals("userID").(uint)

    user, err := h.userService.GetUserByID(c.Context(), userID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "server error"})
    }

    return c.JSON(user)
}
```

**See also**: Endpoint, Fiber, Middleware

### Middleware

**Definition**: Function that executes before or after a handler to perform cross-cutting operations (auth, logging, CORS, etc.).

**Context**: Middlewares are used for JWT authentication, logging, and CORS management.

**Example**:
```go
// Authentication middleware
func AuthMiddleware(tokenService interfaces.TokenService) fiber.Handler {
    return func(c *fiber.Ctx) error {
        token := extractBearerToken(c)
        userID, err := tokenService.ValidateToken(token)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
        }
        c.Locals("userID", userID)
        return c.Next() // Pass to the next handler
    }
}

// Usage
protected := app.Group("/api", AuthMiddleware(tokenService))
protected.Get("/users/me", userHandler.GetCurrentUser)
```

**See also**: Handler, Fiber, Authorization

### REST

**Definition**: Representational State Transfer - architectural style for web APIs using HTTP and its methods (GET, POST, PUT, DELETE).

**Context**: Generated APIs follow REST conventions: resources as URLs, HTTP methods for CRUD, stateless.

**Conventions**:
```
GET    /api/users       # List users
GET    /api/users/:id   # User details
POST   /api/users       # Create a user
PUT    /api/users/:id   # Update a user
DELETE /api/users/:id   # Delete a user
```

**See also**: Endpoint, HTTP Methods, CRUD

### Router

**Definition**: Component that maps URLs and HTTP methods to the appropriate handlers.

**Context**: Fiber Router is configured in `internal/adapters/http/routes.go`.

**Example**:
```go
func SetupRoutes(app *fiber.App, handlers *Handlers) {
    // API group
    api := app.Group("/api")

    // Authentication subgroup
    auth := api.Group("/auth")
    auth.Post("/register", handlers.Auth.Register)
    auth.Post("/login", handlers.Auth.Login)

    // Users subgroup (protected)
    users := api.Group("/users", AuthMiddleware(handlers.TokenService))
    users.Get("/me", handlers.User.GetCurrentUser)
    users.Put("/me", handlers.User.UpdateCurrentUser)
}
```

**See also**: Fiber, Endpoint, Handler

---

## See Also

<i class="material-icons info">info</i> **Additional documentation**:
- [Generated Project Guide](/go-starter-kit/generated-project-guide/) - Complete guide for generated projects
- [CLI Architecture](/go-starter-kit/cli-architecture/) - Project generator architecture
- [Usage Guide](/go-starter-kit/usage/) - Using create-go-starter

<i class="material-icons">menu_book</i> **External resources**:
- [Go Documentation](https://go.dev/doc/) - Official Go documentation
- [Fiber Documentation](https://docs.gofiber.io/) - HTTP framework
- [GORM Documentation](https://gorm.io/docs/) - ORM for Go
- [JWT RFC 7519](https://datatracker.ietf.org/doc/html/rfc7519) - JWT specification
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/) - Original article by Alistair Cockburn
