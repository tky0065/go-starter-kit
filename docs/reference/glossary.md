# Glossaire Technique

Ce glossaire définit les termes techniques utilisés dans go-starter-kit et les projets générés.

---

## Architecture

### Adapter

**Définition** : Composant qui implémente une interface (port) pour connecter le domaine métier à une technologie externe (base de données, API, HTTP, etc.).

**Contexte** : Dans go-starter-kit, les adapters se trouvent dans `internal/adapters/` et incluent les handlers HTTP et les repositories de base de données.

**Exemple** :
```go
// Adapter HTTP (handler)
type UserHandler struct {
    userService interfaces.UserService
}

// Adapter de base de données (repository)
type UserRepository struct {
    db *gorm.DB
}
```

**Voir aussi** : Port, Hexagonal Architecture, Repository

### Domain

**Définition** : Couche centrale contenant la logique métier pure, indépendante des détails techniques d'infrastructure.

**Contexte** : Le domaine se trouve dans `internal/domain/` et contient les services métier (UserService, AuthService). Les entités métier sont dans `internal/models/`.

**Exemple** :
```go
// Service du domaine
type userService struct {
    userRepo interfaces.UserRepository
}

func (s *userService) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error) {
    // Logique métier pure
}
```

**Voir aussi** : Hexagonal Architecture, Service, Models

### Hexagonal Architecture

**Définition** : Pattern architectural (aussi appelé "Ports and Adapters") qui isole le domaine métier des détails techniques en utilisant des interfaces (ports) et des implémentations (adapters).

**Contexte** : go-starter-kit utilise l'architecture hexagonale pour garantir la séparation des préoccupations et faciliter les tests.

**Structure** :
```
internal/
├── models/       # Entités du domaine
├── domain/       # Logique métier
├── interfaces/   # Ports (interfaces)
├── adapters/     # Adapters (HTTP, DB)
└── infrastructure/ # Configuration
```

**Voir aussi** : Port, Adapter, Domain, Dependency Injection

### Models

**Définition** : Entités du domaine métier partagées entre toutes les couches de l'application.

**Contexte** : Les models sont définis dans `internal/models/` pour éviter les dépendances circulaires. Ils représentent les objets métier fondamentaux (User, RefreshToken, etc.).

**Exemple** :
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

**Voir aussi** : Domain, Entity

### Port

**Définition** : Interface qui définit un contrat pour interagir avec le domaine métier ou une ressource externe.

**Contexte** : Les ports sont définis dans `internal/interfaces/` et incluent UserService, UserRepository, AuthService, etc.

**Exemple** :
```go
// Port de service
type UserService interface {
    CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error)
    GetUserByID(ctx context.Context, id uint) (*models.User, error)
}

// Port de repository
type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    FindByEmail(ctx context.Context, email string) (*models.User, error)
}
```

**Voir aussi** : Adapter, Interface, Hexagonal Architecture

### Repository

**Définition** : Pattern qui encapsule la logique d'accès aux données et fournit une interface de collection pour manipuler les entités.

**Contexte** : Les repositories sont implémentés dans `internal/adapters/repositories/` et utilisent GORM pour l'accès à PostgreSQL.

**Exemple** :
```go
type UserRepository struct {
    db *gorm.DB
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}
```

**Voir aussi** : Adapter, GORM, Port

### Service

**Définition** : Composant qui contient la logique métier et orchestre les opérations du domaine.

**Contexte** : Les services sont implémentés dans `internal/domain/` et utilisent les repositories pour persister les données.

**Exemple** :
```go
type userService struct {
    userRepo interfaces.UserRepository
    logger   *zap.Logger
}

func (s *userService) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error) {
    // Validation métier
    if req.Email == "" {
        return nil, errors.New("email requis")
    }

    // Hash du mot de passe
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

    // Création via repository
    user := &models.User{Email: req.Email, Password: string(hashedPassword)}
    err = s.userRepo.Create(ctx, user)

    return user, err
}
```

**Voir aussi** : Domain, Repository, Business Logic

---

## Concepts Go

### Channel

**Définition** : Type natif Go permettant la communication entre goroutines de manière thread-safe.

**Contexte** : Utilisé pour la communication asynchrone et la synchronisation entre goroutines dans les projets générés.

**Exemple** :
```go
// Canal pour signaler l'arrêt gracieux
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

<-quit // Bloque jusqu'à réception d'un signal
```

**Voir aussi** : Goroutine, Concurrency

### Context

**Définition** : Type Go qui transporte les deadlines, signaux d'annulation et valeurs request-scoped à travers les frontières d'API.

**Contexte** : Utilisé dans toutes les fonctions de service et repository pour gérer les timeouts et annulations.

**Exemple** :
```go
func (s *userService) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
    return s.userRepo.FindByID(ctx, id)
}

// Utilisation avec timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
user, err := userService.GetUserByID(ctx, 1)
```

**Voir aussi** : Goroutine, Timeout

### Goroutine

**Définition** : Thread léger géré par le runtime Go permettant l'exécution concurrente.

**Contexte** : Utilisé pour démarrer le serveur HTTP de manière non-bloquante lors du graceful shutdown.

**Exemple** :
```go
// Démarrage du serveur dans une goroutine
go func() {
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        logger.Fatal("Server failed", zap.Error(err))
    }
}()

// Attente du signal d'arrêt
<-quit
```

**Voir aussi** : Channel, Concurrency

### Interface

**Définition** : Type Go qui définit un ensemble de signatures de méthodes. Une valeur satisfait une interface si elle implémente toutes ses méthodes.

**Contexte** : Fondamental dans l'architecture hexagonale de go-starter-kit. Les interfaces définissent les ports.

**Exemple** :
```go
// Définition d'interface (port)
type UserService interface {
    CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error)
}

// Implémentation implicite
type userService struct {}

func (s *userService) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error) {
    // Implémentation
}
```

**Voir aussi** : Port, Adapter, Duck Typing

### Pointer

**Définition** : Variable qui stocke l'adresse mémoire d'une autre variable plutôt que sa valeur.

**Contexte** : Largement utilisé pour éviter les copies de structures volumineuses et permettre les modifications in-place.

**Exemple** :
```go
// Retourner un pointeur pour éviter la copie
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
    var user models.User
    err := r.db.WithContext(ctx).First(&user, id).Error
    return &user, err // Retourne un pointeur
}

// Receiver pointeur pour modifier la structure
func (u *User) SetEmail(email string) {
    u.Email = email // Modifie l'original
}
```

**Voir aussi** : Struct, Method Receiver

### Struct

**Définition** : Type composite Go qui groupe des champs de différents types sous un seul nom.

**Contexte** : Utilisé pour définir les entités (models), services, repositories, handlers, et configurations.

**Exemple** :
```go
// Struct d'entité
type User struct {
    ID        uint   `gorm:"primaryKey" json:"id"`
    Email     string `gorm:"uniqueIndex;not null" json:"email"`
    Password  string `gorm:"not null" json:"-"`
}

// Struct de service
type userService struct {
    userRepo interfaces.UserRepository
    logger   *zap.Logger
}
```

**Voir aussi** : Pointer, Interface, Tag

### Struct Tag

**Définition** : Métadonnées attachées aux champs d'une struct, utilisées par reflection pour configurer le comportement (serialization, validation, ORM).

**Contexte** : Utilisé pour GORM (mapping DB), JSON (serialization), et validation.

**Exemple** :
```go
type User struct {
    ID        uint   `gorm:"primaryKey" json:"id"`
    Email     string `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
    Password  string `gorm:"not null" json:"-"` // Exclu du JSON
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
```

**Voir aussi** : Struct, GORM, JSON Marshaling

---

## Base de Données

### GORM

**Définition** : ORM (Object-Relational Mapping) Go qui facilite les opérations de base de données en mappant les structs Go aux tables SQL.

**Contexte** : GORM est l'ORM par défaut dans les projets générés par go-starter-kit. Il gère les migrations, queries, et relations.

**Exemple** :
```go
// Initialisation
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

// Migration automatique
db.AutoMigrate(&models.User{}, &models.RefreshToken{})

// Query
var user models.User
db.Where("email = ?", "user@example.com").First(&user)
```

**Voir aussi** : ORM, Migration, PostgreSQL

### Migration

**Définition** : Processus de versioning et application de changements au schéma de base de données de manière contrôlée.

**Contexte** : GORM AutoMigrate est utilisé au démarrage de l'application pour synchroniser le schéma DB avec les models Go.

**Exemple** :
```go
// internal/infrastructure/database.go
func InitDB(cfg *config.Config) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})

    // Migration automatique
    if err := db.AutoMigrate(
        &models.User{},
        &models.RefreshToken{},
    ); err != nil {
        return nil, err
    }

    return db, nil
}
```

**Voir aussi** : GORM, Schema, AutoMigrate

### ORM

**Définition** : Object-Relational Mapping - technique de programmation qui convertit les données entre systèmes de types incompatibles (objets vs tables relationnelles).

**Contexte** : GORM est l'ORM utilisé dans go-starter-kit pour abstraire l'accès SQL et mapper les structs Go aux tables PostgreSQL.

**Voir aussi** : GORM, Repository, PostgreSQL

### PostgreSQL

**Définition** : Système de gestion de base de données relationnelle open-source, robuste et conforme aux standards SQL.

**Contexte** : Base de données par défaut pour les projets générés. Configurée via Docker ou installation locale.

**Configuration** :
```bash
# Via Docker
docker run -d --name postgres \
  -e POSTGRES_DB=myproject \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine

# Connection string dans .env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/myproject?sslmode=disable
```

**Voir aussi** : GORM, Migration, Database URL

### Transaction

**Définition** : Ensemble d'opérations de base de données qui s'exécutent de manière atomique (tout ou rien).

**Contexte** : Utilisé pour garantir la cohérence des données lors d'opérations multiples.

**Exemple** :
```go
func (s *userService) CreateUserWithProfile(ctx context.Context, req CreateUserRequest) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Créer l'utilisateur
        user := &models.User{Email: req.Email}
        if err := tx.Create(user).Error; err != nil {
            return err // Rollback automatique
        }

        // Créer le profil
        profile := &models.Profile{UserID: user.ID}
        if err := tx.Create(profile).Error; err != nil {
            return err // Rollback automatique
        }

        return nil // Commit automatique
    })
}
```

**Voir aussi** : GORM, ACID, Rollback

---

## Sécurité

### Authentication

**Définition** : Processus de vérification de l'identité d'un utilisateur (qui êtes-vous ?).

**Contexte** : Implémenté via JWT dans go-starter-kit. L'utilisateur fournit email/password, reçoit un access token JWT.

**Exemple** :
```go
func (s *authService) Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error) {
    // Vérifier l'email
    user, err := s.userRepo.FindByEmail(ctx, req.Email)

    // Vérifier le mot de passe
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        return nil, errors.New("identifiants invalides")
    }

    // Générer JWT
    token, err := s.tokenService.GenerateToken(user.ID)
    return &models.AuthResponse{AccessToken: token}, nil
}
```

**Voir aussi** : Authorization, JWT, bcrypt

### Authorization

**Définition** : Processus de vérification des permissions d'un utilisateur authentifié (que pouvez-vous faire ?).

**Contexte** : Implémenté via middleware JWT qui vérifie le token et extrait l'user ID pour vérifier les permissions.

**Exemple** :
```go
// Middleware d'autorisation
func AuthMiddleware(tokenService interfaces.TokenService) fiber.Handler {
    return func(c *fiber.Ctx) error {
        token := extractBearerToken(c)
        userID, err := tokenService.ValidateToken(token)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{"error": "non autorisé"})
        }
        c.Locals("userID", userID)
        return c.Next()
    }
}
```

**Voir aussi** : Authentication, Middleware, JWT

### bcrypt

**Définition** : Algorithme de hashing de mots de passe adaptatif et résistant au brute-force grâce à un coût de calcul configurable.

**Contexte** : Utilisé pour hasher les mots de passe avant stockage en base de données.

**Exemple** :
```go
import "golang.org/x/crypto/bcrypt"

// Hasher un mot de passe
hashedPassword, err := bcrypt.GenerateFromPassword(
    []byte(plainPassword),
    bcrypt.DefaultCost, // Coût 10
)

// Vérifier un mot de passe
err := bcrypt.CompareHashAndPassword(
    []byte(hashedPassword),
    []byte(plainPassword),
)
if err != nil {
    // Mot de passe invalide
}
```

**Voir aussi** : Hash, Authentication, Password Security

### Bearer Token

**Définition** : Type de token d'authentification transmis dans le header HTTP `Authorization: Bearer <token>`.

**Contexte** : JWT est transmis comme Bearer token dans les requêtes HTTP pour authentifier l'utilisateur.

**Exemple** :
```http
GET /api/users/me HTTP/1.1
Host: localhost:8080
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

```go
// Extraction du token Bearer
func extractBearerToken(c *fiber.Ctx) string {
    authHeader := c.Get("Authorization")
    if strings.HasPrefix(authHeader, "Bearer ") {
        return strings.TrimPrefix(authHeader, "Bearer ")
    }
    return ""
}
```

**Voir aussi** : JWT, Authorization Header, HTTP

### Hash

**Définition** : Fonction cryptographique unidirectionnelle qui transforme des données en une empreinte de taille fixe.

**Contexte** : Utilisé via bcrypt pour hasher les mots de passe. Les hash sont unidirectionnels (on ne peut pas retrouver le mot de passe original).

**Exemple** :
```go
// Hash de mot de passe (bcrypt)
hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
// Résultat: $2a$10$N9qo8uLOickgx2ZMRZoMye...

// Impossible de retrouver "password123" depuis le hash
// On peut seulement vérifier:
bcrypt.CompareHashAndPassword(hash, []byte("password123")) // nil (match)
bcrypt.CompareHashAndPassword(hash, []byte("wrongpass"))   // error (no match)
```

**Voir aussi** : bcrypt, Encryption, One-way Function

### JWT

**Définition** : JSON Web Token - standard ouvert (RFC 7519) pour créer des tokens d'accès signés qui contiennent des claims JSON.

**Contexte** : Utilisé pour l'authentification stateless. Le token contient l'user ID et est signé avec JWT_SECRET.

**Structure** : `header.payload.signature`

**Exemple** :
```go
import "github.com/golang-jwt/jwt/v5"

// Génération
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

**Voir aussi** : Bearer Token, Authentication, Claims

### Refresh Token

**Définition** : Token de longue durée utilisé pour obtenir de nouveaux access tokens sans re-authentification.

**Contexte** : Stocké en base de données (table refresh_tokens) et associé à un utilisateur. Permet de régénérer des JWT courts.

**Exemple** :
```go
type RefreshToken struct {
    ID        uint      `gorm:"primaryKey"`
    Token     string    `gorm:"uniqueIndex;not null"`
    UserID    uint      `gorm:"not null"`
    ExpiresAt time.Time `gorm:"not null"`
    CreatedAt time.Time
}

// Utilisation
func (s *authService) RefreshAccessToken(ctx context.Context, refreshToken string) (*models.AuthResponse, error) {
    // Vérifier le refresh token
    rt, err := s.tokenRepo.FindByToken(ctx, refreshToken)
    if err != nil || rt.ExpiresAt.Before(time.Now()) {
        return nil, errors.New("refresh token invalide")
    }

    // Générer nouveau access token
    accessToken, err := s.tokenService.GenerateToken(rt.UserID)
    return &models.AuthResponse{AccessToken: accessToken}, nil
}
```

**Voir aussi** : JWT, Token Rotation, Authentication

---

## HTTP/API

### Endpoint

**Définition** : URL spécifique exposée par une API pour effectuer une opération particulière.

**Contexte** : Les endpoints sont définis dans `internal/adapters/http/routes.go` et mappés aux handlers.

**Exemple** :
```go
// Définition des endpoints
func SetupRoutes(app *fiber.App, userHandler *UserHandler, authHandler *AuthHandler) {
    api := app.Group("/api")

    // Endpoints publics
    api.Post("/auth/register", authHandler.Register)
    api.Post("/auth/login", authHandler.Login)

    // Endpoints protégés
    protected := api.Group("", AuthMiddleware(tokenService))
    protected.Get("/users/me", userHandler.GetCurrentUser)
    protected.Put("/users/me", userHandler.UpdateCurrentUser)
}
```

**Voir aussi** : Handler, Router, REST

### Fiber

**Définition** : Framework HTTP web haute performance pour Go, construit sur fasthttp, offrant un routeur rapide et des fonctionnalités middleware avec une API inspirée d'Express.js.

**Contexte** : Framework HTTP par défaut dans go-starter-kit. Utilisé pour définir les routes, handlers, et middlewares.

**Exemple** :
```go
import "github.com/gofiber/fiber/v2"

app := fiber.New()

// Route simple
app.Get("/ping", func(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{"message": "pong"})
})

// Route avec paramètre
app.Get("/users/:id", func(c *fiber.Ctx) error {
    id := c.Params("id")
    return c.JSON(fiber.Map{"user_id": id})
})

app.Listen(":8080")
```

**Voir aussi** : Handler, Router, Middleware

### Handler

**Définition** : Fonction qui traite une requête HTTP pour un endpoint spécifique.

**Contexte** : Les handlers sont définis dans `internal/adapters/http/` et appellent les services du domaine.

**Exemple** :
```go
type UserHandler struct {
    userService interfaces.UserService
    logger      *zap.Logger
}

func (h *UserHandler) GetCurrentUser(c *fiber.Ctx) error {
    userID := c.Locals("userID").(uint)

    user, err := h.userService.GetUserByID(c.Context(), userID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "erreur serveur"})
    }

    return c.JSON(user)
}
```

**Voir aussi** : Endpoint, Fiber, Middleware

### Middleware

**Définition** : Fonction qui s'exécute avant ou après un handler pour effectuer des opérations transversales (auth, logging, CORS, etc.).

**Contexte** : Les middlewares sont utilisés pour l'authentification JWT, le logging, et la gestion CORS.

**Exemple** :
```go
// Middleware d'authentification
func AuthMiddleware(tokenService interfaces.TokenService) fiber.Handler {
    return func(c *fiber.Ctx) error {
        token := extractBearerToken(c)
        userID, err := tokenService.ValidateToken(token)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{"error": "non autorisé"})
        }
        c.Locals("userID", userID)
        return c.Next() // Passe au handler suivant
    }
}

// Utilisation
protected := app.Group("/api", AuthMiddleware(tokenService))
protected.Get("/users/me", userHandler.GetCurrentUser)
```

**Voir aussi** : Handler, Fiber, Authorization

### REST

**Définition** : Representational State Transfer - style architectural pour les APIs web utilisant HTTP et ses méthodes (GET, POST, PUT, DELETE).

**Contexte** : Les APIs générées suivent les conventions REST : ressources comme URLs, méthodes HTTP pour CRUD, stateless.

**Conventions** :
```
GET    /api/users       # Liste des utilisateurs
GET    /api/users/:id   # Détails d'un utilisateur
POST   /api/users       # Créer un utilisateur
PUT    /api/users/:id   # Mettre à jour un utilisateur
DELETE /api/users/:id   # Supprimer un utilisateur
```

**Voir aussi** : Endpoint, HTTP Methods, CRUD

### Router

**Définition** : Composant qui mappe les URLs et méthodes HTTP aux handlers appropriés.

**Contexte** : Fiber Router est configuré dans `internal/adapters/http/routes.go`.

**Exemple** :
```go
func SetupRoutes(app *fiber.App, handlers *Handlers) {
    // Groupe d'API
    api := app.Group("/api")

    // Sous-groupe authentification
    auth := api.Group("/auth")
    auth.Post("/register", handlers.Auth.Register)
    auth.Post("/login", handlers.Auth.Login)

    // Sous-groupe utilisateurs (protégé)
    users := api.Group("/users", AuthMiddleware(handlers.TokenService))
    users.Get("/me", handlers.User.GetCurrentUser)
    users.Put("/me", handlers.User.UpdateCurrentUser)
}
```

**Voir aussi** : Fiber, Endpoint, Handler

---

## Voir Aussi

<i class="material-icons info">info</i> **Documentation complémentaire** :
- [Guide du projet généré](/go-starter-kit/generated-project-guide/) - Guide complet pour les projets générés
- [Architecture CLI](/go-starter-kit/cli-architecture/) - Architecture du générateur de projet
- [Guide d'utilisation](/go-starter-kit/usage/) - Utilisation de create-go-starter

<i class="material-icons">menu_book</i> **Ressources externes** :
- [Documentation Go](https://go.dev/doc/) - Documentation officielle Go
- [Documentation Fiber](https://docs.gofiber.io/) - Framework HTTP
- [Documentation GORM](https://gorm.io/docs/) - ORM pour Go
- [JWT RFC 7519](https://datatracker.ietf.org/doc/html/rfc7519) - Spécification JWT
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/) - Article original par Alistair Cockburn
