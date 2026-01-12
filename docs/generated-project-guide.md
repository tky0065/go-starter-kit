# Guide des projets générés

> Guide complet pour développer, tester et déployer les projets créés avec create-go-starter

## Table des matières

1. [Architecture](#architecture)
2. [Configuration](#configuration)
3. [Développement](#développement)
4. [API Reference](#api-reference)
5. [Tests](#tests)
6. [Base de données](#base-de-données)
7. [Sécurité](#sécurité)
8. [Déploiement](#déploiement)
9. [Monitoring & Logging](#monitoring--logging)
10. [Bonnes pratiques](#bonnes-pratiques)

---

## Architecture

### Architecture hexagonale (Ports & Adapters)

Les projets générés suivent l'architecture hexagonale, également appelée "Ports and Adapters".

**Principe fondamental**: Le domaine métier (business logic) est au centre et ne dépend de rien. Toutes les dépendances pointent vers le domaine.

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

### Diagramme d'architecture complète (Mermaid)

Le diagramme suivant montre l'architecture hexagonale complète avec tous les composants et leurs interactions :

```mermaid
flowchart TB
    subgraph External["Monde Externe"]
        Client["Client HTTP<br/>(Web, Mobile, API)"]
        DB[("PostgreSQL<br/>Database")]
    end

    subgraph Adapters["Adapters Layer"]
        direction TB
        subgraph Inbound["Inbound Adapters (Entree)"]
            Handlers["Handlers<br/>AuthHandler<br/>UserHandler"]
            Middleware["Middleware<br/>AuthMiddleware<br/>ErrorHandler"]
        end
        subgraph Outbound["Outbound Adapters (Sortie)"]
            RepoImpl["Repository GORM<br/>UserRepository"]
        end
    end

    subgraph Core["Core Business (Hexagone)"]
        direction TB
        Models["Models (Entites)<br/>User<br/>RefreshToken"]
        Domain["Domain Services<br/>UserService<br/>Business Logic"]
        Interfaces["Interfaces/Ports<br/>UserRepository<br/>UserService"]
        Errors["Domain Errors<br/>NotFound<br/>Validation<br/>Conflict"]
    end

    subgraph Infrastructure["Infrastructure Layer"]
        Server["Fiber Server<br/>Routes et Config"]
        DBConn["Database Connection<br/>GORM Setup"]
        Config["Configuration<br/>Environment vars"]
    end

    subgraph Packages["Packages Reutilisables (pkg/)"]
        Auth["Auth Package<br/>JWT Generation<br/>Token Parsing"]
        Logger["Logger Package<br/>Zerolog Config"]
        ConfigPkg["Config Package<br/>Env Loading"]
    end

    Client -->|"HTTP Request"| Server
    Server -->|"Route"| Handlers
    Handlers --> Middleware
    Handlers -->|"Appelle"| Domain
    Domain -->|"Utilise"| Interfaces
    Domain -->|"Utilise"| Models
    Domain -->|"Retourne"| Errors
    RepoImpl -.->|"Implemente"| Interfaces
    RepoImpl -->|"Utilise"| Models
    RepoImpl -->|"Query"| DBConn
    DBConn -->|"SQL"| DB
    Handlers -->|"Utilise"| Auth
    Server -->|"Utilise"| Config
    Domain -->|"Utilise"| Logger
```

### Flux d'une requete HTTP (Sequence Diagram)

Ce diagramme montre le parcours complet d'une requete HTTP a travers l'architecture :

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
    S->>M: Route vers Handler
    M->>M: Validation (si protege)
    M->>H: Requete validee
    
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
    P->>R: Appel implementation
    
    rect rgb(240, 255, 240)
        Note over R: Repository Layer
        R->>DB: INSERT INTO users...
        DB-->>R: User cree (ID)
    end
    
    R-->>SVC: User entity
    SVC-->>H: User + nil error
    H->>H: Generate JWT tokens
    H-->>C: HTTP 201 + JSON Response
```

### Principe de l'Inversion de Dependances

Le coeur de l'architecture hexagonale repose sur l'**Inversion de Dependances** :

```mermaid
flowchart LR
    subgraph Traditional["Approche Traditionnelle"]
        direction TB
        T_Handler["Handler"] --> T_Service["Service"]
        T_Service --> T_Repo["Repository"]
        T_Repo --> T_DB["Database"]
    end

    subgraph Hexagonal["Architecture Hexagonale"]
        direction TB
        H_Handler["Handler"]
        H_Service["Service"]
        H_Interface["Interface<br/>(Port)"]
        H_Repo["Repository<br/>(Adapter)"]
        H_DB["Database"]
        
        H_Handler --> H_Service
        H_Service --> H_Interface
        H_Repo -.->|"implemente"| H_Interface
        H_Repo --> H_DB
    end
```

**Avantages de cette approche** :

| Aspect | Sans Hexagonal | Avec Hexagonal |
|--------|----------------|----------------|
| **Testabilite** | Difficile (depend de la DB) | Facile (mock des interfaces) |
| **Changement de DB** | Modifications partout | Seulement le repository |
| **Changement de framework** | Refactoring complet | Seulement les handlers |
| **Logique metier** | Dispersee | Centralisee dans le domain |

### Structure des fichiers et responsabilites

```mermaid
flowchart TD
    subgraph CMD["cmd/"]
        Main["main.go<br/>Bootstrap fx.New()"]
    end

    subgraph Internal["internal/"]
        subgraph Models["models/"]
            User["user.go<br/>Entites GORM"]
        end
        
        subgraph Domain["domain/"]
            DErrors["errors.go<br/>Erreurs metier"]
            subgraph UserDomain["user/"]
                Service["service.go<br/>Logique metier"]
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
                Routes["routes.go<br/>Routes centralisees"]
            end
            subgraph MiddlewarePkg["middleware/"]
                AuthM["auth_middleware.go"]
                ErrorM["error_handler.go"]
            end
            subgraph RepoPkg["repository/"]
                UserRepo["user_repository.go<br/>Implementation GORM"]
            end
        end
        
        subgraph Infra["infrastructure/"]
            DBPkg["database/<br/>Connexion GORM"]
            ServerPkg["server/<br/>Config Fiber"]
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

**Flux de données**:

1. **Requête HTTP** → Handler (adapters/handlers)
2. **Handler** → Appelle le Service via l'interface (domain)
3. **Service** → Exécute la logique métier, appelle le Repository via l'interface
4. **Repository** → Persiste dans la DB (infrastructure)
5. **Retour** → Remonte jusqu'au Handler qui retourne la réponse HTTP

**Avantages**:

- **Testabilité**: Le domaine peut être testé sans DB ni HTTP
- **Flexibilité**: Changement de DB (PostgreSQL → MySQL) ou framework (Fiber → Gin) facile
- **Maintenabilité**: Séparation claire des responsabilités
- **Évolutivité**: Ajout de nouvelles fonctionnalités sans casser l'existant

### Stack technique

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
- **Maintenant**: `interfaces` → `models` ← `domain/user` (✅ pas de cycle)
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

Le serveur crée l'application Fiber et gère le lifecycle. Les routes sont enregistrées via `handlers.Module` qui appelle `http.RegisterRoutes()`.

```go
func NewServer(logger zerolog.Logger, db *gorm.DB) *fiber.App {
    app := fiber.New(fiber.Config{
        AppName:      "mon-projet",
        ErrorHandler: middleware.ErrorHandler,
    })

    logger.Info().Msg("Fiber server initialized with centralized error handler")

    // Routes are registered via handlers.Module -> http.RegisterRoutes()

    return app
}

// registerHooks registers lifecycle hooks for server startup and shutdown
func registerHooks(lifecycle fx.Lifecycle, app *fiber.App, logger zerolog.Logger) {
    lifecycle.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            port := config.GetEnv("APP_PORT", "3000")
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

## Configuration

### Variables d'environnement

Le fichier `.env` contient toute la configuration:

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
JWT_SECRET=          # À REMPLIR!
JWT_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=168h
```

### Générer un JWT_SECRET sécurisé

**CRITIQUE**: Générez toujours un secret fort:

```bash
openssl rand -base64 32
```

Exemple de résultat:
```
XqR7nP3vY2kL9wH4sT6mU8jC1bN5aD0f
```

Ajoutez-le dans `.env`:
```
JWT_SECRET=XqR7nP3vY2kL9wH4sT6mU8jC1bN5aD0f
```

### Configuration par environnement

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
JWT_SECRET=<secret-depuis-vault>
```

#### Production

```bash
APP_ENV=production
DB_HOST=prod-db.example.com
DB_SSLMODE=require
DB_PASSWORD=<secret-depuis-secrets-manager>
JWT_SECRET=<secret-depuis-secrets-manager>
```

**Best practice**: Utiliser des secrets managers:

- **AWS**: Secrets Manager, Parameter Store
- **GCP**: Secret Manager
- **Kubernetes**: Secrets
- **HashiCorp**: Vault

### Configuration PostgreSQL

#### Option 1: PostgreSQL local

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

Si un `docker-compose.yml` est généré:

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

Lancer:

```bash
docker-compose up -d postgres
```

#### Vérification de connexion

```bash
# Avec psql
psql -h localhost -U postgres -d mon-projet

# Ou tester depuis l'app
make run
# Vérifier les logs: "Database connected successfully"
```

---

## Développement

### Workflow quotidien

**1. Lancer la base de données**

```bash
# Docker
docker start postgres
# ou
docker-compose up -d postgres

# Local
brew services start postgresql  # macOS
sudo systemctl start postgresql # Linux
```

**2. Lancer l'application**

```bash
make run
```

Ou avec hot-reload (si `air` installé):

```bash
# Installer air
go install github.com/cosmtrek/air@latest

# Lancer avec hot-reload
air
```

**3. Développer**

- Modifier le code
- Sauvegarder (auto-reload avec air)
- Vérifier les logs

**4. Tester**

```bash
# Tests unitaires
make test

# Tests avec coverage
make test-coverage

# Ouvrir le rapport
open coverage.html
```

**5. Linter**

```bash
make lint
```

### Commandes Makefile

| Commande | Description |
|----------|-------------|
| `make help` | Afficher l'aide |
| `make run` | Lancer l'app |
| `make build` | Build binaire |
| `make test` | Tests avec race detector |
| `make test-coverage` | Tests + rapport HTML |
| `make lint` | golangci-lint |
| `make clean` | Nettoyer artifacts |
| `make docker-build` | Build image Docker |
| `make docker-run` | Run conteneur Docker |

### Ajouter une nouvelle fonctionnalite

Cette section vous guide pas a pas pour ajouter une nouvelle entite/fonctionnalite en respectant l'architecture hexagonale.

#### Vue d'ensemble des 9 etapes

```mermaid
flowchart LR
    A["1. Model"] --> B["2. Interface"]
    B --> C["3. Repository"]
    B --> D["4. Service"]
    C --> E["5. Module fx"]
    D --> E
    D --> F["6. Handler"]
    F --> G["7. Routes"]
    A --> H["8. Migration"]
    E --> I["9. Bootstrap"]
    G --> I
```

#### Checklist rapide

Utilisez cette checklist pour ne rien oublier :

| Etape | Fichier a creer/modifier | Depend de | Status |
|-------|-------------------------|-----------|--------|
| 1. Model | `internal/models/<entity>.go` | - | [ ] |
| 2. Interface | `internal/interfaces/<entity>_repository.go` | Etape 1 | [ ] |
| 3. Repository | `internal/adapters/repository/<entity>_repository.go` | Etapes 1, 2 | [ ] |
| 4. Service | `internal/domain/<entity>/service.go` | Etapes 1, 2 | [ ] |
| 5. Module fx | `internal/domain/<entity>/module.go` | Etapes 3, 4 | [ ] |
| 6. Handler | `internal/adapters/handlers/<entity>_handler.go` | Etapes 1, 4 | [ ] |
| 7. Routes | `internal/infrastructure/server/server.go` (modifier) | Etape 6 | [ ] |
| 8. Migration | `internal/infrastructure/database/database.go` (modifier) | Etape 1 | [ ] |
| 9. Bootstrap | `cmd/main.go` (modifier) | Etape 5 | [ ] |

#### Diagramme des dependances entre fichiers

Ce diagramme montre l'ordre de creation des fichiers et leurs dependances :

```mermaid
flowchart TD
    subgraph Step1["Etape 1 - Foundation"]
        Model["models/product.go<br/>Entite GORM"]
    end

    subgraph Step2["Etape 2 - Abstraction"]
        Interface["interfaces/product_repository.go<br/>Port (contrat)"]
    end

    subgraph Step34["Etapes 3 et 4 - Implementation"]
        Repo["repository/product_repository.go<br/>Adapter GORM"]
        Service["domain/product/service.go<br/>Business Logic"]
    end

    subgraph Step5["Etape 5 - DI"]
        Module["domain/product/module.go<br/>fx.Module"]
    end

    subgraph Step6["Etape 6 - HTTP"]
        Handler["handlers/product_handler.go<br/>REST endpoints"]
    end

    subgraph Step789["Etapes 7, 8, 9 - Integration"]
        Routes["server/server.go<br/>Ajouter routes"]
        Migration["database/database.go<br/>AutoMigrate"]
        Main["cmd/main.go<br/>Ajouter module"]
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

### Exemple complet : Entite Product

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

## API Reference

### Endpoints disponibles

#### Health Check

```
GET /health
```

**Response** (200):
```json
{
  "status": "ok"
}
```

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

**Exemple curl**:
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

**Exemple curl**:
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

**Exemple curl**:
```bash
REFRESH_TOKEN="<refresh_token_from_login>"
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\":\"$REFRESH_TOKEN\"}"
```

---

### Users (Protected)

Tous les endpoints users requièrent un JWT token valide.

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

**Exemple curl**:
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

**Exemple curl**:
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

**Note**: Utilise soft delete (DeletedAt), les données restent en DB.

---

### Workflow complet avec l'API

```bash
# 1. Register
REGISTER_RESP=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}')

# Extraire access_token
ACCESS_TOKEN=$(echo $REGISTER_RESP | jq -r '.data.access_token')
REFRESH_TOKEN=$(echo $REGISTER_RESP | jq -r '.data.refresh_token')

# 2. List users (avec token)
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer $ACCESS_TOKEN"

# 3. Update user
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"email":"updated@example.com"}'

# 4. Quand access token expire (15min), utiliser refresh token
NEW_ACCESS=$(curl -s -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\":\"$REFRESH_TOKEN\"}" | jq -r '.data.access_token')

# 5. Continuer avec nouveau token
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer $NEW_ACCESS"
```

---

## Tests

### Organisation des tests

Les tests sont co-localisés avec le code source:

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

### Exécuter les tests

```bash
# Tous les tests
make test

# Tests avec coverage
make test-coverage

# Ouvrir le rapport HTML
open coverage.html  # macOS
xdg-open coverage.html  # Linux

# Tests d'un package spécifique
go test -v ./internal/domain/user

# Test spécifique
go test -run TestRegister ./internal/adapters/handlers

# Tests avec race detector (détection de race conditions)
go test -race ./...
```

### Types de tests

#### 1. Tests unitaires

Testent une fonction ou méthode isolée, avec mocks.

**Exemple: Test du service**

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

#### 2. Tests d'intégration

Testent plusieurs composants ensemble, avec DB réelle (SQLite in-memory).

**Exemple: Test handler avec DB**

```go
// internal/adapters/handlers/auth_handler_integration_test.go
func TestAuthHandler_RegisterIntegration(t *testing.T) {
    // Setup DB in-memory
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

#### 3. Tests table-driven

Pour tester plusieurs cas avec une structure commune:

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

### Best practices pour les tests

1. **Utilisez testify/assert** pour assertions claires:
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

3. **Mock les dépendances externes**:
   - DB (sauf pour tests d'intégration)
   - APIs externes
   - Services tiers

4. **Tests d'intégration avec SQLite in-memory**:
   ```go
   db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
   ```

5. **Clean up après chaque test**:
   ```go
   t.Cleanup(func() {
       db.Exec("DELETE FROM users")
   })
   ```

6. **Noms de tests descriptifs**:
   ```go
   func TestUserService_Register_WhenEmailAlreadyExists_ReturnsConflictError(t *testing.T)
   ```

### Coverage

Objectif: **> 80% coverage**

```bash
# Générer rapport
make test-coverage

# Voir coverage par package
go test -cover ./...

# Output:
# ok      mon-projet/internal/domain/user         0.123s  coverage: 85.7% of statements
# ok      mon-projet/internal/adapters/handlers   0.234s  coverage: 92.3% of statements
```

---

## Base de données

### Migrations

Le projet utilise GORM AutoMigrate pour simplifier les migrations en développement:

```go
// internal/infrastructure/database/database.go
db.AutoMigrate(
    &models.User{},
    &models.RefreshToken{},
)
```

**AutoMigrate**:
- Crée les tables si elles n'existent pas
- Ajoute les colonnes manquantes
- Crée les indexes
- **NE supprime PAS** de colonnes ou tables

**Pour la production**, considérer une solution de migrations versionnées:

- **golang-migrate/migrate**: Migrations SQL ou Go
- **pressly/goose**: Migrations up/down
- **GORM Migrator avancé**: API programmatique

**Exemple avec golang-migrate**:

```bash
# Installer migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Créer migration
migrate create -ext sql -dir migrations -seq create_users_table

# Fichiers créés:
# migrations/000001_create_users_table.up.sql
# migrations/000001_create_users_table.down.sql

# Run migrations
migrate -path migrations -database "postgresql://user:pass@localhost/dbname?sslmode=disable" up
```

### Modèles GORM

Conventions et patterns:

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

**Tags GORM importants**:

| Tag | Description |
|-----|-------------|
| `primarykey` | Clé primaire |
| `uniqueIndex` | Index unique |
| `index` | Index simple |
| `not null` | Colonne NOT NULL |
| `default:value` | Valeur par défaut |
| `size:255` | Taille de colonne |
| `type:varchar(100)` | Type SQL custom |
| `foreignKey:UserID` | Clé étrangère |
| `references:ID` | Référence FK |

**Conventions**:

- **Soft deletes**: `DeletedAt gorm.DeletedAt`
- **Timestamps auto**: `CreatedAt`, `UpdatedAt`
- **JSON hiding**: `json:"-"` pour password
- **Index sur FK**: Toujours indexer les clés étrangères

### Queries avancées

#### Pagination

```go
var users []models.User
limit := 10
offset := 20

db.Limit(limit).Offset(offset).Find(&users)
```

#### Filtering

```go
// Where simple
db.Where("email = ?", "user@example.com").First(&user)

// Where avec multiple conditions
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

### Performance tips

1. **Index les colonnes fréquemment requêtées**:
   ```go
   Email string `gorm:"uniqueIndex"`
   ```

2. **Éviter N+1 queries avec Preload**:
   ```go
   // ❌ N+1
   for _, user := range users {
       db.Model(&user).Association("RefreshTokens").Find(&tokens)
   }

   // ✅ Single query
   db.Preload("RefreshTokens").Find(&users)
   ```

3. **Select seulement les colonnes nécessaires**:
   ```go
   db.Select("id, email").Find(&users)
   ```

4. **Utiliser les connections pools**:
   ```go
   sqlDB, _ := db.DB()
   sqlDB.SetMaxIdleConns(10)
   sqlDB.SetMaxOpenConns(100)
   sqlDB.SetConnMaxLifetime(time.Hour)
   ```

---

## Sécurité

### Authentification JWT

#### Flow complet

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

#### Stockage des tokens (côté client)

**Access Token** (courte durée: 15min):
- **Recommandé**: En mémoire (variable JavaScript)
- **Pas de localStorage** (vulnérable à XSS)
- Perdu au refresh de page → Utiliser refresh token

**Refresh Token** (longue durée: 7j):
- **Option 1**: httpOnly cookie (le plus sécurisé)
- **Option 2**: localStorage (si pas de XSS risk)

**Exemple React**:

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

### Protection des routes

Middleware d'authentification:

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

**Utilisation**:

```go
// Protected routes
users := api.Group("/users")
users.Use(authMiddleware.Authenticate())  // Middleware appliqué
users.Get("/", userHandler.List)
```

**Dans le handler, récupérer user ID**:

```go
func (h *UserHandler) List(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    // Utiliser userID pour vérifier permissions, etc.
}
```

### Validation des entrées

go-playground/validator v10:

```go
type RegisterRequest struct {
    Email    string `json:"email" validate:"required,email,max=255"`
    Password string `json:"password" validate:"required,min=8,max=72"`
}

// Dans handler
validate := validator.New()
if err := validate.Struct(req); err != nil {
    return domain.NewValidationError("Invalid input", "VALIDATION_ERROR", err)
}
```

**Tags de validation courants**:

| Tag | Description |
|-----|-------------|
| `required` | Champ obligatoire |
| `email` | Format email valide |
| `min=N` | Longueur/valeur minimum |
| `max=N` | Longueur/valeur maximum |
| `len=N` | Longueur exacte |
| `gte=N` | Greater than or equal |
| `lte=N` | Less than or equal |
| `alpha` | Lettres seulement |
| `alphanum` | Lettres + chiffres |
| `numeric` | Chiffres seulement |
| `uuid` | Format UUID |
| `url` | Format URL |

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

### Hashage des mots de passe

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

**DefaultCost** = 10 (2^10 iterations) - Bon équilibre sécurité/performance

**Note**: Le hashage de mot de passe est géré dans le service (business logic), pas dans l'entité. L'entité `models.User` stocke le `PasswordHash` qui est toujours haché.

**Utilisation**:

```go
// Register (dans le service)
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

### Checklist sécurité

Production checklist:

- [ ] **JWT_SECRET fort**: Généré avec `openssl rand -base64 32`
- [ ] **HTTPS en production**: Toujours utiliser TLS
- [ ] **Rate limiting**: Implémenter avec fiber/limiter
  ```go
  import "github.com/gofiber/fiber/v2/middleware/limiter"

  app.Use(limiter.New(limiter.Config{
      Max:        100,
      Expiration: 1 * time.Minute,
  }))
  ```
- [ ] **CORS configuré**: Restreindre les origins
  ```go
  import "github.com/gofiber/fiber/v2/middleware/cors"

  app.Use(cors.New(cors.Config{
      AllowOrigins: "https://your-frontend.com",
      AllowHeaders: "Origin, Content-Type, Accept, Authorization",
  }))
  ```
- [ ] **Validation stricte**: Tous les inputs validés
- [ ] **SQL Injection**: GORM le prévient automatiquement
- [ ] **XSS**: Échapper les outputs HTML (si templates)
- [ ] **Logs sans secrets**: Jamais logger passwords, tokens
  ```go
  // ❌ BAD
  logger.Info().Str("password", password).Msg("User login")

  // ✅ GOOD
  logger.Info().Str("email", email).Msg("User login")
  ```
- [ ] **Environment variables**: Secrets dans .env (pas dans code)
- [ ] **DB SSL**: `DB_SSLMODE=require` en production
- [ ] **Helmet headers**: Implémenter security headers
  ```go
  import "github.com/gofiber/fiber/v2/middleware/helmet"

  app.Use(helmet.New())
  ```
- [ ] **Timeout requests**: Éviter DoS
  ```go
  app := fiber.New(fiber.Config{
      ReadTimeout:  10 * time.Second,
      WriteTimeout: 10 * time.Second,
      IdleTimeout:  30 * time.Second,
  })
  ```

---

## Déploiement

### Docker

#### Build de l'image

```bash
make docker-build
```

Ou manuellement:

```bash
docker build -t mon-projet:latest .
```

Le Dockerfile généré utilise un build multi-stage:

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

**Avantages**:
- Image finale légère (~15-20MB vs ~1GB)
- Sécurité (image alpine minimale)
- Binaire statique (pas de dépendances)

#### Run avec Docker

```bash
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PASSWORD=postgres \
  -e JWT_SECRET=<votre_secret> \
  mon-projet:latest
```

**Note**: `host.docker.internal` permet d'accéder à localhost depuis Docker.

#### Docker Compose

Si `docker-compose.yml` est généré:

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

**Lancer**:

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

Manifests basiques pour déploiement K8s:

#### Secret

```yaml
# k8s/secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: mon-projet-secret
type: Opaque
stringData:
  jwt-secret: "<votre_secret_base64>"
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
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /health
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

#### Déployer PostgreSQL (StatefulSet)

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

### CI/CD avec GitHub Actions

Le workflow généré (`.github/workflows/ci.yml`):

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
2. **Test**: Tests avec PostgreSQL (service container)
3. **Build**: Vérification build

### Déploiement en production

#### Checklist pré-déploiement

- [ ] Tous les tests passent (`make test`)
- [ ] Lint passe (`make lint`)
- [ ] Variables d'environnement configurées
- [ ] JWT_SECRET généré (fort, aléatoire)
- [ ] DB_SSLMODE=require
- [ ] Migrations DB exécutées
- [ ] Health check fonctionne
- [ ] Logs configurés
- [ ] Monitoring en place

#### Plateformes recommandées

**1. Google Cloud Run** (le plus simple):

```bash
# Build et push image
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

**4. Kubernetes** (le plus flexible):

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

### Logging avec zerolog

#### Configuration

Le logger est configuré dans `pkg/logger/logger.go`:

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

#### Utilisation

Injection via fx:

```go
type UserService struct {
    logger zerolog.Logger
}

func NewUserService(logger zerolog.Logger) *UserService {
    return &UserService{logger: logger}
}
```

**Logging structuré**:

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

#### Niveaux de log

| Niveau | Usage |
|--------|-------|
| **Debug** | Informations détaillées pour debugging |
| **Info** | Événements importants (user login, etc.) |
| **Warn** | Comportements anormaux non-critiques |
| **Error** | Erreurs nécessitant attention |
| **Fatal** | Erreurs critiques (app exit) |

#### Best practices

**✅ BON - Structured logging**:

```go
logger.Info().
    Str("user_id", userID).
    Str("action", "login").
    Dur("duration", elapsed).
    Msg("User logged in")
```

**❌ MAUVAIS - String formatting**:

```go
logger.Info().Msgf("User %s logged in after %v", userID, elapsed)
```

**✅ BON - Pas de secrets**:

```go
logger.Info().Str("email", email).Msg("User login attempt")
```

**❌ MAUVAIS - Logging secrets**:

```go
logger.Info().Str("password", password).Msg("Login")  // NEVER!
```

### Monitoring (recommandations)

Pour la production, intégrer:

#### 1. Prometheus + Grafana

**Metrics à collecter**:
- Request rate (requests/sec)
- Response time (p50, p95, p99)
- Error rate (%)
- Database connection pool
- CPU/Memory usage

**Implémenter avec fiber/prometheus**:

```go
import "github.com/gofiber/adaptor/v2"
import "github.com/prometheus/client_golang/prometheus/promhttp"

app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
```

#### 2. Jaeger / OpenTelemetry

**Distributed tracing** pour suivre les requêtes à travers les services.

#### 3. Sentry

**Error tracking** en temps réel:

```go
import "github.com/getsentry/sentry-go"

sentry.Init(sentry.ClientOptions{
    Dsn: os.Getenv("SENTRY_DSN"),
})

// Capture errors
sentry.CaptureException(err)
```

#### 4. APM (Application Performance Monitoring)

- **New Relic**: APM complet
- **Datadog**: Monitoring + logs
- **Elastic APM**: Open source

### Health checks

L'endpoint `/health` est crucial pour:

- Load balancers
- Kubernetes probes
- Monitoring tools

**Amélioré**:

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

---

## Bonnes pratiques

### Architecture

**1. Domain isolation**

Le domaine ne doit **jamais** importer d'autres packages:

```go
// ❌ BAD - Domain importing adapter
package user

import "mon-projet/internal/adapters/repository"  // NO!

// ✅ GOOD - Domain only imports interfaces
package user

import "mon-projet/internal/interfaces"
```

**2. Single Responsibility Principle**

Chaque composant a une seule responsabilité:

- **Handlers**: Parse + validate + call service
- **Services**: Business logic uniquement
- **Repositories**: Data access uniquement

**3. Dependency Injection**

Toujours via fx.Provide, pas de variables globales:

```go
// ❌ BAD - Global variable
var db *gorm.DB

// ✅ GOOD - Injection
type UserService struct {
    db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
    return &UserService{db: db}
}
```

### Code style

**1. gofmt**

Toujours formater:

```bash
go fmt ./...
```

Ou configurer l'IDE pour formater à la sauvegarde.

**2. golangci-lint**

Respecter les règles:

```bash
make lint
```

**3. Documentation GoDoc**

Pour les exports publics:

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

**4. Error handling explicite**

Toujours gérer les erreurs, ne pas utiliser `panic`:

```go
// ❌ BAD
user := getUserByID(id)  // What if error?

// ✅ GOOD
user, err := getUserByID(id)
if err != nil {
    return nil, fmt.Errorf("failed to get user: %w", err)
}
```

### Naming conventions

**Interfaces**:
- Suffixe `-er` ou `-Service`
- Exemples: `UserRepository`, `AuthService`, `Logger`

**Repositories**:
- Suffixe `-Repository`
- Exemples: `UserRepository`, `ProductRepository`

**Handlers**:
- Suffixe `-Handler`
- Exemples: `AuthHandler`, `UserHandler`

**Constructeurs**:
- Préfixe `New`
- Exemples: `NewUserService`, `NewAuthHandler`

**Méthodes privées**:
- lowerCamelCase
- Exemples: `hashPassword`, `validateEmail`

### Error handling patterns

**Wrap errors avec contexte**:

```go
// ✅ GOOD
if err != nil {
    return fmt.Errorf("failed to create user %s: %w", email, err)
}
```

**Domain errors pour logique métier**:

```go
if user == nil {
    return domain.NewNotFoundError("User not found", "USER_NOT_FOUND", nil)
}
```

**Ne pas gérer les HTTP status dans le service**:

```go
// ❌ BAD - Service returning HTTP status
func (s *UserService) GetByID(id uint) (int, *User, error) {
    return 404, nil, errors.New("not found")
}

// ✅ GOOD - Service returning domain error
func (s *UserService) GetByID(id uint) (*User, error) {
    return nil, domain.NewNotFoundError("User not found", "USER_NOT_FOUND", nil)
}
```

### Testing best practices

**1. Coverage > 80%**

```bash
go test -cover ./...
```

**2. Tests table-driven**

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

**3. Noms descriptifs**

```go
func TestUserService_Register_WhenEmailAlreadyExists_ReturnsConflictError(t *testing.T)
```

**4. Setup/teardown avec t.Cleanup()**

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

**1. GORM - Éviter N+1 queries**

```go
// ❌ N+1 problem
for _, user := range users {
    db.Model(&user).Association("Posts").Find(&posts)
}

// ✅ Single query with Preload
db.Preload("Posts").Find(&users)
```

**2. Context - Toujours passer context.Context**

```go
func (s *UserService) GetByID(ctx context.Context, id uint) (*User, error) {
    return s.repo.FindByID(ctx, id)
}
```

**3. Database indexes**

```go
Email string `gorm:"uniqueIndex"`  // Index sur colonnes fréquemment requêtées
```

**4. Connection pooling**

```go
sqlDB, _ := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### Sécurité recap

- [ ] Valider toutes les entrées utilisateur
- [ ] Jamais logger de passwords ou tokens
- [ ] Rate limiting sur endpoints publics
- [ ] HTTPS en production
- [ ] JWT secret fort (32+ caractères)
- [ ] Bcrypt pour passwords
- [ ] Mettre à jour les dépendances régulièrement

```bash
# Vérifier vulnérabilités
go list -json -m all | nancy sleuth
```

---

## Conclusion

Ce guide couvre tous les aspects du développement avec les projets générés par `create-go-starter`. Pour aller plus loin:

- **Exemples de code**: Tous les patterns sont dans le code généré
- **Tests**: Regardez les fichiers `*_test.go` pour des exemples
- **Documentation officielle**:
  - [Fiber](https://docs.gofiber.io/)
  - [GORM](https://gorm.io/docs/)
  - [fx](https://uber-go.github.io/fx/)
  - [zerolog](https://github.com/rs/zerolog)

**Bon développement!** 🚀

Si vous rencontrez des problèmes ou avez des questions, consultez:
- [Issues GitHub](https://github.com/tky0065/go-starter-kit/issues)
- [Discussions GitHub](https://github.com/tky0065/go-starter-kit/discussions)
