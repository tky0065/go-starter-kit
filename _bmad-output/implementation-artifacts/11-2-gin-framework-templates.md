# Story 11.2: Gin Framework Templates

Status: ready-for-dev

## Story

As a développeur,
I want générer un projet avec Gin framework,
So that je puisse utiliser Gin au lieu de Fiber pour mon API.

## Acceptance Criteria

**Given** `--framework=gin` est spécifié
**When** le projet est généré
**Then** les handlers utilisent la syntaxe Gin
**And** le middleware est adapté pour Gin
**And** le projet compile et démarre avec Gin

## Tasks / Subtasks

- [ ] Task 1: Create Gin-specific templates infrastructure (AC: #1, #2, #3)
  - [ ] Create templates_gin.go file with Gin template functions
  - [ ] Implement GinProjectTemplates struct similar to ProjectTemplates
  - [ ] Add Gin version constants (gin-gonic/gin v1.10.0 - latest stable)
  - [ ] Create factory function NewGinProjectTemplates(projectName, database string)

- [ ] Task 2: Implement Gin server and main templates (AC: #1, #3)
  - [ ] Create GinServerTemplate() - Gin router initialization with fx
  - [ ] Create GinMainGoTemplate() - Main entry point using Gin
  - [ ] Implement Gin middleware registration pattern
  - [ ] Add Gin-specific graceful shutdown logic

- [ ] Task 3: Implement Gin handler templates (AC: #1)
  - [ ] Create GinAuthHandlerTemplate() - Login, Register, Refresh with gin.Context
  - [ ] Create GinUserHandlerTemplate() - User CRUD with gin.Context
  - [ ] Create GinHealthHandlerTemplate() - Health checks for Gin
  - [ ] Adapt Swagger annotations for Gin handlers

- [ ] Task 4: Implement Gin middleware templates (AC: #2)
  - [ ] Create GinAuthMiddlewareTemplate() - JWT validation for Gin
  - [ ] Create GinErrorMiddlewareTemplate() - Error handling with gin.Context
  - [ ] Create GinCORSMiddlewareTemplate() - CORS configuration for Gin
  - [ ] Create GinLoggerMiddlewareTemplate() - Request logging with zerolog

- [ ] Task 5: Update generator to support Gin (AC: #1, #2, #3)
  - [ ] Update buildFullFileList() to use Gin templates when framework=gin
  - [ ] Update generateFullTemplateFiles() to route to Gin template builder
  - [ ] Implement getGinTemplateFiles() function
  - [ ] Update go.mod template to include gin-gonic/gin dependency

- [ ] Task 6: Create Gin-specific tests and documentation (AC: #3)
  - [ ] Add E2E test for Gin project generation
  - [ ] Add smoke test to verify Gin project compiles and runs
  - [ ] Create docs/frameworks/gin.md with Gin-specific guide
  - [ ] Update README.md with Gin framework example

## Dev Notes

### Gin Framework Analysis

**Gin vs Fiber - Key Differences:**

| Aspect | Fiber (Current) | Gin (New) |
|--------|----------------|-----------|
| Context | `*fiber.Ctx` | `*gin.Context` |
| Router Creation | `fiber.New()` | `gin.Default()` or `gin.New()` |
| Route Registration | `app.Get("/path", handler)` | `router.GET("/path", handler)` |
| Group Routes | `app.Group("/api")` | `router.Group("/api")` |
| Middleware | `app.Use(middleware)` | `router.Use(middleware)` |
| JSON Response | `c.JSON(fiber.Map{...})` | `c.JSON(http.StatusOK, gin.H{...})` |
| Param Binding | `c.BodyParser(&dto)` | `c.ShouldBindJSON(&dto)` |
| Path Params | `c.Params("id")` | `c.Param("id")` |
| Status Codes | `fiber.StatusOK` | `http.StatusOK` |
| Error Handling | `c.Status(500).JSON(...)` | `c.JSON(500, ...)` or `c.AbortWithStatusJSON` |

**Gin Version:** gin-gonic/gin v1.10.0 (latest stable as of 2026)

### Architecture Compliance

**From architecture.md:**
- All business logic MUST remain in `/internal/domain` (framework-agnostic)
- Only `/internal/adapters/handlers` and `/internal/infrastructure/server` depend on framework
- Dependency injection via `fx` MUST work with Gin

**Critical Patterns to Maintain:**
1. **Hexagonal Architecture:** Domain layer MUST NOT import Gin
2. **Route Grouping:** `/api/v1` prefix MUST be preserved
3. **Error Handling:** Unified error response format MUST be maintained
4. **Validation:** go-playground/validator MUST continue to work
5. **Swagger:** swaggo/swag MUST generate correct docs for Gin

### Gin Server Template Implementation

**Example Gin Server Template (templates_gin.go):**

```go
func (t *GinProjectTemplates) GinServerTemplate() string {
	return `package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"` + t.ProjectName + `/pkg/logger"
)

// GinServer wraps the Gin router
type GinServer struct {
	router *gin.Engine
	config *Config
	logger logger.Logger
}

// NewGinServer creates a new Gin server instance with dependency injection
func NewGinServer(
	lc fx.Lifecycle,
	config *Config,
	log logger.Logger,
) *GinServer {
	// Set Gin mode based on environment
	if config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	router := gin.New()

	server := &GinServer{
		router: router,
		config: config,
		logger: log,
	}

	// Register lifecycle hooks
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				addr := fmt.Sprintf(":%d", config.Port)
				log.Info("Starting Gin server on " + addr)

				srv := &http.Server{
					Addr:    addr,
					Handler: router,
				}

				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatal("Server failed to start: " + err.Error())
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Shutting down Gin server gracefully...")
			// Graceful shutdown logic here
			return nil
		},
	})

	return server
}

// GetRouter returns the Gin router instance
func (s *GinServer) GetRouter() *gin.Engine {
	return s.router
}
`
}
```

### Gin Handler Template Example

**Auth Handler with Gin:**

```go
func (t *GinProjectTemplates) GinAuthHandlerTemplate() string {
	return `package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"` + t.ProjectName + `/internal/domain/user"
	"` + t.ProjectName + `/pkg/logger"
)

// AuthHandler handles authentication endpoints for Gin
type AuthHandler struct {
	userService user.Service
	logger      logger.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(
	userService user.Service,
	log logger.Logger,
) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		logger:      log,
	}
}

// Register godoc
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration data"
// @Success 201 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request data",
			"code":    "INVALID_REQUEST",
		})
		return
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
			"code":    "VALIDATION_ERROR",
		})
		return
	}

	// Call service
	authResp, err := h.userService.Register(c.Request.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
			"code":    "REGISTRATION_FAILED",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   authResp,
	})
}
`
}
```

### Gin Middleware Implementation

**JWT Middleware for Gin:**

```go
func (t *GinProjectTemplates) GinAuthMiddlewareTemplate() string {
	return `package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"` + t.ProjectName + `/pkg/auth"
	"` + t.ProjectName + `/pkg/logger"
)

// AuthMiddleware validates JWT tokens for Gin
type AuthMiddleware struct {
	jwtManager auth.JWTManager
	logger     logger.Logger
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(
	jwtManager auth.JWTManager,
	log logger.Logger,
) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
		logger:     log,
	}
}

// Authenticate returns a Gin middleware function that validates JWT
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Missing authorization header",
				"code":    "UNAUTHORIZED",
			})
			return
		}

		// Extract token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Invalid authorization format",
				"code":    "UNAUTHORIZED",
			})
			return
		}

		token := parts[1]
		claims, err := m.jwtManager.ValidateAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Invalid or expired token",
				"code":    "UNAUTHORIZED",
			})
			return
		}

		// Store user ID in context
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
`
}
```

### Generator Integration

**Update generator.go to support Gin:**

```go
func buildFullFileList(projectPath, projectName, database, observabilityLevel, framework string) []FileGenerator {
	var templates interface{}

	// Create appropriate templates based on framework
	switch framework {
	case "gin":
		templates = NewGinProjectTemplates(projectName, database)
	case "fiber":
		templates = NewProjectTemplatesWithDatabase(projectName, database)
	default:
		// Fallback to Fiber
		templates = NewProjectTemplatesWithDatabase(projectName, database)
	}

	// Build file list using framework-specific templates
	files := []FileGenerator{
		{Path: filepath.Join(projectPath, "go.mod"), Content: templates.GoModTemplate()},
		{Path: filepath.Join(projectPath, "cmd/main.go"), Content: templates.MainGoTemplate()},
		// ... rest of files
	}

	return files
}
```

### Testing Requirements

**E2E Test for Gin:**

```go
func TestGenerateGinProject(t *testing.T) {
	projectName := "test-gin-project"
	defer os.RemoveAll(projectName)

	// Generate project with Gin framework
	err := run(projectName, "full", "postgres", "none", "gin")
	if err != nil {
		t.Fatalf("Failed to generate Gin project: %v", err)
	}

	// Verify Gin dependencies in go.mod
	goModContent, err := os.ReadFile(filepath.Join(projectName, "go.mod"))
	if err != nil {
		t.Fatalf("Failed to read go.mod: %v", err)
	}

	if !strings.Contains(string(goModContent), "github.com/gin-gonic/gin") {
		t.Error("go.mod should contain gin-gonic/gin dependency")
	}

	// Verify Gin imports in server
	serverContent, err := os.ReadFile(filepath.Join(projectName, "internal/infrastructure/server/server.go"))
	if err != nil {
		t.Fatalf("Failed to read server.go: %v", err)
	}

	if !strings.Contains(string(serverContent), "github.com/gin-gonic/gin") {
		t.Error("server.go should import gin-gonic/gin")
	}
}
```

**Smoke Test:**

```bash
# In scripts/smoke_test.sh, add Gin validation:
echo "Testing Gin framework project generation..."
./create-go-starter test-gin-project --framework=gin --database=postgres

cd test-gin-project
go mod tidy
go build -o app ./cmd/main.go

# Verify binary was created
if [ ! -f app ]; then
  echo "FAIL: Gin project did not compile"
  exit 1
fi

echo "SUCCESS: Gin project compiles correctly"
```

### File Structure Notes

**New Files to Create:**
1. `cmd/create-go-starter/templates_gin.go` - All Gin templates
2. `docs/frameworks/gin.md` - Gin-specific documentation

**Files to Modify:**
1. `cmd/create-go-starter/generator.go` - Framework routing logic
2. `cmd/create-go-starter/main.go` - Remove "not implemented" error for Gin
3. `cmd/create-go-starter/smoke_test.go` - Add Gin E2E tests
4. `README.md` - Add Gin examples
5. `docs/usage.md` - Document Gin option

**Files NOT to Modify:**
- Domain logic (`internal/domain/**`) - Must remain framework-agnostic
- Business services - No framework dependencies
- Repository interfaces - Framework-agnostic

### Project Context Rules

**From project-context.md:**
- **Dependency Injection:** ALL components MUST be registered with `fx`
- **Error Handling:** Use named errors from `internal/domain/errors.go`
- **Swagger:** Handlers MUST include swag annotations
- **Linting:** MUST pass golangci-lint without warnings

**Gin-Specific Rules to Add:**
1. Use `gin.H{}` for JSON responses (Gin idiom)
2. Use `c.ShouldBindJSON()` for request binding
3. Use `c.Param()` for path parameters
4. Use `c.AbortWithStatusJSON()` for middleware errors
5. Set `gin.SetMode(gin.ReleaseMode)` in production

### Migration from Fiber to Gin

**Handler Method Signature Changes:**

| Fiber | Gin |
|-------|-----|
| `func (h *Handler) Method(c *fiber.Ctx) error` | `func (h *Handler) Method(c *gin.Context)` |
| `return c.JSON(response)` | `c.JSON(200, response)` |
| `c.BodyParser(&dto)` | `c.ShouldBindJSON(&dto)` |

**Router Setup Changes:**

| Fiber | Gin |
|-------|-----|
| `app := fiber.New()` | `router := gin.New()` |
| `app.Get("/path", handler)` | `router.GET("/path", handler)` |
| `api := app.Group("/api")` | `api := router.Group("/api")` |

### References

**Source Documentation:**
- [Source: _bmad-output/planning-artifacts/epics.md#Epic-11-Story-11.2]
- [Source: _bmad-output/planning-artifacts/architecture.md#Multi-Framework-Support]
- [Source: cmd/create-go-starter/templates.go] - Fiber template patterns to adapt
- [Gin Documentation: https://gin-gonic.com/docs/]
- [Gin GitHub: https://github.com/gin-gonic/gin]

**Dependencies:**
- Requires Story 11.1 (Framework Selection Flag) to be completed
- Blocks Story 11.4 (Framework-Agnostic Abstraction)

## Dev Agent Record

### Agent Model Used

_To be filled by dev agent_

### Debug Log References

_To be filled by dev agent_

### Completion Notes List

- [ ] Gin templates created in templates_gin.go
- [ ] Server, handlers, and middleware templates work with Gin
- [ ] Generated Gin project compiles without errors
- [ ] E2E tests pass for Gin project generation
- [ ] Smoke test validates Gin project runs
- [ ] Documentation updated with Gin examples

### File List

_To be filled by dev agent_
