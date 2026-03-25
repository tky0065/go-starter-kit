# Story 11.4: Framework-Agnostic Abstraction

Status: ready-for-dev

## Story

As a développeur,
I want que la logique métier soit indépendante du framework,
So that je puisse migrer facilement entre frameworks.

## Acceptance Criteria

**Given** un projet est généré
**When** j'examine le code
**Then** la couche domain est 100% indépendante du framework
**And** seuls les adapters dépendent du framework
**And** un guide de migration est fourni

## Tasks / Subtasks

- [ ] Task 1: Audit framework dependencies across all templates (AC: #1, #2)
  - [ ] Verify `/internal/domain/**` has ZERO framework imports
  - [ ] Verify `/internal/models/**` has ZERO framework imports
  - [ ] Verify `/internal/interfaces/**` has ZERO framework imports
  - [ ] Identify any framework leakage in business logic

- [ ] Task 2: Create framework-agnostic adapter interfaces (AC: #1, #2)
  - [ ] Define HTTPContext interface abstraction in `/internal/interfaces`
  - [ ] Define HTTPRequest and HTTPResponse abstractions
  - [ ] Create adapter wrappers for Fiber, Gin, Echo contexts
  - [ ] Update handlers to use abstraction layer (optional pattern)

- [ ] Task 3: Ensure service layer independence (AC: #1)
  - [ ] Verify all services accept `context.Context` (not framework context)
  - [ ] Verify services return domain errors (not framework errors)
  - [ ] Verify services use DTOs from `/internal/models` (not framework types)
  - [ ] Add linting rules to prevent framework imports in domain

- [ ] Task 4: Create migration guide documentation (AC: #3)
  - [ ] Create docs/migration/framework-migration-guide.md
  - [ ] Document step-by-step migration process (Fiber → Gin, Fiber → Echo)
  - [ ] Provide code comparison examples for each framework
  - [ ] Document common pitfalls and solutions

- [ ] Task 5: Add architectural validation tests (AC: #1, #2)
  - [ ] Create arch_test.go to validate layer boundaries
  - [ ] Test that domain has no framework imports
  - [ ] Test that interfaces are framework-agnostic
  - [ ] Add CI check for architectural compliance

- [ ] Task 6: Create framework comparison documentation (AC: #3)
  - [ ] Create docs/frameworks/comparison.md
  - [ ] Compare Fiber vs Gin vs Echo (performance, features, ecosystem)
  - [ ] Provide decision matrix for choosing framework
  - [ ] Document trade-offs and recommendations

## Dev Notes

### Hexagonal Architecture Enforcement

**Current Architecture (from architecture.md):**

```
go-starter-kit/
├── internal/
│   ├── domain/              # CORE - MUST be framework-agnostic
│   │   └── user/           # Business logic, services
│   ├── models/             # CORE - Shared entities (User, RefreshToken)
│   ├── interfaces/         # PORTS - Contracts (MUST be framework-agnostic)
│   ├── adapters/           # ADAPTERS - Framework-specific (handlers, middleware)
│   │   ├── handlers/       # ✅ CAN import framework
│   │   ├── middleware/     # ✅ CAN import framework
│   │   └── repository/     # ✅ CAN import GORM
│   └── infrastructure/     # INFRASTRUCTURE
│       ├── database/       # ✅ CAN import GORM
│       └── server/         # ✅ CAN import framework
```

**Dependency Rule:**
```
Domain ← Interfaces ← Adapters ← Infrastructure
   ↑                                 ↓
   └─────────────────────────────────┘
        (Dependency Inversion)
```

**Critical Rules:**
1. **Domain layer** MUST NOT import: Fiber, Gin, Echo, or any HTTP framework
2. **Interfaces layer** MUST define contracts using standard library types only
3. **Adapters layer** CAN import framework-specific packages
4. **Services** MUST accept `context.Context` (not `*fiber.Ctx`, `*gin.Context`, or `echo.Context`)

### Framework Independence Patterns

**❌ BAD - Framework Leakage into Domain:**

```go
// internal/domain/user/service.go
package user

import (
	"github.com/gofiber/fiber/v2"  // ❌ FORBIDDEN
)

type Service interface {
	Register(c *fiber.Ctx) error  // ❌ Framework dependency in domain
}
```

**✅ GOOD - Framework-Agnostic Domain:**

```go
// internal/domain/user/service.go
package user

import (
	"context"  // ✅ Standard library only
)

type Service interface {
	Register(ctx context.Context, email, password, name string) (*AuthResponse, error)  // ✅ Pure Go types
}
```

**Handler Adaptation (Fiber Example):**

```go
// internal/adapters/handlers/auth_handler.go (Fiber)
package handlers

import (
	"github.com/gofiber/fiber/v2"  // ✅ Framework imports allowed in adapters
	"your-project/internal/domain/user"
)

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(...)
	}

	// Call framework-agnostic service
	authResp, err := h.userService.Register(c.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		return c.Status(500).JSON(...)
	}

	return c.JSON(authResp)
}
```

**Handler Adaptation (Gin Example):**

```go
// internal/adapters/handlers/auth_handler.go (Gin)
package handlers

import (
	"github.com/gin-gonic/gin"  // ✅ Framework imports allowed in adapters
	"your-project/internal/domain/user"
)

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, ...)
		return
	}

	// Same framework-agnostic service call
	authResp, err := h.userService.Register(c.Request.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		c.JSON(500, ...)
		return
	}

	c.JSON(200, authResp)
}
```

**Handler Adaptation (Echo Example):**

```go
// internal/adapters/handlers/auth_handler.go (Echo)
package handlers

import (
	"github.com/labstack/echo/v4"  // ✅ Framework imports allowed in adapters
	"your-project/internal/domain/user"
)

func (h *AuthHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, ...)
	}

	// Same framework-agnostic service call
	authResp, err := h.userService.Register(c.Request().Context(), req.Email, req.Password, req.Name)
	if err != nil {
		return c.JSON(500, ...)
	}

	return c.JSON(200, authResp)
}
```

### Service Layer Best Practices

**Service Interface Design:**

```go
// internal/interfaces/user_service.go
package interfaces

import (
	"context"
	"your-project/internal/models"
)

type UserService interface {
	// ✅ Uses standard library context
	// ✅ Returns domain models
	// ✅ No framework dependencies
	Register(ctx context.Context, email, password, name string) (*models.AuthResponse, error)
	Login(ctx context.Context, email, password string) (*models.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResponse, error)
	GetUserByID(ctx context.Context, userID uint) (*models.User, error)
}
```

**Service Implementation:**

```go
// internal/domain/user/service.go
package user

import (
	"context"
	"errors"

	"your-project/internal/interfaces"
	"your-project/internal/models"
	"your-project/pkg/auth"
)

type service struct {
	repo       interfaces.UserRepository
	jwtManager auth.JWTManager
}

func NewService(repo interfaces.UserRepository, jwtManager auth.JWTManager) interfaces.UserService {
	return &service{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

func (s *service) Register(ctx context.Context, email, password, name string) (*models.AuthResponse, error) {
	// ✅ Pure business logic, no framework dependencies
	// Validation, hashing, database operations, JWT generation
	// Returns domain models only
}
```

### Architectural Validation Tests

**Create arch_test.go:**

```go
// cmd/create-go-starter/arch_test.go
package main

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestDomainLayerIndependence verifies that the domain layer has NO framework imports
func TestDomainLayerIndependence(t *testing.T) {
	domainPath := "../../templates/internal/domain"

	forbiddenImports := []string{
		"github.com/gofiber/fiber",
		"github.com/gin-gonic/gin",
		"github.com/labstack/echo",
	}

	err := filepath.Walk(domainPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		for _, forbidden := range forbiddenImports {
			if strings.Contains(string(content), forbidden) {
				t.Errorf("Domain layer file %s contains forbidden import: %s", path, forbidden)
			}
		}

		return nil
	})

	if err != nil {
		t.Fatalf("Failed to walk domain directory: %v", err)
	}
}

// TestInterfacesLayerIndependence verifies that interfaces are framework-agnostic
func TestInterfacesLayerIndependence(t *testing.T) {
	interfacesPath := "../../templates/internal/interfaces"

	forbiddenImports := []string{
		"github.com/gofiber/fiber",
		"github.com/gin-gonic/gin",
		"github.com/labstack/echo",
	}

	err := filepath.Walk(interfacesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		for _, forbidden := range forbiddenImports {
			if strings.Contains(string(content), forbidden) {
				t.Errorf("Interfaces layer file %s contains forbidden import: %s", path, forbidden)
			}
		}

		return nil
	})

	if err != nil {
		t.Fatalf("Failed to walk interfaces directory: %v", err)
	}
}

// TestServiceSignatures verifies that service methods use context.Context
func TestServiceSignatures(t *testing.T) {
	interfacesPath := "../../templates/internal/interfaces"

	err := filepath.Walk(interfacesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, "_service.go") {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		// Verify interface methods accept context.Context as first parameter
		// (Implementation would use AST traversal to check method signatures)

		return nil
	})

	if err != nil {
		t.Fatalf("Failed to parse service interfaces: %v", err)
	}
}
```

### Migration Guide Structure

**docs/migration/framework-migration-guide.md:**

```markdown
# Framework Migration Guide

## Overview

This guide explains how to migrate a go-starter-kit project from one framework to another.

## Migration Paths

### Fiber → Gin
### Fiber → Echo
### Gin → Echo
### Gin → Fiber
### Echo → Fiber
### Echo → Gin

## Step-by-Step Migration Process

### 1. Update Dependencies

### 2. Modify Server Initialization

### 3. Update Handlers

### 4. Update Middleware

### 5. Update Tests

### 6. Verify Migration

## Code Comparison Examples

### Handler Method Comparison

### Middleware Comparison

### Router Setup Comparison

## Common Pitfalls

## Testing Your Migration

## Rollback Strategy
```

### Framework Comparison Matrix

**docs/frameworks/comparison.md:**

| Feature | Fiber | Gin | Echo |
|---------|-------|-----|------|
| **Performance** | Fastest | Fast | Fast |
| **Memory Usage** | Low | Medium | Low |
| **Middleware Ecosystem** | Good | Excellent | Excellent |
| **Learning Curve** | Easy | Easy | Medium |
| **Express.js-like API** | ✅ Yes | Partial | No |
| **Built-in Validation** | No | Yes (binding) | Yes (binding) |
| **Context Type** | Struct | Struct | Interface |
| **Error Handling** | Return error | Void/abort | Return error |
| **Community Size** | Growing | Large | Large |
| **Maturity** | Newer | Mature | Very Mature |
| **Best For** | High-performance APIs | General-purpose | Minimalist APIs |

**When to Choose:**

- **Fiber:** Maximum performance, Express.js familiarity, modern features
- **Gin:** Mature ecosystem, extensive middleware, production-proven
- **Echo:** Minimalist design, testable (interface context), built-in features

### Linting Rules for Architecture Enforcement

**Add to .golangci.yml (generated project):**

```yaml
linters-settings:
  depguard:
    rules:
      main:
        deny:
          - pkg: "github.com/gofiber/fiber"
            desc: "Fiber should only be used in adapters and infrastructure layers"
          - pkg: "github.com/gin-gonic/gin"
            desc: "Gin should only be used in adapters and infrastructure layers"
          - pkg: "github.com/labstack/echo"
            desc: "Echo should only be used in adapters and infrastructure layers"
        files:
          - "**/internal/domain/**/*.go"
          - "**/internal/interfaces/**/*.go"
          - "**/internal/models/**/*.go"
```

### Testing Requirements

**Architectural Compliance Tests:**

```go
func TestArchitecturalBoundaries(t *testing.T) {
	tests := []struct {
		name         string
		layer        string
		allowedDeps  []string
		forbiddenDeps []string
	}{
		{
			name:         "Domain layer has no framework deps",
			layer:        "internal/domain",
			allowedDeps:  []string{"context", "errors", "time"},
			forbiddenDeps: []string{"fiber", "gin", "echo"},
		},
		{
			name:         "Interfaces layer is framework-agnostic",
			layer:        "internal/interfaces",
			allowedDeps:  []string{"context"},
			forbiddenDeps: []string{"fiber", "gin", "echo"},
		},
		{
			name:         "Models layer has no framework deps",
			layer:        "internal/models",
			allowedDeps:  []string{"gorm", "time"},
			forbiddenDeps: []string{"fiber", "gin", "echo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify layer dependencies
		})
	}
}
```

### Project Context Rules

**From project-context.md:**
- **Hexagonal Architecture:** Domain MUST NOT import framework packages
- **Interfaces:** MUST be defined in `/internal/interfaces`
- **Services:** MUST accept `context.Context` (not framework context)
- **Error Handling:** Use named errors from `internal/domain/errors.go`

**New Rules to Add:**
1. **Domain Purity:** `/internal/domain` MUST pass arch_test.go
2. **Interface Contracts:** All service interfaces MUST use standard library types
3. **Adapter Isolation:** Only adapters MAY import framework packages
4. **Migration Path:** All frameworks MUST support same business logic

### References

**Source Documentation:**
- [Source: _bmad-output/planning-artifacts/epics.md#Epic-11-Story-11.4]
- [Source: _bmad-output/planning-artifacts/architecture.md#Hexagonal-Architecture]
- [Source: _bmad-output/project-context.md#Framework-Specific-Rules]
- [Hexagonal Architecture Pattern](https://alistair.cockburn.us/hexagonal-architecture/)

**Dependencies:**
- Requires Story 11.1 (Framework Selection Flag) to be completed
- Requires Story 11.2 (Gin Templates) to be completed
- Requires Story 11.3 (Echo Templates) to be completed

**Related Documentation:**
- Clean Architecture by Robert C. Martin
- Hexagonal Architecture by Alistair Cockburn
- Domain-Driven Design by Eric Evans

## Dev Agent Record

### Agent Model Used

_To be filled by dev agent_

### Debug Log References

_To be filled by dev agent_

### Completion Notes List

- [ ] Domain layer verified 100% framework-independent
- [ ] Interfaces layer verified framework-agnostic
- [ ] Service signatures use context.Context exclusively
- [ ] Architectural validation tests pass (arch_test.go)
- [ ] Migration guide completed with examples for all frameworks
- [ ] Framework comparison documentation created
- [ ] Linting rules added to prevent framework leakage
- [ ] CI checks enforce architectural boundaries

### File List

_To be filled by dev agent_
