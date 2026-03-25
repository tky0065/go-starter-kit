# Story 11.3: Echo Framework Templates

Status: ready-for-dev

## Story

As a développeur,
I want générer un projet avec Echo framework,
So that je puisse utiliser Echo au lieu de Fiber pour mon API.

## Acceptance Criteria

**Given** `--framework=echo` est spécifié
**When** le projet est généré
**Then** les handlers utilisent Echo
**And** le middleware est adapté
**And** le projet compile et démarre avec Echo

## Tasks / Subtasks

- [ ] Task 1: Create Echo-specific templates infrastructure (AC: #1, #2, #3)
  - [ ] Create templates_echo.go file with Echo template functions
  - [ ] Implement EchoProjectTemplates struct similar to GinProjectTemplates
  - [ ] Add Echo version constants (labstack/echo/v4 v4.12.0 - latest stable)
  - [ ] Create factory function NewEchoProjectTemplates(projectName, database string)

- [ ] Task 2: Implement Echo server and main templates (AC: #1, #3)
  - [ ] Create EchoServerTemplate() - Echo router initialization with fx
  - [ ] Create EchoMainGoTemplate() - Main entry point using Echo
  - [ ] Implement Echo middleware registration pattern
  - [ ] Add Echo-specific graceful shutdown logic

- [ ] Task 3: Implement Echo handler templates (AC: #1)
  - [ ] Create EchoAuthHandlerTemplate() - Login, Register, Refresh with echo.Context
  - [ ] Create EchoUserHandlerTemplate() - User CRUD with echo.Context
  - [ ] Create EchoHealthHandlerTemplate() - Health checks for Echo
  - [ ] Adapt Swagger annotations for Echo handlers

- [ ] Task 4: Implement Echo middleware templates (AC: #2)
  - [ ] Create EchoAuthMiddlewareTemplate() - JWT validation for Echo
  - [ ] Create EchoErrorMiddlewareTemplate() - Error handling with echo.Context
  - [ ] Create EchoCORSMiddlewareTemplate() - CORS configuration (use echo/middleware)
  - [ ] Create EchoLoggerMiddlewareTemplate() - Request logging with zerolog

- [ ] Task 5: Update generator to support Echo (AC: #1, #2, #3)
  - [ ] Update buildFullFileList() to use Echo templates when framework=echo
  - [ ] Update generateFullTemplateFiles() to route to Echo template builder
  - [ ] Implement getEchoTemplateFiles() function
  - [ ] Update go.mod template to include labstack/echo dependency

- [ ] Task 6: Create Echo-specific tests and documentation (AC: #3)
  - [ ] Add E2E test for Echo project generation
  - [ ] Add smoke test to verify Echo project compiles and runs
  - [ ] Create docs/frameworks/echo.md with Echo-specific guide
  - [ ] Update README.md with Echo framework example

## Dev Notes

### Echo Framework Analysis

**Echo vs Fiber vs Gin - Key Differences:**

| Aspect | Fiber (Current) | Gin | Echo (New) |
|--------|----------------|-----|------------|
| Context | `*fiber.Ctx` | `*gin.Context` | `echo.Context` (interface) |
| Router Creation | `fiber.New()` | `gin.Default()` | `echo.New()` |
| Route Registration | `app.Get("/path", handler)` | `router.GET("/path", handler)` | `e.GET("/path", handler)` |
| Group Routes | `app.Group("/api")` | `router.Group("/api")` | `e.Group("/api")` |
| Middleware | `app.Use(middleware)` | `router.Use(middleware)` | `e.Use(middleware)` |
| JSON Response | `c.JSON(fiber.Map{...})` | `c.JSON(200, gin.H{...})` | `c.JSON(200, map[string]interface{}{...})` |
| Param Binding | `c.BodyParser(&dto)` | `c.ShouldBindJSON(&dto)` | `c.Bind(&dto)` |
| Path Params | `c.Params("id")` | `c.Param("id")` | `c.Param("id")` |
| Status Codes | `fiber.StatusOK` | `http.StatusOK` | `http.StatusOK` |
| Error Handling | `c.Status(500).JSON(...)` | `c.JSON(500, ...)` | `return c.JSON(500, ...)` or `echo.NewHTTPError` |
| Return Value | `return nil` | void | `return nil` or `return err` |

**Echo Version:** labstack/echo/v4 v4.12.0 (latest stable as of 2026)

**Echo Philosophy:**
- Minimalist and performance-focused
- Extensive built-in middleware (CORS, JWT, recover, logger)
- Context is an interface (highly testable)
- Handlers return errors (better error propagation)

### Architecture Compliance

**From architecture.md:**
- All business logic MUST remain in `/internal/domain` (framework-agnostic)
- Only `/internal/adapters/handlers` and `/internal/infrastructure/server` depend on framework
- Dependency injection via `fx` MUST work with Echo

**Critical Patterns to Maintain:**
1. **Hexagonal Architecture:** Domain layer MUST NOT import Echo
2. **Route Grouping:** `/api/v1` prefix MUST be preserved
3. **Error Handling:** Unified error response format MUST be maintained
4. **Validation:** go-playground/validator MUST continue to work
5. **Swagger:** swaggo/swag MUST generate correct docs for Echo

### Echo Server Template Implementation

**Example Echo Server Template (templates_echo.go):**

```go
func (t *EchoProjectTemplates) EchoServerTemplate() string {
	return `package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
	"` + t.ProjectName + `/pkg/logger"
)

// EchoServer wraps the Echo router
type EchoServer struct {
	echo   *echo.Echo
	config *Config
	logger logger.Logger
}

// NewEchoServer creates a new Echo server instance with dependency injection
func NewEchoServer(
	lc fx.Lifecycle,
	config *Config,
	log logger.Logger,
) *EchoServer {
	// Create Echo instance
	e := echo.New()

	// Hide Echo banner (we have our own logging)
	e.HideBanner = true
	e.HidePort = true

	// Built-in middleware
	e.Use(middleware.Recover())

	server := &EchoServer{
		echo:   e,
		config: config,
		logger: log,
	}

	// Register lifecycle hooks
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				addr := fmt.Sprintf(":%d", config.Port)
				log.Info("Starting Echo server on " + addr)

				if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
					log.Fatal("Server failed to start: " + err.Error())
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Shutting down Echo server gracefully...")
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			return e.Shutdown(shutdownCtx)
		},
	})

	return server
}

// GetEcho returns the Echo instance
func (s *EchoServer) GetEcho() *echo.Echo {
	return s.echo
}
`
}
```

### Echo Handler Template Example

**Auth Handler with Echo:**

```go
func (t *EchoProjectTemplates) EchoAuthHandlerTemplate() string {
	return `package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"` + t.ProjectName + `/internal/domain/user"
	"` + t.ProjectName + `/pkg/logger"
)

// AuthHandler handles authentication endpoints for Echo
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
func (h *AuthHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Invalid request data",
			"code":    "INVALID_REQUEST",
		})
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
			"code":    "VALIDATION_ERROR",
		})
	}

	// Call service
	authResp, err := h.userService.Register(c.Request().Context(), req.Email, req.Password, req.Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
			"code":    "REGISTRATION_FAILED",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status": "success",
		"data":   authResp,
	})
}
`
}
```

### Echo Middleware Implementation

**JWT Middleware for Echo:**

```go
func (t *EchoProjectTemplates) EchoAuthMiddlewareTemplate() string {
	return `package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"` + t.ProjectName + `/pkg/auth"
	"` + t.ProjectName + `/pkg/logger"
)

// AuthMiddleware validates JWT tokens for Echo
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

// Authenticate returns an Echo middleware function that validates JWT
func (m *AuthMiddleware) Authenticate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"status":  "error",
					"message": "Missing authorization header",
					"code":    "UNAUTHORIZED",
				})
			}

			// Extract token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"status":  "error",
					"message": "Invalid authorization format",
					"code":    "UNAUTHORIZED",
				})
			}

			token := parts[1]
			claims, err := m.jwtManager.ValidateAccessToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"status":  "error",
					"message": "Invalid or expired token",
					"code":    "UNAUTHORIZED",
				})
			}

			// Store user ID in context
			c.Set("userID", claims.UserID)
			return next(c)
		}
	}
}
`
}
```

**Error Handling Middleware for Echo:**

```go
func (t *EchoProjectTemplates) EchoErrorMiddlewareTemplate() string {
	return `package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"` + t.ProjectName + `/pkg/logger"
)

// ErrorHandler returns a custom error handler for Echo
func ErrorHandler(log logger.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		message := "Internal server error"
		errorCode := "INTERNAL_ERROR"

		// Check if it's an Echo HTTP error
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = he.Message.(string)
			errorCode = "HTTP_ERROR"
		}

		// Log the error
		log.Error("Request error: " + err.Error())

		// Send JSON error response
		if !c.Response().Committed {
			c.JSON(code, map[string]interface{}{
				"status":  "error",
				"message": message,
				"code":    errorCode,
			})
		}
	}
}
`
}
```

### Echo-Specific Features

**Built-in Middleware Available:**
- `middleware.CORS()` - Cross-Origin Resource Sharing
- `middleware.JWT()` - JWT authentication (alternative to custom)
- `middleware.Logger()` - Request logging
- `middleware.Recover()` - Panic recovery
- `middleware.RequestID()` - Request ID generation
- `middleware.Secure()` - Security headers

**Custom Error Handler:**
Echo allows setting a custom error handler for centralized error handling:

```go
e.HTTPErrorHandler = ErrorHandler(log)
```

### Generator Integration

**Update generator.go to support Echo:**

```go
func buildFullFileList(projectPath, projectName, database, observabilityLevel, framework string) []FileGenerator {
	var templates interface{}

	// Create appropriate templates based on framework
	switch framework {
	case "echo":
		templates = NewEchoProjectTemplates(projectName, database)
	case "gin":
		templates = NewGinProjectTemplates(projectName, database)
	case "fiber":
		fallthrough
	default:
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

**E2E Test for Echo:**

```go
func TestGenerateEchoProject(t *testing.T) {
	projectName := "test-echo-project"
	defer os.RemoveAll(projectName)

	// Generate project with Echo framework
	err := run(projectName, "full", "postgres", "none", "echo")
	if err != nil {
		t.Fatalf("Failed to generate Echo project: %v", err)
	}

	// Verify Echo dependencies in go.mod
	goModContent, err := os.ReadFile(filepath.Join(projectName, "go.mod"))
	if err != nil {
		t.Fatalf("Failed to read go.mod: %v", err)
	}

	if !strings.Contains(string(goModContent), "github.com/labstack/echo/v4") {
		t.Error("go.mod should contain labstack/echo/v4 dependency")
	}

	// Verify Echo imports in server
	serverContent, err := os.ReadFile(filepath.Join(projectName, "internal/infrastructure/server/server.go"))
	if err != nil {
		t.Fatalf("Failed to read server.go: %v", err)
	}

	if !strings.Contains(string(serverContent), "github.com/labstack/echo/v4") {
		t.Error("server.go should import labstack/echo/v4")
	}

	// Verify handler return errors (Echo pattern)
	authHandlerContent, err := os.ReadFile(filepath.Join(projectName, "internal/adapters/handlers/auth_handler.go"))
	if err != nil {
		t.Fatalf("Failed to read auth_handler.go: %v", err)
	}

	if !strings.Contains(string(authHandlerContent), "return c.JSON") {
		t.Error("Echo handlers should return c.JSON() results")
	}
}
```

**Smoke Test:**

```bash
# In scripts/smoke_test.sh, add Echo validation:
echo "Testing Echo framework project generation..."
./create-go-starter test-echo-project --framework=echo --database=postgres

cd test-echo-project
go mod tidy
go build -o app ./cmd/main.go

# Verify binary was created
if [ ! -f app ]; then
  echo "FAIL: Echo project did not compile"
  exit 1
fi

# Verify Echo-specific middleware
if ! grep -q "middleware.Recover()" internal/infrastructure/server/server.go; then
  echo "FAIL: Echo server should use middleware.Recover()"
  exit 1
fi

echo "SUCCESS: Echo project compiles correctly"
```

### File Structure Notes

**New Files to Create:**
1. `cmd/create-go-starter/templates_echo.go` - All Echo templates
2. `docs/frameworks/echo.md` - Echo-specific documentation

**Files to Modify:**
1. `cmd/create-go-starter/generator.go` - Framework routing logic (add Echo case)
2. `cmd/create-go-starter/main.go` - Remove "not implemented" error for Echo
3. `cmd/create-go-starter/smoke_test.go` - Add Echo E2E tests
4. `README.md` - Add Echo examples
5. `docs/usage.md` - Document Echo option

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

**Echo-Specific Rules to Add:**
1. Handlers MUST return errors (`func(c echo.Context) error`)
2. Use `c.Bind()` for request binding
3. Use `c.Request().Context()` to get request context
4. Use `echo.NewHTTPError()` for HTTP errors
5. Set custom error handler with `e.HTTPErrorHandler`
6. Use built-in `middleware.CORS()` and `middleware.Recover()`

### Migration from Fiber to Echo

**Handler Method Signature Changes:**

| Fiber | Echo |
|-------|------|
| `func (h *Handler) Method(c *fiber.Ctx) error` | `func (h *Handler) Method(c echo.Context) error` |
| `return c.JSON(response)` | `return c.JSON(200, response)` |
| `c.BodyParser(&dto)` | `c.Bind(&dto)` |
| `c.Context()` | `c.Request().Context()` |

**Router Setup Changes:**

| Fiber | Echo |
|-------|------|
| `app := fiber.New()` | `e := echo.New()` |
| `app.Get("/path", handler)` | `e.GET("/path", handler)` |
| `api := app.Group("/api")` | `api := e.Group("/api")` |
| `app.Use(middleware)` | `e.Use(middleware)` |

**Middleware Pattern:**

| Fiber | Echo |
|-------|------|
| `func middleware(c *fiber.Ctx) error` | `func(next echo.HandlerFunc) echo.HandlerFunc` |
| `return c.Next()` | `return next(c)` |

### References

**Source Documentation:**
- [Source: _bmad-output/planning-artifacts/epics.md#Epic-11-Story-11.3]
- [Source: _bmad-output/planning-artifacts/architecture.md#Multi-Framework-Support]
- [Source: cmd/create-go-starter/templates_gin.go] - Gin template patterns to adapt
- [Echo Documentation: https://echo.labstack.com/]
- [Echo GitHub: https://github.com/labstack/echo]

**Dependencies:**
- Requires Story 11.1 (Framework Selection Flag) to be completed
- Can be implemented in parallel with Story 11.2 (Gin templates)
- Blocks Story 11.4 (Framework-Agnostic Abstraction)

## Dev Agent Record

### Agent Model Used

_To be filled by dev agent_

### Debug Log References

_To be filled by dev agent_

### Completion Notes List

- [ ] Echo templates created in templates_echo.go
- [ ] Server, handlers, and middleware templates work with Echo
- [ ] Generated Echo project compiles without errors
- [ ] E2E tests pass for Echo project generation
- [ ] Smoke test validates Echo project runs
- [ ] Documentation updated with Echo examples
- [ ] Echo error handler properly set up

### File List

_To be filled by dev agent_
