# Security

<div class="navigation">
  <a href="index.md"><i class="material-icons">arrow_back</i> Guide Index</a>
</div>

---

## JWT Authentication

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
- **Not localStorage** (vulnerable to XSS)
- Lost on page refresh → Use refresh token

**Refresh Token** (long-lived: 7d):
- **Option 1**: httpOnly cookie (most secure)
- **Option 2**: localStorage (if no XSS risk)

**React Example**:

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
| `alphanum` | Letters + numbers |
| `numeric` | Numbers only |
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
- [ ] **Request timeout**: Prevent DoS
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

**Previous**: [Database](database.md)  
**Next**: [Deployment](deployment.md)  
**Index**: [Guide Index](index.md)
