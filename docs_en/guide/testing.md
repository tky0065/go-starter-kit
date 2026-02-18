# Testing

<div class="navigation">
  <a href="index.md"><i class="material-icons">arrow_back</i> Guide Index</a>
</div>

---

## Test Organization

Tests are co-located with the source code:

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

### Test Types

#### 1. Unit Tests

Test an isolated function or method, using mocks.

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

### Best Practices for Testing

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


---

## Navigation

**Previous**: [API Reference](api-reference.md)  
**Next**: [Database](database.md)  
**Index**: [Guide Index](index.md)
