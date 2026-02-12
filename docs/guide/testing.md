# Tests

<div class="navigation">
  <a href="index.md"><i class="material-icons">arrow_back</i> Guide Index</a>
</div>

---

## Organisation des tests

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


---

## Navigation

**Previous**: [API Reference](api-reference.md)  
**Next**: [Base de données](database.md)  
**Index**: [Guide Index](index.md)
