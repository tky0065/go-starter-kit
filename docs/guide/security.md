# Sécurité

<div class="navigation">
  <a href="index.md"><i class="material-icons">arrow_back</i> Guide Index</a>
</div>

---

## Authentification JWT

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

  // :material-check-circle: GOOD
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


---

## Navigation

**Previous**: [Base de données](database.md)  
**Next**: [Déploiement](deployment.md)  
**Index**: [Guide Index](index.md)
